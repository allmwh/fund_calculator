package main

import (
	"fmt"
	"sort"
	"sync"
	"strconv"
	"strings"
	"github.com/ernestosuarez/itertools"
	"github.com/gin-gonic/gin"
)

type Info struct {
	Total   int
	Element []int
}

var wg sync.WaitGroup

func main() {

	//gin
	r := gin.Default()
	r.LoadHTMLGlob("./front_end/*.html")

	//index
	r.GET("/", func(c *gin.Context) {
		c.HTML(200,"index.html",nil)
	})

	//calculate and post
	r.POST("/submit", func(c *gin.Context) {
		totalCost, _ := strconv.Atoi(c.PostForm("total_cost"))
		useMoneyRaw := c.PostForm("use_money")
		useMoneyStr := strings.Split(useMoneyRaw, " ")

		use_money := make([]int, len(useMoneyStr))
		for i, s := range useMoneyStr {
			use_money[i], _ = strconv.Atoi(s)
		}

		info := bestPair(totalCost, use_money)

		c.HTML(200, "submit.html", gin.H{
			"Total":      info.Total,
			"Element":    info.Element,
			"Total_cost": totalCost,
			"Remain":     totalCost - info.Total,
		})
	})

	// //run gin server
	r.Run(":8010")
	// result := Info{}
	// result = bestPair(20000, []int{12, 32, 52, 500, 700, 600, 300, 32, 52, 500, 700, 600, 300, 32, 52, 500, 700, 600, 300, 32, 52, 500, 700, 500, 600})
	// fmt.Println(result)
}

//吐回最好的pair
func bestPair(totalCost int, useMoney []int) Info {

	ch := make(chan []Info, len(useMoney))

	wg.Add(len(useMoney))
	for i := 1; i < len(useMoney)+1; i++ {
		go func(i int, useMoney []int) {
			defer wg.Done()

			possComb := []Info{}
			for v := range itertools.CombinationsInt(useMoney, i) {
				sort.Sort(sort.Reverse(sort.IntSlice(v)))
				info := Info{
					Total:   sum(v),
					Element: v,
				}
				possComb = append(possComb, info)
			}
			fmt.Printf("goroutine %v done\n",i)
			ch <- possComb
		}(i, useMoney)
	}
	wg.Wait()
	close(ch)

	infoSlice := []Info{}
	for possComb := range ch {
		infoSlice = append(infoSlice, possComb...)
	}

	sort.Slice(infoSlice, func(i, j int) bool { return infoSlice[i].Total > infoSlice[j].Total })

	for _, i := range infoSlice {
		if i.Total <= totalCost {
			return i
		}
	}
	return Info{}

}

func sum(array []int) int {
	result := 0
	for _, v := range array {
		result += v
	}
	return result
}
