package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Immortanbird/api_using_go/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) DownloadHandler(c *gin.Context) {
	filePath := "./test_material.pdf"

	// 1. Open the file normally
	file, err := os.Open(filePath)
	if err != nil {
		zap.L().Warn(
			"File not found",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer file.Close()

	// 2. Get file info (needed for Content-Length header)
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get file info"})
		return
	}

	// 3. Define the Speed Limit
	// Let's set it to 500 KB/s
	const KB = 1024
	speedLimit := 500 * KB

	// 4. Wrap the file in our ThrottledReader
	// We pass c.Request.Context() so that if the user cancels the download,
	// the limiter stops waiting and releases resources.
	throttledFile := utils.NewThrottledReader(file, speedLimit, c.Request.Context())

	// 5. Serve the stream using Gin
	// headers: content length, content type, file name
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, fileInfo.Name()),
	}

	c.DataFromReader(
		http.StatusOK,
		fileInfo.Size(),
		"application/octet-stream", // Or verify the actual MIME type
		throttledFile,              // Pass our custom reader here
		extraHeaders,
	)
}
