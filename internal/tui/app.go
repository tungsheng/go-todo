package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tungsheng/go-todo/internal/model"
	"github.com/tungsheng/go-todo/internal/storage"
)

type mode int

const (
	modeList mode = iota
	modeNew
	modeEdit
	modeConfirmDelete
)

type Model struct {
	storage     *storage.Storage
	todos       []model.Todo
	cursor      int
	mode        mode
	input       textinput.Model
	timeFilter  string // "today", "week", "month"
	timeFilters []string
	width       int
	err         error
}

func New(s *storage.Storage) (*Model, error) {
	ti := textinput.New()
	ti.Placeholder = "Enter todo title..."
	ti.CharLimit = 100
	ti.Width = 40

	m := &Model{
		storage:     s,
		input:       ti,
		timeFilter:  "today",
		timeFilters: []string{"today", "week", "month"},
	}
	m.refreshTodos()
	if m.err != nil {
		return nil, m.err
	}
	return m, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case modeList:
			return m.updateList(msg)
		case modeNew, modeEdit:
			return m.updateInput(msg)
		case modeConfirmDelete:
			return m.updateConfirm(msg)
		}
	}

	return m, nil
}

func (m Model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "j", "down":
		if m.cursor < len(m.todos)-1 {
			m.cursor++
		}

	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
		}

	case "n":
		m.mode = modeNew
		m.input.SetValue("")
		m.input.Focus()
		return m, textinput.Blink

	case "e":
		if len(m.todos) > 0 {
			m.mode = modeEdit
			m.input.SetValue(m.todos[m.cursor].Title)
			m.input.Focus()
			return m, textinput.Blink
		}

	case "d":
		if len(m.todos) > 0 {
			m.mode = modeConfirmDelete
		}

	case " ", "s":
		if len(m.todos) > 0 {
			todo := &m.todos[m.cursor]
			todo.Status = todo.Status.Next()
			if err := m.storage.Update(todo); err != nil {
				m.err = err
			}
			m.refreshTodos()
		}

	case "t":
		// Cycle time tag filter
		for i, f := range m.timeFilters {
			if f == m.timeFilter {
				m.timeFilter = m.timeFilters[(i+1)%len(m.timeFilters)]
				break
			}
		}
		m.refreshTodos()
		m.cursor = 0

	case "x":
		// Toggle close status
		if len(m.todos) > 0 {
			todo := &m.todos[m.cursor]
			todo.Status = todo.Status.ToggleClosed()
			if err := m.storage.Update(todo); err != nil {
				m.err = err
			}
			m.refreshTodos()
		}
	}

	return m, nil
}

func (m Model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		title := strings.TrimSpace(m.input.Value())
		if title != "" {
			if m.mode == modeNew {
				_, err := m.storage.Create(title, model.TimeTag(m.timeFilter))
				if err != nil {
					m.err = err
				}
			} else if m.mode == modeEdit && len(m.todos) > 0 {
				todo := &m.todos[m.cursor]
				todo.Title = title
				if err := m.storage.Update(todo); err != nil {
					m.err = err
				}
			}
			m.refreshTodos()
		}
		m.mode = modeList
		m.input.Blur()
		return m, nil

	case "esc":
		m.mode = modeList
		m.input.Blur()
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		if len(m.todos) > 0 {
			if err := m.storage.Delete(m.todos[m.cursor].ID); err != nil {
				m.err = err
			}
			m.refreshTodos()
			if m.cursor >= len(m.todos) && m.cursor > 0 {
				m.cursor--
			}
		}
		m.mode = modeList

	case "n", "N", "esc":
		m.mode = modeList
	}

	return m, nil
}

func (m *Model) refreshTodos() {
	var err error
	m.todos, err = m.storage.ListFiltered(m.timeFilter)
	if err != nil {
		m.err = err
	}
}

func (m Model) renderTodoLine(i int) string {
	todo := m.todos[i]
	cursor := "  "
	style := itemStyle
	if i == m.cursor {
		cursor = "> "
		style = selectedItemStyle
	}

	icon := StatusStyle(todo.Status).Render(todo.Status.Icon())
	line := fmt.Sprintf("%s%s %s", cursor, icon, todo.Title)
	return style.Render(line) + "\n"
}

func (m Model) View() string {
	var b strings.Builder

	// Header
	title := titleStyle.Render("go-todo")
	filter := timeTagStyle.Render(model.TimeTag(m.timeFilter).Label())
	combined := lipgloss.JoinHorizontal(lipgloss.Center, title, filter)

	header := headerStyle.Width(m.width / 2).Render(combined)

	b.WriteString(header)
	b.WriteString("\n\n")

	// Separate active and closed todos
	var activeTodos, closedTodos []int
	for i, todo := range m.todos {
		if todo.Status == model.StatusClosed {
			closedTodos = append(closedTodos, i)
		} else {
			activeTodos = append(activeTodos, i)
		}
	}

	if len(m.todos) == 0 {
		b.WriteString(itemStyle.Render("No todos yet. Press 'n' to create one."))
		b.WriteString("\n")
	} else {
		// Render active todos
		for _, i := range activeTodos {
			b.WriteString(m.renderTodoLine(i))
		}

		// Render closed todos with separator
		if len(closedTodos) > 0 {
			if len(activeTodos) > 0 {
				b.WriteString(separatorStyle.Render("── closed items ──"))
				b.WriteString("\n")
			}
			for _, i := range closedTodos {
				b.WriteString(m.renderTodoLine(i))
			}
		}
	}

	switch m.mode {
	case modeNew:
		b.WriteString("\n")
		b.WriteString(inputStyle.Render("New todo: " + m.input.View()))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("enter: save  esc: cancel"))

	case modeEdit:
		b.WriteString("\n")
		b.WriteString(inputStyle.Render("Edit: " + m.input.View()))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("enter: save  esc: cancel"))

	case modeConfirmDelete:
		if len(m.todos) > 0 {
			b.WriteString("\n")
			b.WriteString(confirmStyle.Render(fmt.Sprintf("Delete '%s'? (y/n)", m.todos[m.cursor].Title)))
		}

	default:
		b.WriteString("\n")
		help := "n:new  e:edit  d:delete  s:status  t:tag  x:close  q:quit"
		b.WriteString(helpStyle.Render(help))
	}

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#E74C3C")).Render("Error: " + m.err.Error()))
	}

	return b.String()
}
