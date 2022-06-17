package subcommand

import (
	"github.com/morphysm/famed-github-backend/internal/devtoolkit"
	"github.com/morphysm/famed-github-backend/internal/server"
	"github.com/rotisserie/eris"
)

type Server struct {
	DevToolkit *devtoolkit.DevToolkit
	Server     *server.Server
}

func NewServer(devtoolkit *devtoolkit.DevToolkit) (*Server, error) {
	server, err := server.NewServer(devtoolkit)
	if err != nil {
		return nil, eris.Wrap(err, "failed to instantiate http server")
	}

	return &Server{
		DevToolkit: devtoolkit,
		Server:     server,
	}, nil
}

func (a *Server) Start() error {
	// Starts HTTP server
	err := a.Server.Start()
	if err != nil {
		return eris.New("failed to start server")
	}

	return nil
}
