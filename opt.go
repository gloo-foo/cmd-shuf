package command

type Count int
type InputRange string
type RandomSource string

type EchoFlag bool

const (
	Echo   EchoFlag = true
	NoEcho EchoFlag = false
)

type ZeroFlag bool

const (
	Zero   ZeroFlag = true
	NoZero ZeroFlag = false
)

type RepeatFlag bool

const (
	Repeat   RepeatFlag = true
	NoRepeat RepeatFlag = false
)

type flags struct {
	Count        Count
	InputRange   InputRange
	RandomSource RandomSource
	Echo         EchoFlag
	Zero         ZeroFlag
	Repeat       RepeatFlag
}

func (c Count) Configure(flags *flags)        { flags.Count = c }
func (i InputRange) Configure(flags *flags)   { flags.InputRange = i }
func (r RandomSource) Configure(flags *flags) { flags.RandomSource = r }
func (e EchoFlag) Configure(flags *flags)     { flags.Echo = e }
func (z ZeroFlag) Configure(flags *flags)     { flags.Zero = z }
func (r RepeatFlag) Configure(flags *flags)   { flags.Repeat = r }
