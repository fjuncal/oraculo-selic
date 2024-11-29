package repositories

import (
	"database/sql"
	"oraculo-selic/db"
	"oraculo-selic/models"
)

type CenarioRepository struct {
	DB      *sql.DB
	PassoDB *db.DB // Acesso ao método SavePassoTeste
}

// NewCenarioRepository cria uma nova instância de CenarioRepository
func NewCenarioRepository(db *sql.DB, passoDB *db.DB) *CenarioRepository {
	return &CenarioRepository{
		DB:      db,
		PassoDB: passoDB,
	}
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
			c.id AS cenario_id,
			c.TXT_DESCRICAO AS cenario_descricao,
			c.TXT_TP_CENARIO AS cenario_tipo,
			c.DT_INCL AS cenario_data_incl,
			cp.id_cenario,
			cp.id_passo_teste,
			cp.ordenacao,
			pt.id AS passo_teste_id,
			pt.TXT_DESCRICAO AS passo_teste_descricao,
			pt.TXT_TP_PASSO_TESTE AS passo_teste_tipo,
			pt.TXT_COD_MSG AS passo_teste_codigo,
			pt.TXT_MSG_DOC_XML AS passo_teste_xml,
			pt.TXT_MSG AS passo_teste_string_selic,
			pt.TXT_CT_CED AS passo_teste_conta_cedente,
			pt.TXT_CT_CESS AS passo_teste_conta_cessionario,
			pt.TXT_NUM_OP AS passo_teste_num_operacao,
			pt.TXT_EMISSOR AS passo_teste_emissor,
			pt.VAL_FIN AS passo_teste_valor_financeiro,
			pt.VAL_PU AS passo_teste_preco_unitario,
			pt.DT_INCL AS passo_teste_data_inclusao
		FROM CENARIOS c
		LEFT JOIN CENARIOS_PASSOS_TESTES cp ON c.id = cp.id_cenario
		LEFT JOIN PASSOS_TESTES pt ON cp.id_passo_teste = pt.id
		ORDER BY c.id, cp.ordenacao
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		cenarioMap = make(map[int]*models.Cenario)
		cenarios   []models.Cenario
	)

	for rows.Next() {
		var (
			cenarioID                 int
			cenarioDescricao          string
			cenarioTipo               string
			cenarioDataIncl           sql.NullTime
			idCenario                 sql.NullInt64
			idPassoTeste              sql.NullInt64
			ordenacao                 sql.NullInt64
			passoTesteID              sql.NullInt64
			passoTesteDescricao       sql.NullString
			passoTesteTipo            sql.NullString
			passoTesteCodigo          sql.NullString
			passoTesteXML             sql.NullString
			passoTesteString          sql.NullString
			passoTesteCedente         sql.NullString
			passoTesteCessionario     sql.NullString
			passoTesteNumOperacao     sql.NullString
			passoTesteEmissor         sql.NullString
			passoTesteValorFinanceiro sql.NullFloat64
			passoTestePrecoUnitario   sql.NullFloat64
			passoTesteDataIncl        sql.NullTime
		)

		err := rows.Scan(
			&cenarioID,
			&cenarioDescricao,
			&cenarioTipo,
			&cenarioDataIncl,
			&idCenario,
			&idPassoTeste,
			&ordenacao,
			&passoTesteID,
			&passoTesteDescricao,
			&passoTesteTipo,
			&passoTesteCodigo,
			&passoTesteXML,
			&passoTesteString,
			&passoTesteCedente,
			&passoTesteCessionario,
			&passoTesteNumOperacao,
			&passoTesteEmissor,
			&passoTesteValorFinanceiro,
			&passoTestePrecoUnitario,
			&passoTesteDataIncl,
		)
		if err != nil {
			return nil, err
		}

		// Verificar se o cenário já existe no mapa
		cenario, exists := cenarioMap[cenarioID]
		if !exists {
			cenario = &models.Cenario{
				ID:                   cenarioID,
				Descricao:            cenarioDescricao,
				Tipo:                 cenarioTipo,
				DataInclusao:         cenarioDataIncl.Time.Format("2006-01-02 15:04:05"),
				PassosTestes:         []models.PassoTeste{},
				CenariosPassosTestes: []models.CenariosPassosTestes{},
			}
			cenarioMap[cenarioID] = cenario
		}

		// Adicionar PassoTeste ao cenário, se houver
		if passoTesteID.Valid {
			passoTeste := models.PassoTeste{
				ID:               int(passoTesteID.Int64),
				Descricao:        passoTesteDescricao.String,
				TipoPassoTeste:   passoTesteTipo.String,
				CodigoMsg:        passoTesteCodigo.String,
				MsgDocXML:        passoTesteXML.String,
				Msg:              passoTesteString.String,
				ContaCedente:     passoTesteCedente.String,
				ContaCessionario: passoTesteCessionario.String,
				NumeroOperacao:   passoTesteNumOperacao.String,
				Emissor:          passoTesteEmissor.String,
				ValorFinanceiro:  passoTesteValorFinanceiro.Float64,
				ValorPU:          passoTestePrecoUnitario.Float64,
				DataInclusao:     passoTesteDataIncl.Time.Format("2006-01-02 15:04:05"),
			}
			cenario.PassosTestes = append(cenario.PassosTestes, passoTeste)
		}

		// Adicionar relação ao CenariosPassosTestes, se existir
		if idCenario.Valid && idPassoTeste.Valid {
			cenario.CenariosPassosTestes = append(cenario.CenariosPassosTestes, models.CenariosPassosTestes{
				CenarioID:    int(idCenario.Int64),
				PassoTesteID: int(idPassoTeste.Int64),
				Ordenacao:    int(ordenacao.Int64),
			})
		}
	}

	// Converter o mapa para um slice
	for _, cenario := range cenarioMap {
		cenarios = append(cenarios, *cenario)
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

//func (repo *CenarioRepository) GetAllWithPassosTestes() ([]models.Cenario, error) {
//	rows, err := repo.DB.Query(`
//		SELECT
//			c.id AS cenario_id,
//			c.txt_descricao AS cenario_descricao,
//			c.txt_tp_cenario AS cenario_tipo,
//			c.dt_incl AS cenario_data_inclusao,
//			pt.id AS passo_teste_id,
//			pt.txt_descricao AS passo_teste_descricao,
//			pt.txt_cod_msg AS passo_teste_codigo_msg,
//			pt.txt_tp_passo_teste AS passo_teste_tipo,
//			pt.dt_incl AS passo_teste_data_inclusao
//		FROM cenarios c
//		LEFT JOIN cenarios_passos_testes cpt ON c.id = cpt.id_cenario
//		LEFT JOIN passos_testes pt ON cpt.id_passo_teste = pt.id
//		ORDER BY c.id, cpt.ordenacao
//	`)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var cenarios []models.Cenario
//	cenarioMap := make(map[int]*models.Cenario)
//
//	for rows.Next() {
//		var (
//			cenarioID           int
//			cenarioDescricao    string
//			cenarioTipo         string
//			cenarioDataInclusao sql.NullTime
//			passoTesteID        sql.NullInt64
//			passoTesteDescricao sql.NullString
//			passoTesteCodigoMsg sql.NullString
//			passoTesteTipo      sql.NullString
//			passoTesteDataIncl  sql.NullTime
//		)
//
//		err := rows.Scan(
//			&cenarioID,
//			&cenarioDescricao,
//			&cenarioTipo,
//			&cenarioDataInclusao,
//			&passoTesteID,
//			&passoTesteDescricao,
//			&passoTesteCodigoMsg,
//			&passoTesteTipo,
//			&passoTesteDataIncl,
//		)
//		if err != nil {
//			return nil, err
//		}
//
//		// Evita duplicação de cenários
//		cenario, exists := cenarioMap[cenarioID]
//		if !exists {
//			cenario = &models.Cenario{
//				ID:           cenarioID,
//				Descricao:    cenarioDescricao,
//				Tipo:         cenarioTipo,
//				DataInclusao: cenarioDataInclusao.Time.Format("2006-01-02T15:04:05Z07:00"),
//				PassosTestes: []models.PassoTeste{},
//			}
//			cenarioMap[cenarioID] = cenario
//			cenarios = append(cenarios, *cenario)
//		}
//
//		// Adiciona passo teste se existir
//		if passoTesteID.Valid {
//			cenario.PassosTestes = append(cenario.PassosTestes, models.PassoTeste{
//				ID:             int(passoTesteID.Int64),
//				Descricao:      passoTesteDescricao.String,
//				CodigoMsg:      passoTesteCodigoMsg.String,
//				TipoPassoTeste: passoTesteTipo.String,
//				DataInclusao:   passoTesteDataIncl.Time.Format("2006-01-02T15:04:05Z07:00"),
//			})
//		}
//	}
//
//	return cenarios, nil
//}
