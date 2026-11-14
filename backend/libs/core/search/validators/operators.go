package validators

import "grovehub/libs/core/search/contracts"

// OperatorSet creates the allowlist used by resource-owned filter fields.
func OperatorSet(operators ...contracts.Operator) map[contracts.Operator]struct{} {
    result := make(map[contracts.Operator]struct{}, len(operators))
    for _, operator := range operators { result[operator] = struct{}{} }
    return result
}

func EqualityOperators() map[contracts.Operator]struct{} {
    return OperatorSet(contracts.OperatorEquals, contracts.OperatorNotEquals, contracts.OperatorIn, contracts.OperatorNotIn)
}

func TextOperators() map[contracts.Operator]struct{} {
    return OperatorSet(
        contracts.OperatorEquals, contracts.OperatorNotEquals, contracts.OperatorIn, contracts.OperatorNotIn,
        contracts.OperatorExists, contracts.OperatorNotExists, contracts.OperatorContains, contracts.OperatorStartsWith,
    )
}

func TimeOperators() map[contracts.Operator]struct{} {
    return OperatorSet(contracts.OperatorEquals, contracts.OperatorGreaterThan, contracts.OperatorGTE, contracts.OperatorLessThan, contracts.OperatorLTE)
}

func IsPartialOperator(operator contracts.Operator) bool {
    return operator == contracts.OperatorContains || operator == contracts.OperatorStartsWith
}

// ValidValueCount is shared by validation and query execution:
// EXISTS takes no values; IN takes at least one; other operators take one.
func ValidValueCount(operator contracts.Operator, count int) bool {
    switch operator {
    case contracts.OperatorExists, contracts.OperatorNotExists:
        return count == 0
    case contracts.OperatorIn, contracts.OperatorNotIn:
        return count > 0
    default:
        return count == 1
    }
}
