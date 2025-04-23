package errlog

import (
	"fmt"
	"log"
	"net/http"
)

func Check(msg string, err error, isCritical bool) {
	report, reportError := NewReport()
	if reportError != nil {
		log.Fatalln(reportError)
		return
	}
	defer report.File.Close()
	if err != nil {
		report.Error.Println(fmt.Sprintf("%s : %v", msg, err))
		if isCritical {
			log.Fatalf("[CRITICAL] %s : %v\n", msg, err)
		} else {
			log.Panicf("[PANIC] %s : %v\n", msg, err)
		}
	}
}

func Post(msg string) {
	report, reportError := NewReport()
	if reportError != nil {
		log.Fatalln(reportError)
		return
	}
	report.Info.Println(msg)
	log.Println(msg)
}

func HTTPCheck(writer http.ResponseWriter, msg string, err int) {
	report, reportError := NewReport()
	if reportError != nil {
		log.Fatalln(reportError)
		return
	}
	defer report.File.Close()
	report.Error.Println(fmt.Sprintf("%s : %d", msg, err))
	http.Error(writer, msg, err)
}
