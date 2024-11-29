package utils

import (
	"encoding/xml"
	"fmt"
)

// Estrutura base do XML
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

type SISMSG struct {
	Content interface{} `xml:",any"` // Elemento dinâmico baseado no código da mensagem
}

// Estrutura genérica para mensagens como SEL1052, SEL1054
type GenericMessage struct {
	XMLName  xml.Name   `xml:"-"` // Nome do elemento dinâmico (ex.: SEL1052, SEL1054)
	Elements []XMLField `xml:",any"`
}
type XMLField struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// GerarMensagem recebe dados genéricos para gerar uma mensagem baseada no canal.
func GerarMensagem(canal string, codigoMsg string, dados map[string]interface{}) (string, error) {
	if canal == "IOS" {
		// Gera mensagem posicional (exemplo simplificado)
		return gerarStringPosicional(dados), nil
	}

	// Para mensageria, gera XML dinâmico
	content := GenericMessage{
		XMLName:  xml.Name{Local: codigoMsg},
		Elements: []XMLField{},
	}

	// Preenche os elementos dinâmicos
	for key, value := range dados {
		content.Elements = append(content.Elements, XMLField{
			XMLName: xml.Name{Local: key},
			Value:   fmt.Sprintf("%v", value),
		})
	}

	doc := Doc{
		Xmlns: "http://www.bcb.gov.br/SPB/" + codigoMsg + ".xsd",
		BCMSG: BCMSG{
			IdentdDestinatario: "00038121",
			DomSist:            "SPB01",
		},
		SISMSG: SISMSG{
			Content: content,
		},
	}

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
