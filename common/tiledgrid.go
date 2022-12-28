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
	resourceDirectory       = "res/"
	resourceLevelsDirectory = "res/levels/"
	groundLayer             = "ground"
	objectsLayer            = "objects"
	imageLayer              = "image"
	backgroundLayer         = "background"
)

type TiledGrid struct {
	Layers            []*Layer            `json:"layers"`
	TileSetReferences []*TileSetReference `json:"tilesets"`
	TileSet           *TileSet
	TileMap           map[int]*TileData
	GroundLayer       *Layer
	ObjectLayer       *Layer
	BackgroundImage   string
}

type Layer struct {
	Data    []int         `json:"Data"`
	Height  int           `json:"height"`
	Width   int           `json:"width"`
	Name    string        `json:"name"`
	Image   string        `json:"image"`
	Objects []TiledObject `json:"objects"`
}

type TiledObject struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	X          int               `json:"x"`
	Y          int               `json:"y"`
	Width      int               `json:"width"`
	Height     int               `json:"height"`
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
	println("new tiled grid ", fileName)
	var tiledGrid TiledGrid

	levelFile, err := os.Open(filepath.Join(resourceLevelsDirectory, fileName+".json"))
	if err != nil {
		log.Fatal("opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(levelFile)
	if err = jsonParser.Decode(&tiledGrid); err != nil {
		log.Fatal("parsing config file", err.Error())
	}
	for _, l := range tiledGrid.Layers {
		if l.Name == groundLayer {
			tiledGrid.GroundLayer = l
		}
		if l.Name == objectsLayer {
			tiledGrid.ObjectLayer = l
		}
		if l.Name == imageLayer {
			tiledGrid.BackgroundImage = l.Image
		}
	}

	tiledGrid.TileSet = loadTileSet(resourceLevelsDirectory, tiledGrid.TileSetReferences[0])

	tiledGrid.TileMap = map[int]*TileData{}
	for _, tile := range tiledGrid.TileSet.Tiles {

		td := &TileData{}
		for _, prop := range tile.Properties {
			if prop.Name == "block" && prop.Value != nil {
				td.Block = (prop.Value).(bool)
			}
			if prop.Name == "platform" && prop.Value != nil {
				td.Platform = (prop.Value).(bool)
			}
			if prop.Name == "ladder" && prop.Value != nil {
				td.Ladder = (prop.Value).(bool)
			}
			if prop.Name == "damage" && prop.Value != nil {
				td.Damage = (prop.Value).(bool)
			}
		}
		tileId := tile.Id
		tiledGrid.TileMap[tileId] = td
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

			ts := tg.TileSet

			op := &ebiten.DrawImageOptions{}
			px, py := float64(((i)%layer.Width)*TileSize), float64(((i)/layer.Width)*TileSize)
			op.GeoM.Translate(px, py)
			op.GeoM.Scale(Scale, Scale)

			sx := ((tileIndex - ts.FirstGid) % ts.numTilesX) * TileSize
			sy := ((tileIndex - ts.FirstGid) / ts.numTilesX) * TileSize

			camera.DrawImage(ts.image.SubImage(image.Rect(sx, sy, sx+TileSize, sy+TileSize)).(*ebiten.Image), op)
		}
	}
}

type ObjectData struct {
	Name       string
	ObjectType string
	X          int
	Y          int
	W          int
	H          int
	Properties []*ObjectProperty
}

type ObjectProperty struct {
	Name    string
	ObjType string
	Value   interface{}
}

func (tg *TiledGrid) GetObjectData() []*ObjectData {
	var ods []*ObjectData

	if tg.ObjectLayer == nil {
		return []*ObjectData{}
	}

	for _, obj := range tg.ObjectLayer.Objects {
		od := &ObjectData{
			Name:       obj.Name,
			ObjectType: obj.Type,
			X:          obj.X,
			Y:          obj.Y,
			W:          obj.Width,
			H:          obj.Height,
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

	return ods
}

type TileData struct {
	Block    bool
	Platform bool
	Ladder   bool
	Damage   bool
}

var EmptyTile = &TileData{}

func (tg *TiledGrid) GetTileData(x int, y int) *TileData {

	index := (y * tg.GroundLayer.Width) + x
	if index < 0 || index >= len(tg.GroundLayer.Data) {
		return EmptyTile
	}
	if x < 0 || y < 0 {
		return EmptyTile
	}
	tileSetIndex := tg.GroundLayer.Data[index]
	if tileSetIndex == 0 {
		return EmptyTile
	}
	result := tg.TileMap[tileSetIndex-tg.TileSet.FirstGid]
	if result == nil {
		return EmptyTile
	}
	return result
}
