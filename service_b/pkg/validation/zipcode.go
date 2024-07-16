package validation

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
)

func ValidadeZipCode(ctx context.Context, w http.ResponseWriter, marshalled []byte) error {
	ctx, span := otel.GetTracerProvider().Tracer("service_a").Start(ctx, "service_a")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, "POST", os.Getenv("SERVICE_A_URL")+"/zipcode", bytes.NewBuffer(marshalled))
	if err != nil {
		return err
	}

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("invalid zipcode")
	}

	defer res.Body.Close()

	return nil
}
