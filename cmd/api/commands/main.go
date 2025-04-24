package commands

import (
	"github.com/spf13/cobra"
	"http_io_bound/cmd/api/commands/create"
)

var root = &cobra.Command{
	Use:   "api",
	Short: "CLI for HTTP API",
}

func Execute() error {
	root.AddCommand(health, create.Cmd, clean, status, result, list)
	return root.Execute()
}
