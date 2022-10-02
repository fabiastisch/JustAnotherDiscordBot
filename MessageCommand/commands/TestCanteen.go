package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"image"
	"image/color"
	"image/png"
	"justAnotherDiscordBot/Picture"
	"justAnotherDiscordBot/clients/canteenClient"
	"log"
	"os"
	"strconv"
	"time"
)

type TestCanteen struct {
}

func (f TestCanteen) ReactOnMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		return
	}

	if message.Content == f.Name() {
		//session.ChannelMessageSend(message.ChannelID, "bar")
		img := GetCanteenPic(time.Now())

		reader, writer, err := os.Pipe()
		defer reader.Close()

		go func() {
			// close the writer, so the reader knows there's no more data
			defer writer.Close()
			if err != nil {
				log.Panic(err)
			}
			//err = png.Encode(writer, img.GetImage())
			err = png.Encode(writer, img.GetImage())

			if err != nil {
				writer.Close()
				log.Panicln(err)
			}
		}()

		session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content:    "",
			Embeds:     nil,
			TTS:        false,
			Components: nil,
			Files: []*discordgo.File{
				{
					Name:        "welcome.png",
					ContentType: "image/png",
					Reader:      reader,
				},
			},
			AllowedMentions: nil,
			Reference:       nil,
			File:            nil,
			Embed:           nil,
		})

		return
	}
}

func (f TestCanteen) Name() string {
	return "canteen"
}

func (e TestCanteen) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "canteen",
		Description: "Erhalte Informationen zu den kommenden Mahlzeiten",
	}
}

func GetCanteenPic(date time.Time) *Picture.Picture {
	menu, err := canteenClient.GetCanteenMenu(date)
	if err != nil {
		return nil
	}

	return buildImage(menu, date)
}

func GetCanteenImageReader(date time.Time) *os.File {
	img := GetCanteenPic(date)
	if img == nil {
		return nil
	}
	reader, writer, err := os.Pipe()

	go func() {
		// close the writer, so the reader knows there's no more data
		defer writer.Close()
		if err != nil {
			log.Panic(err)
		}
		//err = png.Encode(writer, img.GetImage())
		err = png.Encode(writer, img.GetImage())

		if err != nil {
			writer.Close()
			log.Panicln(err)
		}
	}()
	return reader
}

func (e TestCanteen) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	img := GetCanteenPic(time.Now())
	reader, writer, err := os.Pipe()
	defer reader.Close()

	go func() {
		// close the writer, so the reader knows there's no more data
		defer writer.Close()
		if err != nil {
			log.Panic(err)
		}
		//err = png.Encode(writer, img.GetImage())
		err = png.Encode(writer, img.GetImage())

		if err != nil {
			writer.Close()
			log.Panicln(err)
		}
	}()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			TTS:             false,
			Content:         "",
			Components:      nil,
			Embeds:          nil,
			AllowedMentions: nil,
			Files: []*discordgo.File{
				{
					Name:        "welcome.png",
					ContentType: "image/png",
					Reader:      reader,
				},
			},
			Flags:    0,
			Choices:  nil,
			CustomID: "",
			Title:    "",
		},
	})
}

func buildImage(menu *canteenClient.Menu, date time.Time) *Picture.Picture {
	countMeals := len(menu.Meal)
	//log.Printf("Menu: %d\n", len(menu.Meal))

	if countMeals == 0 {
		return nil
	}
	if countMeals == 1 && menu.Meal[0].Category == "Feiertag" { //API returns a Holiday message in a meal
		return nil
	}
	mealYSize := 190
	//mealXImgSize := 190
	sizeX := 1000
	sizeY := countMeals * (mealYSize + 28)
	img := Picture.New(sizeX, sizeY+100)
	//img.Background(color.RGBA{B: 100, A: 255})

	topOffset := 0
	// y := bottomline
	//img.AddLabelCenterHorizontal("Willkommen "+fmt.Sprint(menu.Date), 0, color.White)
	topMargin := 0
	topOffset += topMargin
	height := img.AddLabelCenterHorizontal(
		fmt.Sprintf("Mensaplan %s, %s.%s.%s", date.Weekday(), menu.Date.Day, menu.Date.Month, menu.Date.Year), topMargin, color.White, 40)
	topOffset += height.Ceil() + 20
	//lineHeight := 1
	//img.DrawLine(topOffset, color.Black, lineHeight)
	//topOffset += lineHeight
	skipIndexes := map[int]bool{}
	for i, meal := range menu.Meal {
		if meal.Category == "Hinweis" {
			topOffset += AddMealToImage(meal, img, topOffset+len(skipIndexes)*mealYSize) + mealYSize
			skipIndexes[i] = true
			menu.Meal = append(menu.Meal[:i], menu.Meal[i+1:]...)
		}
	}

	for i, meal := range menu.Meal {
		/*if skipIndexes[i] {
			continue
		}*/
		topOffset += AddMealToImage(meal, img, topOffset+i*mealYSize)
	}

	return &img
}

func AddMealToImage(meal canteenClient.Meal, img Picture.Picture, posY int) int {
	//log.Println(i)
	//topOffset += img.AddLabelCenterHorizontalWithOffset(meal.Category, mealXImgSize-sizeX, topOffset+i*mealYSize, color.White, 26).Ceil() + 2 // 26p height + 2
	//topOffset += img.AddLabelCenterHorizontalWithOffset(meal.Category, mealXImgSize-sizeX, topOffset+i*mealYSize, color.White, 26).Ceil() + 2 // 26p height + 2

	//topOffset += img.AddLabel(0, topOffset+i*mealYSize, 26, color.White, meal.Category).Ceil() // 26p height
	x := 0 //i * 200 % sizeX
	y := posY
	mealImgSize := image.Point{}
	y += img.AddLabel(20, y, 26, color.White, meal.Category).Ceil() + 2

	if meal.Img == "true" {
		mealImg := Picture.GetImageFromURL(meal.ImgSmall)
		mealImgSize = mealImg.Bounds().Size()
		img.DrawImage(mealImg, x, y)
	}
	//y += img.AddLabelCenterHorizontalWithOffset(meal.Category, mealImgSize.X-img.GetImage().Bounds().Size().X, y, color.White, 26).Ceil() + 2 // 26p height + 2

	//log.Println(fmt.Sprintf("Draw Img: %d | %d", x, y))

	img.AddLabelCenterHorizontalWithOffset(meal.Description.Value, mealImgSize.X*1, y, color.White, 30)
	img.AddLabelCenterHorizontalWithOffsetBottom(GetPrices(meal), mealImgSize.X*1, posY+(+1)*mealImgSize.Y-25, color.White, 26)

	// Icons
	mealXOffset := 0
	addImage := func(pic image.Image) {
		img.DrawImage(pic, mealImgSize.X+mealXOffset, y+mealImgSize.Y-pic.Bounds().Size().Y)
		mealXOffset += pic.Bounds().Size().X
	}

	if ok, _ := strconv.ParseBool(meal.Alcohol); ok {
		pic, _ := Picture.GetImageFromFilePath("Picture/food/alkohol.png")
		addImage(pic)
	}
	if ok, _ := strconv.ParseBool(meal.Beef); ok {
		pic, _ := Picture.GetImageFromFilePath("Picture/food/rind.png")
		addImage(pic)
	}
	if ok, _ := strconv.ParseBool(meal.Pig); ok {
		pic, _ := Picture.GetImageFromFilePath("Picture/food/schwein.png")
		addImage(pic)
	}
	if ok, _ := strconv.ParseBool(meal.Vegetarian); ok {
		pic, _ := Picture.GetImageFromFilePath("Picture/food/vegetarisch.png")
		addImage(pic)
	}
	return y - posY
}

func GetPrices(meal canteenClient.Meal) string {
	student := "S"
	employee := "M"
	guest := "G"

	studentPrice := ""
	guestPrice := ""
	employeePrice := ""

	for _, s := range meal.Price {
		switch s.Group {
		case student:
			studentPrice = s.GroupPrice
		case guest:
			guestPrice = s.GroupPrice
		case employee:
			employeePrice = s.GroupPrice
		default:

			break
		}
	}
	if studentPrice == "" || employeePrice == "" || guestPrice == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s €", studentPrice, employeePrice, guestPrice)
}
