package ability

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Server struct {
	port   int
	router *mux.Router
}

func NewServer(port int) *Server {
	server := new(Server)
	server.port = port
	server.router = mux.NewRouter()
	return server
}

func (s *Server) RegisterIntent(intent string, process func(request Request, response *Response)) (err error) {
	s.router.HandleFunc("/resolve/" + intent, func(w http.ResponseWriter, r *http.Request) {
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

func readRequest(r *http.Request) (request Request, err error){
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

func (s *Server) Serve(){
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router))
}