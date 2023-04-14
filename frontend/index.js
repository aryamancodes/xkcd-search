function search(prompt) {
    var xmlhttp = new XMLHttpRequest();
    //check for empty string before calling api
    xmlhttp.open("GET", "http://localhost:8080/search?q="+prompt + "&autocorrect=true")
    xmlhttp.send()
}


function suggest(prompt) {
    var xmlhttp = new XMLHttpRequest();
    //Use debouncing -- don't want to call api on every single keystroke
    xmlhttp.open("GET", "http://localhost:8080/suggest?q="+prompt)
    xmlhttp.send()
}

/* references:
    xkcd.com
    https://relevantxkcd.appspot.com/
    xkcd-search.typesense.org/
    https://catche.co/xkcd
*/
