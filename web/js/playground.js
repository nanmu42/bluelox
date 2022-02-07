// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
In the absence of any formal way to specify interfaces in JavaScript,
here's a skeleton implementation of a playground transport.

        function Transport() {
                // Set up any transport state (eg, make a websocket connection).
                return {
                        Run: function(body, output, options) {
                                // Compile and run the program 'body' with 'options'.
				// Call the 'output' callback to display program output.
                                return {
                                        Kill: function() {
                                                // Kill the running program.
                                        }
                                };
                        }
                };
        }

	// The output callback is called multiple times, and each time it is
	// passed an object of this form.
        var write = {
                Kind: 'string', // 'start', 'stdout', 'stderr', 'end'
                Body: 'string'  // content of write or end status message
        }

	// The first call must be of Kind 'start' with no body.
	// Subsequent calls may be of Kind 'stdout' or 'stderr'
	// and must have a non-null Body string.
	// The final call should be of Kind 'end' with an optional
	// Body string, signifying a failure ("killed", for example).

	// The output callback must be of this form.
	// See PlaygroundOutput (below) for an implementation.
        function outputCallback(write) {
        }
*/

function PlaygroundOutput(el) {
    'use strict';

    // retrieve DOM element from jQuery element
    el = el.get(0)
    return function(write) {
        if (write.Kind === 'start') {
            el.innerHTML = '';
            return;
        }

        var cl = 'system';
        if (write.Kind === 'stdout' || write.Kind === 'stderr') cl = write.Kind;

        var m = write.Body;
        if (write.Kind === 'end') {
            m = '\nProgram exited' + (m ? ': ' + m : '.');
        }

        if (m.indexOf('IMAGE:') === 0) {
            // TODO(adg): buffer all writes before creating image
            var url = 'data:image/png;base64,' + m.substr(6);
            var img = document.createElement('img');
            img.src = url;
            el.appendChild(img);
            return;
        }

        // ^L clears the screen.
        var s = m.split('\x0c');
        if (s.length > 1) {
            el.innerHTML = '';
            m = s.pop();
        }

        m = m.replace(/&/g, '&amp;');
        m = m.replace(/</g, '&lt;');
        m = m.replace(/>/g, '&gt;');

        var needScroll = el.scrollTop + el.offsetHeight === el.scrollHeight;

        var span = document.createElement('span');
        span.className = cl;
        span.innerHTML = m;
        el.appendChild(span);

        if (needScroll) el.scrollTop = el.scrollHeight - el.offsetHeight;
    };
}

(function() {
    function lineHighlight(error) {
        const regex = / at line ([0-9]+)/g
        let r = regex.exec(error)
        while (r) {
            $('.lines div')
                .eq(r[1] - 1)
                .addClass('lineerror');
            r = regex.exec(error);
        }
    }
    function highlightOutput(wrappedOutput) {
        return function(write) {
            if (write.Body) lineHighlight(write.Body);
            wrappedOutput(write);
        };
    }
    function lineClear() {
        $('.lineerror').removeClass('lineerror');
    }

    // opts is an object with these keys
    //  codeEl - code editor element
    //  outputPreEl - program output pre element
    //  runEl - run button element
    //  fmtEl - fmt button element (optional)
    //  toysEl - toys select element (optional)
    //  enableShortcuts - whether to enable shortcuts
    function playground(opts) {
        var code = $(opts.codeEl);

        // autoindent helpers.
        function insertTabs(n) {
            // Without the n > 0 check, Safari cannot type a blank line at the bottom of a playground snippet.
            // See go.dev/issue/49794.
            if (n > 0) {
                document.execCommand('insertText', false, '\t'.repeat(n));
            }
        }
        function autoindent(el) {
            var curpos = el.selectionStart;
            var tabs = 0;
            while (curpos > 0) {
                curpos--;
                if (el.value[curpos] === '\t') {
                    tabs++;
                } else if (tabs > 0 || el.value[curpos] === '\n') {
                    break;
                }
            }
            setTimeout(function() {
                insertTabs(tabs);
            }, 1);
        }

        function keyHandler(e) {
            if (!opts.enableShortcuts) return;

            if (e.keyCode === 9 && !e.ctrlKey) {
                // tab (but not ctrl-tab)
                insertTabs(1);
                e.preventDefault();
                return false;
            }
            if (e.keyCode === 13) {
                // enter
                if (e.ctrlKey) {
                    // +ctrl
                    run();
                    e.preventDefault();
                    return false;
                }
                if (e.altKey) {
                    // +alt
                    fmt();
                    e.preventDefault();
                } else {
                    autoindent(e.target);
                }
            }
            return true;
        }
        code.unbind('keydown').bind('keydown', keyHandler);
        var output = $(opts.outputPreEl).empty();
        window.writeOutput = PlaygroundOutput(output)

        function body() {
            return $(opts.codeEl).val();
        }
        function setBody(text) {
            $(opts.codeEl).val(text);
        }

        function setError(error) {
            lineClear();
            lineHighlight(error);
            output
                .empty()
                .addClass('error')
                .text(error);
        }

        async function sleep(ms) {
            return new Promise((resolve => {
                setTimeout(()=>{
                    resolve()
                }, ms)
            }))
        }

        async function runOnly() {
            try {
                lineClear()
                await window.loxstop()
                output.removeClass('error').text('Running...')
                await sleep(1) // wait for DOM update
                window.writeOutput = highlightOutput(writeOutput)
                await window.loxrun(body())
            } catch (e) {
                setError(e)
            }
        }

        function fmtAnd(run) {
            // TODO: rewrite me
            run();
        }

        function fmt() {
            fmtAnd(function(){});
        }

        function run() {
            fmtAnd(runOnly);
        }

        $(opts.runEl).click(run);
        $(opts.fmtEl).click(fmt);

        if (opts.toysEl !== null) {
            $(opts.toysEl).bind('change', function() {
                var toy = $(this).val();
                $.ajax('/toy/' + toy, {
                    processData: false,
                    type: 'GET',
                    complete: function(xhr) {
                        if (xhr.status !== 200) {
                            alert('Server error; try again.');
                            return;
                        }
                        setBody(xhr.responseText);
                        run();
                    },
                });
            });
        }

        return {

        }
    }

    window.playground = playground;
})();
