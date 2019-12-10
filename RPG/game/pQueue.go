package game

type priorityPosition struct {
	Position
	priority int
}

type pQueue []priorityPosition

func (pq pQueue) push(position Position, priority int) pQueue {
	var newNode = priorityPosition{position, priority}
	pq = append(pq, newNode)
	newNodeIndex := len(pq) - 1
	parentIndex, parent := pq.parent(newNodeIndex)

	for newNode.priority < parent.priority && newNodeIndex != 0 {
		pq.swap(newNodeIndex, parentIndex)
		newNodeIndex = parentIndex
		parentIndex, parent = pq.parent(newNodeIndex)
	}
	return pq
}

func (pq pQueue) pop() (pQueue, Position) {
	var result = pq[0].Position
	pq[0] = pq[len(pq)-1]
	pq = pq[:len(pq)-1]

	if len(pq) == 0 {
		return pq, result
	}
	var (
		index                          = 0
		node                           = pq[index]
		leftExists, leftIndex, left    = pq.left(index)
		rightExists, rightIndex, right = pq.right(index)
	)
	for (leftExists && node.priority > left.priority) ||
		(rightExists && node.priority > right.priority) {

		if !rightExists || left.priority <= right.priority {
			pq.swap(index, leftIndex)
			index = leftIndex
		} else {
			pq.swap(index, rightIndex)
			index = rightIndex
		}

		leftExists, leftIndex, left = pq.left(index)
		rightExists, rightIndex, right = pq.right(index)

	}
	return pq, result

}

func (pq pQueue) swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq pQueue) parent(i int) (int, priorityPosition) {
	var index = (i - 1) / 2
	return index, pq[i]
}

func (pq pQueue) left(i int) (bool, int, priorityPosition) {
	var index = i*2 + 1
	if index < len(pq) {
		return true, index, pq[index]
	}
	return false, 0, priorityPosition{}
}

func (pq pQueue) right(i int) (bool, int, priorityPosition) {
	var index = i*2 + 2
	if index < len(pq) {
		return true, index, pq[index]
	}
	return false, 0, priorityPosition{}
}
