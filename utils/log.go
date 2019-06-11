package utils

import (
	"flag"
	"strconv"

	"github.com/golang/glog"
)

const (
	lfatal = iota
	lerror
	lwarning
	linfo
	ldebug
)

func SetLogger(debugLog bool) {
	var level glog.Level = ldebug
	if debugLog == false {
		level = lerror
	}
	flag.Set("logtostderr", "true")
	flag.Set("v", strconv.Itoa(int(level)))
	flag.Parse()
}

func LogFlush() {
	glog.Flush()
}

func InfoLog(args ...interface{}) {
	glog.V(linfo).Info(args...)
}

func InfoLogf(format string, args ...interface{}) {
	glog.V(linfo).Infof(format, args...)
	glog.Flush()
}

func WarningLog(args ...interface{}) {
	glog.V(lwarning).Info(args...)
}

func WarningLogf(format string, args ...interface{}) {
	glog.V(lwarning).Infof(format, args...)
}

func ErrorLog(args ...interface{}) {
	glog.V(lerror).Info(args...)
}

func ErrorLogf(format string, args ...interface{}) {
	glog.V(lerror).Infof(format, args...)
}

func FatalLog(args ...interface{}) {
	glog.V(lfatal).Info(args...)
}

func FataLogf(format string, args ...interface{}) {
	glog.V(lfatal).Infof(format, args...)
}
