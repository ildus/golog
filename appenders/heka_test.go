package appenders

import (
	"./heka_appender"
	"bytes"
	"github.com/ildus/golog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHekaId(t *testing.T) {
	appender := Heka(golog.Conf{})
	assert.Equal(t, "github.com/ildus/golog/appenders/heka", appender.Id())
}

func TestProtobufEmitter(t *testing.T) {
	buf := new(bytes.Buffer)
	pe := heka_appender.NewProtobufEmitter(buf, "2", "example.com", "test-json-emitter")
	pe.UseMockFuncs = true
	expected := []byte{
		0x1e, 0x2,
		// Header.
		0x8, 0x75,
		0x1f,
		// Message.
		0xa, 0x10, 0xd1, 0xc7, 0xc7, 0x68, 0xb1, 0xbe, 0x4c, 0x70, 0x93, 0xa6, 0x9b,
		0x52, 0x91, 0xd, 0x4b, 0xaa, 0x10, 0x80, 0xc0, 0xe1, 0xd8, 0xda, 0xfd, 0xbb,
		0xba, 0x11, 0x1a, 0x4, 0x74, 0x65, 0x73, 0x74, 0x22, 0x11, 0x74, 0x65, 0x73,
		0x74, 0x2d, 0x6a, 0x73, 0x6f, 0x6e, 0x2d, 0x65, 0x6d, 0x69, 0x74, 0x74, 0x65,
		0x72, 0x28, 0x6, 0x32, 0x5, 0x48, 0x6f, 0x77, 0x64, 0x79, 0x3a, 0x1, 0x32,
		0x40, 0xd2, 0x9, 0x4a, 0xb, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e,
		0x63, 0x6f, 0x6d, 0x52, 0xa, 0xa, 0x1, 0x61, 0x10, 0x0, 0x1a, 0x0, 0x22,
		0x1, 0x62, 0x52, 0xa, 0xa, 0x1, 0x63, 0x10, 0x0, 0x1a, 0x0, 0x22, 0x1, 0x64,
		0x52, 0xa, 0xa, 0x1, 0x65, 0x10, 0x0, 0x1a, 0x0, 0x22, 0x1, 0x66,
	}
	err := pe.Emit(int32(golog.INFO), "test", "Howdy",
		heka_appender.LogFields{"c": "d", "a": "b", "e": "f"})
	if err != nil {
		t.Errorf("Error marshaling framed log message: %s", err)
	}
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Errorf("Malformed framed log message: got %#v; want %#v",
			buf.Bytes(), expected)
	}
}
