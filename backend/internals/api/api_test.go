package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	gin.SetMode(gin.TestMode)
	srv := httptest.NewServer(NewRouter())
	t.Cleanup(srv.Close)
	return srv
}

func createRoom(t *testing.T, srv *httptest.Server, owner string) (string, string) {
	t.Helper()
	body, _ := json.Marshal(map[string]any{"owner": map[string]string{"username": owner, "role": "Developer"}})
	res, err := http.Post(srv.URL+"/api/", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("create room: status %d", res.StatusCode)
	}
	var out struct {
		RoomNumber string `json:"roomNumber"`
		Token      string `json:"token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if out.Token == "" {
		t.Fatal("owner token missing")
	}
	return out.RoomNumber, out.Token
}

func joinURL(srv *httptest.Server, room, name, token string) string {
	q := url.Values{"username": {name}, "role": {"Developer"}, "token": {token}}
	return "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/join/" + room + "?" + q.Encode()
}

type client struct {
	t     *testing.T
	conn  *websocket.Conn
	token string
}

// dial joins the room and waits for the session token.
func dial(t *testing.T, srv *httptest.Server, room, name, token string) *client {
	t.Helper()
	conn, res, err := websocket.DefaultDialer.Dial(joinURL(srv, room, name, token), nil)
	if err != nil {
		status := 0
		if res != nil {
			status = res.StatusCode
		}
		t.Fatalf("dial %s: %v (status %d)", name, err, status)
	}
	t.Cleanup(func() { conn.Close() })

	c := &client{t: t, conn: conn}
	msg := c.next(func(m map[string]any) bool { return m["type"] == "session" })
	c.token = msg["token"].(string)
	return c
}

// next reads messages until one matches.
func (c *client) next(match func(map[string]any) bool) map[string]any {
	c.t.Helper()
	c.conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			c.t.Fatalf("read: %v", err)
		}
		var m map[string]any
		if err := json.Unmarshal(data, &m); err != nil {
			c.t.Fatalf("bad json %s: %v", data, err)
		}
		if match(m) {
			return m
		}
	}
}

// nextState waits for a room state where cond holds on the game object.
func (c *client) nextState(cond func(game map[string]any) bool) map[string]any {
	c.t.Helper()
	msg := c.next(func(m map[string]any) bool {
		p, ok := m["payload"].(map[string]any)
		return ok && cond(p["game"].(map[string]any))
	})
	return msg["payload"].(map[string]any)["game"].(map[string]any)
}

func (c *client) send(msgType string, payload any) {
	c.t.Helper()
	if err := c.conn.WriteJSON(map[string]any{"type": msgType, "payload": payload}); err != nil {
		c.t.Fatal(err)
	}
}

func dialStatus(t *testing.T, srv *httptest.Server, room, name, token string, header http.Header) int {
	t.Helper()
	conn, res, err := websocket.DefaultDialer.Dial(joinURL(srv, room, name, token), header)
	if err == nil {
		conn.Close()
		return http.StatusSwitchingProtocols
	}
	if res == nil {
		t.Fatalf("dial failed without response: %v", err)
	}
	return res.StatusCode
}

func TestOwnerNameCannotBeTakenWithoutToken(t *testing.T) {
	srv := newTestServer(t)
	room, ownerToken := createRoom(t, srv, "owner")

	if got := dialStatus(t, srv, room, "owner", "", nil); got != http.StatusConflict {
		t.Fatalf("impersonating owner: got %d, want 409", got)
	}
	if got := dialStatus(t, srv, room, "owner", "wrong", nil); got != http.StatusConflict {
		t.Fatalf("wrong token: got %d, want 409", got)
	}

	res, err := http.Get(srv.URL + "/api/check/" + room + "?username=owner")
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusConflict {
		t.Fatalf("check endpoint: got %d, want 409", res.StatusCode)
	}

	owner := dial(t, srv, room, "owner", ownerToken)
	if owner.token != ownerToken {
		t.Fatal("owner should keep their token")
	}
}

func TestJoinValidation(t *testing.T) {
	srv := newTestServer(t)
	room, _ := createRoom(t, srv, "owner")

	if got := dialStatus(t, srv, room, "", "", nil); got != http.StatusBadRequest {
		t.Fatalf("empty name: got %d, want 400", got)
	}
	if got := dialStatus(t, srv, "00000000x", "alice", "", nil); got != http.StatusNotFound {
		t.Fatalf("missing room: got %d, want 404", got)
	}
	header := http.Header{"Origin": {"https://evil.example"}}
	if got := dialStatus(t, srv, room, "alice", "", header); got != http.StatusForbidden {
		t.Fatalf("bad origin: got %d, want 403", got)
	}
}

func TestReconnectWithTokenReplacesStaleConnection(t *testing.T) {
	srv := newTestServer(t)
	room, _ := createRoom(t, srv, "owner")

	first := dial(t, srv, room, "alice", "")
	if got := dialStatus(t, srv, room, "alice", "", nil); got != http.StatusConflict {
		t.Fatalf("second alice without token: got %d, want 409", got)
	}

	second := dial(t, srv, room, "alice", first.token)
	if second.token != first.token {
		t.Fatal("reconnect should keep the same token")
	}

	// The old connection is closed by the server...
	first.conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	for {
		if _, _, err := first.conn.ReadMessage(); err != nil {
			break
		}
	}

	// ...and its cleanup must not remove the new one.
	time.Sleep(200 * time.Millisecond)
	if _, err := roomManager.GetRoom(room); err != nil {
		t.Fatal(err)
	}
	r, _ := roomManager.GetRoom(room)
	if r.ToSerializable().Users["alice"].Username != "alice" {
		t.Fatal("alice should still be in the room after the old connection closed")
	}
}

func TestVotesStayPrivateUntilReveal(t *testing.T) {
	srv := newTestServer(t)
	room, ownerToken := createRoom(t, srv, "owner")

	owner := dial(t, srv, room, "owner", ownerToken)
	alice := dial(t, srv, room, "alice", "")

	owner.send("start_voting", nil)
	alice.nextState(func(g map[string]any) bool { return g["status"] == "playing" })
	alice.send("vote", map[string]string{"vote": "8"})

	aliceGame := alice.nextState(func(g map[string]any) bool {
		return g["votes"].(map[string]any)["alice"] == "8"
	})
	if len(aliceGame["votes"].(map[string]any)) != 1 {
		t.Fatal("alice should see only her own vote")
	}

	ownerGame := owner.nextState(func(g map[string]any) bool {
		return g["voted"].(map[string]any)["alice"] == true
	})
	if len(ownerGame["votes"].(map[string]any)) != 0 {
		t.Fatalf("owner saw votes before reveal: %v", ownerGame["votes"])
	}

	// Non-owners can't reveal, and get an error reply.
	alice.send("reveal_votes", nil)
	alice.next(func(m map[string]any) bool { return m["error"] != nil })

	owner.send("reveal_votes", nil)
	revealed := owner.nextState(func(g map[string]any) bool { return g["revealed"] == true })
	if revealed["votes"].(map[string]any)["alice"] != "8" {
		t.Fatal("votes should be visible after reveal")
	}
}
