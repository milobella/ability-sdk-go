package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/milobella/ability-sdk-go/internal/logging"
	"github.com/milobella/ability-sdk-go/pkg/model"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
}

type routingRule struct {
	condition func(request *model.Request) (result bool)
	process   func(request *model.Request, response *model.Response)
}

type Server struct {
	port  int
	name  string
	rules []routingRule
	e     *echo.Echo
}

// New creates an ability server
func New(name string, port int) *Server {
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
	request := new(model.Request)
	if err = c.Bind(request); err != nil {
		return
	}

	response := new(model.Response)
	s.processRules(request, response)

	return c.JSON(http.StatusOK, response)
}

func (s *Server) processRules(request *model.Request, response *model.Response) {
	logrus.Debugf("Received request: %v", request)
	defer func() {
		logrus.Debugf("Sent response : %v", response)
	}()

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
	response.Nlg.Sentence = "Error processing the request given to the server."
	return
}

// Register : Create a routing rule based on condition on request.
func (s *Server) Register(condition func(request *model.Request) (result bool), process func(request *model.Request, response *model.Response)) {
	s.rules = append(s.rules, routingRule{condition: condition, process: process})
}

// Serve : Start the server
func (s *Server) Serve() {
	// Run the echo server
	logrus.Fatal(s.e.Start(fmt.Sprintf(":%d", s.port)))
}
