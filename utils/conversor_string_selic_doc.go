package utils

import (
	"encoding/xml"
	"fmt"
)

// Doc Estrutura base do XML
type Doc struct {
	XMLName xml.Name `xml:"DOC"`
	Xmlns   string   `xml:"xmlns,attr"`
	BCMSG   BCMSG    `xml:"BCMSG"`
	SISMSG  SISMSG   `xml:"SISMSG"`
}

type BCMSG struct {
	IdentdDestinatario string `xml:"IdentdDestinatario"`
	DomSist            string `xml:"DomSist"`
}

// SISMSG ajustada para receber um conteúdo genérico diretamente
type SISMSG struct {
	XMLName xml.Name        `xml:"SISMSG"`
	Content *GenericMessage // Aponta diretamente para a mensagem genérica
}

// GenericMessage Estrutura específica para mensagens como SEL1052, SEL1054
type GenericMessage struct {
	XMLName xml.Name `xml:""` // Permite a definição dinâmica do nome do elemento
	Emi     string   `xml:"Emi"`
	NUOp    string   `xml:"NUOp"`
	CtCed   string   `xml:"ctCed"`
	CtCes   string   `xml:"ctCes"`
	VlrFin  string   `xml:"VlrFinanc"`
	Pu      string   `xml:"Pu"`
}

// GerarMensagem recebe dados genéricos para gerar uma mensagem baseada no canal.
func GerarMensagem(canal string, codigoMsg string, dados map[string]interface{}) (string, error) {
	if canal == "IOS" {
		// Gera mensagem posicional
		return gerarStringPosicional(dados), nil
	}

	// Adiciona o prefixo `SEL` ao código da mensagem
	codigoComPrefixo := "SEL" + codigoMsg

	// Cria a estrutura para o corpo da mensagem
	content := &GenericMessage{
		XMLName: xml.Name{Local: codigoComPrefixo},
		Emi:     fmt.Sprintf("%v", dados["Emissor"]),
		NUOp:    fmt.Sprintf("%v", dados["Número Comando"]),
		CtCed:   fmt.Sprintf("%v", dados["Conta Cedente"]),
		CtCes:   fmt.Sprintf("%v", dados["Conta Cessionária"]),
		VlrFin:  fmt.Sprintf("%v", dados["Valor Financeiro"]),
		Pu:      fmt.Sprintf("%v", dados["PU"]),
	}

	// Monta o documento completo
	doc := Doc{
		Xmlns: "http://www.bcb.gov.br/SPB/" + codigoComPrefixo + ".xsd",
		BCMSG: BCMSG{
			IdentdDestinatario: "00038121",
			DomSist:            "SPB01",
		},
		SISMSG: SISMSG{
			Content: content, // Usa o conteúdo diretamente
		},
	}

	// Serializa para XML
	xmlBytes, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao gerar XML: %v", err)
	}

	return xml.Header + string(xmlBytes), nil
}

// GerarStringPosicional cria uma string posicional baseada nos dados fornecidos.
func gerarStringPosicional(dados map[string]interface{}) string {
	// Cria a string posicional concatenando os valores do mapa
	return fmt.Sprintf(
		"%s%s%s%s%s%s%s",
		fmt.Sprintf("%v", "SSEIN"), // Converte explicitamente para string
		fmt.Sprintf("%v", dados["Conta Cedente"]),
		fmt.Sprintf("%v", dados["Conta Cessionária"]),
		fmt.Sprintf("%v", dados["Emissor"]),
		fmt.Sprintf("%v", dados["PU"]),
		fmt.Sprintf("%v", dados["Valor Financeiro"]),
		fmt.Sprintf("%v", "000000000000000000000"),
	)
}
