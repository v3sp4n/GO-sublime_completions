package main

import (
	"fmt"
	"net/http"
	"io"
	"strings"
	"regexp"
	"encoding/json"
	"io/ioutil"
)



func main() {

	var url string
	var trigger string

	fmt.Println("put url here(example:https://pkg.go.dev/os)")
	fmt.Scan(&url)
	fmt.Println("input trigger(example:os.->os.Readlink)")
	fmt.Scan(&trigger)

	fmt.Println("url",url)
	fmt.Println("trigger",trigger)

	result := []map[string]string{}
	// result = append(result,map[string]string{})

	r,err := http.Get(url)

	if err != nil {
		fmt.Println("ERR>",err)
	}	
	body,err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("ERR>",err)
	}
	for _,v := range strings.Split(string(body),"\n") {
		find,err := regexp.Match("<pre>func.+</pre>",[]byte(v))
		if find && err == nil {
			re := regexp.MustCompile("<pre>(.+)</pre>")
			m := re.FindStringSubmatch(v)
			if len(m) >= 1 && len(regexp.MustCompile(`func \w+\(`).FindStringSubmatch(m[1])) >= 1 {

				funca := regexp.MustCompile(`<a href=\"/\S+#\S+\">`).ReplaceAllString(m[1],"")
				funca = regexp.MustCompile(`<a href=\"/\S+\">`).ReplaceAllString(m[1],"")
				funca = strings.ReplaceAll(funca,"</a>","")

				retSplit := strings.Split(funca,")")

				var funcaReturn string
				if len(retSplit[len(retSplit)-1]) >= 2 {
				    funcaReturn = retSplit[len(retSplit)-1]
				} else {
				    funcaReturn = retSplit[len(retSplit)-2] + ")"
				}
				funcaReturn = strings.Trim(funcaReturn," ")
				funca = strings.ReplaceAll(funca,funcaReturn,"")
				args := ""

				a := regexp.MustCompile(`(\S+)\s*\((.+)\)`).FindStringSubmatch(funca)
				if len(a) >= 1{
					args = a[1]
					funca = strings.ReplaceAll(funca,a[0],a[1]+"(${1:"+a[2]+"})")
				}
				// fmt.Println(k,funca,"RETURN:",funcaReturn)

				triggerMatch := regexp.MustCompile(`^(\S+.\S+)\(`).FindStringSubmatch(strings.ReplaceAll(trigger+funca,"func ",""))
				if len(triggerMatch) >= 1 {
					fmt.Println("add",funca)
					result = append(result,map[string]string{
						"trigger": triggerMatch[1],
						"contents": funca,
						"annotation": args,//func args
						"details": "return " + funcaReturn,//func return
					})
				} else {
					fmt.Println("triggerMatch==0",funca,args,funcaReturn)
				}

			} else {
				fmt.Println("ERROR MATCH",v)
			}
		}
	}
	j,err := json.Marshal(result)
	if err != nil {
		fmt.Println("ERR>",err)
	}	
	fmt.Println("\n\n\n")
	fmt.Printf(string(j))
	fmt.Println("\n\n\n")

	ioutil.WriteFile(strings.ReplaceAll(trigger,".","")+".json", j, 0777)

	fmt.Println("check " + strings.ReplaceAll(trigger,".","")+".json" + " file")

	var stop string
	fmt.Scan(&a)
}