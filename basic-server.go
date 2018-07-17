package main

import (
  "log"
  "net/http"
  "io/ioutil"
  "html/template"
  "regexp"
  "fmt"
)

//TODO:Add a handler to make the web root redirect to /view/FrontPage.
//TODO:Spruce up the page templates by making them valid HTML and adding some CSS rules.
//TODO:Implement inter-page linking by converting instances of [PageName] to
//TODO:<a href="/view/PageName">PageName</a>. (hint: you could use regexp.ReplaceAllFunc to do this)


// the struct that we use to hold all the data for pages: title and body
type Page struct {
  Title string
  Body []byte
}

// saves a page (writes to a file in the data subdirectory)
func (p *Page) save() error {
  // includes data/ prefix to nav to data directory
  filename := "data/" + p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body, 0600)
}

// loads a page
func loadPage(title string) (*Page, error) {
  // includes the data/ prefix to navigate to the data directory
  filename := "data/" + title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

// // THIS PROBSBLY doesn't work
// //Supposed to take the title string (which should always be empty) and return
// // the homepage (which should be in the data directory tbh)
// func rootHandler(w http.ResponseWriter, r *http.Request, title string) {
//   fmt.Println("They got to rootHandler")
// }

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    return
  }
  renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  err := p.save()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}


// if we add any templates we add them to this call to ParseFiles
// We must remember to add tmpl/ as a prefix to all filenames cause that's where
// they are being stored
var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  // don't have to add data prefix cause this is about which, view or edit?
  err := templates.ExecuteTemplate(w, tmpl+".html", p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

// need to include an empty / into this logic, then it'll work
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
      // so if the path isn't valid but ALSO the path is nothing,
      // then redirects to homepage, should do through an if below
      // ie if r = "" then redirect to view/homepage.html or something
      http.NotFound(w, r)
      return
    }
    fn(w, r, m[2])
  }
}

func main() {
  fmt.Println("The program started")
  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))


  fmt.Println("the program is about to start litening on 8080")
  log.Fatal(http.ListenAndServe(":0808", nil))
  fmt.Println("the program is ending after this")
}
