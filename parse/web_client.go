package parse

import (
	"encoding/json"

	"github.com/he-lium/sokoban"
)

// for generating JSON representation of server-to-client actions

// gameInit represents data to be initially sent to clients
type gameInit struct {
	NPlayers  int   `json:"num_players"` // count of all players
	Me        int   `json:"me"`          // index of current player
	GameBoard board `json:"board"`       // initial board
}

// InitBoardJSON generates JSON file for initial state of the game
func InitBoardJSON(nPlayers int, curr int, b *sokoban.Board) ([]byte, error) {
	g := gameInit{nPlayers, curr, convertFromBoard(b)}
	return json.Marshal(g)
}

type actionResult struct {
	Player int  `json:"player"`
	Valid  bool `json:"move_valid"`
}

// ActionResult generates JSON for the result of a player's action
func ActionResult(player int, valid bool) []byte {
	j, _ := json.Marshal(actionResult{player, valid})
	return j
}

type opponentAction struct {
	Player    int    `json:"player"`
	Action    string `json:"action"`
	Direction string `json:"direction"`
}

// OpponentAction generates JSON for a move an opponent player has made
func OpponentAction(player int, a sokoban.Action) []byte {
	opp := opponentAction{
		player,
		sokoban.ActionTypeToStr(a.Type),
		sokoban.DirectionToStr(a.Direction),
	}
	j, _ := json.Marshal(opp)
	return j
}

// WinResultJSON generates JSON for a player win
func WinResultJSON(player int) []byte {
	a := opponentAction{player, "win", "?"}
	j, _ := json.Marshal(a)
	return j
}
