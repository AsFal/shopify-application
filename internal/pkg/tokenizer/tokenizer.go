package tokenizer

type Tokenizer interface {
	Process(string) ([]string, error)
}
