package page

type DirtyListNode struct {
	pg        *Page
	oldestCMN int32
}

func (n *DirtyListNode) GetOldestCMN() int32 {
	return n.oldestCMN
}

func (n *DirtyListNode) GetPage() *Page {
	return n.pg
}

func CreateDirtyListNode(pg *Page, oldestCMN int32) *DirtyListNode {
	return &DirtyListNode{pg, oldestCMN}
}
