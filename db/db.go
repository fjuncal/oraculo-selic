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

// NewDb cria uma nova conexão com o banco de dados
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

func (db *DB) SaveMessage(message *models.Mensagem) error {
	query := `
		INSERT INTO mensagens (txt_cod_msg, txt_canal, txt_msg_doc_xml, txt_msg, txt_status, dt_incl) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	err := db.Conn.QueryRow(query,
		message.CodigoMensagem,
		message.Canal,
		message.XML,
		message.StringSelic,
		message.Status,
		message.DataInclusao,
	).Scan(&message.ID)

	if err != nil {
		return err
	}
	log.Printf("Mensagem salva com ID: %d\n", message.ID)
	return nil
}

func (db *DB) SaveCenario(cenario *models.Cenario) error {
	query := `
        INSERT INTO cenarios (
            TXT_DESCRICAO, TXT_TP_CENARIO, TXT_CANAL, TXT_COD_MSG,
            TXT_MSG_DOC_XML, TXT_MSG, TXT_CT_CED, TXT_CT_CESS, 
            TXT_NUM_OP, TXT_EMISSOR, VAL_FIN, VAL_PU
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id
    `
	err := db.Conn.QueryRow(query,
		cenario.Descricao,
		cenario.TipoCenario,
		cenario.Canal,
		cenario.CodigoMsg,
		cenario.MsgDocXML,
		cenario.Msg,
		cenario.ContaCedente,
		cenario.ContaCessionario,
		cenario.NumeroOperacao,
		cenario.Emissor,
		cenario.ValorFinanceiro,
		cenario.ValorPU,
	).Scan(&cenario.ID)

	if err != nil {
		return fmt.Errorf("erro ao salvar cenário: %v", err)
	}
	log.Printf("Cenário salvo com ID: %d\n", cenario.ID)
	return nil
}

func (db *DB) GetCenarios() ([]models.Cenario, error) {
	query := `SELECT id, TXT_DESCRICAO, TXT_TP_CENARIO, TXT_CANAL, TXT_COD_MSG, 
                     TXT_MSG_DOC_XML, TXT_MSG, TXT_CT_CED, TXT_CT_CESS, TXT_NUM_OP, 
                     TXT_EMISSOR, VAL_FIN, VAL_PU, DT_INCL 
              FROM cenarios`

	rows, err := db.Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cenários: %v", err)
	}
	defer rows.Close()

	var cenarios []models.Cenario
	for rows.Next() {
		var cenario models.Cenario
		if err := rows.Scan(
			&cenario.ID,
			&cenario.Descricao,
			&cenario.TipoCenario,
			&cenario.Canal,
			&cenario.CodigoMsg,
			&cenario.MsgDocXML,
			&cenario.Msg,
			&cenario.ContaCedente,
			&cenario.ContaCessionario,
			&cenario.NumeroOperacao,
			&cenario.Emissor,
			&cenario.ValorFinanceiro,
			&cenario.ValorPU,
			&cenario.DataInclusao,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear cenário: %v", err)
		}
		cenarios = append(cenarios, cenario)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro na iteração de cenários: %v", err)
	}

	return cenarios, nil
}
