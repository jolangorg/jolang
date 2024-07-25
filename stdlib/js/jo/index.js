String.prototype.equals = function(s){
    return this.toString() === s.toString();
}

function assert(cond) {
    if (!cond){
        throw new Error("assert failed");
    }
}