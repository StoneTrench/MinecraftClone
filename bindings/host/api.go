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
