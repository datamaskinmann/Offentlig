package Modul11;

public final class Metric {
    public Metric() {
        throw new UnsupportedOperationException("Cannot instantiate Metric");
    }

    static String[] a = {"90", "50"};



    public static class Centimeter implements LengthUnit {
        private double amount;

        public Centimeter(double amount) {
            this.amount = amount;
        }

        public double getAmount() { return amount; }

        @Override
        public double toInches() {
            return amount*0.393700787;
        }

        @Override
        public double toFeet() {
            return amount*0.032808399;
        }

        @Override
        public double toMeters() {  return amount*0.01; }

        @Override
        public double toYards() { return amount*0.010936133; }
    }

    public static class Meter implements LengthUnit {
        private double amount;

        public Meter(double amount) {
            this.amount = amount;
        }

        public double getAmount() { return amount; }

        @Override
        public double toCM() { return amount*100; }

        @Override
        public double toYards() { return amount*1.0936133; }

        @Override
        public double toFeet() { return amount*3.2808399; }

        @Override
        public double toInches() { return amount*39.3700787; }
    }
}
