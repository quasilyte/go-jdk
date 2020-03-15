package arrays2;

import testutil.T;

public class Test {
    public static void run(int x) {
        int[] lit1 = new int[]{1, 2, 3};
        T.printIntArray(lit1);

        int n = lit1.length;
        T.printInt(n);
        for (int i = 0; i < n; i++) {
            T.printInt(lit1[i]);
        }

        arrayLoop1(lit1);
        arrayLoop2(lit1);
        arrayLoop3(lit1);
    }

    public static void arrayLoop1(int[] xs) {
        int n = xs.length;
        T.printInt(xs.length);
        T.printInt(n);
        for (int i = 0; i < n; i++) {
            T.printInt(xs[i]);
        }
    }

    public static void arrayLoop2(int[] xs) {
        for (int i = 0; i < xs.length; i++) {
            T.printInt(xs[i]);
        }
    }

    public static void arrayLoop3(int[] xs) {
        for (int x : xs) {
            T.printInt(x);
        }
    }
}
