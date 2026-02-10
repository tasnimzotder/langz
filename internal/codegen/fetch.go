package codegen

import (
	"fmt"
	"strings"

	"github.com/tasnimzotder/langz/internal/ast"
)

type fetchOptions struct {
	URL     string
	Method  string
	Body    string
	Headers []headerPair
	Timeout string
	Retries string
}

type headerPair struct {
	Key   string
	Value string
}

// parseFetchOptions extracts fetch configuration from the FuncCall AST node.
// URL comes from the first positional arg; everything else from keyword args.
func (g *Generator) parseFetchOptions(call *ast.FuncCall) fetchOptions {
	opts := fetchOptions{}
	if len(call.Args) > 0 {
		opts.URL = g.genExpr(call.Args[0])
	}
	for _, kw := range call.KwArgs {
		switch kw.Key {
		case "method":
			opts.Method = g.genRawValue(kw.Value)
		case "body":
			opts.Body = g.genExpr(kw.Value)
		case "timeout":
			opts.Timeout = g.genRawValue(kw.Value)
		case "retries":
			opts.Retries = g.genRawValue(kw.Value)
		case "headers":
			if m, ok := kw.Value.(*ast.MapLiteral); ok {
				for i, key := range m.Keys {
					opts.Headers = append(opts.Headers, headerPair{
						Key:   key,
						Value: g.genRawValue(m.Values[i]),
					})
				}
			}
		}
	}
	return opts
}

// genFetchAssignment generates multi-line curl for: name = fetch(...)
func (g *Generator) genFetchAssignment(name string, call *ast.FuncCall) {
	opts := g.parseFetchOptions(call)
	g.emitFetchBlock(opts)
	g.writeln(fmt.Sprintf(`%s="$_body"`, name))
}

// genFetchStatement generates multi-line curl for standalone: fetch(...)
func (g *Generator) genFetchStatement(call *ast.FuncCall) {
	opts := g.parseFetchOptions(call)
	g.emitFetchBlock(opts)
}

// buildCurlCmd assembles the curl command string from fetch options.
func buildCurlCmd(opts fetchOptions) string {
	parts := []string{`curl -s -w "%{http_code}"`}

	if opts.Method != "" {
		parts = append(parts, fmt.Sprintf("-X %s", opts.Method))
	}

	for _, h := range opts.Headers {
		parts = append(parts, fmt.Sprintf(`-H "%s: %s"`, h.Key, h.Value))
	}

	if opts.Body != "" {
		parts = append(parts, fmt.Sprintf("-d %s", opts.Body))
	}

	if opts.Timeout != "" {
		parts = append(parts, fmt.Sprintf("--max-time %s", opts.Timeout))
	}

	parts = append(parts, `-D "$_tmp_headers"`)
	parts = append(parts, `-o "$_tmp_body"`)
	parts = append(parts, opts.URL)

	return strings.Join(parts, " ")
}

// emitCurlCore writes the tmpfile setup, curl call, and cleanup.
func (g *Generator) emitCurlCore(opts fetchOptions) {
	g.writeln(`_tmp_headers=$(mktemp)`)
	g.writeln(`_tmp_body=$(mktemp)`)
	g.writeln(fmt.Sprintf(`_status=$(%s) || true`, buildCurlCmd(opts)))
	g.writeln(`_body=$(cat "$_tmp_body")`)
	g.writeln(`_headers=$(cat "$_tmp_headers")`)
	g.writeln(`rm -f "$_tmp_headers" "$_tmp_body"`)
}

// emitFetchBlock writes the curl block, optionally wrapped in a retry loop.
func (g *Generator) emitFetchBlock(opts fetchOptions) {
	if opts.Retries != "" {
		g.writeln(`_fetch_attempt=0`)
		g.writeln(fmt.Sprintf(`_fetch_max=%s`, opts.Retries))
		g.writeln(`while [ "$_fetch_attempt" -lt "$_fetch_max" ]; do`)
		g.indent++
		g.writeln(`_fetch_attempt=$((_fetch_attempt + 1))`)
		g.emitCurlCore(opts)
		g.writeln(`if [ "$_status" -ge 200 ] && [ "$_status" -lt 300 ]; then`)
		g.indent++
		g.writeln(`break`)
		g.indent--
		g.writeln(`fi`)
		g.writeln(`sleep 1`)
		g.indent--
		g.writeln(`done`)
	} else {
		g.emitCurlCore(opts)
	}
}
