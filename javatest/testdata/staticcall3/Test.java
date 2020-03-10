package staticcall3;

import testutil.T;

public class Test {
    public static void run(int x) {
        // TODO: uncomment a test when long slots are fixed.

        T.printInt(ii_i(1000, 2000));
        //T.printInt(il_i(1000, 2000));
        T.printInt(li_i(1000, 2000));
        T.printLong(ii_l(1000, 2000));
    }

    private static int ii_i(int a1, int a2) {
        T.printInt(a1);
        T.printInt(a2);
        return 23;
    }

    private static int il_i(long a1, int a2) {
        T.printLong(a1);
        T.printInt(a2);
        return 1;
    }

    private static int li_i(int a1, long a2) {
        T.printInt(a1);
        T.printLong(a2);
        return 2;
    }

    private static long ii_l(int a1, int a2) {
        T.printInt(a1);
        T.printInt(a2);
        return 3;
    }
}
