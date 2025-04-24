package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"http_io_bound/internal/errlog"
	"net/http"
)

var health = &cobra.Command{
	Use:   "health",
	Short: "Check service health",
	RunE: func(cmd *cobra.Command, args []string) error {
		response, err := http.Get("http://localhost:8080/health")
		if err != nil {
			return err
		}
		defer response.Body.Close()
		if response.StatusCode == 200 {
			errlog.Post("[HEALTH] OK")
		} else {
			errlog.Post(fmt.Sprintf("[HEALTH][FAIL] %s", response.Status))
		}
		return nil
	},
}
