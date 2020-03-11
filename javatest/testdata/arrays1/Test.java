package arrays1;

import testutil.T;

public class Test {
    public static void run(int x) {
        // We have a few forced GC runs to ensure that
        // our objects (arrays) are not collected.

        int[] xs = new int[4];
        T.printInt(xs.length);

        T.GC();

        T.printInt(xs[0]);
        xs[0] = 15;
        T.printInt(xs[0]);

        T.printInt(xs[1]);
        xs[1] = xs[0];
        T.printInt(xs[1]);

        T.printInt(xs.length);

        T.GC();

        int[] ys = new int[10];
        ys[0] = 1;
        ys[1] = 2;
        ys[3] = 3;
        T.printInt(ys[0]);
        T.printInt(ys[1]);
        T.printInt(ys[2]);
        T.printInt(ys[9]);

        T.GC();

        int[] zs = makeArray();
        zs[0] = -1;
        zs[1] = -2;
        zs[3] = -3;
        T.printInt(zs[0]);
        T.printInt(zs[1]);
        T.printInt(zs[2]);
        T.printInt(zs[9]);

        // Run operations again after GC.
        T.GC();

        T.printInt(xs.length);

        T.printInt(xs[0]);
        xs[0] = 15;
        T.printInt(xs[0]);

        T.printInt(xs[1]);
        xs[1] = xs[0];
        T.printInt(xs[1]);

        T.printInt(xs.length);

        T.printInt(ys[0]);
        T.printInt(ys[1]);
        T.printInt(ys[2]);
        T.printInt(ys[9]);

        T.printInt(zs[0]);
        T.printInt(zs[1]);
        T.printInt(zs[2]);
        T.printInt(zs[9]);
    }

    private static int[] makeArray() {
        return new int[10];
    }

    // TODO: uncomment when variable indexes are supported.
    // private static void indexVarAndSize(int length) {
    //     int[] a = new int[length];
    //     T.printInt(a.length);
    //     int i1 = 0;
    //     T.printInt(a[i1]);
    //     a[i1] = 10;
    //     T.printInt(a[i1]);
    //     int i2 = 5;
    //     T.printInt(a[i2]);
    //     a[i2] = 20;
    //     T.printInt(a[i2]);
    // }
}
