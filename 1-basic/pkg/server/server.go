package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var indexPage = `
<!DOCTYPE html>
<html>
	<body>
		<h1 style="text-align:center;" > User Database </h1>
		<p style="text-align:center;" > Welcome to the user database. </p>
	</body>
</html>
`

// user represents the JSON value that's sent as a response to a request.
type user struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// userinfo is the information that is stored per user.
type userinfo struct {
	email string
	age   int
}

// Server is an HTTP server.
type Server struct {
	users map[string]userinfo // key -> username
}

// New is a new server.
func New() *Server {
	return &Server{
		users: make(map[string]userinfo),
	}
}

// HandleIndex handles the index path "/".
func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(indexPage))
}

// HandleCreateUsers handles the path "/users/create".
// Create -> Post/Put.
func (s *Server) HandleCreateUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodPut:
		// Check that the input is JSON.
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Could not read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError) // HTTP 500
			return
		}
		defer r.Body.Close()

		// Unmarshal the body.
		var u user
		err = json.Unmarshal(body, &u)
		if err != nil {
			log.Printf("Could not unmarshal request body: %v", err)
			w.WriteHeader(http.StatusBadRequest) // HTTP 400
			return
		}

		log.Printf("Create User: %v", u.Name)
		s.users[u.Name] = userinfo{
			email: u.Email,
			age:   u.Age,
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed) // HTTP 415
	}
}

// HandleUsers handles the `/users/{name}` request.
func (s *Server) HandleUsers(w http.ResponseWriter, r *http.Request) {
	// Fetch the name from the params. Common for all methods of this route.
	params := mux.Vars(r)
	name := params["name"]
	u, ok := s.users[name]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		ret := user{
			Name:  name,
			Email: u.email,
			Age:   u.age,
		}
		msg, err := json.Marshal(ret)
		if err != nil {
			log.Printf("Could not marshal request: %v", err)
			w.WriteHeader(http.StatusInternalServerError) // HTTP 500
			return
		}
		log.Printf("Get user: %s", name)
		w.Header().Add("Content-Type", "application/json")
		w.Write(msg)
	case http.MethodPatch:
		// Partial updates.
		// Check that the input is JSON.
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Could not read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError) // HTTP 500
			return
		}
		defer r.Body.Close()

		// Unmarshal the body.
		var u user
		err = json.Unmarshal(body, &u)
		if err != nil {
			log.Printf("Could not unmarshal request body: %v", err)
			w.WriteHeader(http.StatusBadRequest) // HTTP 400
			return
		}

		log.Printf("Update user: %s", name)

		userinfo := s.users[name]
		if u.Age != 0 {
			userinfo.age = u.Age
		}
		if u.Email != "" {
			userinfo.email = u.Email
		}
		s.users[name] = userinfo
	case http.MethodDelete:
		log.Printf("Delete user: %s", name)
		delete(s.users, name)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed) // HTTP 415
	}
}
