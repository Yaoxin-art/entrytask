package zerorpc

import (
	log "github.com/sirupsen/logrus"
	"net"
	"reflect"
)

type Server struct {
	addr  string
	funcs map[string]reflect.Value
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, funcs: make(map[string]reflect.Value)}
}

func (s *Server) Register(rpcName string, foo interface{}) {
	if _, ok := s.funcs[rpcName]; ok {
		log.Infof("foo exists, rpcName: %s \n", rpcName)
		return
	}
	fooVal := reflect.ValueOf(foo)
	s.funcs[rpcName] = fooVal
}

func (s *Server) Run() {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("start rpc server with addr:%s failure. \n", s.addr)
	}
	log.Infof("server start ok, addr:%s ", s.addr)
	for {
		conn, errAccept := listen.Accept()
		if errAccept != nil {
			log.Fatalf("accept err:%v \n", errAccept)
		}
		serverSession := NewSession(conn)
		buf, errRead := serverSession.Read()
		if errRead != nil {
			log.Fatalf("read err:%v \n", errRead)
		}
		rpcData, errDecode := decode(buf)
		if errDecode != nil {
			log.Errorf("decode err:%v \n", errDecode)
			continue
		}
		foo, ok := s.funcs[rpcData.Name]
		if !ok {
			log.Errorf("rpc foo not support:%s \n", rpcData.Name)
			continue
		}

		inArgs := make([]reflect.Value, 0, len(rpcData.Args))
		for _, arg := range rpcData.Args {
			inArgs = append(inArgs, reflect.ValueOf(arg))
		}

		out := foo.Call(inArgs)
		outArgs := make([]interface{}, 0, len(out))
		for _, o := range out {
			outArgs = append(outArgs, o.Interface())
		}

		resRpcData := RpcData{rpcData.Name, outArgs}
		resBytes, errEncode := encode(resRpcData)
		if errEncode != nil {
			log.Errorf("encode server invoke response err:%v \n", errEncode)
			continue
		}
		errWrite := serverSession.Write(resBytes)
		if errWrite != nil {
			log.Errorf("write err:%v \n", errWrite)
			continue
		}
	}
}
