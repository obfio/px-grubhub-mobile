package px

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	http "github.com/bogdanfinn/fhttp"
	"github.com/justhyped/OrderedForm"
	uuid "github.com/satori/go.uuid"
)

// GetConfig does the px-conf request and sets c.Config to a *Config struct
func (c *Client) GetConfig() error {
	body := fmt.Sprintf(`{"device_os_version":"%v","device_os_name":"Android","device_model":"%s","sdk_version":"%s","app_version":"%s","app_id":"%s"}`, c.DeviceData.AndroidVersion, strings.ToUpper(c.DeviceData.Model), c.SDKVersion, c.AppVersion, c.PXID)
	req, err := http.NewRequest("POST", "https://px-conf.perimeterx.net/api/v1/mobile", strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header = c.FormatHeaders(`Host: px-conf.perimeterx.net|Accept-Charset: UTF-8|Accept: */*|Content-Type: application/json|Connection: close`)
	req.Header.Set("User-Agent", "PerimeterX Android SDK/"+c.SDKVersion[1:])
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("http resp code: %v", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	f := &Config{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		return err
	}
	c.Config = f
	return nil
}

// GetInstructions - first PX request that gets instructions, mainly the APPC instruction used for the second request
func (c *Client) GetInstructions() (*Instructions, error) {
	p, err := json.Marshal([]*Payload{c.PXPayload})
	if err != nil {
		return nil, err
	}
	body := new(OrderedForm.OrderedForm)
	body.Set("payload", base64.StdEncoding.EncodeToString(p))
	body.Set("uuid", c.UUIDv4)
	body.Set("appId", c.PXID)
	body.Set("tag", "mobile")
	body.Set("ftag", "22")
	finalResult := strings.ReplaceAll(strings.ReplaceAll(body.URLEncode(), "%3D", "="), "\\u0026", "&")
	c.Payloads = append(c.Payloads, finalResult)
	req, err := http.NewRequest("POST", fmt.Sprintf("https://collector-%s.perimeterx.net/api/v1/collector/mobile", strings.ToLower(c.PXID)), strings.NewReader(finalResult))
	if err != nil {
		return nil, err
	}
	req.Header = c.FormatHeaders(`Accept-Charset: UTF-8|Accept: */*|Content-Type: application/x-www-form-urlencoded; charset=utf-8`)
	req.Header.Set("User-Agent", "PerimeterX Android SDK/"+c.SDKVersion[1:])
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http resp code: %v", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	f := &Instructions{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// GetCookie - last PX request, this request should return the instructions that give you the `bake` aka _px2 cookie/header
func (c *Client) GetCookie() (*Instructions, error) {
	p, err := json.Marshal([]*Payload{c.PXPayload})
	if err != nil {
		return nil, err
	}
	body := new(OrderedForm.OrderedForm)
	body.Set("payload", base64.StdEncoding.EncodeToString(p))
	body.Set("uuid", c.UUIDv4)
	body.Set("appId", c.PXID)
	body.Set("tag", "mobile")
	body.Set("ftag", "22")
	body.Set("sid", c.SID)
	body.Set("vid", c.VID)
	finalResult := strings.ReplaceAll(strings.ReplaceAll(body.URLEncode(), "%3D", "="), "\\u0026", "&")
	c.Payloads = append(c.Payloads, finalResult)
	req, err := http.NewRequest("POST", fmt.Sprintf("https://collector-%s.perimeterx.net/api/v1/collector/mobile", strings.ToLower(c.PXID)), strings.NewReader(finalResult))
	if err != nil {
		return nil, err
	}
	req.Header = c.FormatHeaders(`Accept-Charset: UTF-8|Accept: */*|Content-Type: application/x-www-form-urlencoded; charset=utf-8`)
	req.Header.Set("User-Agent", "PerimeterX Android SDK/"+c.SDKVersion[1:])
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http resp code: %v", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	f := &Instructions{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (c *Client) getAuth() (*GrubHubAuth, error) {
	req, err := http.NewRequest("POST", "https://api-gtm.grubhub.com/auth/anon", strings.NewReader(`{"brand":"GRUBHUB","client_id":"ghandroid_Ujtwar5s9e3RYiSNV31X41y2hsK6Kh1Uv7JDrkpS","scope":"anonymous"}`))
	if err != nil {
		return nil, err
	}
	req.Header = c.FormatHeaders(`Host: api-gtm.grubhub.com|Vary: Accept-Encoding|Accept: */*|X-Px-Authorization: 1|Content-Type: application/json; charset=utf-8`)
	req.Header.Set("User-Agent", fmt.Sprintf("Grubhub/2023.25 (%s; Android %v)", c.DeviceData.Model, c.DeviceData.AndroidVersion))
	req.Header.Set("X-Gh-Features", fmt.Sprintf("0=phone;1=grubhub 2023.25.0;2=Android %v;", c.DeviceData.AndroidVersion))
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http resp code: %v", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	f := &GrubHubAuth{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

var (
	i      = 0
	locker sync.RWMutex
)

func (c *Client) logDebug(debug *LoginDebug) {
	b, _ := json.MarshalIndent(debug, "", "	")
	locker.Lock()
	os.WriteFile(fmt.Sprintf("./tests/%v.json", i), b, 0666)
	i++
	locker.Unlock()
}

// TestCookie uses the grubhub login endpoint to test if our PX cookies pass or not
func (c *Client) TestCookie(cookie, email, password string) error {
	debug := &LoginDebug{Payloads: c.Payloads}
	auth, err := c.getAuth()
	if err != nil {
		debug.Error = err.Error()
		c.logDebug(debug)
		return err
	}
	if auth.Session.AuthToken == "" {
		debug.Error = "no auth token found"
		c.logDebug(debug)
		return errors.New("no auth token found")
	}
	body := fmt.Sprintf(`{"brand":"GRUBHUB","client_id":"ghandroid_Ujtwar5s9e3RYiSNV31X41y2hsK6Kh1Uv7JDrkpS","email":"%s","password":"%s","exclusive_session":false,"scope":"diner","device_id":"%s","device_public_key":null,"metadata_map":null}`, email, password, uuid.Must(uuid.NewV4()))
	req, err := http.NewRequest("POST", "https://api-gtm.grubhub.com/auth/login", strings.NewReader(body))
	if err != nil {
		debug.Error = err.Error()
		c.logDebug(debug)
		return err
	}
	req.Header = c.FormatHeaders(`Host: api-gtm.grubhub.com|Vary: Accept-Encoding|Accept: */*|Content-Type: application/json; charset=utf-8`)
	req.Header.Set("User-Agent", fmt.Sprintf("Grubhub/2023.25 (%s; Android %v)", c.DeviceData.Model, c.DeviceData.AndroidVersion))
	req.Header.Set("X-Gh-Features", fmt.Sprintf("0=phone;1=grubhub 2023.25.0;2=Android %v;", c.DeviceData.AndroidVersion))
	req.Header.Set("X-Px-Authorization", "2:"+cookie)
	req.Header.Set("Authorization", "Bearer "+auth.Session.AuthToken)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		debug.Error = err.Error()
		c.logDebug(debug)
		return err
	}
	defer resp.Body.Close()
	debug.StatusCode = resp.StatusCode
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		debug.Error = err.Error()
		c.logDebug(debug)
		return err
	}
	debug.Body = string(b)
	debug.CookieUsed = cookie
	c.logDebug(debug)
	return nil
}
