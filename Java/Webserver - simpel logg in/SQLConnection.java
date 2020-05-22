import java.sql.*;

public class SQLConnection {

    private String URL;
    private int Port;
    private String DatabaseName;
    private String Username;
    private String Password;
    private boolean isValid;

    private Connection conn;

    /**
     *
     * @param URL the url of the database
     * @param Port the port of the database
     * @param DatabaseName the database-name
     * @param Username the username
     * @param Password the password
     */
    public SQLConnection(String URL, int Port, String DatabaseName, String Username, String Password) {
        this.URL = URL;
        this.Port = Port;
        this.DatabaseName = DatabaseName;
        this.Username = Username;
        this.Password = Password;

        try{
            Class.forName("com.mysql.jdbc.Driver");
            this.conn = DriverManager.getConnection(
                    String.format("jdbc:mysql://%s:%d/%s",
                            this.URL,
                            this.Port,
                            this.DatabaseName), this.Username, this.Password);
            isValid = true;
        } catch (ClassNotFoundException | SQLException ex) {
            System.out.println(
                    ex instanceof ClassNotFoundException ?
                            "Could not find class:\n" :
                            "Error while connecting to SQL:\n"
                    + ex.getMessage()
            );
        }
    }

    public boolean isValid() {
        return isValid;
    }

    /**
     *
     * @param Query the SQL query you want to execute
     * @return a String array with the Query result
     */
    public ResultSet doQuery(String Query) {
        try {
            Statement statement = this.conn.createStatement();
            return statement.executeQuery(Query);
        } catch (SQLException ex) {
            System.out.println("Error while executing query:" + ex.getMessage());
            return null;
        }
    }
}
