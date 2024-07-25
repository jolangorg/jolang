import {numericTypes, int, float, boolean} from "./types.js";
import {Enum} from "./Enum.js";

export {int, float, boolean};

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

export function suitable(args, ...types){
    if (args.length !== types.length){
        return false;
    }
    let i = 0;
    for (let t of types){
        if (numericTypes.indexOf(t) > -1){
            if (typeof args[i] === "number"){
                continue;
            }else{
                return false
            }
        }
        if (t === boolean){
            if (typeof args[i] === "boolean"){
                continue;
            }else{
                return false
            }
        }
        args[i].constructor
    }
    let result = false;

}

export {Enum};