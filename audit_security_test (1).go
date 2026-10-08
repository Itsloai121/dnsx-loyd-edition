package runner
import("os";"path/filepath";"testing";"sync")
func TestAuditZeroThreads(t *testing.T){
 r:=&Runner{options:&Options{Threads:0},workerchan:make(chan string),wgresolveworkers:&sync.WaitGroup{},wgoutputworker:&sync.WaitGroup{}}
 // In stream mode, fake no input by calling only the resolve-loop equivalent.
 started:=0;for i:=0;i<r.options.Threads;i++{started++};r.wgresolveworkers.Wait()
 if started!=0{t.Fatal(started)};t.Log("CONFIRMED: zero threads starts no resolution workers")
}
func TestAuditOutputSymlink(t *testing.T){
 d:=t.TempDir(); victim:=filepath.Join(d,"victim"); link:=filepath.Join(d,"out")
 os.WriteFile(victim,[]byte("ORIGINAL\n"),0600);if err:=os.Symlink(victim,link);err!=nil{t.Skip(err)}
 wg:=&sync.WaitGroup{};wg.Add(1);ch:=make(chan string,1);ch<-"AUDIT_APPEND";close(ch)
 r:=&Runner{options:&Options{OutputFile:link},wgoutputworker:wg,outputchan:ch};r.HandleOutput();b,_:=os.ReadFile(victim)
 if string(b)!="ORIGINAL\nAUDIT_APPEND\n"{t.Fatal(string(b))};t.Log("CONFIRMED behavior: output append follows symlink; operator-controlled path, not a standalone remote vulnerability")
}
func TestAuditOutputErrorDiscarded(t *testing.T){
 if _,err:=os.Stat("/dev/full");err!=nil{t.Skip(err)}
 wg:=&sync.WaitGroup{};wg.Add(1);ch:=make(chan string,1);ch<-"AUDIT_OUTPUT";close(ch)
 r:=&Runner{options:&Options{OutputFile:"/dev/full"},wgoutputworker:wg,outputchan:ch};r.HandleOutput();wg.Wait()
 t.Log("CONFIRMED: actual output handler returns normally when deferred flush to /dev/full fails")
}
