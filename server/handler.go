package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func (s *server) Routes() {
	s.router.HandleFunc("/", s.CounterHandler(s.communication))
}

func (s *server) CounterHandler(com communication) http.HandlerFunc {
	var init sync.Once
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		init.Do(s.initialize)

		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		requestTimestamp := time.Now().Truncate(s.precision)
		s.logger.Printf("RequestTimestamp: '%v'\n", requestTimestamp.Format(time.RFC3339))

		com.exchangeTimestamp <- requestTimestamp
		totalRequestsSoFar := <-com.exchangeRequestCount
		s.logger.Printf("Response '%v'\n", totalRequestsSoFar)

		response := Response{
			timestamp:    totalRequestsSoFar.Timestamp,
			RequestCount: totalRequestsSoFar.TotalRequestsWithinTimeframe,
		}
		encodedCache, err := json.Marshal(response)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, ResponseError{ErrorMsg: err.Error()}.ToJSON())
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(encodedCache))
	})
}
