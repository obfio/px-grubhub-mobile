package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/obfio/px-grubhub-mobile/px"
)

type input struct {
	PXID         string `json:"pxid"`
	Proxy        string `json:"proxy"`
	SDKVersion   string `json:"sdkVersion"`
	AppVersion   string `json:"appVersion"`
	AppName      string `json:"appName"`
	PackageName  string `json:"packageName"`
	IsInstantApp bool   `json:"isInstantApp"`
}

type output struct {
	Cookie string `json:"cookie"`
	Error  string `json:"error"`
}

func bake(c *gin.Context) {
	o := &output{}
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		o.Error = err.Error()
		c.AbortWithStatusJSON(500, o)
		return
	}
	i := &input{}
	err = json.Unmarshal(b, &i)
	if err != nil {
		o.Error = err.Error()
		c.AbortWithStatusJSON(500, o)
		return
	}
	client, err := px.MakeClient(i.PXID, i.Proxy, i.SDKVersion, i.AppVersion, i.AppName, i.PackageName, i.IsInstantApp, devices[rand.Intn(len(devices))])
	if err != nil {
		o.Error = err.Error()
		c.AbortWithStatusJSON(500, o)
		return
	}
	err = client.GetConfig()
	if err != nil {
		o.Error = err.Error()
		c.AbortWithStatusJSON(500, o)
		return
	}
	client.PXPayload.Populate(client.DeviceData, levels[fmt.Sprint(client.DeviceData.AndroidVersion)], client.SDKVersion, client.AppVersion, client.AppName, client.PackageName, client.IsInstantApp)
	err = client.PXPayload.UUIDSection()
	if err != nil {
		o.Error = err.Error()
		c.AbortWithStatusJSON(500, o)
		return
	}
	instructions, err := client.GetInstructions()
	if err != nil {
		o.Error = err.Error()
		c.AbortWithStatusJSON(500, o)
		return
	}
	if len(instructions.Do) < 4 {
		for i := 0; i < 3; i++ {
			client.PXPayload.Populate(client.DeviceData, levels[fmt.Sprint(client.DeviceData.AndroidVersion)], client.SDKVersion, client.AppVersion, client.AppName, client.PackageName, client.IsInstantApp)
			err = client.PXPayload.UUIDSection()
			if err != nil {
				o.Error = err.Error()
				c.AbortWithStatusJSON(500, o)
				return
			}
			instructions, err = client.GetInstructions()
			if err != nil {
				o.Error = err.Error()
				c.AbortWithStatusJSON(500, o)
				return
			}
			if len(instructions.Do) == 4 {
				break
			}
		}
		if len(instructions.Do) != 4 {
			o.Error = "unable to get instructions"
			c.AbortWithStatusJSON(500, o)
			return
		}
	}
	for _, instruction := range instructions.Do {
		parts := strings.Split(instruction, "|")
		if len(parts) < 1 {
			continue
		}
		if parts[0] == "sid" {
			client.SID = parts[1]
			continue
		}
		if parts[0] == "vid" {
			client.VID = parts[1]
			continue
		}
		if parts[0] == "appc" && parts[1] == "2" {
			client.PXPayload.AppcInstruction(instruction)
			continue
		}
	}
	instructions, err = client.GetCookie()
	if err != nil {
		o.Error = err.Error()
		c.AbortWithStatusJSON(500, o)
		return
	}
	for _, instruction := range instructions.Do {
		parts := strings.Split(instruction, "|")
		if len(parts) < 1 {
			continue
		}
		if parts[0] == "bake" {
			o.Cookie = parts[3]
		}
	}
	if o.Cookie == "" {
		o.Error = "unable to get cookie"
		c.AbortWithStatusJSON(500, o)
		return
	}
	err = client.TestCookie(o.Cookie, randStr(20)+"@gmail.com", randStr(10)+"3463456!")
	if err != nil {
		o.Error = err.Error()
		c.AbortWithStatusJSON(500, o)
		return
	}
	c.JSON(200, o)
}

func randStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"[rand.Intn(52)]
	}
	return string(b)
}
