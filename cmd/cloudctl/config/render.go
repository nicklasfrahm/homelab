package config

import "github.com/spf13/cobra"

func RenderCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build <directory>",
		Short: "Build configuration into static files",
		Long: `Build configuration into static files
that can be served by a web server.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(ValidateCommand())

	return cmd
}
