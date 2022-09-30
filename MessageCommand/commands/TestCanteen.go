package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
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
		img := GetCanteenPic()

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

func GetCanteenPic() *Picture.Picture {
	menu, err := canteenClient.GetCanteenMenu(time.Now())
	if err != nil {
		return nil
	}

	return buildImage(menu)
}

func GetCanteenImageReader() *os.File {
	img := GetCanteenPic()
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
	img := GetCanteenPic()
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

func buildImage(menu *canteenClient.Menu) *Picture.Picture {
	countMeals := len(menu.Meal)
	mealYSize := 190
	mealXImgSize := 190
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
		fmt.Sprintf("Mensaplan %s.%s.%s", menu.Date.Day, menu.Date.Month, menu.Date.Year), topMargin, color.White, 40)
	topOffset += height.Ceil() + 20
	//lineHeight := 1
	//img.DrawLine(topOffset, color.Black, lineHeight)
	//topOffset += lineHeight

	for i, meal := range menu.Meal {
		log.Println(i)
		if meal.Img == "true" {

			topOffset += img.AddLabelCenterHorizontalWithOffset(meal.Category, mealXImgSize-sizeX, topOffset+i*mealYSize, color.White, 26).Ceil() + 2 // 26p height + 2

			//topOffset += img.AddLabel(0, topOffset+i*mealYSize, 26, color.White, meal.Category).Ceil() // 26p height
			mealImg := Picture.GetImageFromURL(meal.ImgSmall)
			x := 0 //i * 200 % sizeX
			y := topOffset + i*mealYSize

			log.Println(fmt.Sprintf("Draw Img: %d | %d", x, y))
			img.DrawImage(mealImg, x, y)

			img.AddLabelCenterHorizontalWithOffset(meal.Description.Value, mealXImgSize*1, y, color.White, 30)
			img.AddLabelCenterHorizontalWithOffsetBottom(GetPrices(meal), mealXImgSize*1, topOffset+(i+1)*mealYSize-25, color.White, 26)

			mealXOffset := 0
			if ok, _ := strconv.ParseBool(meal.Alcohol); ok {
				pic, _ := Picture.GetImageFromFilePath("Picture/food/alkohol.png")
				img.DrawImage(pic, mealXImgSize+mealXOffset, topOffset+(i+1)*mealYSize-pic.Bounds().Size().Y)
				mealXOffset += pic.Bounds().Size().X
			}
			if ok, _ := strconv.ParseBool(meal.Beef); ok {
				pic, _ := Picture.GetImageFromFilePath("Picture/food/rind.png")
				img.DrawImage(pic, mealXImgSize+mealXOffset, topOffset+(i+1)*mealYSize-pic.Bounds().Size().Y)
				mealXOffset += pic.Bounds().Size().X
			}
			if ok, _ := strconv.ParseBool(meal.Pig); ok {
				pic, _ := Picture.GetImageFromFilePath("Picture/food/schwein.png")
				img.DrawImage(pic, mealXImgSize+mealXOffset, topOffset+(i+1)*mealYSize-pic.Bounds().Size().Y)
				mealXOffset += pic.Bounds().Size().X
			}
			if ok, _ := strconv.ParseBool(meal.Vegetarian); ok {
				pic, _ := Picture.GetImageFromFilePath("Picture/food/vegetarisch.png")
				img.DrawImage(pic, mealXImgSize+mealXOffset, topOffset+(i+1)*mealYSize-pic.Bounds().Size().Y)
				mealXOffset += pic.Bounds().Size().X
			}
		}

	}

	return &img
}

func GetPrices(meal canteenClient.Meal) string {
	student := "S"
	employee := "M"
	guest := "G"

	for _, s := range meal.Price {
		switch s.Group {
		case student:
			student = s.GroupPrice
		case guest:
			guest = s.GroupPrice
		case employee:
			employee = s.GroupPrice
		}

	}
	return fmt.Sprintf("%s/%s/%s €", student, employee, guest)
}
