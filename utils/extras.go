/*
Joyrex YSM - Manager for Youtube Subscriptions
Copyright (C) 2025 Eric Stacey <ejstacey@joyrex.net>

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

*/

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	gap "github.com/muesli/go-app-paths"
	cp "github.com/otiai10/copy"
)

func FileDirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

<<<<<<< HEAD
func InstallApp() {
	userScope := gap.NewScope(gap.User, "ysm")

	result, err := FileDirExists(filepath.Clean("./templates/default.tmpl"))
	HandleError(err, "Checking for existence of install template data.")
	if !result {
		fmt.Printf("The ./templates/default.tmpl file is not found! If you run --install, you should do it from the uncompressed directory of the downloaded package file.\n")
		os.Exit(1)
	}

	result, err = FileDirExists(filepath.Clean("./html"))
	HandleError(err, "Checking for existence of install html data.")
	if !result {
		fmt.Printf("The ./html directory is not found! If you run --install, you should do it from the uncompressed directory of the downloaded package file.\n")
		os.Exit(1)
	}

	result, err = FileDirExists(filepath.Clean("./settings.sample.json.tmpl"))
	HandleError(err, "Checking for existence of install settings sample.")
	if !result {
		fmt.Printf("The ./settings.sample.json.tmpl file is not found! If you run --install, you should do it from the uncompressed directory of the downloaded package file.\n")
		os.Exit(1)
	}

	// set up base directories
	baseDir, err := userScope.DataPath("")
	HandleError(err, "Could not determine data path for installation!")
	result, err = FileDirExists(baseDir)
	HandleError(err, "Checking for existence of user data path.")
	if !result {
		os.MkdirAll(baseDir, 0755)
	}

	baseDir, err = userScope.ConfigPath("")
	HandleError(err, "Could not determine config path for installation!")
	result, err = FileDirExists(baseDir)
	HandleError(err, "Checking for existence of user config path.")
	if !result {
		os.MkdirAll(baseDir, 0755)
	}

	baseDir, err = os.UserHomeDir()
	HandleError(err, "Could not determine homedir path for installation!")
	binDir := filepath.Clean(baseDir + "/.local/bin")
	result, err = FileDirExists(binDir)
	HandleError(err, "Checking for existence of homedir bin path.")
	if !result {
		os.MkdirAll(binDir, 0755)
	}

	// set up templates.
	templateDir, err := userScope.DataPath("templates")
	HandleError(err, "Could not determine user template path for installation!")
	result, err = FileDirExists(templateDir)
	HandleError(err, "Checking for existence of user data path.")
	if !result {
		os.MkdirAll(templateDir, 0755)
	}
	err = cp.Copy(filepath.Clean("./templates"), templateDir)
	HandleError(err, "Could not copy templates directory to "+templateDir)
	fmt.Println("Copied source templates to: " + templateDir)

	// set up html dir.
	htmlDir, err := userScope.DataPath("html")
	HandleError(err, "Could not determine user html path for installation!")
	result, err = FileDirExists(htmlDir)
	HandleError(err, "Checking for existence of user data path.")
	if !result {
		os.MkdirAll(htmlDir, 0755)
	}
	err = cp.Copy(filepath.Clean("./html"), htmlDir)
	HandleError(err, "Could not copy html directory to "+htmlDir)
	fmt.Println("Copied source html files to: " + htmlDir)

	// set up a default settings file
	configFile, err := userScope.ConfigPath("settings.json")
	HandleError(err, "Could not determine user config path for installation!")
	settings := struct {
		DataDir     string
		OutputDir   string
		TemplateDir string
	}{
		filepath.Dir(templateDir),
		htmlDir,
		templateDir,
	}

	var t *template.Template
	input, err := os.ReadFile(filepath.Clean("./settings.sample.json.tmpl"))
	HandleError(err, "Unable to open settings template.")

	t = template.Must(template.New("default").Parse(string(input)))

	result, err = FileDirExists(configFile)
	HandleError(err, "Checking for existence of config file.")
	if result {
		timeRec := time.Now()
		timeStamp := timeRec.Format(time.DateOnly) + "-" + strings.ReplaceAll(timeRec.Format(time.TimeOnly), ":", "-")
		backupName := configFile + "-backup-" + timeStamp
		err = cp.Copy(configFile, backupName)
		HandleError(err, "Could not create backup of settings file "+configFile+" to "+backupName)
		fmt.Println("settings.json already exists. Created backup copy at " + backupName)
	}

	fo, err := os.Create(configFile)
	HandleError(err, "Unable to open settings file.")

	// close fo on exit and check for its returned error
	defer func() {
		err := fo.Close()
		HandleError(err, "Unable to close settings file.")
	}()

	err = t.Execute(fo, settings)
	HandleError(err, "Unable to write settings file.")
	fmt.Println("Created a default settings file at: " + configFile)

	err = cp.Copy(filepath.Clean("./ysm"), filepath.Clean(binDir+"/ysm"))
	HandleError(err, "Could not copy executable to "+binDir)
	fmt.Println("Copied ysm executable to: " + binDir)

	fmt.Println("")
	fmt.Println("Checking for " + binDir + " in your path.")
	path := os.Getenv("PATH")
	var paths []string
	if runtime.GOOS == "windows" {
		paths = strings.Split(path, ";")
	} else {
		paths = strings.Split(path, ":")
	}
	var found bool = false
	for _, path := range paths {
		abs, err := filepath.Abs(path)
		HandleError(err, "Could not get absolute location of "+abs)
		if abs == binDir {
			found = true
		}
	}
	if found {
		fmt.Println("Found " + binDir + " in path. Yay!")
	} else {
		fmt.Println(binDir + " not found in your PATH. Do some searching to find out how to add " + binDir + " to your PATH.")
		fmt.Println("Alternatively, just run ysm with: " + binDir + string(filepath.Separator) + "ysm")
	}
}

func VerifyInstallation() {
	userScope := gap.NewScope(gap.User, "ysm")

	configFile, err := userScope.ConfigPath("settings.json")
	HandleError(err, "Could not determine user config path.")
	configFileCheck, err := FileDirExists(configFile)
	HandleError(err, "Checking for existence of user config.")
	templatePath, err := userScope.DataPath("templates")
	HandleError(err, "Could not determine templates path.")
	templatePathCheck, err := FileDirExists(templatePath)
	HandleError(err, "Checking for existence of templates data path.")
	outputPath, err := userScope.DataPath("html")
	HandleError(err, "Could not determine output data path.")
	outputPathCheck, err := FileDirExists(outputPath)
	HandleError(err, "Checking for existence of output data path.")
	if !configFileCheck || !templatePathCheck || !outputPathCheck {
		fmt.Println("Could not find required directories. Please run the program with the --install argument.")
		os.Exit(1)
	}
}

=======
>>>>>>> 92162236677c23f90d2c3eea2afb9b48e9787533
func JsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}
