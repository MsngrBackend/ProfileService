package events

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

// ProfilePublisher публикует события обновления профиля в NATS (без изменения REST API).
type ProfilePublisher struct {
	conn *nats.Conn
}

// NewProfilePublisher подключается к NATS если NATS_URL задан; иначе публикации no-op.
func NewProfilePublisher() *ProfilePublisher {
	url := os.Getenv("NATS_URL")
	if url == "" {
		return &ProfilePublisher{}
	}
	conn, err := nats.Connect(url)
	if err != nil {
		log.Printf("nats: connect failed (%v), profile events disabled", err)
		return &ProfilePublisher{}
	}
	log.Println("nats: connected (profile events)")
	return &ProfilePublisher{conn: conn}
}

// Close освобождает соединение.
func (p *ProfilePublisher) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}

// PublishProfileUpdated шлёт событие для подписчиков (kind: "profile" | "avatar").
func (p *ProfilePublisher) PublishProfileUpdated(userID, kind string) {
	if p.conn == nil {
		return
	}
	payload, err := json.Marshal(map[string]string{
		"user_id": userID,
		"kind":    kind,
	})
	if err != nil {
		return
	}
	if err := p.conn.Publish("msngr.profile.updated", payload); err != nil {
		log.Printf("nats: publish profile updated: %v", err)
	}
}
