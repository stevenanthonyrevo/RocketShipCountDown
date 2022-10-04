package main

import (
	"log"
	"net/http"
	"time"
    "fmt"
    "os"
    "context"
    "sync"
    "net/url"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var countdown string
var live string
var messageOut = make(chan string)

func WebSocketConnection() {
    uri := url.URL{Scheme: "ws", Host: "localhost:8888", Path: "/socket",}
    c, resp, err := websocket.DefaultDialer.Dial(uri.String(), nil);
    if err != nil {
		log.Printf("handshake failed %d", resp.StatusCode)
    }
    done := make(chan struct{})

   go func() {
   	defer close(done)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		messageOut <- string(message)
	}
   }()
   x := <- messageOut
   fmt.Println(x)
   live = string(x)
}

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
    fmt.Fprintf(os.Stdout, "default_time: %s\n", timeslot)
    countdown = timeslot
}

func Start() {
    gin.SetMode(gin.ReleaseMode)
    router := gin.Default()

    router.LoadHTMLGlob("templates/*templ.html")
    router.GET("/", rocket)
	router.GET("/socket", serverwebsocket)
    router.StaticFS("/static", http.Dir("static"))
    port := "8888"

    log.Println("Starting http server...")
    if err := router.Run(":" + port); err != nil {
        log.Panicf("error: %s", err)
    }
}

func rocket(c *gin.Context) {
	livecountdownvalue := live
	c.HTML(http.StatusOK, "rocket.templ.html", gin.H{
		"timestamp": livecountdownvalue,
	})
}

func serverwebsocket(c *gin.Context) {
		var w http.ResponseWriter = c.Writer
		var r *http.Request = c.Request

		 conn, err := upgrader.Upgrade(w, r, nil)
  		 if err != nil {
    		log.Println("upgrade failed: ", err)
  		 }

    	 if len(string(countdown)) != 0 {

   		 output := string(countdown)
   		 message := []byte(output)

   		 err = conn.WriteMessage(websocket.TextMessage, message)

   		 if err != nil {
    		 log.Println("write failed:", err)
   		 	 return
		 }
		}
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
			go WebSocketConnection()
		case <- quit:
			ticker.Stop()
			return
        }
     }
   }()
   wg.Wait()
   cancel()
}
