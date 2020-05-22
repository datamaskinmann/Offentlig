package Modul11;

public class Sirkel extends Form  {

    private double Radius;


    public Sirkel(double Radius) {
        super("Sirkel");
        this.Radius = Radius;
    }

    public double getRadius() {
        return this.Radius;
    }

    public double Areal() {
        return Math.PI*Math.pow(Radius, 2);
    }

    @Override
    public boolean equals(Object obj) {
        if(obj.getClass() != Sirkel.class) return false;
        return ((Sirkel)obj).getRadius() == Radius;
    }
}