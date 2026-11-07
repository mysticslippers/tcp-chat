package main

type serverMessageMsg string
type messageSentMsg struct{}
type networkErrorMsg struct {
	err error
}
