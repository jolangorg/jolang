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

    floatToIntBits(f){
        float32Array[0] = f;
        return uint32Array[0];
    }

}