<%@ page import="java.util.Iterator" %><%--
  Created by IntelliJ IDEA.
  User: svenn
  Date: 29/03/2020
  Time: 20:56
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<html>
<head>
    <title>Welcome</title>
</head>
<body>
    <h1>Welcome, <%=request.getAttribute("Username")%></h1>
</body>
</html>
