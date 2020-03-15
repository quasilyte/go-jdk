package arith1;

import testutil.T;

public class Test {
    public static void run(int x) {
        T.printInt(x+1);
        x++;
        x += 1;
        x = x + 1;
        T.printInt(x);

        T.printInt(x-1);
        x--;
        x -= 1;
        x = x - 1;
        T.printInt(x);

        int a = 10;
        int b = 3;
        T.printInt(a-b);
        T.printInt(b-a);
        T.printInt(a+b);
        T.printInt(b+a);

        int y = 1000;
        T.printInt(a/b);
        T.printInt(b/a);
        T.printInt(y/a/b);
        T.printInt(a/2);
        T.printInt(a/3);
        T.printInt(a/5);

        T.printInt(a*b);
        T.printInt(b*a*y);
    }
}
