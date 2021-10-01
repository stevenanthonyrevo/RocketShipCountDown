package main
import(
 "net/http"
 "time"
 "github.com/gin-gonic/gin"
 "log"
)

func rocket(c *gin.Context) {
 c.HTML(http.StatusOK, "rocket.templ.html", gin.H{
    "timestamp": time.Now().Format("2006-01-02 15:04:05.000"),
 })
}

func main() {
 gin.SetMode(gin.ReleaseMode)
 router := gin.New()

 router.LoadHTMLGlob("templates/*templ.html")
 router.GET("/", rocket)
 port :=" 8888"
 if err := router.Run(":" + port); err != nil {
		// Logger
		log.Panicf("error: %s", err)
 }
}
