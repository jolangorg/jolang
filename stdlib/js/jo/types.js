let buf = new ArrayBuffer(4);
let float32Array = new Float32Array(buf);
let uint32Array = new Float32Array(buf);

export class Float extends Number{

    static isNaN(f){
        return isNaN(f)
    }

    static isInfinite(f){
        return !isFinite(f)
    }

    static floatToIntBits(f){
        float32Array[0] = f;
        return uint32Array[0];
    }

}

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