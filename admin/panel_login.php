<?php
	include_once("./database.php");
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Panel Login</title>
	<style>
		body {
			float: center;
			text-align: center;
		}
		form {
			margin-left: 30%;
			float: center;
		}
		input {
			margin: 10px 45% 5px 0%;
			float: center;	
	</style>
</head>
<body>
	<?php
		if (isset($_POST) && !empty($_POST)) {
			$username = $_POST['username'];
			$password = $_POST['password'];

			if (strlen($password) != 0) {
				$query = "SELECT * FROM users WHERE username='" . $username . "' AND password = '" . $password . "';";
				$result = $conn->query($query);

				if ($result -> num_rows > 0) {
					echo "<h1>Login Successful!</h1>";
					echo "\n<!-- JSONSTART";
					while ($row = $result->fetch_assoc()) {
						echo "[ userDetails {'id':'" . $row['id'] . "'},{'username':'" . $row['username'] . "'},{'password':'" . $row['password'] . "'} ]";
					}
					echo "JSONEND -->\n";
				} else {
					echo "<h1>Login Failed..!</h1><br>";
		
				}
			}
		} else {
			echo "<h1>Account Panel Login</h1><br>";
		}
	?>
    <form  method="POST">
		<input name="username" placeholder="Username" type="text">
		<input name="password" placeholder="Password" type="password">
		<input name="submitBtn" type="submit" value="Submit">
	</form>
	<a href="./register.php">Register</a>    
</body>
</html> 
