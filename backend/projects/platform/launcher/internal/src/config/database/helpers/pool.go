package helpers

import ("errors";"os";"strconv";"strings")

func PoolLimits(prefix string, defaultOpen,defaultIdle int)(int,int,error){
 read:=func(suffix string,fallback,min,max int)(int,error){
  raw:=strings.TrimSpace(os.Getenv(prefix+suffix));if raw==""{return fallback,nil}
  n,err:=strconv.Atoi(raw);if err!=nil||n<min||n>max{return 0,errors.New("invalid "+prefix+suffix)}
  return n,nil
 }
 open,err:=read("_MAX_OPEN",defaultOpen,1,100);if err!=nil{return 0,0,err}
 idle,err:=read("_MAX_IDLE",min(defaultIdle,open),0,open);if err!=nil{return 0,0,err}
 return open,idle,nil
}
