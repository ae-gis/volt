package response

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"gitlab.warungpintar.co/back-end/libwp/constant"
)

func TestNew(t *testing.T) {
	data := map[string]interface{}{
		"message": "transaksi telah sukses",
	}
	result := New()
	result.SetData(data)
	assert.Equal(t, result.Data, data)
}

func TestResponseErrors(t *testing.T) {
	err := ErrorContext{
		Code:    constant.StatusCode(constant.StatusInternalError),
		Message: constant.StatusText(constant.StatusInternalError),
		Data: []ErrorData{
			{
				Code:    constant.StatusCode(constant.StatusInternalError),
				Key:     "1000",
				Message: "constraint unique key duplicate",
			},
		},
	}
	result := New()
	result.SetErrors(err)
	assert.Equal(t, result.Errors, err)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", result.Errors.(ErrorContext).Code)
	assert.Equal(t, "Oops something went wrong", result.Errors.(ErrorContext).Message)
}

func TestResponseErrorsJSON(t *testing.T) {
	errCtx := ErrorContext{
		Code:    constant.StatusCode(constant.StatusInternalError),
		Message: constant.StatusText(constant.StatusInternalError),
		Data: []ErrorData{
			{
				Code:    constant.StatusCode(constant.StatusInternalError),
				Key:     "1000",
				Message: "constraint unique key duplicate",
			},
		},
	}
	result := New()
	result.SetErrors(errCtx)

	r, err := http.NewRequest(http.MethodGet, "/", nil)

	assert.NoError(t, err)
	w := httptest.NewRecorder()
	w.WriteHeader(constant.StatusInternalError) // set header code

	WriteJSON(w, r, result) // Write http Body to JSON

	if got, want := w.Code, constant.StatusInternalError; got != want {
		t.Fatalf("status code got: %d, want %d", got, want)
	}

	expected, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, result.Errors, errCtx)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", result.Errors.(ErrorContext).Code)
	assert.Equal(t, "Oops something went wrong", result.Errors.(ErrorContext).Message)
	assert.Equal(t, string(expected), strings.TrimSuffix(string(actual), "\n"))
}

func TestResponseCSV(t *testing.T) {
	rows := make([][]string, 0)
	rows = append(rows, []string{"SO Number", "Nama Warung", "Area", "Fleet Number", "Jarak Warehouse", "Urutan"})
	rows = append(rows, []string{"SO45678", "WPD00011", "Jakarta Selatan", "1", "45.00", "1"})
	rows = append(rows, []string{"SO45645", "WPD001123", "Jakarta Selatan", "1", "43.00", "2"})
	rows = append(rows, []string{"SO45645", "WPD003343", "Jakarta Selatan", "1", "43.00", "3"})

	r, err := http.NewRequest(http.MethodGet, "/csv", nil)

	assert.NoError(t, err)
	w := httptest.NewRecorder()
	w.WriteHeader(constant.StatusSuccess) // set header code

	WriteCSV(w, r, rows, "result-route-fleets") // Write http Body to JSON

	if got, want := w.Code, constant.StatusSuccess; got != want {
		t.Fatalf("status code got: %d, want %d", got, want)
	}

	actual, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, `SO Number,Nama Warung,Area,Fleet Number,Jarak Warehouse,Urutan
SO45678,WPD00011,Jakarta Selatan,1,45.00,1
SO45645,WPD001123,Jakarta Selatan,1,43.00,2
SO45645,WPD003343,Jakarta Selatan,1,43.00,3
`, string(actual))
	assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")

}

func TestResponse(t *testing.T) {
	type args struct {
		httpStatusCode int
		body           string
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: "success response",
			args: args{
				httpStatusCode: 200,
				body:           "\"data\":{\"id\":1}",
			},
			want: events.APIGatewayProxyResponse{
				StatusCode: 200,
				Body:       "\"data\":{\"id\":1}",
				Headers: map[string]string{
					"Content-Type":                "application/json",
					"Access-Control-Allow-Origin": "*",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LambdaResponse(tt.args.httpStatusCode, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Response() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpResponse(t *testing.T) {

	req, err := http.NewRequest("GET", "/check", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := CheckHandler(200, `{data":{"id":1}}"`)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{data":{"id":1}}"`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func CheckHandler(statusCode int, body string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HttpResponse(w, statusCode, body)
	})
}
