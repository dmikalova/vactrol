// Command vactrol is the terminal UI for the Vactrol card game. It opens a menu
// to browse the card database (card explorer) or play an interactive two-player
// hotseat match, all in one program.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	snd := &sender{}
	p := tea.NewProgram(newRootModel(snd), tea.WithAltScreen())
	snd.p = p // let the game's chooser bridge post messages into this program
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "vactrol:", err)
		os.Exit(1)
	}
}

// sender lets background work (the game's chooser, running off the UI thread)
// inject messages into the running program. Its handle is wired up in main once
// the program exists.
type sender struct{ p *tea.Program }

func (s *sender) send(msg tea.Msg) {
	if s.p != nil {
		s.p.Send(msg)
	}
}

// screen identifies which sub-model is active.
type screen int

const (
	screenMenu screen = iota
	screenExplorer
	screenStats
	screenGame
)

// switchMsg asks the root model to change screens.
type switchMsg struct{ to screen }

func gotoScreen(s screen) tea.Cmd {
	return func() tea.Msg { return switchMsg{to: s} }
}

// rootModel owns the active screen and routes messages to the matching sub-model.
type rootModel struct {
	snd      *sender
	screen   screen
	menu     menuModel
	explorer explorerModel
	stats    statsModel
	game     gameModel
	width    int
	height   int
}

func newRootModel(snd *sender) rootModel {
	return rootModel{snd: snd, screen: screenMenu, menu: newMenuModel()}
}

func (m rootModel) Init() tea.Cmd { return nil }

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.menu = m.menu.resize(m.width, m.height)
		m.explorer = m.explorer.resize(m.width, m.height)
		m.stats = m.stats.resize(m.width, m.height)
		m.game = m.game.resize(m.width, m.height)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case switchMsg:
		m.screen = msg.to
		switch msg.to {
		case screenExplorer:
			m.explorer = newExplorerModel(m.width, m.height)
		case screenStats:
			m.stats = newStatsModel(m.width, m.height)
		case screenGame:
			m.game = newGameModel(m.snd, m.width, m.height)
			return m, m.game.Init()
		}
		return m, nil
	}

	var cmd tea.Cmd
	switch m.screen {
	case screenMenu:
		m.menu, cmd = m.menu.Update(msg)
	case screenExplorer:
		m.explorer, cmd = m.explorer.Update(msg)
	case screenStats:
		m.stats, cmd = m.stats.Update(msg)
	case screenGame:
		m.game, cmd = m.game.Update(msg)
	}
	return m, cmd
}

func (m rootModel) View() string {
	switch m.screen {
	case screenExplorer:
		return m.explorer.View()
	case screenStats:
		return m.stats.View()
	case screenGame:
		return m.game.View()
	default:
		return m.menu.View()
	}
}
