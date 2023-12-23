package  services 

import (
	"log"
	"fmt"
	) 

const loggerTeplate = "[info]['%s'] : %s\n"
const loggerErrorTeplate = "[error]['%s'] : %s\n"
type LogerService interface {
	Service  // logger is  a service 
	Log(msg  string) // logs  the  message  in stdout
    Error(msg  string) // logs error  message  in  stdout 
	Fatal(msg  string ) // log Fatal  meassge 
	Sprintf(msg string) string // return  the message that would be logged
} 

type Logger struct {
	ServiceName string
}
/**
   construct  -  create  the  instance  of  logger 
   @Returns error
*/
func (loger * Logger) construct() error {
	log.Printf(loggerTeplate , loger.ServiceName + "Logger" ,  "Service Logger has  succesfully  contructed")
	return nil
}
/**
   Log  -  log a  message  to stdout 
   @Returns  void 
*/
func (loger * Logger) Log(msg string) {
	log.Printf(loggerTeplate , loger.ServiceName ,  msg)
}
/**
   Error  -  log a  message  to stdout 
   @Returns  void 
*/
func (loger * Logger) Error(msg string) {
	log.Printf(loggerErrorTeplate , loger.ServiceName ,  msg)
}
/**
   Fatal  -  log Fatal(kill  process ) a  message  to stdout 
   @Returns  void 
*/
func (loger * Logger) Fatal(msg  string ) {
	log.Fatalf(loggerTeplate , loger.ServiceName ,  msg)

}

/**
   Sprintf  -  sprintf  the  message  
   @Returns  string
*/
func (loger * Logger) Sprintf(msg string) string  {
	return fmt.Sprintf(loggerTeplate , loger.ServiceName ,  msg)
}


// -- Mocks  For  Testing -- 
type mockLogger struct {
	logs []string
	Faults []  string
    errors []  string
}
// mock Log 
func (m *mockLogger) Log(msg string) {
	m.logs = append(m.logs, msg)
}

// mock Error
func (m *mockLogger) Error(msg string) {
	m.errors = append(m.errors, msg)
}
// mock error 
func (m *mockLogger) Fatal(msg string) {
	m.Faults = append(m.Faults, msg)
}
func (m *mockLogger) Sprintf(msg string) string {
	return msg
}
func (m *mockLogger)  construct() error{
	return nil 
}
