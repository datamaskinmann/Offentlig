import javax.servlet.annotation.WebServlet;
import javax.servlet.*;
import javax.servlet.http.*;
import java.io.IOException;
import java.sql.*;

@WebServlet(name="/login", urlPatterns = "/login")
public class login extends HttpServlet {
    protected void doPost(HttpServletRequest request, HttpServletResponse response) throws IOException, ServletException {
        String username = request.getParameter("username");
        String password = request.getParameter("password");

        SQLConnection loginDatabaseConnection = new SQLConnection(
                "localhost",
                3306,
                "users",
                "verifACC",
                "VERIFACC123!\"#"
        );

        if(!loginDatabaseConnection.isValid()) response.sendRedirect("ServerError.jsp");

        ResultSet rs = loginDatabaseConnection.doQuery(
                String.format("SELECT * FROM data WHERE username='%s' AND password='%s'", username, password)
        );

        int rowCount = 0;
        try {
            while(rs.next()) rowCount = rs.getRow();
        }catch (SQLException ex) {
            response.sendRedirect("ServerError.jsp");
        }

        if(rowCount != 1) {
            RequestDispatcher requestDispatcher = request.getRequestDispatcher("index.jsp");
            request.setAttribute("loginSuccess", "no");
            requestDispatcher.forward(request, response);
            return;
        }

        request.setAttribute("Username", username);
        request.setAttribute("loginSuccess", "yes");

        RequestDispatcher requestDispatcher = request.getRequestDispatcher("user.jsp");
        requestDispatcher.forward(request, response);
        response.sendRedirect("user.jsp");
    }
}