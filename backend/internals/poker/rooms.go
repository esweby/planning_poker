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

func (rm *RoomManager) CreateNewRoom(owner User) RoomID {
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
	rm.rooms[roomID] = room

	return roomID
}

func (rm *RoomManager) GetRoom(roomID RoomID) (*Room, error) {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	room, exists := rm.rooms[roomID]
	log.Printf("-- GetRoom %s - exists %t", roomID, exists)
	if !exists {
		return nil, errors.New("room not found")
	}

	return room, nil
}

func (rm *RoomManager) RemoveRoom(roomID RoomID) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	delete(rm.rooms, roomID)
}

func (rm *RoomManager) AddUserToRoom(roomID RoomID, user User, conn *websocket.Conn) error {
	room, err := rm.GetRoom(roomID)
	if err != nil {
		return err
	}

	err = room.AddUser(user, conn)
	if err != nil {
		return err
	}

	log.Printf("- Room (%s): %s joined", roomID, user.Username)
	return nil
}

func (rm *RoomManager) RemoveUserFromRoom(roomID RoomID, username Username) (RoomStatus, error) {
	room, err := rm.GetRoom(roomID)
	if err != nil {
		return RoomEmpty, err
	}

	log.Printf("- Room (%s): %s removed", roomID, username)
	status := room.RemoveUser(username)
	return status, nil
}

func (rm *RoomManager) GetSerializableRoom(roomID RoomID) (SerializableRoom, error) {
	room, err := rm.GetRoom(roomID)
	if err != nil {
		return SerializableRoom{}, err
	}

	return room.ToSerializable(), nil
}

func (rm *RoomManager) cleanupRoutine() {
	ticker := time.NewTicker(60 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rm.mutex.Lock()
		now := time.Now()
		roomsToRemove := make([]RoomID, 0)

		for roomID, room := range rm.rooms {
			room.Mutex.RLock()
			if len(room.Users) == 0 {
				if room.EmptySince.IsZero() {
					room.EmptySince = now
				} else if now.Sub(room.EmptySince) > 5*time.Minute {
					roomsToRemove = append(roomsToRemove, roomID)
				}
			} else {
				room.EmptySince = time.Time{}
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
