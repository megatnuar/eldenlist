package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	stats         map[string]int // Player stats
	quests        []string       // List of quests
	selected      int            // Selected quest index
	completed     map[int]bool   // Track completed quests
	experience    int            // Current experience points
	levelBar      string         // Experience bar
	maxExperience int            // Experience needed to level up
}

var (
	borderStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1).Margin(1)
	headerStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("211"))
	questStyle     = lipgloss.NewStyle().PaddingLeft(2)
	selectedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("36"))
	completedStyle = lipgloss.NewStyle().Strikethrough(true).Foreground(lipgloss.Color("241"))
	barStyle       = lipgloss.NewStyle().Background(lipgloss.Color("34")).Foreground(lipgloss.Color("228"))
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.quests)-1 {
				m.selected++
			}
		case "enter":
			if !m.completed[m.selected] {
				m.completed[m.selected] = true
				m.experience += 10
				if m.experience > m.maxExperience {
					m.experience = m.maxExperience
				}
			}
		}
	}

	// Update the experience bar
	progress := float64(m.experience) / float64(m.maxExperience)
	barWidth := 30
	filled := int(progress * float64(barWidth))
	m.levelBar = fmt.Sprintf("[%s%s]", strings.Repeat("=", filled), strings.Repeat(" ", barWidth-filled))

	return m, nil
}

func (m model) View() string {
	// Stats panel
	stats := borderStyle.Render(headerStyle.Render("Stats") + "\n" +
		fmt.Sprintf("Strength: %d\nDexterity: %d\nIntelligence: %d\nHealth: %d",
			m.stats["Strength"], m.stats["Dexterity"], m.stats["Intelligence"], m.stats["Health"]))

	// Quest picker
	quests := []string{headerStyle.Render("Quests")}
	for i, quest := range m.quests {
		style := questStyle
		if m.completed[i] {
			style = completedStyle
		} else if i == m.selected {
			style = selectedStyle
		}
		quests = append(quests, style.Render(fmt.Sprintf("%d. %s", i+1, quest)))
	}
	questPanel := borderStyle.Render(strings.Join(quests, "\n"))

	// Level progression bar
	levelBar := borderStyle.Render(headerStyle.Render("Level Progression") + "\n" + barStyle.Render(m.levelBar))

	// Layout
	return lipgloss.JoinHorizontal(lipgloss.Top, stats, questPanel, levelBar)
}

func main() {
	initialModel := model{
		stats: map[string]int{
			"Strength":     10,
			"Dexterity":    8,
			"Intelligence": 7,
			"Health":       50,
		},
		quests: []string{
			"Defeat the dragon",
			"Save the villagers",
			"Retrieve the lost artifact",
			"Explore the ancient ruins",
		},
		selected:      0,
		completed:     make(map[int]bool),
		experience:    0,
		maxExperience: 100,
	}

	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting app: %v\n", err)
	}
}
