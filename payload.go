package main

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strings"
)

type Request struct {
	Method string `json:"method"`
	Id     any    `json:"id"`
	Param  any    `json:"params"`
}

type Response struct {
	Id     any `json:"id"`
	Result any `json:"result"`
}

var (
	METHOD = make(map[string]bool)
)

func init() {
	methods := []string{"echo", "evaluate"}
	for _, method := range methods {
		METHOD[method] = true
	}
}

func validateRequest(req []byte) (Request, error) {
	// transform
	var request Request
	err := json.Unmarshal(req, &request)
	if err != nil {
		return request, err
	}
	// validate
	if !METHOD[request.Method] {
		err = errors.New("invalid method")
	}

	// if method is evaluate
	if request.Method == "evaluate" {
		if err := evaluateHandler(&request.Param); err != nil {
			return request, err
		}
	}

	return request, err
}

func transformResponse(req Request) ([]byte, error) {
	// transform
	res := Response{
		Id:     req.Id,
		Result: req.Param,
	}
	// validate
	b, err := json.Marshal(&res)
	// buffer \n at the end
	b = append(b, 10)
	return b, err
}

func evaluateHandler(param any) (err error) {
	expression := struct {
		Expression string `json:"expression"`
	}{}
	b, err := json.Marshal(&param)
	// transform to expression struct
	if err := json.Unmarshal(b, &expression); err != nil {
		return err
	}

	// check whether expression is application
	length := len(expression.Expression)
	if string(expression.Expression[0]) == "(" &&
		string(expression.Expression[length-1]) == ")" {
		lhs, rhs := separateApplication(expression.Expression[1 : length-1])

		// getting rhs char
		rhsChar := getRHSChar(strings.Join(rhs, ""))

		// check pattern abstraction start with !
		if string(lhs[0][0]) == "!" {

			if lhs[0] == lhs[1] {
				// if v1==v2
				expression.Expression = strings.Join(lhs[1:], "")
			} else {
				// fresh var
				var replacement string
				if string(lhs[1][0]) == "!" {
					// generate random
					charset := "abcdefghijklmnopqrstuvwxyz"
					for {
						c := charset[rand.Intn(len(charset))]
						if !rhsChar[string(c)] {
							replacement = string(c)
							break
						}
					}
				}

				if string(rhs[0][0]) == "!" {
					// subs application
					result := strings.Join(lhs[1:], "")

					// 1. check var in LHS
					char1 := string(lhs[0][1])

					// check v1 != v2 for free variable
					var char2 string
					if string(lhs[1][0]) == "!" {
						char2 = string(lhs[1][1])
					}
					if char2 != "" {
						result = strings.ReplaceAll(result, char2, replacement)
					}
					// replace all bounded variable from char1
					expression.Expression = strings.ReplaceAll(result, char1, strings.Join(rhs, ""))
				} else {
					// subs variable
					expression.Expression = strings.Join(rhs, "")
				}
			}

			b, err := json.Marshal(&expression)
			if err != nil {
				return err
			}
			// back to param
			if err := json.Unmarshal(b, &param); err != nil {
				return err
			}
		}
	}
	return err
}

func separateApplication(exp string) (lhs, rhs []string) {
	// parse expression
	var tmp string
	var arg []string

	hasPair := false
	for _, val := range exp {
		if string(val) == "(" || string(val) == "!" {
			hasPair = true
		}

		if tmp == "" {
			tmp += string(val)
		} else if (string(tmp[0]) == "(" && string(val) == ")") ||
			(string(tmp[0]) == "!" && string(val) == ".") {
			tmp += string(val)
			arg = append(arg, tmp)
			tmp = tmp[len(tmp):]
			hasPair = false
		} else {
			if string(val) == " " && tmp != "" && !hasPair {
				arg = append(arg, tmp)
				tmp = tmp[len(tmp):]
			}
			tmp += string(val)
		}
	}
	arg = append(arg, tmp)

	// separate lhs and rhs
	lhsState := true
	for _, val := range arg {
		if string(val[0]) == " " {
			lhsState = false
			rhs = append(rhs, val[1:])
			continue
		}
		if lhsState {
			lhs = append(lhs, val)
			continue
		}
		rhs = append(rhs, val)
	}
	return lhs, rhs
}

func getRHSChar(lhs string) map[string]bool {
	result := make(map[string]bool)
	for _, val := range lhs {
		result[string(val)] = true
	}
	return result
}
