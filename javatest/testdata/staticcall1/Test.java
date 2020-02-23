package staticcall1;

import testutil.T;

public class Test {
    public static void run(int x) {
        T.printInt(add(x, 1));
        T.printInt(Test.add(x, 1));
    }

    public static int add(int x, int y) {
        return x + y;
    }
}
