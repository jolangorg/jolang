import * as jo from "jo";

let B;
let C;
// jo.Imports(
//     import("./B.js").then($m => B = $m.B),
//     import("./C.js").then($ => C = $.C),
// );

jo.Imports(async () => {
    B = (await import("./B.js")).B;
    C = (await import("./C.js")).C;
});

export class A{

    getB(){
        return new B()
    }

    getC(){
        return new C()
    }
}