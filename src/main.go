package main

import (
	"net"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	whois "github.com/likexian/whois-go"
	ping "github.com/sparrc/go-ping"
)

// Ping request payload
type PingPayload struct {
	Target string `json:"target" binding:"required"`
}

func pingHandler(c *gin.Context) {
	var data PingPayload
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

func whoisHandler(c *gin.Context) {
	var data PingPayload
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

func nslookupHandler(c *gin.Context) {
	var data PingPayload
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

func main() {
	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/ping", pingHandler)
	r.POST("/whois", whoisHandler)
	r.POST("/nslookup", nslookupHandler)
	r.Static("/", "/var/www/html")
	r.Run(":80")
}
