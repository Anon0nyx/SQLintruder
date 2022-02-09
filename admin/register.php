<?php
	include_once("./database.php");
	$password = "";
	
	if (isset($_POST)) {
		$username = $_POST['username'];
		$password = $_POST['password'];

		if (strlen($username) != 0) {
			$query = "SELECT * FROM users WHERE username = '" . $username . "'";

			$result = $conn->query($query);

			if ($result->num_rows > 0) {
				while ($row = $result->fetch_assoc()) {
					if ($row['username'] == $username) {
						echo "<h1>Username Taken, Please Try Another...</h1>";
					}
				}
			} else {
				$query = "INSERT INTO users (username, password) VALUES ('" . $username . "','" . $password . "');";

				$result = $conn->query($query);

				if ($result == TRUE) {
					echo "Registered Successfully";
				} else {
					echo "Registration Failed";
				}
			}
		} else {
			echo "<h1>Registration</h1>";
		} 
	}
?>
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Register</title>
	<style>
		body {
			text-align: center;
			float: center;
		}
		form {
			margin-left: 30%;
			float: center;
		}
		input {
			margin: 10px 45% 5px 0%;
			float: center;
		}
	</style>
</head>
<body>
    <form  method="POST">
		<input name="username" placeholder="Username" type="text">
		<input name="password" placeholder="Password" type="text">
		<input name="submitBtn" type="submit" value="Submit">
	</form>
	<a href="./panel_login.php">Login</a> 
</body>
</html> 
