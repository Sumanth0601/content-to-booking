# Content to Booking

A working proof-of-concept that explores two concrete problems independent service providers face: **turning their existing content into bookable listings**, and **understanding what's actually holding their business back**.

Built in Go with a minimal vanilla JS frontend. Uses a free LLM via [OpenRouter](https://openrouter.ai) — no paid API key needed to run it.

---

## The Problem

Independent service providers (tutors, coaches, trainers, wellness practitioners) live across a fragmented set of tools. They post on Instagram to get discovered, take bookings via Calendly or DMs, get paid on Venmo, and communicate over WhatsApp. Nothing talks to anything else.

Two specific friction points this POC addresses:

**1. Content → Booking gap**
A provider writes a great Instagram caption or bio. That content is completely disconnected from any booking flow. To get booked, they have to manually re-enter all their service details into a separate tool — title, description, category, price, duration, tags. Many don't bother. Revenue gets left on the table.

**2. Provider growth blindness**
Independent providers have no visibility into *why* their business is growing or stagnating. They don't know if their 20% cancellation rate is normal, or if their slow response time is costing them repeat bookings. There's no system that looks at their numbers and tells them specifically what to change.

---

## How It Works

### Module 1 — Content → Service Listing Generator

The user pastes any raw content: an Instagram caption, a Twitter bio, a short description they typed in a DM. The backend sends it to an LLM with a structured extraction prompt that instructs the model to act as a service listing generator and return **only valid JSON** matching a fixed schema:

```json
{
  "title": "...",
  "description": "...",
  "category": "...",
  "duration_minutes": 60,
  "price_usd": 75,
  "tags": ["..."],
  "call_to_action": "..."
}
```

The Go backend parses the JSON response and returns it to the frontend, which renders it as a service card. If the model returns malformed JSON (which free models occasionally do), the handler returns a `422` with a clear error rather than crashing.

### Module 2 — AI Growth Coach

The user enters their performance metrics: sessions completed, average star rating, cancellation rate, average response time, and content post frequency. The backend formats these into a structured prompt and sends them to the LLM with instructions to act as a growth coach — analyze the numbers, find the top 3–5 levers, and return **only valid JSON** matching:

```json
{
  "provider_summary": "...",
  "tips": [
    {
      "priority": 1,
      "title": "...",
      "explanation": "...",
      "action": "..."
    }
  ]
}
```

The key constraint in the prompt is that tips must be **specific** — they must reference the actual numbers from the input. The model shouldn't say "post more content"; it should say "your 2 posts in 30 days is well below average — aim for 3 per week starting this week."

---

## Demo

> Both modules call a live LLM — no mocks or hardcoded responses.

### Content → Service Listing

Paste a social caption and get a structured, ready-to-publish listing:

![Listing result](docs/listing-result.png)

### AI Growth Coach

Enter performance data, get numbered metric-specific recommendations:

![Growth coach form](docs/coach-form.png)

![Growth coach results](docs/coach-result.png)

---

## Setup

### 1. Get a free API key

Sign up at [openrouter.ai](https://openrouter.ai) — no credit card required. OpenRouter provides a unified API to hundreds of models, including a generous free tier.

### 2. Run

```bash
OPENROUTER_API_KEY=sk-or-... go run .
```

Then open [http://localhost:8080](http://localhost:8080).

### 3. (Optional) Swap the model

The default is `openai/gpt-oss-20b:free`. Override with any model slug from [openrouter.ai/models?q=free](https://openrouter.ai/models?q=free):

```bash
OPENROUTER_API_KEY=sk-or-... OPENROUTER_MODEL=meta-llama/llama-3.3-70b-instruct:free go run .
```

---

## Example Inputs

### Generate Listing

**Meditation coach**
> Ready to finally quiet the noise? I offer 1:1 guided meditation sessions designed for stressed-out professionals. 30 or 60 min sessions via Zoom. 5 years teaching mindfulness, 200hr certified. Book your free intro call via the link in bio.

**Language tutor**
> Native Spanish speaker offering conversational Spanish lessons for beginners and intermediate learners. I specialize in making grammar actually make sense. Flexible scheduling, all done on video call. Message me to get started!

**Personal trainer**
> I build custom 8-week programs for people who've tried every plan but can't stay consistent. Online coaching + weekly check-ins. Spots limited — DM "READY" if you want to start this week.

### Growth Coach

**Provider A — High cancellation, low content**
- Category: Fitness, Sessions: 28, Rating: 4.6, Cancellation: 22%, Response time: 15 min, Posts: 2

**Provider B — Great content, slow response**
- Category: Wellness, Sessions: 61, Rating: 4.8, Cancellation: 5%, Response time: 180 min, Posts: 18

**Provider C — New provider, low volume**
- Category: Education, Sessions: 8, Rating: 5.0, Cancellation: 0%, Response time: 10 min, Posts: 1

Each set produces meaningfully different recommendations — the model responds to the actual numbers, not generic advice.

---

## Architecture

```
├── main.go                 # HTTP server, routes, CORS middleware
├── handlers/
│   ├── generate_listing.go # POST /api/generate-listing
│   ├── growth_coach.go     # POST /api/growth-coach
│   └── helpers.go          # Shared JSON error helper
├── llm/
│   └── openai.go           # OpenRouter HTTP client (OpenAI-compatible API)
├── models/
│   ├── listing.go          # ServiceListing struct
│   └── coaching.go         # ProviderMetrics + CoachingReport structs
└── static/
    └── index.html          # Two-tab UI (Tailwind CSS CDN + vanilla JS fetch)
```

The LLM client (`llm/openai.go`) makes a direct HTTP call to the OpenRouter API — no SDK, no abstraction. OpenRouter uses the same request/response format as OpenAI, so switching to a paid model or a different provider is a one-line change.

## API

| Method | Path | Body | Returns |
|--------|------|------|---------|
| `POST` | `/api/generate-listing` | `{ "content": "..." }` | `ServiceListing` JSON |
| `POST` | `/api/growth-coach` | `ProviderMetrics` JSON | `CoachingReport` JSON |
| `GET` | `/health` | — | `{ "status": "ok" }` |
