package zerorpc

import (
	"context"
	"fmt"
	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/sirupsen/logrus"
	"net"
	"sync/atomic"
)

type MyPooledObj struct {
	session *Session
}

type MyPooledObjFactory struct {
	serverAddr string
}

var count uint32 = 0
func (f *MyPooledObjFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	conn, err := net.Dial("tcp", f.serverAddr)
	if err != nil {
		logrus.Errorf("create conn err:%v, ctx:%v", err, ctx)
		return nil, err
	}
	atomic.AddUint32(&count, 1)
	fmt.Printf("make obj count:%d \n", count)
	return pool.NewPooledObject(
			&MyPooledObj{
				session: NewSession(conn),
			}),
		nil
}

func (f *MyPooledObjFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	// do destroy
	myObj := object.Object.(*MyPooledObj)
	logrus.Debugf("sessoin in poll destroyed, obj:%s, ctx:%v", myObj.session.conn.LocalAddr(), ctx)
	return myObj.session.Close()
}

func (f *MyPooledObjFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	// do validate
	myObj := object.Object.(*MyPooledObj)
	logrus.Infof("sessoin in pool destroyed, obj:%s, ctx:%v, session destoryed:%v", myObj.session.conn.LocalAddr(), ctx, myObj.session.destroyed)
	return !myObj.session.destroyed
}

func (f *MyPooledObjFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	// do activate
	myObj := object.Object.(*MyPooledObj)
	logrus.Debugf("session in pool activate, obj:%s, ctx:%v", myObj.session.conn.LocalAddr(), ctx)
	return nil
}

func (f *MyPooledObjFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	// do passivate(put into idle list)
	myObj := object.Object.(*MyPooledObj)
	logrus.Debugf("session in pool passivate, obj:%s, ctx:%v", myObj.session.conn.LocalAddr(), ctx)
	return nil
}
