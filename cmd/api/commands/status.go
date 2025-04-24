package commands

import (
	"github.com/spf13/cobra"
	"http_io_bound/internal/errlog"
	"io"
	"net/http"
)

var status = &cobra.Command{
	Use:   "status [task_id]",
	Short: "Get status of a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		response, err := http.Get("http://localhost:8080/tasks/" + id)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		errlog.Post(string(body))
		return nil
	},
}
