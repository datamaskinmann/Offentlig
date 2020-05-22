package Modul11;

import java.util.ArrayList;

public class Main {
    public static void main(String[] args) {
        Cow cow = new Cow(14);
        Cat cat = new Cat(4);

        Firkant firkant = new Firkant(4);
        Sirkel sirkel = new Sirkel(0.5);
        Trekant trekant = new Trekant(4, 4);

        ArrayList<Printable> printables = new ArrayList<>();
        ArrayList<Form> forms = new ArrayList<>();

        forms.add(firkant);
        forms.add(sirkel);
        forms.add(trekant);

        printables.add(cow);
        printables.add(cat);
        printables.add(firkant);
        printables.add(sirkel);
        printables.add(trekant);

        for(Printable p : printables) {
            Print(p);
        }

        System.out.println("Summen av arealene: " + sumArea(forms));

        Imperial.Yard y = new Imperial.Yard(40);
        System.out.println(y.getAmount() + " yards is " + y.toMeters() + " meters");

        Metric.Meter m = new Metric.Meter(40);
        System.out.println(m.getAmount() + " meters is " + m.toYards() + " yards");

        Metric.Centimeter cm = new Metric.Centimeter(100);
        System.out.println(cm.getAmount() + " cm is " + cm.toInches() + " inches");

        Imperial.Inch i = new Imperial.Inch(20);
        System.out.println(i.getAmount() + " inches is " + i.toCM() + " cm");


    }

    public static <T extends Printable> void Print(T obj) {
        obj.Print();
    }

    public static double sumArea(ArrayList<? extends Form> list) {
        double sumBuffer = 0.0;
        for(Form f : list) {
            sumBuffer+=f.Areal();
        }
        return sumBuffer;
    }
}