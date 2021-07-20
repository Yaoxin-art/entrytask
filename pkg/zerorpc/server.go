package zerorpc

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"reflect"
	"strings"
	"time"
)

type Server struct {
	addr     string
	facades  map[string]reflect.Value
	listener net.Listener
}

// NewServer 创建服务
func NewServer(addr string) *Server {
	return &Server{addr: addr, facades: make(map[string]reflect.Value)}
}

// Register 注册暴露的接口
func (s *Server) Register(methodId string, foo interface{}) {
	if _, ok := s.facades[methodId]; ok {
		// if exist, ignore
		return
	}
	fooVal := reflect.ValueOf(foo)
	s.facades[methodId] = fooVal
	logrus.Infof("service methodId:%s for:%T", methodId, foo)
}

func (s *Server) Run() {
	var errListen error
	s.listener, errListen = net.Listen("tcp", s.addr)
	if errListen != nil {
		logrus.Fatalf("start zerorpc server with addr:%s failure, err:%v \n", s.addr, errListen)
	}
	logrus.Infof("server start ok, bind at:%s", s.addr)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
			logrus.Fatalf("accept err:%v", err)
		}
		session := NewSession(conn)
		go s.handleNewConn(session)
	}
}

func (s *Server) Close() {
	err := s.listener.Close()
	if err != nil {
		logrus.Errorf("zerorpc server close err:%v", err)
	}
}

func (s *Server) handleNewConn(session *Session) {
	defer session.Close()
	logrus.Infof("Handle conn from addr:%v", session.remoteAddr)
	for {
		// read
		data, err := session.Read()
		// init response
		var result Result                               // init rcp result
		result.RequestTimestamp = time.Now().UnixNano() // get request timestamp
		if err != nil {
			if err.Error() == "EOF" { // ignore
				continue
			}
			if strings.Contains(err.Error(), "connect reset by peer") { // disconnect
				break
			}
			result.Err = fmt.Sprintf("zerorpc request from addr:%s read err:%v", session.remoteAddr, err)

			result.ResponseTimestamp = time.Now().UnixNano() // invoke end timestamp
			resBytes, err := encode(result)
			if err != nil {
				// todo: ignore it ?
				logrus.Errorf("encode zerorpc result err:%v, from addr:%v", err, session.remoteAddr)
			}
			errWrite := session.Write(resBytes)
			if errWrite != nil {
				if strings.Contains(errWrite.Error(), "broken pipe") {
					continue
				}
				logrus.Errorf("zerorpc write err:%v", errWrite)
			}
			continue
		}
		//logrus.Infof("read from invocation ok, data length:%d", len(data))
		// decode
		invocation, err := decodeInvocation(data)
		//logrus.Infof("read from invocation ok, invocation:%v", invocation)
		if err != nil {
			logrus.Errorf("decode err:%v", err)
			result.Err = fmt.Sprintf("zerorpc request from addr:%s decode err:%v", session.remoteAddr, err)

			result.ResponseTimestamp = time.Now().UnixNano() // invoke end timestamp
			resBytes, err := encode(result)
			if err != nil {
				logrus.Errorf("encode zerorpc result err:%v, from addr:%v", err, session.remoteAddr)
			}
			errWrite := session.Write(resBytes)
			if errWrite != nil {
				if strings.Contains(errWrite.Error(), "broken pipe") {
					continue
				}
				logrus.Errorf("zerorpc write err:%v", errWrite)
			}
			continue
		}
		logrus.Debugf("get request from addr:%s, methodId:%s at:%s", session.remoteAddr, invocation.MethodId, time.Now().Format(time.RFC3339))
		// invoke
		foo, ok := s.facades[invocation.MethodId]
		if !ok {
			// request methodId not exist
			result.Err = fmt.Sprintf("zerorpc request methodId:%s not exists", invocation.MethodId)

			result.ResponseTimestamp = time.Now().UnixNano() // invoke end timestamp
			resBytes, err := encode(result)
			if err != nil {
				logrus.Errorf("encode zerorpc result err:%v, from addr:%v", err, session.remoteAddr)
			}
			errWrite := session.Write(resBytes)
			if errWrite != nil {
				if strings.Contains(errWrite.Error(), "broken pipe") {
					continue
				}
				logrus.Errorf("zerorpc write err:%v", errWrite)
			}
			continue
		}
		logrus.Infof("get request from invocation for method:%s : %T", invocation.MethodId, foo)
		inArgs := make([]reflect.Value, 0, len(invocation.Args))
		for _, arg := range invocation.Args {
			inArgs = append(inArgs, reflect.ValueOf(arg))
		}
		out := foo.Call(inArgs) // invoke
		logrus.Infof("call method success, method:%T, out:%v", foo, out)
		result.ResponseTimestamp = time.Now().UnixNano() // invoke end timestamp
		outArgs := make([]interface{}, 0, len(out))
		for _, o := range out {
			outArgs = append(outArgs, o.Interface())
		}
		result.Args = outArgs
		result.StatusCode = 1 // success code
		// encode
		res, err := encode(result)
		if err != nil {
			logrus.Errorf("encode zerorpc result err:%v, from addr:%v", err, session.remoteAddr)
		}
		errWrite := session.Write(res)
		// write
		if errWrite != nil {
			logrus.Errorf("write zerorpc result err:%v, from addr:%s", err, session.remoteAddr)
		}
		logrus.Debugf("zerorpc invoke success, from addr:%s, methodId:%s, result:%v", session.remoteAddr, invocation.MethodId, result)
	}
}
