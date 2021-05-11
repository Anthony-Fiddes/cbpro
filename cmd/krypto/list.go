package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	listCommand  = "list"
	listLength   = 10
	listTemplate = "%-10s%-10s%-10s\n"
)

func askToContinue() bool {
	for {
		fmt.Print("Continue? (y/n) ")
		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		input := s.Text()
		input = strings.ToLower(input)
		input = strings.TrimSpace(input)
		if input == "y" || input == "" {
			return true
		} else if input == "n" {
			return false
		}
	}
}

func list(args []string) error {
	// TODO: Add the ability to filer results
	products, err := client.GetProducts()
	if err != nil {
		return err
	}

	header := strings.Builder{}
	_, err = header.WriteString(fmt.Sprintf(
		listTemplate,
		"ID",
		"Base",
		"Quote",
	))
	if err != nil {
		return err
	}

	width := header.Len()
	for i := 0; i < width; i++ {
		header.WriteRune('=')
	}
	fmt.Println(header.String())

	for i, p := range products {
		fmt.Printf(
			listTemplate,
			p.ID,
			p.BaseCurrency,
			p.QuoteCurrency,
		)
		if (i+1)%listLength == 0 {
			fmt.Println()
			c := askToContinue()
			if !c {
				break
			}
			fmt.Println()
			fmt.Println(header.String())
		}
	}
	return nil
}
