/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import "Data_Bank/fabric-deploy-tools/fabricNetwork/cmd"







func main() {
	cmd.Execute()
	// var client *simplessh.Client
    // var err error

    // if client, err = simplessh.ConnectWithPassword("10.1.0.5", "zhouj", "Zjy498072205#"); err != nil {
    //     panic(err)
    // }
	// defer client.Close()
	// out,err:=client.Exec("/home/zhouj/fabric-deploy-tools/shell/createChannel.sh")
	// if err!=nil{
	// 	fmt.Println(string(out))
	// 	fmt.Printf("%v\n",err)
	// }else{
	// 	fmt.Println(string(out))
	// }

	// client,_:=utils.Dial("zhouj","Zjy498072205#","10.1.0.5:22")
	// stdout,err:=utils.RunCommand(client, "/home/zhouj/fabric-deploy-tools/shell/createChannel.sh",false)
	// if err!=nil{
	// 	fmt.Printf("Failed!!!!!\n%s\n%v",stdout.String(),err)
	// 	os.Exit(-1)
	// }
	// fmt.Println("Success!!!!!", stdout.String())
}
