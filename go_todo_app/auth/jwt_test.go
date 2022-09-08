package auth

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/koralle/go-web-application-development/go_todo_app/clock"
	"github.com/koralle/go-web-application-development/go_todo_app/entity"
	"github.com/koralle/go-web-application-development/go_todo_app/store"
	"github.com/koralle/go-web-application-development/go_todo_app/testutil/fixture"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type FixedTommorowClocker struct{}

func (c FixedTommorowClocker) Now() time.Time {
	return clock.FixedClocker{}.Now().Add(24 * time.Hour)
}

func TestEmbed(t *testing.T) {
	want := []byte("-----BEGIN PUBLIC KEY-----")
	if !bytes.Contains(rawPublicKey, want) {
		t.Errorf("want %s, but got %s", want, rawPublicKey)
	}

	want = []byte("-----BEGIN PRIVATE KEY-----")
	if !bytes.Contains(rawPrivateKey, want) {
		t.Errorf("want %s, but got %s", want, rawPublicKey)
	}

}

func TestJWTer(t *testing.T) {
	ctx := context.Background()

	moq := &StoreMock{}
	wantID := entity.UserID(10)
	u := fixture.User(&entity.User{ID: wantID})

	moq.SaveFunc = func(ctx context.Context, key string, userID entity.UserID) error {
		if userID != wantID {
			t.Errorf("want %d, but got %d", wantID, userID)
		}

		return nil
	}

	sut, err := NewJWTer(moq, clock.RealClocker{})
	if err != nil {
		t.Fatal(err)
	}

	got, err := sut.GenerateToken(ctx, *u)
	if err != nil {
		t.Fatalf("not want err: %v", err)
	}
	if len(got) == 0 {
		t.Errorf("token is empty")
	}

}

func TestJWTer_GetToken(t *testing.T) {
	t.Parallel()

	c := clock.FixedClocker{}

	want, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer("github.com/koralle/go-web-application-development/go_todo_app").
		Subject("access_token").
		IssuedAt(c.Now()).
		Expiration(c.Now().Add(30*time.Minute)).
		Claim(RoleKey, "test").
		Claim(UserNameKey, "test_user").
		Build()

	if err != nil {
		t.Fatal(err)
	}

	pkey, err := jwk.ParseKey(rawPrivateKey, jwk.WithPEM(true))
	if err != nil {
		t.Fatal(err)
	}

	signed, err := jwt.Sign(want, jwt.WithKey(jwa.RS256, pkey))
	if err != nil {
		t.Fatal(err)
	}

	userID := entity.UserID(20)

	ctx := context.Background()
	moq := &StoreMock{}
	moq.LoadFunc = func(ctx context.Context, key string) (entity.UserID, error) {
		return userID, nil
	}

	sut, err := NewJWTer(moq, c)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"https://github.com/koralle",
		nil,
	)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signed))
	got, err := sut.GetToken(ctx, req)
	if err != nil {
		t.Fatalf("want no error, but got %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetToken() got = %v, want %v", got, want)
	}

}

func TestJWTer_GEtToken_NG(t *testing.T) {
	t.Parallel()

	c := clock.FixedClocker{}

	want, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer("github.com/koralle/go-web-application-development/go_todo_app").
		Subject("access_token").
		IssuedAt(c.Now()).
		Expiration(c.Now().Add(30*time.Minute)).
		Claim(RoleKey, "test").
		Claim(UserNameKey, "test_user").
		Build()

	if err != nil {
		t.Fatal(err)
	}

	pkey, err := jwk.ParseKey(rawPrivateKey, jwk.WithPEM(true))
	if err != nil {
		t.Fatal(err)
	}

	signed, err := jwt.Sign(want, jwt.WithKey(jwa.RS256, pkey))
	if err != nil {
		t.Fatal(err)
	}

	type moq struct {
		userID entity.UserID
		err    error
	}

	tasks := map[string]struct {
		c   clock.Clocker
		moq moq
	}{
		"expire": {
			c: FixedTommorowClocker{},
		},
		"notFoundInStore": {
			c: clock.FixedClocker{},
			moq: moq{
				err: store.ErrNotFound,
			},
		},
	}

	for n, tt := range tasks {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			moq := &StoreMock{}
			moq.LoadFunc = func(ctx context.Context, key string) (entity.UserID, error) {
				return tt.moq.userID, tt.moq.err
			}

			sut, err := NewJWTer(moq, tt.c)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodGet, "https://github.com/koralle", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signed))

			got, err := sut.GetToken(ctx, req)
			if err == nil {
				t.Errorf("want error, but got nil")
			}

			if got != nil {
				t.Errorf("want nil, but got %v", got)
			}
		})
	}

}
