package loggermiddleware

import (
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		rangeStatus := checkStatusCode(lrw.statusCode)
		coloredLog := func(col color.Attribute) {
			colorFunc := color.New(col).SprintFunc()
			log.Printf("%s %s %s %vms - %s", colorFunc(r.Method), colorFunc(lrw.statusCode), r.URL, time.Since(start).Milliseconds(), w.Header().Get("Content-Length"))
		}

		switch rangeStatus {
		case success:
			coloredLog(color.FgGreen)
		case redirect:
			coloredLog(color.FgBlue)
		case clientError:
			coloredLog(color.FgRed)
		case serverError:
			coloredLog(color.BgYellow)
		}

	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

type rangeStatusCode int

const (
	info rangeStatusCode = iota
	success
	redirect
	clientError
	serverError
)

func checkStatusCode(statusCode int) rangeStatusCode {
	switch {
	case statusCode >= 200 && statusCode <= 299:
		return success
	case statusCode >= 300 && statusCode <= 399:
		return redirect
	case statusCode >= 400 && statusCode <= 499:
		return clientError
	case statusCode >= 500 && statusCode <= 599:
		return serverError
	default:
		return info
	}

}
