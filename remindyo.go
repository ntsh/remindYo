package remindYo

import (
    "fmt"
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
    fmt.Fprint(w,user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)    
    }

    c := appengine.NewContext(r)
    rem := Reminder{
                UserName: user,
                TimeStamp: time.Now(),
                Delivered: 0,
        }

    key := datastore.NewIncompleteKey(c, "Reminder", nil)
    _, err = datastore.Put(c, key, &rem)
    if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

}