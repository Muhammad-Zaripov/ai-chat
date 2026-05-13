package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	domain "github.com/udevs/ai-chat/internal/domain/chat"
)

const defaultBaseURL = "https://api.openai.com"

type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		http:    &http.Client{Timeout: 120 * time.Second},
	}
}

type responsesRequest struct {
	Model              string `json:"model"`
	Input              string `json:"input"`
	PreviousResponseID string `json:"previous_response_id,omitempty"`
}

// Minimal shape we care about from the Responses API.
// `output_text` is the convenience field OpenAI emits alongside the
// structured `output` array.
type responsesResponse struct {
	ID         string `json:"id"`
	OutputText string `json:"output_text"`
	Output     []struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func (c *Client) SendMessage(ctx context.Context, model, previousResponseID, userInput string) (domain.AIReply, error) {
	body, err := json.Marshal(responsesRequest{
		Model:              model,
		Input:              userInput,
		PreviousResponseID: previousResponseID,
	})
	if err != nil {
		return domain.AIReply{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/responses", bytes.NewReader(body))
	if err != nil {
		return domain.AIReply{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.AIReply{}, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.AIReply{}, err
	}

	if resp.StatusCode >= 400 {
		// Try to parse a structured OpenAI error; fall back to raw body.
		var parsed responsesResponse
		if json.Unmarshal(raw, &parsed) == nil && parsed.Error != nil {
			return domain.AIReply{}, fmt.Errorf("openai: %s (%s)", parsed.Error.Message, parsed.Error.Type)
		}
		return domain.AIReply{}, fmt.Errorf("openai: http %d: %s", resp.StatusCode, string(raw))
	}

	var parsed responsesResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return domain.AIReply{}, fmt.Errorf("openai: decode response: %w", err)
	}

	text := parsed.OutputText
	if text == "" {
		// Fall back to walking `output[].content[]` for the first text block.
		for _, item := range parsed.Output {
			for _, ct := range item.Content {
				if ct.Type == "output_text" && ct.Text != "" {
					text = ct.Text
					break
				}
			}
			if text != "" {
				break
			}
		}
	}
	if parsed.ID == "" {
		return domain.AIReply{}, errors.New("openai: response missing id")
	}
	return domain.AIReply{ResponseID: parsed.ID, Output: text}, nil
}
