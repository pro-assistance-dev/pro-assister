package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/pro-assistance-dev/sprob/helpers/token"
)

type (
	Claim  string
	Claims []Claim
)

const (
	ClaimUserID    Claim = "user_id"
	ClaimDomainIDS Claim = "domain_ids"
)

var ClaimsSlice = Claims{ClaimUserID, ClaimDomainIDS}

func (item Claim) Inject(r *http.Request, h *token.Token) error {
	d, err := h.ExtractTokenMetadata(r, item)
	if err != nil {
		return err
	}
	*r = *r.WithContext(context.WithValue(r.Context(), item.String(), d)) //nolint:all
	// ctx = context.WithValue(ctx, claim, d)
	return nil
}

func (items Claims) Inject(r *http.Request, h *token.Token) (err error) {
	for i := range items {
		err = items[i].Inject(r, h)
		if err != nil {
			break
		}
	}
	return err
}

func (item Claim) String() string {
	return string(item)
}

func (item Claim) Split() []string {
	return strings.Split(item.String(), ",")
}

func (item Claim) FromContext(ctx context.Context) string {
	return ctx.Value(item.String()).(string)
}

func (item Claim) FromContextSlice(ctx context.Context) []string {
	return strings.Split(item.FromContext(ctx), ",")
}
