package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/labstack/gommon/log"
)

func TestNN(t *testing.T) {

	// 获取公钥
	input := "your input text"

	for i := 0; i < 2; i++ {
		// 1. 生成 MD5 哈希值
		hash := md5.New()
		hash.Write([]byte(input))
		hashBytes := hash.Sum(nil)

		// 2. 使用 Hex 编码将二进制哈希值转换为固定的字符串输出
		hashHex := hex.EncodeToString(hashBytes)

		fmt.Println("MD5 + Hex 编码后的固定输出:", hashHex)
	}
}

func TestOutputHtml(t *testing.T) {
	//遍历制定文件夹下的文件,以文件夹为名称，将内部图片文件插入到img标签中
	OutputHtml(`H:\wechatfiles\WeChat Files\yuyuhaso\FileStorage\File\2024-11\中华优秀传统文化概要11111\中华优秀传统文化概要`)
	log.Info("done")
}
func OutputHtml(path string) {
	indexfile, _ := os.Create("index.html")
	indexfile.WriteString("<body style=\"margin:0px; overflow: hidden;\">")
	//遍历文件夹目录
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if f.IsDir() {
			//遍历文件夹，按目录名称输出.html
			log.Info(f.Name())
			OutputHtml := f.Name() + ".html"
			//创建文件
			file, _ := os.Create(OutputHtml)

			//进入到文件夹 遍历jpg文件
			files2, _ := ioutil.ReadDir(path + "\\" + f.Name())
			for _, f2 := range files2 {
				if f2.Name()[len(f2.Name())-4:] == ".jpg" {
					file.WriteString("<img width=\"100%\" src=\"" + f.Name() + "\\" + f2.Name() + "\">")
					indexfile.WriteString("<img width=\"100%\" src=\"" + f.Name() + "\\" + f2.Name() + "\">")
				}
			}

			defer file.Close()

		} else {
			//遍历文件
			log.Info(f.Name())
		}

	}
	indexfile.WriteString("</body>")
	indexfile.Close()

}
