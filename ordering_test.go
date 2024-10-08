package goiter

import (
    "cmp"
    "slices"
    "testing"
)

func TestOrder(t *testing.T) {
    // case 1: asc
    actual1 := make([]int, 0, 3)
    for v := range Order(SliceElems([]int{4, 1, 6})) {
        actual1 = append(actual1, v)
    }
    expect1 := []int{1, 4, 6}
    if !slices.Equal(expect1, actual1) {
        t.Fatal("expect:", expect1, "actual:", actual1)
    }

    // case 2: desc
    actual2 := make([]string, 0, 3)
    for v := range Order(MapVals(map[int]string{1: "bob", 2: "alice", 3: "eve"}), true) {
        actual2 = append(actual2, v)
    }
    expect2 := []string{"eve", "bob", "alice"}
    if !slices.Equal(expect2, actual2) {
        t.Fatal("expect:", expect2, "actual:", actual2)
    }

    // won't panic
    for _ = range Order(SliceElems([]int{1, 2, 3})) {
        break
    }
}

func TestOrder2V1(t *testing.T) {
    input := map[string]int{
        "bob":   20,
        "eve":   18,
        "alice": 22,
    }

    // case 1: asc
    actual1E1 := []string{}
    actual1E2 := []int{}
    for v1, v2 := range Order2V1(Map(input)) {
        actual1E1 = append(actual1E1, v1)
        actual1E2 = append(actual1E2, v2)
    }
    expect1E1 := []string{"alice", "bob", "eve"}
    if !slices.Equal(expect1E1, actual1E1) {
        t.Fatal("expect:", expect1E1, "actual:", actual1E1)
    }
    expect1E2 := []int{22, 20, 18}
    if !slices.Equal(expect1E2, actual1E2) {
        t.Fatal("expect:", expect1E2, "actual:", actual1E2)
    }

    // case 2: desc
    actual2E1 := []string{}
    actual2E2 := []int{}
    for v1, v2 := range Order2V1(Map(input), true) {
        actual2E1 = append(actual2E1, v1)
        actual2E2 = append(actual2E2, v2)
    }
    expect2E1 := []string{"eve", "bob", "alice"}
    if !slices.Equal(expect2E1, actual2E1) {
        t.Fatal("expect:", expect2E1, "actual:", actual2E1)
    }
    expect2E2 := []int{18, 20, 22}
    if !slices.Equal(expect2E2, actual2E2) {
        t.Fatal("expect:", expect2E2, "actual:", actual2E2)
    }

    // won't panic
    for _, _ = range Order2V1(Map(input)) {
        break
    }
}

func TestOrder2V2(t *testing.T) {
    input := map[string]int{
        "bob":   20,
        "eve":   18,
        "alice": 22,
    }

    // case 1: asc
    actual1E1 := []string{}
    actual1E2 := []int{}
    for v1, v2 := range Order2V2(Map(input)) {
        actual1E1 = append(actual1E1, v1)
        actual1E2 = append(actual1E2, v2)
    }
    expect1E1 := []string{"eve", "bob", "alice"}
    if !slices.Equal(expect1E1, actual1E1) {
        t.Fatal("expect:", expect1E1, "actual:", actual1E1)
    }
    expect1E2 := []int{18, 20, 22}
    if !slices.Equal(expect1E2, actual1E2) {
        t.Fatal("expect:", expect1E2, "actual:", actual1E2)
    }

    // case 2: desc
    actual2E1 := []string{}
    actual2E2 := []int{}
    for v1, v2 := range Order2V2(Map(input), true) {
        actual2E1 = append(actual2E1, v1)
        actual2E2 = append(actual2E2, v2)
    }
    expect2E1 := []string{"alice", "bob", "eve"}
    if !slices.Equal(expect2E1, actual2E1) {
        t.Fatal("expect:", expect2E1, "actual:", actual2E1)
    }
    expect2E2 := []int{22, 20, 18}
    if !slices.Equal(expect2E2, actual2E2) {
        t.Fatal("expect:", expect2E2, "actual:", actual2E2)
    }
}

func TestOrderBy(t *testing.T) {
    type person struct {
        name string
        age  int
    }
    input := []person{
        {"alice", 22},
        {"bob", 20},
        {"eve", 21},
    }
    actual := []person{}
    for each := range SliceElems(input).OrderBy(func(a, b person) int { return cmp.Compare(a.age, b.age) }) {
        actual = append(actual, each)
    }
    expect := []person{
        {"bob", 20},
        {"eve", 21},
        {"alice", 22},
    }
    if !slices.Equal(expect, actual) {
        t.Fatal("expect:", expect, "actual:", actual)
    }
}

func TestStableOrderBy(t *testing.T) {
    type person struct {
        name string
        age  int
    }
    input := []person{
        {"alice", 22},
        {"bob", 20},
        {"eve", 20},
    }
    actual := []person{}
    for each := range SliceElems(input).StableOrderBy(func(a, b person) int { return cmp.Compare(a.age, b.age) }) {
        actual = append(actual, each)
    }
    expect := []person{
        {"bob", 20},
        {"eve", 20},
        {"alice", 22},
    }
    if !slices.Equal(expect, actual) {
        t.Fatal("expect:", expect, "actual:", actual)
    }
}

func TestOrderBy2(t *testing.T) {
    type person struct {
        name string
        age  int
    }
    input := map[string]int{
        "bob":   20,
        "eve":   30,
        "alice": 25,
    }
    actual := []person{}
    for v1, v2 := range Map(input).OrderBy(func(a, b *Combined[string, int]) int { return cmp.Compare(a.V2, b.V2) }) {
        actual = append(actual, person{name: v1, age: v2})
    }
    expect := []person{
        {"bob", 20},
        {"alice", 25},
        {"eve", 30},
    }
    if !slices.Equal(expect, actual) {
        t.Fatal("expect:", expect, "actual:", actual)
    }
}

func TestStableOrderBy2(t *testing.T) {
    type person struct {
        name string
        age  int
    }
    input := []person{
        {"bob", 25},
        {"eve", 30},
        {"alice", 25},
    }
    actual := []person{}
    for _, v2 := range Slice(input).StableOrderBy(func(a, b *Combined[int, person]) int { return cmp.Compare(a.V2.age, b.V2.age) }) {
        actual = append(actual, v2)
    }
    expect := []person{
        {"bob", 25},
        {"alice", 25},
        {"eve", 30},
    }
    if !slices.Equal(expect, actual) {
        t.Fatal("expect:", expect, "actual:", actual)
    }
}
