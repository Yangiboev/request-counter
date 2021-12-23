package server

import (
	"time"
)

type Cache struct {
	RequestCount
	TotalRequestsWithinTimeframe int
}

func NewCache(timestamp time.Time, totalAccumulated int) Cache {
	requestCount := RequestCount{Timestamp: timestamp}
	requestCount.Increment()
	return Cache{
		RequestCount:                 requestCount,
		TotalRequestsWithinTimeframe: totalAccumulated + 1,
	}
}

func (c *Cache) Increment() {
	c.RequestCount.Increment()
	c.TotalRequestsWithinTimeframe++
}
