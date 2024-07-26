import {PrintStream} from "../io/PrintStream.js";

export class System {

    static out = new class extends PrintStream{
        printf() {
            console.log(...arguments);
        }
    }

    static nanoTime(){
        return (performance.now() * 1000 + performance.timeOrigin * 1000);
    }

}