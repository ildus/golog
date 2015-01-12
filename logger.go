package golog

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type LogLevel int32

// syslog severity levels.
const (
	EMERGENCY LogLevel = iota
	ALERT
	CRITICAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

var levelNames = map[LogLevel]string{
	EMERGENCY: "EMERGENCY",
	ALERT:     "ALERT",
	CRITICAL:  "CRITICAL",
	ERROR:     "ERROR",
	WARNING:   "WARNING",
	NOTICE:    "NOTICE",
	INFO:      "INFO",
	DEBUG:     "DEBUG",
}

func (l LogLevel) String() string {
	return levelNames[l]
}

var (
	// limit when logger name will be normalized
	// normalized names are shown in console using stdout appender
	maxnamelen = 20
	curnamelen = 7

	// supported name separators
	separators []byte = []byte{'/', '.', '-'}
)

// Representing one Log instance
type Log struct {
	// date and time of log
	Time time.Time `json:"time"`

	// logged message
	Message string `json:"message"`

	// log level
	Level LogLevel `json:"level"`

	// additional data sent to log
	// this part should be handled by appenders
	// appender can decide to ignore data or to store it on specific way
	Data []interface{} `json:"data"`

	// id of process which made log
	Pid int `json:"pid"`

	// logger instance
	Logger *Logger `json:"logger"`
}

// Representing one logger instance
// Logger can have multiple appenders, it can enable it,
// or disable it. Also you can define level which will be specific to this logger.
type Logger struct {
	// list of appenders
	appenders []Appender

	// is logged disabled
	disabled bool

	// name of logger
	// logger name will be shown in stdout appender output
	// also it can be used to enable/disable logger
	Name string `json:"name"`

	// minimum level of log to be shown
	Level LogLevel `json:"-"`

	// if this flag is set to true, in case any errors in appender
	// appender should panic. This also depends on appender implementation,
	// so appender can decide to ignore or to accept information in this flag
	DoPanic bool `json:"-"`
}

// Making and sending log entry to appenders if log level is appropriate.
func (l *Logger) Log(lvl LogLevel, msg interface{}, data []interface{}) {
	if l.disabled {
		return
	}

	if lvl <= l.Level {
		log := Log{
			Time:    time.Now(),
			Message: l.toString(msg),
			Level:   lvl,
			Data:    data,
			Logger:  l,
			Pid:     os.Getpid(),
		}

		for _, appender := range l.appenders {
			appender.Append(log)
		}
	}
}

func (l *Logger) toString(object interface{}) string {
	return fmt.Sprintf("%v", object)
}

// method will normalize names if they are too big or too short
// normal name length if defined by namelen variable
func (l *Logger) normalizeName() {
	length := len(l.Name)

	// name is ok as it is
	if length == maxnamelen || length == curnamelen {
		return
	}

	// name is too short, add some spaces
	if length < curnamelen {
		l.normalizeNameLen()
		return
	}

	// name is too long
	// do best to normalize it

	var (
		normalized string
		parts      []string
		separator  byte
	)

	// try split long name using different separators
	// this first one which can split name into smaller parts will be used
	for _, sep := range separators {
		parts = strings.Split(l.Name, string(sep))
		if len(parts) > 1 {
			separator = sep
			break
		}
	}

	// if we sucesufully splitted string into multiple parts
	if len(parts) > 1 {
		appendSeparator := true

		for i, str := range parts {
			// if part length is bigger than zero
			switch len(str) {
			case 0:
				appendSeparator = false
				break
			case 1:
				normalized += str[:1]
				break
			case 2:
				normalized += str[:2]
				break
			default:
				normalized += str[:3]
				break
			}

			if appendSeparator && (i != (len(parts) - 1)) {
				normalized += string(separator)
			}
		}

		// if still to long
		if len(normalized) > maxnamelen {
			normalized = normalized[:maxnamelen]
		}
	} else {
		length := len(l.Name)
		if length > maxnamelen {
			normalized = l.Name[:maxnamelen]
		} else {
			normalized = l.Name[0:length]
		}
	}

	l.Name = normalized
	if len(normalized) >= curnamelen {
		curnamelen = len(normalized)
	} else {
		l.normalizeNameLen()
	}
}

// if name is still to short we will add spaces
func (l *Logger) normalizeNameLen() {
	length := len(l.Name)
	missing := curnamelen - length
	for i := 0; i < missing; i++ {
		l.Name += " "
	}
}

// Making log with DEBUG level.
func (l *Logger) Debug(msg interface{}, data ...interface{}) {
	l.Log(DEBUG, msg, data)
}

// Making log with INFO level.
func (l *Logger) Info(msg interface{}, data ...interface{}) {
	l.Log(INFO, msg, data)
}

// Making log with WARN level.
func (l *Logger) Warn(msg interface{}, data ...interface{}) {
	l.Log(WARNING, msg, data)
}

// Making log with ERROR level.
func (l *Logger) Error(msg interface{}, data ...interface{}) {
	l.Log(ERROR, msg, data)
}

// Making log with CRITICAL level.
func (l *Logger) Fatal(msg interface{}, data ...interface{}) {
	l.Log(CRITICAL, msg, data)
	osExit(1)
}

// Making formatted log with DEBUG level.
func (l *Logger) Debugf(msg string, params ...interface{}) {
	l.Log(DEBUG, fmt.Sprintf(msg, params...), nil)
}

// Making formatted log with INFO level.
func (l *Logger) Infof(msg string, params ...interface{}) {
	l.Log(INFO, fmt.Sprintf(msg, params...), nil)
}

// Making formatted log with WARN level.
func (l *Logger) Warnf(msg string, params ...interface{}) {
	l.Log(WARNING, fmt.Sprintf(msg, params...), nil)
}

// Making formatted log with ERROR level.
func (l *Logger) Errorf(msg string, params ...interface{}) {
	l.Log(ERROR, fmt.Sprintf(msg, params...), nil)
}

// Making formatted log with CRITICAL level.
func (l *Logger) Fatalf(msg string, params ...interface{}) {
	l.Log(CRITICAL, fmt.Sprintf(msg, params...), nil)
	osExit(1)
}

// When you want to send logs to another appender,
// you should create instance of appender and call this method.
// Method is expecting appender instance to be passed
// to this method. At the end passed appender will receive logs
func (l *Logger) Enable(appender Appender) {
	l.appenders = append(l.appenders, appender)
}

// If you want to disable logs from some appender you can use this method.
// You have to call method either with appender instance,
// or you can pass appender Id as argument.
// If appender is found, it will be removed from list of appenders of this logger,
// and all other further logs won't be received by this appender.
func (l *Logger) Disable(target interface{}) {
	var id string
	var appender Appender

	switch object := target.(type) {
	case string:
		id = object
	case Appender:
		appender = object
	default:
		l.Warn("Error while disabling logger. Cannot cast to target type.")
		return
	}

	for i, app := range l.appenders {
		// if we can find the same appender reference
		// or we can extract and match id from appender
		// or we can match received id string argument with one of appender's id
		if (appender != nil && (app == appender || appender.Id() == app.Id())) || id == app.Id() {
			var toAppend []Appender

			if len(l.appenders) >= i+1 {
				toAppend = l.appenders[i+1:]
			}

			l.appenders = append(l.appenders[:i], toAppend...)
			return
		}
	}
}
