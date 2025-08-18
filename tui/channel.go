/*
Joyrex YSM - Manager for Youtube Subscriptions
Copyright (C) 2025 Eric Stacey <ejstacey@joyrex.net>

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

*/

package tui

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"sort"
	"strings"

	"gitea.joyrex.net/ejstacey/ysm/channel"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var channels channel.Channels

var channelModifyKeyList = map[string]key.Binding{
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
		key.WithHelp("<esc>", "back out to channel view"),
	),
	"enterKey": key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<enter>", "save changes"),
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

type channelModifyNotesKeyMap struct {
	NextKey key.Binding
	PrevKey key.Binding
	EscKey  key.Binding
}

func (k channelModifyNotesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextKey, k.PrevKey, k.EscKey}
}
func (k channelModifyNotesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextKey, k.PrevKey, k.EscKey},
		{},
	}
}

func newChannelModifyNotesKeyMap() *channelModifyNotesKeyMap {
	return &channelModifyNotesKeyMap{
		NextKey: channelModifyKeyList["nextKey"],
		PrevKey: channelModifyKeyList["prevKey"],
		EscKey:  channelModifyKeyList["escKey"],
	}
}

type channelModifyTagSelectKeyMap struct {
	NextKey  key.Binding
	PrevKey  key.Binding
	SpaceKey key.Binding
	LeftKey  key.Binding
	RightKey key.Binding
	EscKey   key.Binding
}

func (k channelModifyTagSelectKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextKey, k.PrevKey, k.LeftKey, k.RightKey, k.SpaceKey, k.EscKey}
}
func (k channelModifyTagSelectKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextKey, k.PrevKey, k.EscKey},
		{k.LeftKey, k.RightKey, k.SpaceKey},
	}
}

func newChannelModifyTagSelectKeyMap() *channelModifyTagSelectKeyMap {
	return &channelModifyTagSelectKeyMap{
		NextKey:  channelModifyKeyList["nextKey"],
		PrevKey:  channelModifyKeyList["prevKey"],
		SpaceKey: channelModifyKeyList["spaceKey"],
		LeftKey:  channelModifyKeyList["leftKey"],
		RightKey: channelModifyKeyList["rightKey"],
		EscKey:   channelModifyKeyList["escKey"],
	}
}

type channelModifySubmitKeyMap struct {
	NextKey  key.Binding
	PrevKey  key.Binding
	EnterKey key.Binding
	EscKey   key.Binding
}

func (k channelModifySubmitKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextKey, k.PrevKey, k.EnterKey, k.EscKey}
}

func (k channelModifySubmitKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextKey, k.PrevKey, k.EscKey},
		{k.EnterKey},
	}
}

func newChannelModifySubmitKeyMap() *channelModifySubmitKeyMap {
	return &channelModifySubmitKeyMap{
		NextKey:  channelModifyKeyList["nextKey"],
		PrevKey:  channelModifyKeyList["prevKey"],
		EnterKey: channelModifyKeyList["enterKey"],
		EscKey:   channelModifyKeyList["escKey"],
	}
}

type channelListItemDelegate struct{}

func (d channelListItemDelegate) Height() int                               { return 4 }
func (d channelListItemDelegate) Spacing() int                              { return 0 }
func (d channelListItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d channelListItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(channel.Channel)
	if !ok {
		return
	}

	descLines := strings.Split(item.Description(), "\n")

	var channelTags = make(map[string]string)

	var b strings.Builder
	for _, tagId := range item.Tags() {
		tmpTag := tags.ById()[tagId]

		var style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#" + tmpTag.FgColour())).
			Background(lipgloss.Color("#" + tmpTag.BgColour()))

		tagName := tmpTag.Name()
		channelTags[tagName] = style.Render(tmpTag.Name())
	}

	sortedTags := slices.Sorted(maps.Keys(channelTags))
	for _, tagName := range sortedTags {
		tagOutput := channelTags[tagName]
		b.WriteString(tagOutput)
		b.WriteRune(' ')
	}

	var out = b.String()
	if out == "" {
		out = "<none>"
	}

	str := fmt.Sprintf("%s\n%s\n%s\n", item.Name(), descLines[0], "tags: "+out)

	fn := blurredListStyle.Render
	if index == m.Index() {
		fn = focusedListStyle.Render
	}

	fmt.Fprint(w, fn(str))
}

func (m Model) generateChannelItems(untaggedFilter bool) []list.Item {
	var items []list.Item

	keys := make([]string, 0, len(m.channels.ByName()))
	for k := range m.channels.ByName() {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		var channel = m.channels.ByName()[key]
		if !untaggedFilter || len(channel.Tags()) == 0 {
			items = append(items, channel)
		}
	}

	// os.WriteFile("debug-items.log", []byte(dump.Format(items)), 0644)

	return items
}

func (m Model) createChannelModifyHeader(channel channel.Channel) []string {
	channelModifyHeaders := make([]string, 2)

	channelModifyHeaders[0] = fmt.Sprintf("Name: %s\n", channel.Name())
	channelModifyHeaders[1] = fmt.Sprintf("Description: %s\n\n", channel.Description())

	return channelModifyHeaders
}

func (m Model) createChannelModifyForm(channel channel.Channel) []textinput.Model {
	channelModifyInputs := make([]textinput.Model, 1)

	t := textinput.New()
	t.Cursor.Style = cursorStyle
	t.CharLimit = 4000
	t.Placeholder = "notes about channel"
	t.SetValue(channel.Notes())
	t.Focus()
	t.PromptStyle = focusedButtonStyle
	t.TextStyle = focusedButtonStyle

	channelModifyInputs[0] = t

	return channelModifyInputs
}

func (m Model) updateChannelModifyInput(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.channelModifyInputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.channelModifyInputs {
		m.channelModifyInputs[i], cmds[i] = m.channelModifyInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
