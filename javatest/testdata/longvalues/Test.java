package longvalues;

import testutil.T;

public class Test {
    public static void run(int x) {
        T.printLong(x);

        T.printLong(0);
        T.printLong(100);
        T.printLong(-100);
        T.printLong(0xff);
        T.printLong(0xffffff);
        T.printLong(-0xffffff);

        long v1 = 120;
        long v2 = -120;
        long v3 = -139438;
        long v4 = 192394122;
        T.printLong(v1);
        T.printLong(v2);
        T.printLong(v3);
        T.printLong(v4);
    }
}
