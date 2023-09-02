package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/obfio/px-grubhub-mobile/px"
)

var (
	devices = []*px.Device{}
	levels  = make(map[string]int)
)

func init() {
	rand.Seed(time.Now().UnixNano())
	b, err := os.ReadFile("devices.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &devices)
	if err != nil {
		panic(err)
	}
	b, err = os.ReadFile("levels.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &levels)
	if err != nil {
		panic(err)
	}
}

func main() {
	router := gin.New()
	router.Use(gin.Recovery())
	px := router.Group("px")
	px.POST("bake", bake)
	s := &http.Server{
		Addr:           "127.0.0.1:7356",
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("httpServe: %s", err)
	}
}
