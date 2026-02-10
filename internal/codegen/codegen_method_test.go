package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethodReplace(t *testing.T) {
	output := body(compile(`
name = "hello world"
result = name.replace("world", "go")
print(result)
`))

	assert.Contains(t, output, `"${name//world/go}"`)
}

func TestMethodContains(t *testing.T) {
	output := body(compile(`
name = "hello world"
if name.contains("world") {
  print("found")
}
`))

	assert.Contains(t, output, `[[ "$name" == *"world"* ]]`)
}

func TestMethodStartsWith(t *testing.T) {
	output := body(compile(`
path = "/usr/bin"
if path.starts_with("/usr") {
  print("yes")
}
`))

	assert.Contains(t, output, `[[ "$path" == "/usr"* ]]`)
}

func TestMethodEndsWith(t *testing.T) {
	output := body(compile(`
file = "script.sh"
if file.ends_with(".sh") {
  print("shell script")
}
`))

	assert.Contains(t, output, `[[ "$file" == *".sh" ]]`)
}
