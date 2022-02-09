<?php
	$username = "root";
	$password = "CSisGreat!7951";
	$server = "127.0.0.1";
	$db_name = "capstone";

	$conn = new mysqli($server, $username, $password, $db_name);

	if ($conn->connect_error) {
		die("Connection failed: " . $conn->connect_error);
	}
?>
