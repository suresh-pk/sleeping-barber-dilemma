package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	shopStartTime   = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)  // Shop opening time
	shopEndTime     = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 24, 0, 0, 0, time.Local) // Shop closing time
	waitingRoomSize = 5                                                                                           // Maximum number of customers in the waiting room
	cuttingTime     = 10 * time.Second                                                                            // Time after which the customer leaves if not served
)

type Barber struct {
	ID        int
	Name      string
	Available bool
}

func main() {
	fmt.Println("Application starts ...")
	db, err := sql.Open("postgres", "postgres://postgres:admin@localhost/barbershop?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	barbers, err := fetchBarbers(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := range barbers {
		barbers[i].Available = true
	}

	var wg sync.WaitGroup
	waitingRoom := make(chan string, waitingRoomSize)

	for i := 0; i < waitingRoomSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case customerName := <-waitingRoom:
					for {
						barber := getAvailableBarber(barbers)
						if barber != nil {
							fmt.Printf("%s is cutting the hair of %s\n", barber.Name, customerName)
							time.Sleep(cuttingTime)
							fmt.Printf("%s has finished haircut\n", customerName)
							barber.Available = true
							break
						} else {
							time.Sleep(1 * time.Second)
						}
					}
				}
			}
		}()
	}

	for {
		var customerName string
		fmt.Print("Enter the customner name [press ctrl+c to exit]:")
		fmt.Scanln(&customerName)

		if len(customerName) > 0 {
			currentTime := time.Now()
			if currentTime.Before(shopStartTime) || currentTime.After(shopEndTime) {
				fmt.Println("The shop is closed.", customerName, "is leaving.")
				continue
			}
			if len(waitingRoom) < waitingRoomSize {
				if !isBarberAvailable(barbers) {
					fmt.Println("No barbers availble.", customerName, "will stay in the waiting room.")
					waitingRoom <- customerName
					continue
				}
				select {
				case waitingRoom <- customerName:
				default:
					fmt.Println("Waiting room is full.", customerName, "is leaving.")
				}
			} else {
				fmt.Println("Waiting room is full.", customerName, "is leaving.")
			}
		}
	}

	wg.Wait()
}

func fetchBarbers(db *sql.DB) ([]Barber, error) {
	var barbers []Barber
	rows, err := db.Query("SELECT id, name FROM barbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var barber Barber
		err := rows.Scan(&barber.ID, &barber.Name)
		if err != nil {
			return nil, err
		}
		barbers = append(barbers, barber)
	}
	return barbers, nil
}

func getAvailableBarber(barbers []Barber) *Barber {
	for i := range barbers {
		if barbers[i].Available {
			barbers[i].Available = false
			return &barbers[i]
		}
	}
	return nil
}

func isBarberAvailable(barbers []Barber) bool {
	for _, b := range barbers {
		if b.Available {
			return true
		}
	}
	return false
}
