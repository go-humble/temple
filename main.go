package main

import (
	"fmt"
	"github.com/albrow/prtty"
	"github.com/spf13/cobra"
)

const (
	version = "temple version X.X.X (develop)"
)

func main() {
	var cmdBuild = &cobra.Command{
		Use:   "build <src> <dest>",
		Short: "Compile the templates in the src directory into a package in the dest directory.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				prtty.Error.Fatal("temple build requires exactly 2 arguments: the src directory and the dest directory.")
			}
			if err := build(args[0], args[1], cmd.Flag("package").Value.String()); err != nil {
				prtty.Error.Fatal(err)
			}
		},
	}
	cmdBuild.Flags().String("package", "templates", "The name of the generated package.")

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Print the current version number.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	var rootCmd = &cobra.Command{
		Use:   "temple",
		Short: "A command line tool for sharing go templates between client and server.",
	}
	rootCmd.AddCommand(cmdBuild, cmdVersion)
	if err := rootCmd.Execute(); err != nil {
		prtty.Error.Fatal(err)
	}
}

func build(src, dest, packageName string) error {
	prtty.Info.Println("--> building")
	prtty.Default.Printf("    src: %s", src)
	prtty.Default.Printf("    dest: %s", dest)
	prtty.Default.Printf("    package: %s", packageName)
	return fmt.Errorf("build is not yet implemented!")
}
