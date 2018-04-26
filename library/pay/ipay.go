package pay

type IPayResp interface {
	Success() bool
	LocalOrderId() string
	TransId() string
	Price() float32
	Signature() string
	Status() string
}