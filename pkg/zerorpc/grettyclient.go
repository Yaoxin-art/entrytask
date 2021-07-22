package zerorpc

import (
	"context"
	"github.com/jolestar/go-commons-pool/v2"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

type GrettyClient struct {
	serverAddr  string
	services    map[string]interface{}
	sessionPool pool.ObjectPool
}

func NewGrettyClient(serverAddr string, coreSize, maxSize int) *GrettyClient {
	// todo: test conn
	//factory := func() (net.Conn, error) {
	//	return net.Dial("tcp", serverAddr)
	//}
	ctx := context.Background()
	config := pool.ObjectPoolConfig{
		MaxTotal: maxSize,
		MaxIdle:  coreSize,
	}
	sessionPool := pool.NewObjectPool(ctx, &MyPooledObjFactory{
		serverAddr: serverAddr,
	}, &config)

	return &GrettyClient{
		serverAddr:  serverAddr,
		services:    make(map[string]interface{}),
		sessionPool: *sessionPool,
	}
}

// ConfigRemoteMethod 配置远程方法及提供者信息
// methodId:		远程方法Id
// fooPtr:			远程方法声明
func (pc *GrettyClient) ConfigRemoteMethod(methodId string, fooPtr interface{}) {
	if _, ok := pc.services[methodId]; !ok {
		// not config
		pc.services[methodId] = fooPtr
	}
	pc.pooledCallRPC(methodId, fooPtr)
	logrus.Infof("config methodId:%s -> %T, provider:%s success", methodId, fooPtr, pc.serverAddr)
}

// CallRPC 定义rpc代理
func (pc *GrettyClient) pooledCallRPC(methodId string, fooPtr interface{}) {
	fn := reflect.ValueOf(fooPtr).Elem()
	foo := func(args []reflect.Value) []reflect.Value {
		inArgs := make([]interface{}, 0, len(args))
		for _, arg := range args {
			inArgs = append(inArgs, arg.Interface())
		}
		// get remote session
		ctx := context.Background()
		pooledObj, err := pc.sessionPool.BorrowObject(ctx)
		defer func() {
			err := pc.sessionPool.ReturnObject(ctx, pooledObj)
			if err != nil {
				logrus.Errorf("GrettyClient put session err:%v", err)
			}
		}()
		if err != nil {
			logrus.Errorf("invoke remote method:%s get session form poll err:%v", methodId, err)
			return zeroValueFnOut(fooPtr)
		}
		session := pooledObj.(*MyPooledObj).session
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

// dialer session factory
//func dialer(serverAddr string) (*pool.Session, error) {
//	conn, err := net.Dial("tcp", serverAddr)
//	if err != nil {
//		return nil, err
//	}
//	return pool.NewSession(conn), nil
//}

func (pc *GrettyClient) Destroy() {
	logrus.Infof("Rpc GrettyClient destroy...")
}
