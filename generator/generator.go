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
	"time"

	"gitea.joyrex.net/ejstacey/ysm/channel"
	"gitea.joyrex.net/ejstacey/ysm/tag"
	"gitea.joyrex.net/ejstacey/ysm/utils"
)

type Generator struct {
	Channels         []channel.ExportChannel
	Tags             []tag.ExportTag
	Title            string
	OutputFile       string
	TemplateFile     string
	GenerateDateTime string
}

var t *template.Template

func (g Generator) LoadTemplateFile() {
	input, err := os.ReadFile(g.TemplateFile)
	utils.HandleError(err, "Unable to open template.")

	t = template.Must(template.New("default").Parse(string(input)))
}

func (g Generator) GenerateOutputFile() {
	fo, err := os.Create(g.OutputFile)
	utils.HandleError(err, "Unable to open output file.")

	// close fo on exit and check for its returned error
	defer func() {
		err := fo.Close()
		utils.HandleError(err, "Unable to close output file.")
	}()

	g.GenerateDateTime = time.Now().Format(time.UnixDate)

	err = t.Execute(fo, g)
	if err != nil {
		panic(err)
	}
}
