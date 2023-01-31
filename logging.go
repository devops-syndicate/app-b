package main

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

const TRACE_ID_KEY = "trace_id"
const SPAN_ID_KEY = "span_id"
const TIMESTAMP_KEY = "timestamp"
const MESSAGE_KEY = "message"

func initLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: TIMESTAMP_KEY,
			logrus.FieldKeyMsg:  MESSAGE_KEY,
		}})
	logrus.AddHook(&EnrichWithTracingHook{})
	logrus.SetLevel(logrus.InfoLevel)
}

type EnrichWithTracingHook struct {
}

func (t *EnrichWithTracingHook) Fire(e *logrus.Entry) error {
	span := trace.SpanFromContext(e.Context)
	if span.SpanContext().HasTraceID() {
		if _, ok := e.Data[TRACE_ID_KEY]; !ok {
			e.Data[TRACE_ID_KEY] = span.SpanContext().TraceID()
		}
	}
	if span.SpanContext().HasSpanID() {
		if _, ok := e.Data[SPAN_ID_KEY]; !ok {
			e.Data[SPAN_ID_KEY] = span.SpanContext().SpanID()
		}
	}
	return nil
}

func (t *EnrichWithTracingHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
