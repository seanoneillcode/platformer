package common

import (
	"bytes"
	"encoding/json"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	resourceDirectory = "res/"
)

type TiledGrid struct {
	Layers            []*Layer            `json:"layers"`
	TileSetReferences []*TileSetReference `json:"tilesets"`
	TileSet           []*TileSet
}

type Layer struct {
	Data    []int         `json:"Data"`
	Height  int           `json:"height"`
	Width   int           `json:"width"`
	Objects []TiledObject `json:"objects"`
}

type TiledObject struct {
	Name       string            `json:"Name"`
	Type       string            `json:"type"`
	X          int               `json:"x"`
	Y          int               `json:"y"`
	Properties []*TileConfigProp `json:"properties"`
}

type TileSetReference struct {
	Source   string `json:"source"`
	FirstGid int    `json:"firstgid"`
}

type TileSet struct {
	ImageFileName string `json:"image"`
	ImageWidth    int    `json:"imagewidth"`
	ImageHeight   int    `json:"imageheight"`
	numTilesX     int
	numTilesY     int
	FirstGid      int
	Tiles         []*TileConfig `json:"tiles"`
	image         *ebiten.Image
}

type TileConfig struct {
	Id         int               `json:"id"`
	Properties []*TileConfigProp `json:"properties"`
}

type TileConfigProp struct {
	Name  string      `json:"Name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func NewTileGrid(fileName string) *TiledGrid {
	var tiledGrid TiledGrid

	levelDirectory := filepath.Join(resourceDirectory, fileName+"/")

	levelFile, err := os.Open(filepath.Join(levelDirectory, "level.json"))
	if err != nil {
		log.Fatal("opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(levelFile)
	if err = jsonParser.Decode(&tiledGrid); err != nil {
		log.Fatal("parsing config file", err.Error())
	}

	tiledGrid.TileSet = []*TileSet{}
	for _, ref := range tiledGrid.TileSetReferences {
		tiledGrid.TileSet = append(tiledGrid.TileSet, loadTileSet(levelDirectory, ref))
	}

	return &tiledGrid
}

func loadTileSet(levelDirectory string, ref *TileSetReference) *TileSet {
	tileSetConfigFile, err := os.Open(filepath.Join(levelDirectory, ref.Source))
	if err != nil {
		log.Fatal("opening config file", err.Error())
	}

	var tileSet TileSet
	jsonParser := json.NewDecoder(tileSetConfigFile)
	if err = jsonParser.Decode(&tileSet); err != nil {
		log.Fatal("parsing config file", err.Error())
	}
	tileSet.numTilesX = tileSet.ImageWidth / TileSize
	tileSet.numTilesY = tileSet.ImageHeight / TileSize

	b, err := ioutil.ReadFile(filepath.Join(levelDirectory, tileSet.ImageFileName))
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	tileSet.image = ebiten.NewImageFromImage(img)
	tileSet.FirstGid = ref.FirstGid
	return &tileSet
}

func (tg *TiledGrid) Draw(camera Camera) {
	for _, layer := range tg.Layers {
		for i, tileIndex := range layer.Data {
			if tileIndex == 0 {
				continue
			}

			ts := tg.getTileSetForIndex(tileIndex)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(((i)%layer.Width)*TileSize), float64(((i)/layer.Height)*TileSize))
			op.GeoM.Scale(Scale, Scale)

			sx := ((tileIndex - ts.FirstGid) % ts.numTilesX) * TileSize
			sy := ((tileIndex - ts.FirstGid) / ts.numTilesX) * TileSize

			camera.DrawImage(ts.image.SubImage(image.Rect(sx, sy, sx+TileSize, sy+TileSize)).(*ebiten.Image), op)
		}
	}
}

func (tg *TiledGrid) getTileSetForIndex(index int) *TileSet {
	for i, tileSet := range tg.TileSet {
		if i == len(tg.TileSet)-1 || tg.TileSet[i+1].FirstGid > index {
			return tileSet
		}
	}
	// should never happen
	return nil
}

type ObjectData struct {
	Name       string
	ObjectType string
	X          int
	Y          int
	Properties []*ObjectProperty
}

type ObjectProperty struct {
	Name    string
	ObjType string
	Value   interface{}
}

func (tg *TiledGrid) GetObjectData() []*ObjectData {
	var ods []*ObjectData
	for _, l := range tg.Layers {
		for _, obj := range l.Objects {
			od := &ObjectData{
				Name:       obj.Name,
				ObjectType: obj.Type,
				X:          obj.X,
				Y:          obj.Y,
				Properties: []*ObjectProperty{},
			}
			for _, p := range obj.Properties {
				od.Properties = append(od.Properties, &ObjectProperty{
					Name:    p.Name,
					ObjType: p.Type,
					Value:   p.Value,
				})
			}
			ods = append(ods, od)
		}
	}
	return ods
}

type TileData struct {
	X     int
	Y     int
	Block bool
}

func (tg *TiledGrid) GetTileData(x int, y int) *TileData {
	td := TileData{
		X: x,
		Y: y,
	}
	index := (y * tg.Layers[0].Width) + x

	if index < 0 || index >= len(tg.Layers[0].Data) {
		// no tile here
		td.Block = true
		return &td
	}

	if x < 0 || y < 0 {
		// no tile here
		td.Block = true
		return &td
	}

	tileSetIndex := tg.Layers[0].Data[index]
	if tileSetIndex == 0 {
		td.Block = true
		return &td
	}

	ts := tg.getTileSetForIndex(tileSetIndex)

	for _, tile := range ts.Tiles {
		if tile.Id == tileSetIndex-ts.FirstGid {
			for _, prop := range tile.Properties {
				if prop.Name == "block" && prop.Value != nil {
					td.Block = (prop.Value).(bool)
				}
				break
			}
			break
		}
	}

	return &td
}
