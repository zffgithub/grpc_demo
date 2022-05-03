package main

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
)

type Queue struct {
	topic string
	list  *list.List
	mutex sync.Mutex
}

var multi_queue map[string]*Queue

func init() {
	multi_queue = make(map[string]*Queue)
}

func GetQueue(topic string) *Queue {
	if _, ok := multi_queue[topic]; ok {
		fmt.Println("Exist Queue Len: ", multi_queue[topic].Len())
		return multi_queue[topic]
	} else {
		multi_queue[topic] = &Queue{
			topic: topic,
			list:  list.New(),
		}
		fmt.Println("Make New Queue Len: ", multi_queue[topic].Len())
		return multi_queue[topic]
	}
}

func (queue *Queue) Push(data interface{}) {
	if data == nil {
		return
	}
	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	queue.list.PushBack(data)
	queue.Show()
}

func (queue *Queue) Pop() (interface{}, error) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	if element := queue.list.Front(); element != nil {
		queue.list.Remove(element)
		return element.Value, nil
	} else {
		fmt.Println("Pop Error:")
	}
	return nil, errors.New("pop failed")
}

func (queue *Queue) Clear() {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	for element := queue.list.Front(); element != nil; {
		elementNext := element.Next()
		queue.list.Remove(element)
		element = elementNext
	}
}

func (queue *Queue) Len() int {
	// fmt.Println(time.Now())
	queue.Show()
	return queue.list.Len()
}

func (queue *Queue) Show() {
	for item := queue.list.Front(); item != nil; item = item.Next() {
		fmt.Println("show value:", item.Value)
	}
}
