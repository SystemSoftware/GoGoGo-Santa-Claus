/*
Author: Manuel Eiden
Date: 2015-02-03
*/

package main

import "fmt"
import "math/rand"
import "time"

var N_ELVES = 10
var N_REINDEER = 9
var MIN_ELVES_FOR_HELP = 3
var MIN_REINDEERS_FOR_SLEIGHT = 9

func SantaClaus(elve_problem chan bool, elve_release chan bool, reind_back_from_holiday chan bool, reind_ready_for_tour chan bool, santa_tour_start chan bool){
	fmt.Println("Santa: Hohoho here's Santa")

	// There´s no while in Go :)
	elvs := 0
	reinds := 0

	for{
		select {
		// Reindeers are back from Holiday. Santa hichts them to the sligh
		case <- reind_back_from_holiday:
			reinds++
			if reinds == MIN_REINDEERS_FOR_SLEIGHT {
				fmt.Println("Santa: Prepairing Sleight")
				for i := 0; i < reinds; i++ {
					reind_ready_for_tour <- true
					<- santa_tour_start
				}
				fmt.Println("Santa: Hohoho! Gift-Time :-)")
				reinds = 0
			}
		//Problem with elves
		case <- elve_problem:
			elvs++
			if elvs == MIN_ELVES_FOR_HELP {
				fmt.Println("Santa: Helping Elves")
				for i := 0; i < elvs; i++ {
					elve_release <- true
				}
				elvs = 0
			}
		}
	}

}

func Elve(number int, elve_problem chan bool, elve_release chan bool){

	fmt.Println("Heres Elve ", number)
	for true{
		need_help := rand.Int() % 100 < 10
		if need_help{
			fmt.Println("Elve ", number, " is waiting for Santa´s help")

			//Write the Problem in the Problem-Queue
			elve_problem <- true
			//Wait for Reaction of Santa
			<-elve_release
			fmt.Println("Elve", number, "gets Help from Santa")
		}
		fmt.Println("Elve is ", number, " is working")
		waitingtime := time.Duration(rand.Int() % 5 + 2)
		time.Sleep(waitingtime * time.Second)
	}
}

func Reindeer(number int, reind_back_from_holiday chan bool, reind_ready_for_tour chan bool, santa_tour_start chan bool){

	fmt.Println("Here´s Reindeer ", number)
	for true {
		// Set the reindeer ready to go
		reind_back_from_holiday <- true
		// Wait for for Santa to be hitched
		<-reind_ready_for_tour
		// Give Santa the Signal to be ready
		santa_tour_start <- true
		fmt.Println("Reindeer ", number, " getting hitched")
		// Goto Holiday
		time.Sleep(20 * time.Second)
	}
}

func main(){

	elve_problem := make(chan bool)
	elve_release := make(chan bool)

	reind_back_from_holiday := make(chan bool)
	reind_ready_for_tour := make(chan bool)

	santa_tour_start := make(chan bool)

	go SantaClaus(elve_problem, elve_release, reind_back_from_holiday, reind_ready_for_tour, santa_tour_start)

	for i := 1; i <= N_REINDEER; i++{
		go Reindeer(i, reind_back_from_holiday, reind_ready_for_tour, santa_tour_start)
	}
	for i := 1; i <= N_ELVES; i++{
		go Elve(i, elve_problem, elve_release)
	}

	// Don´t end Programm
	i := 0
	fmt.Scan(&i)
}
