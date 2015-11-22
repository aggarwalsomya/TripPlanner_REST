package main

import (
"fmt"

)

type perms struct {
	perm []string
}

var _ = fmt.Println;
var ret_g []perms

func generatePerm(startId string, locations []string){
	printPermutations(startId, locations, 0, len(locations))
}

func printPermutations(startId string, c []string, start int, inputSize int) {
		if start == inputSize -1 {
			// copy the whole c (which contains one of the permutation in perms)
			var fp perms
			fp.perm = []string {startId}
			fp.perm = append(fp.perm, c...)
			fp.perm = append(fp.perm, startId)
			ret_g = append(ret_g, fp)
//			fmt.Println("perm is ", fp.perm)
		}else {
			for i := start; i < inputSize; i++ {
				temp := c[start]
				c[start] = c[i]
				c[i] = temp

				printPermutations(startId, c, start + 1, inputSize)
				
				temp = c[start]
				c[start] = c[i]
				c[i] = temp
			}
		}
	}
