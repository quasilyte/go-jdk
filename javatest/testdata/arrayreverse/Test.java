package arrayreverse;

import testutil.T;

public class Test {
    public static void run(int x) {
        test0();
        test1();
        test2();
        test5();
        test7();
    }

    public static void test0() {
        T.printIntArray(reverse(new int[]{}));
    }

    public static void test1() {
        T.printIntArray(reverse(new int[]{1}));
    }

    public static void test2() {
        T.printIntArray(reverse(new int[]{1, 2}));
        T.printIntArray(reverse(new int[]{2, 1}));
    }

    public static void test5() {
        T.printIntArray(reverse(new int[]{1, 2, 3, 4, 5}));
        T.printIntArray(reverse(new int[]{4, 3, 2, 1, 1}));
    }

    public static void test7() {
        T.printIntArray(reverse(new int[]{1, 1, 1, 2, 1, 1, 1}));
    }

    public static int[] reverse(int[] array) {
        for (int i = 0; i < array.length/2; i++) {
            int temp = array[i];
            array[i] = array[array.length-i-1];
            array[array.length-i-1] = temp;
        }
        return array;
    }
}
