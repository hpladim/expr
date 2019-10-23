package expr

import (
	"fmt"
)

//Extensive compare of scalar values
//Will do a compare of the scalareExpr supporting the following operands: '==', '!=', '>=','>','<=','<'
func compare(env *Environment, operand string, l ScalarExpr, r ScalarExpr) (Expression, error) {

	if r.Value() == nil || r.Value() == env.Null() || l.Value() == nil || l.Value() == env.Null() {
		return env.False(), nil
	}

	switch t := l.Value().(type) {
	case bool:
		vl, lok := l.Value().(bool)
		vr, rok := r.Value().(bool)
		if lok && rok {
			switch operand {
			case "==":
				if vl == vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "!=":
				if vl != vr {
					return env.True(), nil
				}
				return env.False(), nil
			case ">=":
				if vl != vr {
					if vl {
						return env.True(), nil
					}
					return env.False(), nil
				}
				return env.True(), nil
			case ">":
				if vl != vr {
					if vl {
						return env.True(), nil
					}
					return env.False(), nil
				}
				return env.False(), nil
			case "<=":
				if vl != vr {
					if vl {
						return env.False(), nil
					}
					return env.True(), nil
				}
				return env.True(), nil
			case "<":
				if vl != vr {
					if vl {
						return env.False(), nil
					}
					return env.True(), nil
				}
				return env.False(), nil
			default:
				return env.False(), fmt.Errorf("Operand not supported: %s", operand)
			}
		}

	case int:
		vl, lok := l.Value().(int)
		vr, rok := r.Value().(int)
		if lok && rok {
			switch operand {
			case "==":
				if vl == vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "!=":
				if vl != vr {
					return env.True(), nil
				}
				return env.False(), nil
			case ">=":
				if vl >= vr {
					return env.True(), nil
				}
				return env.False(), nil
			case ">":
				if vl > vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "<=":
				if vl <= vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "<":
				if vl < vr {
					return env.True(), nil
				}
				return env.False(), nil

			default:
				return env.False(), fmt.Errorf("Operand not supported: %s", operand)
			}
		}
	case uint:
		vl, lok := l.Value().(uint)
		vr, rok := r.Value().(uint)
		if lok && rok {
			switch operand {
			case "==":
				if vl == vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "!=":
				if vl != vr {
					return env.True(), nil
				}
				return env.False(), nil
			case ">=":
				if vl >= vr {
					return env.True(), nil
				}
				return env.False(), nil
			case ">":
				if vl > vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "<=":
				if vl <= vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "<":
				if vl < vr {
					return env.True(), nil
				}
				return env.False(), nil

			default:
				return env.False(), fmt.Errorf("Operand not supported: %s", operand)
			}
		}
	case int64:
		vl, lok := l.Value().(int64)
		vr, rok := r.Value().(int64)
		if lok && rok {
			switch operand {
			case "==":
				if vl == vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "!=":
				if vl != vr {
					return env.True(), nil
				}
				return env.False(), nil
			case ">=":
				if vl >= vr {
					return env.True(), nil
				}
				return env.False(), nil
			case ">":
				if vl > vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "<=":
				if vl <= vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "<":
				if vl < vr {
					return env.True(), nil
				}
				return env.False(), nil

			default:
				return env.False(), fmt.Errorf("Operand not supported: %s", operand)
			}
		}
	case uint64:
		vl, lok := l.Value().(uint64)
		vr, rok := r.Value().(uint64)
		if lok && rok {
			switch operand {
			case "==":
				if vl == vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "!=":
				if vl != vr {
					return env.True(), nil
				}
				return env.False(), nil
			case ">=":
				if vl >= vr {
					return env.True(), nil
				}
				return env.False(), nil
			case ">":
				if vl > vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "<=":
				if vl <= vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "<":
				if vl < vr {
					return env.True(), nil
				}
				return env.False(), nil

			default:
				return env.False(), fmt.Errorf("Operand not supported: %s", operand)
			}
		}
	case float32:
		vl, lok := l.Value().(float32)
		vr, rok := r.Value().(float32)
		if lok && rok {
			switch operand {
			case "==":
				if vl == vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "!=":
				if vl != vr {
					return env.True(), nil
				}
				return env.False(), nil
			case ">=":
				if vl >= vr {
					return env.True(), nil
				}
				return env.False(), nil
			case ">":
				if vl > vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "<=":
				if vl <= vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "<":
				if vl < vr {
					return env.True(), nil
				}
				return env.False(), nil

			default:
				return env.False(), fmt.Errorf("Operand not supported: %s", operand)
			}
		}
	case float64:
		vl, lok := l.Value().(float64)
		vr, rok := r.Value().(float64)
		if lok && rok {
			switch operand {
			case "==":
				if vl == vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "!=":
				if vl != vr {
					return env.True(), nil
				}
				return env.False(), nil
			case ">=":
				if vl >= vr {
					return env.True(), nil
				}
				return env.False(), nil
			case ">":
				if vl > vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "<=":
				if vl <= vr {
					return env.True(), nil
				}
				return env.False(), nil
			case "<":
				if vl < vr {
					return env.True(), nil
				}
				return env.False(), nil

			default:
				return env.False(), fmt.Errorf("Operand not supported: %s", operand)
			}
		}

	default:
		return env.False(), fmt.Errorf("Scalar datatype not supported: %T", t)
	}
	return env.False(), nil
}
