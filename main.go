package main

import (
	"hcbcCli/hb"
	"time"
	"fmt"
)

func main() {
	api := hb.NewHb()
	api.SetDeviceId([]byte("Admin"))
	api.SetPort("9999")
	api.PrintEnv()

	time.Sleep(time.Second * 5)

	for i := 0; i < 100; i++{
	//	startTime := time.Now()
//api.PutData("name", strconv.FormatUint(uint64(i), 10))
	//fmt.Println("putdata ", time.Since(startTime))
		//st := time.Now()
//api.SendData()
//		fmt.Println(i, "time:", time.Since(startTime))
		//fmt.Println("sendData Time :", time.Since(st))

		//if i%100 == 0 {
			d, _ := api.GetData("name")
			fmt.Println(d)
		//}

		//fmt.Print("name" + string(i) + " : ")
		//data, _ := api.GetData("name" + string(i))
		//fmt.Println(data)
		time.Sleep(time.Millisecond* 100)
	}

	api.Close()
	//time.Sleep(time.Second * 3)

	//data, _ := api.GetData("name")
	//fmt.Println(data)
	//data, _ = api.GetData("myGF")
	//fmt.Println(data)
	//	time.Sleep(1000)

	//	fmt.Println("result : " , api.GetData("name"))

}
