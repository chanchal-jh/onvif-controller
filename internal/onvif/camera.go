package onvifservice

// CameraInfo represents camera device information
type CameraInfo struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Firmware     string `json:"firmware"`
	SerialNumber string `json:"serial_number"`
	HardwareID   string `json:"hardware_id"`
}

// ConnectCamera connects to camera and returns device info
func ConnectCamera(
	ip,
	username,
	password string,
) (*CameraInfo, error) {

	client, _, ctx, cancel, err := GetInitializedClient(
		ip,
		username,
		password,
	)
	if err != nil {
		return nil, err
	}
	defer cancel()

	info, err := client.GetDeviceInformation(ctx)
	if err != nil {
		return nil, err
	}

	return &CameraInfo{
		Manufacturer: info.Manufacturer,
		Model:        info.Model,
		Firmware:     info.FirmwareVersion,
		SerialNumber: info.SerialNumber,
		HardwareID:   info.HardwareID,
	}, nil
}
