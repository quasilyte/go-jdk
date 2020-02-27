package bench;

import benchutil.B;

public class Bench {
    public static void nop() {
    }

    public static void callGoNop() {
        B.nop();
        B.nop();
        B.nop();
        B.nop();
        B.nop();
    }
}
