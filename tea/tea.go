package tea

import (
	"errors"

	"github.com/charmbracelet/bubbletea"

	"github.com/starudream/go-lib/core/v2/utils/signalutil"
)

func Run[T tea.Model](t T) (T, error) {
	p := tea.NewProgram(t, tea.WithAltScreen(), tea.WithoutSignalHandler())
	go func() { <-signalutil.Defer(p.Kill).Done() }()
	m, err := p.Run()
	return m.(T), err
}

func Return(err error) error {
	if err == nil || errors.Is(err, tea.ErrProgramKilled) {
		return nil
	}
	return err
}
