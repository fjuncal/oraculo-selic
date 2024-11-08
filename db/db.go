package db

import (
	"database/sql"
	_ "github.com/lib/pq" // Importa o driver do PostgreSQL de forma anônima para registrá-lo
	"log"
	"oraculo-selic/models"
)

type DB struct {
	Conn *sql.DB
}

func NewDb(databaseURL string) (*DB, error) {
	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	//ping é para verificar se a conexão foi aberta com sucesso
	if err = conn.Ping(); err != nil {
		return nil, err
	}
	log.Println("Conectado ao banco de dados com sucesso!")
	return &DB{Conn: conn}, nil
}

func (db *DB) SaveMessage(message *models.Message) error {
	query := `INSERT INTO messages (content, status) VALUES ($1, $2) RETURNING id`
	err := db.Conn.QueryRow(query, message.Content, message.Status).Scan(&message.ID)
	if err != nil {
		return err
	}
	log.Printf("Mensagem salva com ID: %d\n", message.ID)
	return nil
}
