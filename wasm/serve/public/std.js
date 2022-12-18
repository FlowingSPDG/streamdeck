// WebAssembly.instantiateStreamingがない場合のポリフィル
if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

// main.wasmにビルドされたGoのプログラムを読み込む
const go = new Go();
let mod, inst;
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(async(result) => {
    mod = result.module;
    inst = result.instance;
    await go.run(inst);
    inst = await WebAssembly.instantiate(mod, go.importObject);
});