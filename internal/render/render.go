package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/RistoFlink/bookings/internal/config"
	"github.com/RistoFlink/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}
var app *config.AppConfig
var pathToTemplates = "./templates"

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds a data to all the pages if necessary
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders templates using "html/template"
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	// create a template cache
	// get the template cache from the app config
	// if UseCache is true -> read the information from the template cache. If not, rebuild the template cache
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

	// old way to parse templates, replaced by what's above
	//parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")
	//err := parsedTemplate.Execute(w, nil)
	//if err != nil {
	//	fmt.Println("error parsing template: ", err)
	//	return
	//}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all files named *page.tmpl from the templates folder ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// range through all found files ending in *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}

// a simple way to make a template cache, commented out to build out the more complex way of doing it
//var tc = make(map[string]*template.Template)
//
//func RenderTemplate(w http.ResponseWriter, t string) {
//	var tmpl *template.Template
//	var err error
//
//	// check to see if the template is already in the cache
//	_, inMap := tc[t]
//	if !inMap {
//		// create the template if not in the cache
//		log.Println("creating template and adding to cache")
//		err = createTemplateCache(t)
//
//		if err != nil {
//			log.Println(err)
//		}
//	} else {
//		// template in cache already
//		log.Println("using template in cache")
//	}
//
//	tmpl = tc[t]
//	err = tmpl.Execute(w, nil)
//
//	if err != nil {
//		log.Println(err)
//	}
//}
//
//func createTemplateCache(t string) error {
//	templates := []string{
//		fmt.Sprintf("./templates/%s", t), "./templates/base.layout.tmpl",
//	}
//	// parse the template
//	// ... takes each entry in a slice and inserts them as individual strings
//	tmpl, err := template.ParseFiles(templates...)
//
//	if err != nil {
//		return err
//	}
//
//	// add template to the cache
//	tc[t] = tmpl
//
//	return nil
//}
