package commands

import (
	"github.com/spf13/cobra"
	"http_io_bound/internal/errlog"
	"io"
	"net/http"
)

var clean = &cobra.Command{
	Use:   "clean",
	Short: "Clear old tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		response, err := http.Get("http://localhost:8080/tasks/clear")
		if err != nil {
			return err
		}
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		errlog.Post(string(body))
		return nil
	},
}
