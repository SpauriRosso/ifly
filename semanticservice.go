package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// SemanticAnalysis is an unverified local model description, separate from measured visual features.
type SemanticAnalysis struct {
	Description string   `json:"description"`
	Subjects    []string `json:"subjects"`
	Uncertain   bool     `json:"uncertain"`
	Source      string   `json:"source"`
}

// SemanticAnalyzer allows the local inference runtime to be replaced without changing the application service.
type SemanticAnalyzer interface {
	Analyze(context.Context, string) (SemanticAnalysis, error)
}

type SemanticService struct {
	engine SemanticAnalyzer
	mu     sync.Mutex
}

// NewSemanticService constructs a loopback-only llama.cpp adapter for the development prototype.
func NewSemanticService() *SemanticService {
	endpoint := os.Getenv("IFLY_VLM_URL")
	if endpoint == "" {
		endpoint = "http://127.0.0.1:11480"
	}
	return &SemanticService{engine: &LlamaServerAnalyzer{endpoint: endpoint, client: &http.Client{Timeout: 75 * time.Second}}}
}

// Analyze validates a bounded PNG and runs one local inference at a time with a fixed timeout.
func (s *SemanticService) Analyze(dataURL string) (SemanticAnalysis, error) {
	if _, err := decodeAnalysisImage(dataURL); err != nil {
		return SemanticAnalysis{}, err
	}
	if !s.mu.TryLock() {
		return SemanticAnalysis{}, errors.New("semantic analysis already in progress")
	}
	defer s.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 75*time.Second)
	defer cancel()
	return s.engine.Analyze(ctx, dataURL)
}

type LlamaServerAnalyzer struct {
	endpoint string
	client   *http.Client
}

// Analyze sends an image to a local llama.cpp multimodal endpoint and validates its structured reply.
func (a *LlamaServerAnalyzer) Analyze(ctx context.Context, dataURL string) (SemanticAnalysis, error) {
	endpoint, err := url.Parse(a.endpoint)
	if err != nil || endpoint.Scheme != "http" || (endpoint.Hostname() != "127.0.0.1" && endpoint.Hostname() != "localhost") || endpoint.Port() == "" || endpoint.Path != "" {
		return SemanticAnalysis{}, errors.New("VLM endpoint must be a local HTTP address")
	}
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"description": map[string]any{"type": "string"},
			"subjects":    map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"uncertain":   map[string]any{"type": "boolean"},
		},
		"required":             []string{"description", "subjects", "uncertain"},
		"additionalProperties": false,
	}
	requestBody := map[string]any{
		"model": "local-vlm", "temperature": 0, "max_tokens": 160, "stream": false,
		"response_format": map[string]any{"type": "json_object", "schema": schema},
		"messages": []any{map[string]any{"role": "user", "content": []any{
			map[string]any{"type": "text", "text": "Describe this drawing literally. Name only visible subjects and actions. If unclear, mark uncertain. Do not invent context or emotion."},
			map[string]any{"type": "image_url", "image_url": map[string]any{"url": dataURL}},
		}}},
	}
	encoded, err := json.Marshal(requestBody)
	if err != nil {
		return SemanticAnalysis{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.endpoint+"/v1/chat/completions", bytes.NewReader(encoded))
	if err != nil {
		return SemanticAnalysis{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := a.client.Do(req)
	if err != nil {
		return SemanticAnalysis{}, fmt.Errorf("local VLM unavailable: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return SemanticAnalysis{}, fmt.Errorf("local VLM returned HTTP %d", response.StatusCode)
	}
	var envelope struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 64*1024)).Decode(&envelope); err != nil || len(envelope.Choices) == 0 {
		return SemanticAnalysis{}, errors.New("invalid local VLM response")
	}
	var result SemanticAnalysis
	if err := json.Unmarshal([]byte(envelope.Choices[0].Message.Content), &result); err != nil {
		return SemanticAnalysis{}, errors.New("local VLM did not return valid JSON")
	}
	result.Description = strings.TrimSpace(result.Description)
	if result.Description == "" || utf8.RuneCountInString(result.Description) > 300 || len(result.Subjects) > 8 {
		return SemanticAnalysis{}, errors.New("local VLM returned invalid semantic fields")
	}
	for i, subject := range result.Subjects {
		result.Subjects[i] = strings.TrimSpace(subject)
		if result.Subjects[i] == "" || utf8.RuneCountInString(result.Subjects[i]) > 60 {
			return SemanticAnalysis{}, errors.New("local VLM returned invalid subjects")
		}
	}
	result.Source = "local VLM (experimental)"
	return result, nil
}
