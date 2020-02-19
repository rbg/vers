package cmd

// Copyright © 2020 Robert B Gordon <rbg@h9k.io>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"

	"github.com/apex/log"
	"github.com/h9k-io/utils/vers/ventry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// setCmd represents the set command
	setCmd = &cobra.Command{
		Use:   "set",
		Short: "Add a new entry to version file",
		Long:  `This will add an entry (or update an existing) to the version file`,
		Run:   set,
	}
)

func init() {
	RootCmd.AddCommand(setCmd)
}

func set(cmd *cobra.Command, args []string) {
	if viper.GetBool(DEBUG) {
		fmt.Println("Set Debug")
		log.SetLevel(log.DebugLevel)
	}
	filename := viper.GetString(VFILE)
	if len(filename) == 0 {
		log.Fatalf("you must supply the .json or .yaml version file pathname (--%s)", VFILE)
	}
	entry := viper.GetString(ENTRY)
	if len(entry) == 0 {
		log.Fatalf("you must supply entry name (--%s)", ENTRY)
	}
	vp, err := ventry.Open(filename, false)
	if err != nil {
		log.Fatalf("Failed to open %s; %s", filename, err)
	}
	defer vp.Close()
	if err := vp.Read(10); err != nil {
		log.Infof("Failed to read %s; %s", filename, err)
		return
	}
	vp.Add(entry, &ventry.Vers{
		Prefix: viper.GetString(PREFIX),
		Major:  viper.GetInt(MAJ),
		Minor:  viper.GetInt(MIN),
		Patch:  viper.GetInt(PATCH),
		Suffix: viper.GetString(SUFFIX),
	})
	if err = vp.Write(3); err != nil {
		log.Infof("Failed to write %s; %s", filename, err)
	}

}
