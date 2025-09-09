package poker

import (
	"context"
	"errors"
	"fmt"
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
			Status:   StatusWaiting,
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
			Seed: 	  connectedUser.User.Seed,
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
		Status:	  r.Game.Status,
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
	delete(r.Game.Voted, username)
	delete(r.Game.Votes, username)
	delete(r.Game.Guesses, username)

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
