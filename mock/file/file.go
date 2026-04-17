package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 定义目标目录
	dir := "/Users/liujun/Downloads/shop/category/2025/03/21" // 替换为你的目录路径

	// 读取目录中的所有条目
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("读取目录失败: %v\n", err)
		return
	}
	// 遍历所有条目
	for _, entry := range entries {
		// 跳过子目录，只处理文件
		if entry.IsDir() {
			continue
		}

		// 获取旧文件名和路径
		oldName := entry.Name()
		oldPath := filepath.Join(dir, oldName)

		// 生成新文件名（示例规则：添加前缀、替换空格）
		newName := strings.ReplaceAll(oldName, "..", ".")
		newPath := filepath.Join(dir, newName)

		// 执行重命名
		err = os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Printf("重命名失败: %s -> %s, 错误: %v\n", oldName, newName, err)
		} else {
			fmt.Printf("重命名成功: %s -> %s\n", oldName, newName)

		}
	}
}
