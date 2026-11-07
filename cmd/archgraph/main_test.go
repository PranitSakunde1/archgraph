package main

import (
  "os"
  "path/filepath"
  "testing"
)

func TestWorkspaceDiscoversNestedRepositories(t *testing.T) {
  root:=t.TempDir()
  for _,name:=range []string{"service-a","group/service-b"} {
    if err:=os.MkdirAll(filepath.Join(root,name,".git"),0755);err!=nil { t.Fatal(err) }
  }
  got:=discover(root)
  if len(got)!=2 { t.Fatalf("expected 2 repositories, got %d: %#v",len(got),got) }
}

func TestScanFindsCodewayIdentityAndSNSCandidate(t *testing.T) {
  root:=t.TempDir()
  if err:=os.Mkdir(filepath.Join(root,".git"),0755);err!=nil {t.Fatal(err)}
  svc:=filepath.Join(root,"services","orders")
  if err:=os.MkdirAll(svc,0755);err!=nil {t.Fatal(err)}
  if err:=os.WriteFile(filepath.Join(svc,"codeway.yaml"),[]byte("service_name: orders-api\n"),0644);err!=nil {t.Fatal(err)}
  src:="package main\nfunc publish() { client.Publish(ctx, input) }\n"
  if err:=os.WriteFile(filepath.Join(svc,"main.go"),[]byte(src),0644);err!=nil {t.Fatal(err)}
  g:=Graph{}
  if err:=scanRepo(root,&g);err!=nil {t.Fatal(err)}
  foundService,foundSNS:=false,false
  for _,n:=range g.Nodes {
    if n.Type=="Service" && n.Name=="orders-api" {foundService=true}
    if n.Type=="AWSOperation" && n.Name=="aws_sns" {foundSNS=true}
  }
  if !foundService {t.Fatal("codeway service identity was not found")}
  if !foundSNS {t.Fatal("SNS candidate was not found")}
  if len(g.Edges)==0 {t.Fatal("expected inferred source-to-service edge")}
}
