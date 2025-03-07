package main
import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)
// 自定义ResponseWriter用于捕获状态码
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}
func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}
func main() {
	addr := ":18125"
	log.Printf("Starting server on %s", addr)
	http.Handle("/", loggingMiddleware(http.HandlerFunc(getIPHandler)))
	log.Fatal(http.ListenAndServe(addr, nil))
}
func getIPHandler(w http.ResponseWriter, r *http.Request) {
	// 获取客户端IP地址
	ip := getClientIP(r)

	// 检查查询参数
	format := r.URL.Query().Get("format")

	if format == "json" {
		// 返回JSON格式
		w.Header().Set("Content-Type", "application/json")
		response := struct {
			IP string `json:"ip"`
		}{
			IP: ip,
		}
		
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// 返回纯文本格式
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(ip))
	}
}
// 日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		defer func() {
			duration := time.Since(start)
			clientIP := getClientIP(r)
			timestamp := start.UTC().Add(8 * time.Hour).Format("02/Jan/2006:15:04:05 -0700")
			
			log.Printf(`%s - - [%s] "%s %s %s" %d %d`,
				clientIP,
				timestamp,
				r.Method,
				r.URL.RequestURI(),
				r.Proto,
				lw.statusCode,
				duration.Microseconds(),
			)
		}()
		next.ServeHTTP(lw, r)
	})
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ", ")
		if len(ips) > 0 {
			// 修复点：使用标准库处理IP和端口
			ipPort := strings.TrimSpace(ips[0])
			ipPart, _, err := net.SplitHostPort(ipPort)
			if err == nil {
				return ipPart
			}
			// 如果SplitHostPort失败，说明可能没有端口号，直接返回
			return ipPort
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
