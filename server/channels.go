package server

import (
	"time"
)

type communication struct {
	state                State
	exchangeTimestamp    chan time.Time
	exchangeRequestCount chan Cache
	exchangePersistence  chan persistenceData
	exchangeAccumulated  chan int
}

func NewCommunication() communication {
	return communication{
		exchangeTimestamp:    make(chan time.Time),
		exchangeRequestCount: make(chan Cache),
		exchangePersistence:  make(chan persistenceData),
		exchangeAccumulated:  make(chan int),
	}
}

type persistenceData struct {
	RequestCount RequestCount
	Reference    RequestCount
}

func NewPersistenceData(cache Cache, timestamp time.Time) persistenceData {
	return persistenceData{
		RequestCount: cache.RequestCount,
		Reference: RequestCount{
			Timestamp: timestamp,
		},
	}
}

func (s *server) startCommunicationProcessor() {
	s.logger.Print("Starting communication processor...")

	s.logger.Print("Starting Persistence-Accumulated exchanger...")
	go func() {
		for {
			persistenceData, ok := <-s.communication.exchangePersistence
			if ok {
				s.communication.state.Past = s.communication.state.Past.AppendToTail(persistenceData.RequestCount)
				s.communication.state.Past = s.communication.state.Past.UpdateTotals(persistenceData.Reference, s.persistenceTimeFrame, s.precision)
				s.communication.exchangeAccumulated <- s.communication.state.Past.TotalAccumulatedRequestCount()
			} else {
				break
			}
		}
	}()

	s.logger.Print("Starting Timestamp-RequestCount exchanger...")
	go func() {
		for {
			requestTimestamp, ok := <-s.communication.exchangeTimestamp
			if ok {
				if s.communication.state.Present.Empty() {
					s.communication.state.Present.Timestamp = requestTimestamp
				}

				if s.communication.state.Present.CompareTimestampWithPrecision(requestTimestamp, s.precision) {
					s.communication.state.Present.Increment()
				} else {
					persistenceUpdate := NewPersistenceData(s.communication.state.Present, requestTimestamp)

					s.communication.exchangePersistence <- persistenceUpdate
					totalAccumulated := <-s.communication.exchangeAccumulated

					s.communication.state.Present = NewCache(requestTimestamp, totalAccumulated)
				}

				s.communication.exchangeRequestCount <- s.communication.state.Present
			} else {
				break
			}
		}
	}()
	s.logger.Print("Communication processor up and running")
}

func (s *server) CloseChannels() {
	close(s.communication.exchangeRequestCount)
	close(s.communication.exchangePersistence)
	close(s.communication.exchangeAccumulated)
}
