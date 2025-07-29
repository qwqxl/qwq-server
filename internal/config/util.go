package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindConfigPath 智能查找配置文件路径
// 优先级：1. 当前工作目录/configs  2. 可执行文件所在目录/configs
func FindConfigPath(filename string) (string, error) {
	// 尝试从当前工作目录查找
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("获取工作目录失败: %v", err)
	}
	currentConfigPath := filepath.Join(cwd, "configs", filename)
	if fileExists(currentConfigPath) {
		return currentConfigPath, nil
	}

	// 尝试从可执行文件所在目录查找
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("获取可执行文件路径失败: %v", err)
	}
	exeDir := filepath.Dir(exePath)
	exeConfigPath := filepath.Join(exeDir, "configs", filename)
	if fileExists(exeConfigPath) {
		return exeConfigPath, nil
	}

	// 全部查找失败
	return "", fmt.Errorf("配置文件 %s 未找到，请检查以下路径：\n- %s\n- %s",
		filename, currentConfigPath, exeConfigPath)
}

// 检查文件是否存在
func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
