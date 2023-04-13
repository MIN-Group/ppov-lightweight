package Client

import (
    "fmt"
    "io"
    "io/ioutil"
    "os"
)

// 判断文件夹是否存在
func FoldExists(foldPath string) (bool, error) {
    _, err := os.Stat(foldPath)
    if err == nil {
		//fmt.Printf("目录存在[%v]\n", foldPath)
        return true, nil
    }
    if os.IsNotExist(err) {
		//fmt.Printf("目录不存在[%v]\n", foldPath)
        // 创建文件夹
        err := os.Mkdir(foldPath, os.ModePerm)
        if err != nil {
            //fmt.Printf("创建目录失败[%v]\n", err)
            return false, err
        } else {
			//fmt.Printf("创建目录成功\n")
			return true, nil
        }
    }
    return false, err
}


//创建公钥/私钥txt文件
func TxtExists(foldPath string, fileName string) (bool, error) {
	_, err := os.Stat("./" + foldPath + "/" + fileName + ".txt")
    if err == nil {
		//fmt.Printf("文件存在!\n")
        return true, nil
    }
    if os.IsNotExist(err) {
		//fmt.Printf("no txt!" + fileName + ".txt\n")
        // 创建文件
        _, err1 := os.Create("./" + foldPath + "/" + fileName + ".txt")
        if err1 != nil {
			//fmt.Printf("创建失败[%v]\n", err)
			return false, err1
        } else {
            //fmt.Printf("创建成功!\n")
			return true, nil
        }
    }
    return false, err
}


//往txt文件里写公钥/私钥
func FileWrite(fileName string, writeString string) (bool, error){
    fileInfo, err := os.OpenFile(fileName, os.O_RDWR, 0666)
    if err == nil {
        _, err1 := io.WriteString(fileInfo, writeString)
        if err1 != nil {
            //fmt.Printf("保存失败")
            //fmt.Println(err1) 
            return false, err1
            
        }else{
            fileInfo.Close()
            //fmt.Printf("保存成功")
            return true, nil
        }
    }else{
        //fmt.Printf("保存失败")
        return false, err
    }
}


//保存公钥和私钥
func KeyTxtWrite(foldPath string, fileName string, writeString string) {
    FoldExists("./" + foldPath)
    TxtExists(foldPath, fileName)
    FileWrite("./" + foldPath + "/"  + fileName + ".txt", writeString)
    
}


// 读取txt中的公钥/私钥
func KeyTxtRead(foldPath string, fileName string) (string) {
    path := "./" + foldPath + "/" + fileName + ".txt"
    //fmt.Println(path)
    _, err := os.Stat(path)
    if err != nil {
        fmt.Println(err)
    }
    if err == nil {
        //fmt.Printf("文件存在[%v]\n", fileName)
        data, _ := ioutil.ReadFile(path)
        keytxt := string(data)
        return keytxt
    }
    if os.IsNotExist(err) {
        //fmt.Printf("文件不存在[%v]\n", fileName)
        return ""
    }
    return ""
}