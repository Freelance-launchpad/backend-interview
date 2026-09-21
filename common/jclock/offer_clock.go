package jclock

import (
	"context"
	"log/slog"
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jenv"
	"github.com/google/uuid"
)

type OfferTimeRetriever interface {
	GetOfferTime(ctx context.Context, offerID uuid.UUID) (*time.Time, error)
}

type ClockGetter func(ctx context.Context, offerID uuid.UUID) Clock

func OfferClock(client OfferTimeRetriever) ClockGetter {
	return func(ctx context.Context, offerID uuid.UUID) Clock {
		if jenv.InProd() {
			return RealClock{}
		}

		t, err := client.GetOfferTime(ctx, offerID)
		if err != nil {
			slog.ErrorContext(ctx, "error retrieving time from client", slog.Any("error", err), slog.Any("offerID", offerID.String()))
			return RealClock{}
		}

		if t == nil {
			return RealClock{}
		}

		return FakeClock{now: *t}
	}
}

func RealClockGetter(_ context.Context, _ uuid.UUID) Clock {
	return RealClock{}
}

func FakeOfferClock(t time.Time) ClockGetter {
	return func(ctx context.Context, offerID uuid.UUID) Clock {
		return FakeClock{now: t}
	}
}
