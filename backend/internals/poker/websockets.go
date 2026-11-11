package poker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// readPump reads messages from one connection until it fails. It is the only
// goroutine that closes cu.Send, which in turn stops writePump.
func (r *Room) readPump(cu *ConnectedUser) {
	defer func() {
		r.removeConn(cu)
		cu.Conn.Close()
		close(cu.Send)
	}()

	cu.Conn.SetReadLimit(maxMessageSize)
	cu.Conn.SetReadDeadline(time.Now().Add(pongWait))
	cu.Conn.SetPongHandler(func(string) error {
		cu.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// ReadMessage blocks, so there is no point checking for shutdown here:
	// Room.Close closes the connection, which makes ReadMessage return.
	for {
		_, message, err := cu.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			return
		}

		if err := r.handleMessage(cu, message); err != nil {
			log.Printf("Error handling message: %v", err)
			// Queue the reply rather than writing it here, so writePump stays
			// the connection's only writer.
			reply, _ := json.Marshal(map[string]string{"error": err.Error()})
			select {
			case cu.Send <- reply:
			default:
				log.Printf("Send buffer full for user: %s, dropping error reply", cu.User.Username)
			}
		}
	}
}

// writePump is the only goroutine that writes to the connection, so writes
// and write deadlines never race.
func (r *Room) writePump(cu *ConnectedUser) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		cu.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-cu.Send:
			cu.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				cu.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := cu.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			cu.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := cu.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (r *Room) handleMessage(cu *ConnectedUser, message []byte) error {
	var wsMessage WSMessage
	if err := json.Unmarshal(message, &wsMessage); err != nil {
		return fmt.Errorf("invalid message format: %v", err)
	}

	switch wsMessage.Type {
	case MessageTypeVote:
		var payload VotePayload
		if err := json.Unmarshal(wsMessage.Payload, &payload); err != nil {
			return fmt.Errorf("invalid vote payload: %v", err)
		}
		return r.SubmitVote(cu.User.Username, payload.Vote)

	case MessageTypeStartVoting:
		if cu.User.Username != r.Owner {
			return errors.New("only room owner can start voting")
		}
		r.StartVoting()
		return nil

	case MessageTypeRevealVotes:
		// Only room owner can reveal votes
		if cu.User.Username != r.Owner {
			return errors.New("only room owner can reveal votes")
		}
		return r.RevealVotes(cu.User.Username)

	case MessageTypeRestartVoting:
		// Only room owner can reset voting
		if cu.User.Username != r.Owner {
			return errors.New("only room owner can reset voting")
		}
		r.RestartVoting()
		return nil

	default:
		return fmt.Errorf("unknown message type: %s", wsMessage.Type)
	}
}

// scheduleBroadcast sends the room state to everyone shortly, coalescing
// bursts of changes into one broadcast. It is safe to call while holding
// r.Mutex: the timer callback releases broadcastMutex before it takes r.Mutex,
// so the two locks are never held in opposite orders.
func (r *Room) scheduleBroadcast() {
	r.broadcastMutex.Lock()
	defer r.broadcastMutex.Unlock()

	if r.broadcastPending {
		return
	}
	r.broadcastPending = true

	time.AfterFunc(broadcastDelay, func() {
		r.broadcastMutex.Lock()
		r.broadcastPending = false
		r.broadcastMutex.Unlock()

		r.broadcastState()
	})
}

func (r *Room) broadcastState() {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	if r.closed {
		return
	}

	// Each user gets their own payload so they can see their own vote
	// before reveal without exposing anyone else's.
	for username, user := range r.Users {
		message, err := json.Marshal(map[string]interface{}{
			"payload": r.serializeForLocked(username),
		})
		if err != nil {
			log.Printf("Error marshaling state for %s: %v", username, err)
			continue
		}

		select {
		case user.Send <- message:
		default:
			log.Printf("Send buffer full for user: %s, may be disconnected", username)
		}
	}
}

// Close disconnects everyone and stops the room accepting joins.
func (r *Room) Close() {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	r.closed = true
	for _, user := range r.Users {
		user.Conn.Close()
	}
}
