package utils

import (
	"go.newcapec.cn/ncttools/nmskit/log"
	"gorm.io/gorm"
	"os"
	"strings"
)

var tableNameList = []string{"base_migration",
	"base_menu_api_rule",
	"base_role_menu",
	"base_role_dept",
	"base_backup_table",
	"base_user_group_member",
	"dev_group_nltype",
	"dev_group_member",
	"dev_code_create_history",
	"dev_runtime_history",
	"dev_collection_member",
	"product_price",
	"product_sale_item",
	"product_device_binding",
	"account_debit_rec",
	"account_debit_rec_bad",
	"account_base_subject_member",
	"idt_nl_assemble_data",
}

func ExecSql(db *gorm.DB, filePath string) error {
	sqlFile, err := readFile(filePath)
	if err != nil {
		log.Error(filePath+" 数据库基础数据初始化脚本读取失败！原因:", err.Error())
		return err
	}
	sqlList := strings.Split(sqlFile, ";")
	for _, sql := range sqlList {
		sql = strings.ReplaceAll(sql, "\r", "")
		sql = strings.ReplaceAll(sql, "\n", "")
		if len(strings.TrimSpace(sql)) == 0 {
			continue
		}
		if err = db.Exec(sql).Error; err != nil {
			if !strings.Contains(err.Error(), "Query was empty") {
				log.Errorf("error sql: %s", sql)
				return err
			}
		}
	}
	return nil
}

func readFile(filePath string) (string, error) {
	if contents, err := os.ReadFile(filePath); err == nil {
		//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
		result := string(contents)
		log.Debug("Use os.ReadFile to read a file:" + result)
		return result, nil
	} else {
		return "", err
	}
}
