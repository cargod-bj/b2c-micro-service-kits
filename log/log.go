package log

// Fatal is equivalent to l.Critical(fmt.Sprint()) followed by a call to os.Exit(1).
func Fatal(args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Fatal(args...)
}

// Fatalf is equivalent to l.Critical followed by a call to os.Exit(1).
func Fatalf(format string, args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Fatalf(format, args...)
}

// Panic is equivalent to l.Critical(fmt.Sprint()) followed by a call to panic().
func Panic(args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Panic(args...)
}

// Panicf is equivalent to l.Critical followed by a call to panic().
func Panicf(format string, args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Panicf(format, args...)
}

// Critical logs a message using CRITICAL as log level.
func Critical(args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Critical(args...)
}

// Criticalf logs a message using CRITICAL as log level.
func Criticalf(format string, args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Criticalf(format, args...)
}

// Error logs a message using ERROR as log level.
func Error(args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Error(args...)
}

// Errorf logs a message using ERROR as log level.
func Errorf(format string, args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Errorf(format, args...)
}

// Warning logs a message using WARNING as log level.
func Warning(args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Warning(args...)
}

// Warningf logs a message using WARNING as log level.
func Warningf(format string, args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Warningf(format, args...)
}

// Notice logs a message using NOTICE as log level.
func Notice(args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Notice(args...)
}

// Noticef logs a message using NOTICE as log level.
func Noticef(format string, args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Noticef(format, args...)
}

// Info logs a message using INFO as log level.
func Info(args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Info(args...)
}

// Infof logs a message using INFO as log level.
func Infof(format string, args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Infof(format, args...)
}

// Debug logs a message using DEBUG as log level.
func Debug(args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Debug(args...)
}

// Debugf logs a message using DEBUG as log level.
func Debugf(format string, args ...interface{}) {
	if !hasInit {
		return
	}
	Logger.Debugf(format, args...)
}
