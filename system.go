package clap

var (
	System = Must("default")
)

func Ok() {
	System.Ok()
}

func Parse(args ...string) error {
	System.Parse(args...)

	return System.Err()
}

func SetName(name string) {
	System.Name = name
}

func SetDescription(description string) {
	System.Description = description
}

func Add(args ...Arg) error {
	System.Add(args...)

	return System.Valid()
}

func Using(elements ...any) error {
	args, err := DeriveAll(elements...)
	if err != nil {
		return err
	}
	System.Add(args...)

	return System.Valid()
}

func Err() error {
	return System.Err()
}

func Valid() error {
	return System.Valid()
}

func RawPositions() []string {
	return System.RawPositions()
}

func RawKeyValues() map[string][]string {
	return System.RawKeyValues()
}

func RawFlags() map[string][]string {
	return System.RawFlags()
}
