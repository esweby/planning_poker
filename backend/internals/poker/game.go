package poker

type VoteValue string // Can be "0", "1", "2", "3", "5", "8", "13", "21", "?"
type Status string    // "waiting", "playing", "finished"
const (
	StatusWaiting  = "waiting"
	StatusPlaying  = "playing"
	StatusFinished = "finished"
)

type Voted = map[Username]bool
type Votes = map[Username]VoteValue

type Game struct {
	Title    string  `json:"title"`
	Status   Status  `json:"status"`
	Voted    Voted   `json:"voted"`
	Votes    Votes   `json:"votes"`
	Revealed bool    `json:"revealed"`
}
