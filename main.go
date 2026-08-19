package main

import (
	"log"
	
	"net/http"
	
)




func main() {
	const port = "8080"
	mux := http.NewServeMux()
	
	ser := &http.Server{
		Addr: ":"+port,
		Handler: mux,
	}

	log.Fatal(ser.ListenAndServe())
}