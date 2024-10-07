package fault

import (
	"fmt"
	"strings"
	"sync"
)

type ErrComponent string

type ResponseErrType string

type ErrorCode string

type BasicFaultsCache struct {
	mx          *sync.Mutex
	basicFaults map[ErrorCode]BasicFault
}

func NewBasicFaultCache(partialFaults map[ErrorCode]BasicFault) BasicFaultsCache {
	return BasicFaultsCache{
		basicFaults: partialFaults,
	}
}

func (fc *BasicFaultsCache) AppendBasicFaults(bfaults map[ErrorCode]BasicFault) {
	defer fc.mx.Unlock()
	fc.mx.Lock()
	for k, v := range bfaults {
		fc.basicFaults[k] = v
	}
}

func (p *BasicFaultsCache) GetBasicFault(code ErrorCode) BasicFault {
	pf := p.basicFaults[code]
	if pf != nil {
		return pf
	}
	defer p.mx.Unlock()
	p.mx.Lock()
	bf := NewBasicFault(code)
	p.basicFaults[code] = bf
	return bf
}

func (e ErrorCode) String() string {
	return string(e)
}

type BasicFault interface {
	error
	ResponseErrType() ResponseErrType
	SetResponseType(r ResponseErrType) BasicFault
	Component() ErrComponent
	SetComponent(c ErrComponent) BasicFault
	Code() ErrorCode
	ToFault(data map[string]any, cause error) Fault
}

type BasicAppError struct {
	code         ErrorCode
	component    ErrComponent
	responseType ResponseErrType
}

func NewBasicFault(
	code ErrorCode,
) BasicFault {
	return &BasicAppError{
		code: code,
	}
}

func (p *BasicAppError) SetComponent(component ErrComponent) BasicFault {
	p.component = component
	return p
}

func (p *BasicAppError) Component() ErrComponent {
	return p.component
}

func (p *BasicAppError) SetResponseType(reponse ResponseErrType) BasicFault {
	p.responseType = reponse
	return p
}

func (p *BasicAppError) ResponseErrType() ResponseErrType {
	return p.responseType
}

func (e *BasicAppError) Code() ErrorCode {
	return e.code
}

func (e *BasicAppError) Error() string {
	return e.code.String()
}

func (e *BasicAppError) ToFault(data map[string]any, cause error) Fault {
	var lCauses []error // ensure not nil
	if cause != nil {
		lCauses = append(lCauses, cause)
	}

	lData := map[string]any{} // ensure not nil
	if data != nil {
		lData = data
	}
	return &AppError{
		BasicAppError: *e,
		data:          lData,
		causes:        lCauses,
	}
}

type Fault interface {
	BasicFault
	Cause() error
	Causes() []error
	Data() map[string]interface{}
	Source() string
	SetSource(s string) Fault
	Retryable() bool
	SetRetryable() Fault
	AppendCause() Fault
	// ToMessageAwareFault(*i18n.Bundle) MessageAwareFault
}

type AppError struct {
	BasicAppError
	source      string
	data        map[string]any // cannot be nil
	causes      []error        // cannot be nil
	retryable   bool
	appendCause bool
}

func (e *AppError) Source() string {
	return e.source
}

func (e *AppError) SetSource(s string) Fault {
	e.source = s
	return e
}

func (e *AppError) Cause() error {
	if len(e.causes) > 0 {
		return e.causes[0]
	}
	return nil
}

func (e *AppError) Causes() []error {
	return e.causes
}

func (e *AppError) Data() map[string]interface{} {
	return e.data
}

func (e *AppError) Retryable() bool {
	return e.retryable
}

func (e *AppError) SetRetryable() Fault {
	e.retryable = true
	return e
}

func (e *AppError) AppendCause() Fault {
	e.appendCause = true
	return e
}

func (e *AppError) Error() string {
	var s string
	s = e.code.String()
	if e.appendCause {
		s += "; " + getCauses(e.causes)
	}
	return s
}

type MessageAwareFault interface {
	Fault
	Message(lang string) string
	String() string
}

type AppErrorMessage struct {
	AppError
}

func (e *AppErrorMessage) ErrorString() string {
	localizedMsg := e.Message("en")
	var s string
	if localizedMsg == "" {
		s = e.code.String()
	} else {
		s = fmt.Sprintf("%s:%s", e.code, e.Message("en"))
	}
	if e.appendCause {
		s += "; " + getCauses(e.causes)
	}
	return s
}

func (e *AppErrorMessage) Message(lang string) string {
	return e.code.String()
}

func (e *AppErrorMessage) String() string {
	return e.Error()
}

func getCauses(errors []error) string {
	var s strings.Builder
	for _, err := range errors {
		s.WriteString(err.Error())
		s.WriteString("; ")
	}
	return s.String()
}
