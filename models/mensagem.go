package models

type Mensagem struct {
	ID             int    `json:"id"`
	CorrelationID  string `json:"correlationId" db:"txt_correl_id"` // Mapeamento para o campo da tabela
	CodigoMensagem string `json:"codigoMensagem"`
	Canal          string `json:"canal"`
	XML            string `json:"xml,omitempty"`
	StringSelic    string `json:"stringSelic,omitempty"`
	Status         string `json:"status"`
	DataInclusao   string `json:"dataInclusao"`
}
