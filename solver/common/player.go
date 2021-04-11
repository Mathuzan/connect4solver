package common

func OppositePlayer(player Player) Player {
	return 1 - player
}

func EndingForPlayer(ending Player, player Player) GameEnding {
	if ending == Empty {
		return Tie
	}
	if ending == NoMove {
		return NoEnding
	}
	if ending == player {
		return Win
	} else {
		return Lose
	}
}
