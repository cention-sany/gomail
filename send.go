package gomail

import (
	"crypto/tls"
	"io"
	"net"
	"net/smtp"
)

func (m *Mailer) getSendMailFunc(ssl bool) SendMailFunc {
	return func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		var c smtpClient
		var err error
		if ssl {
			c, err = sslDial(addr, m.host, m.config)
			if err != nil {
				return err
			} else {
				if "" != m.localhost {
					if err = c.Hello(m.localhost); err != nil {
						c.Close()
						return err
					}
				}
			}
		} else {
			c, err = starttlsDial(addr, m.localhost, m.config)
			if err != nil {
				return err
			}
		}
		defer c.Close()
		m.Closer = c
		if a != nil {
			if ok, _ := c.Extension("AUTH"); ok {
				if err = c.Auth(a); err != nil {
					return err
				}
			}
		}

		if err = c.Mail(from); err != nil {
			return err
		}

		for _, addr := range to {
			if err = c.Rcpt(addr); err != nil {
				return err
			}
		}

		w, err := c.Data()
		if err != nil {
			return err
		}
		_, err = w.Write(msg)
		if err != nil {
			return err
		}
		err = w.Close()
		if err != nil {
			return err
		}

		return c.Quit()
	}
}

// no need to close connection if return non-nil error
func sslDial(addr, host string, config *tls.Config) (smtpClient, error) {
	conn, err := initTLS("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return newClient(conn, host)
}

// no need to close connection if return non-nil error
func starttlsDial(addr, localhost string, config *tls.Config) (smtpClient, error) {
	c, err := initSMTP(addr)
	if err != nil {
		return c, err
	}
	if "" != localhost {
		if err = c.Hello(localhost); err != nil {
			c.Close()
			return c, err
		}
	}
	if ok, _ := c.Extension("STARTTLS"); ok {
		if err = c.StartTLS(config); nil != err {
			c.Close()
			return c, err
		}
	}
	return c, nil
}

var initSMTP = func(addr string) (smtpClient, error) {
	return smtp.Dial(addr)
}

var initTLS = func(network, addr string, config *tls.Config) (*tls.Conn, error) {
	return tls.Dial(network, addr, config)
}

var newClient = func(conn net.Conn, host string) (smtpClient, error) {
	return smtp.NewClient(conn, host)
}

type smtpClient interface {
	Extension(string) (bool, string)
	StartTLS(*tls.Config) error
	Auth(smtp.Auth) error
	Hello(string) error
	Mail(string) error
	Rcpt(string) error
	Data() (io.WriteCloser, error)
	Quit() error
	Close() error
}
