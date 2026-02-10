" LangZ filetype plugin
if exists('b:did_ftplugin')
  finish
endif
let b:did_ftplugin = 1

setlocal commentstring=//\ %s
setlocal comments=://
setlocal tabstop=4
setlocal shiftwidth=4
setlocal expandtab
setlocal smartindent

let b:undo_ftplugin = 'setlocal commentstring< comments< tabstop< shiftwidth< expandtab< smartindent<'
