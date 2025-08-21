package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tyrads-com/tyrads-go-sdk-iframe"
)

var API_KEY = "YOUR_API_KEY"       // Replace with your actual API key
var API_SECRET = "YOUR_API_SECRET" // Replace with your actual API secret
var LANGUAGE = "en"
var AGE = 18                                // Replace with actual age
var GENDER = 1                              // 1 for male, 2 for female
var PUBLISHER_USER_ID = "PUBLISHER_USER_ID" // Replace with actual publisher user ID

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

var tyrAdsSdk = tyrads.NewTyrAdsSdk(
	API_KEY,
	API_SECRET,
	LANGUAGE)

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/", func(c *gin.Context) {
		sign, err := tyrAdsSdk.Authenticate(tyrads.AuthenticationRequest{
			Age:             AGE,
			Gender:          GENDER,
			PublisherUserID: PUBLISHER_USER_ID,
		})

		if err != nil || sign == nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		token := sign.Token
		success := token != ""

		iframeUrl, err := tyrAdsSdk.IframeUrl(token, nil)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate iframe URL: " + err.Error()})
			return
		}

		iframePremiumUrl, err := tyrAdsSdk.IframePremiumWidget(token, nil)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate iframe premium URL: " + err.Error()})
			return
		}

		c.JSON(200, gin.H{"success": success, "token": token, "iframeUrl": iframeUrl, "iframePremiumUrl": iframePremiumUrl})
	})
	r.Run()
}
