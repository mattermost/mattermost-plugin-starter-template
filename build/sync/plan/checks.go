package plan

type NilCheck struct {
	Params struct {
		Echo string `json:"echo"`
	}
}

func (NilCheck) Check() error {
	println("ok")
	return nil
}
