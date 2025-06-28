package internal

import (
	"context"

	"github.com/haivision/srtgo"

	"transcoder/internal/logging"
	"transcoder/transcoder/cmd/config"
)

type Manager struct {
	repository *Repository
	logger     *logging.Logger
}

func NewSrtManager(repository *Repository, logger *logging.Logger) *Manager {
	return &Manager{
		repository: repository,
		logger:     logger,
	}
}

func (s *Manager) HandleNewStream(ctx context.Context, stream *srtgo.SrtSocket, streamId string) {
	streamName, err := s.repository.GetStreamName(ctx, streamId)
	if err != nil {
		s.logger.Error("Problem fetching stream name", logging.Data{"stream": streamId, "error": err})
		stream.Close()
		return
	}

	transcoder, err := NewTranscoder(streamName, config.Cfg.HLS)
	if err != nil {
		s.logger.Error("Problem creating transcoder", logging.Data{"streamName": streamName, "error": err})
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
				s.logger.Warn("Problem reading stream", logging.Data{"streamName": streamName, "error": err})
			}

			if size == 0 {
				stream.Close()
				break
			}

			if err = transcoder.Write(buffer[:size]); err != nil {
				s.logger.Error("Problem transcoding", logging.Data{"streamName": streamName, "error": err})
				break
			}
		}

		if err = transcoder.Stop(); err != nil {
			s.logger.Warn("Problem stopping transcoder", logging.Data{"streamName": streamName, "error": err})
		}
	}()
}
