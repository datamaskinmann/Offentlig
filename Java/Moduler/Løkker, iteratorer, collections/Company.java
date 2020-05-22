package com.company;

import java.util.ArrayList;
import java.util.Iterator;
import java.util.stream.Stream;

public class Company {
    // Collection: ArrayList
    private ArrayList<Employee> employees = new ArrayList<>();
    public Company() {}

    public Company(Employee... employees) {
        for(Employee e : employees) { // For-each loop
            this.employees.add(e);
        }
    }

    public void AddEmployee(Employee e) {
        employees.add(e);
    }

    public void AddEmployees(Employee... e) {
        Stream<Employee> employeeStream = Stream.of(e);
        employeeStream.forEach(x -> employees.add(x)); // Variasjon av for loop med lambda
    }

    public void RemoveEmployee(Employee e) {
        employees.remove(e);
    }

    public void RemoveEmployees(Employee... e) {
        Iterator<Employee> employeeIterator = employees.iterator(); // Iterator

        while(employeeIterator.hasNext()) { // While loop
            Employee emp = employeeIterator.next();
            for(Employee param : e) {
                if(emp.equals(param)) employeeIterator.remove();
            }
        }
    }

    public void PrintEmployees() {
        for(int i = 0; i < employees.size(); i++) { // Standard for loop
            System.out.println(employees.get(i).getName() + " is " + employees.get(i).getAge() + " years old");
        }
    }
}
