
# recraft-client

This repository is one of the part of the recraft project, it contains the code of the client.

**Note: this project is still in development, it may not work properly**

# Example code

First, clone this repository in your GOPATH, then you can start an example client  

    package main
    
    import (
    
    "fmt"
    "github.com/recraft/recraft-client/client"
    )
    
    func  main() {
    // Create client instance
    client := client.NewClient("localhost", 25565)
    
    // Obtain server info
    info, err := client.Status()
     
    if err != nil {
    panic(err)
    }
    
    // Return server data
    data, ismap := info.Description.(map[string]interface{}) 
    
    if ismap {
    
    fmt.Println(data["text"].(string))
    
    } else {
    
    fmt.Println(info.Description.(string))
    
    }
    }

# Project status
At the moment only handshake and status states are supported.


