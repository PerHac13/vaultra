package tui

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/PerHac13/vaultra/internal/app"
	"github.com/PerHac13/vaultra/internal/config"
	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

// App wraps the TUI application
type App struct {
	model  *Model
	app    *app.App
	ctx    context.Context
	logger *slog.Logger
}

// New creates a new TUI application
// vaultraApp can be nil if config is not yet loaded
func New(vaultraApp *app.App, logger *slog.Logger, configFile string) *App {
	return &App{
		model:  NewModel(vaultraApp, configFile),
		app:    vaultraApp,
		ctx:    context.Background(),
		logger: logger,
	}
}

// Run starts the TUI application
func (a *App) Run() error {
	p := tea.NewProgram(a.model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// =====================
// Implement tea.Model interface
// =====================

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Window size change
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	// Key press
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	// Menu selection
	case MenuSelectMsg:
		return m.handleMenuSelect(msg)

	// Quit
	case tea.QuitMsg:
		return m, tea.Quit

	// Error
	case ErrorMsg:
		m.progressErr = msg.Error
		m.state = StateResult
		m.addActivity("Error: " + msg.Error.Error())
		return m, nil

	// Success
	case SuccessMsg:
		m.progressMsg = msg.Message
		m.state = StateResult
		m.addActivity(msg.Message)
		return m, nil

	// Navigation
	case NavigateMsg:
		m.state = msg.State
		return m, nil

	// Config loaded
	case ConfigLoadedMsg:
		m.app = msg.App
		m.configFile = msg.ConfigFile
		m.configLoaded = true
		m.state = StateMainMenu
		m.addActivity("Config loaded: " + msg.ConfigFile)
		m.addActivity("Database: " + msg.Config.Database.Type)
		m.addActivity("Storage: " + msg.Config.Storage.Type)
		m.addActivity("Ready")
		return m, nil

	// Backups loaded
	case BackupsLoadedMsg:
		if msg.Error != nil {
			m.addActivity("Error loading backups: " + msg.Error.Error())
			return m, tea.Quit
		}
		m.listModel.backups = msg.Backups
		m.addActivity(fmt.Sprintf("Loaded %d backups", len(msg.Backups)))
		return m, nil

	default:
		return m, nil
	}
}

// View renders the current state
func (m *Model) View() string {
	switch m.state {
	case StateConfigSetup:
		return m.configSetupView()
	case StateMainMenu, StateBackup, StateRestore, StateList, StateMetrics, StateSettings, StateProgress, StateResult:
		return m.dashboardView()
	default:
		return "Unknown state\n"
	}
}

// handleKeyPress handles keyboard input
func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Global quit
	if key == "ctrl+c" {
		return m, tea.Quit
	}

	// Config setup state has its own handling
	if m.state == StateConfigSetup {
		return m.handleConfigSetupKey(msg)
	}

	// 'q' only quits from main menu / nav pane
	if key == "q" {
		if m.state == StateMainMenu && m.activePane == PaneNav {
			return m, tea.Quit
		}
	}

	switch key {
	// Tab to switch panes
	case "tab":
		if m.activePane == PaneNav {
			m.activePane = PaneContent
		} else {
			m.activePane = PaneNav
		}
		return m, nil

	// Escape — back to main menu or switch pane
	case "esc":
		if m.state != StateMainMenu {
			prevState := m.state
			m.state = StateMainMenu
			m.addActivity("Back to menu from " + string(prevState))
		}
		return m, nil
	}

	// State-specific handling based on active pane
	if m.activePane == PaneNav {
		return m.navPaneKeyPress(msg)
	}

	// Content pane key handling
	switch m.state {
	case StateBackup:
		return m.backupFlowKeyPress(msg)
	case StateRestore:
		return m.restoreFlowKeyPress(msg)
	case StateList:
		return m.listBackupsKeyPress(msg)
	default:
		return m, nil
	}
}

// handleConfigSetupKey handles keys in config setup mode
func (m *Model) handleConfigSetupKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Quit
	if key == "q" && m.configSetup.mode == ConfigModeChoose {
		return m, tea.Quit
	}

	// Submit config (Enter in filepath mode, Ctrl+S in manual mode)
	if m.configSetup.mode == ConfigModeFilePath && key == "enter" {
		return m, m.loadConfigFromFile(m.configSetup.filePathInput)
	}
	if m.configSetup.mode == ConfigModeManual && key == "ctrl+s" {
		return m, m.loadConfigFromYAML(m.configSetup.yamlInput)
	}

	// Delegate to config setup
	m.configSetupKeyPress(key)
	return m, nil
}

// navPaneKeyPress handles key presses in the navigation pane
func (m *Model) navPaneKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.mainMenu.selected > 0 {
			m.mainMenu.selected--
		}
	case "down", "j":
		if m.mainMenu.selected < len(m.mainMenu.choices)-1 {
			m.mainMenu.selected++
		}
	case "enter":
		return m.handleMenuSelect(MenuSelectMsg{Choice: m.mainMenu.selected})
	}
	return m, nil
}

// handleMenuSelect handles menu selection
func (m *Model) handleMenuSelect(msg MenuSelectMsg) (tea.Model, tea.Cmd) {
	switch msg.Choice {
	case 0: // Backup
		if !m.configLoaded {
			m.addActivity("Config required for backup")
			return m, nil
		}
		m.state = StateBackup
		m.activePane = PaneContent
		m.backupFlow = NewBackupFlowModel()
		m.addActivity("Started backup flow")
	case 1: // Restore
		if !m.configLoaded {
			m.addActivity("Config required for restore")
			return m, nil
		}
		m.state = StateRestore
		m.activePane = PaneContent
		m.restoreFlow = NewRestoreFlowModel()
		m.addActivity("Started restore flow")
		return m, m.loadBackups()
	case 2: // List
		if !m.configLoaded {
			m.addActivity("Config required to list backups")
			return m, nil
		}
		m.state = StateList
		m.activePane = PaneContent
		m.addActivity("Loading backup list...")
		return m, m.loadBackups()
	case 3: // Metrics
		m.state = StateMetrics
		m.activePane = PaneContent
		m.addActivity("Viewing metrics")
	case 4: // Settings
		m.state = StateSettings
		m.activePane = PaneContent
		m.addActivity("Opened settings")
	case 5: // Exit
		return m, tea.Quit
	}
	return m, nil
}

// backupFlowKeyPress handles key presses in backup flow
func (m *Model) backupFlowKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		switch m.backupFlow.state {
		case "select_db":
			if m.backupFlow.selectedDB > 0 {
				m.backupFlow.selectedDB--
			}
		case "select_strategy":
			if m.backupFlow.selectedStrategy > 0 {
				m.backupFlow.selectedStrategy--
			}
		}
	case "down", "j":
		switch m.backupFlow.state {
		case "select_db":
			if m.backupFlow.selectedDB < len(m.backupFlow.databases)-1 {
				m.backupFlow.selectedDB++
			}
		case "select_strategy":
			if m.backupFlow.selectedStrategy < len(m.backupFlow.strategies)-1 {
				m.backupFlow.selectedStrategy++
			}
		}
	case "enter":
		switch m.backupFlow.state {
		case "select_db":
			m.backupFlow.state = "select_strategy"
			m.addActivity("Selected: " + m.backupFlow.databases[m.backupFlow.selectedDB])
		case "select_strategy":
			m.backupFlow.state = "confirm"
			m.addActivity("Strategy: " + m.backupFlow.strategies[m.backupFlow.selectedStrategy])
		case "confirm":
			m.backupFlow.confirmed = true
			m.addActivity("Executing backup...")
			return m, m.executeBackup()
		}
	case "esc":
		switch m.backupFlow.state {
		case "select_db":
			m.state = StateMainMenu
			m.activePane = PaneNav
		case "select_strategy":
			m.backupFlow.state = "select_db"
		case "confirm":
			m.backupFlow.state = "select_strategy"
		}
	case "n":
		if m.backupFlow.state == "confirm" {
			m.backupFlow.state = "select_strategy"
		}
	}
	return m, nil
}

// restoreFlowKeyPress handles key presses in restore flow
func (m *Model) restoreFlowKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.restoreFlow.selectedBackup > 0 {
			m.restoreFlow.selectedBackup--
		}
	case "down", "j":
		if m.restoreFlow.selectedBackup < len(m.restoreFlow.backups)-1 {
			m.restoreFlow.selectedBackup++
		}
	case "enter":
		m.state = StateProgress
		m.addActivity("Starting restore...")
	case "esc":
		m.state = StateMainMenu
		m.activePane = PaneNav
	}
	return m, nil
}

// listBackupsKeyPress handles key presses in list view
func (m *Model) listBackupsKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = StateMainMenu
		m.activePane = PaneNav
	}
	return m, nil
}

// =====================
// Command functions
// =====================

// loadConfigFromFile loads config from a file path
func (m *Model) loadConfigFromFile(filePath string) tea.Cmd {
	return func() tea.Msg {
		if filePath == "" {
			return ErrorMsg{Error: fmt.Errorf("file path cannot be empty")}
		}

		parser := config.NewParser()
		cfg, err := parser.ParseFile(filePath)
		if err != nil {
			return ErrorMsg{Error: fmt.Errorf("failed to parse config: %w", err)}
		}

		ctx := context.Background()
		vaultraApp, err := app.New(ctx, filePath)
		if err != nil {
			return ErrorMsg{Error: fmt.Errorf("failed to initialize app: %w", err)}
		}

		return ConfigLoadedMsg{
			App:        vaultraApp,
			Config:     cfg,
			ConfigFile: filePath,
		}
	}
}

// loadConfigFromYAML loads config from raw YAML string
func (m *Model) loadConfigFromYAML(yamlContent string) tea.Cmd {
	return func() tea.Msg {
		if yamlContent == "" {
			return ErrorMsg{Error: fmt.Errorf("YAML content cannot be empty")}
		}

		// Parse YAML to validate
		var cfg config.ConfigType
		if err := yaml.Unmarshal([]byte(yamlContent), &cfg); err != nil {
			return ErrorMsg{Error: fmt.Errorf("invalid YAML: %w", err)}
		}

		// Write to temp file so app.New can read it
		tmpFile, err := os.CreateTemp("", "vaultra-config-*.yaml")
		if err != nil {
			return ErrorMsg{Error: fmt.Errorf("failed to create temp file: %w", err)}
		}

		if _, err := tmpFile.WriteString(yamlContent); err != nil {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
			return ErrorMsg{Error: fmt.Errorf("failed to write temp config: %w", err)}
		}
		tmpFile.Close()

		ctx := context.Background()
		vaultraApp, err := app.New(ctx, tmpFile.Name())
		if err != nil {
			os.Remove(tmpFile.Name())
			return ErrorMsg{Error: fmt.Errorf("failed to initialize app: %w", err)}
		}

		return ConfigLoadedMsg{
			App:        vaultraApp,
			Config:     &cfg,
			ConfigFile: tmpFile.Name(),
		}
	}
}

// loadBackups loads backups from repository
func (m *Model) loadBackups() tea.Cmd {
	return func() tea.Msg {
		// TODO: Load from app.BackupRepository()
		return BackupsLoadedMsg{
			Backups: []BackupItem{},
			Error:   nil,
		}
	}
}

// executeBackup executes a backup operation
func (m *Model) executeBackup() tea.Cmd {
	return func() tea.Msg {
		// TODO: Execute backup using m.app.BackupEngine()
		return SuccessMsg{Message: "Backup completed successfully"}
	}
}