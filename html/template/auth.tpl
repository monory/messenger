<!DOCTYPE html>
<html>
<head>
    <title>Messenger</title>
    <link rel="icon" type="image/x-icon" href="/static/favicon.ico"/>
    <link rel="stylesheet" type="text/css" href="/static/auth.css"/>
    <meta charset="utf-8"/>
</head>
<body>
    <div class="wrapper">
        <form action="auth" method="post">
            <input type="text" placeholder="Username" name="username" class="input-text"/>
            <input type="password" placeholder="Password" name="password" class="input-text"/>
            {{with $x := .Err}}
            <h1 class="error">{{$x}}</h1>
            {{end}}
            <input type="submit" value="Log in" name="loginbutton" class="input-button"/>
            <input type="submit" value="Register" name="registerbutton" class="input-button"/>
        </form>
    </div>
</body>
</html>
