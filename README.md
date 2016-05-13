# parser
--
    import "github.com/MJKWoolnough/parser"

Package parser is a simple helper package for parsing strings, byte slices and
io.Readers

## Usage

```go
var (
	ErrNoState = errors.New("no state")
)
```
Errors

#### type Parser

```go
type Parser struct {
	Tokeniser
}
```


#### func  New

```go
func New(t Tokeniser) Parser
```

#### func (*Parser) Accept

```go
func (p *Parser) Accept(types ...TokenType) bool
```

#### func (*Parser) AcceptRun

```go
func (p *Parser) AcceptRun(types ...TokenType) TokenType
```

#### func (*Parser) Done

```go
func (p *Parser) Done() (Phrase, PhraseFunc)
```

#### func (*Parser) Error

```go
func (p *Parser) Error() (Phrase, PhraseFunc)
```

#### func (*Parser) Except

```go
func (p *Parser) Except(types ...TokenType) bool
```

#### func (*Parser) ExceptRun

```go
func (p *Parser) ExceptRun(types ...TokenType) TokenType
```

#### func (*Parser) Get

```go
func (p *Parser) Get() []Token
```

#### func (*Parser) GetPhrase

```go
func (p *Parser) GetPhrase() (Phrase, error)
```

#### func (*Parser) Len

```go
func (p *Parser) Len() int
```

#### func (*Parser) Peek

```go
func (p *Parser) Peek() Token
```

#### func (*Parser) PhraserState

```go
func (p *Parser) PhraserState(pf PhraseFunc)
```

#### type Phrase

```go
type Phrase struct {
	Type PhraseType
	Data []Token
}
```


#### type PhraseFunc

```go
type PhraseFunc func(*Parser) (Phrase, PhraseFunc)
```


#### type PhraseType

```go
type PhraseType int
```


```go
const (
	PhraseDone PhraseType = -1 - iota
	PhraseError
)
```

#### type Token

```go
type Token struct {
	Type TokenType
	Data string
}
```

Token represents data parsed from the stream.

#### type TokenFunc

```go
type TokenFunc func(*Tokeniser) (Token, TokenFunc)
```

TokenFunc is the type that the worker funcs implement in order to be used by the
tokeniser.

#### type TokenType

```go
type TokenType int
```

TokenType represents the type of token being read.

Negative values are reserved for this package.

```go
const (
	TokenDone TokenType = -1 - iota
	TokenError
)
```
Constants TokenError (-2) and TokenDone (-1)

#### type Tokeniser

```go
type Tokeniser struct {
	Err error
}
```

Tokeniser is

#### func  NewByteTokeniser

```go
func NewByteTokeniser(data []byte) Tokeniser
```
NewByteTokeniser returns a Tokeniser which uses a byte slice.

#### func  NewReaderTokeniser

```go
func NewReaderTokeniser(reader io.Reader) Tokeniser
```
NewReaderTokeniser returns a Tokeniser which uses an io.Reader

#### func  NewStringTokeniser

```go
func NewStringTokeniser(str string) Tokeniser
```
NewStringTokeniser returns a Tokeniser which uses a string.

#### func (*Tokeniser) Accept

```go
func (t *Tokeniser) Accept(chars string) bool
```
Accept returns true if the next character to be read is contained within the
given string.

Upon true, it advances the read position, otherwise the position remains the
same.

#### func (*Tokeniser) AcceptRun

```go
func (t *Tokeniser) AcceptRun(chars string) rune
```
AcceptRun reads from the string as long as the read character is in the given
string.

Returns the rune that stopped the run.

#### func (*Tokeniser) Done

```go
func (t *Tokeniser) Done() (Token, TokenFunc)
```
Done is a TokenFunc that is used to indicate that there are no more tokens to
parse.

#### func (*Tokeniser) Error

```go
func (t *Tokeniser) Error() (Token, TokenFunc)
```
Error represents an error state for the parser.

The error value should be set by calling Tokeniser.SetError and then this func
should be called.

#### func (*Tokeniser) Except

```go
func (t *Tokeniser) Except(chars string) bool
```
Except returns true if the next character to be read is not contained within the
given string. Upon true, it advances the read position, otherwise the position
remains the same.

#### func (*Tokeniser) ExceptRun

```go
func (t *Tokeniser) ExceptRun(chars string) rune
```
ExceptRun reads from the string as long as the read character is not in the
given string.

Returns the rune that stopped the run.

#### func (*Tokeniser) Get

```go
func (t *Tokeniser) Get() string
```
Get returns a string of everything that has been read so far and resets the
string for the next round of parsing.

#### func (*Tokeniser) GetToken

```go
func (t *Tokeniser) GetToken() (Token, error)
```

#### func (*Tokeniser) Len

```go
func (t *Tokeniser) Len() int
```
Len returns the number of bytes that has been read since the last Get.

#### func (*Tokeniser) Peek

```go
func (t *Tokeniser) Peek() rune
```
Peek returns the next rune without advancing the read position.

#### func (*Tokeniser) TokeniserState

```go
func (t *Tokeniser) TokeniserState(tf TokenFunc)
```
