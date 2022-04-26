package main

import (
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	go func() {
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()
		r.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "<h1>Hello World</h1>")
		})
		r.Run(":8080")
	}()
	time.Sleep(time.Second * 3)
	chromePath := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	cnd := exec.Command(chromePath, "--app=http://127.0.0.1:8080/")
	cnd.Start()
	select {}
}
