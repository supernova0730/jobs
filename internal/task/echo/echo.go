package echo

import (
	"errors"
	"fmt"
)

type Success struct{}

func (s *Success) Do() error {
	return nil
}
func (s *Success) Code() string { return "ECHO_SUCCESS" }

type Error struct{}

func (e *Error) Do() error {
	return errors.New("error message")
}
func (e *Error) Code() string { return "ECHO_ERROR" }

type Panic struct{}

func (p *Panic) Do() error {
	var a []int
	fmt.Println(a[0])
	return nil
}
func (p *Panic) Code() string { return "ECHO_PANIC" }
