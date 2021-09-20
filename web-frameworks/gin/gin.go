package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// This is the response struct that will be
// serialized and sent back
type StatusResponse struct {
	Status     string `json:"status"`
	ActionItem string `json:"action_item"`
}

func ActionItemGetHandler(c *gin.Context) {
	// Create response object
	body := &StatusResponse{
		Status:     "Hello world from gin!",
		ActionItem: c.Param("ai"),
	}

	// Send it off to the client
	c.JSON(http.StatusOK, body)
}

type RequestBody struct {
	ActionItem string `json:"ai"`
}

func ActionItemPostHandler(c *gin.Context) {
	// Create empty request body struct used to bind actual body into
	requestBody := &RequestBody{}

	// Bind JSON content of request body to struct created above
	err := c.BindJSON(requestBody)
	if err != nil {
		// Gin automatically returns an error response when the BindJSON operation fails,
		// we simply have to stop this function from continuing to execute
		return
	}

	// Create response object
	body := &StatusResponse{
		Status:     "Hello world from echo!",
		ActionItem: requestBody.ActionItem,
	}

	// And send it off to the requesting client
	c.JSON(http.StatusOK, body)
}

func main() {
	// Create gin router
	router := gin.Default()

	// Set up GET endpoint // for route  /ai/<ai>
	router.GET("/ais/:ai", ActionItemGetHandler)

	// Set up POST endpoint for route /ais
	router.POST("/ais", ActionItemPostHandler)

	// Launch Gin and // handle potential error
	err := router.Run(":8003")
	if err != nil {
		panic(err)
	}
}
