package gocall2;

import testutil.T;

public class Test {
    public static void run(int x) {
        T.printLong(T.ii_l(1, -10));
        T.printLong(T.ii_l(-10, 1));
        T.printLong(T.ii_l(-2, -2));

        T.printInt(T.li_i(1000, 2000));
        T.printInt(T.li_i(2000, 1000));

        T.printInt(T.il_i(1, -10));
        T.printInt(T.il_i(-10, 1));
        T.printInt(T.il_i(-2, -2));

        T.printInt(T.ilil_i(10, 4, -24, -4));
        T.printInt(T.ilil_i(2, 8, 1, 1));
        T.printInt(T.ilil_i(1, 1, 8, 2));
    }
}
