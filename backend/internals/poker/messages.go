package poker

import "encoding/json"

type MessageType string

const (
	MessageTypeJoin        MessageType = "join"
	MessageTypeVote        MessageType = "vote"
	MessageTypeStartVoting MessageType = "start_voting"
	MessageTypeRevealVotes MessageType = "reveal_votes"
	MessageTypeRestartVoting MessageType = "restart_voting"
)

type WSMessage struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type VotePayload struct {
	Vote VoteValue `json:"vote"`
}
