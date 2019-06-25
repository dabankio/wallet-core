package wcg

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type long10 struct {
	_0, _1, _2, _3, _4, _5, _6, _7, _8, _9 int64
}

const (
	p25                int64 = (1 << 25) - 1 /*33554431*/
	p26                int64 = (1 << 26) - 1 /*67108863*/
	nxtSignatureOffset       = 1 + 1 + 4 + 2 + 32 + 8 + 8 + 8 + 32
	nxtSignatureEnd          = nxtSignatureOffset + 64
)

var (
	order = [64]byte{
		237, 211, 245, 92,
		26, 99, 18, 88,
		214, 156, 247, 162,
		222, 249, 222, 20,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 16,
	}

	/* smallest multiple of the order that's >= 2^255 */
	orderTimes8 = [64]byte{
		104, 159, 174, 231,
		210, 24, 147, 192,
		178, 230, 188, 23,
		245, 206, 247, 166,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 128,
	}

	/* constants 2Gy and 1/(2Gy) */
	base2Y = long10{
		39999547, 18689728, 59995525, 1648697, 57546132,
		24010086, 19059592, 5425144, 63499247, 16420658,
	}

	baseR2Y = long10{
		5744, 8160848, 4790893, 13779497, 35730846,
		12541209, 49101323, 30047407, 40071253, 6226132,
	}
)

/* mulaSmall calculate p[m..n+m-1] = q[m..n+m-1] + z * x */
/* n is the size of x */
/* n+m is the size of p and q */
func mulaSmall(p *[64]byte, q *[64]byte, m int, x *[64]byte, n, z int) int {
	var value int
	for i := 0; i < n; i++ {
		value += int(q[i+m]&0xFF) + z*int(x[i]&0xFF)
		p[i+m] = byte(value)
		value >>= 8
	}
	return value
}

/* mula32 calculate p += x * y * z  where z is a small integer
 * x is size 32, y is size t, p is size 32+t
 * y is allowed to overlap with p+32 if you don't care about the upper half  */
func mula32(p *[64]byte, x, y *[64]byte, t, z int) int {
	n := 31
	var w, i int
	for ; i < t; i++ {
		zy := z * int(y[i]&0xFF)
		w += mulaSmall(p, p, i, x, n, zy) + int(p[i+n]&0xFF) + zy*int(x[n]&0xFF)
		p[i+n] = byte(w)
		w >>= 8
	}
	p[i+n] = byte(w + int(p[i+n]&0xFF))
	return w >> 8
}

/* divmod divide r (size n) by d (size t), returning quotient q and remainder r
 * quotient is size n-t+1, remainder is size t
 * requires t > 0 && d[t-1] != 0
 * requires that r[-1] and d[-1] are valid memory locations
 * q may overlap with r+t */
func divmod(q *[64]byte, r *[64]byte, n int, d *[64]byte, t int) {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	rn := 0
	dt := int(d[t-1]&0xFF) << 8
	if t > 1 {
		dt |= int(d[t-2] & 0xFF)
	}
	for n >= t {
		n--
		z := (rn << 16) | (int(r[n]&0xFF) << 8)
		if n > 0 {
			z |= int(r[n-1] & 0xFF)
		}
		z /= dt
		rn += mulaSmall(r, r, n-t+1, d, t, -z)
		q[n-t+1] = byte((z + rn) & 0xFF)
		mulaSmall(r, r, n-t+1, d, t, -rn)
		rn = int(r[n] & 0xFF)
		r[n] = 0
	}
	r[t-1] = byte(rn)
}

// xToY2 calculate Y^2 = X^3 + 486662 X^2 + X
// t is a temporary
func xToY2(t, y2, x *long10) {
	sqr(t, x)
	mulSmall(y2, x, 486662)
	add(t, t, y2)
	t._0++
	mul(y2, t, x)
}

// add Add two numbers.  The inputs must be in reduced form, and the
// output isn't, so to do another addition or subtraction on the output,
// first multiply it by one to reduce it.
func add(xy, x, y *long10) {
	xy._0 = x._0 + y._0
	xy._1 = x._1 + y._1
	xy._2 = x._2 + y._2
	xy._3 = x._3 + y._3
	xy._4 = x._4 + y._4
	xy._5 = x._5 + y._5
	xy._6 = x._6 + y._6
	xy._7 = x._7 + y._7
	xy._8 = x._8 + y._8
	xy._9 = x._9 + y._9
}

/* sub subtract two numbers.  The inputs must be in reduced form, and the
 * output isn't, so to do another addition or subtraction on the output,
 * first multiply it by one to reduce it.
 */
func sub(xy, x, y *long10) {
	xy._0 = x._0 - y._0
	xy._1 = x._1 - y._1
	xy._2 = x._2 - y._2
	xy._3 = x._3 - y._3
	xy._4 = x._4 - y._4
	xy._5 = x._5 - y._5
	xy._6 = x._6 - y._6
	xy._7 = x._7 - y._7
	xy._8 = x._8 - y._8
	xy._9 = x._9 - y._9
}

/* mulSmall Multiply a number by a small integer in range -185861411 .. 185861411.
 * The output is in reduced form, the input x need not be.  x and xy may point
 * to the same buffer. */
func mulSmall(xy, x *long10, y int64) *long10 {
	t := x._8 * y
	xy._8 = t & ((1 << 26) - 1)
	t = (t >> 26) + (x._9 * y)
	xy._9 = t & ((1 << 25) - 1)
	t = 19*(t>>25) + (x._0 * y)
	xy._0 = t & ((1 << 26) - 1)
	t = (t >> 26) + (x._1 * y)
	xy._1 = t & ((1 << 25) - 1)
	t = (t >> 25) + (x._2 * y)
	xy._2 = t & ((1 << 26) - 1)
	t = (t >> 26) + (x._3 * y)
	xy._3 = t & ((1 << 25) - 1)
	t = (t >> 25) + (x._4 * y)
	xy._4 = t & ((1 << 26) - 1)
	t = (t >> 26) + (x._5 * y)
	xy._5 = t & ((1 << 25) - 1)
	t = (t >> 25) + (x._6 * y)
	xy._6 = t & ((1 << 26) - 1)
	t = (t >> 26) + (x._7 * y)
	xy._7 = t & ((1 << 25) - 1)
	t = (t >> 25) + xy._8
	xy._8 = t & ((1 << 26) - 1)
	xy._9 += t >> 26
	return xy
}

/* mul Multiply two numbers.  The output is in reduced form, the inputs need not
 * be. */
func mul(xy, x, y *long10) *long10 {
	/* sahn0:
	 * Using local variables to avoid class access.
	 * This seem to improve performance a bit...
	 */
	x0 := x._0
	x1 := x._1
	x2 := x._2
	x3 := x._3
	x4 := x._4
	x5 := x._5
	x6 := x._6
	x7 := x._7
	x8 := x._8
	x9 := x._9

	y0 := y._0
	y1 := y._1
	y2 := y._2
	y3 := y._3
	y4 := y._4
	y5 := y._5
	y6 := y._6
	y7 := y._7
	y8 := y._8
	y9 := y._9

	t := (x0 * y8) + (x2 * y6) + (x4 * y4) + (x6 * y2) + (x8 * y0) + 2*((x1*y7)+(x3*y5)+(x5*y3)+(x7*y1)) + 38*(x9*y9)
	xy._8 = t & ((1 << 26) - 1)
	t = (t >> 26) + (x0 * y9) + (x1 * y8) + (x2 * y7) + (x3 * y6) + (x4 * y5) + (x5 * y4) + (x6 * y3) + (x7 * y2) + (x8 * y1) + (x9 * y0)
	xy._9 = t & ((1 << 25) - 1)
	t = (x0 * y0) + 19*((t>>25)+(x2*y8)+(x4*y6)+(x6*y4)+(x8*y2)) + 38*((x1*y9)+(x3*y7)+(x5*y5)+(x7*y3)+(x9*y1))
	xy._0 = t & ((1 << 26) - 1)
	t = (t >> 26) + (x0 * y1) + (x1 * y0) + 19*((x2*y9)+(x3*y8)+(x4*y7)+(x5*y6)+(x6*y5)+(x7*y4)+(x8*y3)+(x9*y2))
	xy._1 = t & ((1 << 25) - 1)
	t = (t >> 25) + (x0 * y2) + (x2 * y0) + 19*((x4*y8)+(x6*y6)+(x8*y4)) + 2*(x1*y1) + 38*((x3*y9)+(x5*y7)+(x7*y5)+(x9*y3))
	xy._2 = t & ((1 << 26) - 1)
	t = (t >> 26) + (x0 * y3) + (x1 * y2) + (x2 * y1) + (x3 * y0) + 19*((x4*y9)+(x5*y8)+(x6*y7)+(x7*y6)+(x8*y5)+(x9*y4))
	xy._3 = t & ((1 << 25) - 1)
	t = (t >> 25) + (x0 * y4) + (x2 * y2) + (x4 * y0) + 19*((x6*y8)+(x8*y6)) + 2*((x1*y3)+(x3*y1)) + 38*((x5*y9)+(x7*y7)+(x9*y5))
	xy._4 = t & ((1 << 26) - 1)
	t = (t >> 26) + (x0 * y5) + (x1 * y4) + (x2 * y3) + (x3 * y2) + (x4 * y1) + (x5 * y0) + 19*((x6*y9)+(x7*y8)+(x8*y7)+(x9*y6))
	xy._5 = t & ((1 << 25) - 1)
	t = (t >> 25) + (x0 * y6) + (x2 * y4) + (x4 * y2) + (x6 * y0) + 19*(x8*y8) + 2*((x1*y5)+(x3*y3)+(x5*y1)) + 38*((x7*y9)+(x9*y7))
	xy._6 = t & ((1 << 26) - 1)
	t = (t >> 26) + (x0 * y7) + (x1 * y6) + (x2 * y5) + (x3 * y4) + (x4 * y3) + (x5 * y2) + (x6 * y1) + (x7 * y0) + 19*((x8*y9)+(x9*y8))
	xy._7 = t & ((1 << 25) - 1)
	t = (t >> 25) + xy._8
	xy._8 = t & ((1 << 26) - 1)
	xy._9 += t >> 26
	return xy
}

/* sqr Square a number.  Optimization of  mul25519(x2, x, x)  */
func sqr(out, x *long10) *long10 {
	x0 := x._0
	x1 := x._1
	x_2 := x._2
	x3 := x._3
	x4 := x._4
	x5 := x._5
	x6 := x._6
	x7 := x._7
	x8 := x._8
	x9 := x._9

	t := (x4 * x4) + 2*((x0*x8)+(x_2*x6)) + 38*(x9*x9) + 4*((x1*x7)+(x3*x5))
	out._8 = t & ((1 << 26) - 1)
	t = (t >> 26) + 2*((x0*x9)+(x1*x8)+(x_2*x7)+(x3*x6)+(x4*x5))
	out._9 = t & ((1 << 25) - 1)
	t = 19*(t>>25) + (x0 * x0) + 38*((x_2*x8)+(x4*x6)+(x5*x5)) + 76*((x1*x9)+(x3*x7))
	out._0 = t & ((1 << 26) - 1)
	t = (t >> 26) + 2*(x0*x1) + 38*((x_2*x9)+(x3*x8)+(x4*x7)+(x5*x6))
	out._1 = t & ((1 << 25) - 1)
	t = (t >> 25) + 19*(x6*x6) + 2*((x0*x_2)+(x1*x1)) + 38*(x4*x8) + 76*((x3*x9)+(x5*x7))
	out._2 = t & ((1 << 26) - 1)
	t = (t >> 26) + 2*((x0*x3)+(x1*x_2)) + 38*((x4*x9)+(x5*x8)+(x6*x7))
	out._3 = t & ((1 << 25) - 1)
	t = (t >> 25) + (x_2 * x_2) + 2*(x0*x4) + 38*((x6*x8)+(x7*x7)) + 4*(x1*x3) + 76*(x5*x9)
	out._4 = t & ((1 << 26) - 1)
	t = (t >> 26) + 2*((x0*x5)+(x1*x4)+(x_2*x3)) + 38*((x6*x9)+(x7*x8))
	out._5 = t & ((1 << 25) - 1)
	t = (t >> 25) + 19*(x8*x8) + 2*((x0*x6)+(x_2*x4)+(x3*x3)) + 4*(x1*x5) + 76*(x7*x9)
	out._6 = t & ((1 << 26) - 1)
	t = (t >> 26) + 2*((x0*x7)+(x1*x6)+(x_2*x5)+(x3*x4)) + 38*(x8*x9)
	out._7 = t & ((1 << 25) - 1)
	t = (t >> 25) + out._8
	out._8 = t & ((1 << 26) - 1)
	out._9 += t >> 26
	return out
}

// recip Calculates a reciprocal.  The output is in reduced form, the inputs need not
// be.  Simply calculates  y = x^(p-2)  so it's not too fast.
// When sqrtassist is true, it instead calculates y = x^((p-5)/8)
func recip(y, x *long10, sqrtassist int) {
	t0 := &long10{}
	t1 := &long10{}
	t2 := &long10{}
	t3 := &long10{}
	t4 := &long10{}

	var i int
	/* the chain for x^(2^255-21) is straight from djb's implementation */
	sqr(t1, x)      /*  2 == 2 * 1	*/
	sqr(t2, t1)     /*  4 == 2 * 2	*/
	sqr(t0, t2)     /*  8 == 2 * 4	*/
	mul(t2, t0, x)  /*  9 == 8 + 1	*/
	mul(t0, t2, t1) /* 11 == 9 + 2	*/
	sqr(t1, t0)     /* 22 == 2 * 11	*/
	mul(t3, t1, t2) /* 31 == 22 + 9
	== 2^5   - 2^0	*/
	sqr(t1, t3)     /* 2^6   - 2^1	*/
	sqr(t2, t1)     /* 2^7   - 2^2	*/
	sqr(t1, t2)     /* 2^8   - 2^3	*/
	sqr(t2, t1)     /* 2^9   - 2^4	*/
	sqr(t1, t2)     /* 2^10  - 2^5	*/
	mul(t2, t1, t3) /* 2^10  - 2^0	*/
	sqr(t1, t2)     /* 2^11  - 2^1	*/
	sqr(t3, t1)     /* 2^12  - 2^2	*/
	for i = 1; i < 5; i++ {
		sqr(t1, t3)
		sqr(t3, t1)
	} /* t3 */ /* 2^20  - 2^10	*/
	mul(t1, t3, t2) /* 2^20  - 2^0	*/
	sqr(t3, t1)     /* 2^21  - 2^1	*/
	sqr(t4, t3)     /* 2^22  - 2^2	*/
	for i = 1; i < 10; i++ {
		sqr(t3, t4)
		sqr(t4, t3)
	} /* t4 */ /* 2^40  - 2^20	*/
	mul(t3, t4, t1) /* 2^40  - 2^0	*/
	for i = 0; i < 5; i++ {
		sqr(t1, t3)
		sqr(t3, t1)
	} /* t3 */ /* 2^50  - 2^10	*/
	mul(t1, t3, t2) /* 2^50  - 2^0	*/
	sqr(t2, t1)     /* 2^51  - 2^1	*/
	sqr(t3, t2)     /* 2^52  - 2^2	*/
	for i = 1; i < 25; i++ {
		sqr(t2, t3)
		sqr(t3, t2)
	} /* t3 */ /* 2^100 - 2^50 */
	mul(t2, t3, t1) /* 2^100 - 2^0	*/
	sqr(t3, t2)     /* 2^101 - 2^1	*/
	sqr(t4, t3)     /* 2^102 - 2^2	*/
	for i = 1; i < 50; i++ {
		sqr(t3, t4)
		sqr(t4, t3)
	} /* t4 */ /* 2^200 - 2^100 */
	mul(t3, t4, t2) /* 2^200 - 2^0	*/
	for i = 0; i < 25; i++ {
		sqr(t4, t3)
		sqr(t3, t4)
	} /* t3 */ /* 2^250 - 2^50	*/
	mul(t2, t3, t1) /* 2^250 - 2^0	*/
	sqr(t1, t2)     /* 2^251 - 2^1	*/
	sqr(t2, t1)     /* 2^252 - 2^2	*/
	if sqrtassist != 0 {
		mul(y, x, t2) /* 2^252 - 3 */
	} else {
		sqr(t1, t2)    /* 2^253 - 2^3	*/
		sqr(t2, t1)    /* 2^254 - 2^4	*/
		sqr(t1, t2)    /* 2^255 - 2^5	*/
		mul(y, t1, t0) /* 2^255 - 21	*/
	}
}

// egcd32 Returns x if a contains the gcd, y if b.
// Also, the returned buffer contains the inverse of a mod b,
// as 32-byte signed.
// x and y must have 64 bytes space for temporary use.
// requires that a[-1] and b[-1] are valid memory locations
func egcd32(x *[64]byte, y *[64]byte, a *[64]byte, b *[64]byte) *[64]byte {
	var (
		an, qn, i int
		bn        = 32
	)
	for i = 0; i < 32; i++ {
		x[i], y[i] = 0, 0
	}
	x[0] = 1
	an = numsize(a, 32)
	if an == 0 {
		return y
	}
	temp := &[64]byte{}
	for {
		qn = bn - an + 1
		divmod(temp, b, bn, a, an)
		bn = numsize(b, bn)
		if bn == 0 {
			return x
		}
		mula32(y, x, temp, qn, -1)

		qn = an - bn + 1
		divmod(temp, a, an, b, bn)
		an = numsize(a, an)
		if an == 0 {
			return y
		}
		mula32(x, y, temp, qn, -1)
	}
}

/// montPrep calculate
//  t1 = ax + az
//  t2 = ax - az
func montPrep(t1, t2, ax, az *long10) {
	add(t1, ax, az)
	sub(t2, ax, az)
}

// montAdd calculate A = P + Q   where
//   X(A) = ax/az
//   X(P) = (t1+t2)/(t1-t2)
//   X(Q) = (t3+t4)/(t3-t4)
//   X(P-Q) = dx
//   clobbers t1 and t2, preserves t3 and t4  */
func montAdd(t1, t2, t3, t4, ax, az, dx *long10) {
	mul(ax, t2, t3)
	mul(az, t1, t4)
	add(t1, ax, az)
	sub(t2, ax, az)
	sqr(ax, t1)
	sqr(t1, t2)
	mul(az, t1, dx)
}

//montDbl calculate B = 2 * Q   where
//  X(B) = bx/bz
//  X(Q) = (t3+t4)/(t3-t4)
// clobbers t1 and t2, preserves t3 and t4
func montDbl(t1, t2, t3, t4, bx, bz *long10) {
	sqr(t1, t3)
	sqr(t2, t4)
	mul(bx, t1, t2)
	sub(t2, t1, t2)
	mulSmall(bz, t2, 121665)
	add(t1, t1, bz)
	mul(bz, t1, t2)
}

// numsize return the non empty byte size
func numsize(x *[64]byte, n int) int {
	for n--; n >= 0 && x[n] == 0; n-- {
		// do nothing
	}
	return n + 1
}

// Convert to internal format from little-endian byte format
func unpack(x *long10, m *[64]byte) {
	x._0 = int64((m[0] & 0xFF) | (m[1]&0xFF)<<8 | (m[2]&0xFF)<<16 | ((m[3]&0xFF)&3)<<24)
	x._1 = int64(((m[3]&0xFF)&^3)>>2 | (m[4]&0xFF)<<6 | (m[5]&0xFF)<<14 | ((m[6]&0xFF)&7)<<22)
	x._2 = int64(((m[6]&0xFF)&^7)>>3 | (m[7]&0xFF)<<5 | (m[8]&0xFF)<<13 | ((m[9]&0xFF)&31)<<21)
	x._3 = int64(((m[9]&0xFF)&^31)>>5 | (m[10]&0xFF)<<3 | (m[11]&0xFF)<<11 | ((m[12]&0xFF)&63)<<19)
	x._4 = int64(((m[12]&0xFF)&^63)>>6 | (m[13]&0xFF)<<2 | (m[14]&0xFF)<<10 | (m[15]&0xFF)<<18)
	x._5 = int64((m[16] & 0xFF) | (m[17]&0xFF)<<8 | (m[18]&0xFF)<<16 | ((m[19]&0xFF)&1)<<24)
	x._6 = int64(((m[19]&0xFF)&^1)>>1 | (m[20]&0xFF)<<7 | (m[21]&0xFF)<<15 | ((m[22]&0xFF)&7)<<23)
	x._7 = int64(((m[22]&0xFF)&^7)>>3 | (m[23]&0xFF)<<5 | (m[24]&0xFF)<<13 | ((m[25]&0xFF)&15)<<21)
	x._8 = int64(((m[25]&0xFF)&^15)>>4 | (m[26]&0xFF)<<4 | (m[27]&0xFF)<<12 | ((m[28]&0xFF)&63)<<20)
	x._9 = int64(((m[28]&0xFF)&^63)>>6 | (m[29]&0xFF)<<2 | (m[30]&0xFF)<<10 | (m[31]&0xFF)<<18)
}

// isOverflow Check if reduced-form input >= 2^255-19
func isOverflow(x *long10) bool {
	return ((x._0 > p26-19) &&
		(int64(x._1&x._3&x._5&x._7&x._9) == p25) &&
		(int64(x._2&x._4&x._6&x._8) == p26)) || (x._9 > p25)
}

// isNegative checks if x is "negative", requires reduced input
func isNegative(x *long10) int {
	if isOverflow(x) || x._9 < 0 {
		return int(1 ^ (x._0 & 1))
	}
	return int(0 ^ (x._0 & 1))
}

// pack Convert from internal format to little-endian byte format.  The
// number must be in a reduced form which is output by the following ops:
// unpack, mul, sqr
// set --  if input in range 0 .. p25
// If you're unsure if the number is reduced, first multiply it by 1.
func pack(x *long10, m *[64]byte) {
	var ld, ud, t int64

	if isOverflow(x) {
		ld = 1
	} else {
		ld = 0
	}

	if x._9 < 0 {
		ld--
	}

	ud = ld * -(p25 + 1)
	ld *= 19
	t = ld + x._0 + (x._1 << 26)
	m[0] = byte(t)
	m[1] = byte(t >> 8)
	m[2] = byte(t >> 16)
	m[3] = byte(t >> 24)
	t = (t >> 32) + (x._2 << 19)
	m[4] = byte(t)
	m[5] = byte(t >> 8)
	m[6] = byte(t >> 16)
	m[7] = byte(t >> 24)
	t = (t >> 32) + (x._3 << 13)
	m[8] = byte(t)
	m[9] = byte(t >> 8)
	m[10] = byte(t >> 16)
	m[11] = byte(t >> 24)
	t = (t >> 32) + (x._4 << 6)
	m[12] = byte(t)
	m[13] = byte(t >> 8)
	m[14] = byte(t >> 16)
	m[15] = byte(t >> 24)
	t = (t >> 32) + x._5 + (x._6 << 25)
	m[16] = byte(t)
	m[17] = byte(t >> 8)
	m[18] = byte(t >> 16)
	m[19] = byte(t >> 24)
	t = (t >> 32) + (x._7 << 19)
	m[20] = byte(t)
	m[21] = byte(t >> 8)
	m[22] = byte(t >> 16)
	m[23] = byte(t >> 24)
	t = (t >> 32) + (x._8 << 12)
	m[24] = byte(t)
	m[25] = byte(t >> 8)
	m[26] = byte(t >> 16)
	m[27] = byte(t >> 24)
	t = (t >> 32) + ((x._9 + ud) << 6)
	m[28] = byte(t)
	m[29] = byte(t >> 8)
	m[30] = byte(t >> 16)
	m[31] = byte(t >> 24)
}

// cpy Copy a number
func cpy(out, in *long10) {
	out._0 = in._0
	out._1 = in._1
	out._2 = in._2
	out._3 = in._3
	out._4 = in._4
	out._5 = in._5
	out._6 = in._6
	out._7 = in._7
	out._8 = in._8
	out._9 = in._9
}

// set Set a number to value, which must be in range -185861411 .. 185861411
func set(out *long10, in int64) {
	out._0 = in
	out._1 = 0
	out._2 = 0
	out._3 = 0
	out._4 = 0
	out._5 = 0
	out._6 = 0
	out._7 = 0
	out._8 = 0
	out._9 = 0
}

// core calculate P = kG   and  s = sign(P)/k
func signercore(Px, s, k, Gx *[64]byte) {
	dx := &long10{}
	t1 := &long10{}
	t2 := &long10{}
	t3 := &long10{}
	t4 := &long10{}

	x := []*long10{{}, {}}
	z := []*long10{{}, {}}

	/* unpack the base */
	if Gx != nil {
		unpack(dx, Gx)
	} else {
		set(dx, 9)
	}

	/* 0G = point-at-infinity */
	set(x[0], 1)
	set(z[0], 0)

	/* 1G = G */
	cpy(x[1], dx)
	set(z[1], 1)

	for i := 31; i >= 0; i-- {
		//if i == 0 {
		// i = 0
		//}
		for j := 7; j >= 0; j-- {
			/* swap arguments depending on bit */
			bit1 := (k[i] & 0xFF) >> uint(j) & 1
			bit0 := ^(k[i] & 0xFF) >> uint(j) & 1
			ax := x[bit0]
			az := z[bit0]
			bx := x[bit1]
			bz := z[bit1]

			/* a' = a + b	*/
			/* b' = 2 b	*/
			montPrep(t1, t2, ax, az)
			montPrep(t3, t4, bx, bz)
			montAdd(t1, t2, t3, t4, ax, az, dx)
			montDbl(t1, t2, t3, t4, bx, bz)
		}
	}

	recip(t1, z[0], 0)
	mul(dx, x[0], t1)
	pack(dx, Px)

	/* calculate s such that s abs(P) = G  .. assumes G is std base point */
	if s != nil {
		xToY2(t2, t1, dx)        /* t1 = Py^2  */
		recip(t3, z[1], 0)       /* where Q=P+G ... */
		mul(t2, x[1], t3)        /* t2 = Qx  */
		add(t2, t2, dx)          /* t2 = Qx + Px  */
		t2._0 += 9 + 486662      /* t2 = Qx + Px + Gx + 486662  */
		dx._0 -= 9               /* dx = Px - Gx  */
		sqr(t3, dx)              /* t3 = (Px - Gx)^2  */
		mul(dx, t2, t3)          /* dx = t2 (Px - Gx)^2  */
		sub(dx, dx, t1)          /* dx = t2 (Px - Gx)^2 - Py^2  */
		dx._0 -= 39420360        /* dx = t2 (Px - Gx)^2 - Py^2 - Gy^2  */
		mul(t1, dx, &baseR2Y)    /* t1 = -Py  */
		if isNegative(t1) != 0 { /* sign is 1, so just copy  */
			copy(s[:32], k[:32])
		} else { /* sign is -1, so negate  */
			mulaSmall(s, &orderTimes8, 0, k, 32, -1)
		}

		// (Qx + Px + Gx + 486662)(Px - Gx)^2 - Py^2 - Gy^2

		/* reduce s mod q
		 * (is this needed?  do it just in case, it's fast anyway) */
		//divmod((dstptr) t1, s, 32, order25519, 32);

		/* take reciprocal of s mod q */
		temp2, temp3 := &[64]byte{}, &[64]byte{}

		orderCpy := order
		copy(s[:32], egcd32(temp2, temp3, s, &orderCpy)[:32])
		if (s[31] & 0x80) != 0 {
			mulaSmall(s, s, 0, &order, 32, 1)
		}
	}
}

//sign Signature generation primitive, calculates (x-h)s mod q
//   v  [out] signature value
//   h  [in]  signature hash (of message, signature pub key, and context data)
//   x  [in]  signature private key
//   s  [in]  private key for signing
// returns true on success, false on failure (use different x or h)
func Sign(v, h, x, s *[64]byte) bool {
	// v = (x - h) s  mod q
	var w, i int
	h1, x1 := [64]byte{}, [64]byte{}
	tmp1, tmp2, tmp3 := [64]byte{}, [64]byte{}, [64]byte{}

	copy(h1[:], h[:])
	copy(x1[:], x[:])

	// Reduce modulo group order
	divmod(&tmp3, &h1, 32, &order, 32)
	divmod(&tmp3, &x1, 32, &order, 32)

	// v = x1 - h1
	// If v is negative, add the group order to it to become positive.
	// If v was already positive we don't have to worry about overflow
	// when adding the order because v < order and 2*order < 2^256
	mulaSmall(v, &x1, 0, &h1, 32, -1)
	mulaSmall(v, v, 0, &order, 32, 1)

	// tmp1 = (x-h)*s mod q
	mula32(&tmp1, v, s, 32, 1)
	divmod(&tmp2, &tmp1, 64, &order, 32)

	for w, i = 0, 0; i < 32; i++ {
		v[i] = tmp1[i]
		w |= int(v[i])
	}
	return w != 0
}

//keygen Key-pair generation
//   P  [out] your public key
//   s  [out] your private key for signing
//   k  [out] your private key for key agreement
//   k  [in]  32 random bytes
// s may be NULL if you don't care
//
// WARNING: if s is not NULL, this function has data-dependent timing
func Keygen(P, s, k *[64]byte) {
	clamp(k)
	signercore(P, s, k, nil)
}

//clamp Private key clamping
//   k [out] your private key for key agreement
//   k  [in]  32 random bytes
func clamp(k *[64]byte) {
	k[31] &= 0x7F //0111 1111
	k[31] |= 0x40 //0100 0000
	k[0] &= 0xF8  //1111 1000
}

// SignBytes signs the message with secretPhrase and returns a signature
func signBytes(message []byte, secretPhrase string) (signature []byte, err error) {
	P, s := [64]byte{}, [64]byte{}
	sha := sha256.New()

	var digest, m, x, Y, h, v [64]byte

	if _, err = sha.Write([]byte(secretPhrase)); err != nil {
		return
	}
	copy(digest[:32], sha.Sum(nil))
	Keygen(&P, &s, &digest)

	sha.Reset()
	if _, err = sha.Write(message); err != nil {
		return
	}
	copy(m[:], sha.Sum(nil))

	sha.Reset()
	if _, err = sha.Write(m[:32]); err != nil {
		return
	}
	if _, err = sha.Write(s[:32]); err != nil {
		return
	}
	copy(x[:], sha.Sum(nil))

	Keygen(&Y, nil, &x)

	sha.Reset()
	if _, err = sha.Write(m[:32]); err != nil {
		return
	}
	if _, err = sha.Write(Y[:32]); err != nil {
		return
	}
	copy(h[:], sha.Sum(nil))

	if ok := Sign(&v, &h, &x, &s); !ok {
		panic("Signature generation failed")
	}

	signature = make([]byte, 64)
	copy(signature[:32], v[:32])
	copy(signature[32:], h[:32])

	return
}

// SignBytes signs the unsignedTransactionString with secretPhrase and returns a signature
func SignBytes(unsignedTransactionBytes []byte, secretPhrase string) (signedTransactionBytes []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	signature, err := signBytes(unsignedTransactionBytes, secretPhrase)

	if err != nil {
		// err = errors.Wrap(err, "error occurs during generate signature")
		return
	}

	copy(unsignedTransactionBytes[nxtSignatureOffset:nxtSignatureEnd], signature)
	signedTransactionBytes = unsignedTransactionBytes
	return
}

// SignString signs the unsignedTransactionJSON with secretPhrase and returns a signature
func SignString(unsignedTransactionBytes string, secretPhrase string) (signedTransactionBytes string, err error) {
	transactionBytes, err := hex.DecodeString(unsignedTransactionBytes)
	if err != nil {
		// err = errors.Wrap(err, "error occurs during decode hex string")
		return
	}
	signedBytes, err := SignBytes(transactionBytes, secretPhrase)
	signedTransactionBytes = hex.EncodeToString(signedBytes)
	return
}
