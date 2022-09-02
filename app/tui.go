package app

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

type sessionState uint

const (
	listView sessionState = iota
	viewportView
	textareaView
)

var (
	stickerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	imageStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	audioStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	videoStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	linkStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))

	youStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	defaulStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	groupStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

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
	jid, title, desc string
}

func (i Item) JID() string         { return i.jid }
func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

type model struct {
	state    sessionState
	messages []string
	convJID  string

	list     list.Model
	viewport viewport.Model
	textarea textarea.Model
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
			return m, tea.Quit
		case tea.KeyEnter:
			switch m.state {
			case listView:
				listItem := m.list.SelectedItem().(Item)
				m.convJID = listItem.JID()
				conversation := findById(listItem.JID())
				m.messages = conversation.StylesMessages()
				m.viewport.SetContent(strings.Join(m.messages, "\n"))
				m.viewport.GotoBottom()
				// change to textarea
				m.state = textareaView
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			case textareaView:
				if m.textarea.Value() != "" {
					// Send message
					jid, err := types.ParseJID(m.convJID)
					if err != nil {
						panic(err)
					}
					Client.SendMessage(context.Background(), jid, "", &waProto.Message{
						Conversation: proto.String(m.textarea.Value()),
					})
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
				m.textarea.Blur()
			case viewportView:
				m.state = textareaView
				// change to textarea
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			case textareaView:
				m.state = listView
				m.textarea.Blur()
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
			// if !m.textarea.Focused() {
			// 	cmd = m.textarea.Focus()
			// 	cmds = append(cmds, cmd)
			// } else {
			// 	m.textarea.Blur()
			// }
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

func initialModel(first bool) model {
	if first {
		time.Sleep(time.Second * 5)
	}
	var items []list.Item
	for _, conversation := range findAll() {
		item := Item{jid: conversation.GetId(), title: conversation.Title(), desc: conversation.Desc()}
		items = append(items, item)
	}

	ls := list.New(items, list.NewDefaultDelegate(), 0, 0)
	ls.Title = "My Fave Things"

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
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

var Tea *tea.Program

func NewTui(first bool) {
	Tea = tea.NewProgram(initialModel(first), tea.WithAltScreen())
	if err := Tea.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
