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
	"errors"
	"fmt"

	"github.com/apex/log"
	"github.com/h9k-io/utils/vers/ventry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ( // bumpCmd represents the bump command
	bumpCmd = &cobra.Command{
		Use:   "bump",
		Short: "increment either major, minor or patch version number",
		Long:  "increment either major, minor or patch version number",
		Args: func(cmd *cobra.Command, args []string) error {
			if viper.GetBool(DEBUG) {
				log.SetLevel(log.DebugLevel)
			}
			if len(viper.GetString(VFILE)) == 0 {
				return fmt.Errorf("you must supply the .json or .yaml version file pathname (--%s)", VFILE)
			}
			if len(viper.GetString(ENTRY)) == 0 {
				return fmt.Errorf("you must supply the entry name (--%s)", ENTRY)
			}
			switch what := viper.GetString(BUMP); what {
			case "major":
				fallthrough
			case "minor":
				fallthrough
			case "patch":
				return nil
			}
			return errors.New("valid values for bump is one of the following: `major,minor,patch`")
		},
		Run: bump,
	}
)

func init() {

	bumpCmd.Flags().StringP(BUMP, "i", "", "Increamt value (one of 'major,minor or patch')")
	viper.BindPFlag(BUMP, bumpCmd.Flags().Lookup(BUMP))

	RootCmd.AddCommand(bumpCmd)
}

func bump(cmd *cobra.Command, args []string) {

	vp, err := ventry.Open(viper.GetString(VFILE), false)
	if err != nil {
		log.Fatalf("Open failed on %s; %s", viper.GetString(VFILE), err)
	}
	defer vp.Close()
	if err = vp.Bump(viper.GetString(ENTRY), viper.GetString(BUMP)); err != nil {
		log.Infof("Bump failed on %s; %s", viper.GetString(VFILE), err)
		return
	}
	vp.Print(viper.GetString(ENTRY), "str")
}
