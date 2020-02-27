package staticcall2;

import testutil.T;

public class Test {
    public static void run(int x) {
        T.printInt(fib(0));
        T.printInt(fib(1));
        T.printInt(fib(2));
        T.printInt(fib(3));
        T.printInt(fib(4));
    }

    public static int fib(int n) {
        if (n <= 1) {
            return n;
        }
        return fib(n-1) + fib(n-2);
    }

    // TODO: uncomment when we have multiplication support.
    // public static int factorial(int n) {
    //     if (n == 0) {
    //         return 1;
    //     }
    //     return n * factorial(n-1);
    // }
}
