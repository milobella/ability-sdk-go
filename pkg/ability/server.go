package ability

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

// Rule : routing rule
type Rule struct {
	condition func(request *Request) (result bool)
	process   func(request *Request, response *Response)
}

// Server : server
type Server struct {
	port     int
	router   *mux.Router
	listener net.Listener
	logger   loggo.Logger
	name     string
	rules    []Rule
}

// NewServer ctor
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

	server.router.HandleFunc("/resolve", server.handleResolve).Methods("POST")
	return server
}

func (s *Server) handleResolve(writer http.ResponseWriter, r *http.Request) {
	request, err := readRequest(r)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	response := new(Response)
	s.processRules(request, response)

	if err = writeResponse(writer, response); err!= nil {
		http.Error(writer, err.Error(), 500)
	}

	return
}

func (s *Server) processRules(request *Request, response *Response) {
	// Set the auto reprompt to false by default
	response.AutoReprompt = false

	// Loop on all rules
	for _, rule := range s.rules {
		if rule.condition(request) {
			rule.process(request, response)
			return
		}
	}

	response.Nlg.Sentence = "Error processing the request given to the ability."
	return
}

// RegisterIntent : Ability servers are done to handle some intents. We can simply register one using this method,
// giving the intent as string and a function handler which takes a request and a response.
//
// Deprecated: Use RegisterIntentRule instead.
func (s *Server) RegisterIntent(intent string, process func(request Request, response *Response)) (err error) {
	s.router.HandleFunc("/resolve/"+intent, func(w http.ResponseWriter, r *http.Request) {
		abRequest, err := readRequest(r)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		abResponse := new(Response)
		// Set the auto reprompt to false by default
		abResponse.AutoReprompt = false
		process(*abRequest, abResponse)
		err = writeResponse(w, abResponse)
	}).Methods("POST")

	return
}

// RegisterIntentRule : Create a rule of routing based on intent.
func (s *Server) RegisterIntentRule(intent string, process func(*Request, *Response)) {
	s.RegisterRule(func(request *Request) (result bool) {
		return request.Nlu.BestIntent == intent
	}, process)
}

// RegisterRule : Create a rule of routing based on condition on request.
func (s *Server) RegisterRule(condition func(request *Request) (result bool), process func(request *Request, response *Response)) {
	s.rules = append(s.rules, Rule{condition: condition, process: process})
}

func readRequest(r *http.Request) (request *Request, err error) {
	request = new(Request)
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}
	err = json.Unmarshal(b, request)
	return
}

func writeResponse(w http.ResponseWriter, response *Response) (err error) {
	err = json.NewEncoder(w).Encode(response)
	return
}

// Serve : Start the server
func (s *Server) Serve() {
	done := make(chan bool)
	go http.Serve(s.listener, s.router)
	s.logger.Infof("Successfully started the %s server on port %d !", s.name, s.port)
	<-done
}
