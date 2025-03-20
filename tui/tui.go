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
	"log"
	"maps"
	"math"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/devkvlt/hexer"
	"gogs.joyrex.net/ejstacey/ysm/channel"
	"gogs.joyrex.net/ejstacey/ysm/tag"
	"gogs.joyrex.net/ejstacey/ysm/utils"
)

var (
	// P the current tea program
	P *tea.Program
	// WindowSize store the size of the terminal window
	WindowSize tea.WindowSizeMsg

	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	unsavedColour = lipgloss.Color("#0077FF")
	activeColour  = lipgloss.Color("#0000FF")

	focusedListStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	blurredListStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#5dade2"))

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(activeColour)
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	focusedButtonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(activeColour)
	blurredButtonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	cursorStyle = focusedStyle
	// noStyle     = lipgloss.NewStyle()

	focusedColourPickerModelStyle = lipgloss.NewStyle().
					Width(4).
					Height(1).
					Align(lipgloss.Center, lipgloss.Center)

	colourPickerModelStyle = lipgloss.NewStyle().
				Width(4).
				Height(1).
				Align(lipgloss.Center, lipgloss.Center).
				SetString("    ")

	tagDisplayStyle = lipgloss.NewStyle().
			Height(1).
			Align(lipgloss.Center, lipgloss.Center)

	colours = [18][12]string{
		{
			"#CCFF00",
			"#CCFF33",
			"#CCFF66",
			"#CCFF99",
			"#CCFFCC",
			"#CCFFFF",
			"#FFFFFF",
			"#FFFFCC",
			"#FFFF99",
			"#FFFF66",
			"#FFFF33",
			"#FFFF00",
		},
		{
			"#CCCC00",
			"#CCCC33",
			"#CCCC66",
			"#CCCC99",
			"#CCCCCC",
			"#CCCCFF",
			"#FFCCFF",
			"#FFCCCC",
			"#FFCC99",
			"#FFCC66",
			"#FFCC33",
			"#FFCC00",
		},
		{
			"#CC9900",
			"#CC9933",
			"#CC9966",
			"#CC9999",
			"#CC99CC",
			"#CC99FF",
			"#FF99FF",
			"#FF99CC",
			"#FF9999",
			"#FF9966",
			"#FF9933",
			"#FF9900",
		},
		{
			"#CC6600",
			"#CC6633",
			"#CC6666",
			"#CC6699",
			"#CC66CC",
			"#CC66FF",
			"#FF66FF",
			"#FF66CC",
			"#FF6699",
			"#FF6666",
			"#FF6633",
			"#FF6600",
		},
		{
			"#CC3300",
			"#CC3333",
			"#CC3366",
			"#CC3399",
			"#CC33CC",
			"#CC33FF",
			"#FF33FF",
			"#FF33CC",
			"#FF3399",
			"#FF3366",
			"#FF3333",
			"#FF3300",
		},
		{
			"#CC0000",
			"#CC0033",
			"#CC0066",
			"#CC0099",
			"#CC00CC",
			"#CC00FF",
			"#FF00FF",
			"#FF00CC",
			"#FF0099",
			"#FF0066",
			"#FF0033",
			"#FF0000",
		},
		{
			"#660000",
			"#660033",
			"#660066",
			"#660099",
			"#6600CC",
			"#6600FF",
			"#9900FF",
			"#9900CC",
			"#990099",
			"#990066",
			"#990033",
			"#990000",
		},
		{
			"#663300",
			"#663333",
			"#663366",
			"#663399",
			"#6633CC",
			"#6633FF",
			"#9933FF",
			"#9933CC",
			"#993399",
			"#993366",
			"#993333",
			"#993300",
		},
		{
			"#666600",
			"#666633",
			"#666666",
			"#666699",
			"#6666CC",
			"#6666FF",
			"#9966FF",
			"#9966CC",
			"#996699",
			"#996666",
			"#996633",
			"#996600",
		},
		{
			"#669900",
			"#669933",
			"#669966",
			"#669999",
			"#6699CC",
			"#6699FF",
			"#9999FF",
			"#9999CC",
			"#999999",
			"#999966",
			"#999933",
			"#999900",
		},
		{
			"#66CC00",
			"#66CC33",
			"#66CC66",
			"#66CC99",
			"#66CCCC",
			"#66CCFF",
			"#99CCFF",
			"#99CCCC",
			"#99CC99",
			"#99CC66",
			"#99CC33",
			"#99CC00",
		},
		{
			"#66FF00",
			"#66FF33",
			"#66FF66",
			"#66FF99",
			"#66FFCC",
			"#66FFFF",
			"#99FFFF",
			"#99FFCC",
			"#99FF99",
			"#99FF66",
			"#99FF33",
			"#99FF00",
		},
		{
			"#00FF00",
			"#00FF33",
			"#00FF66",
			"#00FF99",
			"#00FFCC",
			"#00FFFF",
			"#33FFFF",
			"#33FFCC",
			"#33FF99",
			"#33FF66",
			"#33FF33",
			"#33FF00",
		},
		{
			"#00CC00",
			"#00CC33",
			"#00CC66",
			"#00CC99",
			"#00CCCC",
			"#00CCFF",
			"#33CCFF",
			"#33CCCC",
			"#33CC99",
			"#33CC66",
			"#33CC33",
			"#33CC00",
		},
		{
			"#009900",
			"#009933",
			"#009966",
			"#009999",
			"#0099CC",
			"#0099FF",
			"#3399FF",
			"#3399CC",
			"#339999",
			"#339966",
			"#339933",
			"#339900",
		},
		{
			"#006600",
			"#006633",
			"#006666",
			"#006699",
			"#0066CC",
			"#0066FF",
			"#3366FF",
			"#3366CC",
			"#336699",
			"#336666",
			"#336633",
			"#336600",
		},
		{
			"#003300",
			"#003333",
			"#003366",
			"#003399",
			"#0033CC",
			"#0033FF",
			"#3333FF",
			"#3333CC",
			"#333399",
			"#333366",
			"#333333",
			"#333300",
		},
		{
			"#000000",
			"#000033",
			"#000066",
			"#000099",
			"#0000CC",
			"#0000FF",
			"#3300FF",
			"#3300CC",
			"#330099",
			"#330066",
			"#330033",
			"#330000",
		},
	}

	keyList = map[string]key.Binding{
		"cKey": key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "channel view"),
		),
		"tKey": key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "tag view"),
		),
		"pKey": key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "toggle pagination"),
		),
		"hKey": key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "toggle help"),
		),
		"qKey": key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		"nKey": key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new entry"),
		),
		"dKey": key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete entry"),
		),
		"mKey": key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "modify entry"),
		),
		"tabKey": key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("<tab>", "next option"),
		),
		"shiftTabKey": key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("<ahift-tab>", "previous option"),
		),
		"enterKey": key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("<enter>", "select entry"),
		),
		"upKey": key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("<up>", "previous option"),
		),
		"downKey": key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("<down>", "next option"),
		),
		"leftKey": key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("<left>", "previous option"),
		),
		"rightKey": key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("<right>", "next option"),
		),
		"escKey": key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("<esc>", "back out"),
		),
		"spaceKey": key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("<space>", "select item"),
		),
	}

	channels channel.Channels
	tags     tag.Tags
)

const (
	tagDeleteCancelButtonId   int = 0
	tagDeleteDeleteButtonId   int = 1
	tagEntryCreateOperationId int = 0
	tagEntryModifyOperationId int = 1
	tagEntryDeleteOperationId int = 2
)

type channelModifyKeyMap struct {
	spaceKey    key.Binding
	tabKey      key.Binding
	shiftTabKey key.Binding
	enterKey    key.Binding
	upKey       key.Binding
	downKey     key.Binding
	leftKey     key.Binding
	rightKey    key.Binding
	escKey      key.Binding
}

type listKeyMap struct {
	cKey        key.Binding
	dKey        key.Binding
	tKey        key.Binding
	pKey        key.Binding
	mKey        key.Binding
	hKey        key.Binding
	qKey        key.Binding
	nKey        key.Binding
	tabKey      key.Binding
	shiftTabKey key.Binding
	enterKey    key.Binding
	upKey       key.Binding
	downKey     key.Binding
	leftKey     key.Binding
	rightKey    key.Binding
	escKey      key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		cKey:        keyList["cKey"],
		tKey:        keyList["tKey"],
		pKey:        keyList["pKey"],
		hKey:        keyList["hKey"],
		qKey:        keyList["qKey"],
		nKey:        keyList["nKey"],
		dKey:        keyList["dKey"],
		mKey:        keyList["mKey"],
		tabKey:      keyList["tabKey"],
		shiftTabKey: keyList["shiftTabKey"],
		enterKey:    keyList["enterKey"],
		upKey:       keyList["upKey"],
		downKey:     keyList["downKey"],
		leftKey:     keyList["leftKey"],
		rightKey:    keyList["rightKey"],
		escKey:      keyList["escKey"],
	}
}

type Model struct {
	current              string
	list                 list.Model
	channels             channel.Channels
	tags                 tag.Tags
	keys                 *listKeyMap
	help                 help.Model // this doesn't get used on list pages. lists have their own built-in help
	selectedChannel      channel.Channel
	selectedTag          tag.Tag
	selectedChannelId    int
	selectedTagId        int
	selectedTagIds       []int
	tagEntryFocus        int
	tagEntryOperation    int
	tagDeleteFocus       int
	tagDeleteInputs      []string
	tagEntryInputs       []textinput.Model
	channelModifyFocus   int
	channelModifyHeaders []string
	channelModifyInputs  []textinput.Model
	colourPickerX        int
	colourPickerY        int
	selectedBackColour   string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// handle list pages
	switch m.current {
	case "channel", "tag":
		switch m.current {
		case "channel":
			m.selectedChannelId = m.list.Index()
		case "tag":
			m.selectedTagId = m.list.Index()
		}

		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			h, v := appStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)

		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.cKey):
				m.current = "channel"
				width := m.list.Width()
				height := m.list.Height()
				m.list = list.New(m.generateChannelItems(), channelListItemDelegate{}, width, height-1)
				m.list.Title = "YSM - Channel View"
				m.list.Styles.Title = titleStyle
				listKeys := newListKeyMap()
				m.list.AdditionalShortHelpKeys = func() []key.Binding {
					return []key.Binding{
						listKeys.tKey,
						listKeys.pKey,
						listKeys.enterKey,
					}
				}
				m.list.AdditionalFullHelpKeys = func() []key.Binding {
					return []key.Binding{
						listKeys.tKey,
						listKeys.pKey,
						listKeys.enterKey,
					}
				}

				// set selected channel
				if m.current == "channel" {
					if m.selectedChannelId != -1 {
						m.list.Select(m.selectedChannelId)
					}
				}

				return m, nil

			case key.Matches(msg, m.keys.tKey):
				m.current = "tag"
				width := m.list.Width()
				height := m.list.Height()

				m.list = list.New(m.generateTagItems(), tagListItemDelegate{}, width, height-1)
				m.list.Title = "YSM - Tag View"
				m.list.Styles.Title = titleStyle

				listKeys := newListKeyMap()
				m.list.AdditionalShortHelpKeys = func() []key.Binding {
					return []key.Binding{
						listKeys.cKey,
						listKeys.pKey,
						listKeys.nKey,
						listKeys.mKey,
						listKeys.dKey,
						listKeys.enterKey,
					}
				}
				m.list.AdditionalFullHelpKeys = func() []key.Binding {
					return []key.Binding{
						listKeys.cKey,
						listKeys.pKey,
						listKeys.nKey,
						listKeys.mKey,
						listKeys.dKey,
						listKeys.enterKey,
					}
				}

				// set selected tag
				if m.current == "tag" {
					if m.selectedTagId != -1 {
						m.list.Select(m.selectedTagId)
					}
				}
				return m, nil

			case key.Matches(msg, m.keys.nKey):
				if m.current != "tag" {
					return m, nil
				}

				var tag tag.Tag
				m.tagEntryOperation = tagEntryCreateOperationId
				m.tagEntryInputs = m.createTagEntryForm(tag)
				m.current = "tagEntry"
				return m, nil

			case key.Matches(msg, m.keys.pKey):
				switch m.current {
				case "channel", "tag":
					m.list.SetShowPagination(!m.list.ShowPagination())
					return m, nil
				}

			case key.Matches(msg, m.keys.hKey):
				switch m.current {
				case "channel", "tag":
					m.list.SetShowHelp(!m.list.ShowHelp())
					return m, nil
				}

			case key.Matches(msg, m.keys.enterKey, m.keys.mKey):
				switch m.current {
				case "channel":
					channel := m.list.SelectedItem().(channel.Channel)
					m.channelModifyHeaders = m.createChannelModifyHeader(channel)
					m.channelModifyInputs = m.createChannelModifyForm(channel)
					m.selectedChannel = channel
					m.current = "channelModify"
					// m.help = help.New()
					// m.help
					return m, nil
				case "tag":
					tag := m.list.SelectedItem().(tag.Tag)
					m.tagEntryOperation = tagEntryModifyOperationId
					m.tagEntryInputs = m.createTagEntryForm(tag)
					m.selectedTag = tag
					m.current = "tagEntry"
					return m, nil
				}

			case key.Matches(msg, m.keys.dKey):
				switch m.current {
				case "tag":
					tag := m.list.SelectedItem().(tag.Tag)
					m.tagEntryOperation = tagEntryDeleteOperationId
					m.tagDeleteInputs = m.deleteTagEntryForm(tag)
					m.current = "confirmDelete"
					return m, nil
				}

			case key.Matches(msg, m.keys.qKey):
				return m, tea.Quit
			}
		}

		m.list, cmd = m.list.Update(msg)

	case "tagEntry":
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			h, v := appStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)

		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.escKey):
				m.current = "tag"
				return m, nil

			case key.Matches(msg, m.keys.tabKey, m.keys.shiftTabKey, m.keys.enterKey, m.keys.upKey, m.keys.downKey):
				s := msg.String()

				// Did the user press enter while the submit button was focused?
				// If so, create it.
				if s == "enter" {
					if m.tagEntryFocus == (len(m.tagEntryInputs) + 2) {
						if m.tagEntryOperation == tagEntryCreateOperationId {
							var tag tag.Tag
							err := tag.New()
							utils.HandleError(err, "creating new tag")
							err = tag.SetName(m.tagEntryInputs[0].Value())
							utils.HandleError(err, "updating tag name")
							err = tag.SetDescription(m.tagEntryInputs[1].Value())
							utils.HandleError(err, "updating tag description")
							err = tag.SetFgColour(m.tagEntryInputs[2].Value())
							utils.HandleError(err, "updating tag fgcolour")
							err = tag.SetBgColour(m.tagEntryInputs[3].Value())
							utils.HandleError(err, "updating tag bgColour")
						} else {
							tag := m.list.SelectedItem().(tag.Tag)
							err := tag.SetName(m.tagEntryInputs[0].Value())
							utils.HandleError(err, "updating tag name")
							err = tag.SetDescription(m.tagEntryInputs[1].Value())
							utils.HandleError(err, "updating tag description")
							err = tag.SetFgColour(m.tagEntryInputs[2].Value())
							utils.HandleError(err, "updating tag fgcolour")
							err = tag.SetBgColour(m.tagEntryInputs[3].Value())
							utils.HandleError(err, "updating tag bgColour")
						}

						m.tags.LoadEntriesFromDb()
						m.list.SetItems(m.generateTagItems())
						m.current = "tag"
						return m, nil
					} else if m.tagEntryFocus == 3 || m.tagEntryFocus == 5 {
						m.current = "colourPicker"
						if m.tagEntryFocus == 3 {
							m.tagEntryInputs[2].SetValue(strings.ToUpper(m.tagEntryInputs[2].Value()))
							m.selectedBackColour = "#" + m.tagEntryInputs[2].Value()
							var found bool
							found = false
							for row := 0; row < 12; row++ {
								for column := 0; column < len(colours); column++ {
									// colours[m.colourPickerX][m.colourPickerY]
									if colours[column][row] == m.selectedBackColour {
										found = true
										m.colourPickerX = column
										m.colourPickerY = row
									}
								}
							}
							if !found {
								m.colourPickerX = 0
								m.colourPickerY = 0
								m.selectedBackColour = colours[0][0]
							}
						} else {
							m.tagEntryInputs[3].SetValue(strings.ToUpper(m.tagEntryInputs[3].Value()))
							m.selectedBackColour = "#" + m.tagEntryInputs[3].Value()
							var found bool
							found = false
							for row := 0; row < 12; row++ {
								for column := 0; column < len(colours); column++ {
									// colours[m.colourPickerX][m.colourPickerY]
									if colours[column][row] == m.selectedBackColour {
										found = true
										m.colourPickerX = column
										m.colourPickerY = row
									}
								}
							}
							if !found {
								m.colourPickerX = 0
								m.colourPickerY = 0
								m.selectedBackColour = colours[0][0]
							}
						}
						return m, nil
					}
				}

				// Cycle indexes
				if s == "up" || s == "shift+tab" {
					m.tagEntryFocus--
				} else {
					m.tagEntryFocus++
				}

				if m.tagEntryFocus > (len(m.tagEntryInputs) + 2) {
					m.tagEntryFocus = 0
				} else if m.tagEntryFocus < 0 {
					m.tagEntryFocus = (len(m.tagEntryInputs) + 2)
				}

				cmds := make([]tea.Cmd, (len(m.tagEntryInputs) + 2))
				for i := 0; i <= len(m.tagEntryInputs)-1+2; i++ {
					var j int
					// laying this out like this because I'm currently too stoned to
					// keep the different permutations in my head
					if i == 0 {
						j = 0
					} else if i == 1 {
						j = 1
					} else if i == 2 {
						j = 2
						// } else if i == 3 {
						// 	j = 2
					} else if i == 4 {
						j = 3
						// } else if i == 5 {
						// 	j = 3
					}
					if i == m.tagEntryFocus && i != 3 && i != 5 {
						// Set focused state
						cmds[j] = m.tagEntryInputs[j].Focus()
						m.tagEntryInputs[j].PromptStyle = focusedStyle
						m.tagEntryInputs[j].TextStyle = focusedStyle
						continue
					} else {
						// Remove focused state
						if i != 3 && i != 5 {
							m.tagEntryInputs[j].Blur()
							m.tagEntryInputs[j].PromptStyle = blurredStyle
							m.tagEntryInputs[j].TextStyle = blurredStyle
						}
					}
				}

				return m, tea.Batch(cmds...)
			default:
				cmd = m.updateTagEntryInput(msg)
			}
		}

	case "channelModify":
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			h, v := appStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)

		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.escKey):
				m.current = "channel"
				return m, nil

			case key.Matches(msg, m.keys.tabKey, m.keys.shiftTabKey, m.keys.enterKey, m.keys.upKey, m.keys.downKey):
				s := msg.String()

				var totalLength = len(m.channelModifyInputs) + 1

				// Did the user press enter while the submit button was focused?
				// If so, create it.
				if s == "enter" && m.channelModifyFocus == totalLength {
					channel := m.list.SelectedItem().(channel.Channel)
					err := channel.SetNotes(m.channelModifyInputs[0].Value())
					utils.HandleError(err, "updating channel notes")
					m.channelModifyFocus = 0
					// process selected tag list here
					sortedTags := slices.Sorted(maps.Keys(m.tags.ByName))
					var submitIds []int64
					for _, orderId := range m.selectedTagIds {
						tag := m.tags.ByName[sortedTags[orderId]]
						var found bool = false
						for _, currTagId := range m.selectedChannel.Tags() {
							if int64(currTagId) == tag.Id() {
								found = true
							}
						}
						// if it's found in the tags, of the selectedChannel, that means we're
						// removing it
						// otherwise we're adding it
						if !found {
							submitIds = append(submitIds, tag.Id())
						}
					}
					// include existing stuff so we don't lose it.
					for _, currTagId := range m.selectedChannel.Tags() {
						tag := m.tags.ById[currTagId]
						var found bool = false
						for _, orderId := range m.selectedTagIds {
							testTag := m.tags.ByName[sortedTags[orderId]]

							if testTag.Id() == tag.Id() {
								found = true
							}
						}
						// if it's found in the tags, of the selectedChannel, that means we're
						// removing it
						// otherwise we're adding it
						if !found {
							submitIds = append(submitIds, tag.Id())
						}
					}
					err = channel.SetTags(submitIds)
					utils.HandleError(err, "updating channel tags")
					m.selectedTagIds = nil
					m.channels.LoadEntriesFromDb()
					m.list.SetItems(m.generateChannelItems())
					m.current = "channel"
					return m, nil
				}

				// Cycle indexes
				if s == "up" || s == "shift+tab" {
					m.channelModifyFocus--
				} else {
					m.channelModifyFocus++
				}

				if m.channelModifyFocus > totalLength {
					m.channelModifyFocus = 0
				} else if m.channelModifyFocus < 0 {
					m.channelModifyFocus = totalLength
				}

				cmds := make([]tea.Cmd, totalLength)
				switch m.channelModifyFocus {
				// notes field
				case 0:
					cmds[0] = m.channelModifyInputs[0].Focus()
					m.channelModifyInputs[0].PromptStyle = focusedButtonStyle
					m.channelModifyInputs[0].TextStyle = focusedButtonStyle
				// tags
				case 1:
					m.channelModifyInputs[0].Blur()
					m.channelModifyInputs[0].PromptStyle = blurredButtonStyle
					m.channelModifyInputs[0].TextStyle = blurredButtonStyle
				// submit button
				case 2:
					m.channelModifyInputs[0].Blur()
					m.channelModifyInputs[0].PromptStyle = blurredButtonStyle
					m.channelModifyInputs[0].TextStyle = blurredButtonStyle
				}

				return m, tea.Batch(cmds...)
			case key.Matches(msg, m.keys.leftKey, m.keys.rightKey):
				s := msg.String()

				var totalLength = len(m.tags.ById)

				if m.channelModifyFocus == 1 {
					if s == "left" {
						m.selectedTagId--
					} else {
						m.selectedTagId++
					}

					if m.selectedTagId >= totalLength {
						m.selectedTagId = 0
					} else if m.selectedTagId < 0 {
						m.selectedTagId = totalLength - 1
					}
				} else {
					cmd = m.updateChannelModifyInput(msg)
				}
			case key.Matches(msg, keyList["spaceKey"]):
				if m.channelModifyFocus == 1 {
					var newTagIds []int
					var found bool = false
					for _, currTagId := range m.selectedTagIds {
						if currTagId == m.selectedTagId {
							found = true
						} else {
							newTagIds = append(newTagIds, currTagId)
						}
					}
					if !found {
						m.selectedTagIds = append(m.selectedTagIds, m.selectedTagId)
					} else {
						m.selectedTagIds = newTagIds
					}
				} else {
					cmd = m.updateChannelModifyInput(msg)
				}
				return m, tea.Batch(cmd)
			default:
				cmd = m.updateChannelModifyInput(msg)
			}
		}

	case "confirmDelete":
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			h, v := appStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)

		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.escKey, m.keys.qKey):
				m.current = "tag"
				return m, nil

			case key.Matches(msg, m.keys.tabKey, m.keys.shiftTabKey, m.keys.enterKey, m.keys.upKey, m.keys.downKey, m.keys.leftKey, m.keys.rightKey):
				s := msg.String()

				// Did the user press enter while the submit button was focused?
				// If so, create it.
				if s == "enter" {
					if m.tagDeleteFocus == tagDeleteDeleteButtonId {
						tag := m.list.SelectedItem().(tag.Tag)
						m.deleteTagEntryForm(tag)
						m.current = "confirmDelete"

						tag.Delete()
						m.tags.LoadEntriesFromDb()
						m.channels.LoadEntriesFromDb()
						m.list.SetItems(m.generateTagItems())
						m.tagDeleteFocus = tagDeleteCancelButtonId
						m.current = "tag"
						return m, nil
					} else {
						m.current = "tag"
						return m, nil
					}
				}

				// Cycle indexes
				if s == "up" || s == "shift+tab" || s == "left" {
					m.tagDeleteFocus--
				} else {
					m.tagDeleteFocus++
				}

				if m.tagDeleteFocus > tagDeleteDeleteButtonId {
					m.tagDeleteFocus = tagDeleteDeleteButtonId
				} else if m.tagDeleteFocus < tagDeleteCancelButtonId {
					m.tagDeleteFocus = tagDeleteCancelButtonId
				}

				return m, tea.Batch(cmds...)
			}
		}

	case "colourPicker":
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			h, v := appStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)

		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.upKey, m.keys.shiftTabKey):
				m.colourPickerY--

				if m.colourPickerY < 0 {
					m.colourPickerY = len(colours[0]) - 1
				}
				m.selectedBackColour = colours[m.colourPickerX][m.colourPickerY]
				return m, nil

			case key.Matches(msg, m.keys.downKey, m.keys.tabKey):
				m.colourPickerY++

				if m.colourPickerY > len(colours[0])-1 {
					m.colourPickerY = 0
				}
				m.selectedBackColour = colours[m.colourPickerX][m.colourPickerY]
				return m, nil

			case key.Matches(msg, m.keys.leftKey):
				m.colourPickerX--

				if m.colourPickerX < 0 {
					m.colourPickerX = len(colours) - 1
				}
				m.selectedBackColour = colours[m.colourPickerX][m.colourPickerY]
				return m, nil

			case key.Matches(msg, m.keys.rightKey):
				m.colourPickerX++

				if m.colourPickerX > len(colours)-1 {
					m.colourPickerX = 0
				}
				m.selectedBackColour = colours[m.colourPickerX][m.colourPickerY]
				return m, nil

			case key.Matches(msg, m.keys.enterKey):
				m.current = "tagEntry"
				if m.tagEntryFocus == 3 {
					m.tagEntryInputs[2].SetValue(strings.ReplaceAll(m.selectedBackColour, "#", ""))
				} else {
					m.tagEntryInputs[3].SetValue(strings.ReplaceAll(m.selectedBackColour, "#", ""))
				}
				m.colourPickerX = 0
				m.colourPickerY = 0

				return m, nil

			case key.Matches(msg, m.keys.escKey, m.keys.qKey):
				m.current = "tagEntry"
				return m, nil
			}
		}

	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var out string
	switch m.current {
	case "channel", "tag":
		channels = m.channels
		tags = m.tags
		out = m.list.View()

	case "tagEntry":
		var b strings.Builder

		var buttonRef *lipgloss.Style

		// do tagentry by hand
		if m.tagEntryInputs[0].Value() != m.selectedTag.Name() {
			m.tagEntryInputs[0].TextStyle = m.tagEntryInputs[0].TextStyle.Background(unsavedColour)
		}
		b.WriteString(fmt.Sprintf("%24s: %s\n", "tag name", m.tagEntryInputs[0].View()))
		if m.tagEntryInputs[1].Value() != m.selectedTag.Description() {
			m.tagEntryInputs[1].TextStyle = m.tagEntryInputs[1].TextStyle.Background(unsavedColour)
		}
		b.WriteString(fmt.Sprintf("%24s: %s\n", "tag description", m.tagEntryInputs[1].View()))
		if m.tagEntryFocus == 3 {
			buttonRef = &focusedButtonStyle
		} else {
			buttonRef = &blurredButtonStyle
		}
		if m.tagEntryInputs[2].Value() != m.selectedTag.FgColour() {
			m.tagEntryInputs[2].TextStyle = m.tagEntryInputs[2].TextStyle.Background(unsavedColour)
		}
		b.WriteString(fmt.Sprintf("%24s: %s %s\n", "foreground colour (hex)", m.tagEntryInputs[2].View(), buttonRef.Render("[ ColourPicker ]")))
		if m.tagEntryFocus == 5 {
			buttonRef = &focusedButtonStyle
		} else {
			buttonRef = &blurredButtonStyle
		}
		if m.tagEntryInputs[3].Value() != m.selectedTag.BgColour() {
			m.tagEntryInputs[3].TextStyle = m.tagEntryInputs[3].TextStyle.Background(unsavedColour)
		}
		b.WriteString(fmt.Sprintf("%24s: %s %s \n", "background colour (hex)", m.tagEntryInputs[3].View(), buttonRef.Render("[ ColourPicker ]")))

		if m.tagEntryFocus == 6 {
			buttonRef = &focusedButtonStyle
		} else {
			buttonRef = &blurredButtonStyle
		}
		var button = buttonRef.Render("[ Submit ]")
		fmt.Fprintf(&b, "\n%s\n\n\n", button)

		var style = tagDisplayStyle.Width(len(m.tagEntryInputs[0].Value())).Background(lipgloss.Color("#" + m.tagEntryInputs[3].Value())).Foreground(lipgloss.Color("#" + m.tagEntryInputs[2].Value()))

		b.WriteString("Example display: " + style.Render(m.tagEntryInputs[0].Value()))

		out = b.String()

	case "channelModify":
		var b strings.Builder

		for i := range m.channelModifyHeaders {
			b.WriteString(m.channelModifyHeaders[i])
		}

		if m.channelModifyInputs[0].Value() != m.selectedChannel.Notes() {
			m.channelModifyInputs[0].TextStyle = m.channelModifyInputs[0].TextStyle.Background(unsavedColour)
		}
		b.WriteString(m.channelModifyInputs[0].View())
		b.WriteString("\n")

		sortedTags := slices.Sorted(maps.Keys(m.tags.ByName))

		var cellSize int = 1
		for _, tagName := range sortedTags {
			if len(tagName) > cellSize {
				cellSize = len(tagName)
			}
		}
		cellSize += 2
		w, _, _ := term.GetSize(os.Stdout.Fd())
		colCount := int(math.Floor((float64)(w / cellSize)))

		var output string
		var curCol = 0
		for i, tagName := range sortedTags {
			tag := m.tags.ByName[tagName]
			var style = tagDisplayStyle.Width(len(tagName)).Background(lipgloss.Color("#" + tag.BgColour())).Foreground(lipgloss.Color("#" + tag.FgColour())).Margin(1)

			var found = false
			for _, currTagId := range m.selectedChannel.Tags() {
				if int64(currTagId) == tag.Id() {
					found = true
				}
			}

			var unsaved = false
			for _, currTagId := range m.selectedTagIds {
				if currTagId == i {
					unsaved = true
				}
			}

			if found {
				if unsaved {
					style = style.Border(lipgloss.HiddenBorder(), true, true, true, true)
				} else {
					style = style.Border(lipgloss.NormalBorder(), true, true, true, true)
				}
			} else {
				if unsaved {
					style = style.Border(lipgloss.NormalBorder(), true, true, true, true)
				} else {
					style = style.Border(lipgloss.HiddenBorder(), true, true, true, true)
				}
			}

			if unsaved {
				style = style.BorderBackground(unsavedColour)
			}

			if 1 == m.channelModifyFocus {
				if m.selectedTagId == i {
					style = style.BorderBackground(activeColour)
				}
			}

			output = lipgloss.JoinHorizontal(lipgloss.Center, output, style.Render(tagName))

			if (colCount > 2 && curCol > 2 && int(math.Mod(float64(colCount), float64(i+1))) == 0) || i == len(sortedTags)-1 || (colCount <= 2 && curCol+1 == colCount) || curCol+2 > colCount {
				b.WriteString(output)
				output = ""
				curCol = 0
			} else {
				curCol++
			}

		}

		var buttonRef *lipgloss.Style
		if m.channelModifyFocus == len(m.channelModifyInputs)+1 {
			buttonRef = &focusedButtonStyle
		} else {
			buttonRef = &blurredButtonStyle
		}
		var button = buttonRef.Render("[ Submit ]")
		fmt.Fprintf(&b, "\n\n%s\n\n", button)

		out = b.String()

	case "confirmDelete":
		var b strings.Builder

		b.WriteString("Confirm deletion:\n\n")

		for _, str := range m.tagDeleteInputs {
			b.WriteString(str)
		}
		b.WriteRune('\n')

		var buttonRef *lipgloss.Style
		if m.tagDeleteFocus == tagDeleteCancelButtonId {
			buttonRef = &focusedButtonStyle
		} else {
			buttonRef = &blurredButtonStyle
		}
		var button = fmt.Sprintf("[ %s ]", buttonRef.Render("Cancel"))
		fmt.Fprintf(&b, "%s   ", button)

		if m.tagDeleteFocus == tagDeleteDeleteButtonId {
			buttonRef = &focusedButtonStyle
		} else {
			buttonRef = &blurredButtonStyle
		}
		button = fmt.Sprintf("[ %s ]", buttonRef.Render("Delete"))
		fmt.Fprintf(&b, "%s\n", button)

		out = b.String()

	case "colourPicker":
		var b strings.Builder
		for row := 0; row < 12; row++ {
			for column := 0; column < len(colours); column++ {
				if column == m.colourPickerX &&
					row == m.colourPickerY {
					var textColour_64 float64
					textColour, err := hexer.Invert(colours[column][row])
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
					b.WriteString(focusedColourPickerModelStyle.Background(lipgloss.Color(colours[column][row])).Foreground(lipgloss.Color(textColour)).Render("XXXX"))
					// m.selectedTextColour = textColour
				} else {
					b.WriteString(colourPickerModelStyle.Background(lipgloss.Color(colours[column][row])).Render(""))
				}
			}
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf("X (column): %d\n", m.colourPickerX))
		b.WriteString(fmt.Sprintf("Y (row): %d\n", m.colourPickerY))
		// b.WriteString(fmt.Sprintf("textColour: %s\n", m.selectedTextColour))
		b.WriteString(fmt.Sprintf("colour: %s\n", m.selectedBackColour))

		out = b.String()
	}
	return appStyle.Render(out)
}

var m Model

func StartTea(channels channel.Channels, tags tag.Tags) {
	listKeys := newListKeyMap()

	m.channels = channels
	m.tags = tags

	m.list = list.New(m.generateChannelItems(), channelListItemDelegate{}, 0, 0)

	m.list.Title = "YSM - Channel View"
	m.list.Styles.Title = titleStyle
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.tKey,
			listKeys.pKey,
		}
	}
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.tKey,
			listKeys.pKey,
		}
	}

	m.current = "channel"
	m.keys = listKeys

	P = tea.NewProgram(m, tea.WithAltScreen())
	if _, err := P.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
