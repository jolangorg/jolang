export class Enum {

    name;
    constructor(name) {
        this.name = name;
    }

    toString(){
        return this.name;
    }

    static values(){
        return Object.values(this)
    }
}