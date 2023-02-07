// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/binancedata/v1/binancedata.proto

package v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on PullBinanceDataRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *PullBinanceDataRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PullBinanceDataRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PullBinanceDataRequestMultiError, or nil if none found.
func (m *PullBinanceDataRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *PullBinanceDataRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Start

	if len(errors) > 0 {
		return PullBinanceDataRequestMultiError(errors)
	}

	return nil
}

// PullBinanceDataRequestMultiError is an error wrapping multiple validation
// errors returned by PullBinanceDataRequest.ValidateAll() if the designated
// constraints aren't met.
type PullBinanceDataRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PullBinanceDataRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PullBinanceDataRequestMultiError) AllErrors() []error { return m }

// PullBinanceDataRequestValidationError is the validation error returned by
// PullBinanceDataRequest.Validate if the designated constraints aren't met.
type PullBinanceDataRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PullBinanceDataRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PullBinanceDataRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PullBinanceDataRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PullBinanceDataRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PullBinanceDataRequestValidationError) ErrorName() string {
	return "PullBinanceDataRequestValidationError"
}

// Error satisfies the builtin error interface
func (e PullBinanceDataRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPullBinanceDataRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PullBinanceDataRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PullBinanceDataRequestValidationError{}

// Validate checks the field values on PullBinanceDataReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *PullBinanceDataReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PullBinanceDataReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PullBinanceDataReplyMultiError, or nil if none found.
func (m *PullBinanceDataReply) ValidateAll() error {
	return m.validate(true)
}

func (m *PullBinanceDataReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return PullBinanceDataReplyMultiError(errors)
	}

	return nil
}

// PullBinanceDataReplyMultiError is an error wrapping multiple validation
// errors returned by PullBinanceDataReply.ValidateAll() if the designated
// constraints aren't met.
type PullBinanceDataReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PullBinanceDataReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PullBinanceDataReplyMultiError) AllErrors() []error { return m }

// PullBinanceDataReplyValidationError is the validation error returned by
// PullBinanceDataReply.Validate if the designated constraints aren't met.
type PullBinanceDataReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PullBinanceDataReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PullBinanceDataReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PullBinanceDataReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PullBinanceDataReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PullBinanceDataReplyValidationError) ErrorName() string {
	return "PullBinanceDataReplyValidationError"
}

// Error satisfies the builtin error interface
func (e PullBinanceDataReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPullBinanceDataReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PullBinanceDataReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PullBinanceDataReplyValidationError{}

// Validate checks the field values on IntervalMAvgEndPriceDataRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *IntervalMAvgEndPriceDataRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on IntervalMAvgEndPriceDataRequest with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// IntervalMAvgEndPriceDataRequestMultiError, or nil if none found.
func (m *IntervalMAvgEndPriceDataRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *IntervalMAvgEndPriceDataRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Start

	// no validation rules for End

	// no validation rules for M

	// no validation rules for N

	if len(errors) > 0 {
		return IntervalMAvgEndPriceDataRequestMultiError(errors)
	}

	return nil
}

// IntervalMAvgEndPriceDataRequestMultiError is an error wrapping multiple
// validation errors returned by IntervalMAvgEndPriceDataRequest.ValidateAll()
// if the designated constraints aren't met.
type IntervalMAvgEndPriceDataRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m IntervalMAvgEndPriceDataRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m IntervalMAvgEndPriceDataRequestMultiError) AllErrors() []error { return m }

// IntervalMAvgEndPriceDataRequestValidationError is the validation error
// returned by IntervalMAvgEndPriceDataRequest.Validate if the designated
// constraints aren't met.
type IntervalMAvgEndPriceDataRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e IntervalMAvgEndPriceDataRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e IntervalMAvgEndPriceDataRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e IntervalMAvgEndPriceDataRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e IntervalMAvgEndPriceDataRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e IntervalMAvgEndPriceDataRequestValidationError) ErrorName() string {
	return "IntervalMAvgEndPriceDataRequestValidationError"
}

// Error satisfies the builtin error interface
func (e IntervalMAvgEndPriceDataRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sIntervalMAvgEndPriceDataRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = IntervalMAvgEndPriceDataRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = IntervalMAvgEndPriceDataRequestValidationError{}

// Validate checks the field values on IntervalMAvgEndPriceDataReply with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *IntervalMAvgEndPriceDataReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on IntervalMAvgEndPriceDataReply with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// IntervalMAvgEndPriceDataReplyMultiError, or nil if none found.
func (m *IntervalMAvgEndPriceDataReply) ValidateAll() error {
	return m.validate(true)
}

func (m *IntervalMAvgEndPriceDataReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetData() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, IntervalMAvgEndPriceDataReplyValidationError{
						field:  fmt.Sprintf("Data[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, IntervalMAvgEndPriceDataReplyValidationError{
						field:  fmt.Sprintf("Data[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return IntervalMAvgEndPriceDataReplyValidationError{
					field:  fmt.Sprintf("Data[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	for idx, item := range m.GetOperationData() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, IntervalMAvgEndPriceDataReplyValidationError{
						field:  fmt.Sprintf("OperationData[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, IntervalMAvgEndPriceDataReplyValidationError{
						field:  fmt.Sprintf("OperationData[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return IntervalMAvgEndPriceDataReplyValidationError{
					field:  fmt.Sprintf("OperationData[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for OperationOrderTotal

	// no validation rules for OperationWinRate

	// no validation rules for OperationWinAmount

	if len(errors) > 0 {
		return IntervalMAvgEndPriceDataReplyMultiError(errors)
	}

	return nil
}

// IntervalMAvgEndPriceDataReplyMultiError is an error wrapping multiple
// validation errors returned by IntervalMAvgEndPriceDataReply.ValidateAll()
// if the designated constraints aren't met.
type IntervalMAvgEndPriceDataReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m IntervalMAvgEndPriceDataReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m IntervalMAvgEndPriceDataReplyMultiError) AllErrors() []error { return m }

// IntervalMAvgEndPriceDataReplyValidationError is the validation error
// returned by IntervalMAvgEndPriceDataReply.Validate if the designated
// constraints aren't met.
type IntervalMAvgEndPriceDataReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e IntervalMAvgEndPriceDataReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e IntervalMAvgEndPriceDataReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e IntervalMAvgEndPriceDataReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e IntervalMAvgEndPriceDataReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e IntervalMAvgEndPriceDataReplyValidationError) ErrorName() string {
	return "IntervalMAvgEndPriceDataReplyValidationError"
}

// Error satisfies the builtin error interface
func (e IntervalMAvgEndPriceDataReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sIntervalMAvgEndPriceDataReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = IntervalMAvgEndPriceDataReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = IntervalMAvgEndPriceDataReplyValidationError{}

// Validate checks the field values on IntervalMAvgEndPriceDataReply_List with
// the rules defined in the proto definition for this message. If any rules
// are violated, the first error encountered is returned, or nil if there are
// no violations.
func (m *IntervalMAvgEndPriceDataReply_List) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on IntervalMAvgEndPriceDataReply_List
// with the rules defined in the proto definition for this message. If any
// rules are violated, the result is a list of violation errors wrapped in
// IntervalMAvgEndPriceDataReply_ListMultiError, or nil if none found.
func (m *IntervalMAvgEndPriceDataReply_List) ValidateAll() error {
	return m.validate(true)
}

func (m *IntervalMAvgEndPriceDataReply_List) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for StartPrice

	// no validation rules for EndPrice

	// no validation rules for TopPrice

	// no validation rules for LowPrice

	// no validation rules for WithBeforeAvgEndPrice

	// no validation rules for Time

	if len(errors) > 0 {
		return IntervalMAvgEndPriceDataReply_ListMultiError(errors)
	}

	return nil
}

// IntervalMAvgEndPriceDataReply_ListMultiError is an error wrapping multiple
// validation errors returned by
// IntervalMAvgEndPriceDataReply_List.ValidateAll() if the designated
// constraints aren't met.
type IntervalMAvgEndPriceDataReply_ListMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m IntervalMAvgEndPriceDataReply_ListMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m IntervalMAvgEndPriceDataReply_ListMultiError) AllErrors() []error { return m }

// IntervalMAvgEndPriceDataReply_ListValidationError is the validation error
// returned by IntervalMAvgEndPriceDataReply_List.Validate if the designated
// constraints aren't met.
type IntervalMAvgEndPriceDataReply_ListValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e IntervalMAvgEndPriceDataReply_ListValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e IntervalMAvgEndPriceDataReply_ListValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e IntervalMAvgEndPriceDataReply_ListValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e IntervalMAvgEndPriceDataReply_ListValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e IntervalMAvgEndPriceDataReply_ListValidationError) ErrorName() string {
	return "IntervalMAvgEndPriceDataReply_ListValidationError"
}

// Error satisfies the builtin error interface
func (e IntervalMAvgEndPriceDataReply_ListValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sIntervalMAvgEndPriceDataReply_List.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = IntervalMAvgEndPriceDataReply_ListValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = IntervalMAvgEndPriceDataReply_ListValidationError{}

// Validate checks the field values on IntervalMAvgEndPriceDataReply_List2 with
// the rules defined in the proto definition for this message. If any rules
// are violated, the first error encountered is returned, or nil if there are
// no violations.
func (m *IntervalMAvgEndPriceDataReply_List2) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on IntervalMAvgEndPriceDataReply_List2
// with the rules defined in the proto definition for this message. If any
// rules are violated, the result is a list of violation errors wrapped in
// IntervalMAvgEndPriceDataReply_List2MultiError, or nil if none found.
func (m *IntervalMAvgEndPriceDataReply_List2) ValidateAll() error {
	return m.validate(true)
}

func (m *IntervalMAvgEndPriceDataReply_List2) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for StartPrice

	// no validation rules for EndPrice

	// no validation rules for StartTime

	// no validation rules for Time

	// no validation rules for EndTime

	// no validation rules for Type

	// no validation rules for Status

	// no validation rules for Rate

	// no validation rules for CloseEndPrice

	if len(errors) > 0 {
		return IntervalMAvgEndPriceDataReply_List2MultiError(errors)
	}

	return nil
}

// IntervalMAvgEndPriceDataReply_List2MultiError is an error wrapping multiple
// validation errors returned by
// IntervalMAvgEndPriceDataReply_List2.ValidateAll() if the designated
// constraints aren't met.
type IntervalMAvgEndPriceDataReply_List2MultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m IntervalMAvgEndPriceDataReply_List2MultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m IntervalMAvgEndPriceDataReply_List2MultiError) AllErrors() []error { return m }

// IntervalMAvgEndPriceDataReply_List2ValidationError is the validation error
// returned by IntervalMAvgEndPriceDataReply_List2.Validate if the designated
// constraints aren't met.
type IntervalMAvgEndPriceDataReply_List2ValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e IntervalMAvgEndPriceDataReply_List2ValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e IntervalMAvgEndPriceDataReply_List2ValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e IntervalMAvgEndPriceDataReply_List2ValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e IntervalMAvgEndPriceDataReply_List2ValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e IntervalMAvgEndPriceDataReply_List2ValidationError) ErrorName() string {
	return "IntervalMAvgEndPriceDataReply_List2ValidationError"
}

// Error satisfies the builtin error interface
func (e IntervalMAvgEndPriceDataReply_List2ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sIntervalMAvgEndPriceDataReply_List2.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = IntervalMAvgEndPriceDataReply_List2ValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = IntervalMAvgEndPriceDataReply_List2ValidationError{}
