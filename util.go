package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func remove(slice []*player, s int) []*player {
	return append(slice[:s], slice[s+1:]...)
}

func generateRandomBoard() [5][5]int {
	var res [5][5]int
	for i, row := range res {
		for j := range row {
			var rn int
			for rn = rand.Intn(25) + 1; contains(res, rn); rn = rand.Intn(25) + 1 {
			}
			res[i][j] = rn
		}

	}
	return res
}

func contains(arr [5][5]int, val int) bool {
	for _, smarr := range arr {
		for j := range smarr {
			if val == smarr[j] {
				return true
			}
		}
	}
	return false
}
