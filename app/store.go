
import (
	"google.golang.org/appengine/datastore"
	"dr2w/hf/game"
)

type Store struct {

}

func (s Store) WriteGame(g *game.Game) {
	return
}

func (s Store) ReadGames() []game.Game {
	return []game.Game{}
}
