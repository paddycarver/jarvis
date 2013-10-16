package cli

import (
  "testing"
)

// map inputs to the tokens they're expected to produce
var providedStrings = map[string][]token {
  `docs myfile.html` : []token{
    token{code: tokenCommand, value: "docs"},
    token{code: tokenParam, value: "myfile.html"},
  },
  `docs --markdown myfile.md` : []token{
    token{code: tokenCommand, value: "docs"},
    token{code: tokenArg, value: "markdown"},
    token{code: tokenParam, value: "myfile.md"},
  },
  `docs -markdown myfile.md` : []token{
    token{code: tokenCommand, value: "docs"},
    token{code: tokenArg, value: "markdown"},
    token{code: tokenParam, value: "myfile.md"},
  },
  `docs -format=markdown myfile.md`: []token{
    token{code: tokenCommand, value: "docs"},
    token{code: tokenArg, value: "format=markdown"},
    token{code: tokenParam, value: "myfile.md"},
  },
  `docs -format="markdown file" myfile.md`: []token{
    token{code: tokenCommand, value: "docs"},
    token{code: tokenArg, value: `format="markdown file"`},
    token{code: tokenParam, value: "myfile.md"},
  },
  `docs --format=markdown myfile.md`: []token{
    token{code: tokenCommand, value: "docs"},
    token{code: tokenArg, value: "format=markdown"},
    token{code: tokenParam, value: "myfile.md"},
  },
  `docs --format="markdown file" myfile.md`: []token{
    token{code: tokenCommand, value: "docs"},
    token{code: tokenArg, value: `format="markdown file"`},
    token{code: tokenParam, value: "myfile.md"},
  },
  `docs myfile.md --markdown`: []token{
    token{code: tokenCommand, value: "docs"},
    token{code: tokenParam, value: "myfile.md"},
    token{code: tokenArg, value: "markdown"},
  },
  `docs myfile.md -markdown`: []token{
    token{code: tokenCommand, value: "docs"},
    token{code: tokenParam, value: "myfile.md"},
    token{code: tokenArg, value: "markdown"},
  },
  `docs --server`: []token{
    token{code: tokenCommand, value: "docs"},
    token{code: tokenArg, value: "server"},
  },
  `docs -server`: []token{
    token{code: tokenCommand, value: "docs"},
    token{code: tokenArg, value: "server"},
  },
}

// test the predefined inputs to ensure the produce the predefined tokens
func TestProvidedStrings(t *testing.T) {
  for providedString, expectedTokens := range providedStrings {
    lex := NewLexer(providedString)
    expectedPos := 0
    token := lex.nextToken()
    for token.code > tokenEOF {
      expectedToken := expectedTokens[expectedPos]
      if token.code != expectedToken.code {
        t.Errorf("Expected token of type %q, got token of type %q", expectedToken.code, token.code)
      }
      if token.value != expectedToken.value {
        t.Errorf(`Expected token value to be "%s", got "%s"`, expectedToken.value, token.value)
      }
      token = lex.nextToken()
      expectedPos += 1
    }
    t.Logf("Exited on token %s\n", token)
  }
}

func BenchmarkLexer(b *testing.B) {
  b.StopTimer()
  inputs := make([]string, len(providedStrings))
  pos := 0
  for input, _ := range providedStrings {
    inputs[pos] = input
    pos++
  }
  b.StartTimer()
  for i := 0; i < b.N; i++ {
    lex := NewLexer(inputs[i%len(inputs)])
    lex.Lex()
  }
}
