package Modul11;

abstract class Form extends Printable {

    private String formNavn;

    public Form(String formNavn) {
        this.formNavn = formNavn;
    }

    public abstract double Areal();

    protected String getFormNavn() {
        return formNavn;
    }

    @Override
    public void Print() {
        System.out.println(getFormNavn() + " has an area of " + Areal());
    }
}