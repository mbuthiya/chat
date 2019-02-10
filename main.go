package main

import(
	"flag"
	"log"
	"net/http"
	"text/template"
	"os"
	"path/filepath"
	"sync"
	"github.com/mbuthiya/tracer"

	// TODO: change package name to tracer


)


type templateHandler struct{
	once sync.Once
	filename string
	templ *template.Template // represents a single template
}


// this method will load the file compile the template and execute it
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	t.once.Do(func(){
		// Compiling the template file
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates",t.filename)))
	})

	// Executing the template file and write the output to the response writer method 
	t.templ.Execute(w,r)
}


func main(){

	// Create Address flag 
	var addr = flag.String("addr",":8080","The Address of the application")
	flag.Parse()

	r := newRoom() // Create a new room instance
	r.tracer = tracer.New(os.Stdout)


	// We use the handle function that takes in a custom handler
	http.Handle("/", &templateHandler{filename:"chat.html"})
	//template handler is a valid http.Handler because of the ServerHttp

	http.Handle("/room",r)

	// Get the room going
	go r.run()
	
	log.Println("Server is Running on",*addr)
	if err:= http.ListenAndServe(*addr,nil); err !=nil{
		log.Fatal("ListenAndServe:",err)
	}
}