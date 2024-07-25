export class Float extends Number{

    static isNaN(v){
        return isNaN(v)
    }

    static isInfinite(v){
        return !isFinite(v)
    }

}