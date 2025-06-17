package main

import "fmt"

// Logger is an interface that defines a logging method.
type Logger interface {
	Log(message string)
}

// ConsoleLogger implements Logger for console output.
type ConsoleLogger struct{}

// FileLogger implements Logger for file output (simulated).
type FileLogger struct{}

// RemoteLogger implements Logger for remote output (simulated).
type RemoteLogger struct{}

// Log prints a message to the console for ConsoleLogger.
func (logger *ConsoleLogger) Log(message string) {
	fmt.Println("The message from the ConsoleLogger is: ", message)
}

// Log prints a message to the console for FileLogger (simulating file logging).
func (logger *FileLogger) Log(message string) {
	fmt.Println("The message from the FileLogger is: ", message)
}

// Log prints a message to the console for RemoteLogger (simulating remote logging).
func (logger *RemoteLogger) Log(message string) {
	fmt.Println("The message from the RemoteLogger is: ", message)
}

// logAll logs a message using all provided Logger implementations.
func logAll(loggers []Logger, message string) {
	for _, item := range loggers {
		item.Log(message)
	}
}

func main() {
	// Create instances of each logger type.
	cl := ConsoleLogger{}
	fl := FileLogger{}
	rl := RemoteLogger{}

	// Create a slice of Logger interfaces and assign logger instances.
	slice := make([]Logger, 3)
	slice[0] = &cl
	slice[1] = &fl
	slice[2] = &rl

	// Log a message using all loggers.
	logAll(slice, "Not gonna lie, it was playing with my sanity.")
}
