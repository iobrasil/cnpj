package parser_test

import (
	"sort"
	"testing"

	"github.com/iobrasil/cnpj/parser"
)

func TestEstablishment_Run(t *testing.T) {
	filePath := "testdata/TEST.ESTABELE"
	records := parseCSV(t, filePath)

	// Parse expected enterprises from CSV
	expected := make([]any, len(records))
	for i, row := range records {
		expected[i] = parser.EstablishmentData{
			CNPJBasico:            row[0],
			CNPJOrdem:             row[1],
			CNPJDV:                row[2],
			IdentificadorMatriz:   row[3],
			NomeFantasia:          row[4],
			SituacaoCadastral:     row[5],
			DataSituacaoCadastral: row[6],
			MotivoSituacao:        row[7],
			CidadeExterior:        row[8],
			Pais:                  row[9],
			DataInicioAtividade:   row[10],
			CNAEPrincipal:         row[11],
			CNAESecundaria:        row[12],
			TipoLogradouro:        row[13],
			Logradouro:            row[14],
			Numero:                row[15],
			Complemento:           row[16],
			Bairro:                row[17],
			CEP:                   row[18],
			UF:                    row[19],
			Municipio:             row[20],
			DDD1:                  row[21],
			Telefone1:             row[22],
			DDD2:                  row[23],
			Telefone2:             row[24],
			DDDFax:                row[25],
			Fax:                   row[26],
			Email:                 row[27],
			SituacaoEspecial:      row[28],
			DataSituacaoEspecial:  row[29],
		}
	}

	sorter := func(data []any) {
		sort.Slice(data, func(i, j int) bool {
			return data[i].(parser.EstablishmentData).CNPJBasico <
				data[j].(parser.EstablishmentData).CNPJBasico
		})
	}

	validateParser(t, parser.NewEstablishment, sorter, filePath, expected)
}

func BenchmarkEstablishment_Run(b *testing.B) {
	suite := testSuite{}
	for b.Loop() {
		parser.New(parser.NewEstablishment(suite.persist)).Run("testdata/TEST.ESTABELE")
	}
}
