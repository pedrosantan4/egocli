package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

// ============== TYPES ==============
type (
	metricsUpdateMsg struct{}
	commandOutputMsg struct{ output, err string }
	cursorMsg        struct{}
)

// ============== STYLES ==============
var (
	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	resultStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DADADA")).
			Padding(0, 1)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF7F")).
			Bold(true)

	cursorStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#7D56F4")).
			Foreground(lipgloss.Color("#FFFFFF"))
)

// ============== TERMINAL MODEL ==============
type terminalModel struct {
	width, height    int
	cpuLoad, memUsed float64
	goroutines       int
	cmdHistory       []string
	inputBuffer      string
	showCursor       bool
	lastOutput       string
	historyOffset    int
}

func newTerminalModel() *terminalModel {
	return &terminalModel{
		cmdHistory:    make([]string, 0),
		showCursor:    true,
		historyOffset: -1,
	}
}

func (m *terminalModel) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("🖥️ EgoCLI Terminal"),
		tea.EnterAltScreen,
		m.startMetricsTicker(),
		m.startCursorTicker(),
	)
}

func (m *terminalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyInput(msg)

	case metricsUpdateMsg:
		return m.updateMetrics()

	case commandOutputMsg:
		m.lastOutput = formatOutput(msg.output, msg.err)
		return m, nil

	case cursorMsg:
		m.showCursor = !m.showCursor
		return m, m.startCursorTicker()

	default:
		return m, nil
	}
}

func (m *terminalModel) View() string {
	var view strings.Builder
	view.WriteString(m.renderStatusBar())

	if m.lastOutput != "" {
		view.WriteString("\n" + resultStyle.Render(m.lastOutput))
	}
	view.WriteString("\n" + m.renderInputLine())

	return view.String()
}

// ============== MODEL UTILITIES ==============
func (m *terminalModel) updateMetrics() (tea.Model, tea.Cmd) {
	go func() {
		if vmStat, err := mem.VirtualMemory(); err == nil {
			m.memUsed = vmStat.UsedPercent
		}

		if cpuPercent, err := cpu.Percent(0, false); err == nil {
			m.cpuLoad = cpuPercent[0]
		}

		m.goroutines = runtime.NumGoroutine()
	}()

	return m, tea.Tick(metricsInterval, func(time.Time) tea.Msg {
		return metricsUpdateMsg{}
	})
}

func (m *terminalModel) startMetricsTicker() tea.Cmd {
	return tea.Tick(metricsInterval, func(time.Time) tea.Msg {
		return metricsUpdateMsg{}
	})
}

func (m *terminalModel) startCursorTicker() tea.Cmd {
	return tea.Tick(cursorInterval, func(time.Time) tea.Msg {
		return cursorMsg{}
	})
}

func (m *terminalModel) renderStatusBar() string {
	return statusStyle.Render(fmt.Sprintf(
		"🖥 CPU: %.1f%%  │  📦 MEM: %.1f%%  │  🔄 GOROUTINES: %d",
		m.cpuLoad, m.memUsed, m.goroutines,
	))
}

func (m *terminalModel) renderInputLine() string {
	cursor := " "
	if m.showCursor {
		cursor = cursorStyle.Render("▌")
	}
	return inputStyle.Render(fmt.Sprintf("egocli> %s%s", m.inputBuffer, cursor))
}

// ============== INPUT HANDLER ==============
func (m *terminalModel) handleKeyInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "enter":
		return m.processCommand()
	case "up", "down":
		m.navigateCommandHistory(msg.String())
	case "backspace":
		if len(m.inputBuffer) > 0 {
			m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
		}
	default:
		if len(msg.Runes) > 0 {
			m.inputBuffer += string(msg.Runes)
		}
	}
	return m, nil
}

func (m *terminalModel) navigateCommandHistory(direction string) {
	if len(m.cmdHistory) == 0 {
		return
	}

	switch direction {
	case "up":
		m.historyOffset = min(m.historyOffset+1, len(m.cmdHistory)-1)
	case "down":
		m.historyOffset = max(m.historyOffset-1, 0)
	}

	m.inputBuffer = m.cmdHistory[len(m.cmdHistory)-1-m.historyOffset]
}

func (m *terminalModel) processCommand() (tea.Model, tea.Cmd) {
	cmd := strings.TrimSpace(m.inputBuffer)
	m.inputBuffer = ""
	m.cmdHistory = append(m.cmdHistory, cmd)
	m.historyOffset = -1 // Reset history position

	if cmd == "" {
		return m, nil
	}

	parts := strings.Fields(cmd)
	command := parts[0]
	args := parts[1:]

	switch command {
	case "clear":
		m.lastOutput = ""
		return m, nil
	case "exit":
		return m, tea.Quit
	case "gen", "new":
		return m.handleTemplateCommand(command, args)
	default:
		return m.handleDirectCommand(command, args)
	}
}

func (m *terminalModel) handleTemplateCommand(cmdType string, args []string) (tea.Model, tea.Cmd) {
	if len(args) == 0 {
		return m, func() tea.Msg {
			return commandOutputMsg{err: fmt.Sprintf("Especifique um template para %s (ex: %s vpc)", cmdType, cmdType)}
		}
	}
	return m, executeTemplateCommand(cmdType, args[0])
}

func (m *terminalModel) handleDirectCommand(templateName string, args []string) (tea.Model, tea.Cmd) {
	return m, executeTemplateCommand("gen", templateName)
}

// ============== OUTPUT FORMATTER ==============
func formatOutput(output, err string) string {
	if err != "" {
		return errorStyle.Render(fmt.Sprintf("❌ %s", err))
	}
	return successStyle.Render(fmt.Sprintf("✅ %s", output))
}

// ============== COMMAND EXECUTOR ==============
func executeTemplateCommand(commandType, templateName string) tea.Cmd {
	return func() tea.Msg {
		template, exists := Templates[templateName]
		if !exists {
			return commandOutputMsg{err: fmt.Sprintf("Template não encontrado: %s", templateName)}
		}

		outputDir := genDir
		if commandType == "new" {
			outputDir = newDir
		}

		// Use goroutine for concurrent file operations
		var wg sync.WaitGroup
		wg.Add(1)

		errChan := make(chan error, 1)
		go func() {
			defer wg.Done()
			errChan <- saveTemplate(template, outputDir)
		}()

		wg.Wait()
		close(errChan)

		if err := <-errChan; err != nil {
			return commandOutputMsg{err: err.Error()}
		}

		return commandOutputMsg{
			output: fmt.Sprintf("%s criado em %s/%s",
				templateName, outputDir, template.DirName),
		}
	}
}

func saveTemplate(template ModuleTemplate, outputDir string) error {
	targetDir := filepath.Join(outputDir, template.DirName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	filePath := filepath.Join(targetDir, template.FileName)
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("arquivo já existe: %s", filePath)
	}

	return os.WriteFile(filePath, []byte(template.Content), 0644)
}

// ============== COBRA INTEGRATION ==============
var terminalCmd = &cobra.Command{
	Use:   "terminal",
	Short: "Inicia o terminal interativo",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(
			newTerminalModel(),
			tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		)

		if _, err := p.Run(); err != nil {
			fmt.Println("Erro ao iniciar terminal:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(terminalCmd)
}
