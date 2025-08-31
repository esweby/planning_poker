package poker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type RoomID = string

type RoomStatus = string

const (
	RoomInUse      = "roomInUse"
	RoomEmpty      = "roomEmpty"
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Room struct {
	ID         RoomID
	Owner      Username
	Users      ConnectedUsers
	Game       Game
	Mutex      sync.RWMutex
	EmptySince time.Time

	// lifecycle methods
	cancelFunc context.CancelFunc
	ctx        context.Context

	// Debouncing messages
	broadcastPending bool
	broadcastTimer   *time.Timer
	broadcastMutex   sync.Mutex
}

type Rooms = map[RoomID]*Room

type SerializableRoom struct {
	ID    RoomID   `json:"id"`
	Owner Username `json:"owner"`
	Users Users    `json:"users"`
	Game  Game     `json:"game"`
}

func generateRoom(owner User, roomID RoomID) *Room {
	ctx, cancel := context.WithCancel(context.Background())

	return &Room{
		ID:    roomID,
		Owner: owner.Username,
		Users: make(ConnectedUsers),
		Game: Game{
			Title:    "",
			Voted:    make(Voted),
			Votes:    make(Votes),
			Guesses:  make(Guesses),
			Revealed: false,
		},
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func (r *Room) ToSerializable() SerializableRoom {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	users := make(Users)

	for username, connectedUser := range r.Users {
		users[username] = User{
			Username: connectedUser.User.Username,
			Role:     connectedUser.User.Role,
			Score:    connectedUser.User.Score,
		}
	}

	voted := make(Voted)
	for k, v := range r.Game.Voted {
		voted[k] = v
	}

	votes := make(Votes)
	for k, v := range r.Game.Votes {
		votes[k] = v
	}

	guesses := make(Guesses)
	for k, v := range r.Game.Guesses {
		guesses[k] = v
	}

	game := Game{
		Title:    r.Game.Title,
		Voted:    voted,
		Votes:    votes,
		Guesses:  guesses,
		Revealed: r.Game.Revealed,
	}

	return SerializableRoom{
		ID:    r.ID,
		Owner: r.Owner,
		Users: users,
		Game:  game,
	}
}

func (r *Room) AddUser(user User, conn *websocket.Conn) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if _, exists := r.Users[user.Username]; exists {
		return errors.New("user already in room")
	}

	cu := &ConnectedUser{
		User: user,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	r.Users[user.Username] = cu

	go r.readPump(cu)
	go r.writePump(cu)

	go r.scheduleBroadcast()
	return nil
}

func (r *Room) RemoveUser(username Username) RoomStatus {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.Users, username)

	if len(r.Users) == 0 {
		r.EmptySince = time.Now()
		return RoomEmpty
	}

	if username == r.Owner {
		var firstDeveloper Username
		var firstTester Username
		var firstPO Username
		var firstPlanner Username

		r.Owner = ""

		for _, user := range r.Users {
			userRole := user.User.Role

			if userRole == Developer && firstDeveloper == "" {
				firstDeveloper = user.User.Username
			} else if userRole == Tester && firstTester == "" {
				firstTester = user.User.Username
			} else if userRole == PO && firstPO == "" {
				firstPO = user.User.Username
			} else if userRole == Planner && firstPlanner == "" {
				firstPlanner = user.User.Username
			}

			if len(firstDeveloper) > 0 && len(firstTester) > 0 && len(firstPO) > 0 && len(firstPlanner) > 0 {
				break
			}
		}

		if len(firstPlanner) > 0 {
			r.Owner = firstPlanner
		} else if len(firstPO) > 0 {
			r.Owner = firstPO
		} else if len(firstDeveloper) > 0 {
			r.Owner = firstDeveloper
		} else if len(firstTester) > 0 {
			r.Owner = firstTester
		}
	}

	go r.scheduleBroadcast()
	return RoomInUse
}

func (r *Room) ResetVoting() {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	r.Game.Voted = make(Voted)
	r.Game.Votes = make(Votes)
	r.Game.Guesses = make(Guesses)
	r.Game.Revealed = false
	r.Game.Status = StatusWaiting

	go r.scheduleBroadcast()
}

func (r *Room) StartVoting() {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	r.Game.Status = StatusPlaying
	go r.scheduleBroadcast()
}

// Changing votes is allowed
func (r *Room) SubmitVote(username Username, vote VoteValue) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if r.Game.Status != StatusPlaying {
		return fmt.Errorf("game is not in progress")
	}

	r.Game.Voted[username] = true
	r.Game.Votes[username] = vote

	go r.scheduleBroadcast()
	return nil
}

func (r *Room) MakeGuess(username Username, guess Guess) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if r.Game.Status != StatusPlaying {
		return fmt.Errorf("game is not in progress")
	}

	if _, ok := r.Users[guess.GuessingOn]; !ok {
		return fmt.Errorf("user you're guessing on is not currently in the game")
	}

	if username == guess.GuessingOn {
		return fmt.Errorf("you cannot place a bet on yourself")
	}

	r.Game.Guesses[username] = guess

	go r.scheduleBroadcast()
	return nil
}

func (r *Room) RevealVotes(revealer Username) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if revealer != r.Owner {
		return fmt.Errorf("only game owner can reveal votes")
	}

	for username, guess := range r.Game.Guesses {
		if r.Game.Votes[guess.GuessingOn] == guess.Prediction {
			if user, exists := r.Users[username]; exists {
				user.Mutex.Lock()
				user.User.Score += 5
				user.Mutex.Unlock()
			}
		}
	}

	r.Game.Revealed = true
	r.Game.Status = StatusFinished

	go r.scheduleBroadcast()
	return nil
}

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

	case MessageTypeGuess:
		var payload GuessPayload
		if err := json.Unmarshal(wsMessage.Payload, &payload); err != nil {
			return fmt.Errorf("invalid guess payload: %v", err)
		}
		guess := Guess{
			GuessingOn: payload.GuessingOn,
			Prediction: payload.Prediction,
		}
		return r.MakeGuess(cu.User.Username, guess)

	case MessageTypeStartVoting:
		// Only room owner can start voting
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

	case MessageTypeResetVoting:
		// Only room owner can reset voting
		if cu.User.Username != r.Owner {
			return errors.New("only room owner can reset voting")
		}
		r.ResetVoting()
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
		"type":    "state_update",
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
			// If send buffer is full, skip this user
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
