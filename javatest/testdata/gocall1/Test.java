package gocall1;

import testutil.T;

public class Test {
    public static void run(int x) {
        T.printInt(T.isub(0, 0));
        T.printInt(T.isub(1, 0));
        T.printInt(T.isub(0, 1));

        T.printInt(T.isub3(0, 0, 0));
        T.printInt(T.isub3(1, 0, 0));
        T.printInt(T.isub3(0, 1, 0));
        T.printInt(T.isub3(0, 0, 1));
        T.printInt(T.isub3(1, 2, 3));
        T.printInt(T.isub3(3, 2, 1));
        T.printInt(T.isub3(3, 3, 3));
    }
}
