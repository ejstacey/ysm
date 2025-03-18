/*
Joyrex YSM - Manager for Youtube Subscriptions
Copyright (C) 2025 Eric Stacey <ejstacey@joyrex.net>

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

*/

package tag

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"gogs.joyrex.net/ejstacey/ysm/utils"
)

type Tag struct {
	id          int64
	name        string
	description string
	bgColour    string
	fgColour    string
	channels    []string
}

func (t Tag) Id() int64           { return t.id }
func (t Tag) FilterValue() string { return t.name }
func (t Tag) Title() string       { return t.name }
func (t Tag) Name() string        { return t.name }
func (t Tag) Description() string { return t.description }
func (t Tag) BgColour() string    { return t.bgColour }
func (t Tag) FgColour() string    { return t.fgColour }
func (t Tag) Channels() []string  { return t.channels }
func (t *Tag) SetTitle(x string)  { t.SetName(x) }

// creates an entry with a generic name
func (t *Tag) New() error {
	var name string
	name = utils.RandSeq(10)
	for {
		err := t.validateName(0, name)
		if err != nil {
			name = utils.RandSeq(10)
		} else {
			break
		}
	}

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var insertSql = "insert into tags (name) values (:name)"

	insertSth, err := utils.DbConn.PrepareContext(ctx, insertSql)
	if err != nil {
		return err
	}

	var res sql.Result
	res, err = insertSth.ExecContext(ctx, name)
	if err != nil {
		return err
	}

	t.id, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (t *Tag) Delete() error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var deleteSql = "delete from tags where id = :id"

	deleteSth, err := utils.DbConn.PrepareContext(ctx, deleteSql)
	if err != nil {
		return err
	}

	// var res sql.Result
	_, err = deleteSth.ExecContext(ctx, t.id)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tag) SetName(x string) error {
	err := t.validateName(0, x)
	if err != nil {
		return err
	}

	if t.id <= 0 {
		return errors.New("cannot set name, missing id")
	}

	// err = os.WriteFile("debug.log", []byte(dump.Format(t)), 0644)
	// if err != nil {
	// 	panic(err)
	// }

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var updateSql = "update tags set name = :name where id = :id"

	updateSth, err := utils.DbConn.PrepareContext(ctx, updateSql)
	if err != nil {
		return err
	}

	_, err = updateSth.ExecContext(ctx, x, t.id)
	if err != nil {
		return err
	}

	t.name = x

	return nil
}

func (t *Tag) SetDescription(x string) error {
	if t.id <= 0 {
		return errors.New("cannot set description, missing id")
	}

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var updateSql = "update tags set description = :description where id = :id"

	updateSth, err := utils.DbConn.PrepareContext(ctx, updateSql)
	if err != nil {
		return err
	}

	_, err = updateSth.ExecContext(ctx, x, t.id)
	if err != nil {
		return err
	}

	t.description = x

	return nil
}

func (t *Tag) SetBgColour(x string) error {
	if t.id <= 0 {
		return errors.New("cannot set bg colour, missing id")
	}

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var updateSql = "update tags set bgColour = :colour where id = :id"

	updateSth, err := utils.DbConn.PrepareContext(ctx, updateSql)
	if err != nil {
		return err
	}

	_, err = updateSth.ExecContext(ctx, x, t.id)
	if err != nil {
		return err
	}

	t.bgColour = x

	return nil
}

func (t *Tag) SetFgColour(x string) error {
	if t.id <= 0 {
		return errors.New("cannot set fg colour, missing id")
	}

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var updateSql = "update tags set fgColour = :colour where id = :id"

	updateSth, err := utils.DbConn.PrepareContext(ctx, updateSql)
	if err != nil {
		return err
	}

	_, err = updateSth.ExecContext(ctx, x, t.id)
	if err != nil {
		return err
	}

	t.fgColour = x

	return nil
}

func (t *Tag) SetChannels(x []string) error {
	if t.id <= 0 {
		return errors.New("cannot set channels, missing tag id")
	}

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	toAdd := utils.StringDifference(x, t.channels)
	for _, channelId := range toAdd {
		var insertSql = "insert into links (tagId, channelId) values (:tagId, :channelId)"

		insertSth, err := utils.DbConn.PrepareContext(ctx, insertSql)
		if err != nil {
			return err
		}

		_, err = insertSth.ExecContext(ctx, t.id, channelId)
		if err != nil {
			return err
		}
	}

	toDelete := utils.StringDifference(t.channels, x)
	for _, channelId := range toDelete {
		var deleteSql = "delete from links where tagId=:tagId and channelId=:channelId"

		deleteSth, err := utils.DbConn.PrepareContext(ctx, deleteSql)
		if err != nil {
			return err
		}

		_, err = deleteSth.ExecContext(ctx, t.id, channelId)
		if err != nil {
			return err
		}
	}

	t.channels = x

	return nil
}

func (t Tag) validateName(num int, name string) error {
	if name == "" {
		return errors.New("name for new tag must be set")
	}

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var checkSql = "select count(*) from tags where lower(name) == :name and id != :id"

	checkSth, err := utils.DbConn.PrepareContext(ctx, checkSql)
	if err != nil {
		return err
	}

	var checkRows int
	err = checkSth.QueryRowContext(ctx, name, t.id).Scan(&checkRows)
	if err != nil {
		return err
	}

	if checkRows != num {
		return fmt.Errorf("didn't get expected number of rows (got: %d, expected: %d)", checkRows, num)
	}

	return nil
}

type Tags struct {
	ById   map[int64]Tag
	ByName map[string]Tag
}

func (t *Tags) LoadEntriesFromDb() {
	var tagsById = make(map[int64]Tag)
	var tagsByName = make(map[string]Tag)

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var csText = "select id, name, description, bgColour, fgColour from tags"

	csSth, err := utils.DbConn.PrepareContext(ctx, csText)
	utils.HandleError(err, "Unable to prepare sqlSelect "+csText)

	rows, err := csSth.QueryContext(ctx)
	utils.HandleError(err, "Unable to run sqlSelect "+csText)

	defer rows.Close()
	for rows.Next() {
		var tag Tag
		var tmpDescription sql.NullString
		var tmpBgColour sql.NullString
		var tmpFgColour sql.NullString

		err = rows.Scan(&tag.id, &tag.name, &tmpDescription, &tmpBgColour, &tmpFgColour)
		utils.HandleError(err, "Loading existing tag from db (row)")

		if tmpDescription.Valid {
			tag.description = tmpDescription.String
		} else {
			tag.description = ""
		}

		if tmpBgColour.Valid {
			tag.bgColour = tmpBgColour.String
		} else {
			tag.bgColour = "#FF0000"
		}

		if tmpFgColour.Valid {
			tag.fgColour = tmpFgColour.String
		} else {
			tag.fgColour = "#FFFFFF"
		}

		var linkText = "select * from links where tagId = :id"

		linkSth, err := utils.DbConn.PrepareContext(ctx, linkText)
		utils.HandleError(err, "Unable to prepare sqlSelect "+linkText)

		linkRows, err := linkSth.QueryContext(ctx, tag.id)
		utils.HandleError(err, "Unable to run sqlSelect "+linkText)

		defer linkRows.Close()
		for linkRows.Next() {
			var channelId string
			var tagId int

			err = linkRows.Scan(&channelId, &tagId)
			utils.HandleError(err, "Loading existing links from db (row)")

			tag.channels = append(tag.channels, channelId)
		}
		slices.Sort(tag.channels)
		tagsById[tag.id] = tag
		tagsByName[tag.name] = tag
	}
	utils.HandleError(err, "Loading existing tag from db")

	t.ById = tagsById
	t.ByName = tagsByName
}
