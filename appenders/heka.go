package appenders

import (
	"./heka_appender"
	"fmt"
	"github.com/ildus/golog"
	"io"
	"net"
)

type HekaAppender struct {
	Proto      string
	Addr       string
	EnvVersion string
	Type       string
	conn       io.Writer
	emitter    *heka_appender.ProtobufEmitter
}

func (fa *HekaAppender) Id() string {
	return "github.com/ildus/golog/appenders/heka"
}

func (ha *HekaAppender) Append(log golog.Log) {
	if ha.conn == nil {
		if len(ha.Addr) == 0 {
			fmt.Printf("Missing remote host")
			return
		}
		var err error
		ha.conn, err = net.Dial(ha.Proto, ha.Addr)
		if err != nil {
			fmt.Printf("Error dialing host %q: %s", ha.Addr, err)
			return
		}

		ha.emitter = heka_appender.NewProtobufEmitter(ha.conn,
			ha.EnvVersion, "", log.Logger.Name)
	}

	ha.emitter.Emit(int32(log.Level), ha.Type, log.Message, nil)
}

func Heka(cnf golog.Conf) *HekaAppender {
	return &HekaAppender{
		Addr:       cnf["addr"],
		Proto:      cnf["proto"],
		EnvVersion: cnf["env_version"],
		Type:       cnf["message_type"],
	}
}
