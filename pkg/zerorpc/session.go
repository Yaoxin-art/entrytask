package zerorpc

import (
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"strings"
)

type Session struct {
	conn       net.Conn
	remoteAddr interface{}
	destroyed  bool
}

const sizeOfInt = 4

// NewSession 确认一个会话
func NewSession(conn net.Conn) *Session {
	return &Session{conn: conn, remoteAddr: conn.RemoteAddr()}
}

// Write write to remote
func (s *Session) Write(data []byte) error {
	// 前一个字节放置当前数据的大小
	buf := make([]byte, sizeOfInt+len(data))
	// 计算data的大小，并赋值到最前面
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	copy(buf[4:], data)
	_, err := s.conn.Write(buf)
	if err != nil {
		if strings.Contains(err.Error(), "broken pipe") {
			return nil
		}
		logrus.Errorf("write data to:%s err:%v", s.remoteAddr, err)
		s.destroyed = true
		return err
	}
	logrus.Debugf("write data to:%s success, length:%d", s.remoteAddr, len(buf))
	return nil
}

// Read read from remote
func (s *Session) Read() ([]byte, error) {
	topSizeOfInt := make([]byte, sizeOfInt)
	_, errReadHead := io.ReadFull(s.conn, topSizeOfInt)
	if errReadHead != nil {
		if errReadHead.Error() == "EOF" || strings.Contains(errReadHead.Error(), "reset by peer") {	// ignore
			s.destroyed = true
			return nil, errReadHead
		}
		logrus.Errorf("read data length from addr:%s err:%v", s.remoteAddr, errReadHead)
		return nil, errReadHead
	}

	dataLen := binary.BigEndian.Uint32(topSizeOfInt)
	data := make([]byte, dataLen)
	_, errReadData := io.ReadFull(s.conn, data)
	if errReadData != nil {
		logrus.Errorf("read data from addr:%s err:%v", s.remoteAddr, errReadData)
		return nil, errReadData
	}
	logrus.Debugf("read data from addr:%s success, length:%d", s.remoteAddr, dataLen)
	return data, nil
}

func (s *Session) Close() error {
	s.destroyed = true
	err := s.conn.Close()
	if err != nil {
		logrus.Errorf("session close err:%v, remote addr:%v", err, s.remoteAddr)
		return err
	} else {
		logrus.Debugf("session closed, remote addr:%s", s.remoteAddr)
		return nil
	}
}
