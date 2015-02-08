/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package heka_emitter

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/ildus/golog/id"
)

// TryClose closes w if w implements io.Closer.
func TryClose(w io.Writer) (err error) {
	if c, ok := w.(io.Closer); ok {
		err = c.Close()
	}
	return
}

// hekaMessagePool holds recycled Heka message objects and encoding buffers.
var hekaMessagePool = sync.Pool{New: func() interface{} {
	return &hekaMessage{
		header: new(Header),
		msg:    new(Message),
		buf:    proto.NewBuffer(nil),
	}
}}

func newHekaMessage() *hekaMessage {
	return hekaMessagePool.Get().(*hekaMessage)
}

type hekaMessage struct {
	header   *Header
	msg      *Message
	buf      *proto.Buffer
	outBytes []byte
}

func (hm *hekaMessage) free() {
	if cap(hm.outBytes) > 1024 {
		return
	}
	hm.buf.Reset()
	hm.outBytes = hm.outBytes[:0]
	hm.msg.Fields = nil
	hekaMessagePool.Put(hm)
}

func (hm *hekaMessage) marshalFrame() ([]byte, error) {
	msgSize := hm.msg.Size()
	if msgSize > MAX_MESSAGE_SIZE {
		return nil, fmt.Errorf("Message size %d exceeds maximum size %d",
			msgSize, MAX_MESSAGE_SIZE)
	}
	hm.header.SetMessageLength(uint32(msgSize))
	headerSize := hm.header.Size()
	if headerSize > MAX_HEADER_SIZE {
		return nil, fmt.Errorf("Header size %d exceeds maximum size %d",
			headerSize, MAX_HEADER_SIZE)
	}
	totalSize := headerSize + HEADER_FRAMING_SIZE + msgSize
	if cap(hm.outBytes) < totalSize {
		hm.outBytes = make([]byte, totalSize)
	} else {
		hm.outBytes = hm.outBytes[:totalSize]
	}
	hm.outBytes[0] = RECORD_SEPARATOR
	hm.outBytes[1] = byte(headerSize)
	hm.buf.SetBuf(hm.outBytes[HEADER_DELIMITER_SIZE:HEADER_DELIMITER_SIZE])
	if err := hm.buf.Marshal(hm.header); err != nil {
		return nil, err
	}
	hm.outBytes[headerSize+HEADER_DELIMITER_SIZE] = UNIT_SEPARATOR
	hm.buf.SetBuf(hm.outBytes[headerSize+HEADER_FRAMING_SIZE : headerSize+HEADER_FRAMING_SIZE])
	if err := hm.buf.Marshal(hm.msg); err != nil {
		return nil, err
	}
	return hm.outBytes, nil
}

// NewProtobufEmitter creates a Protobuf-encoded Heka log message emitter
func NewProtobufEmitter(writer io.Writer, envVersion,
	hostname, loggerName string) *ProtobufEmitter {

	return &ProtobufEmitter{
		Writer:       writer,
		LogName:      loggerName,
		Pid:          int32(os.Getpid()),
		EnvVersion:   envVersion,
		Hostname:     hostname,
		UseMockFuncs: false,
	}
}

// A ProtobufEmitter emits framed, Protobuf-encoded log messages.
type ProtobufEmitter struct {
	io.Writer
	LogName      string
	Pid          int32
	EnvVersion   string
	Hostname     string
	UseMockFuncs bool
}

// Emit encodes and sends a framed log message.
func (pe *ProtobufEmitter) Emit(level int32, messageType, payload string,
	fields map[string]string) (err error) {

	msgID, err := id.GenerateBytes()
	if err != nil {
		return fmt.Errorf("Error generating Protobuf log message ID: %s", err)
	}

	hm := newHekaMessage()
	defer hm.free()

	if pe.UseMockFuncs {
		hm.msg.SetID(idGenerateBytes())
		hm.msg.SetTimestamp(timeNow().UnixNano())
		hm.msg.SetPid(varosGetPid())
	} else {
		hm.msg.SetID(msgID)
		hm.msg.SetTimestamp(time.Now().UnixNano())
		hm.msg.SetPid(pe.Pid)
	}

	hm.msg.SetType(messageType)
	hm.msg.SetLogger(pe.LogName)
	hm.msg.SetSeverity(level)
	hm.msg.SetPayload(payload)
	hm.msg.SetEnvVersion(pe.EnvVersion)
	hm.msg.SetHostname(pe.Hostname)
	if fields != nil {
		for name, val := range fields {
			hm.msg.AddStringField(name, val)
		}
	}
	hm.msg.SortFields()

	outBytes, err := hm.marshalFrame()
	if err != nil {
		return fmt.Errorf("Error encoding Protobuf log message: %s", err)
	}
	if _, err = pe.Writer.Write(outBytes); err != nil {
		return fmt.Errorf("Error sending Protobuf log message: %s", err)
	}
	return nil
}

// Close closes the underlying write stream. Implements LogEmitter.Close.
func (pe *ProtobufEmitter) Close() error {
	return TryClose(pe.Writer)
}
