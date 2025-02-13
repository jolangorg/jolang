String.prototype.equals = function(s){
    return this.toString() === s.toString();
}

Object.prototype.hashCode = function(){
    return this.toString();
}

export function assert(cond, msg = "assert failed") {
    if (!cond){
        throw new Error(msg);
    }
}

import {int, float, boolean} from "./types.js";
export {int, float, boolean};

import {suitable} from "./suitable.js";
export {suitable};

import {Enum} from "./Enum.js";
export {Enum};

import {NewArray} from "./Array.js";
export {NewArray}

import {Interface} from "./Interface.js";
export {Interface}

//Imports
const importPromises = [];

export function Imports(cb) {
    let importPromise = cb();
    importPromises.push(importPromise);
}

export function Import(importPromise){
    importPromises.push(importPromise);
}

export async function waitImports() {
    await Promise.all(importPromises);
}

export async function init() {
    await waitImports();
}
