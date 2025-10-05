/*
Joyrex YSM - Manager for Youtube Subscriptions
Copyright (C) 2025 Eric Stacey <ejstacey@joyrex.net>

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

*/

package main

import (
	"flag"
	"fmt"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"repo.joyrex.net/ejstacey/ysm/channel"
	"repo.joyrex.net/ejstacey/ysm/tag"
	"repo.joyrex.net/ejstacey/ysm/tui"
	"repo.joyrex.net/ejstacey/ysm/utils"
)

func main() {
	var installFlag = flag.Bool("install", false, "install the package files into their expected locations of your user environment")
	flag.Parse()
	if *installFlag {
		utils.InstallApp()
	} else {
		var channels channel.Channels
		var tags tag.Tags

		// verify basic directories exist
		utils.VerifyInstallation()

		var settings = utils.LoadSettings()

		utils.BackupDbFile(settings.DbFile, settings.BackupCopies)

		utils.InitDb(settings.DbFile)

		fmt.Printf("Loading existing DB channel entries.\n")
		channels.LoadEntriesFromDb()

		if settings.Refresh || len(channels.ById()) == 0 {
			fmt.Printf("No existing DB entries found or Refresh was set to True in settings.json. Loading a fresh list from YouTube.\n\n")
			var newChannels = channel.LoadChannelsYoutube()
			channels.CompareAndUpdateChannelsDb(newChannels)
			channels.LoadEntriesFromDb()
		}

		fmt.Printf("Loading existing DB tag entries.\n")
		tags.LoadEntriesFromDb()

		tui.StartTea(channels, tags, settings)
	}
}

// err := os.WriteFile("debug.log", []byte(dump.Format(tag)), 0644)
// if err != nil {
// 	panic(err)
// }
