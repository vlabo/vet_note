package view

const (
	OpenSelectView int = iota
	OpenPatientView
)

var ViewEventChan chan int = make(chan int, 10)
