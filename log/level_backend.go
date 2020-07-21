package log

import (
	"github.com/op/go-logging"
	"sync"
)

type moduleLeveled struct {
	levels    map[string]logging.Level
	backend   logging.Backend
	formatter logging.Formatter
	once      sync.Once
}

// log深度，可放在log的可变长参数中
type CallDepth struct {
	Depth int
}

// AddModuleLevel wraps a log backend with knobs to have different log levels
// for different modules.
func AddModuleLevel(backend logging.Backend) logging.LeveledBackend {
	var leveled logging.LeveledBackend
	leveled = &moduleLeveled{
		levels:  make(map[string]logging.Level),
		backend: backend,
	}
	return leveled
}

// GetLevel returns the log level for the given module.
func (l *moduleLeveled) GetLevel(module string) logging.Level {
	level, exists := l.levels[module]
	if exists == false {
		level, exists = l.levels[""]
		// no configuration exists, default to debug
		if exists == false {
			level = logging.DEBUG
		}
	}
	return level
}

// SetLevel sets the log level for the given module.
func (l *moduleLeveled) SetLevel(level logging.Level, module string) {
	l.levels[module] = level
}

// IsEnabledFor will return true if logging is enabled for the given module.
func (l *moduleLeveled) IsEnabledFor(level logging.Level, module string) bool {
	return level <= l.GetLevel(module)
}

func (l *moduleLeveled) Log(level logging.Level, calldepth int, rec *logging.Record) (err error) {
	if l.IsEnabledFor(level, rec.Module) {
		// TODO get rid of traces of formatter here. BackendFormatter should be used.
		//rec.formatter = l.getFormatterAndCacheCurrent()
		if rec.Args != nil {
			target := -1
			for i := range rec.Args {
				if d, ok := rec.Args[i].(CallDepth); ok {
					calldepth += d.Depth
					target = i
					break
				}
			}
			// 如果有此对象，则将此对象从arg中删除
			if target >= 0 {
				if target == 0 {
					if len(rec.Args) == 1 {
						rec.Args = nil
					} else {
						rec.Args = rec.Args[target+1:]
					}
				} else if target == (len(rec.Args) - 1) {
					rec.Args = rec.Args[:target]
				} else {
					rec.Args = append(rec.Args[:target], rec.Args[target+1:]...)
				}
			}
		}
		err = l.backend.Log(level, calldepth+1, rec)
	}
	return
}

//func (l *moduleLeveled) getFormatterAndCacheCurrent() logging.Formatter {
//	l.once.Do(func() {
//		if l.formatter == nil {
//			l.formatter = getFormatter()
//		}
//	})
//	return l.formatter
//}
