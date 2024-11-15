package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"oraculo-selic/db"
	"oraculo-selic/models"
)

// PassoTesteController lida com as operações de cenários
type PassoTesteController struct {
	DB *db.DB
}

// NewPassoTesteController cria uma nova instância de PassoTesteController
func NewPassoTesteController(db *db.DB) *PassoTesteController {
	return &PassoTesteController{DB: db}
}

// SavePassoTesteHandler manipula a requisição para salvar um novo passo teste
func (cc *PassoTesteController) SavePassoTesteHandler(w http.ResponseWriter, r *http.Request) {
	var passoTeste models.PassoTeste

	// Decodifica o JSON da requisição para a estrutura PassoTeste
	if err := json.NewDecoder(r.Body).Decode(&passoTeste); err != nil {
		log.Printf("Dados inválidos: %v", err)
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Salva o passo teste no banco de dados
	if err := cc.DB.SavePassoTeste(&passoTeste); err != nil {
		log.Printf("Erro ao salvar passo teste: %v", err)
		http.Error(w, "Erro ao salvar passo teste", http.StatusInternalServerError)
		return
	}

	// Retorna o cenário salvo com status 201 Created
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(passoTeste)
}

// GetPassoTesteHandler manipula a requisição para buscar passo teste
func (cc *PassoTesteController) GetPassoTesteHandler(w http.ResponseWriter, r *http.Request) {
	passoTeste, err := cc.DB.GetPassoTeste()
	if err != nil {
		log.Printf("Erro ao buscar passo teste: %v", err)
		http.Error(w, "Erro ao buscar passo teste", http.StatusInternalServerError)
		return
	}

	// Retorna a lista de passoTeste como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(passoTeste)
}
