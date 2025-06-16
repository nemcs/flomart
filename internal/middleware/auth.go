package middleware

/*
Путь: internal/middleware

Файл: auth.go
Сделай в нём middleware: AuthMiddleware(secret string) func(http.Handler) http.Handler

⮕ Это middleware будет:

читать Authorization хедер

валидировать токен (используя identity.ParseToken)

ложить userID, role в context
*/
