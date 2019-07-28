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

	c := igdb.NewClientEx(igdb.IgdbURL, key, true, nil)

	var err error

	// Retrieve PS4 inter-console exclusives
	p, err := c.Franchises.ListPaginated(
		500,
		igdb.SetFields("*"), // , "name", "slug", "logo"),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	moreItems := true
	var platforms []*igdb.Franchise

	for moreItems {
		platforms, moreItems, err = c.Franchises.GetPaginated(p)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Platforms:")
		for _, v := range platforms {
			// var img string
			//var err error

			/*img, err = v.Logo.SizedURL(igdb.SizeThumb, 2) // resize to largest image available
			if err != nil {
				// fmt.Printf("image error for %s: %s\n", v.Name, err)
				img = ""
			}*/

			fmt.Printf("ID=%d,Name=%s,Slug=%s\n", v.ID, v.Name, v.Slug) /*, img) */
		}
	}

	return
}
