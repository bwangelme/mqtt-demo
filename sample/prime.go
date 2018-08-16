package main

import "fmt"

func counter(c chan int) {
	i := 2
	for {
		c <- i
		i++
	}
}

/*
 * FilterPrime 将素数`prime`的倍数过滤掉，输出的值一定不是`prime`的倍数
 *
 * 图示说明: https://passage-1253400711.cos.ap-beijing.myqcloud.com/2018-08-13-150647.png
 */
func FilterPrime(prime int, listen, send chan int) {
	var i int

	for {
		i = <-listen
		if i%prime != 0 {
			send <- i
		}
	}
}

func sieve() (prime chan int) {
	c := make(chan int)
	go counter(c)

	prime = make(chan int)
	go func() {
		var p int
		var newc chan int

		for {
			p = <-c
			prime <- p
			newc = make(chan int)
			go FilterPrime(p, c, newc)
			c = newc
		}
	}()

	return prime
}

func main() {

	// 这种方法计算素数，每产生一个素数，就需要新建一个goroutine。
	// 例如求前N个素数，空间复杂度为O(N)，时间复杂度为O(M)，M表示素数N的大小
	prime := sieve()
	const N = 100

	var times [N][0]int
	for range times {
		p := <-prime
		fmt.Println(p)
	}
}
