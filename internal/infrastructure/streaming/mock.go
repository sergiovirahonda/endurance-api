package streaming

import (
	"strconv"

	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
)

// NATS Mocks

type NatsMock struct{}

func (nm *NatsMock) RunServer(port string) *server.Server {
	natsPort, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	opts := natsserver.DefaultTestOptions
	opts.Port = natsPort
	return nm.RunServerWithOptions(&opts)
}

func (nm *NatsMock) RunServerWithOptions(opts *server.Options) *server.Server {
	srv := natsserver.RunServer(opts)
	err := srv.EnableJetStream(&server.JetStreamConfig{})
	if err != nil {
		panic(err)
	}
	return srv
}

func (nm *NatsMock) ShutdownServer(server *server.Server) {
	server.Shutdown()
}
