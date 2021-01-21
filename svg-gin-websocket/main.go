package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const tmpl = `
<circle id=circle-{{ .N }} cx="{{ .X }}" cy="{{ .Y }}" r="{{ .R }}" style="fill:{{ .COLOR }};stroke:black" />
`

var colors = []string{"yellow", "red", "blue", "green", "magenta"}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func processApi(c *gin.Context) {
	//Upgrade get request to webSocket protocol
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("error get connection")
		log.Fatal(err)
	}
	defer ws.Close()

	mt, message, err := ws.ReadMessage()
	if err != nil {
		log.Println("error read message")
		log.Fatal(err)
	}

	n := 1
	for {
		circle := generateCircle(n)
		message = []byte(string(circle))
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Println("error write message: " + err.Error())
		}
		dur, _ := time.ParseDuration(fmt.Sprintf("%ds", (rand.Intn(3-1) + 1)))
		time.Sleep(dur)
		// time.Sleep(3 * time.Second)
		n++
	}
}

func main() {
	r := gin.Default()
	r.GET("/text", processApi)

	// static files
	r.Use(static.Serve("/", static.LocalFile("./public", true)))

	r.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	log.Println("... Started server on port: 8000")
	r.Run(":8000")
}

func generateCircle(n int) string {
	x := rand.Intn(20-10) + 10
	y := rand.Intn(11-6) + 6
	r := rand.Intn(6-2) + 2
	c := rand.Intn(4)
	vars := map[string]interface{}{
		"N":     n,
		"X":     x,
		"Y":     y,
		"R":     r,
		"COLOR": colors[c],
	}
	s, err := createTemplate(tmpl, vars)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(s)
	return s
}

func createTemplate(body string, vars interface{}) (string, error) {
	tmpl, err := template.New("template").Parse(body)
	if err != nil {
		return "", nil
	}
	return process(tmpl, vars)
}

func process(t *template.Template, vars interface{}) (string, error) {
	var b bytes.Buffer
	err := t.Execute(&b, vars)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
