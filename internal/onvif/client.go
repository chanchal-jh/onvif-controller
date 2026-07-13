package onvifservice

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/0x524a/onvif-go"
)

const cameraTimeout = 30 * time.Second

// CreateClient builds an ONVIF client for the given endpoint, using a
// permissive TLS config so self-signed camera certs (without IP SANs)
// don't get rejected.
func CreateClient(endpoint, username, password string) (*onvif.Client, error) {
	return onvif.NewClient(
		endpoint,
		onvif.WithCredentials(username, password),
		onvif.WithTimeout(cameraTimeout),
		onvif.WithHTTPClient(insecureHTTPClient(cameraTimeout)),
	)
}

func GetEndpointByIP(ip string) (string, error) {

	devices, err := DiscoverCameras()
	if err != nil {
		return "", fmt.Errorf("discovery failed: %w", err)
	}

	for _, d := range devices {
		if endpointHasIP(d.Endpoint, ip) {
			return d.Endpoint, nil
		}
	}
	return "", fmt.Errorf("%w: %s", ErrCameraNotFound, ip)
}

func endpointHasIP(endpoint, ip string) bool {
	u, err := url.Parse(endpoint)
	if err != nil {
		return false
	}
	return u.Hostname() == ip
}

// insecureHTTPClient returns an http.Client that skips TLS certificate
// verification. ONVIF cameras commonly ship self-signed certs without
// proper IP SANs, so strict verification fails even for legitimate devices
// on a trusted local network.
func insecureHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // local/trusted network only
			},
		},
	}
}

// isHTTPEndpoint reports whether endpoint uses the plain http scheme.
func isHTTPEndpoint(endpoint string) bool {
	u, err := url.Parse(endpoint)
	if err != nil {
		// Can't determine scheme reliably; fall back to a simple prefix check.
		return strings.HasPrefix(endpoint, "http://")
	}
	return u.Scheme == "http"
}

func GetInitializedClient(
	ip,
	username,
	password string,
) (*onvif.Client, string, context.Context, context.CancelFunc, error) {

	endpoint, err := GetEndpointByIP(ip)
	if err != nil {
		return nil, "", nil, nil, err
	}

	client, err := CreateClient(endpoint, username, password)
	if err != nil {
		return nil, "", nil, nil, err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		cameraTimeout,
	)

	if err := client.Initialize(ctx); err != nil {
		cancel()
		return nil, "", nil, nil, err
	}

	return client, endpoint, ctx, cancel, nil
}

func GetProfileToken(
	client *onvif.Client,
	ctx context.Context,
) (string, error) {

	profiles, err := client.GetProfiles(ctx)
	if err != nil {
		return "", err
	}

	if len(profiles) == 0 {
		return "", fmt.Errorf("no profiles found")
	}

	return profiles[0].Token, nil
}
