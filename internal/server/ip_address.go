package server

import (
	"log"
	"net"
	"net/http"
)

// CheckIP is chi middleware function used to decrypt the received body
func CheckIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if Config.TrustedSubnet != nil {
			xRealIP := r.Header.Get("X-Real-IP")
			ip := net.ParseIP(xRealIP)
			if !Config.TrustedSubnet.Contains(ip) {
				log.Print(
					"X-Real-IP is not set or address is not allowed:",
					xRealIP,
				)
				http.Error(rw, "IP not allowed", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(rw, r)
	})
}
