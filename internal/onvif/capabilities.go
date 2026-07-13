package onvifservice

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/use-go/onvif"
	"github.com/use-go/onvif/ptz"
	onvifxsd "github.com/use-go/onvif/xsd/onvif"
)

type PTZSpaces struct {
	PanTiltVelocitySpace string
	ZoomVelocitySpace    string
}

type PTZCapabilities struct {
	PanTilt        bool `json:"pan_tilt"`
	Zoom           bool `json:"zoom"`
	AbsoluteMove   bool `json:"absolute_move"`
	RelativeMove   bool `json:"relative_move"`
	ContinuousMove bool `json:"continuous_move"`

	Spaces PTZSpaces `json:"spaces"`
}

type GetNodesEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		GetNodesResponse struct {
			PTZNode []onvifxsd.PTZNode `xml:"PTZNode"`
		} `xml:"GetNodesResponse"`
	} `xml:"Body"`
}

func IsPTZSupported(device *onvif.Device) bool {

	resp, err := device.CallMethod(
		ptz.GetServiceCapabilities{},
	)

	if err != nil {
		return false
	}

	if resp != nil {
		resp.Body.Close()
	}

	return true
}

func GetPTZCapabilities(
	ip,
	username,
	password string,
) (*PTZCapabilities, error) {

	device, err := onvif.NewDevice(
		onvif.DeviceParams{
			Xaddr:    ip,
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to camera: %w", err)
	}
	if !IsPTZSupported(device) {
		return nil, fmt.Errorf("PTZ not supported")
	}

	resp, err := device.CallMethod(ptz.GetNodes{})
	if err != nil {
		return nil, fmt.Errorf("GetNodes failed: %w", err)
	}
	defer resp.Body.Close()

	rawXML, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var envelope GetNodesEnvelope

	if err := xml.Unmarshal(rawXML, &envelope); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	capabilities := &PTZCapabilities{}

	for _, node := range envelope.Body.GetNodesResponse.PTZNode {

		spaces := node.SupportedPTZSpaces

		capabilities.Spaces.PanTiltVelocitySpace =
			string(spaces.ContinuousPanTiltVelocitySpace.URI)

		capabilities.Spaces.ZoomVelocitySpace =
			string(spaces.ContinuousZoomVelocitySpace.URI)

		if spaces.ContinuousPanTiltVelocitySpace.URI != "" ||
			spaces.ContinuousZoomVelocitySpace.URI != "" {
			capabilities.ContinuousMove = true
		}

		if spaces.AbsolutePanTiltPositionSpace.URI != "" ||
			spaces.AbsoluteZoomPositionSpace.URI != "" {
			capabilities.AbsoluteMove = true
		}

		if spaces.RelativePanTiltTranslationSpace.URI != "" ||
			spaces.RelativeZoomTranslationSpace.URI != "" {
			capabilities.RelativeMove = true
		}

		if spaces.AbsolutePanTiltPositionSpace.URI != "" ||
			spaces.RelativePanTiltTranslationSpace.URI != "" ||
			spaces.ContinuousPanTiltVelocitySpace.URI != "" {
			capabilities.PanTilt = true
		}

		if spaces.AbsoluteZoomPositionSpace.URI != "" ||
			spaces.RelativeZoomTranslationSpace.URI != "" ||
			spaces.ContinuousZoomVelocitySpace.URI != "" {
			capabilities.Zoom = true
		}
	}

	return capabilities, nil
}
