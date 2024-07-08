package util

import "reflect"

func GenerateMultiple[T any](len int, f func() T) []T {
  result := make([]T, 0, len)
  for i := 0; i < len; i++ {
    result = append(result, f())
  }
  return result
}

// ArbitraryCheck could be used for slice that has exactly the same value but the order is different,
// mostly used for data from database using "IN" clause
func ArbitraryCheck[T any](needle, haystack []T, comparator func(*T, *T) bool) bool {
  for _, h := range haystack {
    found := false
    for _, n := range needle {
      if !comparator(&h, &n) {
        continue
      }

      found = true
      if !reflect.DeepEqual(h, n) {
        return false
      }
      break
    }

    if !found {
      return false
    }
  }

  return true
}
