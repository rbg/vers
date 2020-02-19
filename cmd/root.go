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

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	BUMP   = "bump"
	CFG    = "config"
	DEBUG  = "debug"
	ENTRY  = "entry"
	FMT    = "fmt"
	FORCE  = "force"
	MAJ    = "major"
	MIN    = "minor"
	PATCH  = "patch"
	PREFIX = "prefix"
	SUFFIX = "suffix"
	VFILE  = "version-file"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "vers",
	Short: "A simple way to manage versions",
	Long:  `Handle versions....`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().String(CFG, "", "config file (default is $HOME/.vers.yaml)")
	viper.BindPFlag(CFG, RootCmd.PersistentFlags().Lookup(CFG))

	RootCmd.PersistentFlags().BoolP(DEBUG, "d", false, "Turn on debug messages")
	viper.BindPFlag(DEBUG, RootCmd.PersistentFlags().Lookup(DEBUG))

	RootCmd.PersistentFlags().StringP(VFILE, "f", "", "version file to use")
	viper.BindPFlag(VFILE, RootCmd.PersistentFlags().Lookup(VFILE))

	RootCmd.PersistentFlags().IntP(MAJ, "M", 0, "major number (default: 0)")
	viper.BindPFlag(MAJ, RootCmd.PersistentFlags().Lookup(MAJ))

	RootCmd.PersistentFlags().IntP(MIN, "m", 0, "minor number (default: 0)")
	viper.BindPFlag(MIN, RootCmd.PersistentFlags().Lookup(MIN))

	RootCmd.PersistentFlags().IntP(PATCH, "p", 1, "patch number ")
	viper.BindPFlag(PATCH, RootCmd.PersistentFlags().Lookup(PATCH))

	RootCmd.PersistentFlags().String(PREFIX, "v", "prefix ")
	viper.BindPFlag(PREFIX, RootCmd.PersistentFlags().Lookup(PREFIX))

	RootCmd.PersistentFlags().String(SUFFIX, "", "suffix")
	viper.BindPFlag(SUFFIX, RootCmd.PersistentFlags().Lookup(SUFFIX))

	RootCmd.PersistentFlags().StringP(ENTRY, "e", "", "Which entry in version file")
	viper.BindPFlag(ENTRY, RootCmd.PersistentFlags().Lookup(ENTRY))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/h9k")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
		viper.SetConfigName(".vers")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Config file was found but an error occured; %s", err)
		}
	}
}
