package com.company;

public class Main {

    public static void main(String[] args) {
        Company c = new Company();

        Employee Arne = new Employee("Arne", 32);
        Employee Line = new Employee("Line", 36);
        Employee Ragnar = new Employee("Ragnar", 140);

        c.AddEmployees(Arne, Line, Ragnar);
        System.out.println("Alle");
        c.PrintEmployees();
        System.out.println("-------------------------");

        System.out.println("Der fikk Arne fyken");
        c.RemoveEmployee(Arne);
        c.PrintEmployees();
        System.out.println("-------------------------");
        System.out.println("Ragnar og Line får fyken");
        c.RemoveEmployees(Ragnar, Line);
        System.out.println("Resterende medlemmer:");
        c.PrintEmployees();
        System.out.println("-------------------------");

    }
}
