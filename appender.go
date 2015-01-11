package golog

import (
	"fmt"
	"io"
)

// Interface for implementing custom appenders.
type Appender interface {
	// method for injecting log to some source
	// when appender receives Log instance through this method,
	// it should decide what to do with log
	Append(log Log)

	// method will return appender ID
	// it will be used for disabling appenders
	Id() string
}

// Representing stdout appender.
type Stdout struct {
	dateformat string
}

var (
	instance *Stdout
	out      io.Writer
)

// Appending logs to stdout.
func (s *Stdout) Append(log Log) {
	fmt.Printf("%s %s [%s]: %s\n",
		log.Level.String()[:4],
		log.Logger.Name,
		log.Time.Format(s.dateformat),
		log.Message)
}

// Getting Id of stdout appender
// Id of default stdout appender is "github.com/ivpusic/golog/stdout"
func (s *Stdout) Id() string {
	return "github.com/ivpusic/golog/stdout"
}

// Function for creating and returning new stdout appender instance.
func StdoutAppender() *Stdout {
	if instance == nil {
		instance = &Stdout{
			dateformat: "2006-01-02 15:04:05",
		}
	}

	return instance
}
