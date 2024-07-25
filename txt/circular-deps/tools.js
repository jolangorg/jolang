const promises = [];

export async function $import(path) {
    let promise = import(path);
    promises.push(promise);
    return await promise;
}

export async function waitImports() {
    await Promise.all(promises);
}