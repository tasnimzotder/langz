" Vim syntax file for LangZ
if exists('b:current_syntax')
  finish
endif

" Comments
syntax match langzComment "//.*$" contains=@Spell

" Strings with interpolation
syntax region langzString start=/"/ end=/"/ contains=langzInterpolation
syntax match langzInterpolation /{\w\+}/ contained

" Numbers
syntax match langzNumber /\<[0-9]\+\>/

" Control flow keywords
syntax keyword langzKeyword if elif else for in fn return match continue break while

" Logical operators
syntax keyword langzLogical and or

" Boolean constants
syntax keyword langzBoolean true false

" Builtin functions
syntax match langzBuiltin /\<\(print\|exec\|env\|read\|write\|rm\|mkdir\|copy\|move\|chmod\|chown\|glob\|exit\|fetch\|sleep\|append\|hostname\|whoami\|arch\|dirname\|basename\|is_file\|is_dir\|rmdir\|upper\|lower\|os\|args\|range\|exists\|json_get\|trim\|len\|timestamp\|date\)\>\ze\s*(/

" String methods
syntax match langzMethod /\.\<\(replace\|contains\|starts_with\|ends_with\|split\|join\|length\)\>\ze\s*(/

" Operators
syntax match langzOperator /|>\|=>\|->\|==\|!=\|>=\|<=\|+=\|-=\|\*=\|\/=\|[=+\-*/%<>!]/

" Highlighting
highlight default link langzComment Comment
highlight default link langzString String
highlight default link langzInterpolation Special
highlight default link langzNumber Number
highlight default link langzKeyword Keyword
highlight default link langzLogical Keyword
highlight default link langzBoolean Boolean
highlight default link langzBuiltin Function
highlight default link langzMethod Function
highlight default link langzOperator Operator

let b:current_syntax = 'langz'
