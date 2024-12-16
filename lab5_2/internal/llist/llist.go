package llist

import (
	"fmt"
	"strings"
	"sync"
)

type Node struct {
	value int
	next  *Node
}

type LinkedList struct {
	head *Node
	mu   sync.RWMutex
}

func New() *LinkedList {
	return &LinkedList{}
}

func (list *LinkedList) AddIfNotExists(value int) {
	list.mu.RLock()

	current := list.head
	for current != nil {
		if current.value == value {
			list.mu.RUnlock()

			return
		}
		current = current.next
	}

	list.mu.RUnlock()

	list.mu.Lock()
	defer list.mu.Unlock()

	current = list.head
	for current != nil {
		if current.value == value {
			return
		}

		current = current.next
	}

	newNode := &Node{value: value}
	if list.head == nil {
		list.head = newNode

		return
	}

	current = list.head
	for current.next != nil {
		current = current.next
	}

	current.next = newNode
}

func (list *LinkedList) CheckForDuplicates() bool {
	list.mu.RLock()
	defer list.mu.RUnlock()

	seen := make(map[int]bool)

	current := list.head
	for current != nil {
		if seen[current.value] {
			return true
		}

		seen[current.value] = true
		current = current.next
	}

	return false
}

func (list *LinkedList) String() string {
	list.mu.RLock()
	defer list.mu.RUnlock()

	var sb strings.Builder

	current := list.head
	for ; current != nil; current = current.next {
		sb.WriteString(fmt.Sprintf("%d ", current.value))
	}

	return sb.String()
}
