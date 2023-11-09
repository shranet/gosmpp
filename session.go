package gosmpp

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Session represents session for TX, RX, TRX.
type Session struct {
	c Connector

	originalOnClosed func(State)
	settings         Settings

	rebindingInterval time.Duration

	trx atomic.Value // transceivable

	state     int32
	rebinding int32
}

// NewSession creates new session for TX, RX, TRX.
//
// Session will `non-stop`, automatically rebind (create new and authenticate connection with SMSC) when
// unexpected error happened.
//
// `rebindingInterval` indicates duration that Session has to wait before rebinding again.
//
// Setting `rebindingInterval <= 0` will disable `auto-rebind` functionality.
func NewSession(c Connector, settings Settings, rebindingInterval time.Duration) (session *Session, err error) {
	if settings.ReadTimeout <= 0 || settings.ReadTimeout <= settings.EnquireLink {
		return nil, fmt.Errorf("invalid settings: ReadTimeout must greater than max(0, EnquireLink)")
	}

	//if settings.Logger == nil {
	//	settings.Logger = &DefaultLogger{}
	//}

	session = &Session{
		c:                 c,
		rebindingInterval: rebindingInterval,
		originalOnClosed:  settings.OnClosed,
	}

	if rebindingInterval > 0 {
		newSettings := settings
		newSettings.OnClosed = func(state State) {
			switch state {
			case ExplicitClosing:
				return

			default:
				if session.originalOnClosed != nil {
					session.originalOnClosed(state)
				}
				session.rebind()
			}
		}
		session.settings = newSettings
	} else {
		session.settings = settings
	}

	go session.rebind()

	return
}

func (s *Session) bound() *transceivable {
	r, _ := s.trx.Load().(*transceivable)
	return r
}

// Transmitter returns bound Transmitter.
func (s *Session) Transmitter() Transmitter {
	return s.bound()
}

// Receiver returns bound Receiver.
func (s *Session) Receiver() Receiver {
	return s.bound()
}

// Transceiver returns bound Transceiver.
func (s *Session) Transceiver() Transceiver {
	return s.bound()
}

// Close session.
func (s *Session) Close() (err error) {
	if atomic.CompareAndSwapInt32(&s.state, Alive, Closed) {
		err = s.close()
	}
	return
}

func (s *Session) IsAlive() bool {
	if b := s.bound(); b != nil && b.out != nil && b.in != nil {
		outAlive := atomic.LoadInt32(&b.out.aliveState) == Alive
		inAlive := atomic.LoadInt32(&b.in.aliveState) == Alive
		if s.settings.Logger != nil {
			s.settings.Logger.Info("[SESSION]", "in:", inAlive, "out:", outAlive)
		}
		return outAlive && inAlive
	} else {
		if s.settings.Logger != nil {
			if b == nil {
				s.settings.Logger.Info("[SESSION]", "bound is nil")
			} else {
				if b.in == nil {
					s.settings.Logger.Info("[SESSION]", "in is nil")
				}

				if b.out == nil {
					s.settings.Logger.Info("[SESSION]", "out is nil")
				}
			}
		}
	}

	return false
}

func (s *Session) close() (err error) {
	if b := s.bound(); b != nil {
		err = b.Close()
	}
	return
}

func (s *Session) rebind() {
	if s.settings.Logger != nil {
		s.settings.Logger.Info("[SESSION]", "rebind")
	}

	defer func() {
		if s.settings.Logger != nil {
			s.settings.Logger.Info("[SESSION]", "Exit from rebind")
		}
	}()

	if atomic.CompareAndSwapInt32(&s.rebinding, 0, 1) {
		_ = s.close()

		for atomic.LoadInt32(&s.state) == Alive {
			conn, err := s.c.Connect()
			if err != nil {
				if s.settings.OnRebindingError != nil {
					s.settings.OnRebindingError(err)
				}
				time.Sleep(s.rebindingInterval)
			} else {
				if s.settings.Logger != nil {
					s.settings.Logger.Info("[SESSION]", "Store new transceivable")
				}
				// bind to session
				s.trx.Store(newTransceivable(conn, s.settings))

				// reset rebinding state
				atomic.StoreInt32(&s.rebinding, 0)

				if s.settings.OnConnected != nil {
					s.settings.OnConnected()
				}

				return
			}
		}
	}
}
