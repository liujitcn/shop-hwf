package core

import (
	"gitee.com/liujit/shop/server/internal/version"
	"gitee.com/liujit/shop/server/pkg/service/file/biz"
	"go.newcapec.cn/nctcommon/nmslib"
	"go.newcapec.cn/nctcommon/nmslib/basedata"
	nmskitAuthData "go.newcapec.cn/ncttools/nmskit-auth/data"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/cache"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/pprof"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/queue"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/sqldb"
	"sync"
	"time"

	"go.newcapec.cn/ncttools/nmskit/log"
	"go.newcapec.cn/ncttools/nmskit/transport/grpc"

	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

type ShopCore struct {
	SqlDb         sqldb.SqlDb
	Cache         cache.Cache
	Queue         queue.Queue
	FileCase      *biz.FileCase
	UserToken     *nmskitAuthData.UserToken
	grpcServer    *grpc.Server        `desc:"grpc服务器的实例"`
	ServiceStatus basedata.WorkStatus //终端服务状态
	quitChan      chan struct{}       //退出Chan
	closeOnce     sync.Once
	taskTimer     *time.Timer
	rwLock        sync.RWMutex //异步数据锁
}

// NewShopCore create a device service core data struct.
func NewShopCore(
	sqlDb sqldb.SqlDb,
	cache cache.Cache,
	queue queue.Queue,
	pprof pprof.Pprof,
	userToken *nmskitAuthData.UserToken,
	FileCase *biz.FileCase,
) (*ShopCore, func(), error) {
	log.Infof("建立%s服务支持数据...", version.Name)

	// 设置全局变量
	nmslib.Runtime.SetCache(cache)
	nmslib.Runtime.SetQueue(queue)

	// 启动服务监控
	if pprof != nil {
		pprof.Start()
	}

	usc := ShopCore{
		SqlDb:     sqlDb,
		Cache:     cache,
		Queue:     queue,
		UserToken: userToken,
		FileCase:  FileCase,

		grpcServer:    nil,
		ServiceStatus: basedata.WorkStatus_Stop,
		quitChan:      make(chan struct{}),
		closeOnce:     sync.Once{},
		taskTimer:     nil,
		rwLock:        sync.RWMutex{},
	}

	// 启动后台服务
	go usc.Serve()

	cleanup := func() {
		usc.Close()
		if pprof != nil {
			pprof.Stop()
		}
	}
	return &usc, cleanup, nil
}

func (sc *ShopCore) Close() {
	sc.closeOnce.Do(func() {
		if sc.taskTimer != nil {
			sc.taskTimer.Stop()
		}
		close(sc.quitChan)
		log.Warnf("%s shutdown!", version.Name)
	})
}

// Serve 缓存加载和刷新线程
func (sc *ShopCore) Serve() {
	// 启动队列
	sc.Queue.Run()
	//循环处理同步事件
	for {
		select {
		case <-sc.quitChan:
			log.Warnf("%s exit ----------", version.Name)
			return
		}
	}
}

func (sc *ShopCore) SetServiceStatus(st basedata.WorkStatus) {
	if sc.ServiceStatus != st {
		sc.ServiceStatus = st
		log.Infof("%s 服务状态：%s", version.Name, st.String())
		var healthStatus int32 = 0
		if sc.ServiceStatus == basedata.WorkStatus_Running {
			healthStatus = 1 // healthgrpc.HealthCheckResponse_SERVING
		} else {
			healthStatus = 2 // healthgrpc.HealthCheckResponse_NOT_SERVING
		}
		sc.GrpcSetHealthStatus(healthStatus)
	}
}

func (sc *ShopCore) SetGrpcServer(grpcServer *grpc.Server) {
	sc.grpcServer = grpcServer
}

func (sc *ShopCore) GrpcSetHealthStatus(hs int32) {
	log.Info("GrpcSetHealthStatus : ", log.Cyan, healthgrpc.HealthCheckResponse_ServingStatus(hs), log.Reset)
	sc.grpcServer.SetHealthStatus(version.Name, hs)
}
