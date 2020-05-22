package Modul11;

public interface LengthUnit {
    default double toMeters(){return 1.0;}
    default double toYards(){return 1.0;}
    default double toCM(){return 1.0;}
    default double toInches(){return 1.0;}
    default double toFeet(){return 1.0;}
}