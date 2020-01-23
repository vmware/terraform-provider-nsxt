/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package log

import (
	"log"
)

var Log Logger = NewDefaultLogger()

type Logger interface {
	Error(args ...interface{})
	Errorf(string, ...interface{})
	Info(args ...interface{})
	Infof(string, ...interface{})
	Debug(args ...interface{})
	Debugf(string, ...interface{})
}

type DefaultLogger struct {
}

func SetLogger(logger Logger) {
	Log = logger
}

func NewDefaultLogger() DefaultLogger {
	return DefaultLogger{}
}

func (d DefaultLogger) Error(args ...interface{}) {
	log.Println(args...)

}

func (d DefaultLogger) Errorf(a string, args ...interface{}) {
	log.Printf(a, args...)
}

func (d DefaultLogger) Info(args ...interface{}) {
	log.Print(args...)
}

func (d DefaultLogger) Infof(a string, args ...interface{}) {
	log.Printf(a, args...)
}

func (d DefaultLogger) Debug(args ...interface{}) {
	log.Print(args...)
}

func (d DefaultLogger) Debugf(a string, args ...interface{}) {
	log.Printf(a, args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Errorf(a string, args ...interface{}) {
	Log.Errorf(a, args...)
}

func Debug(args ...interface{}) {
	Log.Debug(args...)
}

func Debugf(a string, args ...interface{}) {
	Log.Debugf(a, args...)
}

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Infof(a string, args ...interface{}) {
	Log.Infof(a, args...)
}
