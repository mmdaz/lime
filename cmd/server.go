package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/mmdaz/lime/server"
	"github.com/mmdaz/lime/version"
)

var banner = "license server\nversion " + version.Version + "\nhash:" + version.GitCommit

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start license server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(banner)
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
