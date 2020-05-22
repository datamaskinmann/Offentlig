package Modul11;

public class Firkant extends Form {
    private double sideLengde;

    public Firkant(double sideLengde) {
        super("Firkant");
        this.sideLengde = sideLengde;
    }


    @Override
    public double Areal() {
        return Math.pow(sideLengde, 2);
    }
}