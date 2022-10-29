package main

import (
	"log"
	"net/http"
	_ "see-log-go/router"
)

func main() {
	log.Println("日志查看项目启动:端口3000")
	log.Fatalln("error", http.ListenAndServe(":3000", nil))
}
