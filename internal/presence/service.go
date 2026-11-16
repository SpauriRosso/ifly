package presence

import (
	"io"
	"log"
	"strconv"
	"sync"
	"time"
)

// Service keeps one optional Discord activity in sync with the artwork currently visible in IFLY.
type Service struct {
	appID       string
	mu          sync.Mutex
	artworkOpen bool
	connection  io.ReadWriteCloser
	changed     chan struct{}
	stop        chan struct{}
	stopped     chan struct{}
	closeOnce   sync.Once
}

// NewService takes a Discord application ID, then starts a reconnecting local presence session when the ID is valid.
func NewService(appID string) *Service {
	service := &Service{appID: appID, changed: make(chan struct{}, 1), stop: make(chan struct{}), stopped: make(chan struct{})}
	parsedID, err := strconv.ParseUint(appID, 10, 64)
	if appID == "" || err != nil || parsedID == 0 {
		if appID != "" {
			log.Print("Discord presence disabled: IFLY_DISCORD_APP_ID must be a numeric application ID")
		}
		close(service.stopped)
		return service
	}
	go service.run()
	return service
}

// SetArtworkOpen takes whether an image is visible and requests the corresponding Drawing or Admiring Discord status.
func (service *Service) SetArtworkOpen(open bool) {
	service.mu.Lock()
	if service.artworkOpen == open {
		service.mu.Unlock()
		return
	}
	service.artworkOpen = open
	service.mu.Unlock()
	select {
	case service.changed <- struct{}{}:
	default:
	}
}

// Close stops reconnection, closes the active Discord pipe, and waits for the background session to exit.
func (service *Service) Close() {
	service.closeOnce.Do(func() {
		close(service.stop)
		service.mu.Lock()
		connection := service.connection
		service.mu.Unlock()
		if connection != nil {
			_ = connection.Close()
		}
		<-service.stopped
	})
}

// run keeps retrying the local Discord pipe until IFLY shuts down, including when Discord launches later.
func (service *Service) run() {
	defer close(service.stopped)
	for {
		select {
		case <-service.stop:
			return
		default:
		}
		connection, err := dialDiscord()
		if err == nil {
			if err = service.serve(connection); err != nil {
				select {
				case <-service.stop:
					return
				default:
					log.Printf("Discord presence disconnected: %v", err)
				}
			}
		}
		select {
		case <-service.stop:
			return
		case <-time.After(15 * time.Second):
		}
	}
}

// serve takes one Discord pipe, completes its handshake, and sends state changes until the pipe or application closes.
func (service *Service) serve(connection io.ReadWriteCloser) error {
	service.mu.Lock()
	select {
	case <-service.stop:
		service.mu.Unlock()
		return connection.Close()
	default:
	}
	service.connection = connection
	service.mu.Unlock()
	defer func() {
		service.mu.Lock()
		service.connection = nil
		service.mu.Unlock()
		_ = connection.Close()
	}()
	writer := &rpcWriter{connection: connection}
	if deadline, ok := connection.(interface{ SetReadDeadline(time.Time) error }); ok {
		_ = deadline.SetReadDeadline(time.Now().Add(8 * time.Second))
	}
	if err := handshake(connection, writer, service.appID); err != nil {
		return err
	}
	if deadline, ok := connection.(interface{ SetReadDeadline(time.Time) error }); ok {
		_ = deadline.SetReadDeadline(time.Time{})
	}
	service.mu.Lock()
	open := service.artworkOpen
	service.mu.Unlock()
	if err := setActivity(writer, open); err != nil {
		return err
	}
	readError := make(chan error, 1)
	go func() { readError <- readEvents(connection, writer) }()
	refresh := time.NewTicker(30 * time.Second)
	defer refresh.Stop()
	for {
		select {
		case <-service.changed:
		case <-refresh.C:
		case <-service.stop:
			return nil
		case err := <-readError:
			return err
		}
		service.mu.Lock()
		open = service.artworkOpen
		service.mu.Unlock()
		if err := setActivity(writer, open); err != nil {
			return err
		}
	}
}
