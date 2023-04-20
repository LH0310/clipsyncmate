To build an application in Go that sends a message when the clipboard changes and its content is text, you can follow these steps:

1. Import the necessary packages:

```go
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os/exec"
	"time"

	"github.com/atotto/clipboard"
)
```

Here, we've imported the following packages:
- `bytes`: for buffer handling.
- `fmt`: for console output.
- `image`: for image handling.
- `image/png`: for encoding images to PNG format.
- `os/exec`: for executing shell commands.
- `time`: for setting time intervals.
- `github.com/atotto/clipboard`: for clipboard access.

2. Write the main function, which reads the clipboard content in a loop, checks if it's text, and sends a message if it's different from the previous clipboard content:

```go
func main() {
	var prevContents string

	for {
		// Read clipboard contents
		contents, err := clipboard.ReadAll()
		if err != nil {
			panic(err)
		}

		// Check if contents are different from previous contents and are text
		if contents != prevContents && clipboard.IsText(contents) {
			message := fmt.Sprintf("Clipboard contents changed!\n%s", contents)
			sendMessage(message)
		}

		// Remember current clipboard contents
		prevContents = contents

		// Wait for 1 second before checking again
		time.Sleep(1 * time.Second)
	}
}
```

Here, we've defined `prevContents` to keep track of the previous clipboard contents. In each iteration, we read the current clipboard contents using `clipboard.ReadAll()`, check if they're different from the previous contents and are text using `clipboard.IsText()`, and if so, send a message using a function called `sendMessage()`. We then store the current clipboard contents as `prevContents` and wait for 1 second before checking again.

3. Write the `sendMessage()` function, which encodes an image of the clipboard contents and sends it as a message:

```go
func sendMessage(message string) {
	// Encode an image of the clipboard contents as a PNG
	img, err := clipboard.ReadImage()
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		panic(err)
	}

	// Convert the image to base64 and include it in the message
	encImage := buf.Bytes()

	cmd := exec.Command("echo", "-ne", string(encImage), "|", "base64")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	encStr := string(output)

	message = fmt.Sprintf("%s\n\nClipboard screenshot:\n%s", message, encStr)

	fmt.Println(message)
	// TODO: Send message using your preferred messaging API or service
}
```

Here, we've used `clipboard.ReadImage()` to read an image of the clipboard contents, and then used the `image/png` package to encode it to PNG format. We've then converted the encoded image to base64 using a shell command, and concatenated it with the message. You'll need to replace the last line with code to actually send the message using your preferred messaging API or service.

And that's it! This app will continuously monitor the clipboard contents and send a message whenever it changes and its contents are text.