package poker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (r *Room) readPump(cu *ConnectedUser) {
	defer func() {
		r.RemoveUser(cu.User.Username)
		cu.Conn.Close()
		close(cu.Send)
	}()

	cu.Conn.SetReadLimit(maxMessageSize)
	cu.Conn.SetReadDeadline(time.Now().Add(pongWait))
	cu.Conn.SetPongHandler(func(string) error {
		cu.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-r.ctx.Done():
			return // Room is closing
		default:
			_, message, err := cu.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				return
			}

			if err := r.handleMessage(cu, message); err != nil {
				log.Printf("Error handling message: %v", err)
				errorMsg := map[string]string{"error": err.Error()}
				cu.Mutex.Lock()
				jsonErr := cu.Conn.WriteJSON(errorMsg)
				cu.Mutex.Unlock()
				if jsonErr != nil {
					log.Printf("Error sending error message: %v", jsonErr)
				}
			}
		}
	}
}

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

			// Use a mutex to prevent concurrent writes to the connection
			cu.Mutex.Lock()
			err := cu.Conn.WriteMessage(websocket.TextMessage, message)
			cu.Mutex.Unlock()

			if err != nil {
				return
			}
		case <-ticker.C:
			cu.Mutex.Lock()
			cu.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := cu.Conn.WriteMessage(websocket.PingMessage, nil)
			cu.Mutex.Unlock()

			if err != nil {
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

func (r *Room) scheduleBroadcast() {
	r.broadcastMutex.Lock()
	defer r.broadcastMutex.Unlock()

	if r.broadcastTimer != nil {
		r.broadcastTimer.Stop()
	}

	r.broadcastPending = true
	r.broadcastTimer = time.AfterFunc(100*time.Millisecond, func() {
		r.broadcastMutex.Lock()
		defer r.broadcastMutex.Unlock()

		if r.broadcastPending {
			r.broadcastPending = false
			r.broadcastState()
		}
	})
}


func (r *Room) broadcastState() {
	serializable := r.ToSerializable()
	message, err := json.Marshal(map[string]interface{}{
		"payload": serializable,
	})

	if err != nil {
		log.Printf("Error marshaling state: %v", err)
		return
	}

	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	for _, user := range r.Users {
		select {
		case user.Send <- message:
		default:
			log.Printf("Send buffer full for user: %s, may be disconnected", user.User.Username)
		}
	}
}

func (r *Room) Close() {
	r.cancelFunc()

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	for _, user := range r.Users {
		user.Conn.Close()
		close(user.Send)
	}
}
