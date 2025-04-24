package commands

import (
	"github.com/spf13/cobra"
	"http_io_bound/cmd/taskcli/commands/create"
)

var root = &cobra.Command{
	Use:   "taskcli",
	Short: "CLI for HTTP API",
}

func Execute() error {
	root.AddCommand(health, create.Cmd, clean, status, result, list)
	return root.Execute()
}
