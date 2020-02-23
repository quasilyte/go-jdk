package arith1;

import testutil.T;

public class Test {
    public static void run(int x) {
        T.printInt(x);

        T.printInt(x+1);
        x++;
        T.printInt(x);
    }
}
