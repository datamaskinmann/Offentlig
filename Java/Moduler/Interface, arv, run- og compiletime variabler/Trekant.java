package Modul11;

public class Trekant extends Form {
    private double lengde;
    private double høyde;

    public Trekant(double lengde, double høyde) {
        super("Trekant");
        this.lengde = lengde;
        this.høyde = høyde;
    }

    public double Areal() {
        return (lengde*høyde)/2;
    }
}