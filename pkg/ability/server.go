package ability

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/milobella/ability-sdk-go/internal/logging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
}

type RoutingRule struct {
	condition func(request *Request) (result bool)
	process   func(request *Request, response *Response)
}

type Server struct {
	port  int
	name  string
	rules []RoutingRule
	e     *echo.Echo
}

func NewServer(name string, port int) *Server {
	// Initialize an echo server
	e := echo.New()
	server := new(Server)
	server.name = name
	server.port = port
	server.e = e

	// Register logging middleware
	logging.ApplyMiddleware(e)
	// Register route for metrics
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Register the main route of an ability
	e.POST("/resolve", server.handleResolve)

	return server
}

func (s *Server) handleResolve(c echo.Context) (err error) {
	request := new(Request)
	if err = c.Bind(request); err != nil {
		return
	}

	response := new(Response)
	s.processRules(request, response)

	return c.JSON(http.StatusOK, response)
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

	logrus.WithField("request", request).Error("Didn't find any rule matching the request")
	response.Nlg.Sentence = "Error processing the request given to the ability."
	return
}

// RegisterIntentRule : Create a routing rule based on intent.
func (s *Server) RegisterIntentRule(intent string, process func(*Request, *Response)) {
	s.RegisterRule(func(request *Request) (result bool) {
		return request.Nlu.BestIntent == intent
	}, process)
}

// RegisterRule : Create a routing rule based on condition on request.
func (s *Server) RegisterRule(condition func(request *Request) (result bool), process func(request *Request, response *Response)) {
	s.rules = append(s.rules, RoutingRule{condition: condition, process: func(request *Request, response *Response) {
		logrus.Debugf("Received request: %v", request)
		process(request, response)
		logrus.Debugf("Sent response : %v", response)
	}})
}

// Serve : Start the server
func (s *Server) Serve() {
	// Run the echo server
	logrus.Fatal(s.e.Start(fmt.Sprintf(":%d", s.port)))
}
