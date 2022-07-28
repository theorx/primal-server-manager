package Server

import (
	"primal-server-manager/pkg/Server/Config"
	"primal-server-manager/pkg/Server/Plugin"
)

type Server struct {
	UUID      string
	Plugins   []Plugin.Plugin
	Config    Config.Config
	Files     []Config.File
	Templates struct {
		Name        string
		Description string
	}
}
