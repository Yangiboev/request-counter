package server

import (
	"encoding/json"
	"log"
	"time"
)

type Response struct {
	timestamp    time.Time
	RequestCount int `json:"requestCount"`
}

func NewResponse(timestamp time.Time, requestCount int) Response {
	return Response{timestamp: timestamp, RequestCount: requestCount}
}

func (r Response) Timestamp() time.Time {
	return r.timestamp
}

type ResponseError struct {
	ErrorMsg string
}

func (r ResponseError) ToJSON() string {
	encodedError, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}
	return string(encodedError)
}
