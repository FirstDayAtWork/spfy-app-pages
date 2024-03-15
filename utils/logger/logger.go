package logger

import "sync"

// Log Levels
const (
	InfoLevel = iota
	WarningLevel
	ErrorLevel
)

// Log level names
const (
	InfoLevelName    = "INFO"
	WarningLevelName = "WARNING"
	ErrorLevelName   = "ERROR"
)

// BaseLogger represents a logger that handles logging to a destination
// and is used by a higher level app logger.
type BaseLogger interface {
	Info(msg string)
	Warning(msg string)
	Error(msg string)
}

// logTask represents a contaier with information to be logged
type logTask struct {
	msg   string
	level int
}

// AppLogger represents a higher level logger which is meant to be invoked
// from external code (i.e. a web application)
type AppLogger struct {
	flags         int32
	level         int
	wg            sync.WaitGroup // Needed for graceful shutdown
	LogQueue      chan logTask
	ConsoleLogger BaseLogger
	FileLogger    BaseLogger
	HTTPLogger    BaseLogger
}

func GetNewAppLogger(level, queueSize int, flags int32) AppLogger {
	return AppLogger{
		flags:    flags,
		level:    level,
		LogQueue: make(chan logTask, queueSize),
		// TODO
		ConsoleLogger: nil,
		FileLogger:    nil,
		HTTPLogger:    nil,
	}

}

func (al *AppLogger) enqueueTask(task logTask) {
	al.LogQueue <- task
}

func (al *AppLogger) Info(msg string) {
	if al.level <= InfoLevel {
		al.enqueueTask(logTask{
			msg:   msg,
			level: InfoLevel,
		})
	}
}

func (al *AppLogger) Warning(msg string) {
	if al.level <= WarningLevel {
		al.enqueueTask(logTask{
			msg:   msg,
			level: WarningLevel,
		})
	}
}

func (al *AppLogger) Error(msg string) {
	if al.level <= ErrorLevel {
		al.enqueueTask(logTask{
			msg:   msg,
			level: ErrorLevel,
		})
	}
}

func (al *AppLogger) processInfoTask(task *logTask) {
	if al.ConsoleLogger != nil {
		al.ConsoleLogger.Info(task.msg)
	}
	if al.FileLogger != nil {
		al.FileLogger.Info(task.msg)
	}
}

func (al *AppLogger) processWarningTask(task *logTask) {
	if al.ConsoleLogger != nil {
		al.ConsoleLogger.Warning(task.msg)
	}
	if al.FileLogger != nil {
		al.FileLogger.Warning(task.msg)
	}
	if al.HTTPLogger != nil {
		al.HTTPLogger.Warning(task.msg)
	}
}

func (al *AppLogger) processErrorTask(task *logTask) {
	if al.ConsoleLogger != nil {
		al.ConsoleLogger.Error(task.msg)
	}
	if al.FileLogger != nil {
		al.FileLogger.Error(task.msg)
	}
	if al.HTTPLogger != nil {
		al.HTTPLogger.Error(task.msg)
	}
}

func (al *AppLogger) processTask(task *logTask) {
	switch task.level {
	case InfoLevel:
		al.processInfoTask(task)
	case WarningLevel:
		al.processWarningTask(task)
	case ErrorLevel:
		al.processErrorTask(task)
	}
}

func (al *AppLogger) StartConsuming() {
	al.wg.Add(1)
	go func() {
		defer al.wg.Done()
		for {
			task, ok := <-al.LogQueue
			if !ok {
				continue
			}
			al.processTask(&task)
		}
	}()
}

func (al *AppLogger) StopConsuming() {
	al.wg.Wait()
	close(al.LogQueue)
}

// TODO implement console logger and check whether things work!
// TODO Then implement a file logger
// TODO then implement http logger
// // ConsoleLogger handles logging to console / stdout
// type ConsoleLogger struct {
// 	level       *int // propagated from the app logger
// 	infoLogger  *log.Logger
// 	warnLogger  *log.Logger
// 	errorLogger *log.Logger
// }

// func (cl ConsoleLogger) Info(msg string) {
// 	if cl.level <=
// }
