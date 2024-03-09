package Logger

import (
	"Service"
	"fmt"
	"log"
)

const loggerTeplate = "[info]['%s'] : %s\n"
const loggerErrorTeplate = "[error]['%s'] : %s\n"
const fatalTeplate = "[Fatal]['%s'] : %s\n"

type LoggerService interface {
	Service.Service                 // logger is  a service
	Log(msg string)                 // logs  the  message  in stdout
	Error(msg string)               // logs error  message  in  stdout
	Fatal(msg string)               // log Fatal  meassge
	Sprintf(msg string) string      // return  the message that would be logged
	SprintErrorf(msg string) string // return error  string
}

type Logger struct {
	ServiceName string
}

/*
*

	Construct  -  create  the  instance  of  logger
	@Returns error
*/
func (l *Logger) Construct() error {
	log.Printf(loggerTeplate, l.ServiceName+"Logger", "Service Logger has  succesfully  contructed")
	return nil
}

/*
*

	Log  -  log a  message  to stdout
	@Returns  void
*/
func (l Logger) Log(msg string) {
	log.Printf(loggerTeplate, l.ServiceName, msg)
}

/*
*

	Error  -  log a  message  to stdout
	@Returns  void
*/
func (l Logger) Error(msg string) {
	log.Printf(loggerErrorTeplate, l.ServiceName, msg)
}

/*
*

	Fatal  -  log Fatal(kill  process ) a  message  to stdout
	@Returns  void
*/
func (l Logger) Fatal(msg string) {
	log.Fatalf(fatalTeplate, l.ServiceName, msg)

}

/*
*

	Sprintf  -  sprintf  the  message
	@Returns  string
*/
func (l Logger) Sprintf(msg string) string {
	return fmt.Sprintf(loggerTeplate, l.ServiceName, msg)
}

/*
*

	SprintErrorf  sprintf  the  message with error  template
	@Returns  string
*/
func (l Logger) SprintErrorf(msg string) string {
	return fmt.Sprintf(loggerErrorTeplate, l.ServiceName, msg)
}

// -- Mocks  For  Testing --
type MockLogger struct {
	Logs      []string
	Faults    []string
	ErrorList []string
}

// mock Log
func (m *MockLogger) Log(msg string) {
	m.Logs = append(m.Logs, msg)
}

// mock Error
func (m *MockLogger) Error(msg string) {
	m.ErrorList = append(m.ErrorList, msg)
}

// mock error
func (m *MockLogger) Fatal(msg string) {
	m.Faults = append(m.Faults, msg)
}

func (m *MockLogger) Sprintf(msg string) string {
	return msg
}

func (m *MockLogger) SprintErrorf(msg string) string {
	return msg
}
func (m *MockLogger) Construct() error {
	return nil
}
