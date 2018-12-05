package main

import (
	"flag"
	"fmt"

	"github.com/gotomgo/igdb"
)

var key string

func init() {
	flag.StringVar(&key, "k", "", "API key")
	flag.Parse()
}

func main() {
	if key == "" {
		fmt.Println("No key provided. Please run: topgames -k YOUR_API_KEY")
		return
	}

	c := igdb.NewClientEx(key, true, nil)

	// Composing options set to retrieve top 5 popular results
	byPop := igdb.ComposeOptions(
		igdb.SetFields("name", "cover"),
		igdb.SetOrder("popularity", igdb.OrderDescending),
		igdb.SetFilter("version_parent", igdb.OpNotExists),
	)

	var err error

	// Retrieve PS4 inter-console exclusives
	p, err := c.Games.ListPaginated(
		500,
		byPop, // top 5 popular results
		igdb.SetFilter("platforms", igdb.OpIn, "48"), // only PS4 games
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	moreItems := true
	var PS4 []*igdb.Game

	for moreItems {
		PS4, moreItems, err = c.Games.GetPaginated(p)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Top PS4 Games:")
		for _, v := range PS4 {
			img, err := v.Cover.SizedURL(igdb.Size1080p, 1) // resize to largest image available
			if err != nil {
				// fmt.Printf("image error for %s: %s\n", v.Name, err)
				fmt.Printf("%s - (%s)\n", v.Name, err)
			} else {
				fmt.Printf("%s - %s\n", v.Name, img)
			}
		}
	}

	return
}
