import {$import} from "./tools.js";

let B; $import("./B.js").then($module => B = $module.B);

export class A{

    getB(){
        return new B()
    }

}