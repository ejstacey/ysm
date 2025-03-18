/*
Joyrex YSM - Manager for Youtube Subscriptions
Copyright (C) 2025 Eric Stacey <ejstacey@joyrex.net>

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

*/

package tui

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gogs.joyrex.net/ejstacey/ysm/channel"
	"gogs.joyrex.net/ejstacey/ysm/tag"
)

var channels channel.Channels

type tagListItemDelegate struct{}

func (d tagListItemDelegate) Height() int                               { return 4 }
func (d tagListItemDelegate) Spacing() int                              { return 0 }
func (d tagListItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d tagListItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(tag.Tag)
	if !ok {
		return
	}

	if len(channels.ById) == 0 {
		channels.LoadEntriesFromDb()
	}

	var tmpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#25A065"))

	var b strings.Builder
	for _, channelId := range item.Channels() {
		tmpChan := channels.ById[channelId]

		fmt.Fprintf(&b, "%s", tmpStyle.Render(tmpChan.Name()))
	}
	// err := os.WriteFile("debug.log", []byte(dump.Format(channels)), 0644)
	// if err != nil {
	// 	panic(err)
	// }

	var out = b.String()
	if out == "" {
		out = "no channels"
	}

	str := fmt.Sprintf("%s\n%s\n%s\n", item.Name(), item.Description(), out)

	fn := blurredStyle.Render
	if index == m.Index() {
		fn = focusedStyle.Render
	}

	fmt.Fprint(w, fn(str))
}

// Validator functions to ensure valid input
func HexValidator(s string) error {
	var res, err = regexp.Match(`^[0-9A-Fa-f]*$`, []byte(s))
	if err != nil {
		panic(err)
	}
	if !res {
		return errors.New(s + " has invalid characters")
	} else {
		return nil
	}
	// errDebug := os.WriteFile("debug.log", []byte(dump.Format(err)), 0644)
	// if errDebug != nil {
	// 	panic(errDebug)
	// }
}

func (m Model) generateTagItems() []list.Item {
	var items []list.Item

	keys := make([]string, 0, len(m.tags.ByName))
	for k := range m.tags.ByName {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		var tag = m.tags.ByName[key]
		items = append(items, tag)
	}

	return items
}

func (m Model) createTagEntryForm(tag tag.Tag) []textinput.Model {
	tagEntryInputs := make([]textinput.Model, 4)

	var t textinput.Model
	for i := range tagEntryInputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.PromptStyle = blurredStyle
		t.TextStyle = blurredStyle

		switch i {
		case 0:
			t.Placeholder = "funny"
			t.CharLimit = 64
			t.SetValue(tag.Name())
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "channels that are funny"
			t.CharLimit = 4000
			t.SetValue(tag.Description())
		case 2:
			t.Placeholder = "foreground colour (hex)"
			t.CharLimit = 7
			t.Width = 50
			t.Prompt = " #"
			t.Validate = HexValidator
			t.SetValue(tag.FgColour())
		case 3:
			t.Placeholder = "background color (hex)"
			t.CharLimit = 7
			t.Width = 50
			t.Prompt = " #"
			t.Validate = HexValidator
			t.SetValue(tag.BgColour())
		}

		tagEntryInputs[i] = t
	}
	return tagEntryInputs
}

func (m Model) deleteTagEntryForm(tag tag.Tag) []string {
	tagDeleteInputs := make([]string, 2)

	tagDeleteInputs[0] = fmt.Sprintf("Name: %s\n", tag.Name())
	tagDeleteInputs[1] = fmt.Sprintf("Description: %s\n\n", tag.Description())

	return tagDeleteInputs
}

func (m Model) updateTagEntryInput(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.tagEntryInputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.tagEntryInputs {
		m.tagEntryInputs[i], cmds[i] = m.tagEntryInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
