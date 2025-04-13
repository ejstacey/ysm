/*
Joyrex YSM - Manager for Youtube Subscriptions
Copyright (C) 2025 Eric Stacey <ejstacey@joyrex.net>

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

*/

package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var generatePageKeyList = map[string]key.Binding{
	"nextKey": key.NewBinding(
		key.WithKeys("down", "tab"),
		key.WithHelp("<down>/<tab>", "next option"),
	),
	"prevKey": key.NewBinding(
		key.WithKeys("up", "shift+tab"),
		key.WithHelp("<up>/<shift-tab>", "previous option"),
	),
	"escKey": key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("<esc>", "back out to previous page"),
	),
	"enterKey": key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<enter>", "select button"),
	),
	"leftKey": key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("<left>", "move tag selector to left"),
	),
	"rightKey": key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("<right>", "move tag selector to right"),
	),
	"spaceKey": key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("<space>", "select or unselect tag"),
	),
}

type generatePageInputKeyMap struct {
	NextKey  key.Binding
	PrevKey  key.Binding
	SpaceKey key.Binding
	LeftKey  key.Binding
	RightKey key.Binding
	EscKey   key.Binding
}

func (k generatePageInputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextKey, k.PrevKey, k.LeftKey, k.RightKey, k.SpaceKey, k.EscKey}
}
func (k generatePageInputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextKey, k.PrevKey, k.EscKey},
		{k.LeftKey, k.RightKey, k.SpaceKey},
	}
}

// func newGeneratePageInputKeyMap() generatePageInputKeyMap {
// 	return generatePageInputKeyMap{
// 		NextKey: generatePageKeyList["nextKey"],
// 		PrevKey: generatePageKeyList["prevKey"],
// 		EscKey:  generatePageKeyList["escKey"],
// 	}
// }

type generatePageButtonKeyMap struct {
	NextKey  key.Binding
	PrevKey  key.Binding
	EnterKey key.Binding
	EscKey   key.Binding
}

func (k generatePageButtonKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextKey, k.PrevKey, k.EnterKey, k.EscKey}
}

func (k generatePageButtonKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextKey, k.PrevKey, k.EscKey},
		{k.EnterKey},
	}
}

// func newGeneratePageButtonKeyMap() generatePageButtonKeyMap {
// 	return generatePageButtonKeyMap{
// 		NextKey:  generatePageKeyList["nextKey"],
// 		PrevKey:  generatePageKeyList["prevKey"],
// 		EnterKey: generatePageKeyList["enterKey"],
// 		EscKey:   generatePageKeyList["escKey"],
// 	}
// }

func (m Model) createGeneratePageForm() []textinput.Model {
	generatePageInputs := make([]textinput.Model, 3)

	var t textinput.Model
	for i := range generatePageInputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.PromptStyle = blurredStyle
		t.TextStyle = blurredStyle

		switch i {
		case 0:
			t.Placeholder = "location of template file relative to running directory (or absolute path)"
			t.CharLimit = 256
			t.SetValue("templates/default.tmpl")
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "location of output file relative to running directory (or absolute path)"
			t.CharLimit = 256
			t.SetValue("html/index.html")
		case 2:
			t.Placeholder = "title for the page"
			t.CharLimit = 512
			t.Width = 512
			t.SetValue("My Youtube Subscriptions")
		}

		generatePageInputs[i] = t
	}
	return generatePageInputs
}

func (m Model) updateGeneratePageInput(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.generatePageInputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.generatePageInputs {
		m.generatePageInputs[i], cmds[i] = m.generatePageInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
