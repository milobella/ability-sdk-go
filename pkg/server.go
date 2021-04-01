package pkg

import (
	"fmt"
	"github.com/milobella/ability-sdk-go/internal/logging"
	"github.com/milobella/ability-sdk-go/pkg/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// TODO: read it in the config when move to viper
	logrus.SetLevel(logrus.DebugLevel)

	logrus.SetReportCaller(true)
}

// Rule : routing rule
type Rule struct {
	condition func(request *model.Request) (result bool)
	process   func(request *model.Request, response *model.Response)
}

// Server : server
type Server struct {
	port  int
	name  string
	rules []Rule
	e     *echo.Echo
}

// NewServer ctor
func NewServer(name string, port int) *Server {
	// Initialize an echo server
	e := echo.New()
	server := new(Server)
	server.name = name
	server.port = port
	server.e = e

	// Register logging middleware
	e.Use(logging.InitializeLoggingMiddleware().Handle)
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
func (s *Server) RegisterIntent(intent string, process func(request model.Request, response *model.Response)) error {
	s.e.POST("/resolve/"+intent, func(c echo.Context) (err error) {
		abRequest := new(model.Request)
		if err = c.Bind(abRequest); err != nil {
			return
		}
		abResponse := new(model.Response)
		// Set the auto reprompt to false by default
		abResponse.AutoReprompt = false

		process(*abRequest, abResponse)

		return c.JSON(http.StatusOK, abResponse)
	})

	return nil
}

// RegisterIntentRule : Create a rule of routing based on intent.
func (s *Server) RegisterIntentRule(intent string, process func(*model.Request, *model.Response)) {
	s.RegisterRule(func(request *model.Request) (result bool) {
		return request.Nlu.BestIntent == intent
	}, process)
}

// RegisterRule : Create a rule of routing based on condition on request.
func (s *Server) RegisterRule(condition func(request *model.Request) (result bool), process func(request *model.Request, response *model.Response)) {
	s.rules = append(s.rules, Rule{condition: condition, process: func(request *model.Request, response *model.Response) {
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
