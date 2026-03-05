package tui

import "github.com/charmbracelet/lipgloss"

// Color definitions
var (
	colorPrimary   = lipgloss.Color("212")  // Bright magenta
	colorSuccess   = lipgloss.Color("10")   // Bright green
	colorError     = lipgloss.Color("9")    // Bright red
	colorWarning   = lipgloss.Color("11")   // Bright yellow
	colorInfo      = lipgloss.Color("14")   // Bright cyan
	colorMuted     = lipgloss.Color("240")  // Dark gray
	colorBorder    = lipgloss.Color("8")    // Darker gray
)

// Text styling
var (
	// Title style - centered, bold, primary color
	TitleStyle = lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true).
		Align(lipgloss.Center).
		MarginBottom(1)

	// Menu item style - muted color
	MenuStyle = lipgloss.NewStyle().
		Foreground(colorMuted).
		PaddingLeft(2)

	// Selected menu item style - primary color, bold, with background
	SelectedMenuStyle = lipgloss.NewStyle().
		Foreground(colorPrimary).
		Background(lipgloss.Color("235")).
		Bold(true).
		PaddingLeft(1)

	// Success style - green and bold
	SuccessStyle = lipgloss.NewStyle().
		Foreground(colorSuccess).
		Bold(true)

	// Error style - red and bold
	ErrorStyle = lipgloss.NewStyle().
		Foreground(colorError).
		Bold(true)

	// Warning style - yellow
	WarningStyle = lipgloss.NewStyle().
		Foreground(colorWarning)

	// Info style - cyan
	InfoStyle = lipgloss.NewStyle().
		Foreground(colorInfo)

	// Muted style - gray for less important text
	MutedStyle = lipgloss.NewStyle().
		Foreground(colorMuted)

	// Box style - rounded border with padding
	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(1, 2).
		MarginBottom(1)

	// Large box - for main screens
	LargeBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(2, 4)

	// List item style
	ListItemStyle = lipgloss.NewStyle().
		PaddingLeft(2)

	// Selected list item style
	SelectedListItemStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		Background(lipgloss.Color("235")).
		Foreground(colorPrimary).
		Bold(true)

	// Help text style - small, muted, at bottom
	HelpStyle = lipgloss.NewStyle().
		Foreground(colorMuted).
		MarginTop(1).
		Align(lipgloss.Center)

	// Loading style - blinking effect (will be handled separately)
	LoadingStyle = lipgloss.NewStyle().
		Foreground(colorInfo)

	// Progress bar background
	ProgressBarStyle = lipgloss.NewStyle().
		Foreground(colorPrimary)

	// Stat label (e.g., "Size:" in metrics)
	StatLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorInfo)

	// Stat value
	StatValueStyle = lipgloss.NewStyle().
		Foreground(colorSuccess)

	// Table header
	TableHeaderStyle = lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true).
		PaddingRight(2)

	// Table row
	TableRowStyle = lipgloss.NewStyle().
		PaddingRight(2)

	// =====================
	// Grid / Dashboard styles
	// =====================

	// PaneStyle is the default bordered box for each grid cell
	PaneStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(1, 2)

	// ActivePaneStyle has a highlighted border for the focused pane
	ActivePaneStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorInfo).
		Padding(1, 2)

	// InactivePaneStyle has a dimmed border for unfocused panes
	InactivePaneStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(1, 2)

	// TextInputStyle for config file path input
	TextInputStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(colorInfo).
		Padding(0, 1).
		MarginTop(1).
		MarginBottom(1)

	// TextInputActiveStyle for focused text input
	TextInputActiveStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(colorPrimary).
		Padding(0, 1).
		MarginTop(1).
		MarginBottom(1)

	// ConfigInfoStyle for the config info panel
	ConfigInfoStyle = lipgloss.NewStyle().
		Foreground(colorMuted)

	// ActivityLogStyle for the activity/status panel
	ActivityLogStyle = lipgloss.NewStyle().
		Foreground(colorMuted).
		PaddingLeft(1)

	// DashboardTitleStyle for pane titles within the grid
	DashboardTitleStyle = lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true).
		MarginBottom(1)

	// CursorStyle for text input cursor
	CursorStyle = lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true)

	// PlaceholderStyle for text input placeholder
	PlaceholderStyle = lipgloss.NewStyle().
		Foreground(colorMuted).
		Italic(true)
)

// Theme constants
const (
	// Box characters
	BoxCornerTL = "┌"
	BoxCornerTR = "┐"
	BoxCornerBL = "└"
	BoxCornerBR = "┘"
	BoxHorizontal = "─"
	BoxVertical   = "│"
	
	// Bullets
	BulletSelected = "▶"
	BulletUnselected = " "
	
	// Indicators
	CheckMark = "✓"
	XMark     = "✗"
	Ellipsis  = "…"
	
	// Icons
	IconBackup   = "📦"
	IconRestore  = "📥"
	IconList     = "📋"
	IconMetrics  = "📊"
	IconSettings = "⚙️"
	IconError    = "❌"
	IconSuccess  = "✅"
	IconWarning  = "⚠️"
	IconInfo     = "ℹ️"
)

// Separator creates a horizontal separator line
func Separator(width int) string {
	return lipgloss.NewStyle().
		Foreground(colorBorder).
		Render(
			"─" + lipgloss.NewStyle().Width(width-2).Render("") + "─",
		)
}

// RenderHeader renders a header with icon and title
func RenderHeader(icon, title string) string {
	return TitleStyle.Render(icon + " " + title)
}

// RenderMenu renders a menu item
func RenderMenu(text string, selected bool) string {
	if selected {
		return SelectedMenuStyle.Render(BulletSelected + " " + text)
	}
	return MenuStyle.Render(BulletUnselected + " " + text)
}

// RenderBox renders content in a box
func RenderBox(content string) string {
	return BoxStyle.Render(content)
}

// RenderSuccess renders a success message
func RenderSuccess(message string) string {
	return SuccessStyle.Render(CheckMark + " " + message)
}

// RenderError renders an error message
func RenderError(message string) string {
	return ErrorStyle.Render(XMark + " " + message)
}

// RenderInfo renders an info message
func RenderInfo(message string) string {
	return InfoStyle.Render(IconInfo + " " + message)
}

// RenderStat renders a statistic
func RenderStat(label, value string) string {
	return StatLabelStyle.Render(label) + " " + StatValueStyle.Render(value)
}

// RenderProgress renders a progress indicator
func RenderProgress(current, total int, message string) string {
	percentage := (current * 100) / total
	filled := (percentage) / 10
	empty := 10 - filled

	bar := "["
	for i := 0; i < filled; i++ {
		bar += "="
	}
	for i := 0; i < empty; i++ {
		bar += "-"
	}
	bar += "]"

	return ProgressBarStyle.Render(bar + " " + message + " " + RenderStat("", string(rune(percentage))+"%"))
}

// RenderTable renders a simple table
func RenderTable(headers []string, rows [][]string) string {
	content := ""
	
	// Render headers
	for _, h := range headers {
		content += TableHeaderStyle.Render(h) + " "
	}
	content += "\n"
	
	// Render rows
	for _, row := range rows {
		for _, cell := range row {
			content += TableRowStyle.Render(cell) + " "
		}
		content += "\n"
	}
	
	return content
}

// RenderHelpText renders help text at the bottom
func RenderHelpText(keys string) string {
	return HelpStyle.Render(keys)
}

// RenderEmptyState renders an empty state message
func RenderEmptyState(message string) string {
	return MutedStyle.Render("No data to display\n" + message)
}

// RenderLoading renders a loading indicator
func RenderLoading(message string) string {
	indicators := []string{"|", "/", "-", "\\"}
	// Animation frame would be handled by caller
	return LoadingStyle.Render(indicators[0] + " " + message)
}