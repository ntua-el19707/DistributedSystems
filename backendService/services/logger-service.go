package  services 

import (
	"log"
	"fmt"
	) 

const loggerTeplate = "['%s'] : %s"
type LogerService interface {
	Service 
	Log(msg  string)
	Fatal(msg  string )
	Sprintf(msg string) string
} 

type Logger struct {
	ServiceName string
}

func (loger * Logger) construct() error {
	log.Printf(loggerTeplate , loger.ServiceName + "Logger" ,  "Service Logger has  succesfully  contructed")
	return nil
}
func (loger * Logger) Log(msg string) {
	log.Printf(loggerTeplate , loger.ServiceName ,  msg)

}

func (loger * Logger) Fatal(msg  string ) {
	log.Fatalf(loggerTeplate , loger.ServiceName ,  msg)

}
func (loger * Logger) Sprintf(msg string) string  {
	return fmt.Sprintf(loggerTeplate , loger.ServiceName ,  msg)

}


// -- Mocks  For  Testing -- 
type mockLogger struct {
	logs []string
	Faults []  string
}
func (m *mockLogger) Log(msg string) {
	m.logs = append(m.logs, msg)
}

func (m *mockLogger) Fatal(msg string) {
	m.Faults = append(m.Faults, msg)
}
func (m *mockLogger) Sprintf(msg string) string {
	return msg
}
func (m *mockLogger)  construct() error{
	return nil 
}
