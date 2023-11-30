package tea

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/starudream/go-lib/core/v2/gh"

	"github.com/starudream/github-webhook/github"
)

type Event struct {
	items []github.Event

	keys   eventKey
	styles eventStyle

	paginator paginator.Model
	help      help.Model

	title       string
	status      string
	statusTimer *time.Timer
	cursor      int
	selected    map[int]struct{}

	width      int
	height     int
	itemHeight int
	itemSpace  int
}

func NewEvent() *Event {
	t := &Event{}
	t.items = github.Events

	t.keys = newEventKey()
	t.styles = newEventStyle()

	t.paginator = paginator.New()
	t.paginator.Type = paginator.Dots
	t.paginator.ActiveDot = t.styles.ActivePaginationDot.String()
	t.paginator.InactiveDot = t.styles.InactivePaginationDot.String()
	t.paginator.KeyMap = paginator.KeyMap{
		PrevPage: t.keys.PrevPage,
		NextPage: t.keys.NextPage,
	}

	t.help = help.New()

	t.title = "GitHub WebHook"

	t.selected = map[int]struct{}{}

	t.itemHeight = 1
	t.itemSpace = 1

	t.updatePaginator()

	return t
}

var _ tea.Model = (*Event)(nil)

func (t *Event) Init() tea.Cmd {
	return nil
}

func (t *Event) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch m := msg.(type) {
	case eventStatusTimeout:
		t.status = ""
		if t.statusTimer != nil {
			t.statusTimer.Stop()
		}
	case tea.WindowSizeMsg:
		h, v := t.styles.Screen.GetFrameSize()
		t.SetSize(m.Width-h, m.Height-v)
	case tea.KeyMsg:
		switch {
		case key.Matches(m, t.keys.Quit, t.keys.ForceQuit):
			return t, tea.Quit

		case key.Matches(m, t.keys.ShowFullHelp, t.keys.CloseFullHelp):
			t.help.ShowAll = !t.help.ShowAll
			t.updatePaginator()

		case key.Matches(m, t.keys.CursorUp):
			t.CursorUp()
		case key.Matches(m, t.keys.CursorDown):
			t.CursorDown()
		case key.Matches(m, t.keys.PrevPage):
			t.paginator.PrevPage()
		case key.Matches(m, t.keys.NextPage):
			t.paginator.NextPage()

		case key.Matches(m, t.keys.Select):
			idx := t.Index()
			item := t.items[idx]
			_, selected := t.selected[idx]
			if selected {
				delete(t.selected, idx)
			} else {
				t.selected[idx] = struct{}{}
			}
			status := fmt.Sprintf("%q %s, %d selected in total.", item.Name, gh.Ternary(selected, "deselected", "selected"), len(t.selected))
			cmds = append(cmds, t.NewStatus(status))
		}
	}

	current := t.paginator.ItemsOnPage(len(t.items))
	if t.cursor > current-1 {
		t.cursor = max(0, current-1)
	}

	return t, tea.Batch(cmds...)
}

func (t *Event) View() string {
	var (
		sections []string
		height   = t.height
	)

	sTitle := t.renderTitle()
	sections = append(sections, sTitle)
	height -= lipgloss.Height(sTitle)

	sDescription := t.renderDescription()
	height -= lipgloss.Height(sDescription)

	sPaginator := t.renderPaginator()
	height -= lipgloss.Height(sPaginator)

	sHelp := t.renderHelp()
	height -= lipgloss.Height(sHelp)

	height -= t.styles.List.GetVerticalMargins() + t.styles.List.GetVerticalBorderSize()

	sList := t.renderList(height)
	sections = append(sections, sList)

	sections = append(sections, sDescription)
	sections = append(sections, sPaginator)
	sections = append(sections, sHelp)

	return t.styles.Screen.Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

var _ help.KeyMap = (*Event)(nil)

func (t *Event) ShortHelp() []key.Binding {
	return []key.Binding{
		t.keys.CursorUp,
		t.keys.CursorDown,
		t.keys.PrevPage,
		t.keys.NextPage,

		t.keys.ShowFullHelp,

		t.keys.Quit,
	}
}

func (t *Event) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			t.keys.CursorUp,
			t.keys.CursorDown,
			t.keys.PrevPage,
			t.keys.NextPage,
		},
		{
			t.keys.CloseFullHelp,
		},
		{
			t.keys.Quit,
		},
	}
}

func (t *Event) SetSize(width, height int) {
	t.width = width
	t.height = height

	t.help.Width = width

	t.updatePaginator()
}

func (t *Event) Index() int {
	return t.paginator.Page*t.paginator.PerPage + t.cursor
}

func (t *Event) Cursor() int {
	return t.cursor
}

func (t *Event) CursorUp() {
	t.cursor--

	if t.cursor < 0 && t.paginator.Page == 0 {
		t.cursor = 0
	}

	if t.cursor >= 0 {
		return
	}

	t.paginator.PrevPage()
	t.cursor = t.paginator.ItemsOnPage(len(t.items)) - 1
}

func (t *Event) CursorDown() {
	current := t.paginator.ItemsOnPage(len(t.items))

	t.cursor++

	if t.cursor < current {
		return
	}

	if !t.paginator.OnLastPage() {
		t.paginator.NextPage()
		t.cursor = 0
		return
	}

	t.cursor = current - 1
}

func (t *Event) updatePaginator() {
	index := t.Index()

	height := t.height
	height -= lipgloss.Height(t.renderTitle())
	height -= t.styles.List.GetVerticalPadding() + t.styles.List.GetVerticalBorderSize()
	height -= lipgloss.Height(t.renderDescription())
	height -= lipgloss.Height(t.renderPaginator())
	height -= lipgloss.Height(t.renderHelp())

	t.paginator.PerPage = max(1, height/(t.itemHeight+t.itemSpace))
	t.paginator.SetTotalPages(len(t.items))
	t.paginator.Page = index / t.paginator.PerPage
	t.cursor = index % t.paginator.PerPage
}

func (t *Event) renderTitle() string {
	return t.styles.Title.Render(t.title) + "  " + t.status
}

func (t *Event) renderPaginator() string {
	arabic := fmt.Sprintf("%d/%d", t.paginator.Page+1, t.paginator.TotalPages)
	return t.styles.Pagination.Render(arabic + "  " + t.paginator.View())
}

func (t *Event) renderDescription() string {
	item := t.items[t.Index()]
	actions := make([]string, len(item.Actions))
	for i, action := range item.Actions {
		actions[i] = action.Action
	}
	block := fmt.Sprintf(""+
		"Event:       %s\n"+
		"Description: %s\n"+
		"Actions:     %s",
		item.Name, item.Desc, strings.Join(actions, ", "),
	)
	return t.styles.Description.Height(3).Render(block)
}

func (t *Event) renderList(height int) string {
	buf := &bytes.Buffer{}

	start, end := t.paginator.GetSliceBounds(len(t.items))
	items := t.items[start:end]
	for i, item := range items {
		buf.WriteString(t.renderItem(i+start, item))
		if i < len(items)-1 {
			buf.WriteString(strings.Repeat("\n", 1+t.itemSpace))
		}
	}

	current := t.paginator.ItemsOnPage(len(t.items))
	if current < t.paginator.PerPage {
		n := (t.paginator.PerPage - current) * (t.itemHeight + t.itemSpace)
		buf.WriteString(strings.Repeat("\n", n))
	}

	return t.styles.List.Height(height).Render(buf.String())
}

func (t *Event) renderItem(i int, item github.Event) string {
	current := t.Index() == i
	_, selected := t.selected[i]

	checkbox := "(   )"
	if selected {
		checkbox = "( * )"
	}

	itemStyle := t.styles.NormalItem
	if current {
		itemStyle = t.styles.CurrentItem
	} else if selected {
		itemStyle = t.styles.SelectedItem
	}

	itemKey := itemStyle.Copy().Bold(true).Render(item.Key)

	block := itemStyle.Render(checkbox+" "+item.Name+" [ ") + itemKey + itemStyle.Render(" ]")

	return lipgloss.NewStyle().Height(t.itemHeight).Render(block)
}

func (t *Event) renderHelp() string {
	return t.styles.Help.Render(t.help.View(t))
}

type eventStatusTimeout struct{}

func (t *Event) NewStatus(s string) tea.Cmd {
	t.status = s
	if t.statusTimer != nil {
		t.statusTimer.Stop()
	}
	t.statusTimer = time.NewTimer(time.Second)
	return func() tea.Msg {
		<-t.statusTimer.C
		return eventStatusTimeout{}
	}
}
