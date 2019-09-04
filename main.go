package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	filePath string
)

func parseFile(path string) (map[string]int, error) {
	result := make(map[string]int)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	for _, s := range f.Decls {
		switch v := s.(type) {
		case *ast.FuncDecl:
			var key string
			var recv []string
			var mf *ast.FuncDecl
			var pList []string
			var rList []string
			var recvT string
			mf = v
			p := mf.Type.Pos()
			pos := fset.Position(p)
			posF := pos.Line
			if mf.Recv != nil {
				for _, l := range mf.Recv.List {

					switch t := l.Type.(type) {
					case *ast.StarExpr:
						if id, ok := t.X.(*ast.Ident); ok {
							recvT = id.Name
						}

					}
					for _, i := range l.Names {
						recv = append(recv, i.String())
					}
				}
			}
			for _, l := range mf.Type.Params.List {
				for _, i := range l.Names {
					pList = append(pList, i.Name)
				}
			}

			ml := mf.Type.Results
			if ml != nil {
				for _, l := range ml.List {
					for _, i := range l.Names {
						rList = append(rList, i.Name)
					}
				}
			}
			if mf.Recv != nil {
				var str string
				for _, s := range recv {
					str = str + " " + s
				}
				key += fmt.Sprintf("(%s %s) ", str, recvT)
			}
			key += fmt.Sprintf("%s ", mf.Name)
			key += fmt.Sprintf("%v ", pList)
			if ml != nil {
				key += fmt.Sprintf("%v", rList)
			}
			result[key] = posF
		}
	}
	return result, nil
}

// func getAllFile(path string) ([]string, error) {
// 	var fList []string
// 	rd, err := ioutil.ReadDir(path)
// 	if err != nil {
// 		return nil, err
// 		fmt.Println(err)
// 	}
// 	for _, fi := range rd {
// 		if fi.IsDir() {
// 			fl, err := getAllFile(path + fi.Name())
// 			if err != nil {
// 				return nil, err
// 			}
// 			fList = append(fList, fl...)
// 		} else if strings.HasSuffix(fi.Name(), ".go") {
// 			fList = append(fList, path+fi.Name())
// 		} else {
// 			continue
// 		}
// 	}
// 	return fList, nil

// }
func getAllFile(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix)
	//忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() {
			// 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			getAllFile(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".go")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}
	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := getAllFile(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}
	return files, nil
}

func parsePro(path string, proname string) error {
	var commits []string
	commitfile := "/Users/ges/go/src/parsego/filename"
	cf, err := os.Open(commitfile)
	if err != nil {
		return err
	}
	defer cf.Close()
	br := bufio.NewReader(cf)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		commits = append(commits, string(a))
	}

	fi, err := os.OpenFile("./pro/"+proname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	fList, err := getAllFile(path)
	if err != nil {
		return err
	}
	for _, f := range fList {
		var mc string
		for _, c := range commits {
			if strings.Contains(c, proname) {
				mc = c
			}
		}
		mcList := strings.Split(mc, "-")
		mn := mcList[len(mcList)-1]
		m, err := parseFile(f)
		if err != nil {
			return err
		}
		for key, value := range m {
			kstr := proname + strings.TrimPrefix(f, path) + "|" + key
			v := "{c:" + mn + ",l:" + strconv.Itoa(value) + ",f:" + proname + "}\n"
			result := kstr + " -> " + v
			if _, err := fi.Write([]byte(result)); err != nil {
				return err
			}
		}

	}
	fi.Close()
	return nil
}

func main() {

	filePath := "/Users/ges/gop"
	fd, err := os.Open(filePath)
	defer fd.Close()
	if err != nil {
		log.Fatal(err)
	}

	list, err := fd.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range list {
		filename := filepath.Join(filePath, d.Name())
		parsePro(filename, d.Name())
	}
}

//func main() {
//	filePath = "/Users/ges/gop"
//	fd, err := os.Open(filePath)
//	defer fd.Close()
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	list, err := fd.Readdir(-1)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	for _, d := range list {
//		filename := filepath.Join(filePath, d.Name())
//		fset := token.NewFileSet()
//		pkgs, err := parser.ParseDir(fset, filename, nil, parser.ParseComments)
//		if err != nil {
//			fmt.Println(err)
//		}
//		for key, value := range pkgs {
//			fmt.Println(value.Name, key)
//		}
//	}
//}
