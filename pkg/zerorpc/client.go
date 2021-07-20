package zerorpc

import (
	"github.com/sirupsen/logrus"
	"net"
	"reflect"
	"time"
)

type Client struct {
	providers    map[string]string
	connCache    map[string]*Session
	providerConn map[string]*Session
}

func NewClient() *Client {
	return &Client{providers: make(map[string]string), connCache: make(map[string]*Session), providerConn: make(map[string]*Session)}
}

// Config 配置远程方法及提供者信息
// methodId:		远程方法Id
// providerAddr:	远程服务所在地址，eg："127.0.0.1:9999"
// fooPtr:			远程方法声明
func (c *Client) Config(methodId, providerAddr string, fooPtr interface{}) {
	if _, ok := c.providers[methodId]; !ok {
		// not config
		c.providers[methodId] = providerAddr
	}
	if _, ok := c.connCache[providerAddr]; !ok {
		conn, err := net.Dial("tcp", providerAddr)
		if err != nil {
			logrus.Errorf("connect to provider err:%v, provider addr:%s, methodId:%s", err, providerAddr, methodId)
			return
		}
		c.connCache[providerAddr] = NewSession(conn)
	}
	c.providerConn[methodId] = c.connCache[providerAddr]
	c.callRPC(methodId, fooPtr)
	logrus.Infof("config methodId:%s -> %T, provider:%s success", methodId, fooPtr, providerAddr)
}

// CallRPC 定义rpc代理
func (c *Client) callRPC(methodId string, fooPtr interface{}) {
	fn := reflect.ValueOf(fooPtr).Elem()
	foo := func(args []reflect.Value) []reflect.Value {
		inArgs := make([]interface{}, 0, len(args))
		for _, arg := range args {
			inArgs = append(inArgs, arg.Interface())
		}
		// get remote session
		session := c.providerConn[methodId]
		if session == nil {
			logrus.Errorf("remote method:%s not config", methodId)
			return zeroValueFnOut(fooPtr)
		}
		// encode invocation
		invocation := Invocation{MethodId: methodId, Args: inArgs}
		writeBuf, errEncode := encode(invocation)
		if errEncode != nil {
			logrus.Errorf("encode invocation:%v err:%v", invocation, errEncode)
			return zeroValueFnOut(fooPtr)
		}
		requestTimestamp := time.Now()
		// send request
		errWrite := session.Write(writeBuf)
		if errWrite != nil {
			logrus.Errorf("write invocation:%v err:%v", invocation, errWrite)
			return zeroValueFnOut(fooPtr)
		}
		// wait for response
		// read
		resBuf, errRead := session.Read()
		if errRead != nil {
			logrus.Errorf("read zerorpc response err:%v, invocation:%v", errRead, invocation)
			return zeroValueFnOut(fooPtr)
		}
		responseTimestamp := time.Now()
		logrus.Debugf("zerorpc invocation:%v, invoke at:%s, response at:%s, spent:%d ms", invocation, requestTimestamp.Format(time.RFC3339), responseTimestamp.Format(time.RFC3339), (responseTimestamp.UnixNano()-requestTimestamp.UnixNano())/1000000)
		// decode zerorpc response
		resRpc, errDecode := decodeResult(resBuf)
		if errDecode != nil {
			logrus.Errorf("decode zerorpc response err:%v, invocation:%v", errDecode, invocation)
			return zeroValueFnOut(fooPtr)
		}
		// deal with zerorpc result
		if resRpc.Err != "" {
			logrus.Errorf("zerorpc invocation:%v err from remote:%v", invocation, resRpc.Err)
			return zeroValueFnOut(fooPtr)
		}
		resArgs := make([]reflect.Value, 0, len(resRpc.Args))
		for i, arg := range resRpc.Args {
			if arg == nil {
				// if null, fill with blank
				resArgs = append(resArgs, reflect.Zero(fn.Type().Out(i)))
			} else {
				resArgs = append(resArgs, reflect.ValueOf(arg))
			}
		}
		return resArgs
	}
	fn.Set(reflect.MakeFunc(fn.Type(), foo))
}

// zeroValueFnOut 返回值列表填充空值
func zeroValueFnOut(fooPtr interface{}) []reflect.Value {
	fn := reflect.ValueOf(fooPtr).Elem()
	outNum := fn.Type().NumOut()
	zeroOut := make([]reflect.Value, 0, outNum)
	for i := 0; i < outNum; i++ {
		zeroOut = append(zeroOut, reflect.Zero(fn.Type().Out(i)))
	}
	return zeroOut
}
