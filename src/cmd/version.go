package cmd

import (
	"fmt"

	"github.com/bclicn/color"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Current version of app",
	Run:   runVersion,
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Println(color.BBlack("Version:\t") + color.Blue(version))
	fmt.Println(color.BBlack("Branch:\t\t") + color.Blue(branch))
	fmt.Println(color.BBlack("Commit:\t\t") + color.Blue(commit))
	fmt.Println(color.BBlack("BuildTime:\t") + color.Blue(buildTime))
}
