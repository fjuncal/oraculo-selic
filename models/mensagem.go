package models

type Mensagem struct {
	ID             int    `json:"id"`
	CodigoMensagem string `json:"codigoMensagem"`
	Canal          string `json:"canal"`
	XML            string `json:"xml,omitempty"`
	StringSelic    string `json:"stringSelic,omitempty"`
	Status         string `json:"status"`
	DataInclusao   string `json:"dataInclusao"`
}
