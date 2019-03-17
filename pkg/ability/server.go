package ability

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
	"io/ioutil"
	"net"
	"net/http"
)

type Server struct {
	port     int
	router   *mux.Router
	listener net.Listener
	logger   loggo.Logger
	name     string
}

func NewServer(name string, port int) *Server {
	// Initialize the server object
	server := new(Server)
	server.router = mux.NewRouter()
	server.name = name
	server.port = port

	// Initialize a logger from the given name, it will be used for all server logs
	logger := loggo.GetLogger(name)
	server.logger = logger

	// Initialize the listener
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Criticalf("Error initializing the server : %s", err)
	}
	server.listener = listener
	return server
}

// RegisterIntent : Ability servers are done to handle some intents. We can simply register one using this method,
// giving the intent as string and a function handler which takes a request and a response.
func (s *Server) RegisterIntent(intent string, process func(request Request, response *Response)) (err error) {
	s.router.HandleFunc("/resolve/"+intent, func(w http.ResponseWriter, r *http.Request) {
		abRequest, err := readRequest(r)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		abResponse := new(Response)
		process(abRequest, abResponse)
		writeResponse(w, abResponse)
	}).Methods("POST")

	return
}

func readRequest(r *http.Request) (request Request, err error) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &request)
	return
}

func writeResponse(w http.ResponseWriter, response *Response) (err error) {
	json.NewEncoder(w).Encode(response)
	return
}

// Serve : Start the server
func (s *Server) Serve() {
	done := make(chan bool)
	go http.Serve(s.listener, s.router)
	s.logger.Infof("Successfully started the %s server on port %d !", s.name, s.port)
	<-done
}
