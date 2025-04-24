package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"http_io_bound/internal/errlog"
	"net/http"
	"strings"
)

var stub = &cobra.Command{
	Use:   "stub",
	Short: "Create a task for stub server",
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := map[string]interface{}{
			"type": "stub",
		}
		data, _ := json.Marshal(payload)
		response, err := http.Post("http://localhost:8080/tasks",
			"application/json", bytes.NewReader(data))
		if err != nil {
			return err
		}
		defer response.Body.Close()

		location := response.Header.Get("Location")
		if !strings.HasPrefix(location, "/tasks/") {
			return fmt.Errorf("invalid Location header: %q", location)
		}
		errlog.Post(fmt.Sprintf("Created sleep task:", location[strings.LastIndex(location, "/")+1:]))
		return nil
	},
}
