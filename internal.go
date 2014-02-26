package main

import (
	. "./driver"
	. "fmt"
	. "strconv"
	"time"
)

//Waits for internal and external inputs from elevator control
func Wait_for_input(int_button, ext_button chan int, int_order, ext_order, last_order chan string) {

	_ = int_order
	_ = ext_order

	var floor int

	for {
		select {
		case floor = <-int_button:
			go Send_to_floor(floor, "int")
		case floor = <-ext_button:
			go Send_to_floor(floor, "ext")
		case temp := <-last_order:
			_ = temp
			time.Sleep(50 * time.Millisecond)
		}
	}
}

//Sends elevator to specified floor
func Send_to_floor(floor int, button string) {
	current_floor := Get_floor_sensor()
	Elev_set_door_open_lamp(0)
	Set_stop_lamp(0)

	if current_floor < floor {
		Println("Going up")
		for {
			Speed(150)
			if Get_floor_sensor() == floor {
				Println("I am now at floor: " + Itoa(Get_floor_sensor()))
				Set_stop_lamp(1)
				Elev_set_door_open_lamp(1)
				time.Sleep(25 * time.Millisecond)
				Speed(0)
				if button == "int" {
					Set_button_lamp(BUTTON_COMMAND, floor, 0)
				} else {
					Set_button_lamp(BUTTON_CALL_DOWN, floor, 0)
				}
				return
			}
			time.Sleep(25 * time.Millisecond)
		}
	} else {
		Println("Going down")
		for {
			Speed(-150)
			if Get_floor_sensor() == floor {
				Println("I am now at floor: " + Itoa(Get_floor_sensor()))
				Set_stop_lamp(1)
				Elev_set_door_open_lamp(1)
				time.Sleep(25 * time.Millisecond)
				Speed(0)
				if button == "int" {
					Set_button_lamp(BUTTON_COMMAND, floor, 0)
				} else {
					Set_button_lamp(BUTTON_CALL_UP, floor, 0)
				}
				return
			}
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func KeyboardInput(ch chan int) {
	var a int

	for {
		Scan(&a)
		ch <- a
	}
}

func Ext_order(int_button chan int) {

	i := 0
	for {

		if i < 3 {
			if Get_button_signal(BUTTON_CALL_UP, i) == 1 {
				//Println("Button nr: " + Itoa(i) + " has been pressed!")
				Set_button_lamp(BUTTON_CALL_UP, i, 1)
				int_button <- i
				time.Sleep(300 * time.Millisecond)
			}
		}
		if i > 0 {
			if Get_button_signal(BUTTON_CALL_DOWN, i) == 1 {
				//Println("Button nr: " + Itoa(i) + " has been pressed!")
				Set_button_lamp(BUTTON_CALL_DOWN, i, 1)
				int_button <- i
				time.Sleep(300 * time.Millisecond)
			}
		}
		i++
		i = i % 4
		time.Sleep(25 * time.Millisecond)
	}
}

func Int_order(int_button chan int) {

	i := 0

	for {

		if Get_button_signal(BUTTON_COMMAND, i) == 1 {
			//Println("Button nr: " + Itoa(i) + " has been pressed!")
			int_button <- i
			Set_button_lamp(BUTTON_COMMAND, i, 1)
			time.Sleep(300 * time.Millisecond)
		}

		i++
		i = i % 4
		time.Sleep(25 * time.Millisecond)

	}
}

func Floor_indicator(last_order chan string) {
	Println("executing floor indicator!")
	//_ = last_order
	var floor int
	for {
		floor = Get_floor_sensor()
		if floor != -1 {
			Set_floor_indicator(floor)
			last_order <- Itoa(floor)
			time.Sleep(50 * time.Millisecond)
		}
		time.Sleep(25 * time.Millisecond)
	}
}

func main() {

	//channels
	int_button := make(chan int)
	ext_button := make(chan int)

	int_order := make(chan string)
	ext_order := make(chan string)
	last_order := make(chan string)

	// Initialize
	Init()
	Speed(0)
	Set_stop_lamp(1)

	go Floor_indicator(last_order)
	go Int_order(int_button)
	go Ext_order(ext_button)
	go Wait_for_input(int_button, ext_button, int_order, ext_order, last_order)

	neverQuit := make(chan string)
	<-neverQuit

}
