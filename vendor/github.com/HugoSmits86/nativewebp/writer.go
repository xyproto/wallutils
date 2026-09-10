package nativewebp

import (
    //------------------------------
    //general
    //------------------------------
    "io"
    "bytes"
    "encoding/binary"
    //------------------------------
    //imaging
    //------------------------------
    "image"
    "image/draw"
    "image/color"
    //------------------------------
    //errors
    //------------------------------
    "errors"
)

// CompressionLevel indicates the compression level.
type CompressionLevel int

const (
	DefaultCompression CompressionLevel = 4
	BestSpeed          CompressionLevel = 0
	BestCompression    CompressionLevel = 6
)

// Options holds configuration settings for WebP encoding.
//
// Fields:
//   - UseExtendedFormat:
//     If true, wraps the VP8L frame inside a VP8X container to enable
//     metadata support such as EXIF, ICC color profiles, and XMP.
//     This does not affect image compression itself, as VP8L remains
//     the underlying image encoding format.
//
//   - CompressionLevel:
//     Controls the encoder effort and compression trade-off.
//
//     Higher compression levels may improve file size by enabling more
//     expensive analysis and transform selection steps, at the cost of
//     increased CPU usage and encoding time.
type Options struct {
    UseExtendedFormat   bool
    CompressionLevel    CompressionLevel
}

// Animation holds configuration settings for WebP animations.
//
// It allows encoding a sequence of frames with individual timing and disposal options,
// supporting features like looping and background color settings.
//
// Fields:
//   - Images: A list of frames to be displayed in sequence.
//   - Durations: Timing for each frame in milliseconds, matching the Images slice.
//   - Disposals: Disposal methods for frames after display; 0 = keep, 1 = clear to background.
//   - LoopCount: Number of times the animation should repeat; 0 means infinite looping.
//   - BackgroundColor: Canvas background color in BGRA order, used for clear operations.
type Animation struct {
    Images              []image.Image
    Durations           []uint
    Disposals           []uint
    LoopCount           uint16
    BackgroundColor     uint32
}

// Encode writes the provided image.Image to the specified io.Writer in WebP format.
//
// This function always encodes the image using VP8L (lossless WebP). If `UseExtendedFormat`
// is enabled, it wraps the VP8L frame inside a VP8X container, allowing the use of metadata
// such as EXIF, ICC color profiles, or XMP metadata.
//
// Note: VP8L already supports transparency, so VP8X is **not required** for alpha support.
//
// Parameters:
//   w   - The destination writer where the encoded WebP image will be written.
//   img - The input image to be encoded.
//   o   - Pointer to Options containing encoding settings:
//         - UseExtendedFormat: If true, wraps the image in a VP8X container to enable 
//           extended WebP features like metadata.
//         - CompressionLevel: Controls encoding effort and compression trade-off.
//
// Returns:
//   An error if encoding fails or writing to the io.Writer encounters an issue.
func Encode(w io.Writer, img image.Image, o *Options) error {
    method := getMethodLevel(DefaultCompression)
    if o != nil {
        method = getMethodLevel(o.CompressionLevel)
    }

    stream, hasAlpha, err := writeBitStream(img, method)
    if err != nil {
        return err
    }

    buf := &bytes.Buffer{}

    if o != nil && o.UseExtendedFormat {
        writeChunkVP8X(buf, img.Bounds(), hasAlpha, false)
    }

    buf.Write([]byte("VP8L"))
    binary.Write(buf, binary.LittleEndian, uint32(stream.Len()))
    buf.Write(stream.Bytes())

    w.Write([]byte("RIFF"))
    binary.Write(w, binary.LittleEndian, uint32(4 + buf.Len()))

    w.Write([]byte("WEBP"))
    w.Write(buf.Bytes())

    return nil
}

// EncodeAll writes the provided animation sequence to the specified io.Writer in WebP format.
//
// This function encodes a list of frames as a WebP animation using the VP8X container, which
// supports features like looping, frame timing, disposal methods, and background color settings.
// Each frame is individually compressed using the VP8L (lossless) format.
//
// Note: Even if `UseExtendedFormat` is not explicitly set, animations always use the VP8X container
// because it is required for WebP animation support.
//
// Parameters:
//   w   - The destination writer where the encoded WebP animation will be written.
//   ani - Pointer to Animation containing the frames and animation settings:
//         - Images: List of frames to encode.
//         - Durations: Display times for each frame in milliseconds.
//         - Disposals: Disposal methods after frame display (keep or clear).
//         - LoopCount: Number of times the animation should loop (0 = infinite).
//         - BackgroundColor: Background color for the canvas, used when clearing.
//   o   - Pointer to Options containing additional encoding settings:
//         - UseExtendedFormat: Currently unused for animations, but accepted for consistency.
//         - CompressionLevel: Controls encoding effort and compression trade-off.
//
// Returns:
//   An error if encoding fails or writing to the io.Writer encounters an issue.
func EncodeAll(w io.Writer, ani *Animation, o *Options) error {
    method := getMethodLevel(DefaultCompression)
    if o != nil {
        method = getMethodLevel(o.CompressionLevel)
    }

    frames, alpha, err := writeFrames(ani, method)
    if err != nil {
        return err
    }

    var bounds image.Rectangle
    for _, img := range ani.Images {
        bounds.Max.X = max(img.Bounds().Max.X, bounds.Max.X)
        bounds.Max.Y = max(img.Bounds().Max.Y, bounds.Max.Y)
    }

    buf := &bytes.Buffer{}

    writeChunkVP8X(buf, bounds, alpha, true)

    buf.Write([]byte("ANIM"))
    binary.Write(buf, binary.LittleEndian, uint32(6))
    binary.Write(buf, binary.LittleEndian, uint32(ani.BackgroundColor))
    binary.Write(buf, binary.LittleEndian, uint16(ani.LoopCount))

    buf.Write(frames.Bytes())

    w.Write([]byte("RIFF"))
    binary.Write(w, binary.LittleEndian, uint32(4 + buf.Len()))

    w.Write([]byte("WEBP"))
    w.Write(buf.Bytes())

    return nil
}

func getMethodLevel(lvl CompressionLevel) int {
    switch lvl {
        case BestSpeed:
            return 0
        case DefaultCompression:
            return 4
        case BestCompression:
            return 6
        default:
            return 4
    }
}

func writeChunkVP8X(buf *bytes.Buffer, bounds image.Rectangle, flagAlpha, flagAni bool) {
    buf.Write([]byte("VP8X"))
    binary.Write(buf, binary.LittleEndian, uint32(10))

    var flags byte
    if flagAni {
        flags |= 1 << 1
    }

    if flagAlpha {
        flags |= 1 << 4
    }

    binary.Write(buf, binary.LittleEndian, flags)
    buf.Write([]byte{0x00, 0x00, 0x00})

    dx := bounds.Dx() - 1
    dy := bounds.Dy() - 1

    buf.Write([]byte{byte(dx), byte(dx >> 8), byte(dx >> 16)})
    buf.Write([]byte{byte(dy), byte(dy >> 8), byte(dy >> 16)})
}

func writeFrames(ani *Animation, method int) (*bytes.Buffer, bool, error) {
    if len(ani.Images) == 0 {
        return nil, false, errors.New("must provide at least one image")
    }

    if len(ani.Images) != len(ani.Durations) {
        return nil, false, errors.New("mismatched image and durations lengths")
    }

    if len(ani.Images) != len(ani.Disposals) {
        return nil, false, errors.New("mismatched image and disposals lengths")
    }

    for i := 0; i < len(ani.Images); i++ {
        ani.Durations[i] = min(ani.Durations[i], 1 << 24 - 1)
        ani.Disposals[i] = min(ani.Disposals[i], 1)
    }

    buf := &bytes.Buffer{}
    
    var hasAlpha bool
    for i, img := range ani.Images {
        stream, alpha, err := writeBitStream(img, method)
        if err != nil {
            return nil, false, err
        }
    
        hasAlpha = hasAlpha || alpha

        w := &bitWriter{Buffer: buf}
        w.writeBytes([]byte("ANMF"))
        w.writeBits(uint64(16 + 8 + stream.Len()), 32)
    
        // WebP specs requires frame offsets to be divided by 2
        w.writeBits(uint64(img.Bounds().Min.X / 2), 24)
        w.writeBits(uint64(img.Bounds().Min.Y / 2), 24)
    
        w.writeBits(uint64(img.Bounds().Dx() - 1), 24)
        w.writeBits(uint64(img.Bounds().Dy() - 1), 24)
    
        w.writeBits(uint64(ani.Durations[i]), 24)
        w.writeBits(uint64(ani.Disposals[i]), 1)
        w.writeBits(uint64(0), 1)
        w.writeBits(uint64(0), 6)
    
        w.writeBytes([]byte("VP8L"))
        w.writeBits(uint64(stream.Len()), 32)
        w.Buffer.Write(stream.Bytes())
    }

    return buf, hasAlpha, nil
}

func writeBitStream(img image.Image, method int) (*bytes.Buffer, bool, error) {
    if img == nil {
        return nil, false, errors.New("image is nil")
    }

    if img.Bounds().Dx() < 1 || img.Bounds().Dy() < 1 {
        return nil, false, errors.New("invalid image size")
    }

    if img.Bounds().Dx() > 1 << 14 || img.Bounds().Dy() > 1 << 14 {
        return nil, false, errors.New("invalid image size")
    }

    _, isIndexed := img.(*image.Paletted)

    rgba := image.NewNRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
    draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)

    b := &bytes.Buffer{}
    s := &bitWriter{Buffer: b}

    writeBitStreamHeader(s, rgba.Bounds(), !rgba.Opaque())

    var transforms [4]bool
    transforms[transformPredict] = !isIndexed
    transforms[transformColor] = false
    transforms[transformSubGreen] = !isIndexed
    transforms[transformColorIndexing] = isIndexed

    histoBits := getHistoBits(method, isIndexed, img.Bounds().Dx(), img.Bounds().Dy())
    transBits := getTransformBits(method, histoBits)

    err := writeBitStreamData(s, rgba, 11, histoBits, transBits, transforms)
    if err != nil {
        return nil, false, err
    }
    
    s.alignByte()

    if b.Len() % 2 != 0 {
        b.Write([]byte{0x00})
    }

    return b, !rgba.Opaque(), nil
}

func writeBitStreamHeader(w *bitWriter, bounds image.Rectangle, hasAlpha bool) {
    w.writeBits(0x2f, 8)

    w.writeBits(uint64(bounds.Dx() - 1), 14)
    w.writeBits(uint64(bounds.Dy() - 1), 14)

    if hasAlpha {
        w.writeBits(1, 1)
    } else {
        w.writeBits(0, 1)
    }

    w.writeBits(0, 3)
}

func writeBitStreamData(w *bitWriter, img image.Image, colorBits, histoBits, transBits int, transforms [4]bool) error {
    pixels, err := flatten(img)
    if err != nil {
        return err
    }

    width := img.Bounds().Dx()
    height := img.Bounds().Dy()

    if transforms[transformColorIndexing] {
        w.writeBits(1, 1)
        w.writeBits(3, 2)
       
        pal, pw, err := applyPaletteTransform(&pixels, width, height)
        if err != nil {
            return err
        }

        width = pw
       
        w.writeBits(uint64(len(pal) - 1), 8);
        writeImageData(w, pal, len(pal), 1, false, 0);
    }

    if transforms[transformSubGreen] {
        w.writeBits(1, 1)
        w.writeBits(2, 2)

        applySubtractGreenTransform(pixels)
    }

    if transforms[transformColor] {
        w.writeBits(1, 1)
        w.writeBits(1, 2)

        bw, bh, blocks := applyColorTransform(pixels, width, height, transBits)

        w.writeBits(uint64(transBits - 2), 3);
        writeImageData(w, blocks, bw, bh, false, 0)
    }

    if transforms[transformPredict] {
        w.writeBits(1, 1)
        w.writeBits(0, 2)

        bw, bh, blocks := applyPredictTransform(pixels, width, height, transBits)

        w.writeBits(uint64(transBits - 2), 3);
        writeImageData(w, blocks, bw, bh, false, 0)
    }

    w.writeBits(0, 1) // end of transform
    writeImageData(w, pixels, width, height, true, colorBits)

    return nil
}

func writeImageData(w *bitWriter, pixels []color.NRGBA, width, height int, isRecursive bool, colorBits int) {
    if colorBits > 0 {
        w.writeBits(1, 1)
        w.writeBits(uint64(colorBits), 4) 
    } else {
        w.writeBits(0, 1)
    }

    if isRecursive {
        w.writeBits(0, 1)
    }

    encoded := encodeImageData(pixels, width, height, colorBits)
    histos := computeHistograms(encoded, colorBits)

    var codes [][]huffmanCode
    for i := 0; i < 5; i++ {
        // WebP specs requires Huffman codes with maximum depth of 15
        c := buildhuffmanCodes(histos[i], 15)
        codes = append(codes, c)

        writehuffmanCodes(w, c)
    }

    for i := 0; i < len(encoded); i ++ {
        w.writeCode(codes[0][encoded[i + 0]])
        if encoded[i + 0] < 256 {
            w.writeCode(codes[1][encoded[i + 1]])
            w.writeCode(codes[2][encoded[i + 2]])
            w.writeCode(codes[3][encoded[i + 3]])
            i += 3
        } else if encoded[i + 0] < 256 + 24 {
            cnt := prefixEncodeBits(int(encoded[i + 0]) - 256)
            w.writeBits(uint64(encoded[i + 1]), cnt);

            w.writeCode(codes[4][encoded[i + 2]])

            cnt = prefixEncodeBits(int(encoded[i + 2]))
            w.writeBits(uint64(encoded[i + 3]), cnt);
            i += 3
        }
    }
}

func encodeImageData(pixels []color.NRGBA, width, height, colorBits int) []int {
    head := make([]int, 1 << 18)
    prev := make([]int, len(pixels))
    cache := make([]color.NRGBA, 1 << colorBits)

    encoded := make([]int, len(pixels) * 4)
    cnt := 0

    var distances = []int {
        96,   73,  55,  39,  23,  13,   5,  1,  255, 255, 255, 255, 255, 255, 255, 255,
        101,  78,  58,  42,  26,  16,   8,  2,    0,   3,  9,   17,  27,  43,  59,  79,
        102,  86,  62,  46,  32,  20,  10,  6,    4,   7,  11,  21,  33,  47,  63,  87,
        105,  90,  70,  52,  37,  28,  18,  14,  12,  15,  19,  29,  38,  53,  71,  91,
        110,  99,  82,  66,  48,  35,  30,  24,  22,  25,  31,  36,  49,  67,  83, 100,
        115, 108,  94,  76,  64,  50,  44,  40,  34,  41,  45,  51,  65,  77,  95, 109,
        118, 113, 103,  92,  80,  68,  60,  56,  54,  57,  61,  69,  81,  93, 104, 114,
        119, 116, 111, 106,  97,  88,  84,  74,  72,  75,  85,  89,  98, 107, 112, 117,
    }
    
    for i := 0; i < len(pixels); i++ {
        if i + 2 < len(pixels) {
            h := hash(pixels[i + 0], 18)
            h ^= hash(pixels[i + 1], 18) * 0x9e3779b9
            h ^= hash(pixels[i + 2], 18) * 0x85ebca6b
            h = h % (1 << 18)

            cur := head[h] - 1
            prev[i] = head[h]
            head[h] = i + 1

            dis := 0
            streak := 0
            for j := 0; j < 128; j++ {
                // 1 << 20: sliding window size is 2^20 (1,048,576) per WebP specs.
                // 120: reserved margin for offset adjustments.
                if cur == -1 || i - cur >= 1 << 20 - 120 {
                    break
                }

                l := 0
                // Limit the maximum match length to 4096 pixels per WebP specs.
                for i + l < len(pixels) && l < 4096 {
                    if pixels[i + l] != pixels[cur + l] {
                        break
                    }
                    l++
                }

                if l > streak {
                    streak = l
                    dis = i - cur
                }

                cur = prev[cur] - 1
            }

            // Only use the match if it is at least 3 pixels long per WebP specs.
            if streak >= 3 {
                for j := 0; j < streak; j++ {
                    h := hash(pixels[i + j], colorBits)
                    cache[h] = pixels[i + j]
                }
                
                y := dis / width
                x := dis - y * width
            
                code := dis + 120
                if x <= 8 && y < 8 {
                    code = distances[y * 16 + 8 - x] + 1
                } else if x > width - 8 && y < 7 {
                    code = distances[(y + 1) * 16 + 8 + (width - x)] + 1
                }

                s, l := prefixEncodeCode(streak)
                encoded[cnt + 0] = int(s + 256)
                encoded[cnt + 1] = int(l)

                s, l = prefixEncodeCode(code)
                encoded[cnt + 2] = int(s)
                encoded[cnt + 3] = int(l)
                cnt += 4
    
                i += streak - 1
                continue
            }
        }

        p := pixels[i]
        if colorBits > 0 {
            hash := hash(p, colorBits)

            if i > 0 && cache[hash] == p {
                encoded[cnt] = int(hash + 256 + 24)
                cnt++
                continue
            }

            cache[hash] = p
        }

        encoded[cnt+0] = int(p.G)
        encoded[cnt+1] = int(p.R)
        encoded[cnt+2] = int(p.B)
        encoded[cnt+3] = int(p.A)
        cnt += 4
    }

    return encoded[:cnt]
}

func prefixEncodeCode(n int) (int, int) {
    if n <= 5 {
        return max(0, n - 1), 0
    }

    shift := 0
    rem := n - 1
    for rem > 3 {
        rem >>= 1
        shift += 1
    }

    if rem == 2 {
        return 2 + 2 * shift, n - (2 << shift) - 1
    }

    return 3 + 2 * shift, n - (3 << shift) - 1
}

func prefixEncodeBits(prefix int) int {
    if prefix < 4 {
        return 0
    }

    return (prefix - 2) >> 1
}

func hash(c color.NRGBA, shifts int) uint32 {
    //hash formula including magic number 0x1e35a7bd comes directly from WebP specs!
    x := uint32(c.A) << 24 | uint32(c.R) << 16 | uint32(c.G) << 8 | uint32(c.B)
    return (x * 0x1e35a7bd) >> (32 - min(shifts, 32))
}

func computeHistograms(pixels []int, colorBits int) [][]int {
    c := 0
    if colorBits > 0 {
        c = 1 << colorBits
    }

    histos := [][]int{
        make([]int, 256 + 24 + c),
        make([]int, 256),
        make([]int, 256),
        make([]int, 256),
        make([]int, 40),
    }

    for i := 0; i < len(pixels); i++ {
        histos[0][pixels[i]]++
        if(pixels[i] < 256) {
            histos[1][pixels[i + 1]]++
            histos[2][pixels[i + 2]]++
            histos[3][pixels[i + 3]]++
            i += 3
        } else if pixels[i] < 256 + 24 {
            histos[4][pixels[i + 2]]++
            i += 3
        }
    }

    return histos
}

func getTransformBits(method, histoBits int) int {
    maxBits := 5
    if method < 4 {
        maxBits = 6
    } else if method > 4 {
        maxBits = 4
    } 

    return min(histoBits, maxBits)
}

func getHistoBits(method int, isIndexed bool, width, height int) int {
    bits := 9 - method
    if !isIndexed {
        bits = 7 - method
    }

    subSample := func (size, bits int) int {
        return (size + (1 << bits) - 1) >> bits
    }
 
    bits = min(max(bits, MIN_HUFFMAN_BITS), MAX_HUFFMAN_BITS)
    size := subSample(width, bits) * subSample(height, bits)

    for bits < MAX_HUFFMAN_BITS && size > MAX_HUFF_IMAGE_SIZE {
        bits++
        size = subSample(width, bits) * subSample(height, bits)
    }

    for bits > MIN_HUFFMAN_BITS && size == 1 {
        size = subSample(width, bits - 1) * subSample(height, bits - 1)
        if size != 1 {
            break
        }

        bits--
    }

    return bits
}

func flatten(img image.Image) ([]color.NRGBA, error) {
    w := img.Bounds().Dx()
    h := img.Bounds().Dy()

    rgba, ok := img.(*image.NRGBA)
    if !ok {
        return nil, errors.New("unsupported image format")
    }

    pixels := make([]color.NRGBA, w * h)
    for y := 0; y < h; y++ {
        for x := 0; x < w; x++ {
            i := rgba.PixOffset(x, y)
            s := rgba.Pix[i : i + 4 : i + 4]

            pixels[y * w + x].R = uint8(s[0])
            pixels[y * w + x].G = uint8(s[1])
            pixels[y * w + x].B = uint8(s[2])
            pixels[y * w + x].A = uint8(s[3])
        }
    }

    return pixels, nil
}