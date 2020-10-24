package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	var port = flag.String("p", "80", "Port number")
	var logFile = flag.Bool("l", false, "Create log file or not(default false)")
	flag.Parse()

	fs := http.FileServer(http.Dir("."))

	ip := GetLocalIP()
	fmt.Printf("Listening %s:%s\n", ip, *port)
	log.Fatal(http.ListenAndServe(":"+*port, LoggingHandler(fs, *logFile)))
}

// LoggingHandler write log to file and stdout
func LoggingHandler(next http.Handler, logFile bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		ip := GetIP(r)

		// Don't create a log file
		if logFile == false {
			log.Println(ip, r.Method, r.URL.Path)
			return
		}

		// write into a file
		f, err := os.OpenFile("fs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer f.Close()

		// only write log to file
		// log.SetOutput(f)

		// write log to file and stdout
		wrt := io.MultiWriter(os.Stdout, f)
		log.SetOutput(wrt)
		log.Println(ip, r.Method, r.URL.Path)
	})
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

// GetLocalIP get primary IP address
func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Println(err)
		return ""
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
