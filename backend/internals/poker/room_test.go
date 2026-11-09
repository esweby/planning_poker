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

func TestViewerSeesOnlyOwnVoteUntilReveal(t *testing.T) {
	r := generateRoom(User{Username: "owner"}, "12345678")
	defer r.Close()

	r.StartVoting()
	if err := r.SubmitVote("alice", "8"); err != nil {
		t.Fatal(err)
	}
	if err := r.SubmitVote("bob", "3"); err != nil {
		t.Fatal(err)
	}

	alice := r.ToSerializableFor("alice")
	if alice.Game.Votes["alice"] != "8" {
		t.Fatal("viewer must see their own vote before reveal")
	}
	if _, leaked := alice.Game.Votes["bob"]; leaked {
		t.Fatal("viewer must not see other votes before reveal")
	}
	if len(r.ToSerializableFor("carol").Game.Votes) != 0 {
		t.Fatal("non-voter must see no votes before reveal")
	}

	if err := r.RevealVotes("owner"); err != nil {
		t.Fatal(err)
	}
	if len(r.ToSerializableFor("carol").Game.Votes) != 2 {
		t.Fatal("all votes must be visible after reveal")
	}

	r.RestartVoting()
	if len(r.ToSerializableFor("alice").Game.Votes) != 0 {
		t.Fatal("votes must be cleared and hidden after restart")
	}
}
