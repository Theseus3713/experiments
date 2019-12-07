package apt

import (
	"experiments/experiments/noise"
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

type Node interface {
	Eval(x, y float32) float32
	String() string
	AddRandom(node Node)
	NodeCounts() (nodeCount, nilCount int)
}

type BaseNode struct {
	Parent   Node
	Children []Node
}

type OpLerp struct {
	BaseNode
}

func (op *OpLerp) Eval(x, y float32) float32 {
	var (
		a   = op.Children[0].Eval(x, y)
		b   = op.Children[1].Eval(x, y)
		pct = op.Children[2].Eval(x, y)
	)
	return a + pct*(b-a)
}

func (op *OpLerp) String() string {
	return fmt.Sprint("( Lerp %s %s )", op.Children[1].String(), op.Children[2].String())
}

type OpClip struct {
	BaseNode
}

func (op *OpClip) Eval(x, y float32) float32 {
	var (
		value = op.Children[0].Eval(X, y)
		max   = float32(math.Abs(float64(op.Children[1].Eval(x, y))))
	)
	if value > max {
		return max
	} else if value < -max {
		return -max
	}
	return value
}

func (op *OpClip) String() string {
	return fmt.Sprint("( Clip %s %s )", op.Children[0].String(), op.Children[1].String())
}

//type LeafNode struct{}

//func (leaf *LeafNode) AddRandom(node Node) {
//	//panic(`ERROR: You tried to add a node to a leaf node`)
//	fmt.Println(` add a node to a leaf node`)
//}
//
//func (leaf *LeafNode) NodeCounts() (nodeCount, nilCount int) {
//	return 1, 0
//}

type SingleNode struct {
	Child Node
}

func (single *SingleNode) AddRandom(node Node) {
	if single.Child == nil {
		single.Child = node
	} else {
		single.Child.AddRandom(node)
	}
}

func (single *SingleNode) NodeCounts() (nodeCount, nilCount int) {
	if single.Child == nil {
		return 1, 1
	} else {
		childNodeCount, childNilCount := single.Child.NodeCounts()
		return 1 + childNodeCount, childNilCount
	}
}

type DoubleNode struct {
	LeftChild  Node
	RightChild Node
}

func (double *DoubleNode) AddRandom(node Node) {
	if rand.Intn(2) == 0 {
		if double.LeftChild == nil {
			double.LeftChild = node
		} else {
			double.LeftChild.AddRandom(node)
		}
	} else {
		if double.RightChild == nil {
			double.RightChild = node
		} else {
			double.RightChild.AddRandom(node)
		}
	}
}

func (double *DoubleNode) NodeCounts() (nodeCount, nilCount int) {
	var leftCount, leftNilCount, rightCount, rightNilCount int
	if double.LeftChild == nil {
		leftNilCount = 1
		leftCount = 0
	} else {
		leftCount, leftNilCount = double.LeftChild.NodeCounts()
	}

	if double.RightChild == nil {
		rightNilCount = 1
		rightCount = 0
	} else {
		rightCount, rightNilCount = double.RightChild.NodeCounts()
	}
	return 1 + leftCount + rightCount, leftNilCount + rightNilCount
}

type OpSin struct {
	BaseNode
}

func (op *OpSin) Eval(x, y float32) float32 {
	return float32(math.Sin(float64(op.Child.Eval(x, y))))
}

func (op *OpSin) String() string {
	return fmt.Sprintf(`( Sin %s)`, op.Child.String())
}

type OpCos struct {
	SingleNode
}

func (op *OpCos) Eval(x, y float32) float32 {
	return float32(math.Cos(float64(op.Child.Eval(x, y))))
}

func (op *OpCos) String() string {
	return fmt.Sprintf(`( Cos %s)`, op.Child.String())
}

type OpAtan struct {
	SingleNode
}

func (op *OpAtan) Eval(x, y float32) float32 {
	return float32(math.Atan(float64(op.Child.Eval(x, y))))
}

func (op *OpAtan) String() string {
	return fmt.Sprintf(`( Atan %s)`, op.Child.String())
}

type OpNoise struct {
	DoubleNode
}

func (op *OpNoise) Eval(x, y float32) float32 {
	return 80*noise.Snoise2(op.LeftChild.Eval(x, y), op.RightChild.Eval(x, y)) - 2.0
}

func (op *OpNoise) String() string {
	return fmt.Sprintf(`( SimplexNoise %s %s )`, op.LeftChild.String(), op.RightChild.String())
}

type OpPlus struct {
	DoubleNode
}

func (op *OpPlus) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) + op.RightChild.Eval(x, y)
}

func (op *OpPlus) String() string {
	return fmt.Sprintf(`( + %s %s )`, op.LeftChild.String(), op.RightChild.String())
}

type OpMinus struct {
	DoubleNode
}

func (op *OpMinus) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) - op.RightChild.Eval(x, y)
}

func (op *OpMinus) String() string {
	return fmt.Sprintf(`( - %s %s )`, op.LeftChild.String(), op.RightChild.String())
}

type OpMult struct {
	DoubleNode
}

func (op *OpMult) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) * op.RightChild.Eval(x, y)
}
func (op *OpMult) String() string {
	return fmt.Sprintf(`( * %s %s )`, op.LeftChild.String(), op.RightChild.String())
}

type OpDiv struct {
	DoubleNode
}

func (op *OpDiv) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) / op.RightChild.Eval(x, y)
}
func (op *OpDiv) String() string {
	return fmt.Sprintf(`( / %s %s )`, op.LeftChild.String(), op.RightChild.String())
}

type OpAtan2 struct {
	DoubleNode
}

func (op *OpAtan2) Eval(x, y float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}
func (op *OpAtan2) String() string {
	return fmt.Sprintf(`( OpAtan2 %s %s )`, op.LeftChild.String(), op.RightChild.String())
}

type OpX struct {
	LeafNode
}

func (op *OpX) Eval(x, y float32) float32 {
	return x
}

func (op *OpX) String() string {
	return "X"
}

type OpY struct {
	LeafNode
}

func (op *OpY) Eval(x, y float32) float32 {
	return y
}

func (op *OpY) String() string {
	return "Y"
}

type OpConstant struct {
	LeafNode
	value float32
}

func (op *OpConstant) Eval(x, y float32) float32 {
	return op.value
}

func (op *OpConstant) String() string {
	return strconv.FormatFloat(float64(op.value), 'f', 3, 32)
}

func GetRandomNode() Node {
	switch rand.Intn(9) {
	case 0:
		return &OpPlus{}
	case 1:
		return &OpMinus{}
	case 2:
		return &OpMult{}
	case 3:
		return &OpDiv{}
	case 4:
		return &OpAtan2{}
	case 5:
		return &OpAtan{}
	case 6:
		return &OpCos{}
	case 7:
		return &OpSin{}
	case 8:
		return &OpNoise{}
	}
	panic(`ERROR: Get Random Noise Failed`)
}

func GetRandomLeaf() Node {
	switch rand.Intn(3) {
	case 0:
		return &OpX{}
	case 1:
		return &OpY{}
	case 2:
		return &OpConstant{
			LeafNode: LeafNode{},
			value:    rand.Float32()*2 - 1,
		}
	}
	panic(`ERROR: Get Leaf`)
}
