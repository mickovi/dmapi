package main

import (
	"fmt"
	"net/http"
)

// healthcheckhandler method writes a plain-text response with information about the
// application status, operating enviroment and version.
func (app *application) healthcheckhandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}
