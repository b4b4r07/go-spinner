package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/b4b4r07/go-spinner"
)

const mediumTotal = 100

func main() {
	progressBars, _ := spinner.New()

	progressBars.Println("Below are many progress bars.")
	progressBars.Println("It is best to use the print wrappers to keep output synced up.")
	progressBars.Println("We can switch back to normal fmt after our progress bars are done.\n")

	barProgress1 := progressBars.MakeBar()
	barProgress2 := progressBars.MakeBar()
	barProgress3 := progressBars.MakeBar()
	barProgress4 := progressBars.MakeBar()
	barProgress5 := progressBars.MakeBar()
	barProgress6 := progressBars.MakeBar()

	progressBars.Println()
	progressBars.Println("And we can have blocks of text as we wait for progress bars to complete...")

	go progressBars.Listen()

	wg := &sync.WaitGroup{}
	wg.Add(6)

	go func() {
		for i := 0; i <= mediumTotal; i++ {
			barProgress1("hoge")
			time.Sleep(time.Millisecond * 15)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i <= mediumTotal; i++ {
			barProgress2("hoaa")
			time.Sleep(time.Millisecond * 25)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i <= mediumTotal; i++ {
			barProgress3("unko")
			time.Sleep(time.Millisecond * 12)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i <= mediumTotal; i++ {
			barProgress4("mioso")
			time.Sleep(time.Millisecond * 10)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i <= mediumTotal; i++ {
			barProgress5("hogsdfe")
			time.Sleep(time.Millisecond * 20)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i <= mediumTotal; i++ {
			barProgress6("asdfasdfasdfafdf")
			time.Sleep(time.Millisecond * 10)
		}
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("All Bars Complete")
}
