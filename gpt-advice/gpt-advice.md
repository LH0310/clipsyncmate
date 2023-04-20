要使用 `gin` 和 `gorilla websocket` 实现一个能够同步服务器和客户端字符串的服务，我们需要完成以下几个步骤：

1. 安装 `gin` 和 `gorilla websocket`。可以使用以下命令进行安装：

```
go get -u github.com/gin-gonic/gin
go get -u github.com/gorilla/websocket
```

2. 在 `main` 函数中初始化 `gin` 引擎和 `websocket` 路由。

```go
router := gin.Default()

// 初始化 websocket 路由
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

router.GET("/ws", func(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println(err)
        return
    }
    handleWebSocketConnection(conn)
})
```

在上述代码中，我们初始化了一个 `upgrader` 变量，这是 `gorilla websocket` 库中提供的升级器。`upgrader` 可以将普通的 `HTTP` 请求升级成 `WebSocket` 连接。然后，我们将 `/ws` 路由与 `handleWebSocketConnection` 函数绑定，在该函数中处理 `WebSocket` 连接的逻辑。

3. 实现 `handleWebSocketConnection` 函数，该函数用于处理 `WebSocket` 连接的逻辑。在该函数中，我们可以使用 `gorilla websocket` 库提供的方法来处理消息和错误。

```go
func handleWebSocketConnection(conn *websocket.Conn) {
    // 在连接建立时，向客户端发送当前字符串
    err := conn.WriteMessage(websocket.TextMessage, []byte(currentString))
    if err != nil {
        log.Println(err)
        return
    }

    // 开始进行消息和错误处理
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
        if messageType == websocket.TextMessage {
            currentString = string(message)
            broadcastWebSocketMessage(conn, currentString)
        } else {
            log.Println("unknown message type")
        }
    }
}
```

在上述代码中，我们首先向客户端发送当前字符串，然后使用 `for` 循环处理消息和错误。当接收到客户端发送的消息时，我们将更新当前字符串并使用 `broadcastWebSocketMessage` 函数广播给所有连接的客户端。

4. 实现 `broadcastWebSocketMessage` 函数，该函数用于将消息广播给所有连接的客户端。

```go
func broadcastWebSocketMessage(sender *websocket.Conn, message string) {
    for conn := range websocketConnections {
        if conn != sender {
            err := conn.WriteMessage(websocket.TextMessage, []byte(message))
            if err != nil {
                log.Println(err)
            }
        }
    }
}
```

在上述代码中，我们遍历所有与服务端连接的 `WebSocket` 连接，将消息广播给所有连接的客户端，除了发送者本身。

5. 最后，我们需要启动 `http` 服务器并监听的端口号。

```go
err := router.Run(":8080")
if err != nil {
    log.Fatal(err)
}
```

完整代码如下：

```go
package main

import (
    "log"
    "sync"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var (
    currentString        = "Hello World!"
    websocketConnections = make(map[*websocket.Conn]bool)
    mutex                sync.Mutex
)

func main() {
    router := gin.Default()

    // 初始化 websocket 路由
    var upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
    }

    router.GET("/ws", func(c *gin.Context) {
        conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil {
            log.Println(err)
            return
        }
        handleWebSocketConnection(conn)
    })

    err := router.Run(":8080")
    if err != nil {
        log.Fatal(err)
    }
}

func handleWebSocketConnection(conn *websocket.Conn) {
    // 在连接建立时，向客户端发送当前字符串
    err := conn.WriteMessage(websocket.TextMessage, []byte(currentString))
    if err != nil {
        log.Println(err)
        return
    }

    // 加锁，将连接添加到 websocketConnections 中
    mutex.Lock()
    websocketConnections[conn] = true
    mutex.Unlock()

    // 开始进行消息和错误处理
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            break
        }
        if messageType == websocket.TextMessage {
            currentString = string(message)
            broadcastWebSocketMessage(conn, currentString)
        } else {
            log.Println("unknown message type")
        }
    }

    // 连接关闭时，将连接从 websocketConnections 中删除
    mutex.Lock()
    delete(websocketConnections, conn)
    mutex.Unlock()

    conn.Close()
}

func broadcastWebSocketMessage(sender *websocket.Conn, message string) {
    // 广播给所有连接的客户端
    mutex.Lock()
    for conn := range websocketConnections {
        if conn != sender {
            err := conn.WriteMessage(websocket.TextMessage, []byte(message))
            if err != nil {
                log.Println(err)
            }
        }
    }
    mutex.Unlock()
}
```

通过上述代码可以实现一个简单的 `WebSocket` 服务，它可以接收来自客户端发送的字符串，同时将该字符串广播给所有连接的客户端。