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
"container/list"
"regexp"
"strings"
"sync"
)

type Element list.Element

type ISubscriber interface {
	OnSubscribe(elem *Element)
	OnPublish(channel string,v interface{})
	GetElement()*Element
}


type Subscriber struct {
	elem *Element
}
func(d *Subscriber)OnSubscribe(elem *Element){
	d.elem = elem
}
func(d *Subscriber)GetElement()*Element {
	return d.elem
}


type Server struct {
	exactTopics    map[string]*list.List
	patternTopics  map[string]*list.List
	topicsMutex    sync.Mutex
	done           chan interface{}
	unpublishedMsg chan *msgData
}

func(s *Server) start(){
	go func(){
		for ;; {
			select {
			case <-s.done:
				return
			case msg:=<-s.unpublishedMsg:
				s.dopublish(msg)
			}
		}
	}()
}

func(s *Server) Stop(){
	close(s.done)
}

func(s *Server) Subscribe(topic string, ps ISubscriber)bool{
	s.topicsMutex.Lock()
	defer s.topicsMutex.Unlock()

	if !standardizedTopic(&topic){
		return false
	}

	var topicSubsList *list.List

	if isPatternTopic(topic) == false{
		if s.exactTopics == nil{
			s.exactTopics = make(map[string]*list.List)
		}
		if _,ok := s.exactTopics[topic]; !ok{
			s.exactTopics[topic] = list.New()
		}
		topicSubsList = s.exactTopics[topic]
	}else{
		if s.patternTopics == nil{
			s.patternTopics = make(map[string]*list.List)
		}
		if _,ok := s.patternTopics[topic]; !ok{
			s.patternTopics[topic] = list.New()
		}
		topicSubsList = s.patternTopics[topic]
	}

	if topicSubsList != nil{
		ps.OnSubscribe((*Element)(topicSubsList.PushBack(ps)))
	}
	return true
}

func(s *Server) Unsubscribe(topic string, subs ISubscriber){
	s.topicsMutex.Lock()
	defer s.topicsMutex.Unlock()

	if !standardizedTopic(&topic){
		return
	}

	if s.exactTopics == nil{
		return
	}

	//Unsubscribe topic
	if subs == nil{
		if subsList,ok := s.exactTopics[topic]; ok{
			for subsList.Front()!= nil {
				subsList.Remove(subsList.Front())
			}
			delete(s.exactTopics, topic)
		}
		if subsList,ok := s.patternTopics[topic]; ok{
			for subsList.Front()!= nil {
				subsList.Remove(subsList.Front())
			}
			delete(s.patternTopics, topic)
		}
	}else{//Unsubscribe ISubscriber by topic
		if subsList,ok:= s.exactTopics[topic]; ok{
			subsList.Remove((*list.Element)(subs.GetElement()))
			if subsList.Len()==0{
				delete(s.exactTopics, topic)
			}
		}

		if subsList,ok:= s.patternTopics[topic]; ok{
			subsList.Remove((*list.Element)(subs.GetElement()))
			if subsList.Len()==0{
				delete(s.patternTopics, topic)
			}
		}
	}
}

func(s *Server) Publish(topic string,data interface{}){

	if !standardizedTopic(&topic){
		return
	}

	msg := newMsgData()
	msg.SetString("topic", topic)
	msg.SetValue("data",data)
	s.unpublishedMsg <-msg
}

func(s *Server) dopublish(msg *msgData){
	if msg == nil{
		return
	}

	topic := msg.GetString("topic")
	data := msg.GetValue("data")

	if !standardizedTopic(&topic){
		return
	}

	s.topicsMutex.Lock()
	defer s.topicsMutex.Unlock()


	//exact topics publish
	if subsList,ok:= s.exactTopics[topic]; ok {
		for item := subsList.Front(); item!= nil;{
			item.Value.(ISubscriber).OnPublish(topic,data)
			item = item.Next()
		}
	}

	//pattern topics publish
	for t,subsList := range s.patternTopics{
		regexPattern :=""
		//make regexPattern
		split := strings.Split(t,"*")
		for i:=0;i<len(split);i++{
			regexPattern += "("+ split[i]+")"+".*"
		}
		//match
		if matched,err:=regexp.MatchString(regexPattern, topic); err !=nil || matched == false{
			continue
		}

		for item := subsList.Front(); item!= nil;{
			item.Value.(ISubscriber).OnPublish(topic,data)
			item = item.Next()
		}
	}
}

func(s *Server) GetTopics()[]string{
	s.topicsMutex.Lock()
	defer s.topicsMutex.Unlock()

	i :=0
	topics := make([]string, len(s.exactTopics)+len(s.patternTopics))
	for topic,_ := range s.exactTopics {
		topics[i] = topic
		i++
	}
	for topic,_ := range s.patternTopics {
		topics[i] = topic
		i++
	}
	return topics
}

func(s *Server) GetSubscriberCount(topic string)int{
	s.topicsMutex.Lock()
	defer s.topicsMutex.Unlock()

	count:=0
	subsList,ok := s.exactTopics[topic]
	if ok {
		count += subsList.Len()
	}

	for t,subsList := range s.patternTopics{
		regexPattern :=""
		//make regexPattern
		split := strings.Split(t,"*")
		for i:=0;i<len(split);i++{
			regexPattern += "("+ split[i]+")"+".*"
		}
		//match
		if matched,err:=regexp.MatchString(regexPattern, topic); err !=nil || matched == false{
			continue
		}
		count += subsList.Len()
	}
	return count
}

func standardizedTopic(topic *string)bool{
	s := strings.TrimSpace(*topic)
	if s == "" || s == "*"{
		return false
	}
	topic = &s
	return true
}

func isPatternTopic(topic string) bool{
	if topic[len(topic)-1] != '*'{
		return false
	}
	return true
}

func NewServer(params ...int16)*Server {
	var cachedSize int16= 100
	for _, v := range params {
		cachedSize = v
		break
	}
	s := &Server{}
	s.done = make(chan interface{},0)
	s.unpublishedMsg = make(chan *msgData,cachedSize)
	s.start()
	return s
}