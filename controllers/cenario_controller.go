package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"mime/multipart"
	"net/http"
	"oraculo-selic/db/repositories"
	"oraculo-selic/models"
	"oraculo-selic/utils"
	"strconv"
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

func (cc *CenarioController) UploadPlanilhaHandler(w http.ResponseWriter, r *http.Request) {
	// Parse do arquivo
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("Erro ao receber arquivo: %v", err)
		http.Error(w, "Erro ao processar arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Processar a planilha e obter os cenários
	cenarios, err := cc.processarPlanilha(file)
	if err != nil {
		log.Printf("Erro ao processar planilha: %v", err)
		http.Error(w, "Erro ao processar planilha", http.StatusInternalServerError)
		return
	}

	// Lista para armazenar os cenários salvos
	var cenariosSalvos []models.Cenario

	// Iterar sobre os cenários e salvar no banco
	for _, cenario := range cenarios {
		// Salvar o cenário
		if err := cc.Repo.Save(&cenario); err != nil {
			log.Printf("Erro ao salvar cenário '%s': %v", cenario.Descricao, err)
			http.Error(w, "Erro ao salvar cenários", http.StatusInternalServerError)
			return
		}

		// Salvar os passos testes e criar os relacionamentos
		var relacionamentos []models.CenariosPassosTestes
		for i := range cenario.PassosTestes {
			passo := &cenario.PassosTestes[i] // Atualizamos o passo diretamente no slice
			// Salvar passo teste
			if err := cc.Repo.PassoDB.SavePassoTeste(passo); err != nil {
				log.Printf("Erro ao salvar passo teste '%s': %v", passo.Descricao, err)
				http.Error(w, "Erro ao salvar passos testes", http.StatusInternalServerError)
				return
			}

			// Criar relacionamento
			relacionamento := models.CenariosPassosTestes{
				CenarioID:    cenario.ID,
				PassoTesteID: passo.ID,
				Ordenacao:    i + 1,
			}
			relacionamentos = append(relacionamentos, relacionamento)
		}

		// Salvar os relacionamentos
		if err := cc.Repo.SaveOrUpdateRelacionamentos(relacionamentos); err != nil {
			log.Printf("Erro ao salvar relacionamentos para o cenário '%s': %v", cenario.Descricao, err)
			http.Error(w, "Erro ao salvar relacionamentos", http.StatusInternalServerError)
			return
		}

		// Adicionar o cenário salvo à lista
		cenariosSalvos = append(cenariosSalvos, cenario)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cenarios)
}

func (cc *CenarioController) processarPlanilha(file multipart.File) ([]models.Cenario, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir o arquivo Excel: %v", err)
	}

	// Obtém a lista de abas
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("nenhuma aba encontrada na planilha")
	}

	log.Printf("Planilhas disponíveis: %v", sheets)

	var cenarios []models.Cenario

	// Processa cada aba da planilha
	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			log.Printf("Erro ao ler linhas da aba '%s': %v", sheet, err)
			continue
		}

		// Inicializa variáveis específicas para esta aba
		passosTestes := []models.PassoTeste{}
		headers := make(map[string]int)
		headersCenario := make(map[string]int)

		// Cria um novo cenário para cada aba
		cenario := models.Cenario{
			Descricao: fmt.Sprintf("Cenário gerado da aba %s", sheet),
			Tipo:      "Importação",
		}

		for i, row := range rows {
			if len(row) == 0 {
				// Ignorar linhas vazias
				continue
			}

			if containsCabecalhoCenario(row) {
				headersCenario = make(map[string]int)
				for idx, col := range row {
					headersCenario[col] = idx
				}
				log.Printf("Cabeçalhos de cenário identificados na aba '%s': %v", sheet, headersCenario)
				continue
			}

			if containsLinhaDescricaoCenario(row) {
				cenario.Descricao = getCellValue(row, headersCenario, "Descrição Cenário")
				cenario.Tipo = getCellValue(row, headersCenario, "Tipo Cenário")
				continue
			}

			if i == 0 || containsSeq(row) {
				// Identifica a linha de cabeçalho
				headers = make(map[string]int)
				for idx, col := range row {
					headers[col] = idx
				}
				log.Printf("Cabeçalhos identificados na aba '%s': %v", sheet, headers)
				continue
			}

			if len(row) < len(headers) || headers["Seq."] >= len(row) || row[headers["Seq."]] == "" {
				log.Printf("Linha %d ignorada na aba '%s': %v", i, sheet, row)
				continue
			}

			passo := models.PassoTeste{
				Descricao:        getCellValue(row, headers, "Descrição"),
				TipoPassoTeste:   getCellValue(row, headers, "TipoPassoTeste"),
				Canal:            getCellValue(row, headers, "Canal"),
				CodigoMsg:        getCellValue(row, headers, "Operação"),
				ContaCedente:     getCellValue(row, headers, "Conta Cedente"),
				ContaCessionario: getCellValue(row, headers, "Conta Cessionária"),
				NumeroOperacao:   getCellValue(row, headers, "Número Comando"),
				Emissor:          getCellValue(row, headers, "Transmissor Debito"),
				ValorFinanceiro:  parseFloat(getCellValue(row, headers, "Valor Financeiro")),
				ValorPU:          parseFloat(getCellValue(row, headers, "PU")),
			}

			// Gerar mensagem para o passo teste
			codigoMsg := getCellValue(row, headers, "Operação")
			dados := map[string]interface{}{
				"Emissor":           passo.Emissor,
				"Número Comando":    passo.NumeroOperacao,
				"Conta Cedente":     passo.ContaCedente,
				"Conta Cessionária": passo.ContaCessionario,
				"Valor Financeiro":  passo.ValorFinanceiro,
				"PU":                passo.ValorPU,
			}

			msg, err := utils.GerarMensagem(passo.Canal, codigoMsg, dados)
			if err != nil {
				log.Printf("Erro ao gerar mensagem para passo teste na aba '%s': %v", sheet, err)
				continue
			}

			if passo.Canal == "IOS" {
				passo.Msg = msg
			} else {
				passo.MsgDocXML = msg
			}

			passosTestes = append(passosTestes, passo)
		}

		// Verifica se a aba contém passos testes e adiciona ao cenário
		if len(passosTestes) > 0 {
			cenario.PassosTestes = passosTestes
			cenarios = append(cenarios, cenario)
		} else {
			log.Printf("Nenhum passo teste encontrado na aba '%s'", sheet)
		}
	}

	return cenarios, nil
}

func containsLinhaDescricaoCenario(row []string) bool {
	for _, col := range row {
		if col == "***" {
			return true
		}
	}
	return false
}

func containsCabecalhoCenario(row []string) bool {
	for _, col := range row {
		if col == "Seq.Cenário" {
			return true
		}
	}
	return false

}

// Helper para verificar se a linha contém "Seq."
func containsSeq(row []string) bool {
	for _, col := range row {
		if col == "Seq." {
			return true
		}
	}
	return false
}

// Helper para obter o valor de uma célula com base nos cabeçalhos
func getCellValue(row []string, headers map[string]int, columnName string) string {
	index, exists := headers[columnName]
	if !exists || index >= len(row) {
		return ""
	}
	return row[index]
}

// Helper para converter string para float
func parseFloat(value string) float64 {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return floatValue
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
