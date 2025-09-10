package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fzxs8/duolasdk"
	"github.com/fzxs8/duolasdk/core"
	"github.com/google/uuid"
)

// Prompt represents a saved prompt
type Prompt struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Content     string   `json:"content"`
	Description string   `json:"description"`
	CreatedAt   int64    `json:"createdAt"`
	UpdatedAt   int64    `json:"updatedAt"`
	Models      []string `json:"models"`
	Version     int      `json:"version"`
	Tags        []string `json:"tags"`
	CreatedBy   string   `json:"createdBy"`
}

// Server represents a configured Ollama server
type Server struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	BaseURL    string `json:"base_url"`
	APIKey     string `json:"api_key"`
	IsActive   bool   `json:"is_active"`
	TestStatus string `json:"test_status"`
	Type       string `json:"type"`
}

// PromptPilot manages the prompt engineering functionality
type PromptPilot struct {
	ctx        context.Context
	store      *duolasdk.AppStore
	httpClient *core.HttpCli
	logger     *core.AppLog
}

// NewPromptPilot creates a new PromptPilot instance
func NewPromptPilot(store *duolasdk.AppStore) *PromptPilot {
	logger := core.NewLogger(&core.LoggerOption{
		Type:   "console",
		Level:  "debug",
		Prefix: "PromptPilot",
	})

	return &PromptPilot{
		store:  store,
		logger: logger,
	}
}

// Startup is called when the application starts
func (p *PromptPilot) Startup(ctx context.Context) {
	p.ctx = ctx
	p.logger.Info("PromptPilot started")
}

// SetContext sets the context
func (p *PromptPilot) SetContext(ctx context.Context) {
	p.ctx = ctx
}

// SetHTTPClient sets the HTTP client
func (p *PromptPilot) SetHTTPClient(client *core.HttpCli) {
	p.httpClient = client
}

// GetServers returns the list of configured servers
func (p *PromptPilot) GetServers() ([]Server, error) {
	p.logger.Debug("Getting configured servers")

	// Return sample servers for demonstration
	servers := []Server{
		{
			ID:         "local",
			Name:       "本地Ollama服务",
			BaseURL:    "http://localhost:11434",
			APIKey:     "",
			IsActive:   true,
			TestStatus: "success",
			Type:       "local",
		},
		{
			ID:         "remote1",
			Name:       "远程开发环境",
			BaseURL:    "http://dev.example.com:11434",
			APIKey:     "sk-xxx",
			IsActive:   false,
			TestStatus: "unknown",
			Type:       "remote",
		},
	}

	p.logger.Debug("Returning %d servers", len(servers))
	return servers, nil
}

// GetModelsByServer returns the list of models available on a specific server
func (p *PromptPilot) GetModelsByServer(serverID string) ([]Model, error) {
	p.logger.Debug("Getting models for server: %s", serverID)

	// Simulate API call delay
	time.Sleep(800 * time.Millisecond)

	var models []Model

	// Return different models based on the selected server
	switch serverID {
	case "local":
		models = []Model{
			{
				Name:       "llama2",
				Model:      "llama2",
				ModifiedAt: "2023-10-10T10:00:00Z",
				Size:       3790000000,
				Digest:     "sha256:xxx",
				Details:    map[string]interface{}{},
			},
			{
				Name:       "mistral",
				Model:      "mistral",
				ModifiedAt: "2023-11-15T14:30:00Z",
				Size:       4620000000,
				Digest:     "sha256:yyy",
				Details:    map[string]interface{}{},
			},
			{
				Name:       "neural-chat",
				Model:      "neural-chat",
				ModifiedAt: "2023-12-01T09:15:00Z",
				Size:       5200000000,
				Digest:     "sha256:zzz",
				Details:    map[string]interface{}{},
			},
			{
				Name:       "codellama",
				Model:      "codellama",
				ModifiedAt: "2023-09-20T16:45:00Z",
				Size:       3790000000,
				Digest:     "sha256:aaa",
				Details:    map[string]interface{}{},
			},
			{
				Name:       "llama2-uncensored",
				Model:      "llama2-uncensored",
				ModifiedAt: "2023-10-12T11:20:00Z",
				Size:       3790000000,
				Digest:     "sha256:bbb",
				Details:    map[string]interface{}{},
			},
		}
	case "remote1":
		models = []Model{
			{
				Name:       "llama2:13b",
				Model:      "llama2:13b",
				ModifiedAt: "2023-10-10T10:00:00Z",
				Size:       7600000000,
				Digest:     "sha256:ccc",
				Details:    map[string]interface{}{},
			},
			{
				Name:       "mistral:7b",
				Model:      "mistral:7b",
				ModifiedAt: "2023-11-15T14:30:00Z",
				Size:       4620000000,
				Digest:     "sha256:ddd",
				Details:    map[string]interface{}{},
			},
			{
				Name:       "mixtral:8x7b",
				Model:      "mixtral:8x7b",
				ModifiedAt: "2023-12-01T09:15:00Z",
				Size:       18500000000,
				Digest:     "sha256:eee",
				Details:    map[string]interface{}{},
			},
		}
	default:
		p.logger.Warn("Unknown server ID: %s", serverID)
		return nil, fmt.Errorf("unknown server ID: %s", serverID)
	}

	p.logger.Debug("Returning %d models for server %s", len(models), serverID)
	return models, nil
}

// TestServerConnection tests the connection to a server
func (p *PromptPilot) TestServerConnection(server Server) (bool, error) {
	p.logger.Debug("Testing connection to server: %s (%s)", server.Name, server.BaseURL)

	// Simulate connection test
	time.Sleep(500 * time.Millisecond)

	// For demonstration, assume local server is always available
	if server.Type == "local" {
		p.logger.Debug("Local server is always available")
		return true, nil
	}

	// For remote servers, simulate a random result
	result := server.TestStatus == "success"
	p.logger.Debug("Connection test result for %s: %t", server.Name, result)

	return result, nil
}

// GeneratePrompt generates a prompt based on user input and selected models
func (p *PromptPilot) GeneratePrompt(idea string, models []string, serverId string) (map[string]string, error) {
	p.logger.Debug("Generating prompt for idea: %s, models: %v, server: %s", idea, models, serverId)

	if idea == "" {
		return nil, fmt.Errorf("idea cannot be empty")
	}

	if len(models) == 0 {
		return nil, fmt.Errorf("at least one model must be selected")
	}

	if len(models) > 3 {
		return nil, fmt.Errorf("maximum 3 models allowed for comparison")
	}

	result := make(map[string]string)

	// Generate a sample prompt for each model
	for _, model := range models {
		samplePrompt := fmt.Sprintf("你是一个%s模型，请根据以下要求生成内容：\n\n%s\n\n请以专业且易懂的方式回应，确保内容准确且有深度。", model, idea)
		result[model] = samplePrompt
	}

	p.logger.Debug("Generated prompts for %d models", len(models))

	return result, nil
}

// GeneratePromptStream generates a prompt with streaming output
func (p *PromptPilot) GeneratePromptStream(idea string, model string, serverId string, callback func(string)) error {
	p.logger.Debug("Generating prompt stream for idea: %s, model: %s, server: %s", idea, model, serverId)

	if idea == "" {
		return fmt.Errorf("idea cannot be empty")
	}

	if model == "" {
		return fmt.Errorf("model must be specified")
	}

	// Sample prompt for streaming
	samplePrompt := fmt.Sprintf("你是一个%s模型，请根据以下要求生成内容：\n\n%s\n\n请以专业且易懂的方式回应，确保内容准确且有深度。", model, idea)

	// Simulate streaming by sending characters one by one
	for i := 0; i < len(samplePrompt); i++ {
		callback(string(samplePrompt[i]))

		// Simulate network delay
		time.Sleep(10 * time.Millisecond)
	}

	p.logger.Debug("Finished streaming prompt generation")
	return nil
}

// OptimizePrompt optimizes a prompt based on user feedback
func (p *PromptPilot) OptimizePrompt(content string, feedback string, model string, serverId string) (string, error) {
	p.logger.Debug("Optimizing prompt with feedback: %s, model: %s, server: %s", feedback, model, serverId)

	if content == "" {
		return "", fmt.Errorf("content cannot be empty")
	}

	if feedback == "" {
		return "", fmt.Errorf("feedback cannot be empty")
	}

	// Simulate optimization
	optimizedPrompt := fmt.Sprintf("%s\n\n优化建议：%s\n\n请根据以上建议优化Prompt内容。", content, feedback)

	p.logger.Debug("Prompt optimization completed")

	return optimizedPrompt, nil
}

// SavePrompt saves a prompt to storage
func (p *PromptPilot) SavePrompt(prompt Prompt) error {
	p.logger.Debug("Saving prompt: %s", prompt.Name)

	// If it's a new prompt, generate ID and timestamps
	if prompt.ID == "" {
		prompt.ID = uuid.New().String()
		now := time.Now().UnixNano() / int64(time.Millisecond)
		prompt.CreatedAt = now
		prompt.UpdatedAt = now
		prompt.Version = 1
	} else {
		// Update timestamp for existing prompt
		prompt.UpdatedAt = time.Now().UnixNano() / int64(time.Millisecond)
		prompt.Version++
	}

	// Convert prompt to JSON for storage
	promptJSON, err := json.Marshal(prompt)
	if err != nil {
		p.logger.Error("Failed to marshal prompt: %v", err)
		return err
	}

	// Save to storage using hash storage
	err = p.store.HSet("prompts", prompt.ID, string(promptJSON))
	if err != nil {
		p.logger.Error("Failed to save prompt to storage: %v", err)
		return err
	}

	p.logger.Debug("Prompt saved successfully: %s", prompt.ID)
	return nil
}

// ListPrompts returns all saved prompts
func (p *PromptPilot) ListPrompts() ([]Prompt, error) {
	p.logger.Debug("Listing saved prompts")

	// Get all prompts from storage
	promptsMap, err := p.store.HGetAll("prompts")
	if err != nil {
		p.logger.Error("Failed to get prompts from storage: %v", err)
		return nil, err
	}

	// Parse prompts
	prompts := make([]Prompt, 0, len(promptsMap))
	for _, promptJSON := range promptsMap {
		var prompt Prompt
		if err := json.Unmarshal([]byte(promptJSON), &prompt); err != nil {
			p.logger.Warn("Failed to parse prompt: %v", err)
			continue
		}
		prompts = append(prompts, prompt)
	}

	p.logger.Debug("Returning %d prompts", len(prompts))
	return prompts, nil
}

// GetPrompt returns a specific prompt by ID
func (p *PromptPilot) GetPrompt(id string) (Prompt, error) {
	p.logger.Debug("Getting prompt: %s", id)

	// Get prompt from storage
	promptJSON, err := p.store.HGet("prompts", id)
	if err != nil {
		p.logger.Error("Failed to get prompt from storage: %v", err)
		return Prompt{}, err
	}

	// Parse prompt
	var prompt Prompt
	if err := json.Unmarshal([]byte(promptJSON), &prompt); err != nil {
		p.logger.Error("Failed to parse prompt: %v", err)
		return Prompt{}, err
	}

	p.logger.Debug("Returning prompt: %s", prompt.Name)
	return prompt, nil
}

// DeletePrompt deletes a prompt by ID
func (p *PromptPilot) DeletePrompt(id string) error {
	p.logger.Debug("Deleting prompt: %s", id)

	// Delete prompt from storage
	err := p.store.HDel("prompts", id)
	if err != nil {
		p.logger.Error("Failed to delete prompt from storage: %v", err)
		return err
	}

	p.logger.Debug("Prompt deleted successfully: %s", id)
	return nil
}
