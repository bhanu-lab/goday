package main

import (
	"fmt"
	"strings"
	"time"
)

// Widget represents a dashboard widget
type Widget struct {
	Title    string
	Count    int
	Items    []WidgetItem
	Selected int
	HasError bool
}

type WidgetItem struct {
	Title      string
	Subtitle   string
	Status     string
	URL        string
	HasWorkLog bool
}

// WidgetManager manages all widgets
type WidgetManager struct {
	Widgets      map[string]*Widget
	NewsTagIndex int
	NewsTags     []string
}

func NewWidgetManager() *WidgetManager {
	return &WidgetManager{
		Widgets:      make(map[string]*Widget),
		NewsTagIndex: 0,
	}
}

func (wm *WidgetManager) InitializeWidgets(cfg *Config) {
	// Initialize all widgets with placeholder data exactly as per design
	wm.Widgets["jira"] = &Widget{
		Title: "JIRA",
		Count: 4,
		Items: []WidgetItem{
			{Title: "ENG-421 UI bug", Subtitle: "⏳ 8h", Status: "[w]", URL: "https://jira.com/ENG-421", HasWorkLog: true},
			{Title: "ENG-389 SSO fix", Subtitle: "—", Status: "[w]", URL: "https://jira.com/ENG-389", HasWorkLog: true},
			{Title: "ENG-456 Performance", Subtitle: "⏳ 4h", Status: "", URL: "https://jira.com/ENG-456"},
			{Title: "ENG-123 Bug fix", Subtitle: "⏳ 2h", Status: "", URL: "https://jira.com/ENG-123"},
		},
	}

	wm.Widgets["prs"] = &Widget{
		Title: "PRs",
		Count: 2,
		Items: []WidgetItem{
			{Title: "Add new feature", Subtitle: "2 reviews", Status: "🟡", URL: "https://github.com/pr/123"},
			{Title: "Fix bug in auth", Subtitle: "1 review", Status: "🟢", URL: "https://github.com/pr/124"},
		},
	}

	wm.Widgets["builds"] = &Widget{
		Title: "Builds",
		Count: 1,
		Items: []WidgetItem{
			{Title: "main branch", Subtitle: "Failed", Status: "❌", URL: "https://ci.com/build/456"},
		},
		HasError: true,
	}

	wm.Widgets["commits"] = &Widget{
		Title: "Commits",
		Count: 6,
		Items: []WidgetItem{
			{Title: "feat: add new API", Subtitle: "2 hours ago", Status: "", URL: "https://github.com/commit/abc123"},
			{Title: "fix: resolve auth issue", Subtitle: "4 hours ago", Status: "", URL: "https://github.com/commit/def456"},
			{Title: "docs: update README", Subtitle: "6 hours ago", Status: "", URL: "https://github.com/commit/ghi789"},
		},
	}

	wm.Widgets["calendar"] = &Widget{
		Title: "Calendar",
		Count: 3,
		Items: []WidgetItem{
			{Title: "Team Standup", Subtitle: "9:00 AM", Status: "", URL: "https://calendar.com/event/1"},
			{Title: "Code Review", Subtitle: "2:00 PM", Status: "", URL: "https://calendar.com/event/2"},
			{Title: "Sprint Planning", Subtitle: "4:00 PM", Status: "", URL: "https://calendar.com/event/3"},
		},
	}

	wm.Widgets["slack"] = &Widget{
		Title: "Slack",
		Count: 7,
		Items: []WidgetItem{
			{Title: "general", Subtitle: "New message", Status: "🔴", URL: "https://slack.com/channel/general"},
			{Title: "dev-team", Subtitle: "3 unread", Status: "🔴", URL: "https://slack.com/channel/dev-team"},
		},
	}

	wm.Widgets["todos"] = &Widget{
		Title: "Todos",
		Count: 5,
		Items: []WidgetItem{
			{Title: "Review PR #123", Subtitle: "High priority", Status: "🔴", URL: ""},
			{Title: "Update docs", Subtitle: "Medium priority", Status: "🟡", URL: ""},
			{Title: "Fix test", Subtitle: "Low priority", Status: "🟢", URL: ""},
		},
	}

	wm.Widgets["confluence"] = &Widget{
		Title: "Confluence",
		Count: 2,
		Items: []WidgetItem{
			{Title: "API Documentation", Subtitle: "Updated 2h ago", Status: "", URL: "https://confluence.com/doc/1"},
			{Title: "Architecture Guide", Subtitle: "Updated 1d ago", Status: "", URL: "https://confluence.com/doc/2"},
		},
	}

	wm.Widgets["pagerduty"] = &Widget{
		Title: "PagerDuty",
		Count: 0,
		Items: []WidgetItem{},
	}

	// Initialize Tech News widget
	if cfg != nil && len(cfg.Widgets.News.Tags) > 0 {
		wm.NewsTags = cfg.Widgets.News.Tags
	} else {
		// Default tags if none configured
		wm.NewsTags = []string{"golang", "security", "ai"}
	}

	wm.Widgets["news"] = &Widget{
		Title: "Tech News",
		Count: 0, // Will be updated when real data is fetched
		Items: []WidgetItem{
			{Title: "Loading news...", Subtitle: "Fetching from HN & Dev.to", Status: "", URL: ""},
		},
	}
}

func (wm *WidgetManager) CycleNewsTag() {
	if len(wm.NewsTags) > 0 {
		wm.NewsTagIndex = (wm.NewsTagIndex + 1) % (len(wm.NewsTags) + 1)
	}
}

func (wm *WidgetManager) GetCurrentNewsTag() string {
	if wm.NewsTagIndex == 0 {
		return "All"
	}
	if wm.NewsTagIndex <= len(wm.NewsTags) {
		return wm.NewsTags[wm.NewsTagIndex-1]
	}
	return "All"
}

// Render functions for the grid layout - EXACTLY as per design document
func (wm *WidgetManager) RenderGrid() string {
	// Create the exact layout as shown in the UX snapshot
	row1 := wm.renderRow1()
	row2 := wm.renderRow2()

	separator := "├────────────┼──────────────┼────────────┼────────────┼──────────────────────┤\n"
	return row1 + separator + row2
}

func (wm *WidgetManager) renderRow1() string {
	// First row: JIRA (4) | PRs (2) | Builds (1❌) | Commits (6) | Calendar (3)
	jira := wm.renderWidgetSimple("jira")
	prs := wm.renderWidgetSimple("prs")
	builds := wm.renderWidgetSimple("builds")
	commits := wm.renderWidgetSimple("commits")
	calendar := wm.renderWidgetSimple("calendar")

	return fmt.Sprintf("│ %s │ %s │ %s │ %s │ %s │\n", jira, prs, builds, commits, calendar)
}

func (wm *WidgetManager) renderRow2() string {
	// Second row: Slack (7) | Todos (5) | Confluence (2) | PagerDuty (0) | Tech News (5)
	slack := wm.renderWidgetSimple("slack")
	todos := wm.renderWidgetSimple("todos")
	confluence := wm.renderWidgetSimple("confluence")
	pagerduty := wm.renderWidgetSimple("pagerduty")
	news := wm.renderWidgetSimple("news")

	return fmt.Sprintf("│ %s │ %s │ %s │ %s │ %s │\n", slack, todos, confluence, pagerduty, news)
}

func (wm *WidgetManager) renderWidgetSimple(widgetName string) string {
	widget, exists := wm.Widgets[widgetName]
	if !exists {
		return "                    "
	}

	// Simple title with count
	title := fmt.Sprintf("%s (%d)", widget.Title, widget.Count)
	if widget.HasError {
		title += "❌"
	}

	// Add items (max 2 for compact view, except news which shows more)
	maxItems := 2
	if widgetName == "news" {
		maxItems = 10 // Show more news items
	}

	var content []string
	for i, item := range widget.Items {
		if i >= maxItems {
			break
		}

		// Simple item text
		itemText := fmt.Sprintf("• %s", item.Title)
		if item.Subtitle != "" {
			itemText += fmt.Sprintf(" %s", item.Subtitle)
		}
		if item.Status != "" {
			itemText += fmt.Sprintf(" %s", item.Status)
		}

		// Add clickable indicator for items with URLs
		if item.URL != "" {
			itemText += " ↵"
		}

		content = append(content, itemText)
	}

	// Add overflow indicator
	if len(widget.Items) > maxItems {
		content = append(content, fmt.Sprintf("+%d more…", len(widget.Items)-maxItems))
	}

	// Special handling for Tech News - add tag filter
	if widgetName == "news" && len(content) < maxItems {
		content = append(content, fmt.Sprintf("↳ filter: [%s]", wm.GetCurrentNewsTag()))
	}

	// Combine title and content
	result := title
	if len(content) > 0 {
		result += "\n" + strings.Join(content, "\n")
	}

	// Pad to ensure consistent width (20 chars)
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		if len(line) < 20 {
			lines[i] = line + strings.Repeat(" ", 20-len(line))
		} else if len(line) > 20 {
			lines[i] = line[:17] + "..."
		}
	}

	return strings.Join(lines, "\n")
}

// GetClickableItems returns items that have URLs and can be clicked
func (wm *WidgetManager) GetClickableItems() []WidgetItem {
	var clickable []WidgetItem
	for _, widget := range wm.Widgets {
		for _, item := range widget.Items {
			if item.URL != "" {
				clickable = append(clickable, item)
			}
		}
	}
	return clickable
}

// OpenURL opens a URL in the default browser (placeholder for now)
func OpenURL(url string) error {
	// In a real implementation, this would use os/exec to open the URL
	// For now, just return nil as a placeholder
	return nil
}

// UpdateGitCommitsWidget updates the commits widget with data from Git plugin
func (wm *WidgetManager) UpdateGitCommitsWidget(commits []GitCommit) {
	var items []WidgetItem

	for _, commit := range commits {
		// Format the time as relative time
		timeAgo := formatTimeAgo(commit.Date)

		items = append(items, WidgetItem{
			Title:    commit.Message,
			Subtitle: fmt.Sprintf("%s • %s", timeAgo, commit.Repository),
			Status:   "",
			URL:      "", // Could be enhanced with GitHub URL if available
		})
	}

	if wm.Widgets["commits"] != nil {
		wm.Widgets["commits"].Items = items
		wm.Widgets["commits"].Count = len(items)
	}
}

// UpdateGitHubPRsWidget updates the PRs widget with data from GitHub API
func (wm *WidgetManager) UpdateGitHubPRsWidget(prs []GitPullRequest) {
	var items []WidgetItem

	for _, pr := range prs {
		// Format status based on PR state and draft status
		status := "🟢" // open
		if pr.IsDraft {
			status = "🟡" // draft
		}
		if pr.State == "closed" {
			status = "🔴" // closed
		}

		// Format subtitle with repository and update time
		timeAgo := formatTimeAgo(pr.UpdatedAt)
		subtitle := fmt.Sprintf("%s • %s", pr.Repository, timeAgo)

		items = append(items, WidgetItem{
			Title:    pr.Title,
			Subtitle: subtitle,
			Status:   status,
			URL:      pr.URL,
		})
	}

	if wm.Widgets["prs"] != nil {
		wm.Widgets["prs"].Items = items
		wm.Widgets["prs"].Count = len(items)
	}
}

// formatTimeAgo formats a time as a relative time string
func formatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else {
		return t.Format("Jan 2")
	}
}
