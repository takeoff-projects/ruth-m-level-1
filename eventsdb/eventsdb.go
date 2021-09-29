package eventsdb

import (
	"errors"
    "context"
	"github.com/google/uuid"
	"cloud.google.com/go/datastore"
	"log"
	"os"
)

// Event model
type Event struct {
    Name     string `json:"id"`
    Title  string `json:"title" datastore:"title"`
    Location string `json:"location" datastore:"location"`
    When   string `json:"when" datastore:"when"`
}

//var Events []Event
 
func GetEvents() []Event {
    ctx := context.Background()
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    dsClient, err := datastore.NewClient(ctx, projectID)
    if err != nil  {
        log.Println(err.Error())
        return []Event{}
    }

    var events[]Event
    _, err = dsClient.GetAll(ctx, datastore.NewQuery("Event"), &events)
    if err != nil  {
        log.Println(err.Error())
    }

    return events
}

func InitializeEventsArray(){
    ctx := context.Background()
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    dsClient, err := datastore.NewClient(ctx, projectID)
    if err != nil  {
        log.Println(err.Error())
        return
    }

    events := []*Event{
        {Name: "2944a9cb-ef2d-4632-ac1d-af2b2629d0f2",
         Title: "Dinner",
         Location: "My House",
         When: "Tonight"},
         {Name: "f88f1860-9a5d-423e-820f-9acb4db3030e",
          Title: "Go Programming Lesson",
          Location: "At School",
          When: "Tomorrow"},
         {Name: "4cb393fb-dd19-469e-a52c-22a12c0a98df",
          Title: "Company Picnic",
          Location: "At the Park",
          When: "Saturday"},
    }

    keys := []*datastore.Key{
        datastore.NameKey("Event", events[0].Name, nil),
        datastore.NameKey("Event", events[1].Name, nil),
        datastore.NameKey("Event", events[2].Name, nil),
    }

    _, err = dsClient.PutMulti(ctx, keys, events)
    if err != nil  {
        log.Println(err.Error())
    }
}

func GetEventbyID(key string) (Event, error) {
    ctx := context.Background()
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    dsClient, err := datastore.NewClient(ctx, projectID)
    if err != nil  {
        return Event{}, errors.New(err.Error())
    }

    event := Event{}
    nameKey := datastore.NameKey("Event", key, nil)
    err = dsClient.Get(ctx, nameKey, &event)
    if err != nil  {
        return Event{}, errors.New(err.Error())
    }

    return event, nil
}

func AddEvent(event Event) {
    ctx := context.Background()
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    dsClient, err := datastore.NewClient(ctx, projectID)
    if err != nil  {
        log.Println(err.Error())
        return
    }

	newID := uuid.New().String()
	key := datastore.NameKey("Event", newID, nil)
	event.Name = newID
	_, err = dsClient.Put(ctx, key, &event)
    if err != nil  {
        log.Println(err.Error())
    }
}

func UpdateEvent(updatedEvent Event) {
    ctx := context.Background()
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    dsClient, err := datastore.NewClient(ctx, projectID)
    if err != nil  {
        log.Println(err.Error())
        return
    }

    key := datastore.NameKey("Event", updatedEvent.Name, nil)
    _, err = dsClient.Put(ctx, key, &updatedEvent)
    if err != nil  {
        log.Println(err.Error())
    }
}

func DeleteEvent(id string) {
    ctx := context.Background()
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    dsClient, err := datastore.NewClient(ctx, projectID)
    if err != nil  {
        log.Println(err.Error())
        return
    }

    key := datastore.NameKey("Event", id, nil)
    err = dsClient.Delete(ctx, key)
    if err != nil  {
        log.Println(err.Error())
    }
}