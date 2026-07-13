package onvifservice

// StreamProfile represents camera stream profile
type StreamProfile struct {
	Name        string `json:"name"`
	Token       string `json:"token"`
	StreamURI   string `json:"stream_uri"`
	SnapshotURI string `json:"snapshot_uri"`
}

// GetStreamProfiles returns camera stream profiles
func GetStreamProfiles(
	ip,
	username,
	password string,
) ([]StreamProfile, error) {

	client, _, ctx, cancel, err := GetInitializedClient(
		ip,
		username,
		password,
	)
	if err != nil {
		return nil, err
	}
	defer cancel()

	profiles, err := client.GetProfiles(ctx)
	if err != nil {
		return nil, err
	}
	if len(profiles) == 0 {
		return nil, ErrNoProfiles
	}

	var streamProfiles []StreamProfile

	for _, profile := range profiles {

		streamURI, _ := client.GetStreamURI(
			ctx,
			profile.Token,
		)

		snapshotURI, _ := client.GetSnapshotURI(
			ctx,
			profile.Token,
		)

		stream := StreamProfile{
			Name:  profile.Name,
			Token: profile.Token,
		}

		if streamURI != nil {
			stream.StreamURI = streamURI.URI
		}

		if snapshotURI != nil {
			stream.SnapshotURI = snapshotURI.URI
		}

		streamProfiles = append(
			streamProfiles,
			stream,
		)
	}

	return streamProfiles, nil
}
