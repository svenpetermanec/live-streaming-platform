package guard

import (
	"errors"

	"transcoder/internal/logging"
)

func CapturePanic(logger *logging.Logger) {
	if r := recover(); r != nil {
		var err error

		switch recoverType := r.(type) {
		case string:
			err = errors.New(recoverType)
		case error:
			err = recoverType
		default:
			err = errors.New("unknown panic")
		}

		context := map[string]interface{}{
			"error": err,
		}

		logger.Error("panic", context)
	}
}
