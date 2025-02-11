package main

import (
	"os"

	"github.com/nicklasfrahm/homelab/cmd/labctl/config"
	"github.com/spf13/cobra"
)

var version = "dev"
var help bool

var rootCmd = &cobra.Command{
	Use:   "labctl",
	Short: "CLI to manage my homelab",
	Long: `A command line interface to manage infrastructure in my homelab.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if help {
			cmd.Help()
			os.Exit(0)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
	Version:      version,
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "display help for command")

	rootCmd.AddCommand(config.RootCommand())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
