
// Generated automatically by java_test.go.
// This entry point is used by a host Java implementation.
class Main {
    public static void main(String args[]) {
        switch (args[0]) {
        case "arith1":
            arith1.Test.run(400);
            return;
        case "staticcall1":
            staticcall1.Test.run(0);
            return;
        default:
            System.out.println("unknown package: " + args[0]);
        }
    }
}
