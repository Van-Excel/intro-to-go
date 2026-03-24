package main

import (
	"fmt"
	"math/rand"
)

// figure how to create and use multiple packages under one dir
// revise slices esp make(), append(), copy(), accumulating and the logic behind the object
// revise control structures
type Person struct {
	ID    int
	items int // number of items a person is buying
}

func (p *Person) totalCostOfItems() int64 {
	total := 0
	number := p.items
	for i := 0; i < number; i++ {
		total += rand.Intn(100) + 1

	}
	return int64(total)

}
func (p *Person) processingCycle() int {

	return 10000000 * p.items
}

type Counter struct {
	ID       int
	registry int64     // sum of money collected at a counter
	people   []*Person // people who used the counter
}

func (c *Counter) process(p *Person) {
	pc := p.processingCycle()
	for i := 0; i < pc; i++ {
		//do nothing
	}
	c.people = append(c.people, p)
	c.registry += p.totalCostOfItems() // for each person lets accumulate their total purchase cost
	fmt.Println(c.people)

}

func (c *Counter) ReturnPeoplesID() []int {
	IDslice := make([]int, 0)
	for _, p := range c.people {
		IDslice = append(IDslice, p.ID)

	}
	return IDslice
}

type Queue struct {
	people []*Person // people still waiting in the queue
}

func (q *Queue) NumberOfPeople() int {
	return len(q.people)
}

func (q *Queue) Pop() *Person {
	personRemoved := &Person{}

	if q.NumberOfPeople() > 0 {

		personRemoved, q.people = q.people[0], q.people[1:]
		return personRemoved
	}
	return nil

}

func main() {}
