package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// WeatherProvider fetches weather data from OpenWeatherMap
type WeatherProvider struct {
	APIKey   string
	City     string
	LastData *WeatherData
}

type WeatherData struct {
	Temperature int    `json:"temp"`
	Condition   string `json:"condition"`
	Icon        string `json:"icon"`
}

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
}

func NewWeatherProvider(apiKey, city string) *WeatherProvider {
	return &WeatherProvider{
		APIKey: apiKey,
		City:   city,
	}
}

func (w *WeatherProvider) Fetch() (*WeatherData, error) {
	if w.APIKey == "" || w.APIKey == "YOUR_OWM_API_KEY" {
		return &WeatherData{
			Temperature: 30,
			Condition:   "Clouds",
			Icon:        "☁",
		}, nil
	}

	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s", w.City, w.APIKey)
	resp, err := http.Get(url)
	if err != nil {
		return w.LastData, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return w.LastData, err
	}

	var weatherResp WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return w.LastData, err
	}

	// Return fallback data if the response is invalid
	if weatherResp.Main.Temp == 0 {
		return &WeatherData{
			Temperature: 30,
			Condition:   "Clouds",
			Icon:        "☁",
		}, nil
	}

	icon := "☁"
	condition := "Clouds"
	if len(weatherResp.Weather) > 0 {
		icon = getWeatherIcon(weatherResp.Weather[0].ID)
		condition = weatherResp.Weather[0].Main
	}

	data := &WeatherData{
		Temperature: int(weatherResp.Main.Temp),
		Condition:   condition,
		Icon:        icon,
	}
	w.LastData = data
	return data, nil
}

func getWeatherIcon(id int) string {
	switch {
	case id >= 200 && id < 300:
		return "⛈"
	case id >= 300 && id < 400:
		return "🌧"
	case id >= 500 && id < 600:
		return "🌧"
	case id >= 600 && id < 700:
		return "❄"
	case id >= 700 && id < 800:
		return "🌫"
	case id == 800:
		return "☀"
	case id >= 801 && id < 900:
		return "☁"
	default:
		return "☁"
	}
}

// NewsProvider fetches tech news from multiple sources (Hacker News and Dev.to)
type NewsProvider struct {
	Tags        []string
	CurrentTag  string
	LastData    []NewsItem
	HNClient    *http.Client
	DevToClient *http.Client
}

type NewsItem struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Points      int      `json:"points"`
	Author      string   `json:"author"`
	CreatedAt   int64    `json:"created_at_i"`
	ObjectID    string   `json:"objectID"`
	Source      string   // "hackernews" or "devto"
	Description string   `json:"description"`
	Tags        []string `json:"tag_list"`
}

// Hacker News API response
type HNResponse struct {
	Hits []struct {
		Title     string `json:"title"`
		URL       string `json:"url"`
		Points    int    `json:"points"`
		Author    string `json:"author"`
		CreatedAt int64  `json:"created_at_i"`
		ObjectID  string `json:"objectID"`
	} `json:"hits"`
}

// Dev.to API response
type DevToResponse []struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	User  struct {
		Name string `json:"name"`
	} `json:"user"`
	CreatedAt   string   `json:"published_at"`
	Description string   `json:"description"`
	TagList     []string `json:"tag_list"`
}

func NewNewsProvider(tags []string) *NewsProvider {
	return &NewsProvider{
		Tags:        tags,
		CurrentTag:  "all",
		HNClient:    &http.Client{Timeout: 10 * time.Second},
		DevToClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (n *NewsProvider) SetCurrentTag(tag string) {
	n.CurrentTag = tag
}

func (n *NewsProvider) Fetch() ([]NewsItem, error) {
	var allItems []NewsItem

	// Fetch from Hacker News
	hnItems, err := n.fetchFromHackerNews()
	if err == nil {
		allItems = append(allItems, hnItems...)
	}

	// Fetch from Dev.to
	devToItems, err := n.fetchFromDevTo()
	if err == nil {
		allItems = append(allItems, devToItems...)
	}

	// If we couldn't fetch from either source, return cached data
	if len(allItems) == 0 && len(n.LastData) > 0 {
		return n.LastData, nil
	}

	// Filter by current tag
	filtered := n.filterByCurrentTag(allItems)

	// Limit to 12 items
	if len(filtered) > 12 {
		filtered = filtered[:12]
	}

	n.LastData = filtered
	return filtered, nil
}

func (n *NewsProvider) fetchFromHackerNews() ([]NewsItem, error) {
	query := "story"
	if n.CurrentTag != "all" && n.CurrentTag != "" {
		query = n.CurrentTag
	}

	url := fmt.Sprintf("https://hn.algolia.com/api/v1/search_by_date?tags=story&query=%s&hitsPerPage=15", query)

	resp, err := n.HNClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var hnResp HNResponse
	if err := json.Unmarshal(body, &hnResp); err != nil {
		return nil, err
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

	return items, nil
}

func (n *NewsProvider) fetchFromDevTo() ([]NewsItem, error) {
	url := "https://dev.to/api/articles?per_page=15&top=7"
	if n.CurrentTag != "all" && n.CurrentTag != "" {
		url = fmt.Sprintf("https://dev.to/api/articles?tag=%s&per_page=15&top=7", n.CurrentTag)
	}

	resp, err := n.DevToClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var devToResp DevToResponse
	if err := json.Unmarshal(body, &devToResp); err != nil {
		return nil, err
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

	return items, nil
}

func (n *NewsProvider) filterByCurrentTag(items []NewsItem) []NewsItem {
	if n.CurrentTag == "all" || n.CurrentTag == "" {
		return items
	}

	var filtered []NewsItem
	tagLower := strings.ToLower(n.CurrentTag)

	for _, item := range items {
		// Check title and description for the tag
		titleMatch := strings.Contains(strings.ToLower(item.Title), tagLower)
		descMatch := strings.Contains(strings.ToLower(item.Description), tagLower)

		// Check tags array for Dev.to articles
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

// Scheduler manages widget refresh intervals
type Scheduler struct {
	tasks map[string]*Task
}

type Task struct {
	ID       string
	Interval time.Duration
	LastRun  time.Time
	NextRun  time.Time
	Provider interface{}
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks: make(map[string]*Task),
	}
}

func (s *Scheduler) AddTask(id string, interval time.Duration, provider interface{}) {
	s.tasks[id] = &Task{
		ID:       id,
		Interval: interval,
		LastRun:  time.Now(),
		NextRun:  time.Now().Add(interval),
		Provider: provider,
	}
}

func (s *Scheduler) GetNextTask() *Task {
	var next *Task
	for _, task := range s.tasks {
		if next == nil || task.NextRun.Before(next.NextRun) {
			next = task
		}
	}
	return next
}

func (s *Scheduler) UpdateTask(id string) {
	if task, exists := s.tasks[id]; exists {
		task.LastRun = time.Now()
		task.NextRun = time.Now().Add(task.Interval)
	}
}

func (s *Scheduler) RemoveTask(id string) {
	delete(s.tasks, id)
}

func (s *Scheduler) GetTasks() []*Task {
	tasks := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *Scheduler) GetNextWakeTime() time.Time {
	next := s.GetNextTask()
	if next == nil {
		return time.Now().Add(time.Hour) // Default to 1 hour if no tasks
	}
	return next.NextRun
}

// ParseTTL parses TTL string from config (e.g., "600s", "20s")
func ParseTTL(ttlStr string) time.Duration {
	if ttlStr == "" {
		return 600 * time.Second // Default 10 minutes
	}

	duration, err := time.ParseDuration(ttlStr)
	if err != nil {
		return 600 * time.Second // Default on parse error
	}
	return duration
}
