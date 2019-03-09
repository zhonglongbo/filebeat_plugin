
package main

import (
	"errors"
	"regexp"
)

var (

	delimiterRE = regexp.MustCompile("(?s)(.*?)%\\{([^}]*?)}")
	suffixRE    = regexp.MustCompile("(.+?)(/(\\d{1,2}))?(->)?$")

	skipFieldPrefix      = "?"
	appendFieldPrefix    = "+"
	indirectFieldPrefix  = "&"
	appendIndirectPrefix = "+&"
	indirectAppendPrefix = "&+"
	greedySuffix         = "->"
	pointerFieldPrefix   = "*"

	defaultJoinString = " "

	errParsingFailure            = errors.New("parsing failure")
	errInvalidTokenizer          = errors.New("invalid dissect tokenizer")
	errEmpty                     = errors.New("empty string provided")
	errMixedPrefixIndirectAppend = errors.New("mixed prefix `&+`")
	errMixedPrefixAppendIndirect = errors.New("mixed prefix `&+`")
	errEmptyKey                  = errors.New("empty key")
)
