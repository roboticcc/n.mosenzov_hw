package hw04lrucache

type ListItem struct {
	Value any
	Next  *ListItem
	Prev  *ListItem
}

type List struct {
	front *ListItem
	back  *ListItem
	len   int
}

func NewList() *List {
	return &List{}
}

func (l *List) Len() int {
	return l.len
}

func (l *List) Front() *ListItem {
	return l.front
}

func (l *List) Back() *ListItem {
	return l.back
}

func (l *List) PushFront(v any) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}
	if l.front != nil {
		l.front.Prev = newItem
	} else {
		l.back = newItem
	}
	l.front = newItem
	l.len++
	return newItem
}

func (l *List) PushBack(v any) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}
	if l.back != nil {
		l.back.Next = newItem
	} else {
		l.front = newItem
	}
	l.back = newItem
	l.len++
	return newItem
}

func (l *List) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	l.len--
}

func (l *List) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}
