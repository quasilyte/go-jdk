package testutil;

import java.util.Arrays;

public class T {
    public static void printInt(int x) {
        System.out.println(x);
    }

    public static void printLong(long x) {
        System.out.println(x);
    }

    public static void printIntArray(int[] xs) {
        System.out.println(Arrays.toString(xs));
    }

    public static int isub(int x, int y) {
        return x - y;
    }

    public static int isub3(int x, int y, int z) {
        return x - y - z;
    }

    public static long ii_l(int a1, int a2) {
        return a1 - a2;
    }

    public static int li_i(long a1, int a2) {
        return (int)a1 - a2;
    }

    public static int il_i(int a1, long a2) {
        return a1 - (int)a2;
    }

    public static int ilil_i(int a1, long a2, int a3, long a4) {
        return a1 - (int)a2 - a3 - (int)a4;
    }

    public static void GC() {}
}
