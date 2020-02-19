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
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/h9k-io/utils/vers/ventry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Make a new version file",
	Long:  "Make a new version file",
	Run:   newVers,
}

func init() {

	initCmd.PersistentFlags().Bool(FORCE, false, "force write the file (iff it exists) ")
	viper.BindPFlag(FORCE, initCmd.PersistentFlags().Lookup(FORCE))

	RootCmd.AddCommand(initCmd)
}

func newVers(cmd *cobra.Command, args []string) {
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
		entry = strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	}
	if _, err := os.Stat(filename); err == nil {
		if !viper.GetBool(FORCE) {
			log.Fatalf("File exists already and --force not set")
		}
	}
	vp, err := ventry.Open(filename, true)
	if err != nil {
		log.Fatalf("Failed to open %s; %s", filename, err)
	}
	defer vp.Close()
	vp.Add(entry, &ventry.Vers{
		Prefix: viper.GetString(PREFIX),
		Major:  viper.GetInt(MAJ),
		Minor:  viper.GetInt(MIN),
		Patch:  viper.GetInt(PATCH),
	})
	if err = vp.Write(3); err != nil {
		log.Infof("Failed to write %s; %s", filename, err)
	}
}
