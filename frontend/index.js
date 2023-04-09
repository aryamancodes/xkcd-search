function search(prompt) {
    var xmlhttp = new XMLHttpRequest();
    xmlhttp.open("GET", "http://localhost:8080/search?q="+prompt + "&autocorrect=true")
    xmlhttp.send()
}


function suggest(prompt) {
    var xmlhttp = new XMLHttpRequest();
    xmlhttp.open("GET", "http://localhost:8080/suggest?q="+prompt)
    xmlhttp.send()
}

