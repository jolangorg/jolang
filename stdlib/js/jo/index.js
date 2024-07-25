import {Float} from "./Float.js";
import {Enum} from "./Enum.js";

window.Float = Float;

String.prototype.equals = function(s){
    return this.toString() === s.toString();
}

Object.prototype.hashCode = function(){
    return this.toString();
}

export function assert(cond) {
    if (!cond){
        throw new Error("assert failed");
    }
}

export {Enum};
export {Float};