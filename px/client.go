package px

import (
	"errors"
	"fmt"
	"strings"

	http "github.com/bogdanfinn/fhttp"

	tlsclient "github.com/bogdanfinn/tls-client"
	uuid "github.com/satori/go.uuid"
)

// Client - holds our PX payload struct, proxy, device, and HTTP client
type Client struct {
	PXID         string
	PXPayload    *Payload
	Proxy        string // returned to API
	DeviceData   *Device
	HTTPClient   tlsclient.HttpClient
	SDKVersion   string
	AppVersion   string
	AppName      string
	PackageName  string
	IsInstantApp bool
	Config       *Config
	UUIDv4       string
	SID          string
	VID          string
	Payloads     []string
}

// MakeClient - makes a *Client struct given a proxy and device as well as app/px info
func MakeClient(PXID, proxy, sdkVer, appVer, appName, packName string, isInstantApp bool, d *Device) (*Client, error) {
	if proxy == "" {
		return nil, errors.New("no proxy input")
	}
	opts := []tlsclient.HttpClientOption{
		tlsclient.WithTimeoutSeconds(30),
		tlsclient.WithProxyUrl(proxy),
		tlsclient.WithInsecureSkipVerify(),
	}
	switch d.AndroidVersion {
	case 10:
		opts = append(opts, tlsclient.WithClientProfile(tlsclient.Okhttp4Android10))
		break
	case 11:
		opts = append(opts, tlsclient.WithClientProfile(tlsclient.Okhttp4Android11))
		break
	case 12:
		opts = append(opts, tlsclient.WithClientProfile(tlsclient.Okhttp4Android12))
		break
	case 13:
		opts = append(opts, tlsclient.WithClientProfile(tlsclient.Okhttp4Android13))
		break
	}
	h, err := tlsclient.NewHttpClient(nil, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		PXID:         PXID,
		HTTPClient:   h,
		Proxy:        proxy,
		DeviceData:   d,
		PXPayload:    &Payload{},
		SDKVersion:   sdkVer,
		AppVersion:   appVer,
		AppName:      appName,
		PackageName:  packName,
		IsInstantApp: isInstantApp,
		UUIDv4:       fmt.Sprint(uuid.Must(uuid.NewV4())),
	}, nil
}

// FormatHeaders turns a string of headers seperated by `|` into a http.Header map
func (c *Client) FormatHeaders(h string) http.Header {
	headers := http.Header{}
	for _, header := range strings.Split(h, "|") {
		parts := strings.Split(header, ": ")
		headers.Set(parts[0], parts[1])
	}
	return headers
}
