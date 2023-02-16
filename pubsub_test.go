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

type User struct {
	Subscriber
	name string
}


func (u *User)OnPublish(channel string,data interface{})  {
	fmt.Println(u.name,":",data.(string))
}

func NewUser(name string)*User {
	return &User{name: name}
}

func TestAll(t *testing.T) {
	runner := NewRunner()
	userA := NewUser("Jim")
	userB := NewUser("Tom")
	userC := NewUser("Jack")

	runner.Subscribe("hello",userA)
	runner.Subscribe("hello*",userB)
	runner.Subscribe("helloworld",userC)
	runner.Subscribe("hello",userC)

	//[helloworld hello hello*]
	fmt.Println(runner.GetTopics())

	//Jim : hello
	//Jack : hello
	//Tom : hello
	runner.Publish("hello","hello")
	time.Sleep(time.Millisecond*100)

	//Tom : helloworld
	//Jack : helloworld
	runner.Publish("helloworld","helloworld")
	time.Sleep(time.Millisecond*100)

	runner.Unsubscribe("hello",userA)
	//Jack : hello world
	//Tom : hello world
	runner.Publish("hello","hello world")
	time.Sleep(time.Millisecond*100)

	//[hello* helloworld]
	runner.Unsubscribe("hello",nil)
	fmt.Println(runner.GetTopics())

	//Tom : hello world again
	runner.Publish("hello","hello world again")
	time.Sleep(time.Millisecond*100)

	//[helloworld]
	runner.Unsubscribe("hello*",nil)
	fmt.Println(runner.GetTopics())
	time.Sleep(time.Millisecond*100)

}