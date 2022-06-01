package main

import (
	"sync/atomic"
	"unsafe"
)

type Node[T any] struct {
	Value T
	Next  *Node[T]
	Prev  *Node[T]
}

type Queue[T any] struct {
	Tail *Node[T]
	Head *Node[T]
}

func CompareAndSwapPointer[T any](addr **T, old, new *T) (swapped bool) {
	return atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(addr)),
		unsafe.Pointer(old),
		unsafe.Pointer(new),
	)
}

func LoadPointer[T any](addr **T) (val *T) {
	return (*T)(
		atomic.LoadPointer(
			(*unsafe.Pointer)(unsafe.Pointer(addr)),
		))
}

func StorePointer[T any](addr **T, val *T) {
	atomic.StorePointer(
		(*unsafe.Pointer)(unsafe.Pointer(addr)),
		unsafe.Pointer(val),
	)
}

func MakeQueue[T any]() Queue[T] {
	q := Queue[T]{}
	nd := &Node[T]{}
	q.Tail = nd
	q.Head = nd
	return q
}

func (q *Queue[T]) Push(val T) {
	var tail *Node[T]
	nd := &Node[T]{}
	nd.Value = val
	for {
		tail = LoadPointer(&q.Tail)
		nd.Next = tail
		if CompareAndSwapPointer(&q.Tail, tail, nd) {
			StorePointer(&tail.Prev, nd)
			break
		}
	}
}

func (q *Queue[T]) Pop() (T, bool) {
	var tail, head, firstNodePrev *Node[T]
	var val T
	for {
		head = LoadPointer(&q.Head)
		tail = LoadPointer(&q.Tail)
		firstNodePrev = LoadPointer(&head.Prev)
		if head == LoadPointer(&q.Head) {
			if tail != head {
				if firstNodePrev == nil || firstNodePrev.Next != head {
					q.fixList(tail, head)
					continue
				}
				val = firstNodePrev.Value
				if CompareAndSwapPointer(&q.Head, head, firstNodePrev) {
					return val, true
				}
			} else {
				return val, false
			}
		}
	}
}

func (q *Queue[T]) fixList(tail *Node[T], head *Node[T]) {
	var curNode, curNodeNext *Node[T]
	curNode = tail
	for head == LoadPointer(&q.Head) && curNode != head {
		curNodeNext = curNode.Next
		StorePointer(&curNodeNext.Prev, curNode)
		curNode = curNodeNext
	}
}
