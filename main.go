package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"

	"crisco/sse"
	"crisco/utils"
	"net/http"
)

// Temporary storage, until MongoDB implemented
var urlMap = make(map[string]string)

type urlToShorten struct {
	URL string `json:"url"`
}

type clickEvent struct {
	URL   string `json:"url"`
	Agent string `json:"agent"`
	Time  int64  `json:"time"`
}

func shorten(c *gin.Context) {
	// Generate a map value
	route := utils.RandomString(6)

	// Bind JSON data to struct
	var mapData urlToShorten
	c.BindJSON(&mapData)

	// Build the new URL
	newURL := fmt.Sprintf("%s/%s", location.Get(c), route)

	// Assign mred value in map
	urlMap[route] = mapData.URL
	c.JSON(http.StatusOK, gin.H{"url": newURL})
}

func follow(c *gin.Context) {
	route := c.Param("route")

	dest, ok := urlMap[route]
	if ok {
		// Forward to page
		c.Redirect(http.StatusTemporaryRedirect, urlMap[route])

		msg, err := json.Marshal(clickEvent{
			route,
			c.Request.Header.Get("User-Agent"),
			utils.GetTime(),
		})
		if err != nil {
			sse.URL(route).Submit("{ \"error\": \"Unable to generate message for click\"}")
			return
		}

		sse.URL(route).Submit(string(msg))
	}
}

func subscribe(c *gin.Context) {
	route := c.Param("route")

	_, ok := urlMap[route]

	if !ok {
		c.String(http.StatusNotFound, "Stream not available")
		return
	}

	fmt.Printf("Listener subscribing to: %s\n", route)

	// Open a listener for the route
	listener := sse.OpenListener(route)
	defer sse.CloseListener(route, listener)

	// Stream events to the listener
	c.Stream(func(w io.Writer) bool {
		c.SSEvent("click", <-listener)
		return true
	})
}

func main() {
	r := gin.Default()

	// Used for grabbing the URL origin
	r.Use(location.Default())

	// Set up CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	r.Use(cors.New(config))

	// Routes
	r.GET("/:route", follow)
	r.POST("/shorten", shorten)
	r.GET("/:route/events", subscribe)

	// listen and serve on 0.0.0.0:8080
	// Port can be set with PORT environment variable
	r.Run()
}
