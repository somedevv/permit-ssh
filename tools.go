package main

import (
	"os"
	"os/exec"
)

func CallClear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
