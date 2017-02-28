package main

import (
	"log"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/gin-gonic/gin"
)

func JwtAuthRequired(securityConfig SecurityConfig) gin.HandlerFunc {

	return func(c *gin.Context) {
		//Only validate if JWT Authentification is enabled in config
		if securityConfig.EnableJWTAuthentification {
			authToken := c.Request.Header.Get("Authorization")
			log.Printf("Got Request with Authorization-Header: %s ", authToken)

			payload, headers, err := jose.Decode(authToken, []byte(securityConfig.JWTSharedKey))

			if err != nil {
				log.Printf("Wasn't able do validate token using shared key")
				c.AbortWithStatus(401)
			}

			log.Printf("Got valid JWT Headers:%s Payload:%s", headers, payload)
			c.Set("JWTPayload", payload)

		}
	}
}
