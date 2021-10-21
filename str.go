package omw

func IndexAccessString(str string, index int) string {
  i := 0
  for _, c := range str {
    if i == index {
      return string(c)
    }
    i += 1
  }
  return ""
}
