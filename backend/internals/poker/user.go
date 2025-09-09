package poker

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Username string
type Role string // developer, tester, po, planner

const (
	Developer = "developer"
	Tester    = "tester"
	PO        = "po"
	Planner   = "planner"
)

type ConnectedUser struct {
	User  User
	Conn  *websocket.Conn
	Send  chan []byte
	Mutex sync.Mutex
}

type ConnectedUsers = map[Username]*ConnectedUser

type User struct {
	Username Username `json:"username"`
	Role     Role     `json:"role"`
	Seed 	 string   `json:"seed"`
	Score    int32    `json:"score"`
}

type Users = map[Username]User
