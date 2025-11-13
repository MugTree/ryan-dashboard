package dashboard

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/MugTree/ryan_dashboard/shared"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func TestServer_Start(t *testing.T) {

	t.Run("can start and stop server", func(t *testing.T) {

		err := godotenv.Load("../.env")
		env := EnvVars{
			IsProd:      shared.MustEnvGetBool("IS_PRODUCTION"),
			LogLocation: shared.MustEnv("APP_LOG"),
		}

		host := shared.MustEnv("HOST_TESTING")
		if err != nil {
			panic(err)
		}

		t.Logf("running on %v\n", host)

		/**
		---------------
		Db setup
		*/

		s := NewServer(host, &env)

		go func() {
			err := s.Start()
			if err != nil {
				t.Errorf(`Server Start - Expected nil error, but got "%v" (type %T)`, err, err)
			}

		}()

		defer func() {
			err := s.Stop()
			if err != nil {
				t.Errorf(`Server Stop - Expected nil error, but got "%v" (type %T)`, err, err)
			}
		}()

		// tnos
		time.Sleep(10 * time.Millisecond)

		res, err := http.Get(fmt.Sprintf("http://%v/health", host))
		if err != nil {
			t.Errorf(`Server request - Expected nil error, but got "%v" (type %T)`, err, err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf(`Server status code - Expected 200 but got "%v" (type %T)`, err, err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf(`Body - Expected nil error, but got "%v" (type %T)`, err, err)
		}

		if !strings.HasPrefix(string(body), "<!DOCTYPE html><html>") {
			t.Errorf("\nBody Content wrong %v\n", string(body))
		}

	})

}
