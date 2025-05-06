package main

import (
	"fmt"
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

// fibonacciHandler runs the Fibonacci sequence for the provided num route parameter.
// It emulates a high CPU intensive routes.
func (app *application) fibonacciHandler(w http.ResponseWriter, r *http.Request) {
	num, err := strconv.Atoi(r.PathValue("num"))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if num >= 50 {
		app.badRequestResponse(w, r, fmt.Errorf("numbers >= 50 result in timeouts: %d", num))
		return
	}

	start := time.Now()
	result := fibonacci(num)
	end := time.Now()
	timeToCalculate := end.Sub(start)

	response := struct {
		Num             int
		Result          int
		TimeToCalculate string
	}{
		Num:             num,
		Result:          result,
		TimeToCalculate: timeToCalculate.String(),
	}

	writeJSON(w, 200, response)
}
