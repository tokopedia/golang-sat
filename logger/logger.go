package logger

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type Config struct {
	Logger  *log.Logger
	IsDebug bool
}

func (c *Config) GetLogger() LoggerSource {
	return &Logger{
		logger:  c.Logger,
		isDebug: c.IsDebug,
	}
}

// A LoggerSource is anything do logging http request & response
type LoggerSource interface {
	LogRequest(req *http.Request)
	LogResponse(resp *http.Response)
}

type Logger struct {
	logger  *log.Logger
	isDebug bool
}

func (l *Logger) LogRequest(req *http.Request) {
	if !l.isDebug {
		return
	}

	logReq, err := httputil.DumpRequest(req, true)
	if err == nil {
		l.logger.Println(string(logReq))
	}
}

func (l *Logger) LogResponse(resp *http.Response) {
	if !l.isDebug {
		return
	}

	logResp, err := httputil.DumpResponse(resp, true)
	if err == nil {
		l.logger.Println(string(logResp))
	}
}
