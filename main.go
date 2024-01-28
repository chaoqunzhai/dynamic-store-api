package main

import (
	"go-admin/cmd"
	"go-admin/common/global"
	_"go-admin/global/initialization"
)
var Version string
//go:generate swag init --parseDependency --parseDepth=6

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorizationl
func main() {
	global.Version = Version
	cmd.Execute()
}
