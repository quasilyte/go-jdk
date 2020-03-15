package bubblesort;

import testutil.T;

public class Test {
    public static void run(int x) {
        test0();
        test1();
        test2();
        test5();
    }

    public static void test0() {
        int[] a0 = new int[]{};
        sort(a0);
        T.printIntArray(a0);
    }

    public static void test1() {
        int[] a1 = new int[]{1};
        sort(a1);
        T.printIntArray(a1);
    }

    public static void test2() {
        int[] a2sorted = new int[]{1, 2};
        sort(a2sorted);
        T.printIntArray(a2sorted);

        int[] a2 = new int[]{2, 1};
        sort(a2);
        T.printIntArray(a2);
    }

    public static void test5() {
        int[] a5sorted = new int[]{1, 2, 3, 4, 5};
        sort(a5sorted);
        T.printIntArray(a5sorted);

        int[] a5 = new int[]{3, 2, 1, 5, 4};
        T.printIntArray(a5);
        sort(a5);
        T.printIntArray(a5);
    }

    public static void sort(int[] arr) {
        int n = arr.length;
        for (int i = 0; i < n-1; i++) {
            for (int j = 0; j < n-i-1; j++) {
                if (arr[j+1] <= arr[j]) {
                    int temp = arr[j];
                    arr[j] = arr[j+1];
                    arr[j+1] = temp;
                }
            }
        }
    }
}
