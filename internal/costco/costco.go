package costco

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"path/filepath"
)

func (r *Registry) Checksum(data []byte) string {
	return fmt.Sprintf("%x", sha256.New().Sum(data))
}

// CheckRepository takes in gin Context and a repository, validating that the repository
// exits, updating the gin context with the correct errors if it does not.
func (r *Registry) CheckRepository(c *gin.Context, repo string) bool {
	if !r.store.FindObject(repo) {
		costcoErr := Error{
			Code:    404,
			Message: errRepositoryNotFound.Error(),
			Detail:  fmt.Sprintf("repository %s not found", repo),
		}

		c.Header("Content-Type", "application/json; charset=utf-8")
		c.JSONP(404, &costcoErr)
		return false
	}
	return true
}

// UploadMonolith uploads a layer as a single chunk
func (r *Registry) UploadMonolith(c *gin.Context, repo string) {
	digest := c.Query("digest")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {

	}

	checksum := r.Checksum(body)
	if checksum != digest {
		costcoErr := Error{
			Code:    400,
			Message: DIGEST_INVALID,
			Detail:  fmt.Sprintf("checksum does not match digest"),
		}

		c.Header("Content-Type", "application/json; charset=utf-8")
		c.JSONP(400, &costcoErr)
		return
	}

	// file paths for layers should be blobs/<hash>
	key := filepath.Join("blobs", checksum)

	r.store.PutObject(key, c.Request.Body)
}
