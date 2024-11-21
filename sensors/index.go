package main

import (
    "fmt"
    "math/rand"
    "time"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "encoding/json"
    "os"
)

// Valores mínimos (80% centrais):
// temperatura: id=1398381, valor=11.3
// umidade: id=2097, valor=46.0
// luminosidade: id=893, valor=0.0
// ruido: id=10573, valor=54.9
// eco2: id=920, valor=0.0
// etvoc: id=893, valor=400.0
// latitude: id=3986277, valor=-29.379009
// longitude: id=1649505, valor=-51.191895

// Valores máximos (80% centrais):
// temperatura: id=7392, valor=32.0
// umidade: id=1363587, valor=94.6
// luminosidade: id=40227, valor=806.7
// ruido: id=1125, valor=88.9
// eco2: id=134595, valor=334.0
// etvoc: id=27244, valor=1698.0
// latitude: id=893, valor=0.0
// longitude: id=9511, valor=0.0

func randWithMinMax(min, max float64) float64 {
    return min + rand.Float64() * (max - min)
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
    fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
    fmt.Printf("Connect lost: %v", err)
}

func getClient() mqtt.Client {
    var broker = "mytb"
    var port = 1883
    opts := mqtt.NewClientOptions()
    brokerToken := os.Getenv("BROKER_TOKEN")

    opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
    opts.SetUsername(brokerToken)
    // opts.SetUsername("iciyju6i11eho9bi2u78")

    opts.OnConnect = connectHandler
    opts.OnConnectionLost = connectLostHandler
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        fmt.Println(token.Error())
    }

    return client
}

// func publish(client mqtt.Client) {
//     num := 10
//     for i := 0; i < num; i++ {
//         text := fmt.Sprintf("Message %d", i)
//         token := client.Publish("v1/devices/me/telemetry", byte(1), false, text)
//         token.Wait()
//     }
// }

type Measurement struct {
    Temperature float64 `json:"temperature"`
    Humidity    float64 `json:"humidity"`
    Luminosity  float64 `json:"luminosity"`
    Noise       float64 `json:"noise"`
    Eco2        float64 `json:"eco2"`
    Etvoc       float64 `json:"etvoc"`
}

func publishMessage(client mqtt.Client, measurements Measurement) {
    json, err := json.Marshal(measurements)
    topic := "v1/devices/me/telemetry"
    if err != nil {
        fmt.Println("Error parsing data")
        return
    }

    fmt.Println("Data to send: ", string(json))
    token := client.Publish(topic, 1, false, string(json))
    token.Wait()

    if token.Error() != nil {
        fmt.Println("Error sending data: ", token.Error())
        return
    }

    fmt.Println("Data sent!")
    return
}

func main() {
    fmt.Println("Starting...")

    referenceValues := map[string]map[string]float64{
        "temperatura": {
            "min": 11.3,
            "max": 32.0,
        },
        "umidade": {
            "min": 46.0,
            "max": 94.6,
        },
        "luminosidade": {
            "min": 0.0,
            "max": 806.7,
        },
        "ruido": {
            "min": 54.9,
            "max": 88.9,
        },
        "etvoc": {
            "min": 400.0,
            "max": 1698.0,
        },
        "eco2": {
            "min": 0.0,
            "max": 334.0,
        },
    }

    client := getClient()

    for true {
        temperatura := randWithMinMax(referenceValues["temperatura"]["min"], referenceValues["temperatura"]["max"])
        umidade := randWithMinMax(referenceValues["umidade"]["min"], referenceValues["umidade"]["max"])
        luminosidade := randWithMinMax(referenceValues["luminosidade"]["min"], referenceValues["luminosidade"]["max"])
        ruido := randWithMinMax(referenceValues["ruido"]["min"], referenceValues["ruido"]["max"])
        etvoc := randWithMinMax(referenceValues["etvoc"]["min"], referenceValues["etvoc"]["max"])
        eco2 := randWithMinMax(referenceValues["eco2"]["min"], referenceValues["eco2"]["max"])

        if (!client.IsConnected()) {
            client = getClient()
        }
  
        fmt.Printf("Connected: %t\n", client.IsConnected())

        // fmt.Printf("\nIs connected: %s\n", client.IsConnected())

        fmt.Printf("Measurements\n")
        fmt.Printf("temperatura: %f\n", temperatura)
        fmt.Printf("umidade: %f\n", umidade)
        fmt.Printf("luminosidade: %f\n", luminosidade)
        fmt.Printf("ruido: %f\n", ruido)
        fmt.Printf("etvoc: %f\n", etvoc)
        fmt.Printf("eco2: %f\n", eco2)
        fmt.Printf("===================||=================\n")

        if (!client.IsConnected()) {
            fmt.Println("Broker not connected, data not sent")
        } else {
            measurement := Measurement{
                Temperature: temperatura,
                Humidity: umidade,   
                Luminosity: luminosidade,
                Noise: ruido,
                Eco2: eco2, 
                Etvoc: etvoc,
            }

            publishMessage(client, measurement)
        }

        seconds := rand.Intn(10) 
        time.Sleep(time.Duration(seconds) * time.Second)
        fmt.Printf("Sleeped for %d seconds\n", seconds)
    }

}