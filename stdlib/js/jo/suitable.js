import {numericTypes} from "./types.js";

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

        if (args[i].constructor !== t){
            return false
        }
    }

    return true
}