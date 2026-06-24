package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sumanth/content-to-booking/llm"
	"github.com/sumanth/content-to-booking/models"
)

const listingSystemPrompt = `You are a service listing generator for a human services marketplace.
Given raw content from a service provider (a social media caption, bio, or description), extract and return a structured service listing.
Return ONLY valid JSON with no markdown, no code fences, no explanation — just the raw JSON object.
The JSON must match this exact schema:
{
  "title": string,
  "description": string,
  "category": string,
  "duration_minutes": number,
  "price_usd": number,
  "tags": [string],
  "call_to_action": string
}
Make the title concise and professional. The description should be 2-3 sentences. Choose a category from: Wellness, Fitness, Education, Coaching, Music, Art, or Other. Suggest a reasonable session duration and market-rate price.`

func GenerateListing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	body.Content = strings.TrimSpace(body.Content)
	if body.Content == "" {
		jsonError(w, "content field is required", http.StatusBadRequest)
		return
	}

	raw, err := llm.Complete(listingSystemPrompt, body.Content)
	if err != nil {
		jsonError(w, "LLM error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Strip markdown code fences if the model wraps anyway
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var listing models.ServiceListing
	if err := json.Unmarshal([]byte(raw), &listing); err != nil {
		jsonError(w, "failed to parse LLM response as JSON", http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listing)
}
