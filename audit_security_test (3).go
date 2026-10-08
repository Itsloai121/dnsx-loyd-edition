package dnsx
import("io";"net/http";"net/http/httptest";"testing";"time";"fmt";"github.com/miekg/dns";updateutils "github.com/projectdiscovery/utils/update";fileutil "github.com/projectdiscovery/utils/file";"os";"path/filepath")
func TestAuditDoHSelfSigned(t *testing.T){
 hits:=0;s:=httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){hits++;b,_:=io.ReadAll(r.Body);q:=new(dns.Msg);if err:=q.Unpack(b);err!=nil{t.Error(err);return};a:=new(dns.Msg);a.SetReply(q);rr,_:=dns.NewRR("audit.test. 60 IN A 192.0.2.123");a.Answer=[]dns.RR{rr};wire,_:=a.Pack();w.Header().Set("Content-Type","application/dns-message");w.Write(wire)}));defer s.Close()
 _,err:=http.Get(s.URL);if err==nil{t.Fatal("self-signed baseline unexpectedly trusted")}
 o:=DefaultOptions;o.BaseResolvers=[]string{"doh:"+s.URL};o.QuestionTypes=[]uint16{dns.TypeA};o.MaxRetries=1;o.Timeout=time.Second
 c,err:=New(o);if err!=nil{t.Fatal(err)};res,err:=c.QueryOne("audit.test");if err!=nil{t.Fatal(err)};if hits==0||len(res.A)==0||res.A[0]!="192.0.2.123"{t.Fatal(res)}
 t.Log("CONFIRMED: actual DNSX accepted DNS answer over untrusted self-signed HTTPS")
}
func TestAuditVersionTLS(t *testing.T){
 s:=httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){fmt.Fprint(w,"AUDIT_OK")}));defer s.Close()
 res,err:=updateutils.DefaultHttpClient.Get(s.URL);if err!=nil{t.Fatal(err)};defer res.Body.Close();b,_:=io.ReadAll(res.Body);if string(b)!="AUDIT_OK"{t.Fatal(string(b))};t.Log("CONFIRMED: version-check HTTP client accepts untrusted self-signed HTTPS")
}
func TestAuditTempRace(t *testing.T){
 name,err:=fileutil.GetTempFileName();if err!=nil{t.Fatal(err)};defer os.Remove(name)
 if _,err:=os.Stat(name);!os.IsNotExist(err){t.Fatal("temp path should have been removed")}
 d:=t.TempDir();victim:=filepath.Join(d,"victim");os.WriteFile(victim,[]byte("ORIGINAL"),0600)
 if err:=os.Symlink(victim,name);err!=nil{t.Fatal(err)}
 f,err:=os.Create(name);if err!=nil{t.Fatal(err)};f.WriteString("AUDIT_REPLACED");f.Close();b,_:=os.ReadFile(victim)
 if string(b)!="AUDIT_REPLACED"{t.Fatal(string(b))};t.Log("CONFIRMED controlled interleaving: helper deletes temp path, replacement symlink is followed by os.Create; natural race timing not measured")
}
