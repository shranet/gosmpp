package gosmpp

import (
	"fmt"
	"github.com/linxGnu/gosmpp/data"
	"io"
	"log"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
)

// Transceiver interface.
type Transceiver interface {
	io.Closer
	Submit(pdu.PDU) error
	SystemID() string
}

// Transmitter interface.
type Transmitter interface {
	io.Closer
	Submit(pdu.PDU) error
	SystemID() string
}

// Receiver interface.
type Receiver interface {
	io.Closer
	SystemID() string
}

type DefaultLogger struct {
}

func (l *DefaultLogger) Info(v ...interface{}) {
	l.LevelPrintLn("INFO", v...)
}

func (l *DefaultLogger) Warn(v ...interface{}) {
	l.LevelPrintLn("WARN", v...)
}

func (l *DefaultLogger) Error(v ...interface{}) {
	l.LevelPrintLn("ERROR", v...)
}

func (l *DefaultLogger) Debug(v ...interface{}) {
	l.LevelPrintLn("DEBUG", v...)
}

func (l *DefaultLogger) Printf(format string, v ...interface{}) {
	l.Info(fmt.Sprintf(format, v...))
}

func (l *DefaultLogger) LevelPrintLn(level string, v ...interface{}) {
	params := []interface{}{level}
	params = append(params, v...)

	log.Println(params...)
}

// Settings for TX (transmitter), RX (receiver), TRX (transceiver).
type Settings struct {
	// ReadTimeout is timeout for reading PDU from SMSC.
	// Underlying net.Conn will be stricted with ReadDeadline(now + timeout).
	// This setting is very important to detect connection failure.
	//
	// Must: ReadTimeout > max(0, EnquireLink)
	ReadTimeout time.Duration

	// WriteTimeout is timeout for submitting PDU.
	WriteTimeout time.Duration

	// EnquireLink periodically sends EnquireLink to SMSC.
	// The duration must not be smaller than 1 minute.
	//
	// Zero duration disables auto enquire link.
	EnquireLink time.Duration

	// OnPDU handles received PDU from SMSC.
	//
	// `Responded` flag indicates this pdu is responded automatically,
	// no manual respond needed.
	OnPDU PDUCallback

	// OnReceivingError notifies happened error while reading PDU
	// from SMSC.
	OnReceivingError ErrorCallback

	// OnSubmitError notifies fail-to-submit PDU with along error.
	OnSubmitError PDUErrorCallback

	// OnRebindingError notifies error while rebinding.
	OnRebindingError ErrorCallback

	// OnClosed notifies `closed` event due to State.
	OnClosed ClosedCallback

	OnConnected ConnectedCallback

	Logger   data.SmsboxLogger
	response func(pdu.PDU)
}
