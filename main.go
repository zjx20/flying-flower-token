package main

import (
	"net/http"

	handler "github.com/zjx20/flying-flower-token/api"
)

func main() {
	// 我们用 HandleFunc 方法来注册一个路由和一个处理函数
	// 这个函数的签名是：接收一个路径和一个函数作为参数
	// 这个函数必须有合适的签名（如下面的 handler 函数所示）
	http.HandleFunc("/", handler.Handler)

	// 在定义好我们的服务器后，我们就在 8080 端口上监听和服务
	// 第二个参数是处理器，我们暂时留空，用上面定义的处理函数
	http.ListenAndServe(":8080", nil)
}
