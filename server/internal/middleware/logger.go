package middleware

import (
	"log"
	"net/http"
	"strings"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		connectVersion := r.Header.Get("Connect-Protocol-Version")

		// 2. Check the Content-Type to distinguish between gRPC and gRPC-Web
		contentType := r.Header.Get("Content-Type")

		if connectVersion != "" {
			log.Printf("Protocol Used: Connect (Version %s)", connectVersion)
		} else if contentType == "application/grpc" {
			log.Println("Protocol Used: Standard gRPC")
		} else if contentType == "application/grpc-web" {
			log.Println("Protocol Used: gRPC-Web")
		} else {
			log.Printf("Protocol Used: Unknown (Content-Type: %s)", contentType)
		}

		if strings.HasPrefix(contentType, "application/grpc") {
			log.Println("Protocol Used: Standard gRPC (Binary)")
		} else if strings.Contains(contentType, "application/connect") {
			log.Println("Protocol Used: Connect Protocol")
		}

		// This captures the raw HTTP version (e.g., "HTTP/1.1" or "HTTP/2.0")
		log.Printf("HTTP_VER: [%s] | METHOD: [%s] | PATH: [%s]", r.Proto, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
