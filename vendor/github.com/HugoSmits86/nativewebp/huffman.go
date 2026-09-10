package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "container/heap"
    "sort"
)

const (
    NUM_HUFFMAN_BITS        = 3
    MIN_HUFFMAN_BITS        = 2
    MAX_HUFFMAN_BITS        = (MIN_HUFFMAN_BITS + (1 << NUM_HUFFMAN_BITS) - 1)
    MAX_HUFF_IMAGE_SIZE     = 2600
)

type huffmanCode struct {
    Symbol  int
    Bits    int
    Depth   int
}

type node struct {
    IsBranch    bool
    Weight      int
    Symbol      int
    BranchLeft  *node
    BranchRight *node
}

type nodeHeap []*node
func (h nodeHeap) Len() int             { return len(h) }
func (h nodeHeap) Less(i, j int) bool   { return h[i].Weight < h[j].Weight }
func (h nodeHeap) Swap(i, j int)        { h[i], h[j] = h[j], h[i] }
func (h *nodeHeap) Push(x interface{})  { *h = append(*h, x.(*node)) }
func (h *nodeHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

func buildHuffmanTree(histo []int, maxDepth int) *node {
    sum := 0
    for _, x := range histo {
        sum += x
    }

    minWeight := sum >> (maxDepth - 2)

    nHeap := &nodeHeap{}
    heap.Init(nHeap)

    for s, w := range histo {
        if w > 0 {
            if w < minWeight {
                w = minWeight
            }

            heap.Push(nHeap, &node{
                Weight: w, 
                Symbol: s,
            })
        }
    }
    
    for nHeap.Len() < 1 {
        heap.Push(nHeap, &node{
            Weight: minWeight, 
            Symbol: 0,
        })
    }
    
    for nHeap.Len() > 1 {
        n1 := heap.Pop(nHeap).(*node)
        n2 := heap.Pop(nHeap).(*node)
        heap.Push(nHeap, &node{
            IsBranch: true, 
            Weight: n1.Weight + n2.Weight, 
            BranchLeft: n1, 
            BranchRight: n2,
        })
    }

    return heap.Pop(nHeap).(*node)
}

func buildhuffmanCodes(histo []int, maxDepth int) []huffmanCode {
    codes := make([]huffmanCode, len(histo))

    tree := buildHuffmanTree(histo, maxDepth)
    if !tree.IsBranch {
        codes[tree.Symbol] = huffmanCode{tree.Symbol, 0, -1}
        return codes
    }
    
    var symbols []huffmanCode
    setBitDepths(tree, &symbols, 0)

    sort.Slice(symbols, func(i, j int) bool {
        if symbols[i].Depth == symbols[j].Depth {
            return symbols[i].Symbol < symbols[j].Symbol
        }

        return symbols[i].Depth < symbols[j].Depth
    })

    bits := 0
    prevDepth := 0
    for _, sym := range symbols {
        bits <<= (sym.Depth - prevDepth)
        codes[sym.Symbol].Symbol = sym.Symbol
        codes[sym.Symbol].Bits = bits
        codes[sym.Symbol].Depth = sym.Depth
        bits++

        prevDepth = sym.Depth
    }

    return codes
}

func setBitDepths(node *node, codes *[]huffmanCode, level int) {
    if node == nil {
        return
    }

    if !node.IsBranch {
        *codes = append(*codes, huffmanCode{
            Symbol: node.Symbol,
            Depth: level,
        })

        return
    }

    setBitDepths(node.BranchLeft, codes, level + 1)
    setBitDepths(node.BranchRight, codes, level + 1)
}

func writehuffmanCodes(w *bitWriter, codes []huffmanCode) {
    var symbols [2]int
    
    cnt := 0
    for _, code := range codes {
        if code.Depth != 0 {
            if cnt < 2 {
                symbols[cnt] = code.Symbol
            }

            cnt++
        }

        if cnt > 2 {
            break
        }
    }
    
    if cnt == 0 {
        w.writeBits(1, 1)
        w.writeBits(0, 3)
    } else if cnt <= 2 && symbols[0] < 1 << 8 && symbols[1] < 1 << 8 {
        w.writeBits(1, 1)
        w.writeBits(uint64(cnt - 1), 1)
        if symbols[0] <= 1 {
            w.writeBits(0, 1)
            w.writeBits(uint64(symbols[0]), 1)
        } else {
            w.writeBits(1, 1)
            w.writeBits(uint64(symbols[0]), 8)
        }

        if cnt > 1 {
            w.writeBits(uint64(symbols[1]), 8)
        }
    } else {
        writeFullhuffmanCode(w, codes)
    }
}

func writeFullhuffmanCode(w *bitWriter, codes []huffmanCode) {
    histo := make([]int, 19)
    for _, c := range codes {
        histo[c.Depth]++
    }

    // lengthCodeOrder comes directly from the WebP specs!
    var lengthCodeOrder = []int{
        17, 18, 0, 1, 2, 3, 4, 5, 16, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
    }

    cnt := 0
    for i, c := range lengthCodeOrder {
        if histo[c] > 0 {
            cnt = max(i + 1, 4)
        }
    }

    w.writeBits(0, 1)
    w.writeBits(uint64(cnt - 4), 4)

    lengths := buildhuffmanCodes(histo, 7)
    for i := 0; i < cnt; i++ {
        l := lengths[lengthCodeOrder[i]].Depth
        if l < 0 {
            w.writeBits(uint64(1), 3)
            continue
        }
        
        w.writeBits(uint64(l), 3)
    }

    w.writeBits(0, 1)

    for _, c := range codes {
        w.writeCode(lengths[c.Depth])
    }
}