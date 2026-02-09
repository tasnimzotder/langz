package lsp

// builtinDocs maps builtin function names to their markdown documentation.
var builtinDocs = map[string]string{
	// I/O
	"print":  "```\nprint(args...)\n```\nPrint values to stdout.\n\nTranspiles to `echo`.",
	"write":  "```\nwrite(path, content)\n```\nWrite content to a file.\n\nTranspiles to `echo content > path`.",
	"append": "```\nappend(path, content)\n```\nAppend content to a file.\n\nTranspiles to `echo content >> path`.",
	"read":   "```\nread(path) -> string\n```\nRead file contents.\n\nTranspiles to `$(cat path)`.",

	// File operations
	"rm":    "```\nrm(path)\n```\nRemove a file.\n\nTranspiles to `rm -f path`.",
	"rmdir": "```\nrmdir(path)\n```\nRemove a directory recursively.\n\nTranspiles to `rm -rf path`.",
	"mkdir": "```\nmkdir(path)\n```\nCreate a directory (with parents).\n\nTranspiles to `mkdir -p path`.",
	"copy":  "```\ncopy(src, dst)\n```\nCopy a file.\n\nTranspiles to `cp src dst`.",
	"move":  "```\nmove(src, dst)\n```\nMove/rename a file.\n\nTranspiles to `mv src dst`.",
	"chmod": "```\nchmod(path, mode)\n```\nChange file permissions.\n\nTranspiles to `chmod mode path`.",
	"glob":  "```\nglob(pattern) -> list\n```\nExpand a glob pattern.\n\nTranspiles to `(pattern)`.",

	// File checks
	"exists":  "```\nexists(path) -> bool\n```\nCheck if a path exists.\n\nTranspiles to `[ -e path ]`.",
	"is_file": "```\nis_file(path) -> bool\n```\nCheck if path is a regular file.\n\nTranspiles to `[ -f path ]`.",
	"is_dir":  "```\nis_dir(path) -> bool\n```\nCheck if path is a directory.\n\nTranspiles to `[ -d path ]`.",

	// Execution
	"exec": "```\nexec(command) -> string\n```\nExecute a shell command and capture output.\n\nTranspiles to `$(command)`.",
	"exit": "```\nexit(code)\n```\nExit the script with a status code.\n\nTranspiles to `exit code`.",

	// Environment
	"env": "```\nenv(name) -> string\n```\nGet an environment variable.\n\nTranspiles to `\"${NAME}\"`.",

	// System info
	"os":       "```\nos() -> string\n```\nGet OS name (lowercase).\n\nTranspiles to `$(uname -s | tr '[:upper:]' '[:lower:]')`.",
	"arch":     "```\narch() -> string\n```\nGet CPU architecture.\n\nTranspiles to `$(uname -m)`.",
	"hostname": "```\nhostname() -> string\n```\nGet machine hostname.\n\nTranspiles to `$(hostname)`.",
	"whoami":   "```\nwhoami() -> string\n```\nGet current username.\n\nTranspiles to `$(whoami)`.",

	// Path utilities
	"dirname":  "```\ndirname(path) -> string\n```\nGet directory part of a path.\n\nTranspiles to `$(dirname path)`.",
	"basename": "```\nbasename(path) -> string\n```\nGet filename part of a path.\n\nTranspiles to `$(basename path)`.",

	// String utilities
	"upper": "```\nupper(str) -> string\n```\nConvert string to uppercase.\n\nTranspiles to `$(echo str | tr '[:lower:]' '[:upper:]')`.",
	"lower": "```\nlower(str) -> string\n```\nConvert string to lowercase.\n\nTranspiles to `$(echo str | tr '[:upper:]' '[:lower:]')`.",

	// Networking
	"fetch": "```\nfetch(url) -> string\n```\nHTTP GET a URL.\n\nTranspiles to `$(curl -sf url)`.",

	// Misc
	"args":  "```\nargs() -> list\n```\nGet script arguments.\n\nTranspiles to `(\"$@\")`.",
	"range": "```\nrange(start, end) -> list\n```\nGenerate a numeric sequence.\n\nTranspiles to `$(seq start end)`.",
	"sleep": "```\nsleep(seconds)\n```\nPause execution for N seconds.\n\nTranspiles to `sleep N`.",
}
