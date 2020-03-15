package staticfinal;

import testutil.T;

// TODO: implement putstatic and getstatic.

public class Test {
    private static final int SEVEN = 7;
    private static final int TEN = SEVEN + 3;

    // private static final int[] INTS = {1, 2, 3};

    public static void run(int x) {
        T.printInt(SEVEN + TEN);
        // T.printInt(INTS[0]);
    }
}
