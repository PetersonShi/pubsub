# pubsub

### 简介
一个用go实现的高效订阅发布模块。只需要简单的几行代码，就可以实现一个订阅-发布服务。

### 使用
> 获取pubsub
```go
go get github.com/PetersonShi/pubsub@v0.1.0
```
> 在项目中引用刚才下载的包，使用代码如下
```go

   package main
   
   import (
   	"fmt"
   	"github.com/PetersonShi/pubsub"
   	"time"
   )
   
   //订阅者需要在结构中使用pubsub.Subscriber作为匿名成员
   //然后需要实现func OnPublish(channel string,data interface{})函数
   //that's all
   type SubscriberEntity struct {
   	pubsub.Subscriber
   	name string
   }
   
   //实现消息接收函数
   func (u *SubscriberEntity)OnPublish(channel string,data interface{})  {
   	    //为所欲为
        fmt.Println(u.name,":",data.(string))
   }
   
   
   func newSubscriberEntity(name string)*SubscriberEntity {
   	return &SubscriberEntity{name: name}
   }
   
   
   func main(){
   	//创建并开启一个订阅服务
   	server := pubsub.NewRunner()
   
   	//创建三个订阅者
   	a := newSubscriberEntity("Jim")
   	b := newSubscriberEntity("Tom")
   	c := newSubscriberEntity("Jack")
   
   	//订阅主题-精准
   	server.Subscribe("hello", a)
   	server.Subscribe("helloworld", c)
   	server.Subscribe("hello", c)
   	//订阅主题-模式
   	server.Subscribe("hello*", b)
   
   	//获取所有已订阅的主题
   	//[helloworld hello hello*]
   	fmt.Println(server.GetTopics())
   
   	//向hello主题发送消息“hello”
   	//Jim : hello
   	//Jack : hello
   	//Tom : hello
   	server.Publish("hello","hello")
   	time.Sleep(time.Millisecond*100)
   
   	//向helloworld主题发送消息“helloworld”
   	//Tom : helloworld
   	//Jack : helloworld
   	server.Publish("helloworld","helloworld")
   	time.Sleep(time.Millisecond*100)
   
   	//把a从hello主题中取消订阅
   	server.Unsubscribe("hello", a)
   
   	//向hello主题发送消息“hello world”
   	//Jack : hello world
   	//Tom : hello world
   	server.Publish("hello","hello world")
   	time.Sleep(time.Millisecond*100)
   
   	//取消hello主题所有订阅
   	//[hello* helloworld]
   	server.Unsubscribe("hello",nil)
   	fmt.Println(server.GetTopics())
   
   	//向hello主题发送消息“hello world again”
   	//Tom : hello world again
   	server.Publish("hello","hello world again")
   	time.Sleep(time.Millisecond*100)
   
   	//取消hello*主题所有订阅
   	server.Unsubscribe("hello*",nil)
   
   	//获取所有已订阅的主题
   	//[helloworld]
   	fmt.Println(server.GetTopics())
   	time.Sleep(time.Millisecond*100)
   }
```
### 注意
+ 可以根据实际需求创建多个server;
+ 消息发布是一个异步操作，其它操作如订阅，取消订阅，获取订阅数等都是同步操作。所以如果消息发布后，何时会被处理依赖系统处理能力;
+ 关于模式主题，目前只实现topic*这种模式，后续会继续完善；