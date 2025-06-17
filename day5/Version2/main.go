package main

import "fmt"

type Logger interface {
	Log(message string)
}

type ConsoleLogger struct{}
type FileLogger struct{}
type RemoteLogger struct{}

func (logger *ConsoleLogger) Log(message string) {
	fmt.Println("The message from the ConsoleLogger is: ", message)
}
func (logger *FileLogger) Log(message string) {
	fmt.Println("The message from the FileLogger is: ", message)
}
func (logger *RemoteLogger) Log(message string) {
	fmt.Println("The message from the RemoteLogger is: ", message)
}

func logAll(loggers []Logger, message string) {

	for _, item := range loggers {
		item.Log(message)
	}

}

func main() {

	cl := ConsoleLogger{}
	fl := FileLogger{}
	rl := RemoteLogger{}

	slice := make([]Logger, 3)
	slice[0] = &cl
	slice[1] = &fl
	slice[2] = &rl

	logAll(slice, "Not gonna lie, it was playing with my sanity.")

}
