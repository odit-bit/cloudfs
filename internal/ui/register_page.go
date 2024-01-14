package ui

import (
	"html/template"
	"io"
)

const register_template = `

<head>
    <meta charset="UTF-8">
    <title>User Registration</title>
</head>
<body>
    <style>
        // .container {
        //     text-align: center;
        //     width: 60%;
        //     margin: auto;
        //     padding: 20px;
        //     background-color: #fff;
        //     border-radius: 8px;
        //     box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        //     margin-top: 50px;
        // }

        form {
            text-align: center;
            margin-top: 20px;
        }

        input[type="text"],
        input[type="password"],
        input[type="email"],
        input[type="submit"] {
            #width: 100%;
            padding: 10px;
            margin-bottom: 15px;
            border-radius: 5px;
            border: 1px solid #ccc;
        }
        input[type="submit"] {
            background-color: #4caf50;
            color: white;
            cursor: pointer;
        }
        input[type="submit"]:hover {
            background-color: #45a049;
        }
    </style>

    <div class="container">
        <h2>Registration</h2>
        <form action={{.}} method="post">
            <label for="username">Username</label>
            <input type="text" id="username" name="username" required>

            <label for="password">Password</label>
            <input type="password" id="password" name="password" required>

            <input type="submit" value="Register">
        </form>
    </div>
</body>

`

var registerPage = template.Must(template.New("register.html").Parse(register_template))

func RenderRegisterPage(w io.Writer, endpoint string) error {
	return registerPage.Execute(w, endpoint)
}
