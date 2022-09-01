package app

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
	"go.mau.fi/whatsmeow/types"
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

type Item struct {
	title, desc string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

var items = []list.Item{
	Item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
	Item{title: "Nutella", desc: "It's good on toast"},
	Item{title: "Bitter melon", desc: "It cools you down"},
	Item{title: "Nice socks", desc: "And by that I mean socks without holes"},
	Item{title: "Eight hours of sleep", desc: "I had this once"},
	Item{title: "Cats", desc: "Usually"},
	Item{title: "Plantasia, the album", desc: "My plants love it too"},
	Item{title: "Pour over coffee", desc: "It takes forever to make though"},
	Item{title: "VR", desc: "Virtual reality...what is there to say?"},
	Item{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
	Item{title: "Linux", desc: "Pretty much the best OS"},
	Item{title: "Business school", desc: "Just kidding"},
	Item{title: "Pottery", desc: "Wet clay is a great feeling"},
	Item{title: "Shampoo", desc: "Nothing like clean hair"},
	Item{title: "Table tennis", desc: "It’s surprisingly exhausting"},
	Item{title: "Milk crates", desc: "Great for packing in your extra stuff"},
	Item{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
	Item{title: "Stickers", desc: "The thicker the vinyl the better"},
	Item{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
	Item{title: "Warm light", desc: "Like around 2700 Kelvin"},
	Item{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
	Item{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
	Item{title: "Terrycloth", desc: "In other words, towel fabric"},
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
				newitem := m.list.SelectedItem().(Item)
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

func parseJID(arg string) (types.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			return recipient, false
		} else if recipient.User == "" {
			return recipient, false
		}
		return recipient, true
	}
}
func initialModel() model {
	listGroup := []list.Item{}

	var conversations []Conversation
	err := Driver.Open(Conversation{}).Get().AsEntity(&conversations)
	if err != nil {
		panic(err)
	}

	for _, conv := range conversations {
		if conv.Name != nil {
			localItem := Item{title: conv.GetName(), desc: conv.GetDescription()}
			listGroup = append(listGroup, localItem)
		} else {
			jid, ok := parseJID(conv.GetId())
			if ok {
				contact, _ := Store.GetContact(jid)
				localItem := Item{title: contact.FullName, desc: contact.PushName}
				listGroup = append(listGroup, localItem)
			}
		}
	}

	ls := list.New(listGroup, list.NewDefaultDelegate(), 0, 0)
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

func NewTui() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
