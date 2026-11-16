package semantic

import (
	"context"
	"encoding/json"
	"image"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestLlamaServerAnalyzerAcceptsValidReply verifies the image request and the parsed semantic result.
func TestLlamaServerAnalyzerAcceptsValidReply(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" || r.Method != http.MethodPost {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var request struct {
			MaxTokens      int `json:"max_tokens"`
			ResponseFormat struct {
				Schema struct {
					Properties struct {
						Description struct {
							MaxLength int `json:"maxLength"`
						} `json:"description"`
						Atmosphere struct {
							MaxLength int `json:"maxLength"`
						} `json:"atmosphere"`
						Subjects struct {
							MaxItems int `json:"maxItems"`
						} `json:"subjects"`
					} `json:"properties"`
				} `json:"schema"`
			} `json:"response_format"`
			Messages []struct {
				Content []struct {
					Type     string `json:"type"`
					ImageURL struct {
						URL string `json:"url"`
					} `json:"image_url"`
				} `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
		}
		if len(request.Messages) != 1 || len(request.Messages[0].Content) != 2 || request.Messages[0].Content[1].Type != "image_url" || !strings.HasPrefix(request.Messages[0].Content[1].ImageURL.URL, "data:image/png;base64,") {
			t.Error("request did not contain the image")
		}
		if request.MaxTokens < 200 || request.ResponseFormat.Schema.Properties.Description.MaxLength == 0 || request.ResponseFormat.Schema.Properties.Atmosphere.MaxLength == 0 || request.ResponseFormat.Schema.Properties.Subjects.MaxItems == 0 {
			t.Error("request did not bound the structured model response")
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"description\":\"A tree\",\"atmosphere\":\"soft morning light\",\"subjects\":[\"tree\"],\"uncertain\":false}"}}]}`))
	}))
	defer server.Close()
	result, err := (&LlamaServerAnalyzer{endpoint: server.URL, client: server.Client()}).Analyze(context.Background(), "data:image/png;base64,AA==")
	if err != nil || result.Description != "A tree." || result.Atmosphere != "soft morning light" || len(result.Subjects) != 1 || result.Source == "" {
		t.Fatalf("unexpected semantic result: %+v, %v", result, err)
	}
}

// TestNormalizeSemanticDescription checks that a complete first sentence survives and a clipped caption loses its unfinished tail.
func TestNormalizeSemanticDescription(t *testing.T) {
	if got := normalizeSemanticDescription("A loading icon is centered. The loading 3"); got != "A loading icon is centered." {
		t.Fatalf("unexpected first sentence: %q", got)
	}
	clipped := "A landscape at sunset with a deer standing in the foreground, surrounded by mountains and trees under a gradient-"
	if got := normalizeSemanticDescription(clipped); got != "A landscape at sunset with a deer standing in the foreground." {
		t.Fatalf("unexpected clipped sentence: %q", got)
	}
}

type blockingSemanticAnalyzer struct{ started chan struct{} }

// Analyze waits for cancellation so the service test can verify active request cleanup.
func (a *blockingSemanticAnalyzer) Analyze(ctx context.Context, _ string) (SemanticAnalysis, error) {
	close(a.started)
	<-ctx.Done()
	return SemanticAnalysis{}, ctx.Err()
}

// TestSemanticServiceCancel verifies that changing artwork can stop an active inference request.
func TestSemanticServiceCancel(t *testing.T) {
	engine := &blockingSemanticAnalyzer{started: make(chan struct{})}
	service := &SemanticService{engine: engine}
	dataURL := encodeTestImage(t, image.NewRGBA(image.Rect(0, 0, 2, 2)))
	finished := make(chan error, 1)
	go func() { _, err := service.Analyze(dataURL); finished <- err }()
	select {
	case <-engine.started:
	case <-time.After(time.Second):
		t.Fatal("inference did not start")
	}
	service.Cancel()
	select {
	case err := <-finished:
		if err != context.Canceled {
			t.Fatalf("expected cancellation, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("inference did not stop")
	}
}

// TestLlamaServerAnalyzerRejectsInvalidReply ensures malformed model output does not become a valid interpretation.
func TestLlamaServerAnalyzerRejectsInvalidReply(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"some prose"}}]}`))
	}))
	defer server.Close()
	if _, err := (&LlamaServerAnalyzer{endpoint: server.URL, client: server.Client()}).Analyze(context.Background(), "data:image/png;base64,AA=="); err == nil {
		t.Fatal("expected invalid model output to be rejected")
	}
}

// TestLlamaServerAnalyzerRejectsRemoteEndpoint prevents sending private artwork to a nonlocal service.
func TestLlamaServerAnalyzerRejectsRemoteEndpoint(t *testing.T) {
	if _, err := (&LlamaServerAnalyzer{endpoint: "https://example.com:443", client: http.DefaultClient}).Analyze(context.Background(), "data:image/png;base64,AA=="); err == nil {
		t.Fatal("expected remote endpoint to be rejected")
	}
}
