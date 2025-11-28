package goiter

import (
	"iter"
)

type Seq2X[T1, T2 any] interface {
	~func(yield func(T1, T2) bool)
}

type Iterator2[T1, T2 any] func(yield func(T1, T2) bool)

func (it Iterator2[T1, T2]) Seq() iter.Seq2[T1, T2] {
	return iter.Seq2[T1, T2](it)
}

// PickV1 returns an iterator that yields the first value of each pair.
func (it Iterator2[T1, T2]) PickV1() Iterator[T1] {
	return PickV1(it)
}

// PickV2 returns an iterator that yields the second value of each pair.
func (it Iterator2[T1, T2]) PickV2() Iterator[T2] {
	return PickV2(it)
}

func (it Iterator2[T1, T2]) OrderBy(cmp func(*Combined[T1, T2], *Combined[T1, T2]) int) Iterator2[T1, T2] {
	return Order2By(it, cmp)
}

func (it Iterator2[T1, T2]) StableOrderBy(cmp func(*Combined[T1, T2], *Combined[T1, T2]) int) Iterator2[T1, T2] {
	return StableOrder2By(it, cmp)
}

func (it Iterator2[T1, T2]) Filter(cmp func(T1, T2) bool) Iterator2[T1, T2] {
	return Filter2(it, cmp)
}

func (it Iterator2[T1, T2]) Take(n int) Iterator2[T1, T2] {
	return Take2(it, n)
}

func (it Iterator2[T1, T2]) TakeLast(n int) Iterator2[T1, T2] {
	return TakeLast2(it, n)
}

func (it Iterator2[T1, T2]) Skip(n int) Iterator2[T1, T2] {
	return Skip2(it, n)
}

func (it Iterator2[T1, T2]) SkipLast(n int) Iterator2[T1, T2] {
	return SkipLast2(it, n)
}

func (it Iterator2[T1, T2]) Concat(its ...Iterator2[T1, T2]) Iterator2[T1, T2] {
	return Concat2(it, its...)
}

func (it Iterator2[T1, T2]) Reverse() Iterator2[T1, T2] {
	return Reverse2(it)
}

func (it Iterator2[T1, T2]) Count() int {
	return Count2(it)
}

func (it Iterator2[T1, T2]) Through(f func(T1, T2) (T1, T2)) Iterator2[T1, T2] {
	return Transform2(it, f)
}

func (it Iterator2[T1, T2]) Cache() Iterator2[T1, T2] {
	return Cache2(it)
}

func (it Iterator2[T1, T2]) Once() Iterator2[T1, T2] {
	return Once2(it)
}

func (it Iterator2[T1, T2]) FinishOnce() Iterator2[T1, T2] {
	return FinishOnce2(it)
}
