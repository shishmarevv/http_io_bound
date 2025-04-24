package create

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task of a given type",
	// Здесь нет RunE: обработку делаем в подпакомандах
}

func init() {
	Cmd.AddCommand(
		stub,
	)
}
