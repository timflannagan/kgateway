package pluginsdk

// registry is internal to the SDK; core will read from here.
var registry []Plugin

// Register allows external programs to provide plugins **before**
// they invoke pkg/setup.New().
func Register(p Plugin) {
	registry = append(registry, p)
}
