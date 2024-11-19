package repositories

import (
	"database/sql"
	"oraculo-selic/models"
)

type CenarioRepository struct {
	DB *sql.DB
}

// NewCenarioRepository cria uma nova instância de CenarioRepository
func NewCenarioRepository(db *sql.DB) *CenarioRepository {
	return &CenarioRepository{DB: db}
}

// Save salva um novo cenário no banco de dados
func (repo *CenarioRepository) Save(cenario *models.Cenario) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}

	// Insere o cenário
	err = tx.QueryRow(
		"INSERT INTO CENARIOS (TXT_DESCRICAO, TXT_TP_CENARIO) VALUES ($1, $2) RETURNING id",
		cenario.Descricao, cenario.Tipo,
	).Scan(&cenario.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insere os passos testes associados na tabela de relação
	for _, cenarioPassoTeste := range cenario.CenariosPassosTestes {
		_, err := tx.Exec(
			"INSERT INTO CENARIOS_PASSOS_TESTES (id_cenario, id_passo_teste, ordenacao) VALUES ($1, $2, $3)",
			cenario.ID, cenarioPassoTeste.PassoTesteID, cenarioPassoTeste.Ordenacao,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// GetAll busca todos os cenários com seus passos testes associados
func (repo *CenarioRepository) GetAll() ([]models.Cenario, error) {
	rows, err := repo.DB.Query(`
		SELECT 
			c.id, c.TXT_DESCRICAO, c.TXT_TP_CENARIO, c.DT_INCL,
			COALESCE(cp.id_passo_teste, 0) AS passo_id, COALESCE(cp.ordenacao, 0) AS ordenacao,
			COALESCE(pt.id, 0) AS passo_id, COALESCE(pt.TXT_DESCRICAO, '') AS passo_descricao, 
			COALESCE(pt.TXT_TP_PASSO_TESTE, '') AS tipo_passo_teste, COALESCE(pt.TXT_COD_MSG, '') AS codigo_msg
		FROM CENARIOS c
		LEFT JOIN CENARIOS_PASSOS_TESTES cp ON c.id = cp.id_cenario
		LEFT JOIN PASSOS_TESTES pt ON cp.id_passo_teste = pt.id
		ORDER BY c.id, cp.ordenacao
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cenarios []models.Cenario
	cenarioMap := make(map[int]*models.Cenario)

	for rows.Next() {
		var (
			cenarioID           int
			cenario             models.Cenario
			cenarioPassoTeste   models.CenariosPassosTestes
			passoTeste          models.PassoTeste
			passoTesteID        int
			ordenacao           int
			passoTesteDescricao string
			tipoPassoTeste      string
			codigoMsg           string
		)

		err := rows.Scan(
			&cenarioID,
			&cenario.Descricao, &cenario.Tipo, &cenario.DataInclusao,
			&passoTesteID, &ordenacao,
			&passoTeste.ID, &passoTesteDescricao, &tipoPassoTeste, &codigoMsg,
		)
		if err != nil {
			return nil, err
		}

		// Evita duplicação de cenários
		existingCenario, exists := cenarioMap[cenarioID]
		if !exists {
			cenario.ID = cenarioID
			cenario.CenariosPassosTestes = []models.CenariosPassosTestes{}
			cenario.PassosTestes = []models.PassoTeste{}
			cenarioMap[cenarioID] = &cenario
			cenarios = append(cenarios, cenario)
		} else {
			cenario = *existingCenario
		}

		// Adiciona os passos testes associados, se existir
		if passoTesteID != 0 {
			cenarioPassoTeste = models.CenariosPassosTestes{
				CenarioID:    cenarioID,
				PassoTesteID: passoTesteID,
				Ordenacao:    ordenacao,
			}
			cenario.CenariosPassosTestes = append(cenario.CenariosPassosTestes, cenarioPassoTeste)

			// Adiciona detalhes do passo teste
			passoTeste.ID = passoTesteID
			passoTeste.Descricao = passoTesteDescricao
			passoTeste.TipoPassoTeste = tipoPassoTeste
			passoTeste.CodigoMsg = codigoMsg
			cenario.PassosTestes = append(cenario.PassosTestes, passoTeste)
		}
	}

	return cenarios, nil
}

// SaveOrUpdateRelacionamentos atualiza os relacionamentos entre um cenário e seus passos testes
func (repo *CenarioRepository) SaveOrUpdateRelacionamentos(relacionamentos []models.CenariosPassosTestes) error {
	if len(relacionamentos) == 0 {
		return nil
	}

	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}

	// Remove os relacionamentos existentes para o cenário
	cenarioID := relacionamentos[0].CenarioID
	_, err = tx.Exec("DELETE FROM CENARIOS_PASSOS_TESTES WHERE id_cenario = $1", cenarioID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insere os novos relacionamentos
	for _, rel := range relacionamentos {
		_, err := tx.Exec(
			"INSERT INTO CENARIOS_PASSOS_TESTES (id_cenario, id_passo_teste, ordenacao) VALUES ($1, $2, $3)",
			rel.CenarioID, rel.PassoTesteID, rel.Ordenacao,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
