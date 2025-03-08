package costco

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Authenticator is an interface for the various middleware supported for the
// container registry. All authentication must be RFC 7235 compliant.
type Authenticator interface {
	Authenticate() gin.HandlerFunc
}

// InitBasicClient returns a basic auth client with a username and password
// statically assigned.
func InitBasicClient(username, password string) BasicClient {
	b := BasicClient{
		Username: username,
		Password: password,
	}

	return b
}

// BasicClient implements the Authenticator interface as a basic username and password
// client. BasicClient only supports a single username and password.
type BasicClient struct {
	Username string
	Password string
}

// Authenticate is middle ware used to provide authentication to the registry.
func (b *BasicClient) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, hasAuth := c.Request.BasicAuth()

		// Validate credentials
		if !hasAuth || username != b.Username || password != b.Password {
			// Send RFC 7235 compliant response
			c.Header("WWW-Authenticate", `Basic realm="Restricted Area"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": UNAUTHORIZED})
			return
		}

		c.Next() // Proceed if authentication is successful
	}
}
