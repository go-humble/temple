// Copyright 2015 Alex Browne.
// All rights reserved. Use of this source code is
// governed by the MIT license, which can be found
// in the LICENSE file.

// Temple is a command-line tool for managing go templates which
// supports sharing templates between a client and server. It generates
// code which is compatible with gopherjs and can be compiled to
// javascript to run in the browser.
//
// Version 0.1.1
package main

import (
	"fmt"
	"github.com/albrow/prtty"
	"github.com/go-humble/temple/temple"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

const (
	version = "temple version 0.1.1"
)

var (
	verbose = false
)

// setQuiet effectively causes all loggers to print
// to /dev/null. However, error will still be printed
// out to stderr.
func setQuiet() {
	prtty.AllLoggers.SetOutput(ioutil.Discard)
	prtty.Error.Output = os.Stderr
}

// setVerbose sets all loggers to print to stdout,
// except for the Error logger, which will print to
// stderr.
func setVerbose() {
	prtty.AllLoggers.SetOutput(os.Stdout)
	prtty.Error.Output = os.Stderr
}

func main() {
	// Define build command
	cmdBuild := &cobra.Command{
		Use:   "build <src> <dest>",
		Short: "Compile the templates in the src directory and write generated go code to the dest file.",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				prtty.Error.Fatal("temple build requires exactly 2 arguments: the src directory and the dest file.")
			}
			if verbose {
				setVerbose()
			} else {
				setQuiet()
			}
			partials := cmd.Flag("partials").Value.String()
			layouts := cmd.Flag("layouts").Value.String()
			packageName := cmd.Flag("package").Value.String()
			if err := temple.Build(args[0], args[1], partials, layouts, packageName); err != nil {
				prtty.Error.Fatal(err)
			}
		},
	}
	cmdBuild.Flags().String("partials", "", "(optional) The directory to look for partials. Partials are .tmpl files that are associated with layouts and all other templates.")
	cmdBuild.Flags().String("layouts", "", "(optional) The directory to look for layouts. Layouts are .tmpl files which have access to partials and are associated with all other templates.")
	cmdBuild.Flags().String("package", "", "(optional) The package name for the generated go file. If not provided, the default will be the directory where the go file is created.")
	cmdBuild.Flags().BoolVarP(&verbose, "verbose", "v", false, "If set to true, temple will print out information while building.")

	// Define version command
	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "Print the current version number.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	// Define the root command
	rootCmd := &cobra.Command{
		Use:   "temple",
		Short: "A command line tool for sharing go templates between a client and server.",
		Long: `
A command line tool for sharing go templates between a client and server.
Visit https://github.com/albrow/temple for source code, example usage, documentation, and more.`,
	}
	rootCmd.AddCommand(cmdBuild, cmdVersion)
	if err := rootCmd.Execute(); err != nil {
		prtty.Error.Fatal(err)
	}
}
