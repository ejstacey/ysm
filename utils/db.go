/*
Joyrex YSM - Manager for Youtube Subscriptions
Copyright (C) 2025 Eric Stacey <ejstacey@joyrex.net>

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

*/

package utils

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	cp "github.com/otiai10/copy"
)

var DbConn *sql.DB

func InitDb(dbFile string) {
	var sqlText string
	var db *sql.DB // Database connection pool.

	var err error
	db, err = sql.Open("sqlite3", dbFile)
	HandleError(err, "Unable to open db file")

	sqlText = `
		CREATE TABLE IF NOT EXISTS channels (
			id      		TEXT PRIMARY KEY,
			name     		TEXT NOT NULL UNIQUE,
			description     TEXT,
			notes			TEXT
		);
	`
	_, err = db.Exec(sqlText)
	HandleError(err, "Unable to create channels table")

	sqlText = `
		CREATE TABLE IF NOT EXISTS tags (
			id         		INTEGER PRIMARY KEY AUTOINCREMENT,
			name      		TEXT,
			description     TEXT,
			bgColour			TEXT,
			fgColour			TEXT

		);
	`
	_, err = db.Exec(sqlText)
	HandleError(err, "Unable to create tags table")

	sqlText = `
		CREATE TABLE IF NOT EXISTS links (
			channelId      	TEXT,
			tagId     		INTEGER,
			PRIMARY KEY (channelId, tagId),
			FOREIGN KEY (channelId) REFERENCES channels(id) ON DELETE CASCADE,
			FOREIGN KEY (tagId) REFERENCES tags(id) ON DELETE CASCADE 
		);
	`
	_, err = db.Exec(sqlText)
	HandleError(err, "Unable to create links table")

	DbConn = db
}

func BackupDbFile(dbFile string, backupCopies int) {
	result, err := FileDirExists(dbFile)
	HandleError(err, "Checking for existence of existing dbFile.")
	if result {
		timeRec := time.Now()
		timeStamp := timeRec.Format(time.DateOnly) + "-" + strings.ReplaceAll(timeRec.Format(time.TimeOnly), ":", "-")

		// Make sure backups dir exists
		backupPath := filepath.Clean(filepath.Dir(dbFile)+"/backups/") + string(filepath.Separator)
		backupPathCheck, err := FileDirExists(backupPath)
		HandleError(err, "Checking for existence of backup path.")
		if !backupPathCheck {
			os.MkdirAll(backupPath, 0755)
		}

		// Save the file
		backupName := backupPath + filepath.Base(dbFile) + "-backup-" + timeStamp
		err = cp.Copy(dbFile, backupName)
		HandleError(err, "Could not create backup of database file "+dbFile+" to "+backupName)
		fmt.Println("Created backup copy of database to " + backupName)

		// Clean up anything over backupCopies
		entries, err := os.ReadDir(backupPath)
		HandleError(err, "Could not open backups directory "+backupPath)

		if len(entries) > backupCopies {
			toDelete := backupPath + entries[0].Name()
			err := os.Remove(toDelete)
			HandleError(err, "Could not remove file "+toDelete)

			fmt.Println("Removed old backup file: " + toDelete)
		}
	}
}
