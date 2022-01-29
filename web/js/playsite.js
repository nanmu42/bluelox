window.addEventListener('DOMContentLoaded', async () => {
    initEditor()
    await loadBlueLox()

    window.playground && window.playground({
        'codeEl': '.js-playgroundCodeEl',
        'outputPreEl': '.js-playgroundOutputPreEl',
        'runEl': '.js-playgroundRunEl',
        'fmtEl': '.js-playgroundFmtEl',
        'toysEl': '.js-playgroundToysEl',
        'enableShortcuts': true,
    })

    function initEditor() {
        const code = $('#code')
        code.linedtextarea()
        code.attr('wrap', 'off')
        code.resize(function () {
            code.linedtextarea()
        })
    }

    async function loadBlueLox() {
        if (!WebAssembly.instantiateStreaming) {
            WebAssembly.instantiateStreaming = async (resp, importObject) => {
                const source = await (await resp).arrayBuffer()
                return await WebAssembly.instantiate(source, importObject)
            }
        }

        function loadWasm(path) {
            const go = new Go()

            return new Promise((resolve, reject) => {
                WebAssembly.instantiateStreaming(fetch(path), go.importObject)
                    .then(result => {
                        go.run(result.instance)
                        resolve(result.instance)
                    })
                    .catch(error => {
                        reject(error)
                    })
            })
        }

        try {
            await loadWasm('/js/bluelox.wasm')

        } catch (e) {
            console.log('loading Bluelox wasm:', e)
        }

        console.log('Bluelox wasm is loaded 👋')
    }
})