package model

type Command string

const (
	// movement
	Forward   Command = "forward"
	Back      Command = "back"
	Up        Command = "up"
	Down      Command = "down"
	TurnLeft  Command = "turnLeft"
	TurnRight Command = "turnRight"

	// digging
	Dig     Command = "dig"
	DigUp   Command = "digUp"
	DigDown Command = "digDown"

	// placing
	Place     Command = "place"
	PlaceUp   Command = "placeUp"
	PlaceDown Command = "placeDown"

	// attacking
	Attack     Command = "attack"
	AttackUp   Command = "attackUp"
	AttackDown Command = "attackDown"

	// detection (block presence)
	Detect     Command = "detect"
	DetectUp   Command = "detectUp"
	DetectDown Command = "detectDown"

	// inspection (block data)
	Inspect     Command = "inspect"
	InspectUp   Command = "inspectUp"
	InspectDown Command = "inspectDown"

	// comparison (selected slot vs world block)
	Compare     Command = "compare"
	CompareUp   Command = "compareUp"
	CompareDown Command = "compareDown"

	// sucking items from inventories/ground
	Suck     Command = "suck"
	SuckUp   Command = "suckUp"
	SuckDown Command = "suckDown"

	// dropping items
	Drop     Command = "drop"
	DropUp   Command = "dropUp"
	DropDown Command = "dropDown"

	// refueling
	Refuel Command = "refuel"

	// crafting (crafty turtles only)
	Craft Command = "craft"

	// equipping tools/peripherals
	EquipLeft  Command = "equipLeft"
	EquipRight Command = "equipRight"

	// inventory slot management
	Select     Command = "select"
	TransferTo Command = "transferTo"

	// no-op
	Stop Command = "stop"
)
