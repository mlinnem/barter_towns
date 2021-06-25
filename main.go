package main

import (
	"player"
)

func main() {

	//----setup----

	year := 0

	
	scanner := bufio.NewScanner(os.Stdin)

	player := player.construct()
	worldState := worldState.construct()
	//----main loop----

	for year <= MAX_TIME {
		
		//User decisions
	
		scanner.Scan()
		userCommand := scanner.Text()
		fmt.Println(userCommand.Text())

		//World evolves

		//---Town manager decisions

		//---Town evolves

	}
}
