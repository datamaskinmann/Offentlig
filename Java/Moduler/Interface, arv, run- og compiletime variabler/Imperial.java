package Modul11;

public final class Imperial {
    public static class Inch implements LengthUnit {
        private double amount;

        public Inch(double amount) {
            this.amount = amount;
        }

        public double getAmount() { return amount; }

        @Override
        public double toMeters() { return amount*0.0254; }

        @Override
        public double toYards() { return amount*0.0277777778; }

        @Override
        public double toCM() { return amount*2.54; }

        @Override
        public double toFeet() { return amount*0.0833333333; }

        @Override
        public double toInches() { return 1.0; }
    }

    public static class Yard implements LengthUnit {
        private double amount;

        public Yard(double amount) {
            this.amount = amount;
        }

        public double getAmount() { return amount; }

        @Override
        public double toMeters() { return amount*0.9144; }

        @Override
        public double toCM() { return amount*91.44; }

        @Override
        public double toInches() { return amount*36; }

        @Override
        public double toFeet() { return amount*3; }
    }
}
