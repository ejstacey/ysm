/*
Joyrex YSM - Manager for Youtube Subscriptions
Copyright (C) 2025 Eric Stacey <ejstacey@joyrex.net>

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

*/

package channel

import (
	"context"
	"fmt"
	"slices"
	"time"

	"gogs.joyrex.net/ejstacey/ysm/utils"
)

type Channel struct {
	id          string
	name        string
	description string
	notes       string
	tags        []int64
}

func (c Channel) Id() string               { return c.id }
func (c Channel) FilterValue() string      { return c.name }
func (c Channel) Title() string            { return c.name }
func (c Channel) Name() string             { return c.name }
func (c Channel) Description() string      { return c.description }
func (c Channel) Notes() string            { return c.notes }
func (c Channel) Tags() []int64            { return c.tags }
func (c *Channel) SetDescription(x string) { c.description = x }

func LoadChannelsYoutube() []Channel {
	var service = utils.ConnectYoutube()

	// Set the parameters for the request
	var part []string
	part = append(part, "snippet") // Specify the resource properties you want to include
	var channels []Channel

	call := service.Subscriptions.List(part)
	call.Mine(true)

	// Make the API call
	for {
		response, err := call.Do()
		utils.HandleError(err, "Error retrieving descriptions")

		for _, subInfo := range response.Items {
			var channel Channel
			channel.id = subInfo.Id
			channel.name = subInfo.Snippet.Title
			channel.description = subInfo.Snippet.Description

			// fmt.Printf("%s %s %s\n", Id, Name, Description)
			fmt.Printf(".")
			channels = append(channels, channel)
		}

		if response.NextPageToken != "" {
			call.PageToken(response.NextPageToken)
		} else {
			break
		}

		if len(channels) > 300 {
			break
		}
	}
	fmt.Printf("\n")

	return channels
}

func (c *Channel) SetTags(x []int64) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	toAdd := utils.IntDifference(x, c.tags)
	for _, tagId := range toAdd {
		var insertSql = "insert into links (tagId, channelId) values (:tagId, :channelId)"

		insertSth, err := utils.DbConn.PrepareContext(ctx, insertSql)
		if err != nil {
			return err
		}

		_, err = insertSth.ExecContext(ctx, tagId, c.id)
		if err != nil {
			return err
		}
	}

	toDelete := utils.IntDifference(c.tags, x)
	for _, tagId := range toDelete {
		var deleteSql = "delete from links where tagId=:tagId and channelId=:channelId"

		deleteSth, err := utils.DbConn.PrepareContext(ctx, deleteSql)
		if err != nil {
			return err
		}

		_, err = deleteSth.ExecContext(ctx, tagId, c.id)
		if err != nil {
			return err
		}
	}

	c.tags = x

	return nil
}

func (c *Channel) SetNotes(x string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var updateSql = "update channels set notes = :notes where id = :id"

	updateSth, err := utils.DbConn.PrepareContext(ctx, updateSql)
	if err != nil {
		return err
	}

	_, err = updateSth.ExecContext(ctx, x, c.id)
	if err != nil {
		return err
	}

	c.notes = x

	return nil
}

type Channels struct {
	ById   map[string]Channel
	ByName map[string]Channel
}

func (c *Channels) LoadEntriesFromDb() {
	var channelsById = make(map[string]Channel)
	var channelsByName = make(map[string]Channel)

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var csText = "select id, name, description, ifnull(notes, '') from channels"

	csSth, err := utils.DbConn.PrepareContext(ctx, csText)
	utils.HandleError(err, "Unable to prepare sqlSelect "+csText)

	rows, err := csSth.QueryContext(ctx)
	utils.HandleError(err, "Unable to run sqlSelect "+csText)

	defer rows.Close()
	for rows.Next() {
		var channel Channel
		err = rows.Scan(&channel.id, &channel.name, &channel.description, &channel.notes)
		utils.HandleError(err, "Loading existing channels from db (row)")

		var linkText = "select * from links where channelId = :id"

		linkSth, err := utils.DbConn.PrepareContext(ctx, linkText)
		utils.HandleError(err, "Unable to prepare sqlSelect "+linkText)

		linkRows, err := linkSth.QueryContext(ctx, channel.id)
		utils.HandleError(err, "Unable to run sqlSelect "+linkText)

		defer linkRows.Close()
		for linkRows.Next() {
			var channelId string
			var tagId int64

			err = linkRows.Scan(&channelId, &tagId)
			utils.HandleError(err, "Loading existing links from db (row)")

			channel.tags = append(channel.tags, tagId)
		}
		slices.Sort(channel.tags)
		channelsById[channel.id] = channel
		channelsByName[channel.name] = channel
	}
	utils.HandleError(err, "Loading existing channels from db")
	c.ById = channelsById
	c.ByName = channelsByName
}

func (c *Channels) CompareAndUpdateChannelsDb(newChannels []Channel) {
	for _, newEntry := range newChannels {
		var found = false
		var oldEntry Channel

		for _, dbEntry := range c.ById {
			if newEntry.id == dbEntry.id {
				found = true
				oldEntry = dbEntry
				break
			}
		}

		if !found {
			fmt.Printf("not found, adding to db: %s\n", newEntry.name)
			var sqlText = `
				insert into CHANNELS (
					id,
					name,
					description
				) values (
					?,
					?,
					?
				)
			`
			_, err := utils.DbConn.Exec(sqlText, newEntry.id, newEntry.name, newEntry.description)
			utils.HandleError(err, fmt.Sprintf("Unable to insert new entry to db (%s, %s, %s)", newEntry.id, newEntry.name, newEntry.description))
		} else {
			if (oldEntry.name != newEntry.name) || (oldEntry.description != newEntry.description) {
				fmt.Printf("found, updating db: %s\n", newEntry.name)

				var sqlText = `
					update CHANNELS
					set
						name = :name,
						description = :description
					where
						id = :id
				`
				_, err := utils.DbConn.Exec(sqlText, newEntry.name, newEntry.description, newEntry.id)
				utils.HandleError(err, "Unable to update entry on db")
			}
		}
	}

	for _, dbEntry := range c.ById {
		var found = false

		for _, newEntry := range newChannels {
			if newEntry.id == dbEntry.id {
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("not found, deleting from db: %s\n", dbEntry.name)

			var sqlText = `
				delete from CHANNELS
				where id = :id
			`
			_, err := utils.DbConn.Exec(sqlText, dbEntry.id)
			utils.HandleError(err, "Unable to delete entry on db")
		}
	}
}
