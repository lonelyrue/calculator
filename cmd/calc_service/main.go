package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "unicode"
)

type Request struct {
    Expression string `json:"expression"`
}

type Response struct {
    Result string `json:"result,omitempty"`
    Error  string `json:"error,omitempty"`
}

func Calc(expression string) (float64, error) {
    tokens := tokenize(expression)
    if len(tokens) == 0 {
        return 0, fmt.Errorf("empty expression")
    }

    result, err := evaluate(tokens)
    if err != nil {
        return 0, err
    }
    return result, nil
}

func tokenize(expression string) []string {
    var tokens []string
    var current string

    for _, r := range expression {
        if unicode.IsSpace(r) {
            continue
        }
        if isOperator(r) || r == '(' || r == ')' {
            if current != "" {
                tokens = append(tokens, current)
                current = ""
            }
            tokens = append(tokens, string(r))
        } else if unicode.IsDigit(r) || r == '.' {
            current += string(r)
        } else {
            return nil // некорректный символ
        }
    }
    if current != "" {
        tokens = append(tokens, current)
    }
    return tokens
}

func isOperator(r rune) bool {
    return r == '+' || r == '-' || r == '*' || r == '/'
}

func evaluate(tokens []string) (float64, error) {
    stack := make([]float64, 0, len(tokens))
    ops := make([]string, 0, len(tokens))

    for _, token := range tokens {
        if isOperator(rune(token[0])) {
            for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(token) {
                res, err := applyOp(&stack, &ops)
                if err != nil {
                    return 0, err
                }
                stack = append(stack, res)
            }
            ops = append(ops, token)
        } else if token == "(" {
            ops = append(ops, token)
        } else if token == ")" {
            for len(ops) > 0 && ops[len(ops)-1] != "(" {
                res, err := applyOp(&stack, &ops)
                if err != nil {
                    return 0, err
                }
                stack = append(stack, res)
            }
            if len(ops) == 0 {
                return 0, fmt.Errorf("mismatched parentheses")
            }
            ops = ops[:len(ops)-1] // Удаляем '('
        } else { // предполагаем, что это число
            val, err := strconv.ParseFloat(token, 64)
            if err != nil {
                return 0, fmt.Errorf("invalid number: %s", token)
            }
            stack = append(stack, val)
        }
    }

    for len(ops) > 0 {
        res, err := applyOp(&stack, &ops)
        if err != nil {
            return 0, err
        }
        stack = append(stack, res)
    }

    if len(stack) != 1 {
        return 0, fmt.Errorf("invalid expression")
    }
    return stack[0], nil
}

func precedence(op string) int {
    switch op {
    case "+", "-":
        return 1
    case "*", "/":
        return 2
    default:
        return 0
    }
}

func applyOp(stack *[]float64, ops *[]string) (float64, error) {
    if len(*stack) < 2 || len(*ops) == 0 {
        return 0, fmt.Errorf("insufficient values in expression")
    }

    b := (*stack)[len(*stack)-1]
    *stack = (*stack)[:len(*stack)-1]
    a := (*stack)[len(*stack)-1]
    *stack = (*stack)[:len(*stack)-1]

    op := (*ops)[len(*ops)-1]
    *ops = (*ops)[:len(*ops)-1]

    switch op {
    case "+":
        return a + b, nil
    case "-":
        return a - b, nil
    case "*":
        return a * b, nil
    case "/":
        if b == 0 {
			return 0, fmt.Errorf("division by zero")
        }
        return a / b, nil
    default:
        return 0, fmt.Errorf("invalid operator: %s", op)
    }
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
    var req Request

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
        http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
        return
    }

    result, err := Calc(req.Expression)
    if err != nil {
        http.Error(w, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
        return
    }

    response := Response{Result: strconv.FormatFloat(result, 'f', -1, 64)}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    http.HandleFunc("/api/v1/calculate", calculateHandler)
    http.ListenAndServe(":8080", nil)
}
