package health
type Response struct {Status string `json:"status"`}
type Service interface {Health() Response;Readiness()(Response,bool)}
type service struct {ready func()bool}
var _ Service=(*service)(nil)
func New(ready func()bool) Service{return &service{ready:ready}}
func(*service) Health() Response{return Response{Status:"UP"}}
func(s *service) Readiness()(Response,bool){if s.ready==nil||!s.ready(){return Response{Status:"NOT_READY"},false};return Response{Status:"READY"},true}
