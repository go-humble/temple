package main

import (
	"fmt"
	"github.com/albrow/prtty"
	"github.com/albrow/temple"
	"github.com/spf13/cobra"
)

const (
	version = "temple version X.X.X (develop)"
)

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
			partials := cmd.Flag("partials").Value.String()
			layouts := cmd.Flag("layouts").Value.String()
			if err := temple.Build(args[0], args[1], partials, layouts); err != nil {
				prtty.Error.Fatal(err)
			}
		},
	}
	cmdBuild.Flags().String("partials", "", "(optional) The directory to look for partials. Partials are .tmpl files that are associated with layouts and all other templates.")
	cmdBuild.Flags().String("layouts", "", "(optional) The directory to look for layouts. Layouts are .tmpl files which have access to partials and are associated with all other templates.")

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
