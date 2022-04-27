package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/gin-gonic/gin"
)

// go:embedded_file frontend/dist/index.html
var FS embed.FS

func main() {
	go func() {
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default() // router
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		r.StaticFS("/static", http.FS(staticFiles)) // router
		r.NoRoute(func(c *gin.Context) {            // router
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/static/") {
				reader, err := staticFiles.Open("index.html")
				if err != nil {
					log.Fatal(err)
				}
				defer reader.Close()
				stat, err := reader.Stat()
				if err != nil {
					log.Fatal(err)
				}
				c.DataFromReader(http.StatusOK, stat.Size(), "text/html", reader, nil)
			} else {
				c.Status(http.StatusNotFound)
			}
		})
		r.Run(":8080")
	}()
	c := make(chan os.Signal, 1) // chSignal
	signal.Notify(c, os.Interrupt)

	chromePath := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:8080/static/index.html")
	cmd.Start()

	select {
	case <-c:
		cmd.Process.Kill()
	}
}
