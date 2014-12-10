package golog

import (
	"fmt"
	color "github.com/ivpusic/go-clicolor/clicolor"
	"io"
)

type Appender interface {
	Append(log Log)
	Id() string
}

type Stdout struct {
}

var (
	out      io.Writer
	instance *Stdout
)

func (s *Stdout) Append(log Log) {
	msg := fmt.Sprintf("%s - [%s] - %s", log.Logger.Name, log.Level.Name, log.Message)
	color.Out = out
	color.Print(msg).In(log.Level.color)
}

func (s *Stdout) Id() string {
	return "github.com/ivpusic/golog/stdout"
}

func StdoutAppender() *Stdout {
	if instance == nil {
		instance = &Stdout{}
	}

	return instance
}
