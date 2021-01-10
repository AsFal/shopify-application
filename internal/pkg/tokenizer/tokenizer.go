package tokenizer

type Tokenizer interface {
	process(string) []string
}
