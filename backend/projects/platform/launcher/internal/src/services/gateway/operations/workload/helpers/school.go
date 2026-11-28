package helpers

import "github.com/google/uuid"

func ValidSchoolID(value string)bool{
 id,err:=uuid.Parse(value)
 return err==nil&&id!=uuid.Nil&&id.String()==value
}
