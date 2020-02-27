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
}
