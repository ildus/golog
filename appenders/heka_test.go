package appenders

import (
	"github.com/ivpusic/golog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
}

func TestHekaId(t *testing.T) {
	appender := Heka(golog.Conf{})
	assert.Equal(t, "github.com/ildus/golog/appender/heka", appender.Id())
}
