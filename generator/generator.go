/*
Joyrex YSM - Manager for Youtube Subscriptions
Copyright (C) 2025 Eric Stacey <ejstacey@joyrex.net>

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

*/

package generator

import (
	"html/template"
	"os"

	"gitea.joyrex.net/ejstacey/ysm/channel"
	"gitea.joyrex.net/ejstacey/ysm/tag"
	"gitea.joyrex.net/ejstacey/ysm/utils"
)

type Generator struct {
	Channels []channel.Channel
	Tags     []tag.Tag
	Title    string
}

var t *template.Template

func (g Generator) LoadTemplateFile() {
	input, err := os.ReadFile("templates/default.tmpl")
	utils.HandleError(err, "Unable to open template.")
	t = template.Must(template.New("default").Parse(string(input)))
}

func (g Generator) OutputFile() {
	fo, err := os.Create("output/index.html")
	utils.HandleError(err, "Unable to open output file.")

	// close fo on exit and check for its returned error
	defer func() {
		err := fo.Close()
		utils.HandleError(err, "Unable to close output file.")
	}()

	err = t.Execute(fo, g)
	if err != nil {
		panic(err)
	}
}
