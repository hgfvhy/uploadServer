package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	UPLOAD_DIR   = "uploads"
	MAX_SIZE     = 10 * 1024 * 1204
	MAX_SIZE_DES = "10MB"
	SUCCESS      = 1
	FAILE        = 0
	MAX_LEN      = 10 //最大同时可上传文件数
)

type resp struct {
	Code    int    //1上传成功   0上传失败
	Message string //返回信息
	Path    string //上传成功时的文件地址
	Status  int    //http状态码
}

type mulResp struct {
	Code    int //1上传成功   0上传失败
	Message map[int]string
	Path    map[int]string //上传成功时的文件地址
	Status  int            //http状态码
}

var fileTypeMap sync.Map
var extMap sync.Map

func init() {
	//以下文件头是网上找的，不正确，需要自己测试后手动添加
	fileTypeMap.Store("ffd8ffe000104a464946", "jpg")  //JPEG (jpg)
	fileTypeMap.Store("ffd8ffe1013e45786966", "jpg")  //JPEG (jpg)
	fileTypeMap.Store("89504e470d0a1a0a0000", "png")  //PNG (png)
	fileTypeMap.Store("4749463839610005d002", "gif")  //GIF (gif) 动图
	fileTypeMap.Store("d0cf11e0a1b11ae10000", "xls")  //wps 的doc和xls一致
	fileTypeMap.Store("504b03040a0000000000", "docx") //wps 部分的docx和xlsx一致
	fileTypeMap.Store("504b0304140006000800", "xlsx")
	fileTypeMap.Store("255044462d312e350d0a", "pdf") //Adobe Acrobat (pdf)
	//fileTypeMap.Store("49492a00227105008037", "tif")  //TIFF (tif)
	//fileTypeMap.Store("424d228c010000000000", "bmp")  //16色位图(bmp)
	//fileTypeMap.Store("424d8240090000000000", "bmp")  //24位位图(bmp)
	//fileTypeMap.Store("424d8e1b030000000000", "bmp")  //256色位图(bmp)
	//fileTypeMap.Store("41433130313500000000", "dwg")  //CAD (dwg)
	//fileTypeMap.Store("3c21444f435459504520", "html") //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	//fileTypeMap.Store("3c68746d6c3e0", "html")        //HTML (html)   3c68746d6c3e0  3c68746d6c3e0
	//fileTypeMap.Store("3c21646f637479706520", "htm")  //HTM (htm)
	//fileTypeMap.Store("48544d4c207b0d0a0942", "css")  //css
	//fileTypeMap.Store("696b2e71623d696b2e71", "js")   //js
	//fileTypeMap.Store("7b5c727466315c616e73", "rtf")  //Rich Text Format (rtf)
	//fileTypeMap.Store("38425053000100000000", "psd")  //Photoshop (psd)
	//fileTypeMap.Store("46726f6d3a203d3f6762", "eml")  //Email [Outlook Express 6] (eml)
	//fileTypeMap.Store("d0cf11e0a1b11ae10000", "vsd")  //Visio 绘图
	//fileTypeMap.Store("5374616E64617264204A", "mdb")  //MS Access (mdb)
	//fileTypeMap.Store("252150532D41646F6265", "ps")
	//fileTypeMap.Store("2e524d46000000120001", "rmvb") //rmvb/rm相同
	//fileTypeMap.Store("464c5601050000000900", "flv")  //flv与f4v相同
	//fileTypeMap.Store("00000020667479706d70", "mp4")
	//fileTypeMap.Store("49443303000000002176", "mp3")
	//fileTypeMap.Store("000001ba210001000180", "mpg") //
	//fileTypeMap.Store("3026b2758e66cf11a6d9", "wmv") //wmv与asf相同
	//fileTypeMap.Store("52494646e27807005741", "wav") //Wave (wav)
	//fileTypeMap.Store("52494646d07d60074156", "avi")
	//fileTypeMap.Store("4d546864000000060001", "mid") //MIDI (mid)
	//fileTypeMap.Store("504b0304140000000800", "zip")
	//fileTypeMap.Store("526172211a0700cf9073", "rar")
	//fileTypeMap.Store("235468697320636f6e66", "ini")
	//fileTypeMap.Store("504b03040a0000080000", "jar")
	//fileTypeMap.Store("4d5a9000030000000400", "exe")        //可执行文件
	//fileTypeMap.Store("3c25402070616765206c", "jsp")        //jsp文件
	//fileTypeMap.Store("4d616e69666573742d56", "mf")         //MF文件
	//fileTypeMap.Store("3c3f786d6c2076657273", "xml")        //xml文件
	//fileTypeMap.Store("494e5345525420494e54", "sql")        //xml文件
	//fileTypeMap.Store("7061636b616765207765", "java")       //java文件
	//fileTypeMap.Store("406563686f206f66660d", "bat")        //bat文件
	//fileTypeMap.Store("1f8b0800000000000000", "gz")         //gz文件
	//fileTypeMap.Store("6c6f67346a2e726f6f74", "properties") //bat文件
	//fileTypeMap.Store("cafebabe0000002e0041", "class")      //bat文件
	//fileTypeMap.Store("49545346030000006000", "chm")        //bat文件
	//fileTypeMap.Store("04000000010000001300", "mxp")        //bat文件
	//fileTypeMap.Store("d0cf11e0a1b11ae10000", "wps")        //WPS文字wps、表格et、演示dps都是一样的
	//fileTypeMap.Store("6431303a637265617465", "torrent")
	//fileTypeMap.Store("6D6F6F76", "mov")         //Quicktime (mov)
	//fileTypeMap.Store("FF575043", "wpd")         //WordPerfect (wpd)
	//fileTypeMap.Store("CFAD12FEC5FD746F", "dbx") //Outlook Express (dbx)
	//fileTypeMap.Store("2142444E", "pst")         //Outlook (pst)
	//fileTypeMap.Store("AC9EBD8F", "qdf")         //Quicken (qdf)
	//fileTypeMap.Store("E3828596", "pwl")         //Windows Password (pwl)
	//fileTypeMap.Store("2E7261FD", "ram")         //Real Audio (ram)

	//设置可上传的后缀类型
	extMap.Store("jpg", "jpg")
	extMap.Store("png", "png")
	extMap.Store("gif", "gif")
	extMap.Store("doc", "doc")
	extMap.Store("docx", "docx")
	extMap.Store("xls", "xls")
	extMap.Store("xlsx", "xlsx")
	extMap.Store("pdf", "pdf")
}

func main() {
	//go http.HandleFunc("/upload", uploadHandler)
	go http.HandleFunc("/upload", multipartUploadHandler)
	err := http.ListenAndServe(":9080", nil)
	if err != nil {
		writeLog(err.Error())
		log.Fatal("ListenAndServe: ", err.Error())
	}

}

//单文件接收主流程
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// 处理图片上传
	if r.Method == "POST" {
		f, h, err := r.FormFile("file")
		defer f.Close()
		if err != nil {
			sendMsg(w, FAILE, err.Error(), "")
			return
		}
		filename := h.Filename
		fmt.Println(filename)
		//return
		size := h.Size
		if size == 0 {
			sendMsg(w, FAILE, "该文件为空文件", "")
			return
		}
		if size > MAX_SIZE {
			sendMsg(w, FAILE, "最大可支持上文大小为"+MAX_SIZE_DES, "")
			return
		}

		//获取后缀
		var headerByte []byte
		headerByte = make([]byte, 10)
		_, err = f.ReadAt(headerByte, 0)
		if err != nil {
			sendMsg(w, FAILE, "无法读取文件", "")
			return
		}
		ext := GetFileType(headerByte)
		ext = getExcelExt(ext, filename)

		if !canUpload(ext) {
			sendMsg(w, FAILE, ext+"格式无法上传,如不正确,请联系rebirth添加", "")
			return
		}
		path := r.FormValue("path")
		if path == "" {
			path = UPLOAD_DIR
		}
		id := r.FormValue("id")
		if id == "" {
			id = "0"
		}
		path = getStr(path)
		id = getStr(id)
		fmt.Println(id)
		timestamp := time.Now().Unix()
		timeNow := time.Unix(timestamp, 0)
		timeString := timeNow.Format("2006_01_02")
		path = "./" + path + "/" + timeString + "/" + id

		if isExist(path) == false {
			err := os.MkdirAll(path, 777)
			if err != nil {
				sendMsg(w, FAILE, "文件夹创建失败", "")
				return
			}
		}
		fmt.Println(path)
		unix := strconv.Itoa(int(time.Now().Unix()))
		randStr := GetRandomString(4)
		filename = unix + "_" + randStr + "." + filename
		fmt.Println(filename)

		t, err := os.Create(path + "/" + filename)
		defer t.Close()
		if err != nil {
			sendMsg(w, FAILE, err.Error(), "")
			return
		}
		if _, err := io.Copy(t, f); err != nil {
			sendMsg(w, FAILE, err.Error(), "")
			return
		}
		ip := ClientPublicIP(r)
		if ip == "" {
			ip = ClientIP(r)
		}
		writeLog("上传ip：" + ip + "----" + "文件:" + path + "/" + filename)
		sendMsg(w, SUCCESS, "", path+"/"+filename)
	}

}

//多文件接收主流程
func multipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	// 处理图片上传
	if r.Method == "POST" {
		message := make(map[int]string)
		paths := make(map[int]string)

		//首先创建文件夹
		path := r.FormValue("path")
		if path == "" {
			path = UPLOAD_DIR
		}
		path = getStr(path)
		id := r.FormValue("id")
		if id == "" {
			id = "0"
		}
		id = getStr(id)
		timestamp := time.Now().Unix()
		timeNow := time.Unix(timestamp, 0)
		timeString := timeNow.Format("2006_01_02")
		path = "./" + path + "/" + timeString + "/" + id
		if isExist(path) == false {
			err := os.MkdirAll(path, 777)
			if err != nil {
				message[-1] = "文件夹创建失败"
				sendMulMsg(w, FAILE, message, paths)
				return
			}
		}
		unix := strconv.Itoa(int(time.Now().Unix()))
		ip := ClientPublicIP(r)
		if ip == "" {
			ip = ClientIP(r)
		}

		files := r.MultipartForm.File["file"]
		fmt.Println("files",files)
		lens := len(files)
		fmt.Println("lens",lens)
		if lens > MAX_LEN {
			message[-2] = "最多只可同时上传" + string(MAX_LEN) + "个文件"
			sendMulMsg(w, FAILE, message, paths)
			return
		}

		if lens == 0 {
			message[-3] = "请选择要上传的文件"
			sendMulMsg(w, FAILE, message, paths)
			return
		}

		for i := 0; i < lens; i++ {
			h := files[i]
			file, err := h.Open()
			defer file.Close()
			if err != nil {
				message[i] = err.Error()
				continue
			}
			filename := h.Filename
			size := h.Size
			if size == 0 {
				message[i] = "该文件为空文件"
				continue
			}
			if size > MAX_SIZE {
				message[i] = "最大可支持上文大小为" + MAX_SIZE_DES
				continue
			}
			//获取后缀
			headerByte := make([]byte, 10)
			_, err = file.ReadAt(headerByte, 0)
			if err != nil {
				message[i] = "无法读取文件"
				continue
			}
			ext := GetFileType(headerByte)
			ext = getExcelExt(ext, filename)
			if !canUpload(ext) {
				message[i] = ext + "格式无法上传,如不正确,请联系rebirth添加"
				continue
			}
			randStr := GetRandomString(4)
			filename = unix + "_" + randStr + "." + ext
			t, err := os.Create(path + "/" + filename)
			defer t.Close()
			if err != nil {
				message[i] = err.Error()
				continue
			}
			if _, err := io.Copy(t, file); err != nil {
				message[i] = err.Error()
				continue
			}
			paths[i] = path + "/" + filename
			writeLog("上传ip：" + ip + "----" + "文件:" + path + "/" + filename)
		}

		if len(message) == 0 {
			sendMulMsg(w, SUCCESS, message, paths)
		} else {
			sendMulMsg(w, FAILE, message, paths)
		}
	}

}

/**
 * 写入日志
 * add by rebirth 2019/06/14
 */
func writeLog(error string) {
	f, err := os.OpenFile("filename.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// 完成后延迟关闭
	defer f.Close()
	//设置日志输出到 f
	log.SetOutput(f)
	//写入日志内容
	log.Println(error)
}

// ClientIP 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// ClientPublicIP 尽最大努力实现获取客户端公网 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != "" {
			return ip
		}
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

//统一返回格式
//1上传成功   0上传失败
func sendMsg(w http.ResponseWriter, code int, err string, path string) {
	msg := resp{code, err, path, 200}
	json, _ := json.Marshal(msg)
	w.Write(json)
}

func sendMulMsg(w http.ResponseWriter, code int, err map[int]string, path map[int]string) {
	msg := mulResp{code, err, path, 200}
	json, _ := json.Marshal(msg)
	w.Write(json)
}

//字符串处理
func getStr(path string) string {
	path = strings.Replace(path, "\n", "", -1)  //去除换行符
	path = strings.Replace(path, "./", "", -1)  //去除 ./
	path = strings.Replace(path, "../", "", -1) //去除 ../
	path = strings.Replace(path, " ", "", -1)   //去除 空格
	return path
}

//生成随机字符串
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func isExist(path string) (bool) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

// 获取前面结果字节的二进制
func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}

// 用文件前面几个字节来判断
// fSrc: 文件字节流（就用前面几个字节）
func GetFileType(fSrc []byte) string {
	var fileType string
	fileType = "unKnow"
	fileCode := bytesToHexString(fSrc)
	fmt.Println("文件字符串为：" + fileCode)
	fileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
			strings.HasPrefix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	return fileType
}

func canUpload(ext string) bool {
	can := false
	extMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		if ext == strings.ToLower(k) {
			can = true
			return false
		}
		return true
	})

	return can
}

//解决wps 的文件头相同问题
func getExcelExt(ext string, filename string) string {
	if ext == "doc" || ext == "xls" {
		nameExt := path.Ext(filename)
		if nameExt == ".doc" {
			ext = "doc"
		}
		if nameExt == ".xls" {
			ext = "xls"
		}
	}
	if ext == "docx" || ext == "xlsx" {
		nameExt := path.Ext(filename)
		if nameExt == ".docx" {
			ext = "docx"
		}
		if nameExt == ".xlsx" {
			ext = "xlsx"
		}
	}
	return ext
}

//文件流判断文件头识别类型    https://www.cnblogs.com/enjong/articles/10741244.html
//获取用户ip    https://blog.thinkeridea.com/201903/go/get_client_ip.html
