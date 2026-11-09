package poker

import "testing"

func TestVotesHiddenUntilReveal(t *testing.T) {
	r := generateRoom(User{Username: "owner"}, "12345678")
	defer r.Close()
	r.Game.Votes["alice"] = "8"
	r.Game.Voted["alice"] = true
	hidden := r.ToSerializable()
	if len(hidden.Game.Votes) != 0 || !hidden.Game.Voted["alice"] {
		t.Fatal("unrevealed state must hide values but preserve participation")
	}
	r.Game.Revealed = true
	revealed := r.ToSerializable()
	if revealed.Game.Votes["alice"] != "8" {
		t.Fatal("revealed state must include votes")
	}
	revealed.Game.Votes["alice"] = "13"
	if r.Game.Votes["alice"] != "8" {
		t.Fatal("snapshot must not mutate stored votes")
	}
}
