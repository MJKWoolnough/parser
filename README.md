# parser
--
    import "vimagination.zapto.org/parser"

Package parser is a simple helper package for parsing strings, byte slices and
io.Readers

## Usage

```go
var (
	ErrNoState      = errors.New("no state")
	ErrUnknownError = errors.New("unknown error")
)
```
Errors

#### type Parser

```go
type Parser struct {
	Tokeniser
}
```

Parser is a type used to get tokens or phrases (collection of token) from an an
input

#### func  New

```go
func New(t Tokeniser) Parser
```
New creates a new Parser from the given Tokeniser

#### func (*Parser) Accept

```go
func (p *Parser) Accept(types ...TokenType) bool
```
Accept will accept a token with one of the given types, returning true if one is
read and false otherwise.

#### func (*Parser) AcceptRun

```go
func (p *Parser) AcceptRun(types ...TokenType) TokenType
```
AcceptRun will keep Accepting tokens as long as they match one of the given
types.

It will return the type of the token that made it stop.

#### func (*Parser) Done

```go
func (p *Parser) Done() (Phrase, PhraseFunc)
```
Done is a PhraseFunc that is used to indicate that there are no more phrases to
parse.

#### func (*Parser) Error

```go
func (p *Parser) Error() (Phrase, PhraseFunc)
```
Error represents an error state for the phraser.

The error value should be set in Parser.Err and then this func should be called.

#### func (*Parser) Except

```go
func (p *Parser) Except(types ...TokenType) bool
```
Except will Accept a token that is not one of the types given. Returns true if
it Accepted a token.

#### func (*Parser) ExceptRun

```go
func (p *Parser) ExceptRun(types ...TokenType) TokenType
```
ExceptRun will keep Accepting tokens as long as they do not match one of the
given types.

It will return the type of the token that made it stop.

#### func (*Parser) Get

```go
func (p *Parser) Get() []Token
```
Get retrieves a slice of the Tokens that have been read so far.

#### func (*Parser) GetPhrase

```go
func (p *Parser) GetPhrase() (Phrase, error)
```
GetPhrase runs the state machine and retrieves a single Phrase and possibly an
error

#### func (*Parser) GetToken

```go
func (p *Parser) GetToken() (Token, error)
```
GetToken runs the state machine and retrieves a single Token and possibly an
error.

If a Token has already been 'peek'ed, that token will be returned without
running the state machine

#### func (*Parser) Len

```go
func (p *Parser) Len() int
```
Len returns how many tokens have been read.

#### func (*Parser) Peek

```go
func (p *Parser) Peek() Token
```
Peek takes a look at the upcoming Token and returns it.

#### func (*Parser) PhraserState

```go
func (p *Parser) PhraserState(pf PhraseFunc)
```
PhraserState allows the internal state of the Phraser to be set.

#### type Phrase

```go
type Phrase struct {
	Type PhraseType
	Data []Token
}
```

Phrase represents a collection of tokens that have meaning together

#### type PhraseFunc

```go
type PhraseFunc func(*Parser) (Phrase, PhraseFunc)
```

PhraseFunc is the type that the worker types implement in order to be used by
the Phraser

#### type PhraseType

```go
type PhraseType int
```

PhraseType represnts the type of phrase being read.

Negative values are reserved for this package.

```go
const (
	PhraseDone PhraseType = -1 - iota
	PhraseError
)
```
Constants PhraseError (-2) and PhraseDone (-1)

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

Tokeniser is a state machine to generate tokens from an input

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

The error value should be set in Tokeniser.Err and then this func should be
called.

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
GetToken runs the state machine and retrieves a single token and possible an
error

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
TokeniserState allows the internal state of the Tokeniser to be set
