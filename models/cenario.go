package models

import "time"

type Cenario struct {
	ID               int       `json:"id" db:"id"`
	Descricao        string    `json:"descricao" db:"TXT_DESCRICAO"`
	TipoCenario      string    `json:"tipo_cenario" db:"TXT_TP_CENARIO"`
	Canal            string    `json:"canal" db:"TXT_CANAL"`
	CodigoMsg        string    `json:"codigo_msg" db:"TXT_COD_MSG"`
	MsgDocXML        string    `json:"msg_doc_xml" db:"TXT_MSG_DOC_XML"`
	Msg              string    `json:"msg" db:"TXT_MSG"`
	ContaCedente     string    `json:"conta_cedente" db:"TXT_CT_CED"`
	ContaCessionario string    `json:"conta_cessionario" db:"TXT_CT_CESS"`
	NumeroOperacao   string    `json:"numero_operacao" db:"TXT_NUM_OP"`
	Emissor          string    `json:"emissor" db:"TXT_EMISSOR"`
	ValorFinanceiro  float64   `json:"valor_financeiro" db:"VAL_FIN"`
	ValorPU          float64   `json:"valor_pu" db:"VAL_PU"`
	DataInclusao     time.Time `json:"data_inclusao" db:"DT_INCL"`
}
