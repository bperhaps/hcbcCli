package main

import (
	"hcbcCli/hb"
	"time"
)

func main() {
	api := hb.NewHb()
	api.SetDeviceId([]byte("Admin"))
	api.PrintEnv()

	for i := 0; i < 1; i++ {
		api.PutData("name"+string(i), "minsung")
		api.SetDebug(true)
		api.SendData()
		//fmt.Print("name" + string(i) + " : ")
		//data, _ := api.GetData("name" + string(i))
		//fmt.Println(data)
	}

	time.Sleep(time.Second * 3)

	//data, _ := api.GetData("name")
	//fmt.Println(data)
	//data, _ = api.GetData("myGF")
	//fmt.Println(data)
	//	time.Sleep(1000)

	//	fmt.Println("result : " , api.GetData("name"))

}
