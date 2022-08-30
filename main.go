package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState uint

const (
	listView sessionState = iota
	viewportView
	textareaView
)

var (
	youStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	botStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	viewportStyle = lipgloss.NewStyle().
			Padding(1, 2).
			BorderStyle(lipgloss.NormalBorder())

	modelStyle = lipgloss.NewStyle().
			Padding(1, 2).
			BorderStyle(lipgloss.NormalBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Padding(1, 2).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

var items = []list.Item{
	item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
	item{title: "Nutella", desc: "It's good on toast"},
	item{title: "Bitter melon", desc: "It cools you down"},
	item{title: "Nice socks", desc: "And by that I mean socks without holes"},
	item{title: "Eight hours of sleep", desc: "I had this once"},
	item{title: "Cats", desc: "Usually"},
	item{title: "Plantasia, the album", desc: "My plants love it too"},
	item{title: "Pour over coffee", desc: "It takes forever to make though"},
	item{title: "VR", desc: "Virtual reality...what is there to say?"},
	item{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
	item{title: "Linux", desc: "Pretty much the best OS"},
	item{title: "Business school", desc: "Just kidding"},
	item{title: "Pottery", desc: "Wet clay is a great feeling"},
	item{title: "Shampoo", desc: "Nothing like clean hair"},
	item{title: "Table tennis", desc: "It’s surprisingly exhausting"},
	item{title: "Milk crates", desc: "Great for packing in your extra stuff"},
	item{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
	item{title: "Stickers", desc: "The thicker the vinyl the better"},
	item{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
	item{title: "Warm light", desc: "Like around 2700 Kelvin"},
	item{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
	item{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
	item{title: "Terrycloth", desc: "In other words, towel fabric"},
}

type model struct {
	state    sessionState
	viewport viewport.Model
	messages []string
	textarea textarea.Model
	list     list.Model
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			switch m.state {
			case listView:
				newitem := m.list.SelectedItem().(item)
				m.messages = append(m.messages, botStyle.Render("Bot: ")+newitem.Description())
				m.viewport.SetContent(strings.Join(m.messages, "\n"))
				m.viewport.GotoBottom()
			case textareaView:
				if m.textarea.Value() != "" {
					m.messages = append(m.messages, youStyle.Render("You: ")+m.textarea.Value())
					m.viewport.SetContent(strings.Join(m.messages, "\n"))
					m.textarea.Reset()
					m.viewport.GotoBottom()
				}
			}
		case tea.KeyTab:
			switch m.state {
			case listView:
				m.state = viewportView
			case viewportView:
				m.state = textareaView
			case textareaView:
				m.state = listView
			}
		}
		switch m.state {
		case listView:
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		case viewportView:
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		case textareaView:
			m.textarea, cmd = m.textarea.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		h, v := modelStyle.GetFrameSize()

		m.list.SetSize(getSize(30, msg.Width)-h, getSize(100, msg.Height)-v)

		m.textarea.SetWidth(getSize(70, msg.Width) - h)

		viewportHeight := msg.Height - m.textarea.Height() - v
		m.viewport.Width = getSize(70, msg.Width) - h
		m.viewport.Height = viewportHeight - v

		viewportStyle.Width(getSize(70, msg.Width) - h)
		viewportStyle.Height(viewportHeight - v)

	}

	return m, tea.Batch(cmds...)

}

func getSize(per, n int) int {
	return int(math.Floor(float64((per * n) / 100)))
}

func (m model) View() string {
	listRender := modelStyle.Render(m.list.View())
	viewportRender := viewportStyle.Render(m.viewport.View())
	textareaRender := modelStyle.Render(m.textarea.View())

	switch m.state {
	case listView:
		listRender = focusedModelStyle.Render(m.list.View())
	case viewportView:
		viewportRender = focusedModelStyle.Render(m.viewport.View())
	case textareaView:
		textareaRender = focusedModelStyle.Render(m.textarea.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		listRender,
		lipgloss.JoinVertical(lipgloss.Top,
			viewportRender,
			textareaRender,
		),
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func initialModel() model {

	ls := list.New(items, list.NewDefaultDelegate(), 0, 0)
	ls.Title = "My Fave Things"

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(100)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(100, 10)
	vp.SetContent(`Welcome to the Bubbles multi-line text input!
Try typing any message and pressing ENTER.
If you write a long message, it will automatically wrap :D
	`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		state:    listView,
		list:     ls,
		textarea: ta,
		messages: []string{},
		viewport: vp,
	}
}
