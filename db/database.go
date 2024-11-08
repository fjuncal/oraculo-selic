package db

import (
	"database/sql"
	"fmt"
	"oraculo-selic/models"
)

type DatabaseConnections struct {
	DB1 *sql.DB
	DB2 *sql.DB
	DB3 *sql.DB
}

// NewDatabaseConnections Função para conectar a todas as bases
func NewDatabaseConnections(db1Url, db2Url, db3Url string) (*DatabaseConnections, error) {
	db1, err := sql.Open("postgres", db1Url)
	if err != nil {
		return nil, err
	}

	db2, err := sql.Open("postgres", db2Url)
	if err != nil {
		return nil, err
	}

	db3, err := sql.Open("postgres", db3Url)
	if err != nil {
		return nil, err
	}

	return &DatabaseConnections{
		DB1: db1,
		DB2: db2,
		DB3: db3,
	}, nil
}

func (dbc *DatabaseConnections) SaveMessage(message *models.Message) error {
	query := "INSERT INTO messages (content, status) VALUES ($1, $2) RETURNING id"
	err := dbc.DB1.QueryRow(query, message.Content, message.Status).Scan(&message.ID)
	if err != nil {
		return fmt.Errorf("erro ao salvar mensagem: %v", err)
	}
	return nil
}

func (dbc *DatabaseConnections) Close() {
	dbc.DB1.Close()
	dbc.DB2.Close()
	dbc.DB3.Close()
}
