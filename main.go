package main

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", getIPHandler)
	http.ListenAndServe(":8080", nil)
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

func getClientIP(r *http.Request) string {
	// 检查代理头（适用于反向代理情况）
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// 可能有多个IP，取第一个
		ips := strings.Split(forwarded, ", ")
		if len(ips) > 0 {
			return strings.Split(ips[0], ":")[0] // 去掉端口号
		}
	}

	// 直接连接的情况
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // 如果分割失败返回原始值
	}
	return ip
}
