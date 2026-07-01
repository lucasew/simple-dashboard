package godashboard

import (
	"log"
)

// ReportError is the centralized error reporting function.
// It logs the error with its context and can be expanded in the future
// to report to external error tracking services like Sentry.
func ReportError(err error, contextMsg string) {
	if err == nil {
		return
	}
	log.Printf("[ERROR] %s: %v\n", contextMsg, err)
}
