package main

import "github.com/flacatus/che-inspector/pkg/cmd"

func main() {
	cmd := cmd.NewCheInspectorCobraCommand()
	_ = cmd.Execute()
}
