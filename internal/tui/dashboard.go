package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// dashboardView renders the grid-based dashboard layout
//
// Layout:
//   ┌──────────────┬──────────────┬──────────────────────┐
//   │  Navigation  │  Activity    │  Config Info          │
//   │  Menu        │  Log         │  (db, storage, etc.)  │
//   └──────────────┴──────────────┴──────────────────────┘
//   ┌────────────────────────────────────────────────────┐
//   │  Content Area (changes based on menu selection)    │
//   └────────────────────────────────────────────────────┘
//   [Help bar]
func (m *Model) dashboardView() string {
	// Calculate dimensions
	totalWidth := m.width - 2
	if totalWidth < 80 {
		totalWidth = 80
	}
	totalHeight := m.height - 4
	if totalHeight < 20 {
		totalHeight = 20
	}

	// Top row: 3 panes — nav (25%), activity (35%), config info (40%)
	navWidth := totalWidth * 25 / 100
	activityWidth := totalWidth * 35 / 100
	configInfoWidth := totalWidth - navWidth - activityWidth - 6 // account for borders
	topHeight := totalHeight * 40 / 100
	if topHeight < 10 {
		topHeight = 10
	}

	// Bottom row: full width content area
	contentHeight := totalHeight - topHeight - 4

	// Build panes
	navPane := m.renderNavPane(navWidth, topHeight)
	activityPane := m.renderActivityPane(activityWidth, topHeight)
	configPane := m.renderConfigInfoPane(configInfoWidth, topHeight)

	// Top row — join horizontally
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, navPane, activityPane, configPane)

	// Content area
	contentPane := m.renderContentPane(totalWidth-2, contentHeight)

	// Join everything vertically
	dashboard := lipgloss.JoinVertical(lipgloss.Left, topRow, contentPane)

	// Help bar
	helpText := RenderHelpText("Tab=switch pane • j/k=navigate • Enter=select • Esc=back • q=quit")

	return lipgloss.JoinVertical(lipgloss.Left, dashboard, helpText)
}

// renderNavPane renders the navigation menu pane
func (m *Model) renderNavPane(width, height int) string {
	var s strings.Builder
	s.WriteString(DashboardTitleStyle.Render("📋 Navigation") + "\n\n")

	for i, choice := range m.menuOptions {
		if i == m.mainMenu.selected {
			s.WriteString(SelectedMenuStyle.Render(BulletSelected + " " + choice) + "\n")
		} else {
			s.WriteString(MenuStyle.Render(BulletUnselected + " " + choice) + "\n")
		}
	}

	style := m.getPaneStyle(PaneNav)
	return style.Width(width).Height(height).Render(s.String())
}

// renderActivityPane renders the activity/status log pane
func (m *Model) renderActivityPane(width, height int) string {
	var s strings.Builder
	s.WriteString(DashboardTitleStyle.Render("📊 Activity") + "\n\n")

	if len(m.activityLog) == 0 {
		s.WriteString(MutedStyle.Render("No activity yet"))
	} else {
		// Show most recent entries (latest at bottom)
		maxEntries := height - 5
		if maxEntries < 3 {
			maxEntries = 3
		}
		start := 0
		if len(m.activityLog) > maxEntries {
			start = len(m.activityLog) - maxEntries
		}
		for i := start; i < len(m.activityLog); i++ {
			entry := m.activityLog[i]
			// Truncate long entries
			if len(entry) > width-8 {
				entry = entry[:width-11] + "..."
			}
			s.WriteString(ActivityLogStyle.Render("• " + entry) + "\n")
		}
	}

	style := InactivePaneStyle
	return style.Width(width).Height(height).Render(s.String())
}

// renderConfigInfoPane renders the config info pane
func (m *Model) renderConfigInfoPane(width, height int) string {
	var s strings.Builder
	s.WriteString(DashboardTitleStyle.Render("⚙️  Config") + "\n\n")

	if !m.configLoaded || m.app == nil {
		s.WriteString(WarningStyle.Render("No config loaded") + "\n")
		s.WriteString(MutedStyle.Render("Use Settings to load"))
	} else {
		cfg := m.app.Config()

		s.WriteString(renderConfigStat("App", cfg.App.Name) + "\n")
		s.WriteString(renderConfigStat("Database", cfg.Database.Type) + "\n")

		// Database host info
		if host, ok := cfg.Database.Config["host"]; ok {
			port := cfg.Database.Config["port"]
			s.WriteString(renderConfigStat("  Host", fmt.Sprintf("%v:%v", host, port)) + "\n")
		}
		if dbName, ok := cfg.Database.Config["database"]; ok {
			s.WriteString(renderConfigStat("  DB Name", fmt.Sprintf("%v", dbName)) + "\n")
		}

		s.WriteString(renderConfigStat("Storage", cfg.Storage.Type) + "\n")

		if cfg.Compression.Algorithm != "" {
			s.WriteString(renderConfigStat("Compress", fmt.Sprintf("%s (L%d)", cfg.Compression.Algorithm, cfg.Compression.Level)) + "\n")
		}

		if m.configFile != "" {
			// Truncate long paths
			path := m.configFile
			if len(path) > width-14 {
				path = "..." + path[len(path)-width+17:]
			}
			s.WriteString("\n" + MutedStyle.Render("File: "+path))
		}
	}

	style := InactivePaneStyle
	return style.Width(width).Height(height).Render(s.String())
}

// renderContentPane renders the main content area based on current state
func (m *Model) renderContentPane(width, height int) string {
	var content string

	switch m.state {
	case StateMainMenu:
		content = m.contentWelcomeView()
	case StateBackup:
		content = m.contentBackupView()
	case StateRestore:
		content = m.contentRestoreView()
	case StateList:
		content = m.contentListView()
	case StateMetrics:
		content = m.contentMetricsView()
	case StateSettings:
		content = m.contentSettingsView()
	case StateProgress:
		content = m.contentProgressView()
	case StateResult:
		content = m.contentResultView()
	default:
		content = m.contentWelcomeView()
	}

	style := m.getPaneStyle(PaneContent)
	return style.Width(width).Height(height).Render(content)
}

// contentWelcomeView renders the welcome content for the dashboard
func (m *Model) contentWelcomeView() string {
	var s strings.Builder
	s.WriteString(DashboardTitleStyle.Render("📦 Vaultra Dashboard") + "\n\n")

	if m.configLoaded {
		s.WriteString(SuccessStyle.Render("✓ Configuration loaded and ready") + "\n\n")
		s.WriteString(MutedStyle.Render("Select an action from the navigation menu on the left.") + "\n\n")
		s.WriteString(InfoStyle.Render("Quick actions:") + "\n")
		s.WriteString(MutedStyle.Render("  • Press 'Enter' to select menu item") + "\n")
		s.WriteString(MutedStyle.Render("  • Press 'Tab' to switch between panes") + "\n")
	} else {
		s.WriteString(WarningStyle.Render("⚠ No configuration loaded") + "\n\n")
		s.WriteString(MutedStyle.Render("Please load a configuration to get started.") + "\n")
	}

	return s.String()
}

// contentBackupView renders backup flow in the content pane
func (m *Model) contentBackupView() string {
	var s strings.Builder
	s.WriteString(DashboardTitleStyle.Render(IconBackup+" Backup Database") + "\n\n")

	switch m.backupFlow.state {
	case "select_db":
		s.WriteString("Select Database:\n\n")
		for i, db := range m.backupFlow.databases {
			s.WriteString(RenderMenu(db, i == m.backupFlow.selectedDB) + "\n")
		}
	case "select_strategy":
		s.WriteString(RenderStat("Database:", m.backupFlow.databases[m.backupFlow.selectedDB]) + "\n\n")
		s.WriteString("Select Strategy:\n\n")
		for i, strat := range m.backupFlow.strategies {
			s.WriteString(RenderMenu(strat, i == m.backupFlow.selectedStrategy) + "\n")
		}
	case "confirm":
		s.WriteString(RenderStat("Database:", m.backupFlow.databases[m.backupFlow.selectedDB]) + "\n")
		s.WriteString(RenderStat("Strategy:", m.backupFlow.strategies[m.backupFlow.selectedStrategy]) + "\n\n")
		s.WriteString("Continue with backup? (y/n)\n")
	}

	return s.String()
}

// contentRestoreView renders restore flow in the content pane
func (m *Model) contentRestoreView() string {
	var s strings.Builder
	s.WriteString(DashboardTitleStyle.Render(IconRestore+" Restore from Backup") + "\n\n")

	if len(m.restoreFlow.backups) == 0 {
		s.WriteString(RenderEmptyState("No backups available"))
		return s.String()
	}

	for i, backup := range m.restoreFlow.backups {
		s.WriteString(RenderMenu(backup, i == m.restoreFlow.selectedBackup) + "\n")
	}

	return s.String()
}

// contentListView renders the backup list in the content pane
func (m *Model) contentListView() string {
	var s strings.Builder
	s.WriteString(DashboardTitleStyle.Render(IconList+" Available Backups") + "\n\n")

	if len(m.listModel.backups) == 0 {
		s.WriteString(RenderEmptyState("No backups found"))
		return s.String()
	}

	// Headers
	s.WriteString(
		TableHeaderStyle.Render("ID") + " " +
			TableHeaderStyle.Render("Database") + " " +
			TableHeaderStyle.Render("Size") + " " +
			TableHeaderStyle.Render("Created") + " " +
			TableHeaderStyle.Render("Status") + "\n\n")

	for _, backup := range m.listModel.backups {
		s.WriteString(backup.ID + " " + backup.Database + " " + backup.CreatedAt + " " + backup.Status + "\n")
	}

	return s.String()
}

// contentMetricsView renders metrics in the content pane
func (m *Model) contentMetricsView() string {
	var s strings.Builder
	s.WriteString(DashboardTitleStyle.Render(IconMetrics+" Backup Metrics") + "\n\n")

	s.WriteString(RenderStat("Total Backups:", "5") + "\n")
	s.WriteString(RenderStat("Successful:", "5") + "\n")
	s.WriteString(RenderStat("Failed:", "0") + "\n")
	s.WriteString(RenderStat("Success Rate:", "100%") + "\n")
	s.WriteString(RenderStat("Total Size:", "10.2 GB") + "\n")
	s.WriteString(RenderStat("Avg Duration:", "45s") + "\n")

	return s.String()
}

// contentSettingsView renders settings in the content pane
func (m *Model) contentSettingsView() string {
	var s strings.Builder
	s.WriteString(DashboardTitleStyle.Render(IconSettings+" Settings") + "\n\n")

	s.WriteString(RenderMenu("Reload Configuration", false) + "\n")
	s.WriteString(RenderMenu("Configure Notifications", false) + "\n")
	s.WriteString(RenderMenu("Configure Storage", false) + "\n")
	s.WriteString(RenderMenu("View Configuration", false) + "\n")

	return s.String()
}

// contentProgressView renders progress in the content pane
func (m *Model) contentProgressView() string {
	return RenderProgress(50, 100, m.progressMsg)
}

// contentResultView renders operation result in the content pane
func (m *Model) contentResultView() string {
	var s string
	if m.progressErr != nil {
		s = RenderError(m.progressErr.Error())
	} else {
		s = RenderSuccess(m.progressMsg)
	}
	s += "\n\n" + RenderHelpText("Press Esc to return")
	return s
}

// getPaneStyle returns the appropriate style based on whether the pane is active
func (m *Model) getPaneStyle(pane ActivePane) lipgloss.Style {
	if m.activePane == pane {
		return ActivePaneStyle
	}
	return InactivePaneStyle
}

// renderConfigStat renders a config statistic line
func renderConfigStat(label, value string) string {
	return StatLabelStyle.Render(label+":") + " " + lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(value)
}
