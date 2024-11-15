package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Importa o driver do PostgreSQL de forma anônima para registrá-lo
	"log"
	"oraculo-selic/models"
)

// DB representa a conexão com o banco de dados
type DB struct {
	Conn *sql.DB
}

// SavePassoTeste método para salvar passo teste no banco
func (db *DB) SavePassoTeste(passoTeste *models.PassoTeste) error {
	query := `
        INSERT INTO PASSO_TESTE (
            TXT_DESCRICAO, TXT_TP_PASSO_TESTE, TXT_CANAL, TXT_COD_MSG,
            TXT_MSG_DOC_XML, TXT_MSG, TXT_CT_CED, TXT_CT_CESS, 
            TXT_NUM_OP, TXT_EMISSOR, VAL_FIN, VAL_PU
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id
    `
	err := db.Conn.QueryRow(query,
		passoTeste.Descricao,
		passoTeste.TipoPassoTeste,
		passoTeste.Canal,
		passoTeste.CodigoMsg,
		passoTeste.MsgDocXML,
		passoTeste.Msg,
		passoTeste.ContaCedente,
		passoTeste.ContaCessionario,
		passoTeste.NumeroOperacao,
		passoTeste.Emissor,
		passoTeste.ValorFinanceiro,
		passoTeste.ValorPU,
	).Scan(&passoTeste.ID)

	if err != nil {
		return fmt.Errorf("erro ao salvar passo teste: %v", err)
	}
	log.Printf("Passo teste salvo com ID: %d\n", passoTeste.ID)
	return nil
}

// GetPassoTeste método para buscar passo teste
func (db *DB) GetPassoTeste() ([]models.PassoTeste, error) {
	query := `SELECT id, TXT_DESCRICAO, TXT_TP_PASSO_TESTE, TXT_CANAL, TXT_COD_MSG, 
                     TXT_MSG_DOC_XML, TXT_MSG, TXT_CT_CED, TXT_CT_CESS, TXT_NUM_OP, 
                     TXT_EMISSOR, VAL_FIN, VAL_PU, DT_INCL 
              FROM PASSO_TESTE`

	rows, err := db.Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar passo teste: %v", err)
	}
	defer rows.Close()

	var passosTestes []models.PassoTeste
	for rows.Next() {
		var passoTeste models.PassoTeste
		if err := rows.Scan(
			&passoTeste.ID,
			&passoTeste.Descricao,
			&passoTeste.TipoPassoTeste,
			&passoTeste.Canal,
			&passoTeste.CodigoMsg,
			&passoTeste.MsgDocXML,
			&passoTeste.Msg,
			&passoTeste.ContaCedente,
			&passoTeste.ContaCessionario,
			&passoTeste.NumeroOperacao,
			&passoTeste.Emissor,
			&passoTeste.ValorFinanceiro,
			&passoTeste.ValorPU,
			&passoTeste.DataInclusao,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear passo teste: %v", err)
		}
		passosTestes = append(passosTestes, passoTeste)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro na iteração de passo teste: %v", err)
	}

	return passosTestes, nil
}
