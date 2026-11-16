package presence

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"
	"testing"
	"time"
)

// TestDiscordSessionPublishesArtworkState uses a local fake Discord pipe to verify handshake, both activity labels, and keepalive replies.
func TestDiscordSessionPublishesArtworkState(t *testing.T) {
	application, discord := net.Pipe()
	defer discord.Close()
	_ = discord.SetDeadline(time.Now().Add(3 * time.Second))
	service := &Service{appID: "123456789012345678", changed: make(chan struct{}, 1), stop: make(chan struct{})}
	finished := make(chan error, 1)
	go func() { finished <- service.serve(application) }()

	opcode, payload, err := readFrame(discord)
	if err != nil || opcode != 0 {
		t.Fatalf("missing Discord handshake: opcode=%d error=%v", opcode, err)
	}
	var hello struct {
		Version  int    `json:"v"`
		ClientID string `json:"client_id"`
	}
	if err := json.Unmarshal(payload, &hello); err != nil || hello.Version != 1 || hello.ClientID != service.appID {
		t.Fatalf("invalid handshake: %+v error=%v", hello, err)
	}
	respond := &rpcWriter{connection: discord}
	if err := respond.write(1, map[string]string{"cmd": "DISPATCH", "evt": "READY"}); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"Drawing a piece of art", "Admiring a piece of art"} {
		if expected == "Admiring a piece of art" {
			service.SetArtworkOpen(true)
		}
		opcode, payload, err = readFrame(discord)
		if err != nil || opcode != 1 {
			t.Fatalf("missing activity update: opcode=%d error=%v", opcode, err)
		}
		var request activityRequest
		if err := json.Unmarshal(payload, &request); err != nil || request.Command != "SET_ACTIVITY" || request.Args.PID <= 0 || request.Args.Activity.Type != 0 || request.Args.Activity.Details != expected {
			t.Fatalf("wrong activity for %q: %+v error=%v", expected, request, err)
		}
	}
	if err := respond.write(3, map[string]string{"keepalive": "yes"}); err != nil {
		t.Fatal(err)
	}
	opcode, _, err = readFrame(discord)
	if err != nil || opcode != 4 {
		t.Fatalf("missing Discord keepalive response: opcode=%d error=%v", opcode, err)
	}
	close(service.stop)
	select {
	case err := <-finished:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Discord session did not stop")
	}
}

// TestReadFrameRejectsOversizedPayload checks that untrusted IPC lengths cannot force an unbounded allocation.
func TestReadFrameRejectsOversizedPayload(t *testing.T) {
	header := [8]byte{}
	binary.LittleEndian.PutUint32(header[4:], maxFrameSize+1)
	if _, _, err := readFrame(bytes.NewReader(header[:])); err == nil {
		t.Fatal("oversized Discord frame was accepted")
	}
}
