package remindYo

import (
    "net/http"
    "time"

    "appengine"
    "appengine/datastore"

    "github.com/ntsh/go-yo-gae"

)

type Reminder struct {
        UserName  string
        TimeStamp   time.Time
        Delivered int
}

func init() {
    http.HandleFunc("/getyo", getYoHandler)
    http.HandleFunc("/sendyo", sendYoHandler)
}

func getYoHandler(w http.ResponseWriter, r *http.Request) {
    user := r.URL.Query()["username"][0]

    c := appengine.NewContext(r)
    rem := Reminder{
                UserName: user,
                TimeStamp: time.Now(),
                Delivered: 0,
        }

    key := datastore.NewIncompleteKey(c, "Reminder", nil)
    _, err := datastore.Put(c, key, &rem)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func sendYoHandler(w http.ResponseWriter, r *http.Request) {
    client := getYoClient()
    var user string
    timeNow := time.Now()
    time1hr := timeNow.Add(-59*time.Minute)

    q := datastore.NewQuery("Reminder").Filter("Delivered =",0).Filter("TimeStamp <=", time1hr)

    var reminders []Reminder
    c := appengine.NewContext(r)
    keys, err := q.GetAll(c, &reminders)
    if  err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    for i := range reminders {
        user = reminders[i].UserName
        err := client.YoUser(user, r)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)    
        }
        reminders[i].Delivered = 1
    }
    datastore.PutMulti(c,keys,reminders)
}

func getYoClient() *yo.Client{
    return yo.NewClient(yo_api_token)
}