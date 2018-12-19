package main

import "fmt"

// Errors is a collection of errors
type Errors []error

func (e Errors) Error() string { return fmt.Sprint([]error(e)) }
