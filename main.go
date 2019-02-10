package main

import(
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"

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
	t.templ.Execute(w,nil)
}


func main(){
	// We use the handle function that takes in a custom handler
	http.Handle("/", &templateHandler{filename:"chat.html"})
	//template handler is a valid http.Handler because of the ServerHttp

	if err:= http.ListenAndServe(":8080",nil); err !=nil{
		log.Fatal("ListenAndServe:",err)
	}
}