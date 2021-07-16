package zerorpc

import (
	log "github.com/sirupsen/logrus"
	"net"
	"reflect"
)

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) callRpc(rpcName string, fooPtr interface{}) {
	fn := reflect.ValueOf(fooPtr).Elem()

	foo := func(args []reflect.Value) []reflect.Value {
		inArgs := make([]interface{}, 0, len(args))
		for _, arg := range args {
			inArgs = append(inArgs, arg.Interface())
		}
		clientSession := NewSession(c.conn)

		reqRpc := RpcData{Name: rpcName, Args: inArgs}
		writeBuf, errEncode := encode(reqRpc)
		if errEncode != nil {
			log.Errorf("encode rpc request err:%v \n", errEncode)
			panic(errEncode)
		}

		errWrite := clientSession.Write(writeBuf)
		if errWrite != nil {
			log.Errorf("client write err:%v \n", errWrite)
			panic(errWrite)
		}
		resBuf, errRead := clientSession.Read()
		if errRead != nil {
			log.Errorf("client read err:%v \n", errRead)
			panic(errRead)
		}
		resRpc, errDecode := decode(resBuf)
		if errDecode != nil {
			log.Errorf("client decode response err:%v \n", errDecode)
			panic(errDecode)
		}
		resArgs := make([]reflect.Value, 0, len(resRpc.Args))
		for i, arg := range resRpc.Args {
			if arg == nil {
				// response item nil, fill with real type
				resArgs = append(resArgs, reflect.Zero(fn.Type().Out(i)))
				continue
			}
		}
		return resArgs
	}

	res := reflect.MakeFunc(fn.Type(), foo)
	fn.Set(res)
}
