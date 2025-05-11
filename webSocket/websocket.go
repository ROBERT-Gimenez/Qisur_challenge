package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

type EventManager struct {
	clients  map[*websocket.Conn]bool
	broadcast chan Message
}

func NewEventManager() *EventManager {
	return &EventManager{
		clients:  make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
	}
}

var eventManager = NewEventManager()

func (em *EventManager) AddClient(conn *websocket.Conn) {
	em.clients[conn] = true
}

func (em *EventManager) RemoveClient(conn *websocket.Conn) {
	delete(em.clients, conn)
}

func (em *EventManager) BroadcastMessage(msg Message) {
	for client := range em.clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Println("Error al enviar mensaje a cliente:", err)
			client.Close()
			em.RemoveClient(client)
		}
	}
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
    Type string      `json:"type"`
    Data ProductData `json:"data"`
}

type ProductData struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}


func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Error al actualizar a WebSocket:", err)
        return
    }
    defer conn.Close()

    eventManager.AddClient(conn)
    defer eventManager.RemoveClient(conn)

    go func() {
        for {
            select {
            case msg := <-eventManager.broadcast:
                eventManager.BroadcastMessage(msg)
            }
        }
    }()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error al leer mensaje:", err)
            break
        }

        var message Message
        if err := json.Unmarshal(msg, &message); err != nil {
            log.Println("Error al parsear mensaje:", err)
            continue
        }

        log.Printf("Mensaje recibido: %s\n", message.Data)

        switch message.Type {
        case "create":
            eventManager.BroadcastMessage(Message{
                Type:    "product_created",
                Data: message.Data,
            })
        case "update":
            eventManager.BroadcastMessage(Message{
                Type:    "product_updated",
                Data: message.Data,
            })
        case "delete":
            eventManager.BroadcastMessage(Message{
                Type:    "product_deleted",
                Data: message.Data,
            })
        default:
            log.Printf("Tipo de mensaje desconocido: %s", message.Type)
        }
    }
}

func GetEventManager() *EventManager {
	return eventManager
}