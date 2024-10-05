package merger

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestReadTransactions(t *testing.T) {
	content := `2023-01-01 Transaction 1
  ; txid: abc123
  Account1  $100
  Account2  $-100

2023-01-02 Transaction 2
  Account3  $200
  Account4  $-200

2023-01-01 Transaction 3
  ; txid: def456
  Account5  $300
  Account6  $-300
`

	tmpfile, err := os.CreateTemp("", "test*.journal")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	transactions, err := readTransactions(tmpfile.Name())
	if err != nil {
		t.Fatalf("readTransactions() error = %v", err)
	}

	expected := []Transaction{
		{
			Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			Content: []string{
				"2023-01-01 Transaction 1",
				"  ; txid: abc123",
				"  Account1  $100",
				"  Account2  $-100",
			},
			TxID:     "abc123",
			Original: 1,
		},
		{
			Date: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			Content: []string{
				"2023-01-02 Transaction 2",
				"  Account3  $200",
				"  Account4  $-200",
			},
			TxID:     "",
			Original: 6,
		},
		{
			Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			Content: []string{
				"2023-01-01 Transaction 3",
				"  ; txid: def456",
				"  Account5  $300",
				"  Account6  $-300",
			},
			TxID:     "def456",
			Original: 10,
		},
	}

	if !reflect.DeepEqual(transactions, expected) {
		t.Errorf("readTransactions() = %v, want %v", transactions, expected)
	}
}

func TestReadTransactionsWithPostingTags(t *testing.T) {
	content := `2023-01-01 Transaction 1
  ; txid: abc123
  Account1  $100
  Account2  $-100

2023-01-02 Transaction 2
  ; txid: real
  Account3  $200  ; txid: fake
  Account4  $-200

2023-01-03 Transaction 3
  Account5  $300  ; txid: posting_txid
  Account6  $-300
`

	tmpfile, err := os.CreateTemp("", "test*.journal")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	transactions, err := readTransactions(tmpfile.Name())
	if err != nil {
		t.Fatalf("readTransactions() error = %v", err)
	}

	expected := []Transaction{
		{
			Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			Content: []string{
				"2023-01-01 Transaction 1",
				"  ; txid: abc123",
				"  Account1  $100",
				"  Account2  $-100",
			},
			TxID:     "abc123",
			Original: 1,
		},
		{
			Date: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			Content: []string{
				"2023-01-02 Transaction 2",
				"  ; txid: real",
				"  Account3  $200  ; txid: fake",
				"  Account4  $-200",
			},
			TxID:     "real",
			Original: 6,
		},
		{
			Date: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
			Content: []string{
				"2023-01-03 Transaction 3",
				"  Account5  $300  ; txid: posting_txid",
				"  Account6  $-300",
			},
			TxID:     "",
			Original: 11,
		},
	}

	if !reflect.DeepEqual(transactions, expected) {
		t.Errorf("readTransactions() = %v, want %v", transactions, expected)
	}
}

func TestDeduplicateTransactions(t *testing.T) {
	transactions := []Transaction{
		{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Content: []string{"Transaction 1"}, TxID: "abc123"},
		{Date: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), Content: []string{"Transaction 2"}, TxID: ""},
		{Date: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC), Content: []string{"Transaction 3"}, TxID: "def456"},
		{Date: time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC), Content: []string{"Transaction 4"}, TxID: "abc123"},
	}

	deduplicated := deduplicateTransactions(transactions)

	expected := []Transaction{
		{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Content: []string{"Transaction 1"}, TxID: "abc123"},
		{Date: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), Content: []string{"Transaction 2"}, TxID: ""},
		{Date: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC), Content: []string{"Transaction 3"}, TxID: "def456"},
	}

	if !reflect.DeepEqual(deduplicated, expected) {
		t.Errorf("deduplicateTransactions() = %v, want %v", deduplicated, expected)
	}
}

func TestSortTransactions(t *testing.T) {
	transactions := []Transaction{
		{Date: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), Content: []string{"Transaction 2"}, Original: 2},
		{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Content: []string{"Transaction 1"}, Original: 1},
		{Date: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC), Content: []string{"Transaction 3"}, Original: 3},
		{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Content: []string{"Transaction 4"}, Original: 4},
	}

	sortTransactions(transactions)

	expected := []Transaction{
		{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Content: []string{"Transaction 1"}, Original: 1},
		{Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Content: []string{"Transaction 4"}, Original: 4},
		{Date: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), Content: []string{"Transaction 2"}, Original: 2},
		{Date: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC), Content: []string{"Transaction 3"}, Original: 3},
	}

	if !reflect.DeepEqual(transactions, expected) {
		t.Errorf("After sortTransactions(), got = %v, want %v", transactions, expected)
	}
}

func TestMergeFiles(t *testing.T) {
	content1 := `2023-01-01 Transaction 1
  ; txid: abc123
  Account1  $100
  Account2  $-100

2023-01-02 Transaction 2
  Account3  $200
  Account4  $-200
`

	content2 := `2023-01-01 Transaction 3
  ; txid: def456
  Account5  $300
  Account6  $-300

2023-01-03 Transaction 4
  ; txid: abc123
  Account7  $400
  Account8  $-400
`

	tmpfile1, err := os.CreateTemp("", "test1*.journal")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile1.Name())

	tmpfile2, err := os.CreateTemp("", "test2*.journal")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile2.Name())

	if _, err := tmpfile1.Write([]byte(content1)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile1.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile2.Write([]byte(content2)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile2.Close(); err != nil {
		t.Fatal(err)
	}

	outputFile, err := os.CreateTemp("", "output*.journal")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(outputFile.Name())

	err = MergeFiles([]string{tmpfile1.Name(), tmpfile2.Name()}, outputFile.Name())
	if err != nil {
		t.Fatalf("MergeFiles() error = %v", err)
	}

	output, err := ioutil.ReadFile(outputFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expected := `2023-01-01 Transaction 1
  ; txid: abc123
  Account1  $100
  Account2  $-100

2023-01-01 Transaction 3
  ; txid: def456
  Account5  $300
  Account6  $-300

2023-01-02 Transaction 2
  Account3  $200
  Account4  $-200

`

	if string(output) != expected {
		t.Errorf("MergeFiles() output = %v, want %v", string(output), expected)
	}
}
