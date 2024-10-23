package funcs

const (
	Nothing = 0
	Pause   = 1
)

// TickClock is meant to run as a goroutine, which will
// run on a loop. This loop will update the clock to the "correct"
// time. The clock will update every 100 ms, with the database
// being updated every minute.
func TickClock(c chan int) {
	state := 0
	for {
		select {
		case state = <-c:
			switch state {

			}
		}
	}
}
