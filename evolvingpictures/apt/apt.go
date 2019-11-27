package apt

import (
	"fmt"
	"math"
)

type Node interface {
	Eval(x, y float32) float32
	String() string
}

type LeafNode struct {}

type SingleNode struct {
	Child Node
}

type DoubleNode struct {
	LeftChild Node
	RightChild Node
}

type OpSin SingleNode
func (op *OpSin) Eval(x, y float32) float32 {
	return float32(math.Sin(float64(op.Child.Eval(x, y))))
}

func (op *OpSin) String() string {
	return fmt.Sprintf(`( Sin %s)`, op.Child.String())
}

type OpPlus DoubleNode
func (op *OpPlus) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) + op.RightChild.Eval(x, y)
}
func (op *OpPlus) String() string {
	return fmt.Sprintf(`( + %s %s )`, op.LeftChild.String(), op.RightChild.String())
}

type OpMinus DoubleNode
func (op *OpMinus) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) - op.RightChild.Eval(x, y)
}
func (op *OpMinus) String() string {
	return fmt.Sprintf(`( - %s %s )`, op.LeftChild.String(), op.RightChild.String())
}

type OpMult DoubleNode
func (op *OpMult) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) * op.RightChild.Eval(x, y)
}
func (op *OpMult) String() string {
	return fmt.Sprintf(`( * %s %s )`, op.LeftChild.String(), op.RightChild.String())
}

type OpDiv DoubleNode
func (op *OpDiv) Eval(x, y float32) float32 {
	return op.LeftChild.Eval(x, y) / op.RightChild.Eval(x, y)
}
func (op *OpDiv) String() string {
	return fmt.Sprintf(`( / %s %s )`, op.LeftChild.String(), op.RightChild.String())
}

type OpAtan2 DoubleNode
func (op *OpAtan2) Eval(x, y float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}
func (op *OpAtan2) String() string {
	return fmt.Sprintf(`( OpAtan2 %s %s )`, op.LeftChild.String(), op.RightChild.String())
}

type OpX  LeafNode
func (op *OpX) Eval(x, y float32) float32 {
	return x
}

func (op *OpX) String() string {
	return "X"
}

type OpY LeafNode
func (op *OpY) Eval(x, y float32) float32 {
	return y
}

func (op *OpY) String() string {
	return "Y"
}