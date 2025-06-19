package internal

import (
	"github.com/haivision/srtgo"

	"transcoder/internal/logging"
	"transcoder/transcoder/cmd/config"
)

type Manager struct {
	logger *logging.Logger
}

func NewSrtManager(logger *logging.Logger) *Manager {
	return &Manager{
		logger: logger,
	}
}

func (s *Manager) HandleNewStream(stream *srtgo.SrtSocket, streamID string) {
	transcoder, err := NewTranscoder(streamID, config.Cfg.HLS)
	if err != nil {
		s.logger.Error("Problem creating transcoder", logging.Data{"streamId": streamID, "error": err})
		stream.Close()
		return
	}

	go func() {
		// max payload size https://github.com/Haivision/srt/blob/master/docs/API/API-socket-options.md#SRTO_PAYLOADSIZE
		buffer := make([]byte, 1456)

		for {
			size, err := stream.Read(buffer)
			if err != nil {
				// switch to metric
				s.logger.Warn("Problem reading stream", logging.Data{"streamId": streamID, "error": err})
			}

			if size == 0 {
				stream.Close()
				break
			}

			if err = transcoder.Write(buffer[:size]); err != nil {
				s.logger.Error("Problem transcoding", logging.Data{"streamId": streamID, "error": err})
				break
			}
		}

		if err = transcoder.Stop(); err != nil {
			s.logger.Warn("Problem stopping transcoder", logging.Data{"streamId": streamID, "error": err})
		}
	}()
}
