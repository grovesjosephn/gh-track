package main

import "hab/cmd"

// Version is set by ldflags during build
var version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}