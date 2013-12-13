package cli

import (
	"errors"
  "github.com/paddyforan/jarvis/parse"
	"github.com/paddyforan/jarvis/spec"
  "io"
  "os"
	"strings"
)

var (
	MissingResourceDirError = errors.New("Missing resource directory.")
	MultipleCommandError    = errors.New("Too many commands were passed in the input.")
	NoTokensError           = errors.New("No tokens supplied to command routing.")
	NoCommandError          = errors.New("Input must begin with a command.")
	NoInputError            = errors.New("No input provided.")
	UnknownCommandError     = errors.New("Command not recognised.")
	UnknownTokenError       = errors.New("A token of an unknown type was encountered.")
)

type LexingError string

func (e LexingError) Error() string {
	return string(e)
}

type argMap map[string][]interface{}

func Run(input string) error {
  lex := NewLexer(input)
  tokens, err := lex.Lex()
  if err != nil {
    return err
  }
  return route(tokens)
}

func route(tokens []token) error {
	if len(tokens) < 1 {
		return NoTokensError
	}
	if tokens[0].code == tokenError {
		return LexingError(tokens[0].value)
	}
	if tokens[0].code == tokenEOF {
		return NoInputError
	}
	if tokens[0].code != tokenCommand {
		return NoCommandError
	}
	if tokens[len(tokens)-1].code == tokenEOF {
		tokens = tokens[:len(tokens)-1] // we no longer need EOF
	}
	command := strings.ToLower(tokens[0].value)
	args, err := parseInput(tokens[1:])
	if err != nil {
		return err
	}
	switch command {
	case "spec":
		err = serveSpec(args)
		return err
	case "jsonschema":
		err = serveJSONSchema(args)
		return err
	default:
		return UnknownCommandError
	}
	return nil
}

func parseInput(tokens []token) (argMap, error) {
	if len(tokens) < 1 {
		return argMap{}, nil
	}
	results := argMap{}
	for i := 0; i < len(tokens); i++ {
		switch tokens[i].code {
		case tokenParam:
			results[""] = append(results[""], tokens[i].value)
			continue
		case tokenArg:
      parts := strings.SplitN(tokens[i].value, "=", 2)
			arg := strings.ToLower(parts[0])
      // TODO: Error if there is only one part
      // TODO: Fix lexing so a = is not required
			results[arg] = append(results[arg], parts[1])
			i++
      // TODO: Distinguish between args and flags
			// results[arg] = append(results[arg], true)
			continue
		case tokenEOF:
			return results, nil
		case tokenError:
			return argMap{}, LexingError(tokens[i].value)
		case tokenCommand:
			return argMap{}, MultipleCommandError
		default:
			return argMap{}, UnknownTokenError
		}
	}
	return results, nil
}

type server func(format string, output io.WriteCloser, resources []*parse.Resource) error

func serve(f server, args argMap) error {
	if len(args[""]) < 1 {
		return MissingResourceDirError
	}
  if len(args["root"]) < 1 {
    resourceDir, err := os.Getwd()
		if err != nil {
			return err
		}
    args["root"] = append(args["root"], resourceDir)
	}
  root := args["root"][0].(string)
	resources := map[string]*parse.Resource{}
  for _, dir := range args[""] {
    if root[len(root)-1] != os.PathSeparator {
      root = root + string(os.PathSeparator)
    }
		rmap, err := parse.Parse(root, dir.(string))
    if err != nil {
      return err
    }
		for k, v := range rmap {
			resources[k] = v
		}
	}
	rslice := make([]*parse.Resource, len(resources))
	for _, resource := range resources {
		rslice = append(rslice, resource)
	}
  if len(args["format"]) < 1 {
    args["format"] = append(args["format"], "markdown")
  }
  output := os.Stdout
  var err error
  if len(args["output"]) > 1 {
    output, err = os.Open(args["output"][0].(string))
    if err != nil {
      return err
    }
  }
  return f(args["format"][0].(string), output, rslice)
}

func serveSpec(args argMap) error {
  return serve(spec.Generate, args)
}

func serveJSONSchema(args argMap) error {
  return nil
}
