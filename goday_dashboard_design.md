# GoDay Terminal Dash```
┌────────────────────────────────────────────────────────────────────────────┐
│ Bhanu Reddy  •  Fri 25 Jul 2025 19:45   •  ☁ 30 °C (BLR)   •  R Refresh    │
├────────────────────────────────────────────────────────────────────────────┤d – **At‑a‑Glance** Design (v0.6)

*Author: TBD*\
*Last updated: 25 Jul 2025*

---

## 1 Vision

Deliver a **single‑view, zero‑navigation** terminal dashboard that surfaces everything an engineer needs — tasks, code, builds, messages — **and key personal context** (user name, date, local weather) plus a bite‑sized “Tech News” feed that can be filtered by tags. The entire workspace remains readable at a glance, with optional shortcuts for deeper actions (e.g., Jira work‑log).

**New in v0.6**

- Header bar now shows **user name**, **current date/time**, and **weather icon + temp**.
- Added **Tech News** tile with tag filter (`ai`, `golang`, `security`, …).
- Per‑widget refresh TTLs retained; Weather and Clock tick every 10 min.

---

## 2 UX Snapshot (100‑column reference)

```
┌────────────────────────────────────────────────────────────────────────────┐
│ Bhanu Reddy  •  Fri 25 Jul 2025 19:45   •  ☁ 30 °C (HYD)   •  R Refresh    │
├────────────────────────────────────────────────────────────────────────────┤
│ JIRA (4)   │ PRs (2)     │ Builds (1❌) │ Commits (6) │ Calendar (3)     │
│ • ENG‑421 UI bug          ⏳ 8h  [w]  ↵ opens link                           │
│ • ENG‑389 SSO fix         —     [w]                                         │
│ …                                                               +1 more…    │
├────────────┼──────────────┼────────────┼────────────┼──────────────────────┤
│ Slack (7)  │ Todos (5)   │ Confluence (2) │ PagerDuty (0) │ Tech News (5)  │
│                               ↳ filter: [golang]                           │
└────────────┴──────────────┴────────────┴────────────┴──────────────────────┘
Legend: **[w]** log work; **↵ / click** opens row link; **news tag** toggled with **t**.
```

- **Header bar** is now 1 terminal row but may wrap if width < 90 chars (weather pills drop first).
- **Tech News tile** shows headlines that match the active tag list; press **t** to cycle tags.

---

## 3 Goals / Non‑Goals

| Goals                                                                | Non‑Goals                     |
| -------------------------------------------------------------------- | ----------------------------- |
| One‑screen visibility with larger tiles                              | Multi‑screen dashboards       |
| Per‑widget TTL refresh (Slack 20 s, Confluence 300 s, Weather 600 s) | WebSocket streaming (future)  |
| Click / Enter opens every row link                                   | Inline preview of web pages   |
| Jira work‑log shortcut                                               | Editing Slack or News content |

---

## 4 Functional Requirements

1. **Header info**
   - Display `user.name` from config, `time.Now().Format`, and latest Weather provider value.
   - Weather pill shows emoji icon + integer temperature + city code.
2. **Tech News widget**
   - Reads `tags` list from config (`[ai, devsecops]`).
   - Displays up to 6 headlines whose title OR description contains any tag (case‑insensitive).
   - Press **t** cycles through configured tags; press **T** clears filter to *All*.
   - Each headline row is a link (↵ opens).
3. **Per‑widget TTL** enforced by scheduler; default News 600 s, Weather 600 s.
4. Existing behaviour (clickable rows, Jira work‑log `[w]`, overflow pill) retained.

---

## 5 Non‑Functional

- Header updates time every 60 s independent of other TTLs.
- Weather API errors fall back to “N/A” pill instead of failing build.
- Steady RSS target unchanged (< 45 MB).

---

## 6 Architecture Additions

```
                ┌── ClockProvider (60 s) ──┐
Fetch Hub ──────┤                           ├─► headerModel (name/date/weather)
                ├── WeatherProvider (600 s)┘
                ├── NewsProvider (600 s) ──► TechNewsWidget
                └── Existing providers …
```

- `ClockProvider` is in‑process tick, no network.
- `WeatherProvider` hits OpenWeatherMap `/weather?q={city}&units=metric&appid=…`.
- `NewsProvider` defaults to Hacker News Algolia API (`/search_by_date?tags=story&query=golang`).

---

## 7 YAML Config Snippets

```yaml
user:
  name: "Bhanu Reddy"
  location: "Bengaluru,IN"   # for weather API

ui:
  layout: at_a_glance
  min_width: 100
  tile_height: 7

widgets:
  weather:
    ttl: 600s
    api_key: ${OWM_API_KEY}

  news:
    ttl: 600s
    tags: [golang, security, ai]
    provider: hn        # hn | devto | newsapi

  slack:
    ttl: 20s
  confluence:
    ttl: 300s
  jira:
    ttl: 45s
    log_work: true
```

---

## 8 UI Implementation Notes

- **headerModel** separate Bubble Tea model; gridModel renders it above widgets.
- Weather pill style derived from condition ID → emoji map.
- TechNewsWidget implements `Selectable`, can receive key **t/T** events.
- Scheduler now uses priority queue keyed by next wake time; Clock tick is baked into main `tea.Tick` every second.

---

## 9 Figma Updates

- Add **Header/Bar v2** component: left slot `username`, centre `date/time`, right slot `weather`.
- Create **Pill/Weather** with emoji + temperature style.
- New **Widget/Card‑News** component with tag chip row at top (`Chip/Tag active|inactive`).
- Prototype: **t** key variant toggles active tag overlay.

---

## 10 Risks & Mitigations

| Risk                             | Mitigation                                                              |
| -------------------------------- | ----------------------------------------------------------------------- |
| Weather API quota (60 calls/min) | TTL 600 s and on‑disk cache of last value.                              |
| News headlines very long         | Truncate with `…` but keep full text on hover tooltip (mouse terminals) |

---

*End of v0.6 – adds header with user/date/weather and Tech News widget with tag filter.*

