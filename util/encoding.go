package util

import (
"code.google.com/p/mahonia"
)

func GBKtoUTF8(s string)string{
	
	enc := mahonia.NewDecoder("gbk")
	return enc.ConvertString(s)

}