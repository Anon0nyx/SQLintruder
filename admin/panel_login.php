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
		if (isset($_POST)) {
			$username = $_POST['username'];
			$password = $_POST['password'];

			if (strlen($password) != 0) {
				$query = "SELECT * FROM users WHERE username='" . $username . "' AND password = '" . $password . "';";
				$result = $conn->query($query);

				if ($result -> num_rows > 0) {
					while ($row = $result->fetch_assoc()) {
						echo "<h3>Successful Login</h3>";
						echo "<h1>Welcome: " . $row['username'] . "</h1>";
						echo "<h2>Currently your profile information is:</h2>";
						echo "<h3>ID: " . $row['id'] . "</h3>";
						echo "<h3>Password: " . $row['password'] . "</h3>";
					}
				} else {
					echo "<h1>Login Failed..!</h1>";
		
				}
			} else {
				echo "<h1>Admin Panel Account Login</h1>";
			}
		}
	?>
    <form  method="POST">
		<input name="username" placeholder="Username" type="text">
		<input name="password" placeholder="Password" type="text">
		<input name="submitBtn" type="submit" value="Submit">
	</form>
	<a href="./register.php">Register</a>    
</body>
</html> 
