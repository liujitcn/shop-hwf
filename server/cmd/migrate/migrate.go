package migrate

import (
	"bytes"
	"fmt"
	"gitee.com/liujit/shop/server/cmd/migrate/migration"
	"gitee.com/liujit/shop/server/cmd/migrate/models"
	"gitee.com/liujit/shop/server/internal/version"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/conf"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/config"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/logger"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/sqldb"
	"go.newcapec.cn/ncttools/nmskit-bootstrap/utils"
	"go.newcapec.cn/ncttools/nmskit/log"
	"os"
	"strconv"
	"text/template"
	"time"

	_ "gitee.com/liujit/shop/server/cmd/migrate/migration/version"
	_ "gitee.com/liujit/shop/server/cmd/migrate/migration/version-local"

	"github.com/spf13/cobra"
)

var (
	configPath string
	generate   bool
	goAdmin    bool
	host       string
)
var CmdStart = &cobra.Command{
	Use:     "migrate",
	Short:   "Initialize the database",
	Example: "ShopAdmin migrate -c /configs",
	Run: func(cmd *cobra.Command, args []string) {
		Run()
	},
}

func init() {
	CmdStart.PersistentFlags().StringVarP(&configPath, "config", "c", "./configs", "Start server with provided configuration file")
	CmdStart.PersistentFlags().BoolVarP(&generate, "generate", "g", false, "generate migration file")
	CmdStart.PersistentFlags().BoolVarP(&goAdmin, "ShopAdmin", "a", false, "generate ShopAdmin migration file")
	CmdStart.PersistentFlags().StringVarP(&host, "domain", "d", "*", "select tenant host")
	fmt.Println("ShopAdmin-Migrate Server Initial...")
}

func Run() {
	var err error
	// load configs
	if err = config.LoadBootstrapConfig(configPath); err != nil {
		panic(fmt.Sprintf("load config failed: %v", err))
	}
	bootstrapConfig := config.GetBootstrapConfig()

	// init metadata
	serviceInfo := utils.NewServiceInfo(version.Name, version.Release, "")
	serviceInfo.SetMataData(bootstrapConfig.Metadata)

	// init logger
	logger.NewLoggerProvider(bootstrapConfig.Logger, serviceInfo)

	err = migrateModel(bootstrapConfig.GetData().GetDatabase())
	if err != nil {
		panic(fmt.Sprintf("migrate model failed: %v", err))
	}

	log.Debug("Base-Migrate Server Run...")
	if !generate {
		log.Debug(`start init`)
	} else {
		log.Info(`generate migration file`)
		_ = genFile()
	}
}

func migrateModel(databaseConf *conf.Data_Database) error {
	log.Info("数据库迁移开始...")
	sqlDb, cleanup, err := sqldb.NewGormSqlDb(databaseConf)
	if err != nil {
		return err
	}
	defer cleanup()
	db := sqlDb.GetDb()
	if databaseConf.Driver == "mysql" {
		//初始化数据库时候用
		db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	}
	err = db.Debug().AutoMigrate(&models.Migration{})
	if err != nil {
		return err
	}
	migration.Migrate.SetHostDb(db.Debug())
	migration.Migrate.Migrate()
	log.Info(`数据库基础数据初始化成功!`)
	return err
}

func genFile() error {
	t1, err := template.ParseFiles("template/migrate.template")
	if err != nil {
		return err
	}
	m := map[string]string{}
	m["GenerateTime"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	m["Package"] = "version_local"
	if goAdmin {
		m["Package"] = "version"
	}
	var b1 bytes.Buffer
	err = t1.Execute(&b1, m)
	if goAdmin {
		fileCreate(b1, "./cmd/migrate/migration/version/"+m["GenerateTime"]+"_migrate.go")
	} else {
		fileCreate(b1, "./cmd/migrate/migration/version-local/"+m["GenerateTime"]+"_migrate.go")
	}
	return nil
}

func fileCreate(content bytes.Buffer, name string) {
	file, err := os.Create(name)
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Error(err)
		}
	}(file)
	if err != nil {
		log.Error(err)
	}
	_, err = file.WriteString(content.String())
	if err != nil {
		log.Error(err)
	}
}
