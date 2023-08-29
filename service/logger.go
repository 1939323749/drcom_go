package service

import "log"

type Logger struct {
	loggers map[string]*log.Logger
}

type LoggerInterface interface {
	Info(string)
	Error(string)
}

func (l *Logger) Info(msg string) {
	l.loggers["info"].Println(msg)
}

func (l *Logger) Error(msg string) {
	l.loggers["error"].Println(msg)
}
