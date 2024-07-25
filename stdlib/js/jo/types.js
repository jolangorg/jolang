import {Float} from "java/lang/Float.js";

export function float(v){
    return v * 1;
}

window.Float = Float;

export class Integer extends Number{
}

export function int(v){
    return Math.trunc(v);
}

export function boolean(v){
    return Boolean(v)
}

export const numericTypes = [
    float,
    int,
]