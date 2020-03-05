package staticcall2;

import testutil.T;

public class Test {
    public static void run(int x) {
        T.printInt(fib(0));
        T.printInt(fib(1));
        T.printInt(fib(2));
        T.printInt(fib(3));
        T.printInt(fib(4));

        T.printInt(factorial(0));
        T.printInt(factorial(1));
        T.printInt(factorial(2));
        T.printInt(factorial(3));
        T.printInt(factorial(4));
        T.printInt(factorial(5));
        T.printInt(factorial(6));

        T.printInt(sqrt(0));
        T.printInt(sqrt(3));
        T.printInt(sqrt(4));
        T.printInt(sqrt(9));
        T.printInt(sqrt(36));
        T.printInt(sqrt(96));
        T.printInt(sqrt(100));
        T.printInt(sqrt(256));
    }

    public static int fib(int n) {
        if (n <= 1) {
            return n;
        }
        return fib(n-1) + fib(n-2);
    }

    public static int factorial(int n) {
        if (n < 1) {
            return 1;
        }
        return n * factorial(n-1);
    }

    public static int sqrt(int n) {
        int b = 0;
        while (n >= 0) {
            n = n - b;
            b++;
            n = n - b;
        }
        return b - 1;
    }
}
