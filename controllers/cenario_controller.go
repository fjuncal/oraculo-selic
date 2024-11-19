package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"oraculo-selic/db/repositories"
	"oraculo-selic/models"
)

type CenarioController struct {
	Repo *repositories.CenarioRepository
}

// NewCenarioController cria uma nova instância de CenarioController
func NewCenarioController(repo *repositories.CenarioRepository) *CenarioController {
	return &CenarioController{Repo: repo}
}

// SaveCenarioHandler cria um novo cenário com passos testes associados
func (cc *CenarioController) SaveCenarioHandler(w http.ResponseWriter, r *http.Request) {
	var cenario models.Cenario

	if err := json.NewDecoder(r.Body).Decode(&cenario); err != nil {
		log.Printf("Erro ao decodificar cenário: %v", err)
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	if err := cc.Repo.Save(&cenario); err != nil {
		log.Printf("Erro ao salvar cenário: %v", err)
		http.Error(w, "Erro ao salvar cenário", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cenario)
}

// GetCenariosHandler busca todos os cenários com seus passos testes
func (cc *CenarioController) GetCenariosHandler(w http.ResponseWriter, r *http.Request) {
	cenarios, err := cc.Repo.GetAll()
	if err != nil {
		log.Printf("Erro ao buscar cenários: %v", err)
		http.Error(w, "Erro ao buscar cenários", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cenarios)
}
