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
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

func initPolicy() *[9][9][3]int {
	var policy [9][9][3]int
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			for k := 0; k < 3; k++ {
				policy[i][j][k] = 0
			}
		}
	}
	return &policy
}

func initValue() *[9][9][3]float64 {
	var value [9][9][3]float64
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			for k := 0; k < 3; k++ {
				if !isValid(&[3]int{i, j, k}) {
					value[i][j][k] = -1.0
				} else {
					value[i][j][k] = rand.Float64()
				}
			}
		}
	}
	return &value
}

func getAllAction() *[]int {
	return &[]int{-1, 0, 1}
}

func sort(i, j int) (int, int) {
	if i > j {
		return j, i
	}
	return i, j
}

func step(value *[9][9][3]float64, state *[3]int, action int) float64 {
	gamma := 0.9
	s := state[2] + action
	if s > 2 {
		s = 2
	}
	if s < 0 {
		s = 0
	}
	s0 := state[0] + 3
	s1 := state[1] + 3
	if s0-6 == s || s1-6 == s {
		reward := -200.0
		if s0-6 == s {
			if s1 > 8 {
				return (reward + gamma*value[0][s0][s] + reward + gamma*value[1][s0][s] + reward + gamma*value[2][s0][s]) / 3.0
			}
			a, b := sort(s0, s1)
			return reward + gamma*value[a][b][s]
		}
		if s0 > 8 {
			return (reward + gamma*value[0][s1][s] + reward + gamma*value[1][s1][s] + reward + gamma*value[2][s1][s]) / 3.0
		}
		a, b := sort(s0, s1)
		return reward + gamma*value[a][b][s]
	}
	if s0 <= 8 && s1 <= 8 {
		reward := 0.0
		if s0-6 >= 0 && s0-6 <= 2 {
			reward += 30
		}
		if s1-6 >= 0 && s1-6 <= 2 {
			reward += 30
		}
		a, b := sort(s0, s1)
		return 1.0 * (reward + gamma*value[a][b][s])
	}
	if s0 > 8 && s1 > 8 {
		reward := 50.0 * 2.0
		return (reward + gamma*value[0][1][s] + reward + gamma*value[0][2][s] + reward + gamma*value[1][2][s]) / 3.0
	}
	reward := 50.0
	if s0 > 8 {
		if s1 >= 6 {
			reward += 30.0
		}
		return (reward + gamma*value[0][s1][s] + reward + gamma*value[1][s1][s] + reward + gamma*value[2][s1][s]) / 3.0
	}
	if s0 >= 6 {
		reward += 30.0
	}
	return (reward + gamma*value[0][s0][s] + reward + gamma*value[1][s0][s] + reward + gamma*value[2][s0][s]) / 3.0

}

func getPolicyAction(state *[3]int, policy *[9][9][3]int) int {
	return policy[state[0]][state[1]][state[2]]
}

func evaluation(value *[9][9][3]float64, policy *[9][9][3]int) {
	epsilon := 0.00001
	c := 0
	start := time.Now()
	for {
		c++
		m := 0.0
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if j <= i {
					continue
				}
				for k := 0; k < 3; k++ {
					if !isValid(&[3]int{i, j, k}) {
						continue
					}
					v := value[i][j][k]
					a := policy[i][j][k]
					tmp := step(value, &[3]int{i, j, k}, a)
					value[i][j][k] = tmp
					m = math.Max(m, math.Abs(tmp-v))
				}
			}
		}
		current := time.Now()
		cost := current.Sub(start)
		start = current
		fmt.Printf("%d maximum %f in evaluation cost: %fs\n", c, m, cost.Seconds())
		if m < epsilon {
			break
		}
	}
}

func isValid(state *[3]int) bool {
	s0 := state[0]
	s1 := state[1]
	s := state[2]
	if s0 == s1 {
		return false
	}
	if s0-6 == s || s1-6 == s {
		return false
	}
	return true
}

func improvement(value *[9][9][3]float64, policy *[9][9][3]int) bool {
	changed := false
	c := 0
	start := time.Now()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if j <= i {
				continue
			}
			for k := 0; k < 3; k++ {
				if !isValid(&[3]int{i, j, k}) {
					continue
				}
				actions := getAllAction()
				m := 0.0
				ra := getPolicyAction(&[3]int{i, j, k}, policy)
				ma := ra
				for _, a := range *actions {
					v := step(value, &[3]int{i, j, k}, a)
					if v > m {
						m = v
						ma = a
					}
				}
				if ra != ma {
					changed = true
					c++
					policy[i][j][k] = ma
				}
			}
			current := time.Now()
			cost := current.Sub(start)
			start = current
			fmt.Printf("%d changes in evaluation cost: %fs\n", c, cost.Seconds())
		}
	}
	return changed
}

func output(policy *[9][9][3]int) {
	f, err := os.Create("./policy.txt")
	if err != nil {
		fmt.Println("output err: ", err)
		return
	}
	defer f.Close()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if j <= i {
				continue
			}
			for k := 0; k < 3; k++ {
				if !isValid(&[3]int{i, j, k}) {
					continue
				}
				f.WriteString(fmt.Sprintf("%d, %d, %d : %d\n", i, j, k, policy[i][j][k]))
			}
		}
	}
}

func run() *[9][9][3]int {
	policy := initPolicy()
	value := initValue()
	run := true
	c := 0
	for run {
		c++
		evaluation(value, policy)
		run = improvement(value, policy)
		fmt.Println("----------------------------------", c, "----------------------------------")
	}
	println(value[3][4][0], value[3][4][1], value[3][4][2])
	return policy
}

func train() {
	policy := run()
	output(policy)
}
