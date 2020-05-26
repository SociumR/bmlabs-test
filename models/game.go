package models

type Game struct {
	PointsGained string `json:"points_gained" bson:"point_games"`
	WinStatus    string `json:"win_status" bson:"win_status"`
	GameType     string `json:"game_type" bson:"game_type"`
	Created      string `json:"created" bson:"created"`
}

func (g *Game) Validate() error {
	return nil
}
