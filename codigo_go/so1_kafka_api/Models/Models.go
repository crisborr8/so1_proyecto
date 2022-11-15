package Models


type Match struct {
	Team1 string `json:"team1" bson:"team1"`
	Team2 string `json:"team2" bson:"team2"`
	Score string `json:"score" bson:"score"`
	Phase int `json:"phase" bson:"phase"`
}

