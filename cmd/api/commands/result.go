package commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"http_io_bound/internal/errlog"
	"io"
	"net/http"
)

var result = &cobra.Command{
	Use:   "result [task_id]",
	Short: "Fetch only the result of a completed task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		url := fmt.Sprintf("http://localhost:8080/tasks/%s/result", id)
		response, err := http.Get(url)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		switch response.StatusCode {
		case http.StatusOK:
			var body struct {
				Result string `json:"result"`
			}
			data, _ := io.ReadAll(response.Body)
			if err := json.Unmarshal(data, &body); err != nil {
				return err
			}
			errlog.Post(fmt.Sprintf("Result %s:", body.Result))
			return nil

		case http.StatusConflict:
			return fmt.Errorf("task %s is not finished yet", id)
		case http.StatusNotFound:
			return fmt.Errorf("task %s not found", id)
		default:
			data, _ := io.ReadAll(response.Body)
			return fmt.Errorf("error: %s", string(data))
		}
	},
}
