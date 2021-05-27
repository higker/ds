// Open Source: MIT License
// Author: Jaco Ding <deen.job@qq.com>
// Date: 2021/5/26 - 10:52 下午 - UTC/GMT+08:00

// 自己设计的map 超大型单机并发

package kmap

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

// 为了快速查找建立外部索引 k:1,34 就能快速查找到位置
var _index map[interface{}][2]int

type KMap interface {
	Put(k interface{}, v interface{}) bool
	Get(k interface{}) interface{}
	Debug()
}

type Root struct {
	lastIndex int
	data      []*MapItem
	size      int // 这个确定的
}

type MapItem struct {
	k, v interface{}
}

type Map struct {
	capacity int
	size     int
	index    []*Root
}

// 1. for 初始化
// 2.
func New() KMap {
	m := new(Map)
	m.index = make([]*Root, 10, 10)
	// 初始化索引
	for i := range m.index {
		root := new(Root)
		mapItems := make([]*MapItem, 100, 100)
		root.data = mapItems
		root.size = 0
		m.index[i] = root
	}
	m.size = cap(m.index)
	return m
}

func (m *Map) Hash(key interface{}) int {
	var code int = -1
	switch key.(type) {
	case string:
		code = _stringToCode(key.(string))
	case int, int64:
		// 使用crypto/rand生成随机数 然后 计算哈希
		code = _randomInt(100)
	}
	return code
}

// 通过哈希计算 得到root节点下标
func (m *Map) Index(k interface{}) int {
	return m.Hash(k) % m.size
}

func (m *Map) Put(k interface{}, v interface{}) bool {
	// 已经存在
	if _, ok := _index[k]; ok {
		return false
	}

	// 拿到所在的组，满了重新做一次记录
	bucketIndex := m.Index(k)
	root := m.index[bucketIndex]
	if root.lastIndex == root.size {
		// 容量已经满了
	}

	// 通过尾部指针找到数组当前在哪个位置是空的，把元素插入
	root.data[root.lastIndex] = &MapItem{k: k, v: v}
	// 更新外部索引
	_index[k] = [2]int{bucketIndex, root.lastIndex}
	root.lastIndex++

	return true
}

func (m *Map) Debug() {
	fmt.Println(m.index[9].data[1])
	fmt.Println(m.index[8].data[2])
	fmt.Println(m.index[7].data[3])
	fmt.Println(m.index[6].data[4])
}

func (m *Map) Get(k interface{}) interface{} {
	root := m.index[m.Index(k)]
	for _, ele := range root.data {
		if ele.k == k {
			return ele.v
		}
	}
	return nil
}

func (m *Map) Remove(k interface{}) {
	root := m.index[m.Index(k)]
	for _, ele := range root.data {
		if ele.k == k {
			ele = nil
			return
		}
	}
}

func _stringToCode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func _randomInt(max int) int {
	var n uint16
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return int(n) % max
}
