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
	"log"
	"maps"
	"regexp"
	"slices"
	"sort"
	"strings"

	"gitea.joyrex.net/ejstacey/ysm/tag"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devkvlt/hexer"
)

var tags tag.Tags

const (
	tagDeleteCancelButtonId   int = 0
	tagDeleteDeleteButtonId   int = 1
	tagEntryCreateOperationId int = 0
	tagEntryModifyOperationId int = 1
	tagEntryDeleteOperationId int = 2
)

var tagModifyKeyList = map[string]key.Binding{
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
		key.WithHelp("<esc>", "back out to tag view"),
	),
	"enterKey": key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<enter>", "select button"),
	),
}

type tagModifyInputKeyMap struct {
	NextKey key.Binding
	PrevKey key.Binding
	EscKey  key.Binding
}

func (k tagModifyInputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextKey, k.PrevKey, k.EscKey}
}
func (k tagModifyInputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextKey, k.PrevKey, k.EscKey},
		{},
	}
}

func newTagModifyInputKeyMap() tagModifyInputKeyMap {
	return tagModifyInputKeyMap{
		NextKey: tagModifyKeyList["nextKey"],
		PrevKey: tagModifyKeyList["prevKey"],
		EscKey:  tagModifyKeyList["escKey"],
	}
}

type tagModifyButtonKeyMap struct {
	NextKey  key.Binding
	PrevKey  key.Binding
	EnterKey key.Binding
	EscKey   key.Binding
}

func (k tagModifyButtonKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextKey, k.PrevKey, k.EnterKey, k.EscKey}
}

func (k tagModifyButtonKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextKey, k.PrevKey, k.EscKey},
		{k.EnterKey},
	}
}

func newTagModifyButtonKeyMap() tagModifyButtonKeyMap {
	return tagModifyButtonKeyMap{
		NextKey:  tagModifyKeyList["nextKey"],
		PrevKey:  tagModifyKeyList["prevKey"],
		EnterKey: tagModifyKeyList["enterKey"],
		EscKey:   tagModifyKeyList["escKey"],
	}
}

var tagSubmenuKeyList = map[string]key.Binding{
	"nextKey": key.NewBinding(
		key.WithKeys("down", "tab", "right"),
		key.WithHelp("<down>/<tab>", "next option"),
	),
	"prevKey": key.NewBinding(
		key.WithKeys("up", "shift+tab", "left"),
		key.WithHelp("<up>/<shift-tab>", "previous option"),
	),
	"escKey": key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("<esc>", "back out to tag list"),
	),
	"enterKey": key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<enter>", "select button"),
	),
}

type tagSubmenuKeyMap struct {
	NextKey  key.Binding
	PrevKey  key.Binding
	EnterKey key.Binding
	EscKey   key.Binding
}

func (k tagSubmenuKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextKey, k.PrevKey, k.EnterKey, k.EscKey}
}

func (k tagSubmenuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextKey, k.PrevKey, k.EscKey},
		{k.EnterKey},
	}
}

func newTagSubmenuKeyMap() tagSubmenuKeyMap {
	return tagSubmenuKeyMap{
		NextKey:  tagSubmenuKeyList["nextKey"],
		PrevKey:  tagSubmenuKeyList["prevKey"],
		EnterKey: tagSubmenuKeyList["enterKey"],
		EscKey:   tagSubmenuKeyList["escKey"],
	}
}

type tagListItemDelegate struct{}

func (d tagListItemDelegate) Height() int                               { return 4 }
func (d tagListItemDelegate) Spacing() int                              { return 0 }
func (d tagListItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d tagListItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(tag.Tag)
	if !ok {
		return
	}

	var tagChannels = make(map[string]string)

	var b strings.Builder
	for i, channelId := range item.Channels() {
		tmpChannel := channels.ById()[channelId]

		var colour string
		var row = i + 1
		var col = i + 1
		if row >= 12 {
			row = 0
		}
		if col >= 18 {
			col = 0
		}
		colour = colours[col][row]

		var textColour_64 float64
		textColour, err := hexer.Invert(colour)
		if err != nil {
			log.Fatal(err)
		}
		textColour_64, err = hexer.Lightness(textColour)
		if err != nil {
			log.Fatal(err)
		}
		if textColour_64 > 50 {
			textColour = "#FFFFFF"
		} else {
			textColour = "#000000"
		}

		var style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(textColour)).
			Background(lipgloss.Color(colour))

		channelName := tmpChannel.Name()
		tagChannels[channelName] = style.Render(tmpChannel.Name())
	}

	sortedTags := slices.Sorted(maps.Keys(tagChannels))
	for i, tagName := range sortedTags {
		tagOutput := tagChannels[tagName]
		b.WriteString(tagOutput)
		if i != len(sortedTags)-1 {
			b.WriteString(", ")
		}
	}

	var out = b.String()
	if out == "" {
		out = "<none>"
	}

	var style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#" + item.FgColour())).
		Background(lipgloss.Color("#" + item.BgColour()))

	str := fmt.Sprintf("%s\n%s\n%s\n", style.Render(item.Name()), item.Description(), "channels: "+out)

	fn := blurredListStyle.Render
	if index == m.Index() {
		fn = focusedListStyle.Render
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
}

func (m Model) generateTagItems() []list.Item {
	var items []list.Item

	keys := make([]string, 0, len(m.tags.ByName()))
	for k := range m.tags.ByName() {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		var tag = m.tags.ByName()[key]
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
			t.Placeholder = "ffffff"
			t.CharLimit = 6
			t.Width = 6
			t.Prompt = " #"
			t.Validate = HexValidator
			t.SetValue(tag.FgColour())
		case 3:
			t.Placeholder = "000000"
			t.CharLimit = 6
			t.Width = 6
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
