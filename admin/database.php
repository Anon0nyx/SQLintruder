<?php
	$username = "root";
	$password = "CSisGreat!7951";
	$server = "localhost";
	$db_name = "capstone";

	$conn = new mysqli($server, $username, $password, $db_name);

	if ($conn->connect_error) {
		die("Connection failed: " . $conn->connect_error);
	}
?>
