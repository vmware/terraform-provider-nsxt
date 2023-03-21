package core

import "errors"

var DeserializationError = errors.New("error de-serializing method result")

var UnacceptableContent = errors.New("response content type doesn't match any of the acceptable types")
