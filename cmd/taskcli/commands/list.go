package commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"http_io_bound/internal/errlog"
	"io"
	"net/http"
)

var list = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		response, err := http.Get("http://localhost:8080/tasks")
		if err != nil {
			return err
		}
		defer response.Body.Close()
		data, _ := io.ReadAll(response.Body)
		var tasks []map[string]interface{}
		if err := json.Unmarshal(data, &tasks); err != nil {
			return err
		}
		for _, t := range tasks {
			errlog.Post(fmt.Sprintf("- %s %s\n", t["id"], t["status"]))
		}
		return nil
	},
}
