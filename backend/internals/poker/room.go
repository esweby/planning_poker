package poker

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
)

type RoomID = string

type RoomStatus = string

const (
	RoomInUse         = "roomInUse"
	RoomEmpty         = "roomEmpty"
	writeWait         = 10 * time.Second
	pongWait          = 60 * time.Second
	pingPeriod        = (pongWait * 9) / 10
	maxMessageSize    = 512
	maxUsernameLength = 32
	broadcastDelay    = 100 * time.Millisecond
)

var (
	ErrInvalidUsername = errors.New("username must be 1-32 characters")
	ErrNameTaken       = errors.New("that name is already taken in this room")
	ErrRoomClosed      = errors.New("room is closed")
)

type Room struct {
	ID         RoomID
	Owner      Username
	Users      ConnectedUsers
	Game       Game
	Mutex      sync.RWMutex
	EmptySince time.Time

	// sessions maps each name claimed in this room to the secret token that
	// proves ownership of it. Claims last for the room's lifetime, so nobody
	// can take over a name (including the owner's) after its holder drops.
	sessions map[Username]string
	closed   bool

	// Throttling broadcasts: at most one is pending at a time.
	broadcastPending bool
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
	return &Room{
		ID:       roomID,
		Owner:    owner.Username,
		Users:    make(ConnectedUsers),
		sessions: make(map[Username]string),
		Game: Game{
			Title:    "",
			Status:   StatusWaiting,
			Voted:    make(Voted),
			Votes:    make(Votes),
			Revealed: false,
		},
	}
}

// ValidateUsername reports whether name is acceptable as a display name.
func ValidateUsername(name Username) error {
	trimmed := strings.TrimSpace(string(name))
	if trimmed == "" || trimmed != string(name) || utf8.RuneCountInString(trimmed) > maxUsernameLength {
		return ErrInvalidUsername
	}
	return nil
}

func newToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("crypto/rand failed: %v", err))
	}
	return hex.EncodeToString(b)
}

// claimLocked reserves name and returns its new token. Caller holds r.Mutex.
func (r *Room) claimLocked(name Username) string {
	token := newToken()
	r.sessions[name] = token
	return token
}

// checkJoinLocked validates a join attempt. Caller holds r.Mutex.
func (r *Room) checkJoinLocked(name Username, token string) error {
	if r.closed {
		return ErrRoomClosed
	}
	if err := ValidateUsername(name); err != nil {
		return err
	}
	if claimed, ok := r.sessions[name]; ok {
		if subtle.ConstantTimeCompare([]byte(claimed), []byte(token)) != 1 {
			return ErrNameTaken
		}
	}
	return nil
}

// CheckJoin reports whether a join with this name and token would succeed,
// without changing anything. Used to reject bad joins before the websocket
// upgrade so the client gets a readable HTTP error.
func (r *Room) CheckJoin(name Username, token string) error {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.checkJoinLocked(name, token)
}

// ToSerializable returns a snapshot safe to show to anyone: vote values are
// hidden until they are revealed.
func (r *Room) ToSerializable() SerializableRoom {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.serializeForLocked("")
}

// ToSerializableFor returns a snapshot for one viewer: before reveal it
// includes only that viewer's own vote.
func (r *Room) ToSerializableFor(viewer Username) SerializableRoom {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.serializeForLocked(viewer)
}

// serializeForLocked builds a snapshot for viewer (empty means no viewer).
// The caller must hold r.Mutex.
func (r *Room) serializeForLocked(viewer Username) SerializableRoom {
	users := make(Users)

	for username, connectedUser := range r.Users {
		users[username] = User{
			Username: connectedUser.User.Username,
			Role:     connectedUser.User.Role,
			Score:    connectedUser.User.Score,
			Seed:     connectedUser.User.Seed,
		}
	}

	voted := make(Voted)
	for k, v := range r.Game.Voted {
		voted[k] = v
	}

	votes := make(Votes)
	if r.Game.Revealed {
		for k, v := range r.Game.Votes {
			votes[k] = v
		}
	} else if vote, ok := r.Game.Votes[viewer]; ok && viewer != "" {
		votes[viewer] = vote
	}

	game := Game{
		Title:    r.Game.Title,
		Status:   r.Game.Status,
		Voted:    voted,
		Votes:    votes,
		Revealed: r.Game.Revealed,
	}

	return SerializableRoom{
		ID:    r.ID,
		Owner: r.Owner,
		Users: users,
		Game:  game,
	}
}

// Join connects user to the room. A new name is claimed and given a fresh
// token; a claimed name needs its token, and then replaces any existing
// (possibly stale) connection for that name. It returns the session token.
func (r *Room) Join(user User, token string, conn *websocket.Conn) (string, error) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if err := r.checkJoinLocked(user.Username, token); err != nil {
		return "", err
	}

	if _, claimed := r.sessions[user.Username]; !claimed {
		token = r.claimLocked(user.Username)
	}

	if old, ok := r.Users[user.Username]; ok {
		// Its readPump will exit and see it has been replaced.
		old.Conn.Close()
	}

	cu := &ConnectedUser{
		User: user,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	session, _ := json.Marshal(map[string]string{"type": "session", "token": token})
	cu.Send <- session

	r.Users[user.Username] = cu
	r.EmptySince = time.Time{}

	go r.readPump(cu)
	go r.writePump(cu)

	r.scheduleBroadcast()
	return token, nil
}

// removeConn removes cu if it is still the live connection for its name. A
// connection that has been replaced by a reconnect is ignored, so it cannot
// remove its replacement.
func (r *Room) removeConn(cu *ConnectedUser) RoomStatus {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	username := cu.User.Username
	if r.Users[username] != cu {
		return RoomInUse
	}

	delete(r.Users, username)
	delete(r.Game.Voted, username)
	delete(r.Game.Votes, username)

	if len(r.Users) == 0 {
		r.EmptySince = time.Now()
		return RoomEmpty
	}

	r.scheduleBroadcast()
	return RoomInUse
}

func (r *Room) StartVoting() {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	r.Game.Voted = make(Voted)
	r.Game.Votes = make(Votes)
	r.Game.Revealed = false
	r.Game.Status = StatusPlaying
	r.scheduleBroadcast()
}

func (r *Room) SubmitVote(username Username, vote VoteValue) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if r.Game.Status != StatusPlaying {
		return fmt.Errorf("game is not in progress")
	}

	r.Game.Voted[username] = true
	r.Game.Votes[username] = vote

	r.scheduleBroadcast()
	return nil
}

func (r *Room) RevealVotes(revealer Username) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if revealer != r.Owner {
		return fmt.Errorf("only game owner can reveal votes")
	}

	r.Game.Revealed = true
	r.Game.Status = StatusFinished

	r.scheduleBroadcast()
	return nil
}

func (r *Room) RestartVoting() {
	r.StartVoting()
}
