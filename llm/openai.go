package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const openRouterURL = "https://openrouter.ai/api/v1/chat/completions"

// defaultModel is a free model on OpenRouter. Override with OPENROUTER_MODEL env var.
// Full list of free models: https://openrouter.ai/models?q=free
const defaultModel = "openai/gpt-oss-20b:free"

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Complete sends a system + user prompt to OpenRouter and returns the model reply.
// Set OPENROUTER_MODEL to use a different model (any free or paid model on openrouter.ai).
func Complete(systemPrompt, userPrompt string) (string, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENROUTER_API_KEY environment variable not set")
	}

	model := os.Getenv("OPENROUTER_MODEL")
	if model == "" {
		model = defaultModel
	}

	reqBody := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, openRouterURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("OpenRouter error: %s", chatResp.Error.Message)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenRouter returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from OpenRouter")
	}

	return chatResp.Choices[0].Message.Content, nil
}
