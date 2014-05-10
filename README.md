Go-Chatwork
===========

あんまりやる気ないからふんわり実装

## Example

```go
package main

import (
	"github.com/chobie/go-chatwork"
)

func main() {
	t := &chatwork.Transport{Token: "YourToken"}
	c := chatwork.NewChatwork(t.Client())
	c.SendMessage("ROOMID", "やっほー")
}
```

## License

MIT License
