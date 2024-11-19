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
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cenarios); err != nil {
		log.Printf("Erro ao codificar JSON: %v", err)
		http.Error(w, "Erro ao gerar resposta", http.StatusInternalServerError)
	}
}

// SaveRelacionamentoHandler salva ou atualiza os relacionamentos entre cenários e passos testes
func (cc *CenarioController) SaveRelacionamentoHandler(w http.ResponseWriter, r *http.Request) {
	var relacionamentos []models.CenariosPassosTestes

	// Decodifica o corpo da requisição
	if err := json.NewDecoder(r.Body).Decode(&relacionamentos); err != nil {
		log.Printf("Erro ao decodificar relacionamentos: %v", err)
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Validação básica
	if len(relacionamentos) == 0 || relacionamentos[0].CenarioID == 0 {
		http.Error(w, "Cenário ou passos testes inválidos", http.StatusBadRequest)
		return
	}

	// Atualiza os relacionamentos no banco
	if err := cc.Repo.SaveOrUpdateRelacionamentos(relacionamentos); err != nil {
		log.Printf("Erro ao salvar relacionamentos: %v", err)
		http.Error(w, "Erro ao salvar relacionamentos", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Relacionamentos salvos com sucesso"))
}

//func (cc *CenarioController) GetCenariosWithPassosTestesHandler(w http.ResponseWriter, r *http.Request) {
//	cenarios, err := cc.Repo.GetAllWithPassosTestes()
//	if err != nil {
//		log.Printf("Erro ao buscar cenários com passos testes: %v", err)
//		http.Error(w, "Erro ao buscar cenários", http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(cenarios)
//}
