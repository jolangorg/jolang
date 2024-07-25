String.prototype.equals = function(s){
    return this.toString() === s.toString();
}

Object.prototype.hashCode = function(){
    return this.toString();
}



function assert(cond) {
    if (!cond){
        throw new Error("assert failed");
    }
}