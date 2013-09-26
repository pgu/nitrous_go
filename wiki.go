package main

import (
  "fmt"
  "html/template"
  "io/ioutil"
  "net/http"
  "regexp"
)

type Page struct {
  Title string
  Body  []byte
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

func (p *Page) save() error {
  filename := p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

const lenPath = len("/view/")

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  err := templates.ExecuteTemplate(w, tmpl + ".html", p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[lenPath:]
    if !titleValidator.MatchString(title) {
      http.NotFound(w, r)
      return
    }
    fn(w, r, title)
  }
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/edit/" + title, http.StatusFound)
    return
  }
  renderTemplate(w, "view", p)
  
//  fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  renderTemplate(w, "edit", p)
  
//  fmt.Fprintf(w, "<h1>Editing %s</h1>" +
//              "<form action=\"/save/%s\" method=\"POST\">" +
//              "<textarea name=\"body\">%s</textarea><br>" +
//              "<input type=\"submit\" value=\"Save\">" +
//              "</form>",
//              p.Title, p.Title, p.Body)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  err := p.save()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

//
// http://golang.org/doc/articles/wiki/
// go build wiki.go
// ./wiki
//
func main() {
  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))
  http.ListenAndServe(":8080", nil)
}

// hello world works! http://go-box-33817.euw1.actionbox.io:8080/view/test
func main_helloworld() { // server
  //http.HandleFunc("/view/", viewHandler)
  http.ListenAndServe(":8080", nil)
}
 
func main_page() { // test // create wiki pages
  p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
  p1.save()
  p2, _ := loadPage("TestPage")
  fmt.Println(string(p2.Body))
}


