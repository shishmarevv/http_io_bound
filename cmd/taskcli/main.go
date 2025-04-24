package main

import (
	"http_io_bound/cmd/taskcli/commands"
	"http_io_bound/internal/errlog"
)

func main() {
	err := commands.Execute()
	errlog.Check("CLI Error : ", err, false)
}
