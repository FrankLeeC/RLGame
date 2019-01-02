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
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func read() *[9][9][3]int {
	f, err := os.Open("./policy.txt")
	if err != nil {
		fmt.Println("err: ", err)
		return nil
	}
	defer f.Close()
	var policy [9][9][3]int
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		kv := strings.Split(line, ":")
		a, _ := strconv.Atoi(strings.TrimSpace(kv[1]))
		state := strings.Split(kv[0], ",")
		s0, _ := strconv.Atoi(strings.TrimSpace(state[0]))
		s1, _ := strconv.Atoi(strings.TrimSpace(state[1]))
		s, _ := strconv.Atoi(strings.TrimSpace(state[2]))
		policy[s0][s1][s] = a
	}
	return &policy
}

func randExcept(a int) int {
	r := rand.Intn(2)
	if a == 0 {
		return r + 1
	}
	if a == 2 {
		return r
	}
	if r == 1 {
		return 2
	}
	return 0
}

func randomStart() *[3]int {
	a := rand.Intn(3)
	b := randExcept(a)
	c := rand.Intn(3)
	return &[3]int{a, b, c}
}

func stepPlay(state *[3]int, policy *[9][9][3]int) int {
	action := policy[state[0]][state[1]][state[2]]
	state[2] += action
	if state[2] > 2 {
		state[2] = 2
	}
	if state[2] < 0 {
		state[2] = 0
	}
	state[0] += 3
	state[1] += 3
	if state[0]-6 == state[2] || state[1]-6 == state[2] {
		c0 := false
		if state[0] > 8 {
			c0 = true
			state[0] = rand.Intn(3)
		}
		if state[1] > 8 {
			if c0 {
				state[1] = randExcept(state[0])
			} else {
				state[1] = rand.Intn(3)
			}
		}
		return -200
	}
	if state[0] <= 8 && state[1] <= 8 {
		p := 0
		if state[0] >= 6 {
			p += 30
		}
		if state[1] >= 6 {
			p += 30
		}
		return p
	}
	p := 0
	c0 := false
	if state[0] > 8 {
		c0 = true
		state[0] = rand.Intn(3)
		p += 50
	}
	if state[1] > 8 {
		p += 50
		if c0 {
			state[1] = randExcept(state[0])
		} else {
			state[1] = rand.Intn(3)
		}
	}
	if state[0] > state[1] {
		state[0], state[1] = state[1], state[0]
	}
	return p
}

func clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func print(state *[3]int) {
	s0 := state[0]
	s1 := state[1]
	s2 := state[2]
	var s strings.Builder
	s.WriteString(" ----- ----- ----- \n")
	if s0 == 0 || s1 == 0 {
		s.WriteString("|  O  |")
	} else {
		s.WriteString("|     |")
	}
	if s0 == 1 || s1 == 1 {
		s.WriteString("  O  |")
	} else {
		s.WriteString("     |")
	}
	if s0 == 2 || s1 == 2 {
		s.WriteString("  O  |\n")
	} else {
		s.WriteString("     |\n")
	}
	s.WriteString(" ----- ----- ----- \n")
	if s0 == 3 || s1 == 3 {
		s.WriteString("|  O  |")
	} else {
		s.WriteString("|     |")
	}
	if s0 == 4 || s1 == 4 {
		s.WriteString("  O  |")
	} else {
		s.WriteString("     |")
	}
	if s0 == 5 || s1 == 5 {
		s.WriteString("  O  |\n")
	} else {
		s.WriteString("     |\n")
	}
	s.WriteString(" ----- ----- ----- \n")
	if s0 == 6 || s1 == 6 {
		s.WriteString("|  O  |")
	} else {
		if s2 == 0 {
			s.WriteString("|  X  |")
		} else {
			s.WriteString("|     |")
		}
	}
	if s0 == 7 || s1 == 7 {
		s.WriteString("  O  |")
	} else {
		if s2 == 1 {
			s.WriteString("  X  |")
		} else {
			s.WriteString("     |")
		}
	}
	if s0 == 8 || s1 == 8 {
		s.WriteString("  O  |\n")
	} else {
		if s2 == 2 {
			s.WriteString("  X  |\n")
		} else {
			s.WriteString("     |\n")
		}
	}
	s.WriteString(" ----- ----- ----- \n")
	fmt.Println(s.String())
}

func play() {
	rand.Seed(time.Now().UnixNano())
	goal := 1000
	policy := read()
	state := randomStart()
	point := 0
	clear()
	print(state)
	for point < goal {
		<-time.After(time.Duration(1) * time.Second)
		clear()
		point += stepPlay(state, policy)
		if point < 0 {
			point = 0
		}
		print(state)
		fmt.Println("point:", point)
	}
}
