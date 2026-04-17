package start

import (
	"context"
	"fmt"
	"gitee.com/liujit/shop/server/internal/version"
	"gitee.com/liujit/shop/server/pkg/core"
	"gitee.com/liujit/shop/server/pkg/service/admin/biz"
	"github.com/spf13/cobra"
	"go.newcapec.cn/nctcommon/nmslib/basedata"
	"go.newcapec.cn/ncttools/nmskit"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/bootstrap"
	NmsVue "go.newcapec.cn/ncttools/nmskit-transport/transport/vue"
	"go.newcapec.cn/ncttools/nmskit/log"
	"go.newcapec.cn/ncttools/nmskit/registry"
	"go.newcapec.cn/ncttools/nmskit/transport"
	NmsGrpc "go.newcapec.cn/ncttools/nmskit/transport/grpc"
	NmsHttp "go.newcapec.cn/ncttools/nmskit/transport/http"
)

var (
	//YMAL配置文件路径
	configPath string
)

// CmdStart represents the new command.
var CmdStart = &cobra.Command{
	Use:   "start",
	Short: "Start a template",
	Long:  fmt.Sprintf("Start a template project. Example: NmsKit run %s start", version.Name),
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.Bootstrap(initApp, configPath, version.Name, version.Release)
	},
}

func init() {
	CmdStart.PersistentFlags().StringVarP(&configPath, "configs", "c", "./configs", "Start Server with provided configuration file")
}

func newApp(logger log.Logger, jobCase *biz.BaseJobCase, rr registry.Registrar,
	c *core.ShopCore, gs *NmsGrpc.Server, hs *NmsHttp.Server, vue *NmsVue.Server) *nmskit.App {
	c.SetGrpcServer(gs)
	var srv = []transport.Server{
		gs,
	}
	if hs != nil {
		srv = append(srv, hs)
	}
	if vue != nil {
		srv = append(srv, vue)
	}
	opts := make([]nmskit.Option, 0)
	opts = append(opts, nmskit.BeforeStart(func(ctx context.Context) error {
		log.Infof("BeforeStart %s ...", version.Name)
		// 启动定时任务
		err := jobCase.Init(ctx)
		if err != nil {
			log.Errorf("Init jobCase error: %s", err.Error())
			return err
		}
		return nil
	}))
	opts = append(opts,
		nmskit.BeforeStop(func(_ context.Context) error {
			log.Infof("BeforeStop %s ...", version.Name)
			c.SetServiceStatus(basedata.WorkStatus_Stop)
			c.Close()
			return nil
		}))
	opts = append(opts,
		nmskit.AfterStart(func(_ context.Context) error {
			log.Infof("AfterStart %s ...", version.Name)
			c.SetServiceStatus(basedata.WorkStatus_Running)
			return nil
		}))
	opts = append(opts,
		nmskit.AfterStop(func(_ context.Context) error {
			log.Infof("AfterStop %s ...", version.Name)
			return nil
		}))
	return bootstrap.NewApp(logger, rr, opts, srv...)
}
