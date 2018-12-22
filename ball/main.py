# -*- coding: utf-8 -*-

'''
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
'''

import numpy as np
import random

COUNT = np.zeros([25, 25, 25, 25, 5, 3])
POLICY = np.zeros([25, 25, 25, 25, 5, 3])  # 第1，2，3，4个球的位置，自己的位置  -1: left   1: right  0: hold
Q = np.zeros([25, 25, 25, 25, 5, 3])
LEFT = -1
HOLD = 0
RIGHT = 1


class State:

    def __init__(self, l):
        self.l = l
        self.p = 0

    def get(self, idx):
        if idx < len(self.l):
            return self.l[idx]
        return None

    def set_point(self, p):
        self.p = p

    def point(self):
        return self.p

def get_action(state):
    global POLICY
    return np.argmax(POLICY[state.a][state.b][state.c][state.d][state.e]) - 1  # {-1, 0, 1}

def next_state(e):
    '''
    本体不能超出两边界，否则保持原位置
    '''
    if e < 0:
        return 0
    if e > 4:
        return 4
    return e


def choice(r, c):
    random.seed()
    rs = []
    for _ in range(c):
        t = random.choice(r)
        rs.append(t)
        r.remove(t)
    return rs


def random_start():
    '''
    随机初始化  
    '''
    random.seed()
    l = choice(list(range(15)), 4)  # 初始化四个球
    l.append(random.randint(0, 4))  # 初始化自己
    return State(l)

def random_state():
    '''
    随机一个球
    '''
    random.seed()
    return random.randint(0, 14)

def is_out(s):
    return s > 24


def step(state, action):
    '''
    在state状态下执行action
    返回新的状态，奖励，是否结束
    '''
    e = state.get(4)
    e1 = next_state(e + action)
    l = [state.get(0) + 5, state.get(1) + 5, state.get(2) + 5, state.get(3) + 5, e1]
    reward = 0
    for i in range(4):
        if is_out(l[i]):
            reward += 1
            l[i] = random_state()
    new_state = State(l)
    new_state.set_point(state.point() + reward)
    end, fail = is_end(new_state)
    if end:
        if fail:
            return new_state, -1, True
        return new_state, reward, True
    return new_state, reward, False


def is_end(state):
    '''
    是否终结
    '''
    e = state.get(4)
    for i in range(4):
        a = state.get(i)
        if a - 20 == e:
            return True, True
    return state.point() >= 100, False

class Episode:

    def __init__(self):
        self.s = []
        self.a = []
        self.r = []

    def add_state(self, s):
        self.s.append(s)

    def add_action(self, a):
        self.a.append(a)

    def add_reward(self, r):
        self.r.append(r)

    def states(self):
        return self.s

    def action(self):
       return self.a

    def reward(self):
        return self.r


def generate_episode():
    '''
    生成一个片段
    '''
    episode = Episode()
    state = random_start()
    end = is_end(state)
    while not end:
        action = get_action(state)
        episode.add_state(state)
        episode.add_action(action)
        state, reward, end = step(state, action)
        episode.add_reward(reward)
    return episode


def update_policy(state, action):
    global POLICY
    p = POLICY[state.get(0)][state.get(1)][state.get(2)][state.get(3)][state.get(4)]
    if np.argmax(p) - 1 != action:
        POLICY[state.get(0)][state.get(1)][state.get(2)][state.get(3)][state.get(4)][0] = 0.0
        POLICY[state.get(0)][state.get(1)][state.get(2)][state.get(3)][state.get(4)][1] = 0.0
        POLICY[state.get(0)][state.get(1)][state.get(2)][state.get(3)][state.get(4)][2] = 0.0
        POLICY[state.get(0)][state.get(1)][state.get(2)][state.get(3)][state.get(4)][action + 1] = 1.0


def run():
    global COUNT, Q
    count = 1000000
    for _ in range(count):
        episode = generate_episode()
        l = len(episode.states())
        g = 0.0
        states = episode.states()
        actions = episode.action()
        rewards = episode.reward()
        cache = set()
        for i in range(l):
            _a = actions[i]
            g = 0.9 * g + rewards[i]
            s = states[i]
            a, b, c, d, e = s.get(0), s.get(1), s.get(2), s.get(3), s.get(4)
            key = '%d_%d_%d_%d_%d' % (a, b, c, d, e)
            if key not in cache:
                cache.add(key)
                COUNT[a][b][c][d][e][_a] += 1
                Q[a][b][c][d][e][_a] += (rewards[i] - Q[a][b][c][d][e][_a]) / COUNT[a][b][c][d][e][_a]
                update_policy(s, _a)


def save():
    global POLICY
    with open('./policy.txt', mode='w', encoding='utf-8') as f:
        for a in range(25):
            for b in range(25):
                for c in range(25):
                    for d in range(25):
                        for e in range(5):
                            _a = np.argmax(POLICY[a][b][c][d][e]) - 1
                            f.write('%d %d %d %d %d     %d\n' % (a, b, c, d, e, _a))
                        f.flush()

if __name__ == "__main__":
    run()
    save()