package analysis

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	ZookeeperClusterAddress []string `ini:"zookeeper_cluster_address"`
	InsecurePort            uint     `ini:"insecure_port"`
	RootDirectory           string   `ini:"root_directory"`
}

type analysis struct {
	fieldOfIni string
	valueOfIni string
}

// 将分隔符为=号切割后的两列，返回第一列是字段，第二列是值
func fieldAndValue(s []string) (f string, v string) {
	return s[0], s[1]
}

// 分析切片
func (a *analysis) identifySlice(config *Config) bool {
	// 如果不包含","，认为不是切片值就退出
	if !strings.Contains(a.valueOfIni, ",") {
		return false
	}
	value := strings.Split(a.valueOfIni, ",")
	a.reflectSetToIni(config, value)
	return true

}

// 分析int
func (a *analysis) identifyInt(config *Config) bool {
	valueOfInt, err := strconv.Atoi(a.valueOfIni)
	if err != nil {
		return false
	}
	// 出于练习并没有做到很多类型的支持，只练习一种配置文件，所以直接转成uint来处理
	value := uint(valueOfInt)
	a.reflectSetToIni(config, value)
	return true
}

// 分析string
func (a *analysis) identifyString(config *Config) {
	a.reflectSetToIni(config, a.valueOfIni)
}

// 初始化reflect实例并赋值给*Config
func (a *analysis) reflectSetToIni(config *Config, x interface{}) bool {
	// “=”的0索引是字段名，应该与struct tag名一样
	// 获取*config（必须解引用）反射后的动态类型
	c := reflect.TypeOf(*config)
	// 从c中获取字段数量并遍历
	for i := 0; i < c.NumField(); i++ {
		field := c.Field(i)
		// 从每个字段中获取ini tag
		tag := field.Tag.Get("ini")
		// 如果ini tag与配置中“=”左边的字段名一样
		if tag == a.fieldOfIni {
			// 获取&Config(指针)反射后的动态值，valueOf必须是指针类型不然会panic
			v := reflect.ValueOf(config)
			// 从动态值中获取从动态类型中拿到的field.Name
			fieldOfSlice := v.Elem().FieldByName(field.Name)
			// 先声明一个reflect.Value结构体，对x进行断言后传参数
			// 必须符合被解析的结构体中的字段的值类型，如果是string的那就必须断言为string，然后转成reflect.Value，才能被Set()
			fieldOfSliceValue := reflect.Value{}
			switch x.(type) {
			case []string:
				fieldOfSliceValue = reflect.ValueOf(x)
			case uint:
				fieldOfSliceValue = reflect.ValueOf(x)
			case string:
				fieldOfSliceValue = reflect.ValueOf(x)
			}
			// 只有转成reflect.value，才能Set到reflect.value类型的fieldOfSlice变量中
			fieldOfSlice.Set(fieldOfSliceValue)
		}
	}
	return true
}

// 检查是否为[]标题
func (a *analysis) checkTitle(line string) bool {
	if strings.HasPrefix(line, "[") {
		return true
	} else {
		return false
	}
}

// 将ini配置文件反序列化赋值到结构体，返回结构体指针和error
func UnMarshalWithIniPath(ini string) (*Config, error) {
	a := analysis{}
	c := Config{}
	file, err := os.Open(ini)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("%w", err)
	}
	defer file.Close()

	// 使用bufio的scanner，可以自动识别每行的结尾符号，默认是'\n'，读取到的行的字符串中会把分隔符去掉
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()

		if a.checkTitle(line) {
			continue
		}

		// 把以“=”分割的切片传到每个analysis方法中，作用是直接可以在方法内判断
		// index 0 是字段，也是struct tag名
		// index 1 是值
		ls := strings.Split(line, "=")
		a.fieldOfIni, a.valueOfIni = fieldAndValue(ls)
		// 分析配置文件中所有的切片
		if ok := a.identifySlice(&c); ok {
			continue
		}
		if ok := a.identifyInt(&c); ok {
			continue
		}
		a.identifyString(&c)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("%w", err)
	}
	return &c, nil
}
