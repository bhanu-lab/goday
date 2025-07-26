package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

// BaseNewsPlugin provides common functionality for news plugins
type BaseNewsPlugin struct {
	id            string
	pluginType    string
	name          string
	version       string
	description   string
	author        string
	tags          []string
	currentTag    string
	supportedTags []string
	client        *http.Client
	lastData      []NewsItem
}

// NewBaseNewsPlugin creates a new base news plugin
func NewBaseNewsPlugin(id, name, version, description, author string) *BaseNewsPlugin {
	return &BaseNewsPlugin{
		id:          id,
		pluginType:  "news",
		name:        name,
		version:     version,
		description: description,
		author:      author,
		tags:        []string{},
		currentTag:  "all",
		client:      &http.Client{Timeout: 10 * time.Second},
		lastData:    []NewsItem{},
	}
}

// GetID returns the plugin ID
func (bnp *BaseNewsPlugin) GetID() string {
	return bnp.id
}

// GetType returns the plugin type
func (bnp *BaseNewsPlugin) GetType() string {
	return bnp.pluginType
}

// GetMetadata returns plugin metadata
func (bnp *BaseNewsPlugin) GetMetadata() PluginMetadata {
	return PluginMetadata{
		Name:        bnp.name,
		Version:     bnp.version,
		Description: bnp.description,
		Author:      bnp.author,
		Type:        bnp.pluginType,
		Config: map[string]string{
			"current_tag":    bnp.currentTag,
			"supported_tags": strings.Join(bnp.supportedTags, ","),
		},
	}
}

// SetTags configures the tags for filtering
func (bnp *BaseNewsPlugin) SetTags(tags []string) {
	bnp.tags = tags
}

// GetCurrentTag returns the current active tag
func (bnp *BaseNewsPlugin) GetCurrentTag() string {
	return bnp.currentTag
}

// SetCurrentTag sets the current active tag
func (bnp *BaseNewsPlugin) SetCurrentTag(tag string) {
	bnp.currentTag = tag
}

// GetSupportedTags returns all supported tags
func (bnp *BaseNewsPlugin) GetSupportedTags() []string {
	return bnp.supportedTags
}

// Cleanup performs cleanup
func (bnp *BaseNewsPlugin) Cleanup() error {
	// Close HTTP client if needed
	return nil
}

// filterByCurrentTag filters news items by the current tag
func (bnp *BaseNewsPlugin) filterByCurrentTag(items []NewsItem) []NewsItem {
	if bnp.currentTag == "all" || bnp.currentTag == "" {
		return items
	}

	var filtered []NewsItem
	tagLower := strings.ToLower(bnp.currentTag)

	for _, item := range items {
		// Check title and description for the tag
		titleMatch := strings.Contains(strings.ToLower(item.Title), tagLower)
		descMatch := strings.Contains(strings.ToLower(item.Description), tagLower)

		// Check tags array
		tagMatch := false
		for _, tag := range item.Tags {
			if strings.ToLower(tag) == tagLower {
				tagMatch = true
				break
			}
		}

		if titleMatch || descMatch || tagMatch {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// HackerNewsPlugin implements news fetching from Hacker News
type HackerNewsPlugin struct {
	*BaseNewsPlugin
}

// NewHackerNewsPlugin creates a new Hacker News plugin
func NewHackerNewsPlugin() *HackerNewsPlugin {
	base := NewBaseNewsPlugin(
		"hackernews",
		"Hacker News",
		"1.0.0",
		"Fetches tech news from Hacker News using the Algolia API",
		"GoDay Team",
	)

	base.supportedTags = []string{"all", "golang", "javascript", "python", "rust", "ai", "security", "startup", "programming"}

	return &HackerNewsPlugin{
		BaseNewsPlugin: base,
	}
}

// Initialize sets up the plugin with configuration
func (hn *HackerNewsPlugin) Initialize(config map[string]interface{}) error {
	// Hacker News doesn't require API keys, so just validate config
	if tags, ok := config["tags"].([]string); ok {
		hn.SetTags(tags)
	}
	if currentTag, ok := config["current_tag"].(string); ok {
		hn.SetCurrentTag(currentTag)
	}
	return nil
}

// Fetch retrieves news from Hacker News
func (hn *HackerNewsPlugin) Fetch(ctx context.Context) (interface{}, error) {
	query := "story"
	if hn.currentTag != "all" && hn.currentTag != "" {
		query = hn.currentTag
	}

	url := fmt.Sprintf("https://hn.algolia.com/api/v1/search_by_date?tags=story&query=%s&hitsPerPage=15", query)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return hn.lastData, err
	}

	resp, err := hn.client.Do(req)
	if err != nil {
		return hn.lastData, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return hn.lastData, err
	}

	var hnResp struct {
		Hits []struct {
			Title     string `json:"title"`
			URL       string `json:"url"`
			Points    int    `json:"points"`
			Author    string `json:"author"`
			CreatedAt int64  `json:"created_at_i"`
			ObjectID  string `json:"objectID"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(body, &hnResp); err != nil {
		return hn.lastData, err
	}

	var items []NewsItem
	for _, hit := range hnResp.Hits {
		if hit.URL == "" || hit.Title == "" {
			continue
		}

		items = append(items, NewsItem{
			Title:     hit.Title,
			URL:       hit.URL,
			Points:    hit.Points,
			Author:    hit.Author,
			CreatedAt: hit.CreatedAt,
			ObjectID:  hit.ObjectID,
			Source:    "hackernews",
		})
	}

	// Filter by current tag
	filtered := hn.filterByCurrentTag(items)

	// Limit to 10 items
	if len(filtered) > 10 {
		filtered = filtered[:10]
	}

	hn.lastData = filtered
	return filtered, nil
}

// DevToPlugin implements news fetching from Dev.to
type DevToPlugin struct {
	*BaseNewsPlugin
}

// NewDevToPlugin creates a new Dev.to plugin
func NewDevToPlugin() *DevToPlugin {
	base := NewBaseNewsPlugin(
		"devto",
		"Dev.to",
		"1.0.0",
		"Fetches developer articles from Dev.to API",
		"GoDay Team",
	)

	base.supportedTags = []string{"all", "golang", "javascript", "python", "react", "webdev", "tutorial", "beginners", "productivity"}

	return &DevToPlugin{
		BaseNewsPlugin: base,
	}
}

// Initialize sets up the plugin with configuration
func (dt *DevToPlugin) Initialize(config map[string]interface{}) error {
	if tags, ok := config["tags"].([]string); ok {
		dt.SetTags(tags)
	}
	if currentTag, ok := config["current_tag"].(string); ok {
		dt.SetCurrentTag(currentTag)
	}
	return nil
}

// Fetch retrieves articles from Dev.to
func (dt *DevToPlugin) Fetch(ctx context.Context) (interface{}, error) {
	url := "https://dev.to/api/articles?per_page=15&top=7"
	if dt.currentTag != "all" && dt.currentTag != "" {
		url = fmt.Sprintf("https://dev.to/api/articles?tag=%s&per_page=15&top=7", dt.currentTag)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return dt.lastData, err
	}

	resp, err := dt.client.Do(req)
	if err != nil {
		return dt.lastData, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dt.lastData, err
	}

	var devToResp []struct {
		Title string `json:"title"`
		URL   string `json:"url"`
		User  struct {
			Name string `json:"name"`
		} `json:"user"`
		CreatedAt   string   `json:"published_at"`
		Description string   `json:"description"`
		TagList     []string `json:"tag_list"`
	}

	if err := json.Unmarshal(body, &devToResp); err != nil {
		return dt.lastData, err
	}

	var items []NewsItem
	for _, article := range devToResp {
		if article.URL == "" || article.Title == "" {
			continue
		}

		items = append(items, NewsItem{
			Title:       article.Title,
			URL:         article.URL,
			Author:      article.User.Name,
			Description: article.Description,
			Tags:        article.TagList,
			Source:      "devto",
		})
	}

	// Filter by current tag
	filtered := dt.filterByCurrentTag(items)

	// Limit to 10 items
	if len(filtered) > 10 {
		filtered = filtered[:10]
	}

	dt.lastData = filtered
	return filtered, nil
}

// AggregateNewsPlugin combines multiple news sources
type AggregateNewsPlugin struct {
	*BaseNewsPlugin
	sources []NewsPlugin
}

// NewAggregateNewsPlugin creates a new aggregate news plugin
func NewAggregateNewsPlugin(sources []NewsPlugin) *AggregateNewsPlugin {
	base := NewBaseNewsPlugin(
		"aggregate-news",
		"Aggregate News",
		"1.0.0",
		"Aggregates news from multiple sources",
		"GoDay Team",
	)

	// Collect all supported tags from sources
	var allTags []string
	tagSet := make(map[string]bool)
	tagSet["all"] = true
	allTags = append(allTags, "all")

	for _, source := range sources {
		for _, tag := range source.GetSupportedTags() {
			if !tagSet[tag] {
				tagSet[tag] = true
				allTags = append(allTags, tag)
			}
		}
	}

	base.supportedTags = allTags

	return &AggregateNewsPlugin{
		BaseNewsPlugin: base,
		sources:        sources,
	}
}

// Initialize sets up the plugin with configuration
func (an *AggregateNewsPlugin) Initialize(config map[string]interface{}) error {
	if tags, ok := config["tags"].([]string); ok {
		an.SetTags(tags)
	}
	if currentTag, ok := config["current_tag"].(string); ok {
		an.SetCurrentTag(currentTag)
	}

	// Initialize all source plugins
	for _, source := range an.sources {
		if err := source.Initialize(config); err != nil {
			return fmt.Errorf("failed to initialize source %s: %w", source.GetID(), err)
		}
	}

	return nil
}

// Fetch retrieves news from all sources and aggregates them
func (an *AggregateNewsPlugin) Fetch(ctx context.Context) (interface{}, error) {
	var allItems []NewsItem

	// Set current tag on all sources
	for _, source := range an.sources {
		source.SetCurrentTag(an.currentTag)

		// Fetch from each source
		data, err := source.Fetch(ctx)
		if err != nil {
			// Log error but continue with other sources
			fmt.Printf("Error fetching from source %s: %v\n", source.GetID(), err)
			continue
		}

		if items, ok := data.([]NewsItem); ok {
			allItems = append(allItems, items...)
		}
	}

	// If we couldn't fetch from any source, return cached data
	if len(allItems) == 0 && len(an.lastData) > 0 {
		return an.lastData, nil
	}

	// Filter by current tag (in case sources didn't filter properly)
	filtered := an.filterByCurrentTag(allItems)

	// Limit to 12 items total (more items for better coverage)
	if len(filtered) > 12 {
		filtered = filtered[:12]
	}

	an.lastData = filtered
	return filtered, nil
}

// SetCurrentTag sets the current tag on the aggregate plugin and all sources
func (an *AggregateNewsPlugin) SetCurrentTag(tag string) {
	an.currentTag = tag
	for _, source := range an.sources {
		source.SetCurrentTag(tag)
	}
}

// HackernoonPlugin implements news fetching from Hackernoon RSS feed
type HackernoonPlugin struct {
	*BaseNewsPlugin
	feedParser *gofeed.Parser
}

// NewHackernoonPlugin creates a new Hackernoon RSS plugin
func NewHackernoonPlugin() *HackernoonPlugin {
	base := NewBaseNewsPlugin(
		"hackernoon",
		"Hackernoon",
		"1.0.0",
		"Fetches tech articles from Hackernoon RSS feed",
		"GoDay Team",
	)

	base.supportedTags = []string{"all", "tech", "programming", "blockchain", "ai", "startup", "cybersecurity", "javascript", "python", "golang"}

	return &HackernoonPlugin{
		BaseNewsPlugin: base,
		feedParser:     gofeed.NewParser(),
	}
}

// Initialize sets up the plugin with configuration
func (hn *HackernoonPlugin) Initialize(config map[string]interface{}) error {
	if tags, ok := config["tags"].([]string); ok {
		hn.SetTags(tags)
	}
	if currentTag, ok := config["current_tag"].(string); ok {
		hn.SetCurrentTag(currentTag)
	}
	return nil
}

// Fetch retrieves articles from Hackernoon RSS feed
func (hn *HackernoonPlugin) Fetch(ctx context.Context) (interface{}, error) {
	// Hackernoon RSS feed URL
	url := "https://hackernoon.com/feed"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return hn.lastData, err
	}

	resp, err := hn.client.Do(req)
	if err != nil {
		return hn.lastData, err
	}
	defer resp.Body.Close()

	feed, err := hn.feedParser.Parse(resp.Body)
	if err != nil {
		return hn.lastData, err
	}

	var items []NewsItem
	for _, item := range feed.Items {
		if item.Link == "" || item.Title == "" {
			continue
		}

		// Extract tags from categories
		var tags []string
		for _, category := range item.Categories {
			tags = append(tags, category)
		}

		// Parse published date
		var createdAt int64
		if item.PublishedParsed != nil {
			createdAt = item.PublishedParsed.Unix()
		}

		// Get author name
		author := "Hackernoon"
		if len(item.Authors) > 0 && item.Authors[0].Name != "" {
			author = item.Authors[0].Name
		}

		items = append(items, NewsItem{
			Title:       item.Title,
			URL:         item.Link,
			Author:      author,
			Description: item.Description,
			Tags:        tags,
			Source:      "hackernoon",
			CreatedAt:   createdAt,
		})

		// Limit to 15 items from RSS
		if len(items) >= 15 {
			break
		}
	}

	// Filter by current tag
	filtered := hn.filterByCurrentTag(items)

	// Limit to 10 items
	if len(filtered) > 10 {
		filtered = filtered[:10]
	}

	hn.lastData = filtered
	return filtered, nil
}
