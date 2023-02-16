// License statement
// ---------------------------------------------------------------------------------
//
// Licensed under The MIT License
// For full copyright and license information, please see the MIT-LICENSE.txt
// Redistributions of files must retain the above copyright notice.
//
// @author    shitou<rookiefun@sina.com>
// @copyright shitou<rookiefun@sina.com>
// @license   http://www.opensource.org/licenses/mit-license.php MIT License
// ---------------------------------------------------------------------------------

package pubsub

import (
	"fmt"
	"testing"
	"time"
)

type SubscriberEntity struct {
	Subscriber
	name string
}


func (u *SubscriberEntity)OnPublish(channel string,data interface{})  {
	fmt.Println(u.name,":",data.(string))
}

func newSubscriberEntity(name string)*SubscriberEntity {
	return &SubscriberEntity{name: name}
}

func TestAll(t *testing.T) {
	server := NewServer()
	a := newSubscriberEntity("Jim")
	b := newSubscriberEntity("Tom")
	c := newSubscriberEntity("Jack")

	server.Subscribe("hello", a)
	server.Subscribe("hello*", b)
	server.Subscribe("helloworld", c)
	server.Subscribe("hello", c)

	//[helloworld hello hello*]
	fmt.Println(server.GetTopics())

	//Jim : hello
	//Jack : hello
	//Tom : hello
	server.Publish("hello","hello")
	time.Sleep(time.Millisecond*100)

	//Tom : helloworld
	//Jack : helloworld
	server.Publish("helloworld","helloworld")
	time.Sleep(time.Millisecond*100)

	server.Unsubscribe("hello", a)
	//Jack : hello world
	//Tom : hello world
	server.Publish("hello","hello world")
	time.Sleep(time.Millisecond*100)

	//[hello* helloworld]
	server.Unsubscribe("hello",nil)
	fmt.Println(server.GetTopics())

	//Tom : hello world again
	server.Publish("hello","hello world again")
	time.Sleep(time.Millisecond*100)

	//[helloworld]
	server.Unsubscribe("hello*",nil)
	fmt.Println(server.GetTopics())
	time.Sleep(time.Millisecond*100)

}