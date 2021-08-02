package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"image/jpeg"
	"os"
	"github.com/nfnt/resize"
	
)

var (
	veq string
)


func SendMessageRabbit() {
    conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672")
    ch, _ := conn.Channel()

    q, _ := ch.QueueDeclare(
        "first.queue",
        true,
        false,
        false,
        false, 
        nil)
    msgs, _ := ch.Consume(
        q.Name,
        "",
        true,
        false,
        false,
        false,
        nil)   
    for m := range msgs {
       veq = string(m.Body)
       break
    }
    
}



func resized() {
    imgIn, _ := os.Open(veq)
    imgJpg, _ := jpeg.Decode(imgIn)
    imgIn.Close()

    imgJpg = resize.Resize(800, 0, imgJpg, resize.Bicubic) // <-- Собственно изменение размера картинки

    imgOut, _ := os.Create(veq + "" + "NewImage.jpg")
    jpeg.Encode(imgOut, imgJpg,nil)
    imgOut.Close()
}

func main() {
	SendMessageRabbit()
	resized()
	fmt.Println("Message Received")
	}