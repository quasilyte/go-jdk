package irgen;

class C1 {
    // slots=0
    //   b0 Iret 10
    public static int m1() {
        return 10;
    }

    // slots=2
    //   b0 flags = Icmp r0 0
    //   b0 JumpGtEq label0 flags
    //   b1 r1 = Ineg r0
    //   b1 Iret r1
    // label0:
    //   b2 Iret r0
    public static int abs(int x) {
        if (x < 0) {
            return -x;
        }
        return x;
    }

    // slots=3
    //   b0 flags = Lcmp r0 0
    //   b0 JumpGtEq label0 flags
    //   b1 r2 = Lneg r0
    //   b1 Jump label1
    // label0:
    //   b2 r2 = Lload r0
    // label1:
    //   b3 Lret r2
    public static long labs(long x) {
        return (x < 0) ? -x : x;
    }

    // slots=4
    //   b0 flags = Icmp r0 1
    //   b0 JumpGt label0 flags
    //   b1 Iret r0
    // label0:
    //   b2 r1 = Isub r0 1
    //   b2 r2 = CallStatic fib r1
    //   b2 r1 = Isub r0 2
    //   b2 r3 = CallStatic fib r1
    //   b2 r1 = Iadd r2 r3
    //   b2 Iret r1
    public static int fib(int n) {
        if (n <= 1) {
            return n;
        }
        return fib(n-1) + fib(n-2);
    }

    // slots=3
    //   b0 flags = Icmp r0 1
    //   b0 JumpGtEq label0 flags
    //   b1 Iret 1
    // label0:
    //   b2 r1 = Isub r0 1
    //   b2 r2 = CallStatic factorial r1
    //   b2 r1 = Imul r0 r2
    //   b2 Iret r1
    public static int factorial(int n) {
        if (n < 1) {
            return 1;
        }
        return n * factorial(n-1);
    }

    // slots=4
    //   b0 r1 = Iload 0
    // label0:
    //   b1 flags = Icmp r0 0
    //   b1 JumpLt label1 flags
    //   b2 r2 = Isub r0 r1
    //   b2 r0 = Iload r2
    //   b2 r1 = Iadd r1 1
    //   b2 r2 = Isub r0 r1
    //   b2 r0 = Iload r2
    //   b2 Jump label0
    // label1:
    //   b3 r2 = Iload r1
    //   b3 r3 = Isub r2 1
    //   b3 Iret r3
    public static int sqrt(int n) {
        int b = 0;
        while (n >= 0) {
            n = n - b;
            b++;
            n = n - b;
        }
        return b - 1;
    }

    // slots=1
    //   b0 r0 = NewIntArray 10
    //   b0 Aret r0
    public static int[] newIarray() {
        return new int[10];
    }

    // slots=2
    //   b0 r0 = Iload 128
    //   b0 r1 = NewDoubleArray r0
    //   b0 Aret r1
    public static double[] newDarray() {
        int length = 128;
        return new double[length];
    }
}
