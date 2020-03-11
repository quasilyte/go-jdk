package testutil;

public class T {
    public static native void printInt(int x);
    public static native void printLong(long x);

    public static native int isub(int x, int y);
    public static native int isub3(int x, int y, int z);

    public static native long ii_l(int a1, int a2);
    public static native int li_i(long a1, int a2);
    public static native int il_i(int a1, long a2);
    public static native int ilil_i(int a1, long a2, int a3, long a4);

    public static native void GC();
}
