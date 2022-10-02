package main

import (
	"log"
	"net/http"
	"time"
    "fmt"
    "os"
    "context"
    "sync"

	"github.com/gin-gonic/gin"
)

var countdown string

func RocketCountDownTimer() {
    var timeslot string;
    slot1 := make(chan string)
    timer := time.NewTicker(1 * time.Second)
    go func() {
     for _ = range timer.C {
     now := time.Now()
     slot1 <- now.Format(time.Stamp);
     }
    }()
    timeslot = <-slot1
    fmt.Fprintf(os.Stdout, "%s.\n", timeslot)
    //fmt.Fprintf(os.Stdout, "%s.\n", timer)
    countdown = timeslot
}

func Start() {
    gin.SetMode(gin.ReleaseMode)
    router := gin.New()

    router.LoadHTMLGlob("templates/*templ.html")
    router.GET("/", rocket)
    router.StaticFS("/static", http.Dir("static"))
    port := "8888"

    log.Println("Starting http server...")
    if err := router.Run(":" + port); err != nil {
        // Logger
        log.Panicf("error: %s", err)
    }
}

func rocket(c *gin.Context) {
	fmt.Fprintf(os.Stdout, "%s.\n", countdown)
	countdownvalue := countdown
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	fmt.Fprintf(os.Stdout, "%s.\n", timestamp)
    c.HTML(http.StatusOK, "rocket.templ.html", gin.H{
		"timestamp": countdownvalue,
	})
}

func main() {

   go Start()

   ticker := time.NewTicker(5 * time.Second)
   quit := make(chan struct{})
    _, cancel := context.WithTimeout(context.Background(), 1*time.Second)

   var wg sync.WaitGroup
   wg.Add(1)
   go func() {
    for {
		select {
	 	case <- ticker.C:
			go RocketCountDownTimer()
		case <- quit:
			ticker.Stop()
			return
        }
     }
   }()
   wg.Wait()
   cancel()
}
