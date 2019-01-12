package function

func GetValidString(str string) string {
    var str_buf []byte
    src := []byte(str)
    for _, v := range src {
        if v != 0 {
            str_buf = append(str_buf, v)
        }
    }
    return string(str_buf)
}