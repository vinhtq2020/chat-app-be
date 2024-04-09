package logger

import "time"

type LoggingLevel string

const (
	ERROR = "ERROR"
	WARN  = "WARN"
	INFO  = "INFO"
	DEBUG = "DEBUG"
	TRACE = "TRACE"
)

type LogInfo struct {
	CreatedAt time.Time    `json:"createdAt"`
	Level     LoggingLevel `json:"level"`
	Location  string       `json:"location"`
	MsgId     *string      `json:"msgId"`
	Msg       string       `json:"msg"`
}
