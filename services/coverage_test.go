package services

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	bxhttp "github.com/tigusigalpa/bingx-go/v2/http"
)

func TestServicesExportedMethodsExecute(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"code":0,"data":{}}`)
	}))
	defer srv.Close()

	client := bxhttp.NewBaseHTTPClient("key", "secret", srv.URL, "", "hex")
	services := []interface{}{
		NewAccountService(client), NewContractService(client), NewCopyTradingService(client),
		NewListenKeyService(client), NewMarketService(client), NewSpotAccountService(client),
		NewSpotTradeService(client), NewSubAccountService(client), NewTradeService(client), NewWalletService(client),
	}
	for _, service := range services {
		callExportedMethods(t, service)
	}
}

func callExportedMethods(t *testing.T, service interface{}) {
	t.Helper()
	value := reflect.ValueOf(service)
	for methodIndex := 0; methodIndex < value.NumMethod(); methodIndex++ {
		method := value.Type().Method(methodIndex)
		arguments := make([]reflect.Value, method.Type.NumIn()-1)
		for i := range arguments {
			arguments[i] = reflect.Zero(method.Type.In(i + 1))
		}
		func() {
			defer func() {
				if recovered := recover(); recovered != nil {
					t.Errorf("%s panicked: %v", method.Name, recovered)
				}
			}()
			value.Method(methodIndex).Call(arguments)
		}()
		for i := range arguments {
			arguments[i] = populatedValue(method.Type.In(i + 1))
		}
		func() {
			defer func() { _ = recover() }()
			value.Method(methodIndex).Call(arguments)
		}()
	}
}

func populatedValue(valueType reflect.Type) reflect.Value {
	switch valueType.Kind() {
	case reflect.String:
		return reflect.ValueOf("test")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.New(valueType).Elem()
	case reflect.Float32, reflect.Float64:
		return reflect.New(valueType).Elem()
	case reflect.Bool:
		return reflect.ValueOf(true)
	case reflect.Ptr:
		return reflect.New(valueType.Elem())
	case reflect.Map:
		return reflect.MakeMap(valueType)
	case reflect.Slice:
		return reflect.MakeSlice(valueType, 0, 0)
	default:
		return reflect.Zero(valueType)
	}
}
