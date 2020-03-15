package eratosthenes;

import testutil.T;

public class Test {
    private static final int TRUE = 1;
    private static final int FALSE = 0;

    public static void run(int n) {
        // TODO: change int array to booleans when its implemented.
        // TODO: use *2 instead of *two when immediate imul is ready.

        int prime[] = new int[n+1];
        for (int i = 0; i < n; i++) {
            prime[i] = TRUE;
        }

        for (int p = 2; p * p <= n; p++) {
            if (prime[p] == TRUE) {
                int two = 2;
                for (int i = p * two; i <= n; i += p) {
                    prime[i] = FALSE;
                }
            }
        }

        for (int i = 2; i <= n; i++) {
            if (prime[i] == TRUE) {
                T.printInt(i);
            }
        }
    }
}
