package presence

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

const maxFrameSize = 64 * 1024

type rpcMessage struct {
	Command string `json:"cmd"`
	Event   string `json:"evt"`
}

type rpcWriter struct {
	connection io.Writer
	mu         sync.Mutex
	nonce      atomic.Uint64
}

type activityRequest struct {
	Command string `json:"cmd"`
	Args    struct {
		PID      int `json:"pid"`
		Activity struct {
			Type    int    `json:"type"`
			Details string `json:"details"`
		} `json:"activity"`
	} `json:"args"`
	Nonce string `json:"nonce"`
}

// writeAll takes a stream and payload, then writes every byte or returns the first write error.
func writeAll(connection io.Writer, payload []byte) error {
	for len(payload) > 0 {
		written, err := connection.Write(payload)
		if err != nil {
			return err
		}
		if written == 0 {
			return io.ErrShortWrite
		}
		payload = payload[written:]
	}
	return nil
}

// write takes an IPC opcode and JSON value, then writes one complete Discord frame without interleaving other writers.
func (writer *rpcWriter) write(opcode uint32, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if len(payload) > maxFrameSize {
		return errors.New("Discord frame exceeds maximum size")
	}
	header := [8]byte{}
	binary.LittleEndian.PutUint32(header[:4], opcode)
	binary.LittleEndian.PutUint32(header[4:], uint32(len(payload)))
	writer.mu.Lock()
	defer writer.mu.Unlock()
	if err := writeAll(writer.connection, header[:]); err != nil {
		return err
	}
	return writeAll(writer.connection, payload)
}

// readFrame takes an IPC stream and returns the next bounded Discord opcode and JSON payload.
func readFrame(connection io.Reader) (uint32, []byte, error) {
	header := [8]byte{}
	if _, err := io.ReadFull(connection, header[:]); err != nil {
		return 0, nil, err
	}
	length := binary.LittleEndian.Uint32(header[4:])
	if length > maxFrameSize {
		return 0, nil, errors.New("Discord frame exceeds maximum size")
	}
	payload := make([]byte, length)
	if _, err := io.ReadFull(connection, payload); err != nil {
		return 0, nil, err
	}
	return binary.LittleEndian.Uint32(header[:4]), payload, nil
}

// handshake takes a Discord IPC connection and application ID, then waits for the matching READY event.
func handshake(connection io.ReadWriteCloser, writer *rpcWriter, appID string) error {
	if err := writer.write(0, struct {
		Version  int    `json:"v"`
		ClientID string `json:"client_id"`
	}{1, appID}); err != nil {
		return err
	}
	for {
		opcode, payload, err := readFrame(connection)
		if err != nil {
			return err
		}
		if opcode == 2 {
			return errors.New("Discord closed the handshake")
		}
		if opcode == 3 {
			if err := writer.write(4, json.RawMessage(payload)); err != nil {
				return err
			}
			continue
		}
		if opcode != 1 {
			return fmt.Errorf("unexpected Discord opcode %d", opcode)
		}
		var message rpcMessage
		if err := json.Unmarshal(payload, &message); err != nil {
			return err
		}
		if message.Command == "DISPATCH" && message.Event == "READY" {
			return nil
		}
		if message.Event == "ERROR" {
			return errors.New("Discord rejected the application ID")
		}
	}
}

// setActivity takes an artwork state and publishes its Playing activity for the current IFLY process.
func setActivity(writer *rpcWriter, artworkOpen bool) error {
	request := activityRequest{Command: "SET_ACTIVITY", Nonce: strconv.FormatUint(writer.nonce.Add(1), 10)}
	request.Args.PID = os.Getpid()
	request.Args.Activity.Type = 0
	request.Args.Activity.Details = activityDetails(artworkOpen)
	return writer.write(1, request)
}

// readEvents takes an established Discord IPC session, answers keepalive frames, and reports when the connection closes.
func readEvents(connection io.Reader, writer *rpcWriter) error {
	for {
		opcode, payload, err := readFrame(connection)
		if err != nil {
			return err
		}
		switch opcode {
		case 1:
			var message rpcMessage
			if err := json.Unmarshal(payload, &message); err != nil {
				return err
			}
			if message.Event == "ERROR" {
				return errors.New("Discord rejected a presence update")
			}
		case 2:
			return errors.New("Discord closed the IPC connection")
		case 3:
			if err := writer.write(4, json.RawMessage(payload)); err != nil {
				return err
			}
		}
	}
}
