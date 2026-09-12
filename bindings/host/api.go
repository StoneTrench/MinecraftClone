package host
import "context"
import "github.com/tetratelabs/wazero/api"
type Color struct{R uint8;
G uint8;
B uint8;
A uint8}
type BlockType struct{Name string;
Color Color;
IsSolid bool}
func LogInfo(fn func(m api.Module,msg string)())(HostBinding){
return HostBinding{"log-info",func(ctx context.Context,mod api.Module,msg_str uint32,msg_len uint32)(){
var msg string;
msg=(string)(uReadString(mod, msg_str, msg_len));
fn(mod,msg);
}}
}
func LogWarn(fn func(m api.Module,msg string)())(HostBinding){
return HostBinding{"log-warn",func(ctx context.Context,mod api.Module,msg_str uint32,msg_len uint32)(){
var msg string;
msg=(string)(uReadString(mod, msg_str, msg_len));
fn(mod,msg);
}}
}
func LogError(fn func(m api.Module,msg string)())(HostBinding){
return HostBinding{"log-error",func(ctx context.Context,mod api.Module,msg_str uint32,msg_len uint32)(){
var msg string;
msg=(string)(uReadString(mod, msg_str, msg_len));
fn(mod,msg);
}}
}
func RegisterBlock(fn func(m api.Module,block BlockType)(namespaceid string))(HostBinding){
return HostBinding{"register-block",func(ctx context.Context,mod api.Module,block_Name_f_str uint32,block_Name_f_len uint32,block_IsSolid_f_b uint32,block_Color_f_R_f_w uint32,block_Color_f_G_f_w uint32,block_Color_f_B_f_w uint32,block_Color_f_A_f_w uint32,namespaceid_str uint32,namespaceid_len uint32)(){
_malloc:=mod.ExportedFunction("malloc");
alloc:=func(size uint32)uint32{res,err:=_malloc.Call(ctx,uint64(size));
if err!=nil{panic(err)};
return uint32(res[0])};
var block struct{Name string;
Color Color;
IsSolid bool};
var block_Name_f string;
block_Name_f=(string)(uReadString(mod, block_Name_f_str, block_Name_f_len));
block.Name=block_Name_f;
var block_Color_f struct{R uint8;
G uint8;
B uint8;
A uint8};
var block_Color_f_R_f uint8;
block_Color_f_R_f=(uint8)(block_Color_f_R_f_w);
block_Color_f.R=block_Color_f_R_f;
var block_Color_f_G_f uint8;
block_Color_f_G_f=(uint8)(block_Color_f_G_f_w);
block_Color_f.G=block_Color_f_G_f;
var block_Color_f_B_f uint8;
block_Color_f_B_f=(uint8)(block_Color_f_B_f_w);
block_Color_f.B=block_Color_f_B_f;
var block_Color_f_A_f uint8;
block_Color_f_A_f=(uint8)(block_Color_f_A_f_w);
block_Color_f.A=block_Color_f_A_f;
block.Color=block_Color_f;
var block_IsSolid_f bool;
block_IsSolid_f=(bool)(block_IsSolid_f_b != 0);
block.IsSolid=block_IsSolid_f;
namespaceid:=fn(mod,block);
uWriteString(mod,alloc,namespaceid,namespaceid_str,namespaceid_len);
}}
}
func AppendLang(fn func(m api.Module,name string)())(HostBinding){
return HostBinding{"append-lang",func(ctx context.Context,mod api.Module,name_str uint32,name_len uint32)(){
var name string;
name=(string)(uReadString(mod, name_str, name_len));
fn(mod,name);
}}
}
func ConfigSetInt(fn func(m api.Module,name string,value int64)())(HostBinding){
return HostBinding{"config-set-int",func(ctx context.Context,mod api.Module,name_str uint32,name_len uint32,value_w int64)(){
var name string;
name=(string)(uReadString(mod, name_str, name_len));
var value int64;
value=(int64)(value_w);
fn(mod,name,value);
}}
}
func ConfigSetFloat(fn func(m api.Module,name string,value float64)())(HostBinding){
return HostBinding{"config-set-float",func(ctx context.Context,mod api.Module,name_str uint32,name_len uint32,value_w float64)(){
var name string;
name=(string)(uReadString(mod, name_str, name_len));
var value float64;
value=(float64)(value_w);
fn(mod,name,value);
}}
}
func ConfigSetString(fn func(m api.Module,name string,value string)())(HostBinding){
return HostBinding{"config-set-string",func(ctx context.Context,mod api.Module,name_str uint32,name_len uint32,value_str uint32,value_len uint32)(){
var name string;
name=(string)(uReadString(mod, name_str, name_len));
var value string;
value=(string)(uReadString(mod, value_str, value_len));
fn(mod,name,value);
}}
}
func ConfigGetInt(fn func(m api.Module,name string)(value int64))(HostBinding){
return HostBinding{"config-get-int",func(ctx context.Context,mod api.Module,name_str uint32,name_len uint32,value_w uint32)(){
var name string;
name=(string)(uReadString(mod, name_str, name_len));
value:=fn(mod,name);
uWriteAnyFixed[int64](mod,value,value_w);
}}
}
func ConfigGetFloat(fn func(m api.Module,name string)(value float64))(HostBinding){
return HostBinding{"config-get-float",func(ctx context.Context,mod api.Module,name_str uint32,name_len uint32,value_w uint32)(){
var name string;
name=(string)(uReadString(mod, name_str, name_len));
value:=fn(mod,name);
uWriteAnyFixed[float64](mod,value,value_w);
}}
}
func ConfigGetString(fn func(m api.Module,name string)(value string))(HostBinding){
return HostBinding{"config-get-string",func(ctx context.Context,mod api.Module,name_str uint32,name_len uint32,value_str uint32,value_len uint32)(){
_malloc:=mod.ExportedFunction("malloc");
alloc:=func(size uint32)uint32{res,err:=_malloc.Call(ctx,uint64(size));
if err!=nil{panic(err)};
return uint32(res[0])};
var name string;
name=(string)(uReadString(mod, name_str, name_len));
value:=fn(mod,name);
uWriteString(mod,alloc,value,value_str,value_len);
}}
}
