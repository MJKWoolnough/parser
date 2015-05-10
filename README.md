# strparse
--
    import "github.com/MJKWoolnough/strparse"

Package strparse is a simple helper package for parsing strings

## Usage

#### type Parser

```go
type Parser struct {
	Str string
	// contains filtered or unexported fields
}
```

Parser is a helper with aids with the parsing of formatted strings.

#### func  New

```go
func New(s string) *Parser
```
New returns a new Parser type containg the given string.

#### func (*Parser) Accept

```go
func (p *Parser) Accept(chars string) bool
```
Accept returns true if the next character to be read is contained within the
given string. Upon true, it advances the read position, otherwise the position
remains the same.

#### func (*Parser) AcceptRun

```go
func (p *Parser) AcceptRun(chars string)
```
AcceptRun reads from the string as long as the read character is in the given
string.

#### func (*Parser) Except

```go
func (p *Parser) Except(chars string) bool
```
Except returns true if the next character to be read is not contained within the
given string. Upon true, it advances the read position, otherwise the position
remains the same.

#### func (*Parser) ExceptRun

```go
func (p *Parser) ExceptRun(chars string)
```
ExceptRun reads from the string as long as the read character is not in the
given string.

#### func (*Parser) Get

```go
func (p *Parser) Get() string
```
Get returns a string of everything that has been read so far and resets the
string for the next round of parsing.

#### func (*Parser) Left

```go
func (p *Parser) Left() int
```
Left returns how much of the string is left in the Parser. This includes
everything read since the last Get.

#### func (*Parser) Len

```go
func (p *Parser) Len() int
```
Len returns the current length of the read string.

#### func (*Parser) Peek

```go
func (p *Parser) Peek() rune
```
Peek returns the next rune without advancing the read position.
