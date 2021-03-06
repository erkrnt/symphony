package main

import (
	"github.com/erkrnt/symphony/internal/manager"
	"github.com/sirupsen/logrus"
)

func main() {
	m, err := manager.New()

	if err != nil {
		logrus.Fatal(err)
	}

	go m.Start()

	errorC := <-m.ErrorC

	logrus.Fatal(errorC)
}
