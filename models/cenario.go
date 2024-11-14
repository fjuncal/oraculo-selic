package models

type Cenario struct {
	ID               int     `json:"id" db:"id"`
	Descricao        string  `json:"descricao" db:"TXT_DESCRICAO"`
	TipoCenario      string  `json:"tipoCenario" db:"TXT_TP_CENARIO"`
	Canal            string  `json:"canal" db:"TXT_CANAL"`
	CodigoMsg        string  `json:"codigoMsg" db:"TXT_COD_MSG"`
	MsgDocXML        string  `json:"xml" db:"TXT_MSG_DOC_XML"`
	Msg              string  `json:"stringSelic" db:"TXT_MSG"`
	ContaCedente     string  `json:"contaCedente" db:"TXT_CT_CED"`
	ContaCessionario string  `json:"contaCessionaria" db:"TXT_CT_CESS"`
	NumeroOperacao   string  `json:"numeroOperacaoSelic" db:"TXT_NUM_OP"`
	Emissor          string  `json:"emissor" db:"TXT_EMISSOR"`
	ValorFinanceiro  float64 `json:"valorFinanceiro" db:"VAL_FIN"`
	ValorPU          float64 `json:"precoUnitario" db:"VAL_PU"`
	DataInclusao     string  `json:"dataInclusao" db:"DT_INCL"`
}
