package routes

import (
	"net/http"

	"github.com/MikMuellerDev/smarthome/core/database"
	"github.com/MikMuellerDev/smarthome/server/api"
	"github.com/MikMuellerDev/smarthome/server/middleware"
	"github.com/MikMuellerDev/smarthome/server/templates"
)

// Redirects to the dashboard
func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dash", http.StatusSeeOther)
}

// Serves HTML for the dashboard
func dashGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "dash.html", http.StatusOK)
}

// Serves HTML for rooms
func roomsGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "rooms.html", http.StatusOK)
}

// Serves HTML for reminders
func reminderGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "reminders.html", http.StatusOK)
}

// Serves HTML for user management
func usersGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "users.html", http.StatusOK)
}

// If not user is logged in, it serves the HTML for the login page
// Otherwise the user is redirected to the dashboard
func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.Store.Get(r, "session")
	if err != nil {
		log.Debug("Session exists but could not be decoded: ", err.Error())
	}

	loginValidTemp, loginValidOkTemp := session.Values["valid"]
	loginValid, loginValidOk := loginValidTemp.(bool)

	usernameTemp, usernameTempOk := session.Values["username"]
	usernameSession, usernameSessionOk := usernameTemp.(string)

	if loginValidOkTemp && loginValidOk && loginValid {
		if usernameTempOk && usernameSessionOk && usernameSession != "" {
			_, exists, err := database.GetUserByUsername(usernameSession)
			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				api.Res(w, api.Response{Success: false, Message: "Could not check user validity", Error: "database failure"})
				return
			}
			if exists {
				// The session is valid: redirect to the dashboard
				http.Redirect(w, r, "/dash", http.StatusFound)
				return
			}
		}
	}
	// If no user is logged in, serve the login screen
	templates.ExecuteTemplate(w, "login.html", http.StatusOK)
}

// Serves HTML for profile settings
func userProfileGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "profile.html", http.StatusOK)
}
