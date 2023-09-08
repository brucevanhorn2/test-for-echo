package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "sync"
)

// GetANSIColorCode returns the ANSI color code based on the HTTP method.
func GetANSIColorCode(method string) string {
    switch method {
    case "GET":
        return "\x1b[32m" // Green
    case "POST":
        return "\x1b[33m" // Yellow
    case "DELETE":
        return "\x1b[31m" // Red
    case "PUT":
        return "\x1b[36m" // Cyan
    default:
        return "\x1b[37m" // White
    }
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
    // Read body
    bodyBytes, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error reading body", http.StatusBadRequest)
        return
    }

    // Prepare response data
    responseData := map[string]interface{}{
        "Method": r.Method,
        "Header": r.Header,
        "URL":    r.URL.String(),
        "Body":   string(bodyBytes),
    }

    // Determine the color for the method
    methodColor := GetANSIColorCode(r.Method)
    gray := "\x1b[90m" // Gray
    reset := "\x1b[0m"  // Reset

    // Log the request data to stdout
    log.Printf("%sReceived request: map[Body: %s Header: %s Method: %s%s%s URL: %s%s%s]%s\n",
        gray, string(bodyBytes), r.Header, methodColor, r.Method, reset, reset, r.URL.String(), gray, reset)

    // Encode response data as JSON
    jsonResponse, err := json.Marshal(responseData)
    if err != nil {
        http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
        return
    }

    // Set Content-Type and write the response
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}

func startServer(port int) {
    mux := http.NewServeMux()
    mux.HandleFunc("/", echoHandler)

    address := fmt.Sprintf(":%d", port)
    server := &http.Server{
        Addr:    address,
        Handler: mux,
    }

    fmt.Printf("Echo server is running on port %d\n", port)
    log.Fatal(server.ListenAndServe())
}

func main() {
    var wg sync.WaitGroup

    ports := []int{4000, 5000, 6001}

    for _, port := range ports {
        wg.Add(1)
        go func(p int) {
            defer wg.Done()
            startServer(p)
        }(port)
    }

    wg.Wait()
}
