package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"oraculo-selic/db"
	"oraculo-selic/models"
)

// CenarioController lida com as operações de cenários
type CenarioController struct {
	DB *db.DB
}

// NewCenarioController cria uma nova instância de CenarioController
func NewCenarioController(db *db.DB) *CenarioController {
	return &CenarioController{DB: db}
}

// SaveCenarioHandler manipula a requisição para salvar um novo cenário
func (cc *CenarioController) SaveCenarioHandler(w http.ResponseWriter, r *http.Request) {
	var cenario models.Cenario

	// Decodifica o JSON da requisição para a estrutura Cenario
	if err := json.NewDecoder(r.Body).Decode(&cenario); err != nil {
		log.Printf("Dados inválidos: %v", err)
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Salva o cenário no banco de dados
	if err := cc.DB.SaveCenario(&cenario); err != nil {
		log.Printf("Erro ao salvar cenário: %v", err)
		http.Error(w, "Erro ao salvar cenário", http.StatusInternalServerError)
		return
	}

	// Retorna o cenário salvo com status 201 Created
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cenario)
}

// GetCenariosHandler manipula a requisição para buscar cenários
func (cc *CenarioController) GetCenariosHandler(w http.ResponseWriter, r *http.Request) {
	cenarios, err := cc.DB.GetCenarios()
	if err != nil {
		log.Printf("Erro ao buscar cenários: %v", err)
		http.Error(w, "Erro ao buscar cenários", http.StatusInternalServerError)
		return
	}

	// Retorna a lista de cenários como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cenarios)
}
