# Validator

## About

This is a lightweight extendable validator library for Golang structs.

The library is based on golang tags. It comes with a set of generic stdlib
validation functions and provides a straightforward way to develop your own
validation logic.

## Example

```go
package main

import (
    "github.com/osdrv/validator"

    apiv1
)

type Message struct {
    Id      apiv1.MessageId `validate:"apiv1_message_id"`
    Title   string          `validate:"nonempty, maxlen(255)"`
    Version int             `validate:"gt(0)"`
    UserId  apiv1.UserId    `validate:"apiv1_user_id"`
    Kind    string          `validate:"enum(text, audio, video)"`
    ReplyTo apiv1.MessageId `validate:"optional, apiv1_message_id"`
}

func main() {
    var message Message
    if err := validator.Validate(message); err != nil {
        // Handle validation error
    }
}
```

```go
package apiv1

import (
    "github.com/osdrv/validator"
)

type MessageId = string
type UserId    = string

var (
    MessageIdRegex = regexp.MustCompile("^[_\\-0-9a-zA-Z]{32}:[_\\-0-9a-zA-Z]{31}$")
    UserIdRegex    = regexp.MustCompile("^[_\\-0-9a-zA-Z]{32}$")
)

func ValidateMessageId(v string) (bool, string) {
	return MessageIdRegex.MatchString(v), "does not look like message Id"
}

func ValidateUserId(v string) (bool, string) {
	return UserIdRegex.MatchString(v), "does not look like user Id"
}

func init() {
	validator.Register("apiv1_message_id", ValidateMessageId)
	validator.Register("apiv1_user_id", ValidateUserId)
}

```

## Built-in functions

### Stringer interface

A non-primitive type would be probed for implementing a stringer interface:

```go
type stringer interface {
    String() string
}
```

### STDLib

| Function handle | Accepted arguments | Details |
| --------------- | ------------------ | ------- |
| empty           | No arguments
| enum            | A list of bools, ints (including: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr), strings and stringer interface| |
| eq              | A single argument of type: int(all the flavors above), bool (casted to string), string and stringer interface | |
| gt              | A single argument of type: int(all the flavors above), bool (casted to string), string and stringer interface | |
| gte             | A single argument of type: int(all the flavors above), bool (casted to string), string and stringer interface | |
| len             | A single string or stringer interface | |
| lt              | A single argument of type: int(all the flavors above), bool (casted to string), string and stringer interface | |
| lte             | A single argument of type: int(all the flavors above), bool (casted to string), string and stringer interface | |
| maxlen          | A single int argument | |
| ne              | A single argument of type: int(all the flavors above), bool (casted to string), string and stringer interface | |
| none            | A list of bools, ints (including: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr), strings and stringer interface| |
| nonempty        | A single argument of type: int(all the flavors above), bool (casted to string), string and stringer interface
| optional        | No arguments
| range           | A list of bools, ints (including: int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr), strings and stringer interface| |

## Implementing a custom validation function

### Validator function interface

A validator func should conform to one of the interfaces below:

```go
func (v <Any>, extra ...<Any>) bool
func (v <Any>, extra ...<Any>) (bool, string)
func (v <Any>, extra ...<Any>) (bool, string, bool)
```

In these signatures `<Any>` means any type: the arguments defined in validator
tags would be casted to the corresponding type. For example:

Assuming there is a validation tag that looks like: `validate:"foo(bar, baz)"`,
and the corresponding validation function defined as:
```
validator.Register("foo", func(v interface{}, s1, s2 string) bool {
    // some validation logic returning whether value v is valid or not
})
```

The expected values s1 and s2 would be: "bar" and "baz" correspondingly.

Here is another example that should emphasize the idea:

Assuming the same validation function is being used, `validate:"foo(1, 42)"`
would invoke the function def above with: s1 = "1" and s2 = "42" because the
validator library will guess the types from the function definition. Please pay
attention to the argument types you provide in the actual validation function
and ensure an automatic type cast is possible. If not sure, consider using
`interface{}` and do a type casting manually.

In the function argument list, there is 1 mandatory argument `v`: the factual
value found under the tagged struct field. It can be of any type: the value
won't be type-casted while passing around.

The rest of the arguments is completely optional and depends on the validator
logic. A validator function can therefore be variative (see std::enum and
std::range for more details).

A validator function should return either:
* a single bool, indicating whether the value is valid or not
* value above + a string error message; it is safe to always return an error
  message even if the field is correct: it would be ignored
* values above + chain breaker flag; see Chaining section for more details

### Chaining

By default, all validators are chainable: one can declare a validator chain with
a comma-separated list like: `validate:"gt(5), lt(10), ne(7)"`. A chain invokes
validators in a declaration order. It means all validators in the chain are
conjucted with an AND logic.

In contrast to that, a chain breaker flag can be used as a third return value
from a validator function. `validator.Continue` is assumed by default. If
`validator.Break` is returned by defaul, the chain stops here with a validation
result from the current validator function which returned the chain break
directive.

Returning a chain breaker is somewhat rare. This is how `optional` is
implemented: if a zero-value is provided, it prevents the remainig chain from
execution and returns a valid flag.

## Contributing

If you found an issue, please open an issue in this repository.

If you have an idea or a proposal, please submit a pull request.

If in doubts or just want to say hi, feel free to send me an email.

## License

Distributed under MIT License, please see LICENSE file within the code for more details.

## Author

Oleg Sidorov <me at whitebox.io>, 2020-2021
