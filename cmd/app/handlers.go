package main

import (
	"fmt"
	"net/http"
)

// Handler for GET /health
func (a *app) getHealth(w http.ResponseWriter, _ *http.Request) {
	output := "status: ok\n"
	output = fmt.Sprintf(output+"environment: %s\n", a.config.env)
	output = fmt.Sprintf(output+"versioin: %s\n", version)

	_, err := fmt.Fprint(w, output)
	if err != nil {
		a.logger.Error(err.Error())
	}
}
