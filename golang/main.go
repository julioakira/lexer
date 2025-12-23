package main

import (
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type CustomError struct {
	timestamp time.Time
	position  string
	error     string
}

const (
  // structural tokens
	TKN_BRACE_OPEN     rune   = '{'
	NAME_BRACE_OPEN    string = "BRACE_OPEN"
	TKN_BRACE_CLOSE    rune   = '}'
	NAME_BRACE_CLOSE   string = "BRACE_CLOSE"
	TKN_BRACKET_OPEN   rune   = '['
	NAME_BRACKET_OPEN  string = "BRACKET_OPEN"
	TKN_BRACKET_CLOSE  rune   = ']'
	NAME_BRACKET_CLOSE string = "BRACKET_CLOSE"
	TKN_COLON          rune   = ':'
	NAME_COLON         string = "COLON"
	TKN_COMMA          rune   = ','
	NAME_COMMA         string = "COMMA"
	TKN_QUOTE          rune   = '"'
	NAME_QUOTE         string = "QUOTE"
  // literals tokens
	STRING = "STRING"
	NUMBER = "NUMBER"
  // primitives tokens
	TRUE  = "TRUE"
	FALSE = "FALSE"
	NULL  = "NULL"
)

type Token struct {
	value          string
	name           string
	start_position int
	end_position   int
	width          int
}

func validateNumber(number string) bool {

	if number[0] == '-' {
		if utf8.RuneCountInString(number) == 0 || number[1] == '0' {
			return false
		}
	}

	if number[0] == '0' {
		return false
	}
	return true
}

func tokenizer(input string) ([]Token, *CustomError) {
	runes := []rune(input)
	length := len(runes)
	idx := 0
	tokens := []Token{}

	if length == 0 {
		return tokens, &CustomError{
			timestamp: time.Now(),
			position:  "0",
			error:     "No input to parse",
		}
	}

	for idx < length {
		char := runes[idx]

		if unicode.IsSpace(char) {
			idx++
			continue
		}

		switch char {

		case TKN_BRACE_OPEN:
			tokens = append(tokens, Token{string(TKN_BRACE_OPEN), NAME_BRACE_OPEN, idx, idx + 1, 1})
			idx++

		case TKN_BRACE_CLOSE:
			tokens = append(tokens, Token{string(TKN_BRACE_CLOSE), NAME_BRACE_CLOSE, idx, idx + 1, 1})
			idx++

		case TKN_BRACKET_OPEN:
			tokens = append(tokens, Token{string(TKN_BRACKET_OPEN), NAME_BRACKET_OPEN, idx, idx + 1, 1})
			idx++

		case TKN_BRACKET_CLOSE:
			tokens = append(tokens, Token{string(TKN_BRACKET_CLOSE), NAME_BRACKET_CLOSE, idx, idx + 1, 1})
			idx++

		case TKN_COLON:
			tokens = append(tokens, Token{string(TKN_COLON), NAME_COLON, idx, idx + 1, 1})
			idx++

		case TKN_COMMA:
			tokens = append(tokens, Token{string(TKN_COMMA), NAME_COMMA, idx, idx + 1, 1})
			idx++

		case TKN_QUOTE:
			start := idx
			idx++

			for idx < length {
				if runes[idx] == '\\' {
					idx += 2
					continue
				}
				if runes[idx] == '"' {
					break
				}
				idx++
			}

			if idx >= length {
				return nil, &CustomError{
					timestamp: time.Now(),
					position:  string(fmt.Sprint(start)),
					error:     "Unterminated string",
				}
			}

			value := string(runes[start+1 : idx])
			tokens = append(tokens, Token{
				value:          value,
				name:           STRING,
				start_position: start + 1,
				end_position:   idx,
				width:          idx - start - 1,
			})

			idx++

		default:
			remaining := string(runes[idx:])

			if strings.HasPrefix(remaining, "true") {
				tokens = append(tokens, Token{"true", TRUE, idx, idx + 4, 4})
				idx += 4
				continue
			}

			if strings.HasPrefix(remaining, "false") {
				tokens = append(tokens, Token{"false", FALSE, idx, idx + 5, 5})
				idx += 5
				continue
			}

			if strings.HasPrefix(remaining, "null") {
				tokens = append(tokens, Token{"null", NULL, idx, idx + 4, 4})
				idx += 4
				continue
			}

			if unicode.IsDigit(char) || char == '-' {
				start := idx
				idx++

				isFloat := false
				isExp := false
				expDigits := 0

				for idx < length {
					c := runes[idx]

					switch {
					case unicode.IsDigit(c):
						idx++
						if isExp {
							expDigits++
						}

					case c == '.':
						if isFloat || isExp {
							return nil, &CustomError{
								timestamp: time.Now(),
								position:  string(fmt.Sprint(idx)),
								error:     "Invalid number format",
							}
						}
						isFloat = true
						idx++

					case c == 'e' || c == 'E':
						if isExp {
							return nil, &CustomError{
								timestamp: time.Now(),
								position:  string(fmt.Sprint(idx)),
								error:     "Multiple exponents not allowed",
							}
						}
						isExp = true
						idx++

						if idx < length && (runes[idx] == '+' || runes[idx] == '-') {
							idx++
						}
						expDigits = 0

					default:
						goto numberDone
					}
				}

			numberDone:
				number := string(runes[start:idx])

				if !validateNumber(number) {
					return nil, &CustomError{
						timestamp: time.Now(),
						position:  string(fmt.Sprint(start)),
						error:     "Invalid number",
					}
				}

				if isExp && expDigits == 0 {
					return nil, &CustomError{
						timestamp: time.Now(),
						position:  string(fmt.Sprint(idx)),
						error:     "Exponent requires digits",
					}
				}

				tokens = append(tokens, Token{
					value:          number,
					name:           NUMBER,
					start_position: start,
					end_position:   idx,
					width:          idx - start,
				})
				continue
			}

			return nil, &CustomError{
				timestamp: time.Now(),
				position:  string(fmt.Sprint(idx)),
				error:     "Unexpected character",
			}
		}
	}

	return tokens, nil
}

func main() {
	raw := `{"array":[1,2,3],"boolean":true,"color":"gold","null":null,"number":123,"object":{"a":"b","c":"d"},"string":"Hello World"}`
	fmt.Printf("input Len: %d\n", utf8.RuneCountInString(raw))
	tokenized, err := tokenizer(raw)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", tokenized)
}
