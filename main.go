package main

import (
	"gateway/common"
	"gateway/teleport"
)

func main() {
	go teleport.Run()

	common.HoldOn()
}
