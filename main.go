package main

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
	"github.com/solarlune/paths"
	"golang.org/x/image/font"
	_ "golang.org/x/image/font/sfnt"
	"log"
	"os"
	"path"
	"strings"
)

//go:embed assets/*
var embeddedFiles embed.FS

// CONSTANTS
const (
	map1                   = "Start_Map.tmx"
	map2                   = "Right_Map.tmx"
	map3                   = "Up_Map2.tmx"
	player_width           = 96
	player_height          = 96
	player_fps             = 12
	player_frame_per_sheet = 6
	DOWN                   = 0
	RIGHT                  = 1
	UP                     = 2
	LEFT                   = 4
	SPACE                  = 6
	mimic_width            = 32
	mimic_height           = 32
	sword_width            = 32
	sword_height           = 32
	demon_id_wt            = 160
	demon_id_ht            = 144
	skele_walk_wt          = 48
	skele_walk_ht          = 64
	dragon_width           = 191
	dragon_height          = 161
)

const controls = ("CONTROLS\n Arrow keys to move\n" +
	" J to traverse over objects\n" +
	"SpaceBar to attack/interact\n")

const objectives = ("PLAYER OBJECTIVES:\n -Grab Sword\n -Kill Mimic Chest\n" +
	"-Talk to warlock\n")

const notes = ("NOTE:\nPlayer score and health\n " +
	"in top left corner\n" + "When player dies, score is set to\n zero and health is set back to \n100 + playerscore/2, speed also\n adjust" +
	"based on score after death")

// STUCTS
type mapGame struct {
	map1       tileMap
	map2       tileMap
	map3       tileMap
	player     player
	gmWar      npc
	mimic      npc
	swords     npc
	mystic     npc
	panicSkele npc
	enemyskele npc
	demon      npc
	dragon     npc
	font       font.Face
	font2      font.Face
	sounds     audioPlay
}

type audioPlay struct {
	audioContext     *audio.Context
	playerwalk       *audio.Player
	playerNoSwordATK *audio.Player
	playerSwordATK   *audio.Player
	NPCtalk          *audio.Player
	enemyHit         *audio.Player
	playerExit       *audio.Player
	mimicHit         *audio.Player
	mimicDead        *audio.Player
}

type tileMap struct {
	Level    *tiled.Map
	tileHash map[uint32]*ebiten.Image
	pathMap  *paths.Grid
	path     *paths.Path
}

// GAME DISPLAY FUNCTIONS
func main() {
	Map1, err := tiled.LoadFile(path.Join("assets", "Maps", map1))
	Map2, err := tiled.LoadFile(path.Join("assets", "Maps", map2))
	Map3, err := tiled.LoadFile(path.Join("assets", "Maps", map3))
	windowWidth := 960
	windowHeight := 960
	ebiten.SetWindowSize(windowWidth, windowHeight)
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	ebitenImageMap := makeEbiteImagesFromMap(*Map1)
	start := tileMap{Level: Map1, tileHash: ebitenImageMap}
	ebitenImageMap2 := makeEbiteImagesFromMap(*Map2)
	mid := tileMap{Level: Map2, tileHash: ebitenImageMap2}
	ebitenImageMap3 := makeEbiteImagesFromMap(*Map3)
	gameMap := loadMapFromEmbedded(path.Join("assets", "Maps", map3))
	pathMap := makeSearchMap(gameMap)
	searchablePathMap := paths.NewGridFromStringArrays(pathMap, gameMap.Width, gameMap.Height)
	//fmt.Println(searchablePathMap)
	searchablePathMap.SetWalkable('6', false)
	searchablePathMap.SetWalkable('4', false)
	//fmt.Println(searchablePathMap)
	up := tileMap{Level: Map3, tileHash: ebitenImageMap3, pathMap: searchablePathMap}

	playerImg := LoadEmbeddedImage("NPCs", "player.png")
	myPlayer := player{spriteSheet: playerImg,
		xLoc:        580,
		yLoc:        500,
		playerState: 0,
		SPD:         2,
		ATK:         10,
		sword:       false,
		health:      100.0}

	gmWar := initNPC("GrandmasterWarlockIdle.png", 140, 10, 0)
	mimic := initNPC("HungryMimicIdle.png", 480, 200, 0)
	swords := initNPC("Swords.png", 720, 680, 0)
	mystic := initNPC("charge.png", 140, 150, 0)
	demon := initNPC("demon-idle.png", 730, 230, 0)
	panicSkele := initNPC("skeleton_walking.png", 10, 800, 0)
	enemySkele := initNPC("skeleton.png", 355, 100, 1)
	dragon := initNPC("Enemy_dragon.png", 10, 20, 10000)

	soundContext := audio.NewContext(44100)
	sounds := audioPlay{audioContext: soundContext,
		playerwalk:       LoadWav("Player_walking.wav", soundContext),
		playerNoSwordATK: LoadWav("No_sword_ATK.wav", soundContext),
		playerSwordATK:   LoadWav("Sword_ATK.wav", soundContext),
		NPCtalk:          LoadWav("Talk_To_Warlock.wav", soundContext),
		enemyHit:         LoadWav("ATK_Hit.wav", soundContext),
		playerExit:       LoadWav("Exit_Map.wav", soundContext),
		mimicHit:         LoadWav("Mimic_Sound1.wav", soundContext),
		mimicDead:        LoadWav("Mimic_Sound2.wav", soundContext),
	}

	demo := mapGame{
		map1:       start,
		map2:       mid,
		map3:       up,
		player:     myPlayer,
		gmWar:      gmWar,
		mimic:      mimic,
		swords:     swords,
		mystic:     mystic,
		demon:      demon,
		panicSkele: panicSkele,
		enemyskele: enemySkele,
		dragon:     dragon,
		font:       LoadScoreFont(),
		font2:      LoadControlsfont(),
		sounds:     sounds,
	}
	ebiten.SetWindowTitle("Project 2")
	ebiten.RunGame(&demo)
}

func (tmap tileMap) isTileBarrier(x, y int) bool {
	tileX := (x) / tmap.Level.TileWidth
	tileY := (y) / tmap.Level.TileHeight
	if tileX >= 0 && tileX < tmap.Level.Width && tileY >= 0 && tileY < tmap.Level.Height {
		tileID := tmap.Level.Layers[0].Tiles[tileY*tmap.Level.Width+tileX].ID
		return tileID == 9 || tileID == 16 || tileID == 4
	}
	return true
}

func loadMapFromEmbedded(name string) *tiled.Map {
	embeddedMap, err := tiled.LoadFile(name, tiled.WithFileSystem(embeddedFiles))
	if err != nil {
		fmt.Println("Error loading embedded map:", err)
	}
	return embeddedMap
}

func makeSearchMap(tiledMap *tiled.Map) []string {
	mapAsStringSlice := make([]string, 0, tiledMap.Height) //each row will be its own string
	row := strings.Builder{}
	for position, tile := range tiledMap.Layers[0].Tiles {
		if position%tiledMap.Width == 0 && position > 0 { // we get the 2d array as an unrolled one-d array
			mapAsStringSlice = append(mapAsStringSlice, row.String())
			row = strings.Builder{}
		}
		row.WriteString(fmt.Sprintf("%d", tile.ID))
	}
	mapAsStringSlice = append(mapAsStringSlice, row.String())
	return mapAsStringSlice
}

func LoadEmbeddedImage(folderName string, imageName string) *ebiten.Image {
	embeddedFile, err := embeddedFiles.Open(path.Join("assets", folderName, imageName))
	if err != nil {
		log.Fatal("failed to load embedded image ", imageName, err)
	}
	ebitenImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		fmt.Println("Error loading tile image:", imageName, err)
	}
	return ebitenImage
}

func makeEbiteImagesFromMap(tiledMap tiled.Map) map[uint32]*ebiten.Image {
	idToImage := make(map[uint32]*ebiten.Image)
	for _, tile := range tiledMap.Tilesets[0].Tiles {
		embeddedFile, err := embeddedFiles.Open(path.Join("assets", "tiles", tile.Image.Source))
		if err != nil {
			log.Fatal("failed to load embedded image ", embeddedFile, err)
		}
		ebitenImageTile, _, err := ebitenutil.NewImageFromReader(embeddedFile)
		if err != nil {
			fmt.Println("Error loading tile image:", tile.Image.Source, err)
		}
		idToImage[tile.ID] = ebitenImageTile
	}
	return idToImage
}
func LoadWav(name string, context *audio.Context) *audio.Player {
	playAudio, err := embeddedFiles.Open(path.Join("assets", "audio", name))
	if err != nil {
		fmt.Println("Error loading sound:", err)
	}
	playSound, err := wav.DecodeWithoutResampling(playAudio)
	if err != nil {
		fmt.Println("Error interpreting soundfile:", err)
	}
	soundPlayer, err := context.NewPlayer(playSound)
	if err != nil {
		fmt.Println("Couldn't create sound player:", err)
	}
	return soundPlayer
}
