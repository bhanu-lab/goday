package main

import (
	"context"
	"fmt"
	"time"
)

// Plugin represents a generic plugin interface for all widget types
type Plugin interface {
	// GetID returns a unique identifier for the plugin
	GetID() string

	// GetType returns the plugin type (e.g., "news", "weather", "calendar")
	GetType() string

	// Initialize sets up the plugin with configuration
	Initialize(config map[string]interface{}) error

	// Fetch retrieves data from the plugin source
	Fetch(ctx context.Context) (interface{}, error)

	// GetMetadata returns plugin metadata
	GetMetadata() PluginMetadata

	// Cleanup performs any necessary cleanup
	Cleanup() error
}

// PluginMetadata contains information about a plugin
type PluginMetadata struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Type        string            `json:"type"`
	Config      map[string]string `json:"config"`
}

// NewsPlugin is a specialized interface for news providers
type NewsPlugin interface {
	Plugin

	// SetTags configures the tags to filter news by
	SetTags(tags []string)

	// GetCurrentTag returns the currently active tag
	GetCurrentTag() string

	// SetCurrentTag sets the active tag for filtering
	SetCurrentTag(tag string)

	// GetSupportedTags returns all tags supported by this plugin
	GetSupportedTags() []string
}

// PluginRegistry manages all registered plugins
type PluginRegistry struct {
	plugins    map[string]Plugin
	newsByType map[string][]NewsPlugin
}

// NewPluginRegistry creates a new plugin registry
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins:    make(map[string]Plugin),
		newsByType: make(map[string][]NewsPlugin),
	}
}

// RegisterPlugin adds a plugin to the registry
func (pr *PluginRegistry) RegisterPlugin(plugin Plugin) error {
	if plugin == nil {
		return fmt.Errorf("plugin cannot be nil")
	}

	id := plugin.GetID()
	if _, exists := pr.plugins[id]; exists {
		return fmt.Errorf("plugin with ID %s already registered", id)
	}

	pr.plugins[id] = plugin

	// If it's a news plugin, also register it in the news registry
	if newsPlugin, ok := plugin.(NewsPlugin); ok {
		pluginType := plugin.GetType()
		pr.newsByType[pluginType] = append(pr.newsByType[pluginType], newsPlugin)
	}

	return nil
}

// GetPlugin retrieves a plugin by ID
func (pr *PluginRegistry) GetPlugin(id string) (Plugin, bool) {
	plugin, exists := pr.plugins[id]
	return plugin, exists
}

// GetPluginsByType retrieves all plugins of a specific type
func (pr *PluginRegistry) GetPluginsByType(pluginType string) []Plugin {
	var plugins []Plugin
	for _, plugin := range pr.plugins {
		if plugin.GetType() == pluginType {
			plugins = append(plugins, plugin)
		}
	}
	return plugins
}

// GetNewsPlugins retrieves all news plugins of a specific type
func (pr *PluginRegistry) GetNewsPlugins(pluginType string) []NewsPlugin {
	return pr.newsByType[pluginType]
}

// GetAllNewsPlugins retrieves all registered news plugins
func (pr *PluginRegistry) GetAllNewsPlugins() []NewsPlugin {
	var allPlugins []NewsPlugin
	for _, plugins := range pr.newsByType {
		allPlugins = append(allPlugins, plugins...)
	}
	return allPlugins
}

// UnregisterPlugin removes a plugin from the registry
func (pr *PluginRegistry) UnregisterPlugin(id string) error {
	plugin, exists := pr.plugins[id]
	if !exists {
		return fmt.Errorf("plugin with ID %s not found", id)
	}

	// Cleanup the plugin
	if err := plugin.Cleanup(); err != nil {
		return fmt.Errorf("error cleaning up plugin %s: %w", id, err)
	}

	delete(pr.plugins, id)

	// Remove from news registry if applicable
	if _, ok := plugin.(NewsPlugin); ok {
		pluginType := plugin.GetType()
		plugins := pr.newsByType[pluginType]
		for i, p := range plugins {
			if p.GetID() == id {
				pr.newsByType[pluginType] = append(plugins[:i], plugins[i+1:]...)
				break
			}
		}
	}

	return nil
}

// ListPlugins returns metadata for all registered plugins
func (pr *PluginRegistry) ListPlugins() []PluginMetadata {
	var metadata []PluginMetadata
	for _, plugin := range pr.plugins {
		metadata = append(metadata, plugin.GetMetadata())
	}
	return metadata
}

// PluginManager handles plugin lifecycle and execution
type PluginManager struct {
	registry  *PluginRegistry
	scheduler *PluginScheduler
	config    *PluginConfig
}

// PluginConfig holds configuration for all plugins
type PluginConfig struct {
	Plugins map[string]map[string]interface{} `yaml:"plugins"`
}

// PluginScheduler manages scheduled execution of plugins
type PluginScheduler struct {
	tasks   map[string]*PluginTask
	stopCh  chan struct{}
	running bool
}

// PluginTask represents a scheduled plugin execution
type PluginTask struct {
	ID       string
	Plugin   Plugin
	Interval time.Duration
	LastRun  time.Time
	NextRun  time.Time
	Context  context.Context
	Cancel   context.CancelFunc
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(config *PluginConfig) *PluginManager {
	return &PluginManager{
		registry:  NewPluginRegistry(),
		scheduler: NewPluginScheduler(),
		config:    config,
	}
}

// NewPluginScheduler creates a new plugin scheduler
func NewPluginScheduler() *PluginScheduler {
	return &PluginScheduler{
		tasks:  make(map[string]*PluginTask),
		stopCh: make(chan struct{}),
	}
}

// RegisterPlugin registers a plugin with the manager
func (pm *PluginManager) RegisterPlugin(plugin Plugin) error {
	if err := pm.registry.RegisterPlugin(plugin); err != nil {
		return err
	}

	// Initialize plugin with config if available
	if pm.config != nil && pm.config.Plugins != nil {
		if pluginConfig, exists := pm.config.Plugins[plugin.GetID()]; exists {
			if err := plugin.Initialize(pluginConfig); err != nil {
				return fmt.Errorf("failed to initialize plugin %s: %w", plugin.GetID(), err)
			}
		}
	}

	return nil
}

// SchedulePlugin schedules a plugin for periodic execution
func (pm *PluginManager) SchedulePlugin(pluginID string, interval time.Duration) error {
	plugin, exists := pm.registry.GetPlugin(pluginID)
	if !exists {
		return fmt.Errorf("plugin %s not found", pluginID)
	}

	ctx, cancel := context.WithCancel(context.Background())
	task := &PluginTask{
		ID:       pluginID,
		Plugin:   plugin,
		Interval: interval,
		LastRun:  time.Now(),
		NextRun:  time.Now().Add(interval),
		Context:  ctx,
		Cancel:   cancel,
	}

	pm.scheduler.AddTask(task)
	return nil
}

// GetRegistry returns the plugin registry
func (pm *PluginManager) GetRegistry() *PluginRegistry {
	return pm.registry
}

// GetScheduler returns the plugin scheduler
func (pm *PluginManager) GetScheduler() *PluginScheduler {
	return pm.scheduler
}

// Cleanup shuts down the plugin manager
func (pm *PluginManager) Cleanup() error {
	pm.scheduler.Stop()

	// Cleanup all plugins
	for _, plugin := range pm.registry.plugins {
		if err := plugin.Cleanup(); err != nil {
			fmt.Printf("Error cleaning up plugin %s: %v\n", plugin.GetID(), err)
		}
	}

	return nil
}

// AddTask adds a task to the scheduler
func (ps *PluginScheduler) AddTask(task *PluginTask) {
	ps.tasks[task.ID] = task
}

// RemoveTask removes a task from the scheduler
func (ps *PluginScheduler) RemoveTask(taskID string) {
	if task, exists := ps.tasks[taskID]; exists {
		task.Cancel()
		delete(ps.tasks, taskID)
	}
}

// Start starts the plugin scheduler
func (ps *PluginScheduler) Start() {
	if ps.running {
		return
	}
	ps.running = true

	go ps.run()
}

// Stop stops the plugin scheduler
func (ps *PluginScheduler) Stop() {
	if !ps.running {
		return
	}

	close(ps.stopCh)
	ps.running = false

	// Cancel all tasks
	for _, task := range ps.tasks {
		task.Cancel()
	}
}

// run is the main scheduler loop
func (ps *PluginScheduler) run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ps.stopCh:
			return
		case now := <-ticker.C:
			ps.checkAndExecuteTasks(now)
		}
	}
}

// checkAndExecuteTasks checks for due tasks and executes them
func (ps *PluginScheduler) checkAndExecuteTasks(now time.Time) {
	for _, task := range ps.tasks {
		if now.After(task.NextRun) || now.Equal(task.NextRun) {
			go ps.executeTask(task, now)
		}
	}
}

// executeTask executes a plugin task
func (ps *PluginScheduler) executeTask(task *PluginTask, now time.Time) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Plugin %s panicked: %v\n", task.ID, r)
		}
	}()

	// Update timing
	task.LastRun = now
	task.NextRun = now.Add(task.Interval)

	// Execute plugin
	ctx, cancel := context.WithTimeout(task.Context, 30*time.Second)
	defer cancel()

	_, err := task.Plugin.Fetch(ctx)
	if err != nil {
		fmt.Printf("Plugin %s execution failed: %v\n", task.ID, err)
	}
}
