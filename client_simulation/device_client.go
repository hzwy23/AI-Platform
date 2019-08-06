package main

import (
	"ai-platform/client_simulation/test"
	"time"
)

func main() {
	for {
		test.WriteBroadcast()
		time.Sleep(time.Second * 3)
	}

}