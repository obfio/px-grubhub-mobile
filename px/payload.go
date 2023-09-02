package px

import (
	"fmt"
	"math"
	mrand "math/rand"
	"strings"
	"time"
)

// Payload is our main px payload struct
type Payload struct {
	T string `json:"t"`
	D struct {
		Px1214 string  `json:"PX1214"`
		Px91   int     `json:"PX91"`
		Px92   int     `json:"PX92"`
		Px316  bool    `json:"PX316"`
		Px318  string  `json:"PX318"`
		Px319  string  `json:"PX319"`
		Px320  string  `json:"PX320"`
		Px339  string  `json:"PX339"`
		Px321  string  `json:"PX321"`
		Px323  int64   `json:"PX323"`
		Px322  string  `json:"PX322"`
		Px337  bool    `json:"PX337"`
		Px336  bool    `json:"PX336"`
		Px335  bool    `json:"PX335"`
		Px334  bool    `json:"PX334"`
		Px333  bool    `json:"PX333"`
		Px331  bool    `json:"PX331"`
		Px332  bool    `json:"PX332"`
		Px421  string  `json:"PX421"`
		Px442  string  `json:"PX442"`
		Px317  string  `json:"PX317"`
		Px344  string  `json:"PX344"`
		Px347  string  `json:"PX347"`
		Px343  string  `json:"PX343"`
		Px415  int     `json:"PX415"`
		Px413  string  `json:"PX413"`
		Px416  string  `json:"PX416"`
		Px414  string  `json:"PX414"`
		Px419  string  `json:"PX419"`
		Px418  float64 `json:"PX418"`
		Px420  float64 `json:"PX420"`
		Px340  string  `json:"PX340"`
		Px342  string  `json:"PX342"`
		Px341  string  `json:"PX341"`
		Px348  string  `json:"PX348"`
		Px1159 bool    `json:"PX1159"`
		Px330  string  `json:"PX330"`
		Px345  int64   `json:"PX345"`
		Px351  int64   `json:"PX351"`
		Px326  string  `json:"PX326"`
		Px327  string  `json:"PX327"`
		Px328  string  `json:"PX328"`
		Px259  *int64  `json:"PX259,omitempty"`
		Px256  *string `json:"PX256,omitempty"`
		Px257  *string `json:"PX257,omitempty"`
		Px1208 string  `json:"PX1208"`
	} `json:"d"`
	Cache struct {
		StartTime    int64 `json:"-"`
		TimeStamp    int64 `json:"-"`
		PrevUUIDTime int64 `json:"-"`
	} `json:"-"`
}

// Device is an instance of scraped device data
type Device struct {
	Model          string `json:"model"`
	Manufacturer   string `json:"manufacturer"`
	Device         string `json:"device"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	GPS            bool   `json:"gps"`
	Gyro           bool   `json:"gyro"`
	Accelerometer  bool   `json:"accelerometer"`
	Ethernet       bool   `json:"ethernet"`
	TouchScreen    bool   `json:"touchScreen"`
	NFC            bool   `json:"nfc"`
	WiFi           bool   `json:"wifi"`
	AndroidVersion int    `json:"androidVersion"`
}

const hexChars = "abcdef0123456789"

func randomHex(n int, isPx419 bool) string {
	o := []byte{}
	numCount := 0
	for i := 0; i < n; i++ {
		num := mrand.Intn(16)
		if isPx419 && numCount == 4 && num > 5 {
			i--
			continue
		}
		if num > 5 {
			numCount++
		}
		o = append(o, hexChars[num])
	}
	return string(o)
}

func getkernelVersion(l int) string {
	out := ""
	switch l {
	case 29:
		out = fmt.Sprintf("4.19.%v", mrand.Intn(292))
	case 30:
		out = fmt.Sprintf("5.4.%v", mrand.Intn(254))
	case 31:
		out = fmt.Sprintf("5.10.%v", mrand.Intn(192))
	case 33:
		out = fmt.Sprintf("5.15.%v", mrand.Intn(128))

	}
	return out
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func getRandomTemp() float64 {
	return toFixed(-40.0+mrand.Float64()*100.0, 1)
}

func getRandomVoltage() float64 {
	return toFixed(1.0+mrand.Float64()*5.0, 3)
}

func (p *Payload) getBatteryStatus() string {
	if p.D.Px418 > 40.0 {
		return "overheat"
	}
	if p.D.Px418 < -20.0 {
		return "cold"
	}
	if p.D.Px420 > 4.205 {
		return "over voltage"
	}
	return "good"
}

func (p *Payload) getChargingPortStatus() string {
	if p.D.Px420 > 1.0 {
		return "USB"
	}
	return ""
}

func (p *Payload) getChargingStatus() string {
	if p.D.Px415 > 100 {
		return "full"
	}
	if p.D.Px413 != "" {
		return "charging"
	}
	return "not charging"
}

// Populate fills in all the values we don't need to write a lot of code for
func (p *Payload) Populate(d *Device, l int, sdkVer, appVer, appName, packName string, isInstantApp bool) {
	p.T = "PX315"
	// ! First, we need to have a flag that sees if we've already populated the payload before
	populated := p.Cache.StartTime != 0
	if !populated {
		// remember, we saw this was Long.toHexString(new SecureRandom().nextLong());
		// We're going to just generate a random 16 long hex string for this
		p.D.Px1214 = randomHex(16, false)
		p.D.Px319 = getkernelVersion(l)
		p.Cache.StartTime = time.Now().UnixMilli()
		// this is how many bars you have at the moment, the scraper makes sure that every device supports 5G so we can choose anything from Unknown - 5G
		p.D.Px343 = []string{"Unknown", "2G", "3G", "4G", "5G"}[mrand.Intn(5)]
		// we have 3 options for this, one of them is NA though so lets not use that one.
		p.D.Px317 = []string{"WiFi", "Mobile"}[mrand.Intn(2)]
		// gonna use a few different carriers here, doubt they'd just ban all T-Mobile users but ya never know
		p.D.Px344 = []string{"T-Mobile", "Vodafone", "Msg2Send", "Mobitel", "Cequens", "Vodacom", "MTN", "Meteor", "Android", "Movistar", "Swisscom", "Orange", "Unite", "Oxygen8", "Txtlocal", "TextOver", "Virgin-Mobile", "Aircel", "AT&T", "Cellcom", "BellSouth", "Cleartalk", "Cricket", "DTC", "nTelos", "Esendex", "Kajeet", "LongLines", "MetroPCS", "Nextech", "SMS4Free", "Solavei", "Southernlinc", "Sprint", "Teleflip", "Unicel", "Viaero", "UTBox"}[mrand.Intn(38)]
		// this is your battery level, I'll make this between 20 and 100
		p.D.Px415 = 20 + mrand.Intn(80)
		// there's many options here, however, we will need to make sure whatever we choose lines up with our voltage/temp options.
		// for this reason, we will be setting voltage and temp first.
		// to know what's cold, what's hot, or what's high voltage, we will need to do some research into batterys.
		// after some googling, I see that it's dependant on the battery to tell android that it's too hot, too cold, etc.
		// so we don't have to be percise, just a general number will work here and for this number I'll choose anything above 40C as being too hot with the max being 60C
		// It seems anything below -20C is considered too cold for a battery so that will be our mentric for that, with a low of -40C
		// as for voltage, I saw the number 4.205 mentioned on some random android fourm so I'll just use that, anything above that will be considered too high and the max will be 6.0, lowest will be 1.0
		p.D.Px418 = getRandomTemp()
		p.D.Px420 = getRandomVoltage()
		// we aren't sure if we should prioritize a the temp or voltage so we'll just prioritize the temp
		p.D.Px413 = p.getBatteryStatus()
		// for the charging port status, we're gonna make it only set it to either USB or "", this is because they may validate the model we choose has wireless charging.
		// we will only set it to charging if the voltage is > 3.5
		p.D.Px416 = p.getChargingPortStatus()
		// I don't know enough about batteries to confidently use the discharging option and I wont be using the empty string option because it should never be triggered.
		// Full will only be used if Px415 == 100
		// charging will only be used if Px413 != ""
		// else, not charging
		p.D.Px414 = p.getChargingStatus()
		// I had someone reach out after I posted part 1 and tell me that this value was actually a mistake on the part of PX and was changed in more recent SDK versions.
		// the reason I bring this up is because it means the value is easy to mess up and I've seen people messs it up before when open-sourcing their APis
		// it's important you make everything after the `@` hex characters and 7 characters.
		// Also note bv0.a, this is the class and method that the payload building code is running in but this may not be the same between every APK using this SDK version so it's important you validate you have the correct calss and method name there.
		// edit: we're going to make sure this has at least 3 letters in it lol
		p.D.Px419 = "bv0.a@" + randomHex(7, true)
	} else {
		p.D.Px345++
		p.D.Px351 += time.Now().UnixMilli() - p.Cache.StartTime
	}
	p.D.Px91 = d.Width
	p.D.Px92 = d.Height
	p.D.Px316 = true
	p.D.Px318 = fmt.Sprint(l)
	p.D.Px320 = strings.ToUpper(d.Model)
	p.D.Px339 = d.Manufacturer
	p.D.Px321 = d.Device
	// remember, this timestamp is also what's used to make the UUID so we must add to the cache.
	p.Cache.TimeStamp = time.Now().UnixMilli()
	p.D.Px323 = p.Cache.TimeStamp / 1000
	p.D.Px322 = "Android"
	p.D.Px337 = d.GPS
	p.D.Px336 = d.Gyro
	p.D.Px335 = d.Accelerometer
	// ethernet is hardcoded as false since I don't know how it really works, I'm just assuming most people wont have this feature enabled
	p.D.Px334 = false
	p.D.Px333 = d.TouchScreen
	p.D.Px331 = d.NFC
	p.D.Px332 = d.WiFi
	// These next 2 are basically checking if our device is rooted, most android users wont root their device so lets keep this false.
	p.D.Px421 = "false"
	p.D.Px442 = "false"
	// I'll use a static language here because they might cross-check it to the IP and I'm not sure about that yet
	p.D.Px347 = "[en_US]"
	p.D.Px340 = sdkVer
	p.D.Px342 = appVer
	p.D.Px341 = appName
	p.D.Px348 = packName
	p.D.Px1159 = isInstantApp
	p.D.Px330 = "new_session"
	p.D.Px1208 = "[]"
}
