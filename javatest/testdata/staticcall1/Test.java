package staticcall1;

import testutil.T;

public class Test {
    public static void run(int x) {
        T.printInt(add(x, 1));
        T.printInt(Test.add(x, 1));

        T.printInt(iabs(0));
        T.printInt(iabs(1));
        T.printInt(iabs(-1));
        T.printInt(iabs(13924));

        T.printInt(iabs2(0));
        T.printInt(iabs2(1));
        T.printInt(iabs2(-1));

        T.printInt(iabs3(0));
        T.printInt(iabs3(1));
        T.printInt(iabs3(-1));

        T.printLong(labs(0));
        T.printLong(labs(1));
        T.printLong(labs(-1));
        T.printLong(labs(-15));
        T.printLong(labs(13));
    }

    public static int add(int x, int y) {
        return x + y;
    }

    public static int iabs(int x) {
        if (x < 0) {
            return -x;
        }
        return x;
    }

    public static int iabs2(int x) {
        if (x < 0) {
            x = -x;
        }
        return x;
    }

    public static int iabs3(int x) {
        int result = 0;
        if (x < 0) {
            result = -x;
        } else {
            result = x;
        }
        return result;
    }

    public static long labs(long x) {
        // TODO: rewrite as ternary expr after we fix a bug
        // that generates incorrect IR for it. See #37.
        // return (x < 0) ? -x : x;

        if (x < 0) {
            return -x;
        }
        return x;
    }
}
