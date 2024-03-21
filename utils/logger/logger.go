package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

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

// callerDepth controls the number of stack frames to ascend
// when calling runtime.Caller() to retrieve info on file and line number
const callerDepth = 2

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
	flags         int
	level         int
	wg            sync.WaitGroup
	stopCh        chan struct{}
	LogQueue      chan logTask
	ConsoleLogger BaseLogger
	FileLogger    BaseLogger
	HTTPLogger    BaseLogger
}

func GetNewAppLogger(
	level, queueSize, flags int,
	folder, file string,
) AppLogger {
	return AppLogger{
		flags:    flags,
		level:    level,
		LogQueue: make(chan logTask, queueSize),
		wg:       sync.WaitGroup{},
		stopCh:   make(chan struct{}, 1),
		// TODO
		ConsoleLogger: GetNewConsoleLogger(level, flags),
		FileLogger:    GetNewFileLogger(level, flags, file, folder),
		HTTPLogger:    nil,
	}
}

func (al *AppLogger) addCallerInfo(msg string) string {
	_, file, line, ok := runtime.Caller(callerDepth)
	if !ok {
		fmt.Println("error enriching message, logging info is not full")
		return msg
	}
	return fmt.Sprintf("%s:%d | %s", path.Base(file), line, msg)
}

func (al *AppLogger) enqueueTask(task logTask) {
	al.LogQueue <- task
}

func (al *AppLogger) Info(msg string) {
	if al.level <= InfoLevel {
		al.enqueueTask(logTask{
			msg:   al.addCallerInfo(msg),
			level: InfoLevel,
		})
	}
}

func (al *AppLogger) Warning(msg string) {
	if al.level <= WarningLevel {
		al.enqueueTask(logTask{
			msg:   al.addCallerInfo(msg),
			level: WarningLevel,
		})
	}
}

func (al *AppLogger) Error(msg string) {
	if al.level <= ErrorLevel {
		al.enqueueTask(logTask{
			msg:   al.addCallerInfo(msg),
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
			if ok {
				al.processTask(&task)
				continue
			}
			// Check for stop signal if the queue is empty
			_, ok = <-al.stopCh
			if ok {
				fmt.Println("Got a stop signal")
				return
			}
		}
	}()
}

func (al *AppLogger) StopConsuming() {
	// This does not cause a deadlock, but not all logs go through
	// Read on waitgroups more!
	close(al.LogQueue)
	al.stopCh <- struct{}{}
	close(al.stopCh)
	al.wg.Wait()

}

// TODO implement console logger and check whether things work!
// TODO then implement http logger
// ConsoleLogger handles logging to console / stdout
type ConsoleLogger struct {
	level       int // propagated from the app logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

func GetNewConsoleLogger(level, flags int) *ConsoleLogger {
	return &ConsoleLogger{
		level:       level,
		infoLogger:  log.New(os.Stdout, fmt.Sprintf("%s | ", InfoLevelName), flags),
		warnLogger:  log.New(os.Stdout, fmt.Sprintf("%s | ", WarningLevelName), flags),
		errorLogger: log.New(os.Stdout, fmt.Sprintf("%s | ", ErrorLevelName), flags),
	}
}

func (cl *ConsoleLogger) Info(msg string) {
	if cl.level <= InfoLevel {
		cl.infoLogger.Println(msg)
	}
}

func (cl *ConsoleLogger) Warning(msg string) {
	if cl.level <= WarningLevel {
		cl.warnLogger.Println(msg)
	}
}

func (cl *ConsoleLogger) Error(msg string) {
	if cl.level <= ErrorLevel {
		cl.errorLogger.Println(msg)
	}
}

// FileLogger handles logging to .txt files
type FileLogger struct {
	level           int
	file            string
	folder          string
	infoLogger      *log.Logger
	warnLogger      *log.Logger
	errorLogger     *log.Logger
	destinationFile *os.File
}

func GetNewFileLogger(level, flags int, file, folder string) *FileLogger {
	fl := &FileLogger{
		level:  level,
		file:   file,
		folder: folder,
	}
	err := fl.updateDestination()
	if err != nil {
		panic(err)
	}
	fl.infoLogger = log.New(fl.destinationFile, fmt.Sprintf("%s | ", InfoLevelName), flags)
	fl.warnLogger = log.New(fl.destinationFile, fmt.Sprintf("%s | ", WarningLevelName), flags)
	fl.errorLogger = log.New(fl.destinationFile, fmt.Sprintf("%s | ", ErrorLevelName), flags)
	return fl
}

func GetLogFileName(file string) string {
	currDate := time.Now().UTC()
	return fmt.Sprintf(
		"%s.%s",
		file,
		fmt.Sprintf("%d_%d_%d.log", currDate.Year(), int(currDate.Month()), currDate.Day()),
	)
}

// CreateLogFolder creates a folder for log files if it does not exist
func CreateLogFolder(folder string) error {
	info, err := os.Stat(folder)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return os.MkdirAll(folder, os.ModePerm)
		}
		return err
	}
	if info.IsDir() {
		return nil
	}
	return fmt.Errorf("path %s already exists and it's not a dir", folder)
}

func (fl *FileLogger) updateDestination() error {
	newFile := GetLogFileName(fl.file)
	if fl.destinationFile != nil && path.Base(fl.destinationFile.Name()) == newFile {
		// Date has not yet changed, no updates to destination are needed
		return nil
	}
	// Date has changed so we need a new log file
	if err := CreateLogFolder(fl.folder); err != nil {
		return err
	}
	// At this point we are sure the folder exists
	logFile, err := os.OpenFile(path.Join(fl.folder, newFile), os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	fl.destinationFile = logFile
	return nil
}

func (fl *FileLogger) Info(msg string) {
	if err := fl.updateDestination(); err != nil {
		fmt.Printf("logging error: updating file destination failed. Details: %s\n", err)
	}
	if fl.level <= InfoLevel {
		fl.infoLogger.Println(msg)
	}
}

func (fl *FileLogger) Warning(msg string) {
	if err := fl.updateDestination(); err != nil {
		fmt.Printf("logging error: updating file destination failed. Details: %s\n", err)
	}
	if fl.level <= WarningLevel {
		fl.warnLogger.Println(msg)
	}
}

func (fl *FileLogger) Error(msg string) {
	if err := fl.updateDestination(); err != nil {
		fmt.Printf("logging error: updating file destination failed. Details: %s\n", err)
	}
	if fl.level <= ErrorLevel {
		fl.errorLogger.Println(msg)
	}
}
