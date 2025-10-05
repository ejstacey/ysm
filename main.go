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
<<<<<<< HEAD
		utils.InstallApp()
=======
		result, err := utils.FileDirExists(filepath.FromSlash("./templates/default.tmpl"))
		utils.HandleError(err, "Checking for existence of install template data.")
		if !result {
			fmt.Printf("The ./templates/default.tmpl file is not found! If you run --install, you should do it from the uncompressed directory of the downloaded package file.\n")
			os.Exit(1)
		}

		result, err = utils.FileDirExists(filepath.FromSlash("./html"))
		utils.HandleError(err, "Checking for existence of install html data.")
		if !result {
			fmt.Printf("The ./html directory is not found! If you run --install, you should do it from the uncompressed directory of the downloaded package file.\n")
			os.Exit(1)
		}

		result, err = utils.FileDirExists(filepath.FromSlash("./settings.sample.json.tmpl"))
		utils.HandleError(err, "Checking for existence of install settings sample.")
		if !result {
			fmt.Printf("The ./settings.sample.json.tmpl file is not found! If you run --install, you should do it from the uncompressed directory of the downloaded package file.\n")
			os.Exit(1)
		}

		userScope := gap.NewScope(gap.User, "ysm")

		// set up templates.
		templateDir, err := userScope.DataPath("templates")
		utils.HandleError(err, "Could not determine user template path for installation!")
		result, err = utils.FileDirExists(templateDir)
		utils.HandleError(err, "Checking for existence of user data path.")
		if !result {
			os.MkdirAll(templateDir, 0755)
		}
		err = cp.Copy(filepath.FromSlash("./templates"), templateDir)
		utils.HandleError(err, "Could not copy templates directory to "+templateDir)
		fmt.Println("Copied source templates to: " + templateDir)

		// set up html dir.
		htmlDir, err := userScope.DataPath("html")
		utils.HandleError(err, "Could not determine user html path for installation!")
		result, err = utils.FileDirExists(htmlDir)
		utils.HandleError(err, "Checking for existence of user data path.")
		if !result {
			os.MkdirAll(htmlDir, 0755)
		}
		err = cp.Copy(filepath.FromSlash("./html"), htmlDir)
		utils.HandleError(err, "Could not copy html directory to "+htmlDir)
		fmt.Println("Copied source html files to: " + htmlDir)

		// set up a default settings file
		configFile, err := userScope.ConfigPath("settings.json")
		utils.HandleError(err, "Could not determine user data path for installation!")
		settings := struct {
			DataDir     string
			OutputDir   string
			TemplateDir string
		}{
			utils.JsonEscape(filepath.FromSlash(filepath.Dir(templateDir) + "/")),
			utils.JsonEscape(filepath.FromSlash(htmlDir + "/")),
			utils.JsonEscape(filepath.FromSlash(templateDir + "/")),
		}

		var t *template.Template
		input, err := os.ReadFile("./settings.sample.json.tmpl")
		utils.HandleError(err, "Unable to open settings template.")

		t = template.Must(template.New("default").Parse(string(input)))

		dir := filepath.Dir(configFile)
		result, err = utils.FileDirExists(dir)
		utils.HandleError(err, "Checking for existence of config directory.")
		if !result {
			os.MkdirAll(dir, 0755)
		}

		fo, err := os.Create(configFile)
		utils.HandleError(err, "Unable to open settings file.")

		// close fo on exit and check for its returned error
		defer func() {
			err := fo.Close()
			utils.HandleError(err, "Unable to close settings file.")
		}()

		err = t.Execute(fo, settings)
		if err != nil {
			panic(err)
		}
		fmt.Println("Created a default settings file at: " + configFile)
		fmt.Println("")
		fmt.Println("You should now be able to copy the ysm exectuable wherever you want. Wherever it goes, it should be in your PATH.")
		fmt.Println("For Linux, ~/.local/bin can be a good option.")
		fmt.Println("For other OSs, you are on your own. Sorry!")
>>>>>>> 92162236677c23f90d2c3eea2afb9b48e9787533
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
