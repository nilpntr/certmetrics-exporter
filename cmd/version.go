package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of CertMetrics Exporter",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("CertMetrics Exporter version: %s", version))
	},
}
