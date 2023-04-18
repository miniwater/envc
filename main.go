package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const (
	version = "1.0.0"
	env_var = "Envc"
)

func init() {

}

func main() {
	args := os.Args

	if len(args) == 1 {
		Catalogue()
		return
	}

	switch args[1] {
	case "-v", "--v", "version":
		fmt.Println("v" + version)
	case "creat":
		if len(args) == 3 {
			CrearFolder(args[2])
		} else {
			fmt.Println("请输入要创建的文件名")
		}
	default:
		ModifyEnv(args[1])
	}
}

// 输出目录
func Catalogue() {
	files, _ := ioutil.ReadDir(absolutePath())
	for _, file := range files {
		if file.IsDir() {
			fmt.Println(file.Name())
		}
	}
}

func CrearFolder(file_name string) {
	err := os.Mkdir("./"+file_name, os.ModePerm)
	if err != nil {
		fmt.Println(file_name + "目录已存在")
	}
}

func ModifyEnv(file_name string) {
	dir := absolutePath()

	files, _ := ioutil.ReadDir(dir + "/" + file_name)
	if len(files) == 0 {
		fmt.Println("找不到目录或目录下没有任何可切换版本")
		return
	}

	pathList := getEnv2()
	i := 1

	var folder []string
	for _, file := range files {
		if file.IsDir() {
			folder = append(folder, file.Name())
			for _, path := range pathList {
				if path == dir+"\\"+file_name+"\\"+file.Name() {
					fmt.Print("-")
					break
				}
			}
			fmt.Println("[", i, "]:", file.Name())
			i++
		}
	}

	var num int
	fmt.Println("请输入要切换环境的序号，按0退出")
	fmt.Scanln(&num)
	if num == 0 || num > len(folder) {
		return
	}

	newPath := dir + "\\" + file_name + "\\" + folder[num-1]

	chang := 0
	for i := 0; i < len(pathList); i++ {
		if strings.Contains(pathList[i], dir+"\\"+file_name) {
			pathList[i] = newPath
			chang = 1
		}
	}
	if chang == 0 {
		pathList = append(pathList, newPath)
	}
	pathList = removeDuplicates(pathList)
	newEnv := strings.Join(pathList, ";") + ";"
	cmd := exec.Command("setx", env_var, newEnv)
	CheckErr(cmd.Run())

}

// 获取环境变量
func getEnv() []string {
	env := os.Getenv(env_var)
	pathList := strings.Split(env, ";")
	newStrs := []string{}
	for _, str := range pathList {
		if str != "" {
			newStrs = append(newStrs, str)
		}
	}
	return newStrs
}

// 获取环境变量-注册表
func getEnv2() []string {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE)
	CheckErr(err)
	defer key.Close()

	// 获取所有的字符串值
	values, err := key.ReadValueNames(0)
	CheckErr(err)

	// 获取每个字符串值对应的数据
	env := ""
	for _, name := range values {
		if name == env_var {
			env, _, err = key.GetStringValue(name)
			CheckErr(err)
			break
		}
	}
	pathList := strings.Split(env, ";")
	newStrs := []string{}
	for _, str := range pathList {
		if str != "" {
			newStrs = append(newStrs, str)
		}
	}
	return newStrs
}

// 错误处理函数
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

// 去重
func removeDuplicates(nums []string) []string {
	m := make(map[string]bool)
	result := []string{}
	for _, num := range nums {
		if m[num] == false {
			m[num] = true
			result = append(result, num)
		}
	}
	return result
}

func absolutePath() string {
	exePath, err := os.Executable()
	CheckErr(err)
	return filepath.Dir(exePath)
}
