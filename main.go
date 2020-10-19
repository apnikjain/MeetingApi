package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Meeting struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title        string             `bson:"title" json:"title"`
	Created      string             `bson:"created" json:"created"`
	Start        time.Time          `bson:"start" json:"start"`
	End          time.Time          `bson:"end" json:"end"`
	Participants []Participant      `bson:"participants" json:"participants"`
}

type Participant struct {
	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	RSVP  string `bson:"rsvp" json:"rsvp"`
}

func CreateMeeting(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	now := time.Now()
	var meeting Meeting
	meeting.Created = now.Format("Mon Jan _2 15:04:05 2006")
	_ = json.NewDecoder(request.Body).Decode(&meeting)
	collection1 := client.Database("test3").Collection("meetings")
	ctx1, _ := context.WithTimeout(context.Background(), 5*time.Second)

	collection := client.Database("test3").Collection("meetings")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	fmt.Println(len(meeting.Participants))
	for i := 0; i < len(meeting.Participants); i++ {
		fmt.Println("sdf")
		for cursor.Next(ctx) {
			var meetingg Meeting
			cursor.Decode(&meetingg)
			start := meetingg.Start
			end := meetingg.End
			in_start := meeting.Start
			in_end := meeting.End
			var startOverlap = inTimeSpan(start, end, in_start)
			var endOverlap = inTimeSpan(start, end, in_end)
			fmt.Println(startOverlap, endOverlap)
			if startOverlap || endOverlap {
				for j := 0; j < len(meetingg.Participants); j++ {
					if (meetingg.Participants[j].Email == meeting.Participants[i].Email) && meetingg.Participants[j].RSVP == "Yes" && meeting.Participants[i].RSVP == "Yes" {
						response.WriteHeader(http.StatusInternalServerError)
						response.Write([]byte(`{ "message: Overlaping" }`))
						return
					}
				}
			}

		}
		if err := cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
	}

	result, _ := collection1.InsertOne(ctx1, meeting)
	json.NewEncoder(response).Encode(result)

}

func inTimeSpan(start, end, check time.Time) bool {
	return (check.After(start) && check.Before(end)) || check.Equal(start) || check.Equal(end)
}

func GetMeetingById(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var meetArr []Meeting
	collection := client.Database("test3").Collection("meetings")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var meeting Meeting
		cursor.Decode(&meeting)
		if meeting.ID == id {
			meetArr = append(meetArr, meeting)
		}
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meetArr)
}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var meetArr []Meeting
	collection := client.Database("test3").Collection("meetings")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var meeting Meeting
		cursor.Decode(&meeting)
		meetArr = append(meetArr, meeting)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meetArr)
}

func GetMeetingByParticipant(response http.ResponseWriter, request *http.Request) {

	mailId := request.FormValue("participant")

	start := request.FormValue("start")
	end := request.FormValue("end")

	if mailId != "" {
		response.Header().Set("content-type", "application/json")

		var meetArr []Meeting
		collection := client.Database("test3").Collection("meetings")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var meeting Meeting
			cursor.Decode(&meeting)
			for i := 0; i < len(meeting.Participants); i++ {
				if meeting.Participants[i].Email == mailId {
					meetArr = append(meetArr, meeting)
				}
			}

		}
		if err := cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(response).Encode(meetArr)
	}

	if start != "" && end != "" {
		response.Header().Set("content-type", "application/json")

		var meetArr []Meeting
		collection := client.Database("test3").Collection("meetings")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var meeting Meeting
			cursor.Decode(&meeting)
			var start_time, _ = time.Parse(time.RFC3339, start)
			var end_time, _ = time.Parse(time.RFC3339, end)
			var afterStart = start_time.Before(meeting.Start)
			var beforeEnd = end_time.Before(meeting.End)
			fmt.Println(afterStart, start_time, beforeEnd)
			if afterStart && beforeEnd {
				meetArr = append(meetArr, meeting)
			}

		}
		if err := cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(response).Encode(meetArr)
	}

}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/meetings", CreateMeeting).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/meetings/{id}", GetMeetingById).Methods("GET")
	router.HandleFunc("/meetings", GetMeetingByParticipant).Methods("GET")
	http.ListenAndServe(":12345", router)
}
