package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"text/template"
)


  // written by: Hipolito, Michael, Jahmur, Dennis, Rene 
  // tested by: Hipolito, Michael, Jahmur, Dennis, Rene 
  // debugged by: Hipolito, Michael, Jahmur, Dennis, Rene 
  // etc


//remders our templates (.tmpl files)
func RenderTemplate(w http.ResponseWriter, tmpl string, data *templateData) {

	ts, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	//Executes the template with any template data we pass 
	err = ts.Execute(w, data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)

	}
}

//Deals with server errors 
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	// deal with the error status
	http.Error(w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

//Deals with client errors 
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

//deals with Not found errors 
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
//Checks to see if the user is authenticated
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionsManager.Exists(r.Context(), "authenticatedUserID")
}
//Checks to see if the admin is authenticated 
func (app *application) isAuthenticatedAdmin(r *http.Request) bool {
	return app.sessionsManager.Exists(r.Context(), "authenticatedAdminID")
}
