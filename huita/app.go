package main

import (
	a "app_fnc"
	"fmt"
	//"math"
	"math/big"
	"math/rand"
	"time"
)

func init() {
	InitLog()
}

func main() {
	exectime("test1", func() {

		fromS0 := "0123456789"
		fromS := "0123456789"
		toS := "0123456789abcdefghijklmnopqrstuvwxyz"
		toS = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-+_*=|\\/[]{}()<>,.?!;:$#@%^&~№"
		toS = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
		fromS = "01234567"
		toS = "0123456789"

		c0 := NewByteSeqConverterStr(fromS0, fromS)
		c1 := NewByteSeqConverterStr(fromS, toS)
		c2 := NewByteSeqConverterStr(toS, fromS)

		for i := 0; i < 10; i++ {
			n0 := a.Itoa(rand.Intn(1000000))
			n0 = c0.ConvertStr(n0)
			n1 := c1.ConvertStr(n0)
			n2 := c2.ConvertStr(n1)

			ok := 0
			if n0 == n2 {
				ok = 1
			}
			fmt.Printf("%d  %s -> %s -> %s\n", ok, n0, n1, n2)
		}

		{
			bi := big.NewInt(0)
			for i := 0; i < 36; i++ {
				bi.SetString(a.Itoa(i), 10)
				//fmt.Printf("%d -> %s\n", i, bi.Text(36))
				fmt.Printf("%s", bi.Text(36))
			}
			fmt.Printf("\n\n\n%s", bi.Text(36))
		}

	})
}

func exectime(name string, f func()) {
	LogPrint(name)
	t0 := time.Now()
	f()
	t1 := time.Now()
	LogPrint(a.Sprintf("%s time: %v\nend\n", name, t1.Sub(t0)))
}

//
//
//
//
//

type byteSeqConvert struct {
	fromSeq []uint8
	toSeq   []uint8

	fromSymbToIdx []uint8
	toSymbToIdx   []uint8

	fromSeqLen int
	toSeqLen   int

	fromBits uint8
	toBits   uint8
}

var bigintsymbols []byte
var bigintsymbolsidx []uint8

func NewByteSeqConverterStr(fromSeq, toSeq string) *byteSeqConvert {
	return NewByteSeqConverter([]byte(fromSeq), []byte(toSeq))
}

func (p *byteSeqConvert) getBitsForLen(l int) uint8 {
	switch {
	case l <= 2:
		return 1
	case l > 2 && l <= 4:
		return 2
	case l > 4 && l <= 8:
		return 3
	case l > 8 && l <= 16:
		return 4
	case l > 16 && l <= 32:
		return 5
	case l > 32 && l <= 64:
		return 6
	case l > 64 && l <= 128:
		return 7
	case l > 128:
		return 8
	}
	return 0
}

func NewByteSeqConverter(fromSeq, toSeq []byte) *byteSeqConvert {
	p := new(byteSeqConvert)
	p.fromSeq = fromSeq
	p.toSeq = toSeq
	p.fromSeqLen = len(fromSeq)
	p.toSeqLen = len(toSeq)

	p.fromBits = p.getBitsForLen(p.fromSeqLen)
	p.toBits = p.getBitsForLen(p.toSeqLen)

	p.fromSymbToIdx = make([]uint8, 256)
	for i := 0; i < p.fromSeqLen; i++ {
		smb := uint8(fromSeq[i])
		//LogPrint(a.Sprintf("%d smb: %d", i, smb))
		p.fromSymbToIdx[smb] = uint8(i)
	}

	p.toSymbToIdx = make([]uint8, 256)
	for i := 0; i < p.toSeqLen; i++ {
		smb := uint8(toSeq[i])
		p.toSymbToIdx[smb] = uint8(i)
	}

	if bigintsymbolsidx == nil && p.fromSeqLen != p.toSeqLen && (p.toSeqLen <= 36 || p.fromSeqLen <= 36) {
		bigintsymbols = []byte("0123456789abcdefghijklmnopqrstuvwxyz") //len 36
		bigintsymbolsidx = make([]uint8, 128)
		for i := 0; i < 36; i++ {
			smb := uint8(bigintsymbols[i])
			bigintsymbolsidx[smb] = uint8(i)
		}
	}
	return p
}

func (p *byteSeqConvert) ConvertStr(number string) string {
	return string(p.Convert([]byte(number)))
}

func (p *byteSeqConvert) Convert(number []byte) []byte {
	if p.fromSeqLen == p.toSeqLen {
		newlen := len(number)
		newnum := make([]byte, newlen)
		for i := 0; i < newlen; i++ {
			smb := uint8(number[i])
			idxfrom := p.fromSymbToIdx[smb]
			newnum[i] = p.toSeq[idxfrom]
		}
		return newnum
	}

	if p.fromSeqLen <= 36 && p.toSeqLen <= 36 && p.toSeqLen == -999 {
		bi := big.NewInt(0)
		//fmt.Printf("(0:%s)", string(number))
		//переводим строку из кодировки p.fromSeq в bigintsymbols
		{
			newlen := len(number)
			tmpnum := make([]byte, newlen)
			for i := 0; i < newlen; i++ {
				smb := uint8(number[i])
				idxfrom := p.fromSymbToIdx[smb]
				tmpnum[i] = bigintsymbols[idxfrom]
			}
			//fmt.Printf("(1:%s)", string(tmpnum))
			//полученную строку переводим в big.Int
			bi.SetString(string(tmpnum), p.fromSeqLen)
		}

		//далее переводим big.Int в строку с кодировкой bigintsymbols в base==p.toSeqLen
		newnum := []byte(bi.Text(p.toSeqLen))
		newlen := len(newnum)

		//fmt.Printf("(2:%s)", string(newnum))
		//полученную строку переводим из кодировки bigintsymbols в p.toSeq
		for i := 0; i < newlen; i++ {
			smb := uint8(newnum[i])
			idxfrom := bigintsymbolsidx[smb]
			newnum[i] = p.toSeq[idxfrom]
		}
		//fmt.Printf("(3:%s)", string(newnum))
		return newnum
	}

	//if p.fromSeqLen > p.toSeqLen
	{
		//переводим из большей системы счисления в меньшую
		//для этого вначале переводим большую систему счисления в 256 ричную систему счисления(сс)
		//  для этого собираем биты из большей системы счисления в один массив байт
		//    байты задаем в обратном порядке - с конца!
		//полученый массив байт переводим в экземпляр big.Int который предоставляет удобные math операции
		//далее big.Int переводим в нужную нам сс
		lenn := len(number)
		lent := (lenn*int(p.fromBits))/8 + 1
		bt := make([]byte, lent)

		//bt - байты задаем с конца
		//биты в байтах задаем с начала
		var curbit uint = 0
		curbyte := lent - 1
		setbits := uint(p.fromBits)
		for i := lenn - 1; i >= 0; i-- {
			//for i := 0; i < lenn; i++ {
			smb := uint8(number[i])
			ni := p.fromSymbToIdx[smb]

			var ib uint = 0
			for ; ib < setbits; ib++ {
				if curbit >= 8 {
					curbit = 0
					curbyte--
				}
				if hasBit(ni, ib) {
					bt[curbyte] = setBit(bt[curbyte], curbit)
				}
				curbit++
			}

			s1 := big.NewInt(int64(ni)).Text(2)
			s2 := big.NewInt(0).SetBytes(bt).Text(2)
			s3 := big.NewInt(0).SetBytes(bt).Text(p.fromSeqLen)
			fmt.Printf("%s i:%d ni:%d nib:%s  bnb:%s  bn:%s\n", string(number), i, ni, s1, s2, s3)
		}

		//в bt массив байт с 256ричной сс из p.fromSeqLen сс
		//создаем bn(big.Int) из массива байт bt
		bn := big.NewInt(0).SetBytes(bt)

		fmt.Printf("[%d->%d(bits:%d->%d)][%s/%s->%s]\n\n", p.fromSeqLen, p.toSeqLen, p.fromBits, p.toBits,
			string(number), bn.Text(p.fromSeqLen), bn.Text(p.toSeqLen))

		{ //переводим bn из 256 ричной сс в p.toSeqLen сс
			newlen := (lenn*int(p.fromBits))/int(p.toBits) + 1
			newnum := make([]byte, newlen)

			dv := big.NewInt(int64(p.toSeqLen))
			mod := big.NewInt(0)

			j := newlen
			for {
				j--
				bn, mod = bn.DivMod(bn, dv, mod)
				nj := mod.Uint64() // остаток от деления (bn%dv)
				smb := p.toSeq[nj]
				newnum[j] = smb
				fmt.Printf("%s j:%d smb:%v  bn:%d  dv:%d  mod:%d \n", string(number), j, smb, bn, dv, mod)

				if bn.Cmp(dv) == -1 { // bn < dv
					j--
					nj = bn.Uint64()
					newnum[j] = p.toSeq[nj]
					break
				}
			}
			fmt.Printf("[%d->%d(bits:%d->%d)][%s->%s]\n\n", p.fromSeqLen, p.toSeqLen, p.fromBits, p.toBits,
				string(number), string(newnum))

			fmt.Printf("---------------------------------------------------------\n\n")
			return newnum
		}
	}
	return nil

	//if p.fromSeqLen < p.toSeqLen
	{
		lenn := int64(len(number))
		//newlen := len(number)*math.Ceil(p.toSeqLen/p.fromSeqLen) + 1
		//newnum := make([]byte, newlen)

		b := big.NewInt(0)

		var i int64 = 0
		for ; i < lenn; i++ {
			smb := uint8(number[i])
			ni := p.fromSymbToIdx[smb]
			pn := int64(ni) * pow(int64(p.fromSeqLen), lenn-i)
			b.Add(b, big.NewInt(pn))
		}
		//fmt.Printf("[%s->%s]", number, b.Text(36))

		//далее переводим из 256сс в p.toSeqLen сс
		bb := b.Bytes()
		lenb := int64(len(bb))

		//fmt.Printf("[%d:%s]", lenb, bb)

		var ptoSeqLen uint8 = uint8(p.toSeqLen)
		newlen := int(lenb)*(256/int(ptoSeqLen)) + 1
		newnum := make([]byte, newlen)

		j := 0
		for i = 0; i < lenb; i++ {
			ni := uint8(bb[i])
			for {
				if ni < ptoSeqLen {
					newnum[j] = p.toSeq[ni]
					j++
					break
				}
				nj := ni % ptoSeqLen
				ni /= ptoSeqLen
				newnum[j] = p.toSeq[nj]
				j++
			}
		}
		//fmt.Printf("[%d/%d:%s]", newlen, j, newnum)

		/********
		to := p.toSeqLen
		if to > 36 {
			to = 36
		}
		fmt.Printf("(p.toSeqLen: %d -> %d)", p.toSeqLen, to)
		r := []byte(b.Text(to))
		*********/
		return newnum[:j]
	}

	return nil
}

func (p *byteSeqConvert) setBitsBigIntFromSeq(z *big.Int, fromBit int, n uint8) {
	bn := big.NewInt(int64(n) << uint(fromBit))

	for i := 0; i < int(p.fromBits); i++ {
		z.SetBit(bn, i+fromBit, bn.Bit(i))
	}
}

func (p *byteSeqConvert) setBitsBytesFromSeq(x []byte, fromBit int, n uint8) {
	setbits := uint(p.fromBits)
	nbyte := fromBit / 8
	nbit := uint(fromBit % 8)
	var i uint = 0
	for ; i < setbits; i++ {
		if nbit >= 8 {
			nbit = 0
			nbyte++
		}
		if hasBit(n, i) {
			x[nbyte] = setBit(x[nbyte], nbit)
		}
		nbit++
	}
}

//----------------------------------------------------------
// Sets the bit at pos in the integer n.
func setBit(n uint8, pos uint) uint8 {
	n |= (1 << pos)
	return n
}

// Clears the bit at pos in n.
func clearBit(n uint8, pos uint) uint8 {
	mask := uint8(^(1 << pos))
	n &= mask
	return n
}
func hasBit(n uint8, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

//----------------------------------------------------------

func pow(n int64, pw int64) int64 {
	if pw == 0 {
		return 1
	}
	if pw == 1 {
		return n
	}
	pn := n
	for pw > 1 {
		pw--
		pn *= n
	}
	return pn
}

func powBigInt(n *big.Int, pw *big.Int) {
	//http://stackoverflow.com/questions/29912249/what-is-the-equivalent-for-bigint-powa-in-go
	n.Exp(n, pw, nil)
}
