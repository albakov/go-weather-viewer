package entity

type FormField struct {
	Label, Name, Value string
}

type FormError struct {
	Message string
}

type FormRegistration struct {
	Message, LoginPlaceholder, PasswordPlaceholder string
}

type FormLogin struct {
	Message, LoginPlaceholder, PasswordPlaceholder string
}
