package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

// For measuring high CPU
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}

	return fibonacci(n-1) + fibonacci(n-2)
}

func (app *application) fibonacciHandler(w http.ResponseWriter, r *http.Request) {
	num, err := strconv.Atoi(r.PathValue("num"))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if num >= 50 {
		app.badRequestResponse(w, r, errors.New("numbers larger than 50 result in timeouts"))
	}
	start := time.Now()
	result := fibonacci(num)
	end := time.Now()
	timeToCalculate := end.Sub(start)
	response := map[string]string{
		"result":          strconv.Itoa(result),
		"timeToCalculate": timeToCalculate.String(),
		"num":             r.PathValue("num"),
	}
	writeJSON(w, 200, response)
}
