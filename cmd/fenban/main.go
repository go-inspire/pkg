package main

import "fmt"

// 把学生按照成绩从高到低排序, 然后分配从1班到18班,按照s型号分配学生
func main() {
	classes := [19][53]string{}
	for i := 1; i <= 52; i++ {
		for j := 1; j <= 18; j++ {
			if i%2 == 1 {
				classes[j][i] = fmt.Sprintf("%d:%d", i, (i-1)*18+j)
			} else {
				classes[j][i] = fmt.Sprintf("%d:%d", i, i*18-j+1)
			}

		}
	}

	// 4. 打印班级信息
	for i := 1; i < len(classes); i++ {
		fmt.Printf("第 %d班级: %v\n", i, classes[i][1:])
	}

}
