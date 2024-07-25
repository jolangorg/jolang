class P{

    constructor(v){
        console.log(v);
    }

}

class A{
    pcstack = new class extends P {

        constructor(){
            super(456);
        }

        newInstance() {
            console.log("newInstance");
        }

        newArray(intsize) {
            console.log("newArray", intsize);
        }
    };
}

let a = new A;