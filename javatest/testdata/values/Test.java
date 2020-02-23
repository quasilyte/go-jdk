package values;

import testutil.T;

public class Test {
    public static void run(int x) {
        T.printInt(x);

        T.printInt(0);
        T.printInt(100);
        T.printInt(-100);
        T.printInt(0xff);
        T.printInt(0xffffff);
        T.printInt(-0xffffff);

        int v1 = 120;
        int v2 = -120;
        int v3 = -139438;
        T.printInt(v1);
        T.printInt(v2);
        T.printInt(v3);
    }
}
