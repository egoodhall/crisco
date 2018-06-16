package main

import (
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

// Structure of incoming data for URL shortening
type urlToShorten struct {
	URL string `json:"url"`
}

// Structure of events that are broadcast to
// listeners via SSE's
type clickEvent struct {
	URL   string `json:"url"`
	Agent string `json:"agent"`
	Time  int64  `json:"time"`
}

// Shorten the url supplied in the request
// body's JSON. Will respond with a shortened
// URL that can be used in place of the original
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

// Follow a shortened URL to the destination
// held within the mapped values
func follow(c *gin.Context) {
	route := c.Param("route")

	// Handle forwarding and errors with redirecting
	dest, ok := urlMap[route]
	if ok {
		// Forward to page
		c.Redirect(http.StatusTemporaryRedirect, dest)

		// Generate click event only if there's an existing channel to put it on
		channel, ok := sse.RouteChannels[route]
		if ok {
			agent := c.Request.Header.Get("User-Agent")
			channel.Submit(gin.H{"route": route, "agent": agent, "time": utils.GetTime()})
		}
	} else {
		c.String(http.StatusNotFound, "Unknown route /%s", route)
	}
}

// Subscribe to a broadcaster that dispatches Server-Sent
// Events on clicks of the shortened link.
func subscribe(c *gin.Context) {
	route := c.Param("route")

	// Verify route exists
	_, ok := urlMap[route]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stream not available"})
		return
	}

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

	// Set up CORS (for testing with a React App)
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
