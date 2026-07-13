package onvifservice

import (
	"errors"
	"net/http"
)

var (
	ErrCameraNotFound     = errors.New("no camera discover with ip:")
	ErrNoProfiles         = errors.New("no profiles found")
	ErrPanTiltUnsupported = errors.New("camera does not support pan/tilt")
	ErrZoomUnsupported    = errors.New("camera does not support zoom")
	ErrUnsupportedAction  = errors.New("unsupported action")
)

func GetHTTPStatus(err error) int {
	switch {
	case errors.Is(err, ErrCameraNotFound):
		return http.StatusBadRequest

	case errors.Is(err, ErrNoProfiles):
		return http.StatusBadRequest

	case errors.Is(err, ErrPanTiltUnsupported):
		return http.StatusBadRequest

	case errors.Is(err, ErrZoomUnsupported):
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}
