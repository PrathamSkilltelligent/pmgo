package logger

import (
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// prependEncoder is a custom encoder that adds a prefix before each log entry
type prependEncoder struct {
	zapcore.Encoder
	pool buffer.Pool
	tag  string
}

func (e *prependEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// New log buffer
	buf := e.pool.Get()

	// Prepend a custom tag based on the entry level
	if e.tag != "" {
		buf.AppendString(e.tag + ": ")
	}

	// append @cee prefix for rsyslog
	buf.AppendString("@cee: ")

	// Calling the embedded encoder's EncodeEntry to keep the original encoding format
	consolebuf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	// Write the output into your own buffer
	_, err = buf.Write(consolebuf.Bytes())
	if err != nil {
		return nil, err
	}
	return buf, nil
}
