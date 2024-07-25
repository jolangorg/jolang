import {numericTypes, boolean} from "./types.js";

export function suitable(args, ...types){
    if (args.length !== types.length){
        return false;
    }
    let i = 0;
    for (let t of types){
        if (numericTypes.indexOf(t) > -1){
            if (typeof args[i] === "number"){
                i++;
                continue;
            }else{
                return false
            }
        }

        if (t === boolean){
            if (typeof args[i] === "boolean"){
                i++;
                continue;
            }else{
                return false
            }
        }

        if (!(args[i] instanceof t)) {
            return false
        }
        i++;

        // if (args[i].constructor !== t){
        //     return false
        // }
    }

    return true
}