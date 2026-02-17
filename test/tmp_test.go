package test

import (
	"bluebell/initialize"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"go.uber.org/zap"
)

func Serv(w http.ResponseWriter, r *http.Request) {
	accept, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		zap.L().Debug("failed to accept websocket connection", zap.Error(err))
		return
	}
	defer accept.Close(websocket.StatusInternalError, "internal error")
	for {
		_, msg, err := accept.Read(r.Context())
		if err != nil {
			zap.L().Debug("failed to read message", zap.Error(err))
			return
		}
		zap.L().Debug("received message", zap.ByteString("msg", msg))
		err = accept.Write(r.Context(), websocket.MessageText, msg)
		if err != nil {
			zap.L().Debug("failed to write message", zap.Error(err))
			return
		}
	}
}

func TestSer(t *testing.T) {
	initialize.Initialize()
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", Serv)

	server := httptest.NewServer(mux)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial websocket failed: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "test done")

	messages := []string{"hello", "coder/websocket", "real-time echo"}
	for _, msg := range messages {
		if err := conn.Write(ctx, websocket.MessageText, []byte(msg)); err != nil {
			t.Fatalf("write message failed: %v", err)
		}

		_, got, err := conn.Read(ctx)
		if err != nil {
			t.Fatalf("read echo failed: %v", err)
		}

		if string(got) != msg {
			t.Fatalf("echo mismatch: want %q, got %q", msg, string(got))
		}
		t.Logf("client -> server: %q | server -> client: %q", msg, string(got))
	}

}

func TestClient(t *testing.T) {
	initialize.Initialize()
	pushMessages := []string{
		"server push #1",
		"server push #2",
		"server push #3",
	}
	ackMessages := make([]string, 0, len(pushMessages))
	done := make(chan struct{})

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
		if err != nil {
			t.Errorf("accept websocket failed: %v", err)
			return
		}
		defer conn.Close(websocket.StatusNormalClosure, "server done")

		ctx := r.Context()
		for i, push := range pushMessages {
			if err := conn.Write(ctx, websocket.MessageText, []byte(push)); err != nil {
				t.Errorf("server push #%d failed: %v", i+1, err)
				return
			}
			_, ack, err := conn.Read(ctx)
			if err != nil {
				t.Errorf("server read ack #%d failed: %v", i+1, err)
				return
			}
			ackMessages = append(ackMessages, string(ack))
			t.Logf("server <- client ack #%d: %q", i+1, string(ack))
		}
		close(done)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial websocket failed: %v", err)
	}
	defer client.Close(websocket.StatusNormalClosure, "client done")

	for i, wantPush := range pushMessages {
		_, gotPush, err := client.Read(ctx)
		if err != nil {
			t.Fatalf("client read push #%d failed: %v", i+1, err)
		}
		if string(gotPush) != wantPush {
			t.Fatalf("push mismatch #%d: want %q, got %q", i+1, wantPush, string(gotPush))
		}
		t.Logf("server -> client push #%d: %q", i+1, string(gotPush))

		ack := fmt.Sprintf("ack #%d for %s", i+1, gotPush)
		if err := client.Write(ctx, websocket.MessageText, []byte(ack)); err != nil {
			t.Fatalf("client send ack #%d failed: %v", i+1, err)
		}
	}
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("server did not finish handling acknowledgements in time")
	}

	if len(ackMessages) != len(pushMessages) {
		t.Fatalf("ack count mismatch: want %d, got %d", len(pushMessages), len(ackMessages))
	}
	for i, ack := range ackMessages {
		want := fmt.Sprintf("ack #%d for %s", i+1, pushMessages[i])
		if ack != want {
			t.Fatalf("ack mismatch #%d: want %q, got %q", i+1, want, ack)
		}
	}

}
