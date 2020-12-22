package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/norwoodj/helm-docs/pkg/helm"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"istio.io/pkg/log"
)

type server struct {
	router *mux.Router
}

// Start simply starts the HTTP server
func Start() {
	s := &server{}
	s.router = mux.NewRouter()
	s.routes()
	s.start()
}

func (s *server) start() {
	srv := &http.Server{
		Handler:      s.router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func (s *server) clientError(e error) http.HandlerFunc {
	log.Errorf(e.Error())
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", e.Error())
	}
}

func (s *server) handleSomething() http.HandlerFunc {
	chartSearchRoot := viper.GetString("chart-search-root")
	var fullChartSearchRoot string

	if path.IsAbs(chartSearchRoot) {
		fullChartSearchRoot = chartSearchRoot
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return s.clientError(err)
		}

		fullChartSearchRoot = path.Join(cwd, chartSearchRoot)
	}

	chartDirs, err := helm.FindChartDirectories(fullChartSearchRoot)
	if err != nil {
		s.clientError(err)
		os.Exit(1)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(chartDirs)
	}
}

func (s *server) handleChartInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		p := path.Join(vars["repository"], vars["chart"])
		info, err := helm.ParseChartInformation(p)
		if err != nil {
			s.clientError(err)
		}
		json.NewEncoder(w).Encode(info)
	}
}

func (s *server) handleRenderValues() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		p := path.Join(vars["repository"], vars["chart"])
		info, err := helm.ParseChartInformation(p)
		if err != nil {
			s.clientError(err)
		}
		yaml.NewEncoder(w).Encode(info.ChartValues)
	}
}
