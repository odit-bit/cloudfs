package ui

import (
	"html/template"
	"io"
)

// const loginPage_template = `
// <!DOCTYPE html>
// <html>
// <head>
// 	<meta name="viewport" content="width=device-width, initial-scale=1.0">
// 	<title>cloudfs-login</title>
// </head>
// <body>
// 	<form id="authForm">
// 		<label for="username">Username:</label><br>
// 		<input type="text" id="username" name="username"><br>
// 		<label for="password">Password:</label><br>
// 		<input type="password" id="password" name="password"><br><br>
// 		<input type="submit" value="Login">
// 	</form>

// 	<script>
// 	document.getElementById("authForm").addEventListener("submit", function(event) {
// 		event.preventDefault();

// 		var username = document.getElementById("username").value;
// 		var password = document.getElementById("password").value;

// 		var url =  {{.}}; // URL injection in JavaScript

// 		fetch(url, {
// 			method: 'POST',
// 			headers: {
// 				'Content-Type': 'application/json',
// 				'Authorization': 'Basic ' + btoa(username + ':' + password)
// 			}
// 		})
// 		.then(response => {
// 			if (response.ok) {
// 				console.log('Login successful');
// 				window.location.href = '/'; // Redirect to "/" endpoint upon successful login
// 			} else if (response.status === 401) {
// 				console.error('Unauthorized:', response.status);
// 				window.location.href = '/login-error.html'; // Redirect to login error page for unauthorized access
// 			} else {
// 				console.error('Error:', response.status);
// 				window.location.href = '/error.html'; // Redirect to a general error page for other errors
// 			}
// 		})
// 		.catch(error => {
// 			console.error('Error:', error);
// 			window.location.href = '/error.html'; // Redirect to a general error page for fetch failures
// 		});
// 	});
// 	</script>
// </body>
// </html>
// `

const loginPage_template = `
<!DOCTYPE html>
<html>
<head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>cloudfs-login</title>
	
</head>
<body>
	<style>
		form {
			text-align: center;
			margin-top: 20px;
		}
	</style>

	<div id="loginForm" >
		<form  hx-encoding='multipart/form-data' hx-post={{.}} hx-swap="innerHtml" >
			<label for="username">Username:</label><br>
			<input type="text" id="username" name="username"><br>
			<label for="password">Password:</label><br>
			<input type="password" id="password" name="password"><br><br>
			<input type="submit" value="Login">

			<button hx-get="/register" hx-target="#loginForm" hx-swap="outerHtml" >register</button>
		
			
		</form>
		
	</div>
	


	
	<script>
		document.body.addEventListener('htmx:beforeSwap', function(evt) {
			var status = evt.detail.xhr.status
			if(status === 404){
				// alert the user when a 404 occurs (maybe use a nicer mechanism than alert())
				// alert("Error: Could Not Find Username");
				evt.detail.shouldSwap = true;
				evt.detail.isError = false;
			} else if ( status === 401){
				// allow 422 responses to swap as we are using this as a signal that
				// a form was submitted with bad data and want to rerender with the
				// errors
				//
				// set isError to false to avoid error logging in console
				evt.detail.shouldSwap = true;
				evt.detail.isError = false;
			} else {
				// window.location.href = "/login";
				evt.detail.shouldSwap = true;
				evt.detail.isError = false;
			}
		});	
	</script>
</body>
</html>
`

var loginPage = template.Must(template.New("loginPage.html").Parse(loginPage_template))

func RenderLoginPage(w io.Writer, authEndpoint string) error {
	return loginPage.Execute(w, authEndpoint)
}
