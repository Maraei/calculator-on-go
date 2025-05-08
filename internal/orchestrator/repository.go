package orchestrator

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Expression{}, &Task{}, &User{})
}

func (r *Repository) GetExpressionsByUser(userID uint) ([]Expression, error) {
	var expressions []Expression
	err := r.db.Where("user_id = ?", userID).Find(&expressions).Error
	return expressions, err
}

func (r *Repository) GetExpressionByID(id string) (*Expression, error) {
	var expr Expression
	err := r.db.First(&expr, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("expression not found")
	}
	return &expr, err
}

func (r *Repository) CreateExpression(userID uint, input string) (*Expression, error) {
	expr := &Expression{
		ID:     uuid.New().String(),
		UserID: uint32(userID),
		Input:  input,
		Status: "pending",
	}
	if err := r.db.Create(expr).Error; err != nil {
		return nil, err
	}
	return expr, nil
}

func (r *Repository) UpdateExpressionResult(expressionID string, result float64) error {
	return r.db.Model(&Expression{}).Where("id = ?", expressionID).
		Updates(map[string]interface{}{"status": "completed", "result": result}).Error
}

func (r *Repository) UpdateExpressionError(expressionID string, errorMsg string) error {
	return r.db.Model(&Expression{}).Where("id = ?", expressionID).
		Updates(map[string]interface{}{"status": "error", "error": errorMsg}).Error
}

func (r *Repository) CreateTask(task *Task) error {
	return r.db.Create(task).Error
}

func (r *Repository) GetPendingTask() (*Task, error) {
	var task Task
	err := r.db.Where("status = ?", "pending").First(&task).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("нет ожидающих задач")
	}
	return &task, err
}

func (r *Repository) CompleteTask(id string, result float64) error {
	return r.db.Model(&Task{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": "completed", "result": result}).Error
}

func (r *Repository) FailTask(id string, errorMsg string) error {
	return r.db.Model(&Task{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": "error", "error": errorMsg}).Error
}

func (r *Repository) GetTasksByExpression(expressionID string) ([]Task, error) {
	var tasks []Task
	err := r.db.Where("expression_id = ?", expressionID).Find(&tasks).Error
	return tasks, err
}

func (r *Repository) CreateUser(username, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &User{Username: username, Password: string(hashedPassword)}
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindUserByUsername(username string) (*User, error) {
	var user User
	err := r.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("пользователь не найден")
	}
	return &user, err
}

func (r *Repository) CheckPassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
