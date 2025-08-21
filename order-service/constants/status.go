package constants

type OrderStatus int
type OrderStatusString string

const (
	Pending        OrderStatus = 100
	PendingPayment OrderStatus = 200
	PaymentSuccess OrderStatus = 300
	Expired        OrderStatus = 400

	PendingString        OrderStatusString = "pending"
	PendingPaymentString OrderStatusString = "pending-payment"
	PaymentSuccessString OrderStatusString = "payment-success"
	ExpiredString        OrderStatusString = "expired"
)

var mapStatusStringToInt = map[OrderStatusString]OrderStatus{
	PendingString:        Pending,
	PendingPaymentString: PendingPayment,
	PaymentSuccessString: PaymentSuccess,
	ExpiredString:        Expired,
}

var mapStatusIntToString = map[OrderStatus]OrderStatusString{
	Pending:        PendingString,
	PendingPayment: PendingPaymentString,
	PaymentSuccess: PaymentSuccessString,
	Expired:        ExpiredString,
}

func (p OrderStatus) String() string {
	return string(p.GetStatusString())
}

func (p OrderStatusString) Int() int {
	return int(p.GetStatusInt())
}

func (p OrderStatus) GetStatusString() OrderStatusString {
	return mapStatusIntToString[p]
}

func (p OrderStatusString) GetStatusInt() OrderStatus {
	return mapStatusStringToInt[p]
}
