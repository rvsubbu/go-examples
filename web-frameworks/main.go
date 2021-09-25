package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi"
	"github.com/gorilla/mux"
	"github.com/labstack/echo"
)

// This is the expected body for req
type RequestBody struct {
	AI string `json:"action_item"`
}

// This is the response struct that will be serialized and sent back
type StatusResponse struct {
	Status string `json:"status"`
	AI     string `json:"action_item"`
}

func writeNewline(w http.ResponseWriter) {
	_, _ = w.Write([]byte("\n"))
}

func chiAiGetHandler(w http.ResponseWriter, r *http.Request) {
	// Add Content-Type header to indicate JSON response
	w.Header().Set(
		"Content-Type", "application/json",
	)

	body := StatusResponse{
		Status: "Hello world from chi!",
		AI:     chi.URLParam(r, "ai"),
	}

	serializedBody, _ := json.Marshal(body)
	_, _ = w.Write(serializedBody)
	writeNewline(w)
}

func chiAiPostHandler(w http.ResponseWriter, r *http.Request) {
	// Read complete request body
	rawRequestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Transform into RequestBody struct
	requestBody := &RequestBody{}
	err = json.Unmarshal(rawRequestBody, requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	body := StatusResponse{
		Status: "Hello world from chi!",
		AI:     requestBody.AI,
	}

	serializedBody, _ := json.Marshal(body)
	_, _ = w.Write(serializedBody)
	writeNewline(w)
}

func chiServer(wg sync.WaitGroup) {
	defer wg.Done()

	r := chi.NewRouter()

	r.Get("/ais/{ai}", chiAiGetHandler)
	r.Post("/ais", chiAiPostHandler)

	log.Fatal(http.ListenAndServe(":8001", r))
}

// In addition to echo request handlers using a special context including
// all kinds of utilities, generated errors can be returned to handle them easily
func echoAiGetHandler(e echo.Context) error {
	// Create response object
	body := &StatusResponse{
		Status: "Hello world from echo!",
		AI:     e.Param("ai"),
	}

	// In this case we can return the JSON function with our body as errors
	// thrown by this will be handled
	return e.JSON(http.StatusOK, body)
}

func echoAiPostHandler(e echo.Context) error {
	// Similar to the gin implementation,
	// we start off by creating an empty request body struct
	requestBody := &RequestBody{}

	// Bind body to the request body struct and check for potential errors
	err := e.Bind(requestBody)
	if err != nil {
		// If an error was created by the Bind operation, we can utilize echo's request
		// handler structure and simply return the error so it gets handled accordingly
		return err
	}

	body := &StatusResponse{
		Status: "Hello world from echo!",
		AI:     requestBody.AI,
	}

	return e.JSON(http.StatusOK, body)
}

func echoServer(wg sync.WaitGroup) {
	defer wg.Done()

	// Create echo instance
	e := echo.New()

	// Add endpoint route for /ais/<ai>
	e.GET("/ais/:ai", echoAiGetHandler)

	// Add endpoint route for /ais
	e.POST("/ais", echoAiPostHandler)

	// Start echo and handle errors
	e.Logger.Fatal(e.Start(":8002"))
}

func ginAiGetHandler(c *gin.Context) {
	// Create response object
	body := &StatusResponse{
		Status: "Hello world from gin!",
		AI:     c.Param("ai"),
	}

	// Send it off to the client
	c.JSON(http.StatusOK, body)
	c.Writer.WriteString("\n")
}

func ginAiPostHandler(c *gin.Context) {
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
		Status: "Hello world from gin!",
		AI:     requestBody.AI,
	}

	// And send it off to the requesting client
	c.JSON(http.StatusOK, body)
	c.Writer.WriteString("\n")
}

func ginServer(wg sync.WaitGroup) {
	defer wg.Done()

	// Create gin router
	router := gin.Default()

	// Set up GET endpoint // for route  /ai/<ai>
	router.GET("/ais/:ai", ginAiGetHandler)

	// Set up POST endpoint for route /ais
	router.POST("/ais", ginAiPostHandler)

	// Launch Gin and // handle potential error
	err := router.Run(":8003")
	if err != nil {
		panic(err)
	}
}

func muxAiGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")

	body := StatusResponse{
		Status: "Hello world from mux!",
		AI:     vars["action_item"],
	}

	serializedBody, _ := json.Marshal(body)
	_, _ = w.Write(serializedBody)
	writeNewline(w)
}

func muxAiPostHandler(w http.ResponseWriter, r *http.Request) {
	// Read complete request body
	rawRequestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Transform into RequestBody struct
	requestBody := &RequestBody{}
	err = json.Unmarshal(rawRequestBody, requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	body := StatusResponse{
		Status: "Hello world from mux!",
		AI:     requestBody.AI,
	}

	serializedBody, _ := json.Marshal(body)
	_, _ = w.Write(serializedBody)
	writeNewline(w)
}

func muxServer(wg sync.WaitGroup) {
	r := mux.NewRouter()

	r.HandleFunc("/ais/{ai}", muxAiGetHandler).Methods("GET")
	r.HandleFunc("/ais", muxAiPostHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8004", r))
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	log.Println("chi Server listening on :8001")
	go chiServer(wg) // Running on port 8001

	wg.Add(1)
	log.Println("echo Server listening on :8002")
	go echoServer(wg) // Running on port 8002

	wg.Add(1)
	log.Println("gin Server listening on :8003")
	go ginServer(wg) // Running on port 8003

	wg.Add(1)
	log.Println("mux Server listening on :8004")
	go muxServer(wg) // Running on port 8004

	wg.Wait()
}
