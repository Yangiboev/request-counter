package config

import (
	"flag"
	"time"
)

type Config struct {
	ListenAddress        string
	PersistenceFile      string
	PersistenceTimeFrame time.Duration
	Precision            time.Duration
}

func LoadConfig() Config {
	var (
		config               Config
		precision            string
		persistenceTimeframe string
		err                  error
	)
	flag.StringVar(&config.ListenAddress, "port", ":5001", "Server listens to the port")
	flag.StringVar(&persistenceTimeframe, "timeframe", "60s", "Time frame to calculate hit counts")
	flag.StringVar(&precision, "precision", "100ms", "Timestamps that differ ")
	flag.StringVar(&config.PersistenceFile, "persistence-file", "persistence.bin", "File stores all values after the termination")
	flag.Parse()

	config.PersistenceTimeFrame, err = time.ParseDuration(persistenceTimeframe)
	if err != nil {
		panic(err)
	}

	config.Precision, err = time.ParseDuration(precision)
	if err != nil {
		panic(err)
	}

	return config
}
