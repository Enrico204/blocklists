package dnsbl

import (
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

// zapToClogCore is a zap Core that logs into CoreDNS clog system
type zapToClogCore struct {
	fields []zapcore.Field
}

func (zapToClogCore) Enabled(_ zapcore.Level) bool { return true }
func (n zapToClogCore) With(fields []zapcore.Field) zapcore.Core {
	clone := zapToClogCore{}
	clone.fields = append(clone.fields, n.fields...)
	clone.fields = append(clone.fields, fields...)
	return clone
}
func (n zapToClogCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(entry, n)
}
func (n zapToClogCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	var buf strings.Builder
	for _, field := range append(n.fields, fields...) {
		buf.WriteString(field.Key)
		buf.WriteRune('=')
		buf.WriteString(field.String)
		buf.WriteRune('\t')
	}
	buf.WriteString(entry.Message)

	var msg = buf.String()
	switch entry.Level {
	case zapcore.DPanicLevel:
		fallthrough
	case zapcore.PanicLevel:
		fallthrough
	case zapcore.FatalLevel:
		clog.Fatal(msg)
	case zapcore.ErrorLevel:
		clog.Error(msg)
	case zapcore.WarnLevel:
		clog.Warning(msg)
	case zapcore.InfoLevel:
		clog.Info(msg)
	case zapcore.DebugLevel:
		fallthrough
	default:
		clog.Debug(msg)
	}
	return nil
}
func (zapToClogCore) Sync() error { return nil }

// zapToClog returns a zap logger instance that logs into the CoreDNS clog
func zapToClog() *zap.SugaredLogger {
	zlogger := zap.New(zapToClogCore{})
	return zlogger.Sugar()
}
