package models

type Cenario struct {
	ID                   int                    `json:"id" db:"id"`
	Descricao            string                 `json:"descricao" db:"TXT_DESCRICAO"`
	Tipo                 string                 `json:"tipo" db:"TXT_TP_CENARIO"`
	DataInclusao         string                 `json:"dataInclusao" db:"DT_INCL"`
	CenariosPassosTestes []CenariosPassosTestes `json:"cenariosPassosTestes"`   // Adiciona este campo
	PassosTestes         []PassoTeste           `json:"passosTestes,omitempty"` // Adiciona este campo

}
