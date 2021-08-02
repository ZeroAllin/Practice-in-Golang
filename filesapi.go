
package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "math/rand"
    "time"
    "strings"
    "github.com/streadway/amqp"
)

var (
    FileName   string 
    fullUrlFile string
    randomid string 
)



func RabbitMQ() {
    conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672")
    ch, _ := conn.Channel()

    q, _ := ch.QueueDeclare(
        "first.queue",
        true,
        false,
        false,
        false, 
        nil)

    msg := amqp.Publishing{
        Body: []byte(FileName),
    }

    ch.Publish(
        "",
        q.Name,
        false,
        false,
        msg)
}


func main() {

    fullUrlFile = "https://static3.depositphotos.com/1000575/154/i/600/depositphotos_1549339-stock-photo-lithuania-landscape-panorama.jpg"

    
    buildFileName()
    //fmt.Println(SendFileName())
    
    file := createFile()

    
    putFile(file, httpClient())
    

    RabbitMQ()
}


func random(min, max int) int {
    rand.Seed(time.Now().Unix())
    if min > max {
        return min
    } else {
        return rand.Intn(max-min) + min
    }
}

func putFile(file *os.File, client *http.Client) {
    resp, err := client.Get(fullUrlFile)

    checkError(err)

    defer resp.Body.Close()

    size, err := io.Copy(file, resp.Body)

    defer file.Close()

    checkError(err)

    fmt.Printf("Just Downloaded a file %s with size %d", FileName, size)

}

func buildFileName() {
    rand.Seed(time.Now().UnixNano())
    chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
        "abcdefghijklmnopqrstuvwxyzåäö" +
        "0123456789")
    length := 8
    var b strings.Builder
    for i := 0; i < length; i++ {
        b.WriteRune(chars[rand.Intn(len(chars))])
    }
    FileName = b.String()
}

func httpClient() *http.Client {
    client := http.Client{
        CheckRedirect: func(r *http.Request, via []*http.Request) error {
            r.URL.Opaque = r.URL.Path
            return nil
        },
    }

    return &client
}

func createFile() *os.File {
    file, err := os.Create(FileName)

    checkError(err)
    return file
}

func checkError(err error) {
    if err != nil {
        panic(err)
    }
}

