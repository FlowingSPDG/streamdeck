if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

function fetchLocal(url) {
    return new Promise(function(resolve, reject) {
      var xhr = new XMLHttpRequest
      xhr.onload = function() {
        resolve(new Response(xhr.response))
      }
      xhr.onerror = function() {
        reject(new TypeError('Local request failed'))
      }
      xhr.open('GET', url)
      xhr.responseType = 'blob'
      xhr.send(null)
    })
  }

// main.wasmにビルドされたGoのプログラムを読み込む
const go = new Go();
let mod, inst;
WebAssembly.instantiateStreaming(fetchLocal('main.wasm'), go.importObject).then(async(result) => {
    mod = result.module;
    inst = result.instance;
    await go.run(inst);
    inst = await WebAssembly.instantiate(mod, go.importObject);
});

// wasm.go と連携している
var port = 0,
    uuid = "",
    registerEventName = "",
    Info = {},
    actionInfo = {}

function connectElgatoStreamDeckSocket(inPort, inPropertyInspectorUUID, inRegisterEvent, inInfo, inActionInfo){
    port = parseInt(inPort)
    uuid = inPropertyInspectorUUID
    registerEventName = inRegisterEvent
    Info = inInfo
    actionInfo = inActionInfo
}