package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "os"
  "path/filepath"
  "regexp"
  "sort"
  "strings"
  "time"
)

type Evidence struct { Repository string; File string; Line int; Extractor string; Detail string }
type Node struct { ID string; Type string; Name string; Repository string; Evidence []Evidence }
type Edge struct { From string; To string; Type string; Status string; Evidence []Evidence }
type Graph struct { SchemaVersion string; GeneratedAt time.Time; Repositories []string; Nodes []Node; Edges []Edge; Warnings []string }

var codeExt = map[string]string{".go":"Go",".java":"Java",".py":"Python",".ts":"TypeScript",".tsx":"TypeScript",".js":"JavaScript",".jsx":"JavaScript"}
var signatures = map[string]*regexp.Regexp{
 "aws_sns":regexp.MustCompile(`(?i)(Publish\(|PublishInput|SNSClient|sns\.publish)`),
 "aws_sqs":regexp.MustCompile(`(?i)(SendMessage\(|ReceiveMessage\(|SQSClient)`),
 "aws_eventbridge":regexp.MustCompile(`(?i)(PutEvents\(|EventBridgeClient)`),
 "aws_kinesis":regexp.MustCompile(`(?i)(PutRecord\(|PutRecords\(|KinesisClient)`),
 "aws_lambda":regexp.MustCompile(`(?i)(Invoke\(|LambdaClient|lambda\.invoke)`),
 "aws_dynamodb":regexp.MustCompile(`(?i)(DynamoDbClient|DynamoDBClient|GetItem\(|PutItem\(|UpdateItem\()`),
 "database_mongo":regexp.MustCompile(`(?i)(mongo\.Connect|MongoClient|collection\.(Find|Insert|Update|Delete))`),
 "database_sql":regexp.MustCompile(`(?i)(database/sql|sql\.Open|JdbcTemplate|@Query\(|SELECT\s+.+\s+FROM)`),
 "api_route":regexp.MustCompile(`(?i)(http\.HandleFunc|@GetMapping|@PostMapping|@RequestMapping|app\.(get|post|put|delete)|router\.(get|post|put|delete))`),
}
func main() {
 if len(os.Args)<2 || os.Args[1]!="scan" { fmt.Println("archgraph scan --path <workspace> --output <directory>"); os.Exit(2) }
 fs:=flag.NewFlagSet("scan",flag.ExitOnError)
 root:=fs.String("path",".","workspace path")
 fs.StringVar(root,"workspace",".","workspace path")
 out:=fs.String("output","./.archgraph","output directory")
 _=fs.Parse(os.Args[2:])
 repos:=discover(*root)
 if len(repos)==0 { fmt.Fprintln(os.Stderr,"no Git repositories found under",*root); os.Exit(1) }
 g:=Graph{SchemaVersion:"1.0",GeneratedAt:time.Now().UTC(),Repositories:repos}
 for _,r:=range repos { if err:=scanRepo(r,&g); err!=nil { g.Warnings=append(g.Warnings,fmt.Sprintf("%s: %v",r,err)) } }
 sort.Slice(g.Nodes,func(i,j int)bool{return g.Nodes[i].ID<g.Nodes[j].ID})
 if err:=os.MkdirAll(*out,0755);err!=nil { fatal(err) }
 b,err:=json.MarshalIndent(g,"","  ");if err!=nil{fatal(err)}
 if err=os.WriteFile(filepath.Join(*out,"architecture.json"),append(b,'\n'),0644);err!=nil{fatal(err)}
 md:=fmt.Sprintf("# ArchGraph scan report\n\n- Repositories: %d\n- Nodes: %d\n- Edges: %d\n- Warnings: %d\n",len(g.Repositories),len(g.Nodes),len(g.Edges),len(g.Warnings))
 if err=os.WriteFile(filepath.Join(*out,"architecture.md"),[]byte(md),0644);err!=nil{fatal(err)}
 fmt.Printf("Scan complete: %d repositories, %d nodes, %d relationships, %d warnings. Output: %s\n",len(repos),len(g.Nodes),len(g.Edges),len(g.Warnings),*out)
 if len(g.Warnings)>0 { os.Exit(3) }
}
func discover(root string) []string {
 abs,err:=filepath.Abs(root);if err!=nil{fatal(err)}
 if exists(filepath.Join(abs,".git")) { return []string{abs} }
 var repos []string
 _=filepath.WalkDir(abs,func(p string,d os.DirEntry,e error)error{
  if e!=nil{return nil};if d.IsDir(){
   n:=d.Name();if n==".git" {return filepath.SkipDir}
   if n=="node_modules"||n=="vendor"||n==".archgraph"||n=="dist"||n=="build"||n==".venv" {return filepath.SkipDir}
   if p!=abs&&exists(filepath.Join(p,".git")){repos=append(repos,p);return filepath.SkipDir}
  };return nil
 })
 sort.Strings(repos);return repos
}
func scanRepo(root string,g *Graph) error {
 repo:=filepath.Base(root);addNode(g,"Repository",repo,root,Evidence{root,"",0,"filesystem","Git repository root"})
 services:=map[string]string{}
 var files []string
 err:=filepath.WalkDir(root,func(p string,d os.DirEntry,e error)error{
  if e!=nil{return nil};if d.IsDir(){n:=d.Name();if n==".git"||n=="node_modules"||n=="vendor"||n==".archgraph"||n=="dist"||n=="build"||n==".venv"{return filepath.SkipDir};return nil}
  files=append(files,p);return nil
 });if err!=nil{return err}
 for _,p:=range files {
  base:=filepath.Base(p);if base=="codeway.yaml"||base=="codeway.yml"{
   data,_:=os.ReadFile(p);name:=""
   for _,line:=range strings.Split(string(data),"\n"){parts:=strings.SplitN(strings.TrimSpace(line),":",2);if len(parts)==2{
    k:=strings.ToLower(strings.TrimSpace(parts[0]));v:=strings.Trim(strings.TrimSpace(parts[1]),"\"'")
    if k=="servicename"||k=="service_name"||k=="name"||k=="lambdaname"||k=="lambda_name"||k=="function_name" {name=v;break}
   }}
   if name==""{g.Warnings=append(g.Warnings,p+": codeway identity not recognized");continue}
   id:=addNode(g,"Service",name,root,Evidence{root,rel(root,p),0,"codeway","declared service identity"});services[filepath.Dir(p)]=id
  }
 }
 for _,p:=range files{
  ext:=strings.ToLower(filepath.Ext(p));rel:=rel(root,p)
  if lang,ok:=codeExt[ext];ok{addNode(g,"SourceFile",rel,root,Evidence{root,rel,0,"file-extension",lang})}
  data,err:=os.ReadFile(p);if err!=nil{continue};text:=string(data)
  for name,re:=range signatures{loc:=re.FindStringIndex(text);if loc==nil{continue}
   kind:="DiscoveryFinding";if strings.HasPrefix(name,"aws_"){kind="AWSOperation"};if strings.HasPrefix(name,"database_"){kind="DatabaseAccess"};if name=="api_route"{kind="APIEndpoint"}
   line:=1+strings.Count(text[:loc[0]],"\n");id:=addNode(g,kind,name,root,Evidence{root,rel,line,"deterministic-pattern","candidate match; requires contextual validation"})
   if service:=nearest(p,services);service!=""{g.Edges=append(g.Edges,Edge{service,id,"CONTAINS_OR_USES","inferred",[]Evidence{{root,rel,line,"path-association","source file is beneath codeway directory"}}})}
  }
  if ext==".tf" { re:=regexp.MustCompile(`(?m)^\s*resource\s+"([^"]+)"\s+"([^"]+)"`);for _,m:=range re.FindAllStringSubmatch(text,-1){addNode(g,"AWSResource",m[2],root,Evidence{root,rel,0,"terraform",m[1]})} }
 }
 return nil
}
func addNode(g *Graph,kind,name,repo string,e Evidence)string{
 id:=fmt.Sprintf("%x",[]byte(fmt.Sprintf("%s|%s|%s",repo,kind,name)))
 for i:=range g.Nodes{if g.Nodes[i].ID==id{g.Nodes[i].Evidence=append(g.Nodes[i].Evidence,e);return id}}
 g.Nodes=append(g.Nodes,Node{id,kind,name,repo,[]Evidence{e}});return id
}
func nearest(p string,m map[string]string)string{best:="";depth:=-1;for d,id:=range m{r,e:=filepath.Rel(d,p);if e==nil&&!strings.HasPrefix(r,"..")&&len(d)>depth{best=id;depth=len(d)}};return best}
func rel(root,p string)string{r,_:=filepath.Rel(root,p);return filepath.ToSlash(r)}
func exists(p string)bool{_,e:=os.Stat(p);return e==nil}
func fatal(e error){fmt.Fprintln(os.Stderr,"archgraph:",e);os.Exit(1)}
