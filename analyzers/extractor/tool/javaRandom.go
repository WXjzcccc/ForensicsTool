package tool

import "github.com/donnie4w/go-logger/logger"

// JavaRandom 模拟Java的Random类
type JavaRandom struct {
	seed int64
}

func NewJavaRandom(seed int64) *JavaRandom {
	return &JavaRandom{seed: (seed ^ 0x5DEECE66D) & ((1 << 48) - 1)}
}

func (r *JavaRandom) next(bits int) int32 {
	r.seed = (r.seed*0x5DEECE66D + 0xB) & ((1 << 48) - 1)
	return int32(r.seed >> (48 - bits))
}

func (r *JavaRandom) NextInt64() int64 {
	return int64(r.next(32))<<32 + int64(r.next(32))
}

func (r *JavaRandom) NextInt(n int) int {
	if n <= 0 {
		logger.Errorf("n must be positive")
		panic("n must be positive")
	}

	// 等同于Java的Random.nextInt(n)实现
	if (n & -n) == n { // n是2的幂
		return int((int64(n) * int64(r.next(31))) >> 31)
	}

	var bits, val int32
	for {
		bits = r.next(31)
		val = bits % int32(n)
		if bits-val+(int32(n)-1) >= 0 {
			break
		}
	}
	return int(val)
}
