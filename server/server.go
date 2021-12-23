package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Yangiboev/request-counter/config"
)

type key int

const (
	requestIDKey key = 0
)

func nextRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

type server struct {
	router               *http.ServeMux
	logger               *log.Logger
	communication        communication
	persistenceTimeFrame time.Duration
	precision            time.Duration
	persistenceFile      string
	http.Server
}

func NewServer(env config.Config) *server {
	router := http.NewServeMux()
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	errorLogger := log.New(os.Stderr, "http: ", log.LstdFlags)
	communication := NewCommunication()
	server := &server{
		router:               router,
		logger:               logger,
		communication:        communication,
		persistenceTimeFrame: env.PersistenceTimeFrame,
		precision:            env.Precision,
		persistenceFile:      env.PersistenceFile,
		Server: http.Server{
			Addr:         env.ListenAddress,
			Handler:      tracing(nextRequestID)(logging(logger)(router)),
			ErrorLog:     errorLogger,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}

	return server
}

func (s *server) readStateFromDisk() {
	if _, err := os.Open(s.persistenceFile); err != nil {
		s.logger.Printf("No state file could be found under '%v': %v. Will work on a clean slate.\n", s.persistenceFile, err)
	} else {
		s.logger.Printf("Reading last state from file '%v'...\n", s.persistenceFile)
		s.communication.state, err = ReadFromFile(s.persistenceFile)
		if err != nil {
			s.logger.Printf("Could not read state from file '%v': %v\n", s.persistenceFile, err)
			return
		}
		s.logger.Printf("State restored. Current request count: '%v'\n", s.communication.state.Present.TotalRequestsWithinTimeframe)

		s.logger.Println("Removing file...")
		if err := os.Remove(s.persistenceFile); err != nil {
			s.logger.Printf("Could not remove state file '%v': %v\n", s.persistenceFile, err)
		}
	}
}

func (s *server) PersistState() error {
	s.logger.Printf("Persisting state '%+v' to file '%v'.", s.communication.state, s.persistenceFile)
	if err := s.communication.state.WriteToFile(s.persistenceFile); err != nil {
		return err
	}
	return nil
}

func (s *server) initialize() {
	s.logger.Print("Initialising server with following parameters:")
	s.logger.Printf("Persistence File: '%v'\n", s.persistenceFile)
	s.logger.Printf("Persistence Timeframe: '%v'\n", s.persistenceTimeFrame)
	s.logger.Printf("Precision: '%v'\n", s.precision)
	s.readStateFromDisk()
	s.startCommunicationProcessor()
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
