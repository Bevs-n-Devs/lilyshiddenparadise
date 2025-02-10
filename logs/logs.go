package logs

import "log"

var logChannel = make(chan string) // channel to send logs to

const (
	info   = "INFO: "
	warn   = "WARNING! "
	logErr = "ERROR! "
)

func LogProcessor() {
	for logMessage := range logChannel {
		log.Println(logMessage)
	}
}

// logType: 1 = info, 2 = warning, 3 = error
func Logs(logType int, logMessage string) {
	var loggedMessage string
	switch logType {
	case 1:
		loggedMessage = info + logMessage
	case 2:
		loggedMessage = warn + logMessage
	case 3:
		loggedMessage = logErr + logMessage
	}
	logChannel <- loggedMessage
}
