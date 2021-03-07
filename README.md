### go-context.v2

Library for simple control execution flow
It extends standard context functionality with flow control

Example:

```go
package main

import (
	"fmt"
    "github.com/encobrain/go-context.v2"
    "time"
)

func main () {
    context.Main.Child("main", func (ctx context.Context){
        ctx.PanicHandlerSet(func (ctx context.Context, panicErr interface{}){
            fmt.Printf("Main panic catch: %s\n", panicErr)
        })

        ctx.Child("child", func (ctx context.Context){ 
            
            loop:
            for {
                select {
                case <-time.After(time.Second):
                    fmt.Printf("Child work...\n")
                case <-ctx.Done():
                    break loop
                }       
            }
            
            fmt.Printf("Child done with reason: %s\nFinishing...\n", ctx.Err())
            
            <-time.After(time.Second)

            fmt.Printf("Child finished")
        })
        
        fmt.Printf("Main long execution...\n")

        <-time.After(time.Second*5)
    
        panic("Oops. Something went wrong")
    })    
    
    <-context.Main.ChildsFinished(true)
    
    fmt.Printf("Main and all childs finished\n")
}






```