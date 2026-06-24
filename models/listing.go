package models

type ServiceListing struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	DurationMin int      `json:"duration_minutes"`
	PriceUSD    float64  `json:"price_usd"`
	Tags        []string `json:"tags"`
	CTA         string   `json:"call_to_action"`
}
