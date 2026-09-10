package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "math"
    "slices"
    //------------------------------
    //imaging
    //------------------------------
    "image/color"
    //------------------------------
    //errors
    //------------------------------
    //"log"
    "errors"
)

type transform int

const (
    transformPredict        = transform(0)
    transformColor          = transform(1)
    transformSubGreen       = transform(2)
    transformColorIndexing  = transform(3)     
)

func applyPredictTransform(pixels []color.NRGBA, width, height, transBits int) (int, int, []color.NRGBA) {
    bw := (width + (1 << transBits) - 1) >> transBits
    bh := (height + (1 << transBits) - 1) >> transBits

    blocks := make([]color.NRGBA, bw * bh)
    deltas := make([]color.NRGBA, width * height)
    
    accum := [][]int{
        make([]int, 256),
        make([]int, 256),
        make([]int, 256),
        make([]int, 256),
        make([]int, 40),
    }

    histos := make([][]int, len(accum))
    for i := range accum {
        histos[i] = make([]int, len(accum[i]))
    }

    for y := 0; y < bh; y++ {
        for x := 0; x < bw; x++ {
            mx := min((x + 1) << transBits, width)
            my := min((y + 1) << transBits, height)

            var best int
            var bestEntropy float64
            for i := 0; i < 14; i++ {
                for j := range accum {
                    copy(histos[j], accum[j])
                }

                for tx := x << transBits; tx < mx; tx++ {
                    for ty := y << transBits; ty < my; ty++ {
                        d := applyFilter(pixels, width, tx, ty, i)

                        off := ty * width + tx
                        histos[0][int(uint8(pixels[off].R - d.R))]++
                        histos[1][int(uint8(pixels[off].G - d.G))]++
                        histos[2][int(uint8(pixels[off].B - d.B))]++
                        histos[3][int(uint8(pixels[off].A - d.A))]++
                    }
                }

                var total float64
                for _, histo := range histos {
                    sum := 0
                    sumSquares := 0
                
                    for _, count := range histo {
                        sum += count
                        sumSquares += count * count
                    }
                
                    if sum == 0 {
                        continue
                    }
                
                    total += 1.0 - float64(sumSquares) / (float64(sum) * float64(sum))    
                }

                if i == 0 || total < bestEntropy {
                    bestEntropy = total
                    best = i
                }
            }

            for tx := x << transBits; tx < mx; tx++ {
                for ty := y << transBits; ty < my; ty++ {
                    d := applyFilter(pixels, width, tx, ty, best)

                    off := ty * width + tx
                    deltas[off] = color.NRGBA{
                        R: uint8(pixels[off].R - d.R),
                        G: uint8(pixels[off].G - d.G),
                        B: uint8(pixels[off].B - d.B),
                        A: uint8(pixels[off].A - d.A),
                    }

                    accum[0][int(uint8(pixels[off].R - d.R))]++
                    accum[1][int(uint8(pixels[off].G - d.G))]++
                    accum[2][int(uint8(pixels[off].B - d.B))]++
                    accum[3][int(uint8(pixels[off].A - d.A))]++
                }
            }

            blocks[y * bw + x] = color.NRGBA{0, byte(best), 0, 255}
        }
    }
    
    copy(pixels, deltas)
    
    return bw, bh, blocks
}

func applyFilter(pixels []color.NRGBA, width, x, y, prediction int) color.NRGBA {
    if x == 0 && y == 0 {
        return color.NRGBA{0, 0, 0, 255}
    } else if x == 0 {
        return pixels[(y - 1) * width + x]
    } else if y == 0 {
        return pixels[y * width + (x - 1)]
    }
    
    t := pixels[(y - 1) * width + x]
    l := pixels[y * width + (x - 1)]

    tl := pixels[(y - 1) * width + (x - 1)]
    tr := pixels[(y - 1) * width + (x + 1)]

    avarage2 := func(a, b color.NRGBA) color.NRGBA {
        return color.NRGBA {
            uint8((int(a.R) + int(b.R)) / 2), 
            uint8((int(a.G) + int(b.G)) / 2),  
            uint8((int(a.B) + int(b.B)) / 2),  
            uint8((int(a.A) + int(b.A)) / 2),
        }
    }

    filters := []func(t, l, tl, tr color.NRGBA) color.NRGBA {
        func(t, l, tl, tr color.NRGBA) color.NRGBA { return color.NRGBA{0, 0, 0, 255} },
        func(t, l, tl, tr color.NRGBA) color.NRGBA { return l },
        func(t, l, tl, tr color.NRGBA) color.NRGBA { return t },
        func(t, l, tl, tr color.NRGBA) color.NRGBA { return tr },
        func(t, l, tl, tr color.NRGBA) color.NRGBA { return tl },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            return avarage2(avarage2(l, tr), t)
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            return avarage2(l, tl)
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            return avarage2(l, t)
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            return avarage2(tl, t)
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            return avarage2(t, tr)
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            return avarage2(avarage2(l, tl), avarage2(t, tr))
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA { 
            pr := float64(l.R) + float64(t.R) - float64(tl.R)
            pg := float64(l.G) + float64(t.G) - float64(tl.G)
            pb := float64(l.B) + float64(t.B) - float64(tl.B)
            pa := float64(l.A) + float64(t.A) - float64(tl.A)

            // Manhattan distances to estimates for left and top pixels.
            pl := math.Abs(pa - float64(l.A)) + math.Abs(pr - float64(l.R)) + 
                  math.Abs(pg - float64(l.G)) + math.Abs(pb - float64(l.B))
            pt := math.Abs(pa - float64(t.A)) + math.Abs(pr - float64(t.R)) + 
                  math.Abs(pg - float64(t.G)) + math.Abs(pb - float64(t.B))

            if pl < pt {
                return l
            }

            return t
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            return color.NRGBA{
                uint8(max(min(int(l.R) + int(t.R) - int(tl.R), 255), 0)),
                uint8(max(min(int(l.G) + int(t.G) - int(tl.G), 255), 0)),
                uint8(max(min(int(l.B) + int(t.B) - int(tl.B), 255), 0)),
                uint8(max(min(int(l.A) + int(t.A) - int(tl.A), 255), 0)),
            }
        },
        func(t, l, tl, tr color.NRGBA) color.NRGBA {
            a := avarage2(l, t)

            return color.NRGBA{
                uint8(max(min(int(a.R) + (int(a.R) - int(tl.R)) / 2, 255), 0)),
                uint8(max(min(int(a.G) + (int(a.G) - int(tl.G)) / 2, 255), 0)),
                uint8(max(min(int(a.B) + (int(a.B) - int(tl.B)) / 2, 255), 0)),
                uint8(max(min(int(a.A) + (int(a.A) - int(tl.A)) / 2, 255), 0)),
            }
        },
    }
    
    return filters[prediction](t, l, tl, tr)
}

func applyColorTransform(pixels []color.NRGBA, width, height, transBits int) (int, int, []color.NRGBA) {
    bw := (width + (1 << transBits) - 1) >> transBits
    bh := (height + (1 << transBits) - 1) >> transBits

    blocks := make([]color.NRGBA, bw * bh)
    deltas := make([]color.NRGBA, width * height)
    
    //TODO: analyze block and pick best Color transform Element (CTE)
    cte := color.NRGBA {
        R: 1,   //red to blue
        G: 2,   //green to blue
        B: 3,   //green to red
        A: 255,
    }
    
    for y := 0; y < bh; y++ {
        for x := 0; x < bw; x++ {
            mx := min((x + 1) << transBits, width)
            my := min((y + 1) << transBits, height)

            for tx := x << transBits; tx < mx; tx++ {
                for ty := y << transBits; ty < my; ty++ {
                    off := ty * width + tx

                    r := int(int8(pixels[off].R))
                    g := int(int8(pixels[off].G))
                    b := int(int8(pixels[off].B))
                
                    b -= int(int8((int16(int8(cte.G)) * int16(g)) >> 5))
                    b -= int(int8((int16(int8(cte.R)) * int16(r)) >> 5))
                    r -= int(int8((int16(int8(cte.B)) * int16(g)) >> 5))
                    
                    pixels[off].R = uint8(r & 0xff)
                    pixels[off].B = uint8(b & 0xff)

                    deltas[off] = pixels[off]
                }
            }

            blocks[y * bw + x] = cte
        }
    }
    
    copy(pixels, deltas)
    
    return bw, bh, blocks
}

func applySubtractGreenTransform(pixels []color.NRGBA) {
    for i, _ := range pixels {
        pixels[i].R = pixels[i].R - pixels[i].G
        pixels[i].B = pixels[i].B - pixels[i].G
    }
}

func applyPaletteTransform(pixels *[]color.NRGBA, width, height int) ([]color.NRGBA, int, error) {
    var pal []color.NRGBA
    for _, p := range (*pixels) {
        if !slices.Contains(pal, p) {
            pal = append(pal, p)
        }
   
        if len(pal) > 256 {
            return nil, 0, errors.New("palette exceeds 256 colors")
        }
    }

    size := 1
    if len(pal) <= 2 {
        size = 8
    } else if len(pal) <= 4 {
        size = 4
    } else if len(pal) <= 16 {
        size = 2
    }
    
    pw := (width + size - 1) / size

    packed := make([]color.NRGBA, pw * height)
    for y := 0; y < height; y++ {
        for x := 0; x < pw; x++ {
            pack := 0
            for i := 0; i < size; i++ {
                px := x * size + i
                if px >= width {
                    break
                }

                idx := slices.Index(pal, (*pixels)[y * width + px])
                pack |= int(idx) << (i * (8 / size))
            }

            packed[y * pw + x] = color.NRGBA{G: uint8(pack), A: 255}
        }
    }

    *pixels = packed
    
    for i := len(pal) - 1; i > 0; i-- {
        pal[i] = color.NRGBA{
            R: pal[i].R - pal[i - 1].R,
            G: pal[i].G - pal[i - 1].G,
            B: pal[i].B - pal[i - 1].B,
            A: pal[i].A - pal[i - 1].A,
        }
    }

    return pal, pw, nil
}
