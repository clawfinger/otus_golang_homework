package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	Head *ListItem
	Tail *ListItem

	len int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.Head
}

func (l *list) Back() *ListItem {
	return l.Tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.len == 0 {
		l.Tail = item
	} else {
		head := l.Head
		item.Next = head
		head.Prev = item
	}
	l.Head = item

	l.len++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}
	if l.len == 0 {
		l.Head = item
	} else {
		tail := l.Tail
		tail.Next = item
		item.Prev = tail
	}
	l.Tail = item
	l.len++
	return item
}

func (l *list) Remove(i *ListItem) {
	prev := i.Prev

	if prev == nil {
		l.Head = i.Next
	}

	next := i.Next

	if next == nil {
		l.Tail = i.Prev
	}

	prev.Next = next
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.Head {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
