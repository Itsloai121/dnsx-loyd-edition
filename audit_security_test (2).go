package testutils
import("os";"path/filepath";"testing")
func TestAuditShellExpansion(t *testing.T){
 dir:=t.TempDir(); marker:=filepath.Join(dir,"benign-marker"); stub:=filepath.Join(dir,"dnsx-stub")
 if err:=os.WriteFile(stub,[]byte("#!/bin/sh\ncat\n"),0700);err!=nil{t.Fatal(err)}
 target:="$(printf AUDIT_OK > "+marker+")example.test"
 _,err:=RunDnsxBinaryAndGetResults(target,stub,false,nil);if err!=nil{t.Fatal(err)}
 b,err:=os.ReadFile(marker);if err!=nil{t.Fatal("not reproduced:",err)}
 if string(b)!="AUDIT_OK"{t.Fatal("unexpected marker")};t.Log("CONFIRMED: target shell expansion created local harmless marker")
}
