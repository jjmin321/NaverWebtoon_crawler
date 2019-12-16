//대상사이트 : 네이버웹툰(https://comic.naver.com)

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

//스크래핑 대상 URL

const (
	urlRoot = "https://comic.naver.com/webtoon/weekday.nhn"
)

//에러체크할 함수 선언하여 계속 사용해줌.
func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	//메인 페이지 GET 방식 요청
	response, err := http.Get(urlRoot)
	errCheck(err)

	//항상 들어갔으면 닫아줘야 함.
	defer response.Body.Close()

	//응답 데이터(HTML)
	root, err := html.Parse(response.Body)
	errCheck(err)

	//모든 웹툰들 대상으로 원하는 URL파싱 후 반환하는 함수
	parseMainNodes := func(n *html.Node) bool {
		if n.DataAtom == atom.A && n.Parent != nil {
			fmt.Println("a")
			return scrape.Attr(n, "class") == "title"
		}
		return false
	}

	//웹툰들 URL 추출
	fmt.Println("urlROot : ", urlRoot)
	fmt.Printf("response.Body : %#v", response.Body)
	fmt.Println("root : ", root)
	urlList := scrape.FindAll(root, parseMainNodes)
	fmt.Println(urlList)

	//월요일 웹툰들 전부 긁어와서 for range로 순회
	for idx, link := range urlList {
		//대상 URL 출력
		fmt.Println("c")
		fmt.Println("Mon Link : ", link, idx)
		fmt.Println("TargetURL : ", scrape.Attr(link, "href"))

	}

}
