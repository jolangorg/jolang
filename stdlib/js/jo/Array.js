export class TypedArray extends Array {
    type;
}

export function NewArray(type, dim, ...dims) {
    let arr = new TypedArray(dim);
    arr.type = type;
    if (dims.length === 0){
        return arr;
    }
    for (let i = 0; i < dim; i++){
        let subdim = dims.shift();
        arr[i] = NewArray(type, subdim, ...dims);
    }
    return arr;
}