/*
MIT License

Copyright (c) 2018 Frank Lee

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	up    = 0
	right = 1
	down  = 2
	left  = 3
	gamma = 0.9
)
var transition = readGrid()

func readGrid() *map[int][4]int {
	f, _ := os.Open("./grid.txt")
	defer f.Close()
	m := make(map[int][4]int)
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if line != "" {
			data := strings.Split(line, ":")
			state, _ := strconv.Atoi(strings.TrimSpace(data[0]))
			next := strings.Split(data[1], ",")
			s1, _ := strconv.Atoi(strings.TrimSpace(next[0]))
			s2, _ := strconv.Atoi(strings.TrimSpace(next[1]))
			s3, _ := strconv.Atoi(strings.TrimSpace(next[2]))
			s4, _ := strconv.Atoi(strings.TrimSpace(next[3]))
			m[state] = [4]int{s1, s2, s3, s4}
		}
	}
	return &m
}

func initPolicy() *[100]int {
	var policy [100]int
	for i := 0; i < 100; i++ {
		policy[i] = up
	}
	return &policy
}

func initValue(target int) *[100]float64 {
	var value [100]float64
	for i := 0; i < 100; i++ {
		value[i] = rand.Float64()
	}
	value[target] = 0.0
	return &value
}

func getReward(state, target int) float64 {
	if state == target {
		return 1.0
	}
	return 0.0
}

func evaluate(value *[100]float64, policy *[100]int, target int) {
	epsilon := 0.0001
	c := 0
	start := time.Now()
	for {
		c++
		m := 0.0
		for i := 0; i < 100; i++ {
			ns := (*transition)[i]
			a := policy[i]
			s := ns[a]
			v := value[i]
			nv := getReward(i, target) + gamma*value[s]
			value[i] = nv
			m = math.Max(m, math.Abs(nv-v))
		}
		current := time.Now()
		cost := current.Sub(start)
		start = current
		fmt.Printf("%d max change %f in evaluation cost: %fs\n", c, m, cost.Seconds())
		if m < epsilon {
			break
		}
	}
}

func improve(value *[100]float64, policy *[100]int, target int) bool {
	changed := false
	c := 0
	start := time.Now()
	for i := 0; i < 100; i++ {
		ra := policy[i]
		ma := 0
		mv := 0.0
		for j := 0; j < 4; j++ {
			ns := (*transition)[i][j]
			v := value[ns]
			if v > mv {
				mv = v
				ma = j
			}
		}
		if ra != ma {
			c++
			changed = true
			policy[i] = ma
		}
		current := time.Now()
		cost := current.Sub(start)
		start = current
		fmt.Printf("%d changes in improvement cost: %fs\n", c, cost.Seconds())
	}
	return changed
}

func save(target int, policy *[100]int) {
	f, _ := os.Create("./policy/policy_" + strconv.Itoa(target) + ".txt")
	for i := 0; i < 100; i++ {
		f.WriteString(fmt.Sprintf("%d: %d\n", i, policy[i]))
	}
	f.Close()
	f, _ = os.Create("./policy/gragh_" + strconv.Itoa(target) + ".txt")
	var s [10][10]string
	for i := 0; i < 100; i++ {
		x, y := getLocation(i)
		a := getAction(policy[i])
		s[x][y] = a
	}
	var b strings.Builder
	for i := 0; i < 10; i++ {
		b.WriteString(" ----- ")
	}
	b.WriteString("\n")
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			b.WriteString("|  " + formatNumber(i*10+j) + "" + s[i][j] + "")
		}
		b.WriteString("|\n")
		for i := 0; i < 10; i++ {
			b.WriteString(" ----- ")
		}
		b.WriteString("\n")
	}
	f.WriteString(b.String())
	f.Close()
}

func formatNumber(i int) string {
	if i < 10 {
		return "0" + strconv.Itoa(i)
	}
	return strconv.Itoa(i)
}

func getLocation(i int) (int, int) {
	return i / 10, i % 10
}

func getAction(i int) string {
	if i == 0 {
		return "⬆️"
	}
	if i == 1 {
		return "➡️"
	}
	if i == 2 {
		return "⬇️"
	}
	return "⬅️"
}

func run(target int) {
	run := true
	policy := initPolicy()
	value := initValue(target)
	c := 0
	for run {
		c++
		evaluate(value, policy, target)
		fmt.Println("evaluate")
		run = improve(value, policy, target)
		fmt.Println("improve")
		fmt.Println("-----------------", c, "-----------------")
	}
	save(target, policy)
}

func main() {
	for i := 0; i < 100; i++ {
		run(i)
	}
}
