package loops1;

import testutil.T;

public class Test {
    public static void run(int x) {
        int v = 0;
        for (int i = 0; i < 100; i++) {
            if (i < 50) {
                v++;
            }
        }
        T.printInt(v);

        v = 101;
        int i = 0;
        while (i < 10) {
            T.printInt(i);
            v -= i;
            i += 3;
        }
        T.printInt(v);

        i = 110;
        do {
            i--;
        } while (i >= 100);
        T.printInt(i);
    }
}
