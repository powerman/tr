var wasmSupported = (typeof WebAssembly === "object");
if (wasmSupported) {
    if (!WebAssembly.instantiateStreaming) { // polyfill
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }
    fetch("/app.wasm").then(function(res) {
        if (res.ok) {
            const go = new Go();
            WebAssembly.instantiateStreaming(res, go.importObject).then((result) => {
                go.run(result.instance);
            });
        } else {
            res.text().then(function(txt) {
                var el = document.getElementById("vugu_mount_point");
                el.style = 'font-family: monospace; background: black; color: red; padding: 10px';
                el.innerText = txt;
            })
        }
    })
} else {
    document.getElementById("vugu_mount_point").innerHTML = 'This application requires WebAssembly support. Please upgrade your browser.';
}
