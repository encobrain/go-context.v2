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
        
            ctx.Cancel(fmt.Errorf("Main panic: %s", panicErr))
        })

        ctx.Child("child", func (ctx context.Context){ 
            count := ctx.Value("count").(int)

            loop:
            for i:=0; i<count; i++ {
                select {
                case <-time.After(time.Second):
                    fmt.Printf("Child work... %d\n", i)
                case <-ctx.Done():
                    break loop
                }       
            }
            
            fmt.Printf("Child done with reason: %s\nFinishing...\n", ctx.Err())
            
            <-time.After(time.Second)

            fmt.Printf("Child finished")
        }).Go()
        
        fmt.Printf("Main long execution...\n")

        <-time.After(time.Second*5)
    
        panic("Oops. Something went wrong")
        
    }).ValueSet("count", 10).Go()
    
    <-context.Main.ChildsFinished(true)
    
    fmt.Printf("Main and all childs finished\n")
}






```