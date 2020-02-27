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

        T.printInt(iabs4(0));
        T.printInt(iabs4(1));
        T.printInt(iabs4(-1));

        T.printLong(labs(0));
        T.printLong(labs(1));
        T.printLong(labs(-1));
        T.printLong(labs(-15));
        T.printLong(labs(13));
        T.printLong(lucky7());

        T.printLong(labs2(0));
        T.printLong(labs2(1));
        T.printLong(labs2(-1));
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

    public static int iabs4(int x) {
        return (x < 0) ? -x : x;
    }

    public static long labs(long x) {
        if (x < 0) {
            return -x;
        }
        return x;
    }

    public static long labs2(long x) {
        return (x < 0) ? -x : x;
    }

    public static long lucky7() { return 777; }
}
