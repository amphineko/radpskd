package log

import (
	golog "log"
)

var (
	Debug   *golog.Logger = golog.New(golog.Writer(), "[DEBUG] ", golog.Flags())
	Error   *golog.Logger = golog.New(golog.Writer(), "[ERROR] ", golog.Flags())
	Info    *golog.Logger = golog.New(golog.Writer(), "[INFO] ", golog.Flags())
	Warning *golog.Logger = golog.New(golog.Writer(), "[WARN] ", golog.Flags())
)
