package cv

import (
        "math"
)

func VarH(h []float64) float64 {
    ret := 0.0
    n := 0.0
    for _, x := range h {
        for _, y := range h {
            d := math.Abs(x - y)
            if d > 0.5 {
                d = 1.0 - d
            }
            ret += d
            n += 1.0
        }
    }
    return ret / n
}

func RGBToHSV(r, g, b uint8) (h, s, v float64) {
        fR := float64(r) / 255
        fG := float64(g) / 255
        fB := float64(b) / 255
        max := math.Max(math.Max(fR, fG), fB)
        min := math.Min(math.Min(fR, fG), fB)
        d := max - min
        s, v = 0, max
        if max > 0 {
                s = d / max
        }
        if max == min {
                // Achromatic.
                h = 0
        } else {
                // Chromatic.
                switch max {
                case fR:
                        h = (fG - fB) / d
                        if fG < fB {
                                h += 6
                        }
                case fG:
                        h = (fB-fR)/d + 2
                case fB:
                        h = (fR-fG)/d + 4
                }
                h /= 6
        }
        return
}

func HSVToRGB(h, s, v float64) (r, g, b uint8) {
        var fR, fG, fB float64
        i := math.Floor(h * 6)
        f := h*6 - i
        p := v * (1.0 - s)
        q := v * (1.0 - f*s)
        t := v * (1.0 - (1.0-f)*s)
        switch int(i) % 6 {
        case 0:
                fR, fG, fB = v, t, p
        case 1:
                fR, fG, fB = q, v, p
        case 2:
                fR, fG, fB = p, v, t
        case 3:
                fR, fG, fB = p, q, v
        case 4:
                fR, fG, fB = t, p, v
        case 5:
                fR, fG, fB = v, p, q
        }
        r = uint8((fR * 255) + 0.5)
        g = uint8((fG * 255) + 0.5)
        b = uint8((fB * 255) + 0.5)
        return
}