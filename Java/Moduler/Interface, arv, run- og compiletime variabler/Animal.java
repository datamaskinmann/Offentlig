package Modul11;

public abstract class Animal extends Printable {

    private String name;
    private String sound;

    private int age;

    public Animal(String name, int age, String sound) {
        this.name = name;
        this.age = age;
        this.sound = sound;
    }

    @Override
    public void Print() {
        System.out.println(name + " says " + sound);
    }
}
