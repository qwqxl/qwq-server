package util

import (
	"os"
	"path/filepath"
)

// ExePath 获取可执行文件路径，并可拼接额外路径
// s: 可选的要拼接的路径片段
// 返回: 可执行文件路径或拼接后的完整路径
func ExePath(s ...string) string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// 解析符号链接获取真实路径
	realPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		panic(err)
	}

	// 拼接额外路径
	if len(s) > 0 {
		return filepath.Join(append([]string{realPath}, s...)...)
	}
	return realPath
}

// ExeDir 获取可执行文件所在目录，并可拼接额外路径
// s: 可选的要拼接的路径片段
// 返回: 可执行文件所在目录或拼接后的完整路径
func ExeDir(s ...string) string {
	exePath := ExePath()
	dir := filepath.Dir(exePath)

	if len(s) > 0 {
		return filepath.Join(append([]string{dir}, s...)...)
	}
	return dir
}

// WorkDir 获取当前工作目录，并可拼接额外路径
// s: 可选的要拼接的路径片段
// 返回: 当前工作目录或拼接后的完整路径
func WorkDir(s ...string) string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if len(s) > 0 {
		return filepath.Join(append([]string{wd}, s...)...)
	}
	return wd
}
