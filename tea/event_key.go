package tea

import (
	"github.com/charmbracelet/bubbles/key"
)

type eventKey struct {
	Quit      key.Binding
	ForceQuit key.Binding

	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding

	CursorUp   key.Binding
	CursorDown key.Binding
	PrevPage   key.Binding
	NextPage   key.Binding

	Select key.Binding
}

func newEventKey() eventKey {
	return eventKey{
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c"),
		),

		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more help"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close help"),
		),

		CursorUp: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys("left", "pgup"),
			key.WithHelp("←/pgup", "prev page"),
		),
		NextPage: key.NewBinding(
			key.WithKeys("right", "pgdown"),
			key.WithHelp("→/pgdown", "next page"),
		),

		Select: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter", "select item"),
		),
	}
}
