package ability

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"milobella/oratio/pkg/anima"
	"milobella/oratio/pkg/cerebro"
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


func (s *Server) RegisterIntent(intent string, process func(nlu cerebro.NLU, nlg *anima.NLG)) (err error) {
	s.router.HandleFunc("/resolve/" + intent, func(w http.ResponseWriter, r *http.Request) {
		nlu, err := readNLU(r)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		nlg := new(anima.NLG)
		process(nlu, nlg)
		writeNLG(w, nlg)
	}).Methods("GET")

	return
}

func readNLU(r *http.Request) (nlu cerebro.NLU, err error){
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &nlu)
	return
}

func writeNLG(w http.ResponseWriter, nlg *anima.NLG) (err error) {
	json.NewEncoder(w).Encode(nlg)
	return
}

func (s *Server) Serve(){
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router))
}