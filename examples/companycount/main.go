package main

import (
	"flag"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/gotomgo/igdb"
)

var key string

func init() {
	flag.StringVar(&key, "k", "", "API key")
	flag.Parse()
}

func main() {
	if key == "" {
		fmt.Println("No key provided. Please run: companycount -k YOUR_API_KEY")
		return
	}

	c := igdb.NewClient(key, nil)

	// Count number of US companies
	USA, err := c.Companies.Count(igdb.SetFilter("country", igdb.OpEquals, "840"))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Count number of UK companies
	UK, err := c.Companies.Count(igdb.SetFilter("country", igdb.OpEquals, "826"))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Count number of JP companies
	JP, err := c.Companies.Count(igdb.SetFilter("country", igdb.OpEquals, "392"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("Number of Companies Based In Specific Countries:\n\n")
	fmt.Println("USA -", USA)
	fmt.Println("UK -", UK)
	fmt.Println("Japan -", JP)

	matches, err := c.Companies.List(nil, igdb.SetFields("*"), igdb.SetFilter("changed_company_id", igdb.OpGreaterThan, "0"))
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, match := range matches {
		spew.Dump(match)
	}

	return
}
