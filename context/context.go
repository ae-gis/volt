package context

type Key struct {
        Name string
}

func (c *Key) String() string {
        return "context value " + c.Name
}
