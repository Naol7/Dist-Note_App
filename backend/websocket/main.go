package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins in development
		},
	}

	manager = &ConnectionManager{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan Message),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}

	client *mongo.Client
)

type ConnectionManager struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan Message
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.Mutex
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Note struct {
	ID      string `json:"id" bson:"id"`
	Content string `json:"content" bson:"content"`
	Version int    `json:"version" bson:"version"`
}

type Delta struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
	Start   int    `json:"start"`
	End     int    `json:"end"`
	Text    string `json:"text"`
}

func init() {
	var err error
	uri := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(uri)

	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")
}

func (manager *ConnectionManager) run() {
	for {
		select {
		case conn := <-manager.register:
			manager.mutex.Lock()
			manager.clients[conn] = true
			manager.mutex.Unlock()

			notes := fetchAllNotes()
			message := Message{
				Type:    "init",
				Payload: notes,
			}
			conn.WriteJSON(message)

		case conn := <-manager.unregister:
			manager.mutex.Lock()
			if _, ok := manager.clients[conn]; ok {
				delete(manager.clients, conn)
				conn.Close()
			}
			manager.mutex.Unlock()

		case message := <-manager.broadcast:
			manager.mutex.Lock()
			for conn := range manager.clients {
				err := conn.WriteJSON(message)
				if err != nil {
					log.Printf("Error broadcasting to client: %v", err)
					conn.Close()
					delete(manager.clients, conn)
				}
			}
			manager.mutex.Unlock()
		}
	}
}

func fetchAllNotes() []Note {
	collection := client.Database("notes_db").Collection("notes")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Error fetching notes: %v", err)
		return []Note{}
	}
	defer cursor.Close(context.TODO())

	var notes []Note
	if err = cursor.All(context.TODO(), &notes); err != nil {
		log.Printf("Error decoding notes: %v", err)
		return []Note{}
	}
	return notes
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
		return
	}

	manager.register <- conn

	go func() {
		defer func() {
			manager.unregister <- conn
		}()

		for {
			var message Message
			err := conn.ReadJSON(&message)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Error reading message: %v", err)
				}
				break
			}

			switch message.Type {
			case "create", "update", "delete":
				handleDatabaseOperation(message)
				manager.broadcast <- message
			case "delta":
				handleDeltaOperation(message)
			}
		}
	}()
}

func handleDatabaseOperation(message Message) {
	collection := client.Database("notes_db").Collection("notes")
	ctx := context.TODO()

	switch message.Type {
	case "create":
		if note, ok := message.Payload.(map[string]interface{}); ok {
			_, err := collection.InsertOne(ctx, note)
			if err != nil {
				log.Printf("Error creating note: %v", err)
			}
		}
	case "update":
		if note, ok := message.Payload.(map[string]interface{}); ok {
			_, err := collection.UpdateOne(
				ctx,
				bson.M{"id": note["id"]},
				bson.M{"$set": bson.M{"content": note["content"]}},
			)
			if err != nil {
				log.Printf("Error updating note: %v", err)
			}
		}
	case "delete":
		if note, ok := message.Payload.(map[string]interface{}); ok {
			_, err := collection.DeleteOne(ctx, bson.M{"id": note["id"]})
			if err != nil {
				log.Printf("Error deleting note: %v", err)
			}
		}
	}
}

func handleDeltaOperation(message Message) {
	collection := client.Database("notes_db").Collection("notes")
	ctx := context.TODO()

	if delta, ok := message.Payload.(Delta); ok {
		var note Note
		err := collection.FindOne(ctx, bson.M{"id": delta.ID}).Decode(&note)
		if err != nil {
			log.Printf("Error fetching note: %v", err)
			return
		}

		if note.Version != delta.Version {
			log.Printf("Version conflict for note %s", delta.ID)
			return
		}

		newContent := applyDelta(note.Content, delta)

		_, err = collection.UpdateOne(
			ctx,
			bson.M{"id": delta.ID},
			bson.M{"$set": bson.M{"content": newContent, "version": note.Version + 1}},
		)
		if err != nil {
			log.Printf("Error updating note: %v", err)
			return
		}

		manager.broadcast <- Message{Type: "delta", Payload: delta}
	}
}

func applyDelta(content string, delta Delta) string {
	return content[:delta.Start] + delta.Text + content[delta.End:]
}

func main() {
	go manager.run()

	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("WebSocket server running on :8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal(err)
	}
}
