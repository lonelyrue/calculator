package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCalc(t *testing.T) {
    tests := []struct {
        expression string
        expected   float64
        shouldFail bool
    }{
        {"3 + 4 * 2", 11, false},
        {"(1 + 2) * 3", 9, false},
        {"10 / 2 - 1", 4, false},
        {"3 + (2 * (1 + 1))", 7, false},
        {"10 / 0", 0, true}, // Division by zero
        {"invalid_expression", 0, true}, // Invalid expression
    }

    for _, test := range tests {
        result, err := Calc(test.expression)
        if test.shouldFail {
            if err == nil {
                t.Errorf("Expected an error for expression: %s", test.expression)
            }
        } else {
            if err != nil {
                t.Errorf("Unexpected error for expression: %s - %v", test.expression, err)
            }
            if result != test.expected {
                t.Errorf("For expression %s, expected %f but got %f", test.expression, test.expected, result)
            }
        }
    }
}

func TestCalculateHandler(t *testing.T) {
    tests := []struct {
        body       string
        statusCode int
        result     string
        errMsg     string
    }{
        {`{"expression": "3 + 4 * 2"}`, http.StatusOK, "11", ""},
        {`{"expression": "(1 + 2) * 3"}`, http.StatusOK, "9", ""},
        {`{"expression": "10 / 0"}`, http.StatusUnprocessableEntity, "", "Expression is not valid"},
        {`{"expression": "invalid_expression"}`, http.StatusUnprocessableEntity, "", "Expression is not valid"},
        {`{}`, http.StatusInternalServerError, "", ""},
    }

    for _, test := range tests {
        req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer([]byte(test.body)))
        if err != nil {
            t.Fatal(err)
        }
        w := httptest.NewRecorder()
        calculateHandler(w, req)

        res := w.Result()
        if res.StatusCode != test.statusCode {
            t.Errorf("Expected status code %d but got %d", test.statusCode, res.StatusCode)
        }

        var response Response
        json.NewDecoder(res.Body).Decode(&response)

        if test.errMsg != "" {
            if response.Error != test.errMsg {
                t.Errorf("Expected error message '%s' but got '%s'", test.errMsg, response.Error)
            }
        } else {
            if response.Result != test.result {
                t.Errorf("Expected result '%s' but got '%s'", test.result, response.Result)
            }
        }
    }
}
