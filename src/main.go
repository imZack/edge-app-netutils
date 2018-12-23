package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	fast "github.com/ddo/go-fast"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	whois "github.com/likexian/whois-go"
	ping "github.com/sparrc/go-ping"
)

type TargetPayload struct {
	Target string `json:"target" binding:"required"`
}

// Handler for PING check
func pingHandler(c *gin.Context) {
	var data TargetPayload
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pinger, err := ping.NewPinger(data.Target)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pinger.SetPrivileged(true)
	pinger.Count = 3
	pinger.Run()
	stats := pinger.Statistics()
	c.JSON(http.StatusOK, gin.H{"status": stats})
}

// Handler for Domain Infomation lookup
func whoisHandler(c *gin.Context) {
	var data TargetPayload
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := whois.Whois(data.Target)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"record": result})
}

// Handler for DNS lookup
func nslookupHandler(c *gin.Context) {
	var data TargetPayload
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ips, err := net.LookupIP(data.Target)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ips": ips})
}

// Handler for SpeedTest
func fastHandler(c *gin.Context) {
	fastCom := fast.New()
	err := fastCom.Init()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	urls, err := fastCom.GetUrls()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	KbpsChan := make(chan float64)
	go func() {
		for Kbps := range KbpsChan {
			fmt.Printf("%.2f Kbps %.2f Mbps\n", Kbps, Kbps/1000)
			c.SSEvent("message", fmt.Sprintf("%.2f Kbps %.2f Mbps\n", Kbps, Kbps/1000))
		}

		fmt.Println("done")
	}()

	err = fastCom.Measure(urls, KbpsChan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	api := r.Group("api")
	{
		api.POST("/ping", pingHandler)
		api.POST("/whois", whoisHandler)
		api.POST("/nslookup", nslookupHandler)
		api.POST("/fast", fastHandler)
	}

	staticPath := "/var/www/html"
	if sp := os.Getenv("HTML_PATH"); sp != "" {
		staticPath = sp
	}

	r.Static("/static", staticPath)
	r.Run(":80")
}
