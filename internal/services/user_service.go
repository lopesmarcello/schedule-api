package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lopesmarcello/schedule-api/internal/repositories/pg"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	pool    *pgxpool.Pool
	queries *pg.Queries
}

func NewUserService(pool *pgxpool.Pool) UserService {
	return UserService{
		pool:    pool,
		queries: pg.New(pool),
	}
}

func (us *UserService) CreateUser(ctx context.Context, name, email, password string) (pg.User, error) {
	existing, _ := us.queries.GetUserByEmail(ctx, email)
	if existing.ID != 0 {
		return pg.User{}, errors.New("duplicated e-mail")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return pg.User{}, err
	}

	var slug string
	modifier := 0
	isSlugNotUnique := true

	slug = strings.ToLower(strings.ReplaceAll(name, " ", "-"))

	for isSlugNotUnique {
		exists, _ := us.queries.GetUserBySlug(ctx, slug)
		if exists.ID != 0 {
			modifier++
			slug = slug + "-" + strconv.Itoa(modifier)
		}

		if exists.ID == 0 {
			isSlugNotUnique = false
		}
	}

	params := pg.CreateUserParams{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Slug:         slug,
	}

	dbUser, err := us.queries.CreateUser(ctx, params)
	if err != nil {
		return pg.User{}, err
	}

	return dbUser, nil
}

func (us *UserService) AuthenticateUser(ctx context.Context, email, password string) (pg.User, error) {
	user, err := us.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			fmt.Println("err no rows", err)
			return pg.User{}, errors.New("invalid credentials")
		}
		fmt.Println("err searching for user", err)
		return pg.User{}, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		fmt.Println("err comparing hash", err)
		return pg.User{}, errors.New("invalid credentials")
	}

	return user, nil
}
