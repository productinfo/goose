package logger

import (
	"bytes"
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func buildLogger() (Logger, *bytes.Buffer) {
	buf := bytes.NewBuffer(nil)
	logrusLogger := logrus.New()
	logrusLogger.Out = buf
	logrusLogger.Formatter = &logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	}
	entry := logrus.NewEntry(logrusLogger)

	logger := func(ctx Valuer, err error) *logrus.Entry {
		return ContextLog(ctx, err, entry)
	}

	return logger, buf
}

func TestContextLog(t *testing.T) {
	logger := New("foo")

	origGlobal := logrus.Fields{}
	for k, v := range GlobalFields {
		origGlobal[k] = v
	}

	GlobalFields["testKey"] = "value"

	ctx := context.Background()
	ctx = WithField(ctx, "bar", "baz")
	entry := logger(ctx, nil).WithField("a", "b")
	assert.Equal(t, logrus.Fields{
		"component": "foo",
		"bar":       "baz",
		"a":         "b",
		"testKey":   "value",
	}, entry.Data)

	// Restore
	GlobalFields = origGlobal
}

func TestLogIfError(t *testing.T) {
	ctx := context.Background()
	{
		logger, buf := buildLogger()
		fn := func() error { return nil }
		LogIfError(ctx, fn, logger, "")
		assert.Equal(t, "", buf.String())
	}
	{
		logger, buf := buildLogger()
		fn := func() error { return errors.New("foo") }
		LogIfError(ctx, fn, logger, "msg")
		assert.Equal(t, "level=error msg=msg error=foo\n", buf.String())
	}
	{
		logger, buf := buildLogger()
		fn := func() error { return errors.New("foo") }
		ctx := WithField(ctx, "test", "bar")
		LogIfError(ctx, fn, logger, "msg")
		assert.Equal(t, "level=error msg=msg error=foo test=bar\n", buf.String())
	}
	{
		logger, buf := buildLogger()
		nestedFn := func() error { return errors.New("root_cause") }
		fn := func() error {
			err := nestedFn()
			return errors.Wrap(err, "err_msg")
		}
		LogIfError(ctx, fn, logger, "msg")
		assert.Equal(t, "level=error msg=msg cause=root_cause error=\"err_msg: root_cause\"\n", buf.String())
	}
}
