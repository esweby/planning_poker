package poker

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type RoomManager struct {
	rooms map[RoomID]*Room
	mutex sync.RWMutex
}

var (
	roomManager     *RoomManager
	roomManagerOnce sync.Once
)

func GetRoomManager() *RoomManager {
	roomManagerOnce.Do(func() {
		roomManager = &RoomManager{
			rooms: make(map[RoomID]*Room),
		}
		go roomManager.cleanupRoutine()
	})

	return roomManager
}

// ErrRoomNotFound is returned when no room has the given ID.
var ErrRoomNotFound = errors.New("room not found")

const (
	cleanupInterval = time.Minute
	emptyRoomTTL    = 5 * time.Minute
)

// CreateNewRoom creates a room owned by owner and returns its ID and the
// owner's session token, which the owner must present when joining.
func (rm *RoomManager) CreateNewRoom(owner User) (RoomID, string, error) {
	if err := ValidateUsername(owner.Username); err != nil {
		return "", "", err
	}

	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	var roomID RoomID

	for {
		roomID = Generate8DigitString()
		if _, ok := rm.rooms[roomID]; !ok {
			break
		}
	}

	log.Printf("- Room (%s): %s created room", roomID, owner.Username)
	room := generateRoom(owner, roomID)
	token := room.claimLocked(owner.Username)
	// An unjoined room counts as empty, so it is cleaned up if nobody comes.
	room.EmptySince = time.Now()
	rm.rooms[roomID] = room

	return roomID, token, nil
}

func (rm *RoomManager) GetRoom(roomID RoomID) (*Room, error) {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	room, exists := rm.rooms[roomID]
	log.Printf("-- GetRoom %s - exists %t", roomID, exists)
	if !exists {
		return nil, ErrRoomNotFound
	}

	return room, nil
}

func (rm *RoomManager) RemoveRoom(roomID RoomID) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	delete(rm.rooms, roomID)
}

// JoinRoom connects user to a room and returns their session token.
func (rm *RoomManager) JoinRoom(roomID RoomID, user User, token string, conn *websocket.Conn) (string, error) {
	room, err := rm.GetRoom(roomID)
	if err != nil {
		return "", err
	}

	token, err = room.Join(user, token, conn)
	if err != nil {
		return "", err
	}

	log.Printf("- Room (%s): %s joined", roomID, user.Username)
	return token, nil
}

func (rm *RoomManager) GetSerializableRoom(roomID RoomID) (SerializableRoom, error) {
	room, err := rm.GetRoom(roomID)
	if err != nil {
		return SerializableRoom{}, err
	}

	return room.ToSerializable(), nil
}

func (rm *RoomManager) cleanupRoutine() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rm.mutex.Lock()
		now := time.Now()
		roomsToRemove := make([]RoomID, 0)

		for roomID, room := range rm.rooms {
			room.Mutex.RLock()
			if len(room.Users) == 0 {
				if !room.EmptySince.IsZero() && now.Sub(room.EmptySince) > emptyRoomTTL {
					roomsToRemove = append(roomsToRemove, roomID)
				}
			}
			room.Mutex.RUnlock()
		}

		for _, roomID := range roomsToRemove {
			if room, exists := rm.rooms[roomID]; exists {
				room.Close()
				delete(rm.rooms, roomID)
				fmt.Printf("Removed empty room: %s\n", roomID)
			}
		}
		rm.mutex.Unlock()
	}
}

func (rm *RoomManager) GetAllRooms() map[RoomID]SerializableRoom {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	result := make(map[RoomID]SerializableRoom)
	for roomID, room := range rm.rooms {
		result[roomID] = room.ToSerializable()
	}

	return result
}

func (rm *RoomManager) RoomExists(roomID RoomID) bool {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	_, exists := rm.rooms[roomID]
	return exists
}

func (rm *RoomManager) GetRoomCount() int {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	return len(rm.rooms)
}

func (rm *RoomManager) GetUserCount() int {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	total := 0
	for _, room := range rm.rooms {
		room.Mutex.RLock()
		total += len(room.Users)
		room.Mutex.RUnlock()
	}

	return total
}
