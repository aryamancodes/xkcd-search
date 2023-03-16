package main

import (
	"fmt"
	"xkcd/index"
)

func main() {
	comics := index.FetchAllExplanations()

	tf := index.ComputeAllTermFreq(comics)

	df := index.ComputeAllComicFreq(comics, tf)

	fmt.Printf("\nTERM FREQ LOOKS LIKE: %+v\n", tf)
	fmt.Printf("\nDOC FREQ LOOKS LIKE: %+v\n", df)
}
