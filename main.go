package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"io"

	"github.com/gorilla/mux"
)

type Event struct {
    Name        string  `json:"id"`
    Title       string  `json:"title"`
    Location    string  `json:"location"`
    When        string  `json:"when"`
}

func main() {
	// This just adds some events for testing.
	//eventsdb.InitializeEventsArray()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}
	log.Printf("Port set to: %s", port)

	fs := http.FileServer(http.Dir("assets"))
	myRouter := mux.NewRouter().StrictSlash(true)

	// This serves the static files in the assets folder
	myRouter.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// The rest of the routes
	myRouter.HandleFunc("/", indexHandler)
	myRouter.HandleFunc("/about", aboutHandler)
	myRouter.HandleFunc("/add", addHandler)
	myRouter.HandleFunc("/edit/{id}", editHandler)
	myRouter.HandleFunc("/delete/{id}", deleteHandler)

	log.Printf("Webserver listening on Port: %s", port)
	http.ListenAndServe(":"+port, myRouter)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(os.Getenv("EVENTS_API_URL") + "/events")
	events := []Event{}
	if err != nil {
	    log.Println(err.Error())
	} else {
        defer resp.Body.Close()
        bodyBytes, _ := io.ReadAll(resp.Body)
        json.Unmarshal(bodyBytes, &events)
    }

	data := HomePageData{
		PageTitle: "Home Page",
		Events:    events,
		Count:     len(events),
	}

	var tpl = template.Must(template.ParseFiles("templates/index.html", "templates/layout.html"))

	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	buf.WriteTo(w)
	log.Println("Home Page Served")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := AboutPageData{
		PageTitle: "About Go Website",
	}

	var tpl = template.Must(template.ParseFiles("templates/about.html", "templates/layout.html"))

	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	buf.WriteTo(w)
	log.Println("About Page Served")
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := AddPageData{
			PageTitle: "Add Event",
		}

		var tpl = template.Must(template.ParseFiles("templates/add.html", "templates/layout.html"))

		buf := &bytes.Buffer{}
		err := tpl.Execute(buf, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		buf.WriteTo(w)

		log.Println("Add Page Served")
	} else {
		// Add Event Here
		event := Event{
			Title:    r.FormValue("title"),
			Location: r.FormValue("location"),
			When:     r.FormValue("when"),
		}
		eventBody, _ := json.Marshal(event)
		requestBody := bytes.NewBuffer(eventBody)
		_, err := http.Post(os.Getenv("EVENTS_API_URL") + "/events", "application/json", requestBody)
		if err != nil {
		    log.Println(err.Error())
		}

		// Go back to home page
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		log.Println("Edit Handler")
		id := mux.Vars(r)["id"]
	    resp, err := http.Get(os.Getenv("EVENTS_API_URL") + "/events/" + id)
	    event := Event{}
	    if err != nil {
	        log.Println(err.Error())
	    } else {
            defer resp.Body.Close()
            bodyBytes, _ := io.ReadAll(resp.Body)
            json.Unmarshal(bodyBytes, &event)
        }

		data := EditPageData{
			PageTitle: "Edit Event",
			Event:     event,
		}

		var tpl = template.Must(template.ParseFiles("templates/edit.html", "templates/layout.html"))

		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		buf.WriteTo(w)

		log.Println("Edit Page Served")
	} else {
		// Add Event Here
		event := Event{
			Name:     r.FormValue("id"),
			Title:    r.FormValue("title"),
			Location: r.FormValue("location"),
			When:     r.FormValue("when"),
		}
		eventBody, _ := json.Marshal(event)
   		requestBody := bytes.NewBuffer(eventBody)
   		client := &http.Client{}
     	req, err := http.NewRequest("PUT", os.Getenv("EVENTS_API_URL") + "/events/" + event.Name, requestBody)
   		if err != nil {
            log.Println(err.Error())
       	}
       	_, err = client.Do(req)
       	if err != nil {
       	    log.Println(err.Error())
       	}
		log.Println("Event Updated")

		// Go back to home page
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", os.Getenv("EVENTS_API_URL") + "/events/" + id, nil)
	if err != nil {
	    log.Println(err.Error())
	}
	_, err = client.Do(req)
	if err != nil {
	    log.Println(err.Error())
	}
	log.Println("Event deleted")

	// Go back to home page
	http.Redirect(w, r, "/", http.StatusFound)
}

// HomePageData for Index template
type HomePageData struct {
	PageTitle string
	Events    []Event
	Count     int
}

// AboutPageData for About template
type AboutPageData struct {
	PageTitle string
}

// AddPageData for About template
type AddPageData struct {
	PageTitle string
}

// EditPageData for About template
type EditPageData struct {
	PageTitle string
	Event     Event
}
