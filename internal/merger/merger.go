package merger

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type Transaction struct {
	Date     time.Time
	Content  []string
	TxID     string
	Original int
}

func MergeFiles(inputs []string, output string) error {
	var allTransactions []Transaction

	for _, input := range inputs {
		transactions, err := readTransactions(input)
		if err != nil {
			return err
		}
		allTransactions = append(allTransactions, transactions...)
	}

	deduplicatedTransactions := deduplicateTransactions(allTransactions)
	sortTransactions(deduplicatedTransactions)

	return writeTransactions(output, deduplicatedTransactions)
}

func readTransactions(filename string) ([]Transaction, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var transactions []Transaction
	var currentTransaction []string
	var currentDate time.Time
	var currentTxID string
	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if len(strings.TrimSpace(line)) == 0 {
			if len(currentTransaction) > 0 {
				transactions = append(transactions, Transaction{
					Date:     currentDate,
					Content:  currentTransaction,
					TxID:     currentTxID,
					Original: lineCount - len(currentTransaction),
				})
				currentTransaction = nil
				currentTxID = ""
			}
		} else {
			if len(currentTransaction) == 0 {
				date, err := time.Parse("2006-01-02", strings.Fields(line)[0])
				if err != nil {
					return nil, fmt.Errorf("invalid date format on line %d: %v", lineCount, err)
				}
				currentDate = date
			} else if currentTxID == "" {
				currentTxID = extractTransactionTxID(line)
			}
			currentTransaction = append(currentTransaction, line)
		}
	}

	if len(currentTransaction) > 0 {
		transactions = append(transactions, Transaction{
			Date:     currentDate,
			Content:  currentTransaction,
			TxID:     currentTxID,
			Original: lineCount - len(currentTransaction) + 1,
		})
	}

	return transactions, scanner.Err()
}

func extractTransactionTxID(line string) string {
	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, "; txid:") {
		return strings.TrimSpace(trimmedLine[7:])
	}
	return ""
}

func deduplicateTransactions(transactions []Transaction) []Transaction {
	seen := make(map[string]bool)
	var result []Transaction

	for _, t := range transactions {
		if t.TxID != "" {
			if !seen[t.TxID] {
				seen[t.TxID] = true
				result = append(result, t)
			}
		} else {
			result = append(result, t)
		}
	}

	return result
}

func sortTransactions(transactions []Transaction) {
	sort.Slice(transactions, func(i, j int) bool {
		if transactions[i].Date.Equal(transactions[j].Date) {
			return transactions[i].Original < transactions[j].Original
		}
		return transactions[i].Date.Before(transactions[j].Date)
	})
}

func writeTransactions(filename string, transactions []Transaction) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, t := range transactions {
		for _, line := range t.Content {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}
		_, err := writer.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}
