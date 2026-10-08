package main
import("os";"path/filepath";"testing")
func TestAuditSameLengthMismatch(t *testing.T){
 d:=t.TempDir(); a:=filepath.Join(d,"a");b:=filepath.Join(d,"b")
 os.WriteFile(a,[]byte("#!/bin/sh\nprintf 'EXPECTED\\n'\n"),0700);os.WriteFile(b,[]byte("#!/bin/sh\nprintf 'INCORRECT\\n'\n"),0700)
 *mainDnsxBinary=a;*devDnsxBinary=b
 if err:=runIndividualTestCase("example.test {{binary}} -silent");err!=nil{t.Fatal(err)}
 t.Log("CONFIRMED: differing one-line outputs accepted as equal")
}
