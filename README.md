# parser
--
    import "github.com/MJKWoolnough/parser"

Package parser is a simple helper package for parsing strings, byte slices and
io.Readers

## Usage

#### type Parser

```go
type Parser struct {
}
```

Parser is the wrapper type for the various different parsers.

#### func  NewByteParser

```go
func NewByteParser(data []byte) Parser
```
NewByteParser returns a Parser which parses a byte slice.

#### func  NewReaderParser

```go
func NewReaderParser(reader io.Reader) Parser
```
NewReaderParser returns a Parser which parses a Reader.

#### func  NewStringParser

```go
func NewStringParser(str string) Parser
```
NewStringParser returns a Parser which parses a string.

#### func (Parser) Accept

```go
func (p Parser) Accept(chars string) bool
```
Accept returns true if the next character to be read is contained within the
given string. Upon true, it advances the read position, otherwise the position
remains the same.

#### func (Parser) AcceptRun

```go
func (p Parser) AcceptRun(chars string)
```
AcceptRun reads from the string as long as the read character is in the given
string.

#### func (Parser) Except

```go
func (p Parser) Except(chars string) bool
```
Except returns true if the next character to be read is not contained within the
given string. Upon true, it advances the read position, otherwise the position
remains the same.

#### func (Parser) ExceptRun

```go
func (p Parser) ExceptRun(chars string)
```
ExceptRun reads from the string as long as the read character is not in the
given string.

#### func (Parser) Get

```go
func (p Parser) Get() string
```
Get returns a string of everything that has been read so far and resets the
string for the next round of parsing.

#### func (Parser) Len

```go
func (p Parser) Len() int
```
Len returns the number of bytes that has been read since the last Get.

#### func (Parser) Peek

```go
func (p Parser) Peek() rune
```
Peek returns the next rune without advancing the read position.
