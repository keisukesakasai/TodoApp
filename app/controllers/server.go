package controllers

import (
	"TodoApp/config"
	"fmt"
	"net/http"
)

func StartMainServer() error {
	fmt.Println("start server" + "port: " + config.Config.Port)
	files := http.FileServer((http.Dir(config.Config.Static)))
	http.Handle("/static/", http.StripPrefix("/static/", files))

	http.HandleFunc("/", top)
	return http.ListenAndServe(":"+config.Config.Port, nil)
}
