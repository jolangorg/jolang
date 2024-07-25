export class Float extends Number{

    static isNaN(f){
        return isNaN(f)
    }

    static isInfinite(f){
        return !isFinite(f)
    }

    floatToIntBits(f){
        var buf = new ArrayBuffer(4);
        (new Float32Array(buf))[0] = f;
        return (new Uint32Array(buf))[0];
    }

}