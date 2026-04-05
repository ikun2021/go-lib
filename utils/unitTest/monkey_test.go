package unitTest

import (
	"encoding/json"
	"errors"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/luxun9527/zlog"
	"reflect"
	"testing"
)

func TestCompute_Mock(t *testing.T) {
	// ApplyFunc 第一个参数是目标函数，第二个参数是桩函数
	// 注意：桩函数的签名必须与原函数完全一致
	patches := gomonkey.ApplyFunc(Compute, func(a, b int) int {
		return 100 // 无论输入什么，都返回 100
	})
	defer patches.Reset() // 测试结束必须复原

	res := Compute(1, 2)
	zlog.Infof("%v", res)
	assert.Equal(t, 100, res)
}
func TestValidateJSON_MockError(t *testing.T) {
	// Mock 标准库 json.Unmarshal 让它强制返回错误
	patches := gomonkey.ApplyFunc(json.Unmarshal, func(data []byte, v interface{}) error {
		return errors.New("mock error")
	})
	defer patches.Reset()

	isValid := ValidateJSON([]byte(`{"name": "test"}`))
	assert.False(t, isValid) // 期望返回 false
}
func TestUser_FetchData(t *testing.T) {
	var u *User

	// ApplyMethod 参数：
	// 1. reflect.TypeOf(receiver): 获取接收者的类型（如果是指针接收者，需要传指针）
	// 2. "FetchData": 方法名字符串
	// 3. func(...): 桩函数。注意：桩函数的第一个参数必须是接收者（这里是 *User），后面才是原方法的参数
	patches := gomonkey.ApplyMethod(reflect.TypeOf(u), "FetchData", func(_ *User, url string) (string, error) {
		return "mocked data", nil
	})
	defer patches.Reset()

	user := &User{Name: "Jack"}
	data, err := user.FetchData("http://invalid-url.com")

	assert.NoError(t, err)
	assert.Equal(t, "mocked data", data)
}
