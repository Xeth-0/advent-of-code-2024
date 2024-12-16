package tui

import (
	"fmt"

	"math"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "-"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
	// main style
	mainStyle = lipgloss.NewStyle().MarginLeft(2)

	viewportStyle = lipgloss.NewStyle().MarginLeft(2).MarginRight(2)
)

type Model struct {
	ready         bool
	title         string
	viewport      viewport.Model
	viewportLines []string
	minWidth      int
	windowHeight  int
	windowWidth   int
}

type updateViewport struct {
	content string
	width   int
	height  int
}

type updateViewportLine struct {
	lineNum int
	line    string
}

func initModel(title string) Model {
	return Model{title: title}
}

func (m Model) withViewport(lines []string) Model {
	m.viewportLines = lines
	return m
}

func (m Model) withMinWidth(minWidth int) Model {
	m.minWidth = minWidth
	return m
}

func NewViewportProgram(initModel Model) *tea.Program {
	return tea.NewProgram(
		initModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
}

// Viewport updates
func UpdateViewport(content string, width int) tea.Msg {
	return updateViewport{content: content, width: width}
}

func UpdateViewportLine(lineNum int, line string) tea.Msg {
	return updateViewportLine{lineNum: lineNum, line: line}
}

func minInt(nums ...int) int {
	result := math.MaxInt
	for _, val := range nums {
		if val < result {
			result = val
		}
	}
	return result
}

func maxInt(nums ...int) int {
	result := math.MinInt
	for _, val := range nums {
		if val > result {
			result = val
		}
	}
	return result
}

// Viewport stylized views
func (m Model) headerView() string {
	title := titleStyle.Render(m.title)
	line := strings.Repeat("-", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	line := strings.Repeat("-", max(0, m.viewport.Width))
	return lipgloss.JoinHorizontal(lipgloss.Center, line)
}

// Model
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		k := msg.String()
		if k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())

		if !m.ready {
			m.windowWidth, m.windowHeight = msg.Width, msg.Height
			if m.minWidth == 0 {
				m.minWidth = m.windowWidth
			}
			m.viewport = viewport.New(msg.Width, msg.Height-(headerHeight+footerHeight))
			m.viewport.YPosition = headerHeight
			m.ready = true

			if m.viewportLines != nil {
				m.viewport.SetContent(strings.Join(m.viewportLines, "\n"))
			}

			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - (headerHeight + footerHeight)
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	// s := "A terminal UI? why not."
	// return s
	return mainStyle.Render(fmt.Sprintf("%s\n%s\n%s\n", m.headerView(), viewportStyle.Render(m.viewport.View()), m.footerView()))
}
