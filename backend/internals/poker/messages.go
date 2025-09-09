package poker

import "encoding/json"

type MessageType string

const (
	MessageTypeJoin        MessageType = "join"
	MessageTypeVote        MessageType = "vote"
	MessageTypeGuess       MessageType = "guess"
	MessageTypeStartVoting MessageType = "start_voting"
	MessageTypeRevealVotes MessageType = "reveal_votes"
	MessageTypeResetVoting MessageType = "reset_voting"
)

type WSMessage struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Specific payload structures
type VotePayload struct {
	Vote VoteValue `json:"vote"`
}

type GuessPayload struct {
	GuessingOn Username  `json:"guessingOn"`
	Prediction VoteValue `json:"prediction"`
}
