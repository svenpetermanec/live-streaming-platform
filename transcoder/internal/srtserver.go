package internal

import (
	"context"
	"log"
	"net"

	"github.com/haivision/srtgo"

	"transcoder/internal/logging"
	"transcoder/transcoder/cmd/config"
)

type SrtServer struct {
	socket            *srtgo.SrtSocket
	sourceToStreamMap map[string]string

	manager *Manager
	logger  *logging.Logger
	config  config.SRTConfig
}

func NewSrtServer(manager *Manager, logger *logging.Logger, config config.SRTConfig) *SrtServer {
	options := map[string]string{
		"transtype": config.TransType,
	}

	socket := srtgo.NewSrtSocket("0.0.0.0", uint16(config.Port), options)

	return &SrtServer{
		socket:            socket,
		sourceToStreamMap: make(map[string]string),
		manager:           manager,
		logger:            logger,
		config:            config,
	}
}

func (s *SrtServer) Start(ctx context.Context) {
	defer s.socket.Close()

	err := s.socket.Listen(s.config.BacklogConnections)
	if err != nil {
		log.Fatal(err)
	}

	s.logger.Info("Started listening", nil)

	s.socket.SetListenCallback(
		func(socket *srtgo.SrtSocket, version int, addr *net.UDPAddr, streamId string) bool {
			s.sourceToStreamMap[addr.IP.String()+addr.AddrPort().String()] = streamId
			// check whitelist if needed
			return true
		},
	)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Shutting down", nil)
		default:
			socket, addr, err := s.socket.Accept()
			if err != nil {
				s.logger.Error("Problem accepting connection", logging.Data{"addr": addr, "err": err})
				socket.Close()
			}

			streamId := s.sourceToStreamMap[addr.IP.String()+addr.AddrPort().String()]

			s.manager.HandleNewStream(socket, streamId)
		}
	}
}
