package identity

import (
	"errors"
)

// TODO Много дублей (dev‑const и var‑errors). Предлагаю:
// dev‑строку убирать в logger: logger.Log.Error("db error", slog.Any(logger.Err, err))
// user‑сообщение + http‑код держать в AppError. |
//| 💡 | Коды ("email_already_exists") лучше хранить здесь же, а текст отдавать в handler. |

// dev
const (
	ErrInvalidJSONDev   = "json decode error"
	ErrDBConnectionDev  = "database connection error"
	ErrPasswordHashDev  = "password hashing error"
	ErrTokenGenDev      = "token generation error"
	ErrUserNotFoundDev  = "user not found"
	ErrValidationDev    = "validation failed"
	ErrLoadingEnvDev    = "error loading .env file"
	ErrRunMigrationsDev = "error run migrations file"
	ErrSqlInsertDev     = "insert query failed"
	ErrSqlSelectDev     = "select query failed"
	ErrRunServerDev     = "server run failed"
)

// user
const (
	ErrInvalidJSONMsg        = "Невалидный JSON"
	ErrUnauthorizedMsg       = "Неавторизованный доступ"
	ErrInvalidCredentialsMsg = "Неверный логин или пароль"
	ErrUserAlreadyExistsMsg  = "Пользователь уже существует"
	ErrInternalServerMsg     = "Внутренняя ошибка сервера, попробуйте позже"
	ErrNotFoundMsg           = "Ресурс не найден"
	ErrBadRequestMsg         = "Некорректный запрос"
)

var (
	ErrEmailInvalid       = errors.New("некорректный email")
	ErrInvalidJSON        = errors.New("невалидный JSON")
	ErrInvalidEmailFormat = errors.New("некорректный формат email")
	ErrInvalidRole        = errors.New("указана недопустимая роль")
	ErrEmailAlreadyExists = errors.New("пользователь с таким email уже существует")
	ErrEmailSearch        = errors.New("ошибка при поиске email")
	ErrHashingPassword    = errors.New("ошибка при хэшировании пароля")
	ErrSavingUser         = errors.New("не удалось сохранить пользователя в базу")
	ErrGeneratingJWT      = errors.New("не удалось создать JWT токен")
	ErrInvalidCredentials = errors.New("неверный email или пароль")
	ErrComparingPassword  = errors.New("не удалось проверить пароль")
	ErrUserNotFound       = errors.New("пользователь не найден")
)
