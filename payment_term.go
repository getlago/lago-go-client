package lago

type PaymentTermType string

const (
	PaymentTermTypeDueOnReceipt   PaymentTermType = "due_on_receipt"
	PaymentTermTypeNet            PaymentTermType = "net"
	PaymentTermTypeEndOfMonth     PaymentTermType = "end_of_month"
	PaymentTermTypeNetEndOfMonth  PaymentTermType = "net_end_of_month"
	PaymentTermTypeDaysEndOfMonth PaymentTermType = "days_end_of_month"
	PaymentTermTypeDayOfMonth     PaymentTermType = "day_of_month"
)

type PaymentTerm struct {
	TermType    PaymentTermType `json:"term_type"`
	Days        *int            `json:"days,omitempty"`
	DayOfMonth  *int            `json:"day_of_month,omitempty"`
	MonthOffset *int            `json:"month_offset,omitempty"`
}
