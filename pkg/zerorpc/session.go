package zerorpc

import (
	"encoding/binary"
	"io"
	"net"
)

type Session struct {
	conn net.Conn
}

func NewSession(conn net.Conn) *Session {
	return &Session{conn}
}

func (s Session) Write(data []byte) error {
	buf := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	copy(buf[:4], data)
	_, err := s.conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (s Session) Read() ([]byte, error) {
	header := make([]byte, 4)
	_, errHeader := io.ReadFull(s.conn, header)
	if errHeader != nil {
		return nil, errHeader
	}
	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	_, errData := io.ReadFull(s.conn, data)
	if errData != nil {
		return nil, errData
	}
	return data, nil
}
