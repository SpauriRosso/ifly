package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"changeme/internal/analysis"
	"github.com/wailsapp/wails/v3/pkg/application"
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
	engine  SemanticAnalyzer
	mu      sync.Mutex
	busy    bool
	cancel  context.CancelFunc
	runtime *LocalRuntime
}

// NewSemanticService constructs a loopback-only llama.cpp adapter for the development prototype.
func NewSemanticService() *SemanticService {
	endpoint := os.Getenv("IFLY_VLM_URL")
	var runtime *LocalRuntime
	if endpoint == "" {
		endpoint = "http://127.0.0.1:11480"
		runtime = &LocalRuntime{}
	}
	return &SemanticService{engine: &LlamaServerAnalyzer{endpoint: endpoint, client: &http.Client{Timeout: 75 * time.Second}}, runtime: runtime}
}

// ServiceStartup starts an installed local CPU runtime while keeping structural analysis usable if it is absent.
func (s *SemanticService) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	if s.runtime != nil {
		if err := s.runtime.Start(ctx); err != nil {
			log.Printf("local VLM startup unavailable: %v", err)
		}
	}
	return nil
}

// ServiceShutdown cancels inference and stops only a runtime owned by this application.
func (s *SemanticService) ServiceShutdown() error {
	s.Cancel()
	if s.runtime != nil {
		return s.runtime.Stop()
	}
	return nil
}

// Analyze validates a bounded PNG and runs one local inference at a time with a fixed timeout.
func (s *SemanticService) Analyze(dataURL string) (SemanticAnalysis, error) {
	if _, err := analysis.DecodeImage(dataURL); err != nil {
		return SemanticAnalysis{}, err
	}
	s.mu.Lock()
	if s.busy {
		s.mu.Unlock()
		return SemanticAnalysis{}, errors.New("semantic analysis already in progress")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 75*time.Second)
	s.busy = true
	s.cancel = cancel
	s.mu.Unlock()
	defer func() {
		cancel()
		s.mu.Lock()
		s.busy = false
		s.cancel = nil
		s.mu.Unlock()
	}()
	return s.engine.Analyze(ctx, dataURL)
}

// Cancel stops the current inference request without changing previously completed results.
func (s *SemanticService) Cancel() {
	s.mu.Lock()
	cancel := s.cancel
	s.mu.Unlock()
	if cancel != nil {
		cancel()
	}
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
			"description": map[string]any{"type": "string", "minLength": 1, "maxLength": 140},
			"subjects":    map[string]any{"type": "array", "maxItems": 3, "items": map[string]any{"type": "string", "maxLength": 24}},
			"uncertain":   map[string]any{"type": "boolean"},
		},
		"required":             []string{"description", "subjects", "uncertain"},
		"additionalProperties": false,
	}
	requestBody := map[string]any{
		"model": "local-vlm", "temperature": 0, "max_tokens": 192, "stream": false,
		"response_format": map[string]any{"type": "json_object", "schema": schema},
		"messages": []any{map[string]any{"role": "user", "content": []any{
			map[string]any{"type": "text", "text": "Write one short factual sentence about the visible image. List at most three visible subjects. Use an empty subjects list for abstract patterns. Return only JSON."},
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
	if result.Description == "" || utf8.RuneCountInString(result.Description) > 140 || len(result.Subjects) > 3 {
		return SemanticAnalysis{}, errors.New("local VLM returned invalid semantic fields")
	}
	result.Description = normalizeSemanticDescription(result.Description)
	for i, subject := range result.Subjects {
		result.Subjects[i] = strings.TrimSpace(subject)
		if result.Subjects[i] == "" || utf8.RuneCountInString(result.Subjects[i]) > 24 {
			return SemanticAnalysis{}, errors.New("local VLM returned invalid subjects")
		}
	}
	result.Source = "local VLM (experimental)"
	return result, nil
}

// normalizeSemanticDescription takes a bounded model caption and returns one readable sentence, removing any unfinished tail.
func normalizeSemanticDescription(description string) string {
	if end := strings.IndexAny(description, ".!?"); end >= 0 {
		return strings.TrimSpace(description[:end+1])
	}
	if utf8.RuneCountInString(description) >= 135 || strings.HasSuffix(description, "-") {
		if end := strings.LastIndex(description, ","); end > 0 {
			description = description[:end]
		}
	}
	return strings.TrimRight(strings.TrimSpace(description), " ,-:;") + "."
}
