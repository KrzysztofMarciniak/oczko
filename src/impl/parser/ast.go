package parser

type Node interface{ isNode() }
type Atom struct{ Value string }

func (Atom) isNode() {}

type List struct{ Items []Node }

func (List) isNode() {}
