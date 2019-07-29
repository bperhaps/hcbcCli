package main

import (
	"fmt"
	"hcbcCli/hb"
	"time"
)

func main() {
	api := hb.NewHb()

	ip := []string{
		"192.168.241.131",
		"192.168.241.132",
		"192.168.241.133",
		"192.168.241.134",
	}

	for _, v := range ip {
		api.NodeList = append(api.NodeList, v)
	}

	api.SetDeviceId([]byte("Admin"))
	api.SetPort("9999")
	api.PrintEnv()

	//time.Sleep(time.Second * 5)

	//	startTime := time.Now()
	//api.PutData("name", "minsung")
	//api.SendData()
	////fmt.Println("putdata ", time.Since(startTime))
	////st := time.Now()

	//		fmt.Println(i, "time:", time.Since(startTime))
	//fmt.Println("sendData Time :", time.Since(st))

	//if i%100 == 0 {

	d, _ := api.GetData("name")
	fmt.Println(d)

	//}

	//fmt.Print("name" + string(i) + " : ")
	//data, _ := api.GetData("name" + string(i))
	//fmt.Println(data)
	time.Sleep(time.Millisecond * 100)

	api.Close()
	//time.Sleep(time.Second * 3)

	//data, _ := api.GetData("name")
	//fmt.Println(data)
	//data, _ = api.GetData("myGF")
	//fmt.Println(data)
	//	time.Sleep(1000)

	//	fmt.Println("result : " , api.GetData("name"))

}
