package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sumanth/content-to-booking/llm"
	"github.com/sumanth/content-to-booking/models"
)

const coachingSystemPrompt = `You are a growth coach for independent service providers on a human services platform.
Analyze the provider's performance metrics and return 3-5 specific, prioritized, actionable growth recommendations.
Return ONLY valid JSON with no markdown, no code fences, no explanation — just the raw JSON object.
The JSON must match this exact schema:
{
  "provider_summary": string,
  "tips": [
    {
      "priority": number,
      "title": string,
      "explanation": string,
      "action": string
    }
  ]
}
provider_summary should be 1-2 sentences summarizing the provider's current situation.
Each tip's explanation should say WHY it matters. The action should be a concrete next step they can take this week.
Be specific — reference the actual numbers from their metrics. Do NOT give generic advice.`

func GrowthCoach(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var metrics models.ProviderMetrics
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		jsonError(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if metrics.Category == "" {
		jsonError(w, "category field is required", http.StatusBadRequest)
		return
	}

	userMsg := fmt.Sprintf(
		"Category: %s\nSessions completed: %d\nAverage star rating: %.1f\nCancellation rate: %.1f%%\nAverage response time: %d minutes\nContent posts last 30 days: %d",
		metrics.Category,
		metrics.SessionsCompleted,
		metrics.AvgStarRating,
		metrics.CancellationRatePct,
		metrics.AvgResponseTimeMin,
		metrics.ContentPostsLast30d,
	)

	raw, err := llm.Complete(coachingSystemPrompt, userMsg)
	if err != nil {
		jsonError(w, "LLM error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var report models.CoachingReport
	if err := json.Unmarshal([]byte(raw), &report); err != nil {
		jsonError(w, "failed to parse LLM response as JSON", http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
