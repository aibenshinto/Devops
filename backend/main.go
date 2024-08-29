package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    // Connect to MongoDB
    var err error
    client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongo:27017"))
    if err != nil {
        log.Fatal(err)
    }
    if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB!")

    // HTTP Handlers
    http.HandleFunc("/register", registerHandler)
    http.HandleFunc("/view", viewHandler)
    http.Handle("/", http.FileServer(http.Dir("./frontend")))

    // Start Server
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    collection := client.Database("mydb").Collection("users")
    _, err = collection.InsertOne(context.TODO(), user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "User registered successfully!")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
    collection := client.Database("mydb").Collection("users")
    cur, err := collection.Find(context.TODO(), bson.D{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cur.Close(context.TODO())

    var users []User
    for cur.Next(context.TODO()) {
        var user User
        err := cur.Decode(&user)
        if err != nil {
            log.Fatal(err)
        }
        users = append(users, user)
    }

    json.NewEncoder(w).Encode(users)
}
