package onvifservice

import (
	"fmt"
	"time"

	"github.com/0x524a/onvif-go"
)

const (
	defaultPTZSpeed         = 0.1
	ptzDuration             = 1 * time.Second
	ptzPanTiltVelocitySpace = "http://www.onvif.org/ver10/tptz/PanTiltSpaces/VelocityGenericSpace"
	ptzZoomVelocitySpace    = "http://www.onvif.org/ver10/tptz/ZoomSpaces/VelocityGenericSpace"
)

func HandleAction(req ActionRequest) error {

	switch req.Action {

	case ActionSetPreset:
		return handleSetPreset(req)

	case ActionGotoPreset:
		return handleGotoPreset(req)

	default:
		return handlePTZMove(req)
	}
}

func handleSetPreset(
	req ActionRequest,
) error {
	if req.PresetToken == "" {
		return fmt.Errorf("preset_token is required")
	}

	client, _, ctx, cancel, err := GetInitializedClient(
		req.IP,
		req.Username,
		req.Password,
	)
	if err != nil {
		return err
	}
	defer cancel()

	profileToken, err := GetProfileToken(client, ctx)
	if err != nil {
		return err
	}

	_, err = client.SetPreset(
		ctx,
		profileToken,
		req.PresetName,
		req.PresetToken,
	)

	return err
}

func handleGotoPreset(
	req ActionRequest,
) error {
	if req.PresetToken == "" {
		return fmt.Errorf("preset_token is required")
	}

	client, _, ctx, cancel, err := GetInitializedClient(
		req.IP,
		req.Username,
		req.Password,
	)
	if err != nil {
		return err
	}
	defer cancel()

	profileToken, err := GetProfileToken(client, ctx)
	if err != nil {
		return err
	}

	return client.GotoPreset(
		ctx,
		profileToken,
		req.PresetToken,
		nil,
	)
}

// MovePTZ moves camera in a direction
func handlePTZMove(
	req ActionRequest,
) error {
	client, endpoint, ctx, cancel, err := GetInitializedClient(
		req.IP,
		req.Username,
		req.Password,
	)
	if err != nil {
		return err
	}
	defer cancel()

	profileToken, err := GetProfileToken(client, ctx)
	if err != nil {
		return err
	}
	// use-go/onvif (used by GetPTZCapabilities) only speaks HTTP and will
	// always fail against an HTTPS endpoint, so skip the call entirely for
	// HTTPS cameras and fall back to default ONVIF generic velocity spaces.
	var caps *PTZCapabilities

	if isHTTPEndpoint(endpoint) {
		caps, err = GetPTZCapabilities(req.IP, req.Username, req.Password)
		if err != nil {
			if err.Error() == "PTZ not supported" {
				return err
			}
			fmt.Printf("GetPTZCapabilities failed for %s, falling back to defaults: %v\n", req.IP, err)
			caps = nil
		}
	}

	if caps == nil {
		caps = &PTZCapabilities{
			PanTilt: true,
			Zoom:    true,
			Spaces: PTZSpaces{
				PanTiltVelocitySpace: ptzPanTiltVelocitySpace,
				ZoomVelocitySpace:    ptzZoomVelocitySpace,
			},
		}
	} else {
		if caps.Spaces.PanTiltVelocitySpace == "" {
			caps.Spaces.PanTiltVelocitySpace = ptzPanTiltVelocitySpace
		}
		if caps.Spaces.ZoomVelocitySpace == "" {
			caps.Spaces.ZoomVelocitySpace = ptzZoomVelocitySpace
		}
	}

	switch req.Action {

	case ActionDirectionLeft,
		ActionDirectionRight,
		ActionDirectionUp,
		ActionDirectionDown:

		if !caps.PanTilt {
			return ErrPanTiltUnsupported
		}

	case ActionZoomIn,
		ActionZoomOut:

		if !caps.Zoom {
			return ErrZoomUnsupported
		}
	}

	var velocity *onvif.PTZSpeed
	ptzSpeed := defaultPTZSpeed

	if req.Speed != nil {
		ptzSpeed = *req.Speed
	}

	switch req.Action {

	case ActionDirectionLeft:
		velocity = &onvif.PTZSpeed{
			PanTilt: &onvif.Vector2D{
				X:     -ptzSpeed,
				Y:     0,
				Space: caps.Spaces.PanTiltVelocitySpace,
			},
		}

	case ActionDirectionRight:
		velocity = &onvif.PTZSpeed{
			PanTilt: &onvif.Vector2D{
				X:     ptzSpeed,
				Y:     0,
				Space: caps.Spaces.PanTiltVelocitySpace,
			},
		}

	case ActionDirectionUp:
		velocity = &onvif.PTZSpeed{
			PanTilt: &onvif.Vector2D{
				X:     0,
				Y:     ptzSpeed,
				Space: caps.Spaces.PanTiltVelocitySpace,
			},
		}

	case ActionDirectionDown:
		velocity = &onvif.PTZSpeed{
			PanTilt: &onvif.Vector2D{
				X:     0,
				Y:     -ptzSpeed,
				Space: caps.Spaces.PanTiltVelocitySpace,
			},
		}

	case ActionZoomIn:
		velocity = &onvif.PTZSpeed{
			PanTilt: &onvif.Vector2D{
				X:     0,
				Y:     0,
				Space: caps.Spaces.PanTiltVelocitySpace,
			},
			Zoom: &onvif.Vector1D{
				X:     ptzSpeed,
				Space: caps.Spaces.ZoomVelocitySpace,
			},
		}

	case ActionZoomOut:
		velocity = &onvif.PTZSpeed{
			PanTilt: &onvif.Vector2D{
				X:     0,
				Y:     0,
				Space: caps.Spaces.PanTiltVelocitySpace,
			},
			Zoom: &onvif.Vector1D{
				X:     -ptzSpeed,
				Space: caps.Spaces.ZoomVelocitySpace,
			},
		}

	default:
		return ErrUnsupportedAction
	}

	timeout := "PT1S"

	err = client.ContinuousMove(
		ctx,
		profileToken,
		velocity,
		&timeout,
	)
	if err != nil {
		return err
	}

	time.Sleep(ptzDuration)

	_ = client.Stop(
		ctx,
		profileToken,
		true,
		true,
	)

	return nil
}
