package pkg

type QrLogin interface {
	// qrcode login
	QrLogin()
}
type PwdLogin interface {
	// username and password login
	PwdLogin(name, password string)
}
type Login interface {
	QrLogin
	PwdLogin
}
