package webui

import (
	"device/logs"
	"device/uuid"
	"device/websocket"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

//go:embed PetUI.js
var webPetui []byte

var (
	webPetuiLength string = strconv.FormatInt(int64(len(webPetui)), 10)
	upgrader              = websocket.Upgrader{}
)

// EventFn 函数接口
type EventFn func(arg any)

// Event 事件接口
type Event interface {
	Close()
	IsClose() bool
	Key() string
	GetConn() *websocket.Conn
	GetRequest() *http.Request
	SetEval(eval string)
	CallFunc(name string, value any)
	SetValue(name string, value any)
	SetValueMap(value map[string]any)
	GetValue(name string) (value any)
	DeleteValue(name string)
	SetFunc(name string, fn EventFn)
	DeleteFunc(name string)
}

// value 事件值
type value struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value any    `json:"value"`
}
type event struct {
	key       string
	request   *http.Request
	conn      *websocket.Conn
	isClose   bool
	funclist  sync.Map
	valueList sync.Map
}

func (eve *event) GetConn() *websocket.Conn {
	return eve.conn
}

func (eve *event) GetRequest() *http.Request {
	return eve.request
}

func (eve *event) Key() string {
	return eve.key
}

func (eve *event) Close() {
	eve.isClose = false
}

func (eve *event) IsClose() bool {
	return !eve.isClose
}

func (eve *event) CallFunc(name string, val any) {
	eve.conn.WriteJSON(&value{
		Type:  "CallFunc",
		Name:  name,
		Value: val,
	})
}

func (eve *event) SetValue(name string, val any) {
	eve.valueList.Store(name, val)
	eve.conn.WriteJSON(&value{
		Type:  "SetValue",
		Name:  name,
		Value: val,
	})
}
func (eve *event) SetValueMap(value map[string]any) {
	for key, value := range value {
		eve.SetValue(key, value)
	}
}

func (eve *event) SetEval(eval string) {
	eve.conn.WriteJSON(&value{
		Type:  "$",
		Value: eval,
	})
}

func (eve *event) GetValue(name string) (val any) {
	if v, ok := eve.valueList.Load(name); ok {
		return v
	}
	return nil
}

func (eve *event) DeleteValue(name string) {
	eve.valueList.Delete(name)
	eve.conn.WriteJSON(&value{
		Type: "DeleteValue",
		Name: name,
	})
}

func (eve *event) SetFunc(name string, fn EventFn) {
	eve.funclist.Store(name, fn)
	eve.conn.WriteJSON(&value{
		Type: "SetFunc",
		Name: name,
	})
}

func (eve *event) DeleteFunc(name string) {
	eve.funclist.Delete(name)
	eve.conn.WriteJSON(&value{
		Type: "DeleteFunc",
		Name: name,
	})
}

// PteUI UI框架
type PteUI struct {
	*http.ServeMux
	Pattern  string
	CallList map[string]func(petui Event, value any)
}

func (pteUI *PteUI) Bing(name string, fn func(petui Event, value any)) {
	pteUI.CallList[name] = fn
}

// NewPetUI 添加新界面
func NewPetUI(mux *http.ServeMux, pattern string, initFn func(event Event)) *PteUI {
	path, err := url.JoinPath(pattern)
	logs.Panicf("函数 PetUI.NewPetUI 路径定义错误: %v", err != nil, err)
	petui := &PteUI{ServeMux: mux, Pattern: path, CallList: map[string]func(petui Event, value any){}}
	// 界面初始化
	petui.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		// ID设置
		key, _ := r.Cookie("Key")
		if key == nil {
			id, _ := uuid.NewV4()
			http.SetCookie(w, &http.Cookie{
				Name:  "Key",
				Value: id.String(),
			})
		}
		w.Write([]byte(`<!DOCTYPE html><html><head><meta charset="utf-8">
		<meta content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0;" name="viewport"><script src="` + path + `.js"></script>
		<script>window.onload = ()=>{$=PetUIExpand.Bind({},{ url: PetUIExpand.URL+"ws" });};</script></head><body></body></html>`))
	})
	// 打包JS库
	petui.HandleFunc(path+".js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript;charset=utf-8")
		w.Header().Set("Content-Length", webPetuiLength)
		w.Write(webPetui)
	})
	// 界面通信过程
	petui.HandleFunc(path+"ws", func(w http.ResponseWriter, r *http.Request) {
		// 错误拦截到日志
		logs.Recover()
		// 升级协议
		conn, err := upgrader.Upgrade(w, r, nil)
		logs.Panicf("函数 PetUI.AddUI->Upgrade[IP:%v] 协议升级失败: %v", err != nil, r.RemoteAddr, err)
		defer conn.Close()
		// 事件结构体
		key, _ := r.Cookie("Key")
		eve := &event{request: r, conn: conn, isClose: true, key: key.Value}
		initFn(eve)
		val := &value{}
		val.Type = "onopen"
		for name := range petui.CallList {
			eve.SetValue(name, nil)
		}
		eve.conn.WriteJSON(val)
		for eve.isClose {
			messageType, message, err := conn.ReadMessage()
			if messageType == -1 {
				eve.isClose = false
				return
			}
			logs.Panicf("函数 PetUI.AddUI->ReadMessage[IP:%v] 读取数据失败: %v", err != nil, r.RemoteAddr, err)
			json.Unmarshal(message, val)
			// 处理事件
			switch val.Type {
			case "SetValue":
				eve.valueList.Store(val.Name, val.Value)
				// 拦截前端的数据修改
				if fn, ok := petui.CallList[val.Name]; ok {
					fn(eve, val.Value)
				}
			case "DeleteFunc":
				eve.funclist.Delete(val.Name)
			case "DeleteValue":
				eve.valueList.Delete(val.Name)
			case "CallFunc":
				if fn, ok := eve.funclist.Load(val.Name); ok {
					(fn.(EventFn))(val.Value)
				}
			default:
				logs.Printf("函数 PetUI.AddUI 异常的访问: %v", true, r.RemoteAddr)
			}
		}
	})
	return petui
}
