export class Enum {

    name;

    constructor(name) {
        this.name = name;
    }

    toString(){
        return this.name;
    }

    ordinal(){
        return this.constructor.values().indexOf(this);
    }

    static values(){
        return Object.values(this)
    }
}