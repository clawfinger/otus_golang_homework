package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var errNoConnection = errors.New("connection is not established")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	conn       net.Conn
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &Client{
		address: address, timeout: timeout, in: in, out: out, conn: nil, ctx: ctx, cancelFunc: cancelFunc,
	}
}

func (c *Client) Connect() error {
	var dialer net.Dialer
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	var err error
	c.conn, err = dialer.DialContext(ctx, "tcp", c.address)
	return err
}

func (c *Client) Send() error {
	if c.conn == nil {
		return errNoConnection
	}
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return fmt.Errorf("unable to copyBytes: %w", err)
	}
	return nil
}

func (c *Client) Receive() error {
	if c.conn == nil {
		return errNoConnection
	}
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return fmt.Errorf("unable to copyBytes: %w", err)
	}
	return nil
}

func (c *Client) Close() error {
	c.in.Close()
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
