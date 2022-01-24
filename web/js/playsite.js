window.addEventListener('DOMContentLoaded', () => {
    const code = $('#code')
    code.linedtextarea()
    code.attr('wrap', 'off')
    code.resize(function () {
        code.linedtextarea()
    })
})