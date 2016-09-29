package inspector

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/apprenda/kismatic-platform/pkg/inspector/check"
)

// Server supports the execution of inspector rules from a remote node
type Server struct {
	// The Port the server will listen on
	Port int
	// NodeFacts are the facts that apply to the node where the server is running
	NodeFacts []string
	// RulesEngine for running inspector rules
	rulesEngine *Engine
}

type serverError struct {
	Error string
}

var executeEndpoint = "/execute"
var closeEndpoint = "/close"

// NewServer returns an inspector server that has been initialized
// with the default rules engine
func NewServer(nodeRole string, port int) (*Server, error) {
	s := &Server{
		Port: port,
	}
	distro, err := check.DetectDistro()
	if err != nil {
		return nil, fmt.Errorf("error building server: %v", err)
	}
	s.NodeFacts = []string{nodeRole, string(distro)}
	pkgMgr, err := check.NewPackageManager(distro)
	if err != nil {
		return nil, fmt.Errorf("error building server: %v", err)
	}
	engine := &Engine{
		RuleCheckMapper: DefaultCheckMapper{
			PackageManager: pkgMgr,
		},
	}
	s.rulesEngine = engine
	return s, nil
}

// Start the server
func (s *Server) Start() error {
	mux := http.NewServeMux()
	// Execute endpoint
	mux.HandleFunc(executeEndpoint, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// Decode rules
		rules := []Rule{}
		err := json.NewDecoder(req.Body).Decode(rules)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error decoding rules when processing request: %v", err)
			return
		}
		defer req.Body.Close()
		// Run the rules that we received
		results, err := s.rulesEngine.ExecuteRules(rules, s.NodeFacts)
		if err != nil {
			err = json.NewEncoder(w).Encode(serverError{Error: err.Error()})
			if err != nil {
				log.Printf("error writing server response: %v\n", err)
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(results)
		if err != nil {
			log.Printf("error writing server response: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	// Close endpoint
	mux.HandleFunc(closeEndpoint, func(w http.ResponseWriter, req *http.Request) {
		err := s.rulesEngine.CloseChecks()
		if err != nil {
			log.Printf("error closing checks: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), mux)
}
