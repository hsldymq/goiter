//go:build goexperiment.rangefunc

package goiter

import (
    "iter"
    "math"
    "time"
)

type GeneratorFunc[T any] func() (T, bool)

type GeneratorFunc2[T1, T2 any] func() (T1, T2, bool)

type TInt interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func Range[T TInt](start, stop T) iter.Seq[T] {
    return RangeStep(start, stop, 1)
}

// RangeStep extends the ability to Range function, allowing iteration from any integer and stepping forward or backward in any step.
// It is similar to Python's range function, but with some differences:
//  1. stepSize does not accept negative numbers. Whether iterating forward or backward, stepSize must be positive.
//     so you don't need to consider adjusting the sign of step according to the direction of iteration, It is the absolute value of the step parameter of Python range function.
//  2. Providing a value less than or equal to 0 for stepSize will not return an error, it simply doesn't yield any values.
func RangeStep[T TInt, S TInt](start, stop T, stepSize S) iter.Seq[T] {
    if stepSize <= 0 {
        // 0 will lead to infinite loops
        return Empty[T]()
    }

    step := uint64(stepSize)
    inc := true
    if start > stop {
        inc = false
    }

    if willOverflow(start, step, inc) {
        return func(yield func(T) bool) {
            yield(start)
        }
    }

    return func(yield func(T) bool) {
        curr := start
        for {
            if !yield(curr) {
                return
            }

            if inc {
                next := curr + T(step)
                if next >= stop || next < start || next <= curr {
                    return
                }
                curr = next
            } else {
                next := curr - T(step)
                if next <= stop || next > start || next >= curr {
                    return
                }
                curr = next
            }
        }
    }
}

// RangeTime is similar to RangeStep, but it is specifically used for iterating over time, and it can iterate time forward or backward.
// The interval parameter is its step size, which can be any positive duration.
// Unlike the half-open interval represented by the start and end parameters of RangeStep, the from and to parameters of RangeTime represent a closed interval.
// For example:
//
//	from := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
//	to := time.Date(2021, 1, 1, 5, 0, 0, 0, time.UTC)
//	for t := range RangeTime(from, to, time.Hour) {
//			fmt.Println(t.Format("15:04:05"))
//	}
//	the above code will print:
//		00:00:00
//		01:00:00
//		02:00:00
//		03:00:00
//		04:00:00
//		05:00:00
func RangeTime(from time.Time, to time.Time, interval time.Duration) iter.Seq[time.Time] {
    if interval <= 0 {
        return Empty[time.Time]()
    }

    return func(yield func(time.Time) bool) {
        if from.Before(to) || from.Equal(to) {
            t := from
            for t.Before(to) || t.Equal(to) {
                if !yield(t) {
                    return
                }
                t = t.Add(interval)
            }
        } else {
            t := from
            for t.After(to) || t.Equal(to) {
                if !yield(t) {
                    return
                }
                t = t.Add(-interval)
            }
        }
    }
}

// Sequence takes a generator function and returns an iterator that yields the values generated by the generator.
// This is a general sequence generator function
// For example, you can use it to generate the Fibonacci sequence like this:
//
//	genFib := func(n int) goiter.GeneratorFunc[int] {
//		a, b := 0, 1
//		return func() (int, bool) {
//			if n <= 0 {
//				return 0, false
//			}
//			n--
//			a, b = b, a+b
//			return a, true
//		}
//	}
//	// this will print first 5 Fibonacci numbers: 1 1 2 3 5
//	for v := range goiter.Sequence(genFib(5)) {
//		fmt.Printf("%d ", v)
//	}
func Sequence[T any](generator func() (T, bool)) iter.Seq[T] {
    return func(yield func(T) bool) {
        for {
            v, hasMore := generator()
            if !hasMore {
                return
            }
            if !yield(v) {
                return
            }
        }
    }
}

// Sequence2 is the iter.Seq2 version of Sequence function.
func Sequence2[T1, T2 any](generator func() (T1, T2, bool)) iter.Seq2[T1, T2] {
    return func(yield func(T1, T2) bool) {
        for {
            v1, v2, hasMore := generator()
            if !hasMore {
                return
            }
            if !yield(v1, v2) {
                return
            }
        }
    }
}

// Concat returns an iterator that allows you to traverse multiple iterators in sequence.
func Concat[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
    return func(yield func(T) bool) {
        for _, seq := range seqs {
            for v := range seq {
                if !yield(v) {
                    return
                }
            }
        }
    }
}

// Concat2 returns an iterator that allows you to traverse multiple iterators in sequence.
func Concat2[T1 any, T2 any](seqs ...iter.Seq2[T1, T2]) iter.Seq2[T1, T2] {
    return func(yield func(T1, T2) bool) {
        for _, seq := range seqs {
            for v1, v2 := range seq {
                if !yield(v1, v2) {
                    return
                }
            }
        }
    }
}

func Reverse[T any](seq iter.Seq[T]) iter.Seq[T] {
    return func(yield func(T) bool) {
        var buffer []T
        next, stop := iter.Pull(seq)
        defer stop()
        for {
            v, ok := next()
            if !ok {
                break
            }
            buffer = append(buffer, v)
        }
        for i := len(buffer) - 1; i >= 0; i-- {
            if !yield(buffer[i]) {
                return
            }
        }
    }
}

func Reverse2[T1, T2 any](seq iter.Seq2[T1, T2]) iter.Seq2[T1, T2] {
    return func(yield func(T1, T2) bool) {
        var buffer []*Combined[T1, T2]
        next, stop := iter.Pull2(seq)
        defer stop()
        for {
            v1, v2, ok := next()
            if !ok {
                break
            }
            buffer = append(buffer, &Combined[T1, T2]{V1: v1, V2: v2})
        }
        for i := len(buffer) - 1; i >= 0; i-- {
            if !yield(buffer[i].V1, buffer[i].V2) {
                return
            }
        }
    }
}

func willOverflow[T TInt](v T, step uint64, inc bool) bool {
    tMax := int64(intMax(v))
    tMin := int64(intMin(v))

    if tMax != math.MaxInt64 && step >= uint64(tMax-tMin+1) {
        return true
    }
    if inc && v+T(step) < v {
        return true
    }
    if !inc && v-T(step) > v {
        return true
    }

    return false
}

func intMin[T TInt](v T) T {
    ones := ^T(0)
    if ones < 0 {
        return ^(ones ^ (1 << (countBits(ones) - 1)))
    }
    return 0
}

func intMax[T TInt](v T) T {
    ones := ^T(0)
    if ones < 0 {
        return ones ^ (1 << (countBits(ones) - 1))
    }
    return ones
}

func countBits[T TInt](v T) int {
    v = 1
    for _, bits := range [4]int{8, 16, 32} {
        if v<<bits == 0 {
            return bits
        }
    }
    return 64
}
