package main

import (
	// "io/ioutil"
	// "log"
	// "net/http"

	// "github.com/tidwall/gjson"

	"xkcd/util"
)

func main() {
	util.GetAllComics()
}

// resp, err := http.Get("https://www.explainxkcd.com/wiki/api.php?action=parse&page=2741:_Wish_Interpretation&format=json")
// if err != nil {
// 	log.Fatalln(err)
// }

// body, err := ioutil.ReadAll(resp.Body)
// if err != nil {
// 	log.Fatalln(err)
// }

// explainedHTML := gjson.Get(string(body), "parse.text")
// log.Println(explainedHTML)
//}
