package analysis

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	ZookeeperClusterAddress []string `ini:"zookeeper_cluster_address"`
	InsecurePort            uint     `ini:"insecure_port"`
	RootDirectory           string   `ini:"root_directory"`
}

type analysis struct{}

func (a *analysis) reflectSlice(lineS []string, config *Config) {
	// 如果不包含","，认为不是切片值就退出
	valueOfINI := lineS[1]
	if !strings.Contains(valueOfINI, ",") {
		return
	}
	// “=”的0索引是字段名，应该与struct tag名一样
	fieldOfINI := lineS[0]
	// 获取*config（必须解引用）反射后的动态类型
	c := reflect.TypeOf(*config)
	// 从c中获取字段数量并遍历
	for i := 0; i < c.NumField(); i++ {
		field := c.Field(i)
		// 从每个字段中获取ini tag
		tag := field.Tag.Get("ini")
		// 如果ini tag与配置中“=”左边的字段名一样
		if tag == fieldOfINI {
			// 获取&Config(指针)反射后的动态值，valueOf必须是指针类型不然会panic
			v := reflect.ValueOf(config)
			// 从动态值中获取从动态类型中拿到的field.Name
			fieldOfSlice := v.Elem().FieldByName(field.Name)
			// 以“，”分割为切片
			value := strings.Split(valueOfINI, ",")
			// 将切片转为reflect.value类型
			fieldOfSliceValue := reflect.ValueOf(value)
			// 只有转成reflect.value，才能Set到reflect.value类型的fieldOfSlice变量中
			fieldOfSlice.Set(fieldOfSliceValue)
		}
	}
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
	as := analysis{}
	c := Config{}
	file, err := os.Open(ini)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("%w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()

		if as.checkTitle(line) {
			continue
		}

		// 把以“=”分割的切片传到每个analysis方法中，作用是直接可以在方法内判断
		// index 0 是字段，也是struct tag名
		// index 1 是值
		lineS := strings.Split(line, "=")
		// fmt.Println(lineS)
		// 分析配置文件中所有的切片
		as.reflectSlice(lineS, &c)

	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("%w", err)
	}
	return &c, nil
}
