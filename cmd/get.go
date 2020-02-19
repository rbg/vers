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
	"github.com/apex/log"
	"github.com/h9k-io/utils/vers/ventry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// getCmd represents the get command
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "get version info",
		Long:  `For the given binary get the current version information`,
		Run:   get,
	}
)

func init() {
	getCmd.PersistentFlags().StringP(FMT, "o", "json", "Output format")
	viper.BindPFlag(FMT, getCmd.PersistentFlags().Lookup(FMT))

	RootCmd.AddCommand(getCmd)
}

func get(cmd *cobra.Command, args []string) {

	if viper.GetBool(DEBUG) {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugf("filename: %s, fmt: %s, entry: %s",
		viper.GetString(VFILE), viper.GetString(FMT),
		viper.GetString(ENTRY))

	filename := viper.GetString(VFILE)
	if len(filename) == 0 {
		log.Fatalf("you must supply the .json or .yaml version file pathname (--%s)", VFILE)
	}

	vp, err := ventry.Open(filename, false)
	if err != nil {
		log.Fatalf("Open failed on %s; %s", filename, err)
	}
	defer vp.Close()
	err = vp.Read(10)
	if err != nil {
		log.Fatalf("Read  failed on %s; %s", filename, err)
	}
	entry := viper.GetString(ENTRY)
	if len(entry) != 0 {
		if err := vp.Print(entry, viper.GetString(FMT)); err != nil {
			log.Fatalf("Failed: %s", err)
		}
		return
	}
	vp.Dump(viper.GetString(FMT))
}
