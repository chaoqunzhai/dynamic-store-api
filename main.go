package main

import (
	"go-admin/cmd"
)

//go:generate swag init --parseDependency --parseDepth=6

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
