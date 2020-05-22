<%--
  Created by IntelliJ IDEA.
  User: svenn
  Date: 29/03/2020
  Time: 18:42
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<html class="noScroll">
    <head>
        <link rel="stylesheet" href="style.css">
        <script src="jquery-3.4.1.min.js"></script>
        <title>Title</title>
    </head>

    <body class="noScroll" style="background: url(bgimg.png) no-repeat;background-size: 100%;">
            <div id="userLogin" class="centerText">
                 <h1 style="width: 70%; top: 25%; font-size: 3vw; border-bottom: 0.1vw solid black;">LOG IN</h1>
                <form action="login" method="post">
                    <input id="unm" class="stdIn" type="text" name="username" value="Username..."><br>
                    <input id="psw" class="stdIn" type="password" name="password" value="Password..."><br>
                    <input id="login" type="submit" value="LOG IN">
                </form>
            </div>
    </body>
    <script type="text/javascript">
        document.getElementById("login").disabled = true;

        $("#unm").focusin(function () {
            document.getElementById("unm").setAttribute("value", "");
            document.getElementById("login").disabled = false;
        });

        $("#psw").focusin(function () {
            document.getElementById("psw").setAttribute("value", "");
            document.getElementById("login").disabled = false;
        });

        let loginSuccess = "<%=request.getAttribute("loginSuccess")%>";
        if(loginSuccess == "no") alert("Login failed! Please try again <:^(");
    </script>
</html>
