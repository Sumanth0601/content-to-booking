package models

type ProviderMetrics struct {
	Category            string  `json:"category"`
	SessionsCompleted   int     `json:"sessions_completed"`
	AvgStarRating       float64 `json:"avg_star_rating"`
	CancellationRatePct float64 `json:"cancellation_rate_pct"`
	AvgResponseTimeMin  int     `json:"avg_response_time_minutes"`
	ContentPostsLast30d int     `json:"content_posts_last_30_days"`
}

type CoachingTip struct {
	Priority    int    `json:"priority"`
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Action      string `json:"action"`
}

type CoachingReport struct {
	ProviderSummary string        `json:"provider_summary"`
	Tips            []CoachingTip `json:"tips"`
}
