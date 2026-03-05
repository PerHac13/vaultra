package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// configSetupView renders the config setup screen
func (m *Model) configSetupView() string {
	width := m.width - 4
	if width < 60 {
		width = 60
	}

	var s strings.Builder

	s.WriteString(RenderHeader(IconSettings, "Vaultra — Configuration Setup") + "\n\n")
	s.WriteString(MutedStyle.Render("No configuration loaded. Choose how to provide your config:") + "\n\n")

	switch m.configSetup.mode {
	case ConfigModeChoose:
		s.WriteString(m.configSetupChooseView())
	case ConfigModeFilePath:
		s.WriteString(m.configSetupFilePathView(width))
	case ConfigModeManual:
		s.WriteString(m.configSetupManualView(width))
	}

	if m.configSetup.errorMsg != "" {
		s.WriteString("\n" + RenderError(m.configSetup.errorMsg) + "\n")
	}

	s.WriteString("\n" + RenderHelpText("j/k navigate • Enter select • Esc back • q quit"))

	return LargeBoxStyle.Width(width).Render(s.String())
}

// configSetupChooseView renders the mode selection
func (m *Model) configSetupChooseView() string {
	var s strings.Builder

	options := []struct {
		icon string
		name string
		desc string
	}{
		{"📁", "Load from file path", "Provide a path to an existing YAML config file"},
		{"✏️ ", "Type YAML manually", "Enter configuration directly in the editor"},
	}

	for i, opt := range options {
		if i == m.configSetup.selected {
			s.WriteString(SelectedMenuStyle.Render(BulletSelected + " " + opt.icon + " " + opt.name) + "\n")
			s.WriteString("     " + InfoStyle.Render(opt.desc) + "\n\n")
		} else {
			s.WriteString(MenuStyle.Render(BulletUnselected + " " + opt.icon + " " + opt.name) + "\n")
			s.WriteString("     " + MutedStyle.Render(opt.desc) + "\n\n")
		}
	}

	return s.String()
}

// configSetupFilePathView renders the file path input
func (m *Model) configSetupFilePathView(width int) string {
	var s strings.Builder

	s.WriteString(DashboardTitleStyle.Render("📁 Load Configuration File") + "\n\n")
	s.WriteString(MutedStyle.Render("Enter the path to your YAML config file:") + "\n")

	// Render text input field
	inputContent := m.configSetup.filePathInput
	if inputContent == "" {
		inputContent = PlaceholderStyle.Render("e.g., ./configs/example-postgres.yaml")
	} else {
		inputContent = inputContent + CursorStyle.Render("█")
	}

	inputWidth := width - 10
	if inputWidth < 40 {
		inputWidth = 40
	}
	s.WriteString(TextInputActiveStyle.Width(inputWidth).Render(inputContent) + "\n")

	s.WriteString("\n" + InfoStyle.Render("Tip: ") + MutedStyle.Render("Use relative or absolute paths. Example configs in ./configs/"))

	return s.String()
}

// configSetupManualView renders the manual YAML input
func (m *Model) configSetupManualView(width int) string {
	var s strings.Builder

	s.WriteString(DashboardTitleStyle.Render("✏️  Type YAML Configuration") + "\n\n")
	s.WriteString(MutedStyle.Render("Enter your YAML configuration below:") + "\n")

	// Render YAML editor
	editorContent := m.configSetup.yamlInput
	if editorContent == "" {
		editorContent = PlaceholderStyle.Render("app:\n  name: \"My Backup\"\n  log_level: info\n\ndatabase:\n  type: postgres\n  config:\n    host: localhost\n    port: 5432\n    ...")
	} else {
		editorContent = editorContent + CursorStyle.Render("█")
	}

	editorWidth := width - 10
	if editorWidth < 40 {
		editorWidth = 40
	}

	editorStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(1, 2).
		Width(editorWidth).
		Height(14)

	s.WriteString(editorStyle.Render(editorContent) + "\n")

	s.WriteString("\n" + InfoStyle.Render("Enter: ") + MutedStyle.Render("new line") +
		"  " + InfoStyle.Render("Ctrl+S: ") + MutedStyle.Render("submit") +
		"  " + InfoStyle.Render("Esc: ") + MutedStyle.Render("back"))

	return s.String()
}

// configSetupKeyPress handles key presses in config setup
func (m *Model) configSetupKeyPress(key string) {
	switch m.configSetup.mode {
	case ConfigModeChoose:
		switch key {
		case "up", "k":
			if m.configSetup.selected > 0 {
				m.configSetup.selected--
			}
		case "down", "j":
			if m.configSetup.selected < 1 {
				m.configSetup.selected++
			}
		case "enter":
			if m.configSetup.selected == 0 {
				m.configSetup.mode = ConfigModeFilePath
			} else {
				m.configSetup.mode = ConfigModeManual
			}
			m.configSetup.errorMsg = ""
		}

	case ConfigModeFilePath:
		switch key {
		case "backspace":
			if len(m.configSetup.filePathInput) > 0 {
				m.configSetup.filePathInput = m.configSetup.filePathInput[:len(m.configSetup.filePathInput)-1]
			}
		case "esc":
			m.configSetup.mode = ConfigModeChoose
			m.configSetup.errorMsg = ""
		default:
			// For special keys that shouldn't be typed
			if len(key) == 1 {
				m.configSetup.filePathInput += key
			} else if key == "space" {
				m.configSetup.filePathInput += " "
			}
		}

	case ConfigModeManual:
		switch key {
		case "backspace":
			if len(m.configSetup.yamlInput) > 0 {
				m.configSetup.yamlInput = m.configSetup.yamlInput[:len(m.configSetup.yamlInput)-1]
			}
		case "esc":
			m.configSetup.mode = ConfigModeChoose
			m.configSetup.errorMsg = ""
		case "enter":
			m.configSetup.yamlInput += "\n"
		default:
			if len(key) == 1 {
				m.configSetup.yamlInput += key
			} else if key == "space" {
				m.configSetup.yamlInput += " "
			} else if key == "tab" {
				m.configSetup.yamlInput += "  "
			}
		}
	}
}
