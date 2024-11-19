package models

type CenariosPassosTestes struct {
	CenarioID    int `json:"cenarioId" db:"id_cenario"`
	PassoTesteID int `json:"passoTesteId" db:"id_passo_teste"`
	Ordenacao    int `json:"ordenacao" db:"ordenacao"`
}
