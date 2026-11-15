package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestLlamaServerAnalyzerAcceptsValidReply verifies the image request and the parsed semantic result.
func TestLlamaServerAnalyzerAcceptsValidReply(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" || r.Method != http.MethodPost {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var request struct {
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
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"description\":\"A tree\",\"subjects\":[\"tree\"],\"uncertain\":false}"}}]}`))
	}))
	defer server.Close()
	result, err := (&LlamaServerAnalyzer{endpoint: server.URL, client: server.Client()}).Analyze(context.Background(), "data:image/png;base64,AA==")
	if err != nil || result.Description != "A tree" || len(result.Subjects) != 1 || result.Source == "" {
		t.Fatalf("unexpected semantic result: %+v, %v", result, err)
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
