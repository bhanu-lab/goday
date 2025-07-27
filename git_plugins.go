package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// GitCommit represents a Git commit
type GitCommit struct {
	Hash       string    `json:"hash"`
	Message    string    `json:"message"`
	Author     string    `json:"author"`
	Date       time.Time `json:"date"`
	Repository string    `json:"repository"`
}

// GitPullRequest represents a GitHub Pull Request
type GitPullRequest struct {
	Number     int       `json:"number"`
	Title      string    `json:"title"`
	State      string    `json:"state"`
	Author     string    `json:"author"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Repository string    `json:"repository"`
	URL        string    `json:"url"`
	IsDraft    bool      `json:"draft"`
	Mergeable  *bool     `json:"mergeable"`
}

// LocalGitCommitsPlugin fetches commits from local Git repositories
type LocalGitCommitsPlugin struct {
	id           string
	pluginType   string
	name         string
	version      string
	description  string
	author       string
	gitUser      string
	gitEmail     string
	repositories []string
	client       *http.Client
	lastData     []GitCommit
}

// NewLocalGitCommitsPlugin creates a new local Git commits plugin
func NewLocalGitCommitsPlugin() *LocalGitCommitsPlugin {
	// Get Git user configuration
	gitUser := getGitConfig("user.name")
	gitEmail := getGitConfig("user.email")

	return &LocalGitCommitsPlugin{
		id:          "local-git-commits",
		pluginType:  "git",
		name:        "Local Git Commits",
		version:     "1.0.0",
		description: "Fetches recent commits from local Git repositories",
		author:      "GoDay Team",
		gitUser:     gitUser,
		gitEmail:    gitEmail,
		client:      &http.Client{Timeout: 10 * time.Second},
		lastData:    []GitCommit{},
	}
}

// getGitConfig retrieves Git configuration value
func getGitConfig(key string) string {
	cmd := exec.Command("git", "config", "--global", key)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// GetID returns the plugin ID
func (lgc *LocalGitCommitsPlugin) GetID() string {
	return lgc.id
}

// GetType returns the plugin type
func (lgc *LocalGitCommitsPlugin) GetType() string {
	return lgc.pluginType
}

// GetMetadata returns plugin metadata
func (lgc *LocalGitCommitsPlugin) GetMetadata() PluginMetadata {
	return PluginMetadata{
		Name:        lgc.name,
		Version:     lgc.version,
		Description: lgc.description,
		Author:      lgc.author,
		Type:        lgc.pluginType,
		Config: map[string]string{
			"git_user":  lgc.gitUser,
			"git_email": lgc.gitEmail,
		},
	}
}

// Initialize sets up the plugin with configuration
func (lgc *LocalGitCommitsPlugin) Initialize(config map[string]interface{}) error {
	if repos, ok := config["repositories"].([]string); ok {
		lgc.repositories = repos
	} else {
		// Default to current directory and common dev locations
		lgc.repositories = []string{
			".",
			"~/Development",
			"~/Projects",
			"~/src",
			"~/code",
			"~/workspace",
		}
	}
	return nil
}

// Fetch retrieves recent Git commits from local repositories
func (lgc *LocalGitCommitsPlugin) Fetch(ctx context.Context) (interface{}, error) {
	var allCommits []GitCommit

	for _, repoPath := range lgc.repositories {
		// Expand home directory
		if strings.HasPrefix(repoPath, "~/") {
			home, _ := os.UserHomeDir()
			repoPath = filepath.Join(home, repoPath[2:])
		}

		commits, err := lgc.getCommitsFromRepo(ctx, repoPath)
		if err != nil {
			// Log error but continue with other repositories
			fmt.Printf("Error fetching commits from %s: %v\n", repoPath, err)
			continue
		}
		allCommits = append(allCommits, commits...)
	}

	// Filter commits by the configured Git user
	var userCommits []GitCommit
	for _, commit := range allCommits {
		if commit.Author == lgc.gitUser || strings.Contains(commit.Author, lgc.gitUser) {
			userCommits = append(userCommits, commit)
		}
	}

	// Sort by date (most recent first) and limit to 10
	if len(userCommits) > 1 {
		for i := 0; i < len(userCommits)-1; i++ {
			for j := i + 1; j < len(userCommits); j++ {
				if userCommits[i].Date.Before(userCommits[j].Date) {
					userCommits[i], userCommits[j] = userCommits[j], userCommits[i]
				}
			}
		}
	}

	if len(userCommits) > 10 {
		userCommits = userCommits[:10]
	}

	lgc.lastData = userCommits
	return userCommits, nil
}

// getCommitsFromRepo fetches commits from a specific repository
func (lgc *LocalGitCommitsPlugin) getCommitsFromRepo(ctx context.Context, repoPath string) ([]GitCommit, error) {
	// Check if it's a Git repository
	gitDir := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("not a git repository: %s", repoPath)
	}

	// Get recent commits (last 20 commits)
	cmd := exec.CommandContext(ctx, "git", "-C", repoPath, "log", "--oneline", "--format=%H|%s|%an|%ad", "--date=iso", "-20")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git log: %w", err)
	}

	var commits []GitCommit
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	repoName := filepath.Base(repoPath)
	if repoName == "." {
		// Get current directory name
		pwd, _ := os.Getwd()
		repoName = filepath.Base(pwd)
	}

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 4 {
			continue
		}

		hash := parts[0]
		message := parts[1]
		author := parts[2]
		dateStr := parts[3]

		date, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr)
		if err != nil {
			// Try alternative format
			date, err = time.Parse("2006-01-02 15:04:05", dateStr[:19])
			if err != nil {
				continue
			}
		}

		commits = append(commits, GitCommit{
			Hash:       hash[:8], // Short hash
			Message:    message,
			Author:     author,
			Date:       date,
			Repository: repoName,
		})
	}

	return commits, nil
}

// Cleanup performs cleanup
func (lgc *LocalGitCommitsPlugin) Cleanup() error {
	return nil
}

// GitHubPRsPlugin fetches Pull Requests from GitHub for the configured user
type GitHubPRsPlugin struct {
	id          string
	pluginType  string
	name        string
	version     string
	description string
	author      string
	githubToken string
	githubUser  string
	client      *http.Client
	lastData    []GitPullRequest
}

// NewGitHubPRsPlugin creates a new GitHub PRs plugin
func NewGitHubPRsPlugin() *GitHubPRsPlugin {
	// Try to get GitHub token from environment
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		githubToken = os.Getenv("GH_TOKEN")
	}

	// Get GitHub username from Git config or environment
	githubUser := getGitConfig("github.user")
	if githubUser == "" {
		githubUser = os.Getenv("GITHUB_USER")
	}
	if githubUser == "" {
		// Fallback to Git user name
		githubUser = getGitConfig("user.name")
	}

	return &GitHubPRsPlugin{
		id:          "github-prs",
		pluginType:  "git",
		name:        "GitHub Pull Requests",
		version:     "1.0.0",
		description: "Fetches Pull Requests from GitHub for the configured user",
		author:      "GoDay Team",
		githubToken: githubToken,
		githubUser:  githubUser,
		client:      &http.Client{Timeout: 15 * time.Second},
		lastData:    []GitPullRequest{},
	}
}

// GetID returns the plugin ID
func (gpr *GitHubPRsPlugin) GetID() string {
	return gpr.id
}

// GetType returns the plugin type
func (gpr *GitHubPRsPlugin) GetType() string {
	return gpr.pluginType
}

// GetMetadata returns plugin metadata
func (gpr *GitHubPRsPlugin) GetMetadata() PluginMetadata {
	return PluginMetadata{
		Name:        gpr.name,
		Version:     gpr.version,
		Description: gpr.description,
		Author:      gpr.author,
		Type:        gpr.pluginType,
		Config: map[string]string{
			"github_user":      gpr.githubUser,
			"has_github_token": fmt.Sprintf("%t", gpr.githubToken != ""),
		},
	}
}

// Initialize sets up the plugin with configuration
func (gpr *GitHubPRsPlugin) Initialize(config map[string]interface{}) error {
	if token, ok := config["github_token"].(string); ok && token != "" {
		gpr.githubToken = token
	}
	if user, ok := config["github_user"].(string); ok && user != "" {
		gpr.githubUser = user
	}
	return nil
}

// Fetch retrieves Pull Requests from GitHub
func (gpr *GitHubPRsPlugin) Fetch(ctx context.Context) (interface{}, error) {
	if gpr.githubUser == "" {
		return gpr.lastData, fmt.Errorf("GitHub user not configured")
	}

	// Search for PRs created by the user
	url := fmt.Sprintf("https://api.github.com/search/issues?q=type:pr+author:%s+is:open&sort=updated&per_page=10", gpr.githubUser)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return gpr.lastData, err
	}

	// Add GitHub token if available
	if gpr.githubToken != "" {
		req.Header.Set("Authorization", "token "+gpr.githubToken)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := gpr.client.Do(req)
	if err != nil {
		return gpr.lastData, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return gpr.lastData, err
	}

	var searchResult struct {
		Items []struct {
			Number int    `json:"number"`
			Title  string `json:"title"`
			State  string `json:"state"`
			User   struct {
				Login string `json:"login"`
			} `json:"user"`
			CreatedAt  time.Time `json:"created_at"`
			UpdatedAt  time.Time `json:"updated_at"`
			HTMLURL    string    `json:"html_url"`
			Draft      bool      `json:"draft"`
			Repository struct {
				Name string `json:"name"`
			} `json:"repository"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &searchResult); err != nil {
		return gpr.lastData, err
	}

	var prs []GitPullRequest
	for _, item := range searchResult.Items {
		prs = append(prs, GitPullRequest{
			Number:     item.Number,
			Title:      item.Title,
			State:      item.State,
			Author:     item.User.Login,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
			Repository: item.Repository.Name,
			URL:        item.HTMLURL,
			IsDraft:    item.Draft,
		})
	}

	gpr.lastData = prs
	return prs, nil
}

// Cleanup performs cleanup
func (gpr *GitHubPRsPlugin) Cleanup() error {
	return nil
}
