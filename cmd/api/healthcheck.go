package main

import (
	"net/http"
)

// healthcheckhandler method writes a plain-text response with information about the
// application status, operating enviroment and version.
func (app *application) healthcheckhandler(w http.ResponseWriter, r *http.Request) {
	// Declare an envelope map containing the data for the response.
	env := envelope{
		"status":      "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
