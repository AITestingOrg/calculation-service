package db

import (
	"gopkg.in/mgo.v2"
	"log"
	"github.com/AITestingOrg/calculation-service/models"
	"net/http"
	"encoding/json"
)

type TripDao struct {
	Server string
	Database string
}
var MgoSession *mgo.Session

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	log.Printf( "{message: %q\nheader: }", message, w)
}

// Establish connection with db
func init() {
	session, err := mgo.Dial("mongo:27017")

	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}
	session.SetMode(mgo.Monotonic, true)
	MgoSession = session
}

// Find a trip by it id
//func (m *TripDao) FindById(id string) (models.Cost, error) {
//	var cost Cost
//	err := MgoSession.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&cost)
//	return trip, err
//}

// Insert trip into database
func addTrip(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Clone()

		var cost models.Cost
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&cost)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body",  http.StatusBadRequest)
			return
		}

		c := session.DB("costData").C("costs")

		err = c.Insert(cost)
		if err != nil {
			if mgo.IsDup(err) {
				ErrorWithJSON(w, "Book with ISBN already exists", http.StatusBadRequest)
				return
			}
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed insert book: ", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path+"/"+cost.UserId)
		w.WriteHeader(http.StatusCreated)
	}
}