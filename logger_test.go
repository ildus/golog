package golog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testAppender struct {
	count      int
	errorCount int
	warnCount  int
	msg        string
}

func (s *testAppender) Append(log Log) {
	s.msg = log.Message
	s.count += 1

	if log.Level == WARNING {
		s.warnCount += 1
	}

	if log.Level < WARNING {
		s.errorCount += 1
	}
}

func cleanupTest() {
	useStdFuncs()
	loggers = map[string]*Logger{}

	Default = &Logger{
		Name:  "default",
		Level: DEBUG,
	}

	Default.Enable(StdoutAppender())
}

func (s *testAppender) Id() string {
	return "github.com/ildus/golog/test"
}

// func TestDefaultPresent(t *testing.T) {
// 	assert.NotNil(t, Default, "default logged should be authomatically created")
// }

func TestEnable(t *testing.T) {
	defer cleanupTest()
	ta := &testAppender{}

	oldcount := len(Default.appenders)
	Default.Enable(ta)
	assert.True(t, oldcount == len(Default.appenders)-1)

	Default.Info("testenable: some msg")
	assert.Exactly(t, 1, ta.count)

	Default.Info("testenable: some msg")
	assert.Exactly(t, 2, ta.count)
}

func TestDisableByInstance(t *testing.T) {
	defer cleanupTest()
	ta := &testAppender{}

	oldcount := len(Default.appenders)
	Default.Enable(ta)
	assert.True(t, oldcount == len(Default.appenders)-1)

	Default.Disable(ta)
	assert.True(t, oldcount == len(Default.appenders))

	Default.Info("some msg", "err1", 423)
	assert.Exactly(t, 0, ta.count)

	Default.Info("some msg")
	assert.Exactly(t, 0, ta.count)

	Default.Enable(ta)
	Default.Info("some msg")
	assert.Exactly(t, 1, ta.count)
}

func TestDisableById(t *testing.T) {
	defer cleanupTest()
	ta := &testAppender{}

	oldcount := len(Default.appenders)
	Default.Enable(ta)
	assert.True(t, oldcount == len(Default.appenders)-1)

	Default.Disable(ta.Id())
	assert.True(t, oldcount == len(Default.appenders))

	Default.Info("some msg")
	assert.Exactly(t, 0, ta.count)

	Default.Info("some msg")
	assert.Exactly(t, 0, ta.count)
}

func TestDisableInvalid(t *testing.T) {
	defer cleanupTest()

	ta := &testAppender{}
	Default.Enable(ta)

	Default.Disable(123)
	assert.Exactly(t, 1, ta.warnCount)
}

func TestLogCalls(t *testing.T) {
	mockFuncs()
	defer cleanupTest()

	defer func() {
		if r := recover(); r != nil {
		}
	}()

	ta := &testAppender{}
	Default.Enable(ta)

	Default.Debug("some msg")
	Default.Debug("some msg")

	Default.Info("some msg")
	Default.Info("some msg")

	Default.Warn("some msg")
	Default.Warn("some msg")

	Default.Error("some msg")
	Default.Error("some msg")

	Default.Fatal("some msg")
	Default.Fatal("some msg")

	assert.Exactly(t, 10, ta.count)

	ta.msg = ""
	Default.Debugf("some %s", "message")
	assert.Equal(t, "some message", ta.msg)

	ta.msg = ""
	Default.Infof("some %d %s", 3, "message")
	assert.Equal(t, "some 3 message", ta.msg)

	ta.msg = ""
	Default.Warnf("some %s", "message")
	assert.Equal(t, "some message", ta.msg)

	ta.msg = ""
	Default.Errorf("some %d %s", 3, "message")
	assert.Equal(t, "some 3 message", ta.msg)

	ta.msg = ""
	Default.Errorf("some %s message", "panic")
	assert.Equal(t, "some panic message", ta.msg)
}

func TestLogCallsWithLevel(t *testing.T) {
	defer cleanupTest()

	ta := &testAppender{}
	Default.Enable(ta)

	Default.Level = WARNING

	Default.Debug("some msg")
	Default.Debug("some msg")

	Default.Info("some msg")
	Default.Info("some msg")

	Default.Warn("some msg")
	Default.Warn("some msg")

	Default.Error("some msg")
	Default.Error("some msg")

	assert.Exactly(t, 4, ta.count)

	ta.count = 0
	Default.Level = DEBUG

	Default.Debug("some msg")
	Default.Debug("some msg")

	Default.Info("some msg")
	Default.Info("some msg")

	Default.Warn("some msg")
	Default.Warn("some msg")

	Default.Error("some msg")
	Default.Error("some msg")

	assert.Exactly(t, 8, ta.count)
}

func normalizeNameLenInTest(name string) string {
	length := len(name)
	missing := curnamelen - length

	for i := 0; i < missing; i++ {
		name += " "
	}

	return name
}

func TestNormalizeName(t *testing.T) {
	// name is too long
	l := GetLogger("s.o.m.e.r.e.a.l.l.y.l.o.n.g.n.a.m.e.t.e.s.t.n.a.m.e.")
	l.Debug(l.Name)
	assert.Equal(t, normalizeNameLenInTest("s.o.m.e.r.e.a.l.l.y."), l.Name)

	l = GetLogger("github.com/ildus/golog")
	l.Debug(l.Name)
	assert.Equal(t, normalizeNameLenInTest("git/ild/gol"), l.Name)

	l = GetLogger("github.com.ildus.golog")
	l.Debug(l.Name)
	assert.Equal(t, normalizeNameLenInTest("git.com.ild.gol"), l.Name)

	// name is too short
	l = GetLogger("main")
	l.Debug(l.Name)
	assert.Equal(t, normalizeNameLenInTest("main"), l.Name)

	// name is correct
	rightName := ""
	for i := 0; i < curnamelen; i++ {
		rightName += "a"
	}

	l = GetLogger(rightName)
	l.Debug(l.Name)
	assert.Equal(t, rightName, l.Name)
}
