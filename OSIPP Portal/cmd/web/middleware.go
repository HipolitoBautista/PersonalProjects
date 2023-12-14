package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"


  // written by: Hipolito, Michael, Jahmur, Dennis, Rene 
  // tested by: Hipolito, Michael, Jahmur, Dennis, Rene 
  // debugged by: Hipolito, Michael, Jahmur, Dennis, Rene 

)

//Setting security to our header 
func securityHeadersMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)

	})

}

//Log requests made 
func (app *application) logRequestMiddlewear(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//when the request comes to me
		start := time.Now()
		app.infoLog.Printf("%s - %s %s %s ", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
		//when the responce comes to me
		app.infoLog.Printf("Request took %v", time.Since(start))

	})
}

//Trace what made us panic 
func (app *application) recoverPanicMiddlewear(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {

			if err := recover(); err != nil {
				w.Header().Set("connection", "close")
				trace := fmt.Sprintf("%s \n%s", err, debug.Stack())
				app.errorLog.Output(2, trace)
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

//User is required to be authneticated or he is redirected 
func (app *application) requireAuthnticationMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		w.Header().Add("cache-control", "no-store")
		next.ServeHTTP(w, r)
	})
}

//Admin is required to be authneticated or he is redirected 
func (app *application) requireAdminAuthnticationMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticatedAdmin(r) {
			http.Redirect(w, r, "/adminLogin", http.StatusSeeOther)
			return
		}
		w.Header().Add("cache-control", "no-store")
		next.ServeHTTP(w, r)
	})
}

//NoSurf for forms 
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(
		http.Cookie{
			HttpOnly: true,
			Path:     "/",
			Secure:   true,
		})
	return csrfHandler
}
