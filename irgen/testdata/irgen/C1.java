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
}
