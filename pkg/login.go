package pkg

type Login interface {
	QrLogin()

	PwdLogin(name, password string)
}
