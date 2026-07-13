package onvifservice

import (
	"context"
	"time"

	"github.com/0x524a/onvif-go/discovery"
)

const (
	defaultTimeout    = 10 * time.Second
	defaultRetryDelay = 5 * time.Second
)

// CameraDevice represents a discovered ONVIF camera
type CameraDevice struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}

// DiscoverCameras discovers ONVIF cameras on the network
func DiscoverCameras() ([]CameraDevice, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		defaultTimeout,
	)
	defer cancel()

	opts := &discovery.DiscoverOptions{}

	devices, err := discovery.DiscoverWithOptions(
		ctx,
		defaultRetryDelay,
		opts,
	)
	if err != nil {
		return nil, err
	}

	var cameras []CameraDevice

	for _, device := range devices {
		camera := CameraDevice{
			Name:     device.GetName(),
			Endpoint: device.GetDeviceEndpoint(),
		}

		cameras = append(cameras, camera)
	}

	return cameras, nil
}

// DiscoverCamerasWithInterface discovers cameras using a specific network interface
func DiscoverCamerasWithInterface(
	networkInterface string,
) ([]CameraDevice, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		defaultTimeout,
	)
	defer cancel()

	opts := &discovery.DiscoverOptions{
		NetworkInterface: networkInterface,
	}

	devices, err := discovery.DiscoverWithOptions(
		ctx,
		defaultRetryDelay,
		opts,
	)
	if err != nil {
		return nil, err
	}

	var cameras []CameraDevice

	for _, device := range devices {
		camera := CameraDevice{
			Name:     device.GetName(),
			Endpoint: device.GetDeviceEndpoint(),
		}

		cameras = append(cameras, camera)
	}

	return cameras, nil
}

// ListInterfaces returns available network interfaces
func ListInterfaces() ([]discovery.NetworkInterface, error) {
	return discovery.ListNetworkInterfaces()
}
