package appenders

import (
	"fmt"
	"github.com/ildus/golog"
	"github.com/ildus/golog/heka_emitter"
	"io"
	"net"
)

type HekaAppender struct {
	Proto      string
	Addr       string
	EnvVersion string
	Type       string
	conn       io.Writer
	emitter    *heka_emitter.ProtobufEmitter
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

		ha.emitter = heka_emitter.NewProtobufEmitter(ha.conn,
			ha.EnvVersion, "", log.Logger.Name)
	}

	// if additional data contains maps or errors collect them into one map
	var logFields map[string]string
	if log.Data != nil {
		for _, item := range log.Data {
			switch item.(type) {
			case map[string]string:
				{
					for key, val := range item.(map[string]string) {
						logFields[key] = val
					}
				}
			case error:
				{
					logFields["error"] = item.(error).Error()
				}
			}
		}
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
