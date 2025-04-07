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
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

func (m Model) generateChannelItems() []list.Item {
	var items []list.Item

	keys := make([]string, 0, len(m.channels.ByName()))
	for k := range m.channels.ByName() {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		var channel = m.channels.ByName()[key]
		items = append(items, channel)
	}

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
