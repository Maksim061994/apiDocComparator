package apiserver

import (
	"bytes"
	"io"
	"net/http"

	"github.com/api-doc-compare/internal/app/docx"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

func New(config *Config) *APIServer {
	var logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	return &APIServer{
		config: config,
		logger: logger,
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.configureRouter()
	s.logger.Info("starting api server")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/docs/compare", s.handleInputDocs()).Methods("POST")
	s.router.Use(s.loggingMiddleware)
}

func (s *APIServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("connection string: ", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) handleInputDocs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20)
		// processing input files
		s.logger.Info("start compare documents")
		filedocNew, _, err := r.FormFile("docNew")
		if err != nil {
			s.logger.Error(err)
		}
		defer filedocNew.Close()
		var docNew bytes.Buffer
		io.Copy(&docNew, filedocNew)

		filedocOld, _, err := r.FormFile("docOld")
		if err != nil {
			s.logger.Error(err)
		}
		defer filedocOld.Close()
		var docOld bytes.Buffer
		io.Copy(&docOld, filedocOld)

		typeResponse := r.FormValue("typeResponse")
		bodyResponse := docx.Comparator(docNew, docOld, typeResponse)
		w.Write([]byte(bodyResponse))
	}
}
