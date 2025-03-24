package parser_test

import (
	"encoding/csv"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/iobrasil/cnpj/parser"
)

type testSuite struct {
	output []any
}

func (t *testSuite) persist(data map[uint64]any) error {
	if t.output == nil {
		t.output = make([]any, len(data))
	}

	i := 0
	for _, d := range data {
		t.output[i] = d
		i++
	}
	return nil
}

func parseCSV(t *testing.T, filePath string) [][]string {
	t.Helper()
	// Setup test data
	file, err := os.Open(filePath)
	require.NoError(t, err, "Failed to open test file")
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	require.NoError(t, err, "Failed to read CSV records")

	return records
}

func validateParser(t *testing.T, pp func(persist func(map[uint64]any) error) parser.Parser, sorter func(data []any), filePath string, expected []any) {
	t.Helper()
	// Run the parser
	suite := testSuite{}
	p := parser.New(pp(suite.persist))
	p.Run(filePath)

	// Process results
	actual := deref(suite.output)

	// Sort both slices for comparison

	sorter(expected)
	sorter(actual)

	// Assertions
	assert.Equal(t, len(expected), len(actual), "Number of enterprises doesn't match")
	assert.Equal(t, expected, actual, "Enterprise data doesn't match expected values")
}

func deref(items []any) []any {
	var result []any
	for _, item := range items {
		val := reflect.ValueOf(item)

		if val.Kind() == reflect.Ptr {
			result = append(result, val.Elem().Interface()) // Dereference the pointer
		} else {
			result = append(result, item)
		}
	}
	return result
}
