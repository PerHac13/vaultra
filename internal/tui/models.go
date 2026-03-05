package tui

import (
	"github.com/PerHac13/vaultra/internal/app"
	"github.com/PerHac13/vaultra/internal/config"
)

// AppState represents different screens in the TUI
type AppState string

const (
	StateConfigSetup AppState = "config_setup"
	StateMainMenu    AppState = "main_menu"
	StateBackup      AppState = "backup"
	StateRestore     AppState = "restore"
	StateProgress    AppState = "progress"
	StateResult      AppState = "result"
	StateList        AppState = "list"
	StateMetrics     AppState = "metrics"
	StateSettings    AppState = "settings"
)

// ActivePane represents which pane has focus in the dashboard
type ActivePane string

const (
	PaneNav     ActivePane = "nav"
	PaneContent ActivePane = "content"
)

// ConfigSetupMode represents the config input mode
type ConfigSetupMode string

const (
	ConfigModeChoose   ConfigSetupMode = "choose"
	ConfigModeFilePath ConfigSetupMode = "filepath"
	ConfigModeManual   ConfigSetupMode = "manual"
)

// Model is the root model for the TUI application
type Model struct {
	// Application reference (nil until config loaded)
	app *app.App

	// Current state
	state AppState

	// Config management
	configFile   string
	configLoaded bool
	configSetup  *ConfigSetupModel

	// Dashboard grid
	activePane  ActivePane
	activityLog []string

	// Menu state
	menuChoice   int
	menuOptions  []string
	selectedDB   string
	selectedBackup string

	// Progress tracking
	progressMsg  string
	progressErr  error
	progressDone bool

	// Dimensions
	width  int
	height int

	// Sub-models for different flows
	mainMenu    *MainMenuModel
	backupFlow  *BackupFlowModel
	restoreFlow *RestoreFlowModel
	listModel   *ListBackupsModel
}

// ConfigSetupModel handles the config setup screen
type ConfigSetupModel struct {
	mode          ConfigSetupMode
	selected      int // 0 = file path, 1 = manual input
	filePathInput string
	yamlInput     string
	cursorPos     int
	errorMsg      string
}

// MainMenuModel handles the main menu navigation
type MainMenuModel struct {
	choices  []string
	selected int
}

// BackupFlowModel handles the backup workflow
type BackupFlowModel struct {
	state            string
	selectedDB       int
	selectedStrategy int
	databases        []string
	strategies       []string
	confirmed        bool
}

// RestoreFlowModel handles the restore workflow
type RestoreFlowModel struct {
	state          string
	selectedBackup int
	backups        []string
	confirmed      bool
	dryRun         bool
}

// ListBackupsModel displays available backups
type ListBackupsModel struct {
	backups  []BackupItem
	selected int
}

// BackupItem represents a single backup in the list
type BackupItem struct {
	ID        string
	Database  string
	Size      int64
	CreatedAt string
	Status    string
}

// NewModel creates a new TUI model
func NewModel(vaultraApp *app.App, configFile string) *Model {
	menuOptions := []string{
		"Backup Database",
		"Restore from Backup",
		"List Backups",
		"View Metrics",
		"Settings",
		"Exit",
	}

	startState := StateConfigSetup
	configLoaded := false

	if vaultraApp != nil {
		startState = StateMainMenu
		configLoaded = true
	}

	m := &Model{
		app:          vaultraApp,
		state:        startState,
		configFile:   configFile,
		configLoaded: configLoaded,
		configSetup:  NewConfigSetupModel(),
		activePane:   PaneNav,
		activityLog:  []string{"Vaultra TUI started"},
		menuOptions:  menuOptions,
		width:        120,
		height:       40,
	}

	// Initialize sub-models
	m.mainMenu = NewMainMenuModel(m.menuOptions)
	m.backupFlow = NewBackupFlowModel()
	m.restoreFlow = NewRestoreFlowModel()
	m.listModel = NewListBackupsModel()

	if configLoaded {
		m.activityLog = append(m.activityLog, "Config loaded: "+configFile)
		m.activityLog = append(m.activityLog, "Ready")
	}

	return m
}

// NewConfigSetupModel creates a new config setup model
func NewConfigSetupModel() *ConfigSetupModel {
	return &ConfigSetupModel{
		mode:     ConfigModeChoose,
		selected: 0,
	}
}

// NewMainMenuModel creates a new main menu model
func NewMainMenuModel(choices []string) *MainMenuModel {
	return &MainMenuModel{
		choices:  choices,
		selected: 0,
	}
}

// NewBackupFlowModel creates a new backup flow model
func NewBackupFlowModel() *BackupFlowModel {
	return &BackupFlowModel{
		state:            "select_db",
		selectedDB:       0,
		selectedStrategy: 0,
		databases:        []string{"PostgreSQL", "MySQL", "MongoDB"},
		strategies:       []string{"Full Backup", "Incremental Backup"},
	}
}

// NewRestoreFlowModel creates a new restore flow model
func NewRestoreFlowModel() *RestoreFlowModel {
	return &RestoreFlowModel{
		state:          "list_backups",
		selectedBackup: 0,
		backups:        []string{}, // Will be populated from repository
		dryRun:         false,
	}
}

// NewListBackupsModel creates a new list backups model
func NewListBackupsModel() *ListBackupsModel {
	return &ListBackupsModel{
		backups:  []BackupItem{},
		selected: 0,
	}
}

// addActivity adds an entry to the activity log
func (m *Model) addActivity(msg string) {
	m.activityLog = append(m.activityLog, msg)
	// Keep only last 10 entries
	if len(m.activityLog) > 10 {
		m.activityLog = m.activityLog[len(m.activityLog)-10:]
	}
}

// Message types for the TUI

// MenuSelectMsg indicates a menu item was selected
type MenuSelectMsg struct {
	Choice int
}

// BackupConfirmMsg indicates backup should execute
type BackupConfirmMsg struct {
	Database string
	Strategy string
}

// BackupCompleteMsg indicates backup finished
type BackupCompleteMsg struct {
	Success bool
	Error   error
	Message string
}

// RestoreConfirmMsg indicates restore should execute
type RestoreConfirmMsg struct {
	BackupID string
	DryRun   bool
}

// RestoreCompleteMsg indicates restore finished
type RestoreCompleteMsg struct {
	Success bool
	Error   error
	Message string
}

// BackupsLoadedMsg indicates backups were loaded
type BackupsLoadedMsg struct {
	Backups []BackupItem
	Error   error
}

// NavigateMsg indicates navigation between screens
type NavigateMsg struct {
	State AppState
}

// ErrorMsg represents an error that occurred
type ErrorMsg struct {
	Error error
}

// SuccessMsg represents a successful operation
type SuccessMsg struct {
	Message string
}

// ConfigLoadedMsg indicates config was successfully loaded
type ConfigLoadedMsg struct {
	App        *app.App
	Config     *config.ConfigType
	ConfigFile string
}