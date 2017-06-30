package errors

/**
*    Copyright (C) 2017 Ethan Frey
**/

import (
	"github.com/pkg/errors"

	abci "github.com/tendermint/abci/types"
)

const defaultErrCode = abci.CodeType_InternalError

type stackTracer interface {
	error
	StackTrace() errors.StackTrace
}

type TMError interface {
	stackTracer
	ErrorCode() abci.CodeType
	Message() string
}

type tmerror struct {
	stackTracer
	code abci.CodeType
	msg  string
}

func (t tmerror) ErrorCode() abci.CodeType {
	return t.code
}

func (t tmerror) Message() string {
	return t.msg
}

// Result converts any error into a abci.Result, preserving as much info
// as possible if it was already a TMError
func Result(err error) abci.Result {
	tm := Wrap(err)
	return abci.Result{
		Code: tm.ErrorCode(),
		Log:  tm.Message(),
	}
}

// Wrap safely takes any error and promotes it to a TMError
func Wrap(err error) TMError {
	// nil or TMError are no-ops
	if err == nil {
		return nil
	}
	// and check for noop
	tm, ok := err.(TMError)
	if ok {
		return tm
	}

	return WithCode(err, defaultErrCode)
}

// WithCode adds a stacktrace if necessary and sets the code and msg,
// overriding the state if err was already TMError
func WithCode(err error, code abci.CodeType) TMError {
	// add a stack only if not present
	st, ok := err.(stackTracer)
	if !ok {
		st = errors.WithStack(err).(stackTracer)
	}
	// and then wrap it with TMError info
	return tmerror{
		stackTracer: st,
		code:        code,
		msg:         err.Error(),
	}
}

// New adds a stacktrace if necessary and sets the code and msg,
// overriding the state if err was already TMError
func New(msg string, code abci.CodeType) TMError {
	// create a new error with stack trace and attach a code
	st := errors.New(msg).(stackTracer)
	return tmerror{
		stackTracer: st,
		code:        code,
		msg:         msg,
	}
}
