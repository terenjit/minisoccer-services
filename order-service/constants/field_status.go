package constants

type FieldStatusString string

const (
	AvailableFieldStatus FieldStatusString = "available"
	BookedFieldStatus    FieldStatusString = "booked"
)

func (p FieldStatusString) String() string {
	return string(p)
}
