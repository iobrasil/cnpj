package parser_test

import (
	"sort"
	"testing"

	"github.com/iobrasil/cnpj/parser"
)

func TestEnterprise_Run(t *testing.T) {
	filePath := "testdata/TEST.EMPRECSV"
	records := parseCSV(t, filePath)

	expected := make([]any, len(records))
	for i, row := range records {
		expected[i] = parser.EnterpriseData{
			BasicCNPJ:                row[0],
			CorporateName:            row[1],
			LegalNature:              row[2],
			ResponsibleQualification: row[3],
			SocialCapital:            row[4],
			CompanySize:              row[5],
			FederativeEntity:         row[6],
		}
	}

	sorter := func(data []any) {
		sort.Slice(data, func(i, j int) bool {
			return data[i].(parser.EnterpriseData).BasicCNPJ <
				data[j].(parser.EnterpriseData).BasicCNPJ
		})
	}

	validateParser(t, parser.NewEnterprise, sorter, filePath, expected)
}

func BenchmarkEnterprise_Run(b *testing.B) {
	suite := testSuite{}
	for b.Loop() {
		parser.New(parser.NewEnterprise(suite.persist)).Run("testdata/TEST.EMPRECSV")
	}
}
