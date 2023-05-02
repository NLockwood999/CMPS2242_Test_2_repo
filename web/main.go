//filename: main.go

package main

import (
	"log"
	"net/http"
)

//middleware creation
func middlewareA(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//this is executing on the way down
		log.Println("Executing middleware A")
		next.ServeHTTP(w, r)
		//this is execting on the way up
		log.Println("Executing maddleware A again")
	})
}

//middleware creation
func middlewareB(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//this is executing on the way down
		log.Println("Executing middleware B")
		if r.URL.Path == "/cherry" {
			return
		}
		next.ServeHTTP(w, r)
		//this is execting on the way up
		log.Println("Executing middleware B again")
	})
}

func ourHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing the custom handler .............")
	w.Write([]byte("carrots"))
}

//create the handler function
/*func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello I'm under the water\n"))
	w.Write([]byte("please help me\n"))
}*/

func main() {
	mux := http.NewServeMux()
	mux.Handle("/check", middlewareA(middlewareB(http.HandlerFunc(ourHandler))))
	//             key value

	log.Print("starting server on 4888")
	err := http.ListenAndServe(":4888", mux)
	log.Fatal(err)
}

// client->   req   web server->   middlware->   router->    handler
// client<-   req   <-web server   <-middlware   <-router    handler
