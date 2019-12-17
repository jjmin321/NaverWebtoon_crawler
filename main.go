//대상사이트 : 네이버웹툰(https://comic.naver.com)

package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

//스크래핑 대상 URL

const (
	urlRoot    = "https://comic.naver.com/webtoon/weekday.nhn"
	urlSubRoot = "https://comic.naver.com"
)

//동기화를 위한 작업 그룹 선언
var wg sync.WaitGroup
var mutex = &sync.Mutex{}

//class가 title인(모든) 웹툰들 대상으로 원하는 URL파싱 후 반환하는 함수
func parseMainNodes(n *html.Node) bool {
	if n.DataAtom == atom.A && n.Parent != nil {
		return scrape.Attr(n, "class") == "title"
	}
	return false
}

//<td class="title"> 최근 10화 웹툰 사이트들 태그 긁어오는 함수
func parseSubNodes(n *html.Node) bool {
	// return n.Parent.DataAtom == atom.Td && scrape.Attr(n, "class") == "title"
	if n.DataAtom == atom.A && n.Parent != nil {
		return scrape.Attr(n.Parent, "class") == "title"
	}
	return false
}

//<div class="..." id="topTotalStarPoint"> 가 있는 코드들 긁어오는 함수
func parseStarNodes(n *html.Node) bool {
	return n.DataAtom == atom.Dd && scrape.Attr(n, "class") == "total"
}

//에러체크할 함수 선언하여 계속 사용해줌.
func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//최근 10화 웹툰들의 평점과 참여자 다 띄워줘야함.
func scrapRatingType(href, fn string) {
	//작업 종료 알림
	defer wg.Done()

	//GET 방식 요청
	response, err := http.Get(urlSubRoot + href)

	//에러체크
	errCheck(err)

	//코드 읽어오고나서 닫기
	defer response.Body.Close()

	//html 바디 코드 root에 넣기
	root, err := html.Parse(response.Body)

	//에러체크
	errCheck(err)

	//파일 스크림 생성(열기) -> 파일명, 옵션, 권한
	file, err := os.OpenFile("/Users/jejeongmin/documents/go/src/NaverWebtoon_crawler/Scrape/"+fn+".txt", os.O_CREATE|os.O_RDWR, os.FileMode(0777))
	//에러체크
	errCheck(err)

	//메소드 종료 시 파일 닫아야 함
	defer file.Close()

	//쓰기 버퍼 선언
	w := bufio.NewWriter(file)

	for _, g := range scrape.FindAll(root, parseStarNodes) {
		//URL 및 해당 데이터 출력
		fmt.Println("result : ", scrape.Text(g))
		//파싱 데이터 -> 버퍼에 기록
		w.WriteString(scrape.Text(g) + "\r\n")
	}
	w.Flush()
}

//링크랑 웹툰제목
func scrapContents(href, fn string) {
	//작업 종료 알림
	defer wg.Done()

	//GET 방식 요청
	response, err := http.Get(urlSubRoot + href)

	//에러체크
	errCheck(err)

	//response에는 값이 있는데 왜 response.Body에는 값이 없을까?

	//코드 읽어오고나서 닫기
	defer response.Body.Close()

	//html 바디 코드 root에 넣기
	root, err := html.Parse(response.Body)

	//에러체크
	errCheck(err)

	//파일 스크림 생성(열기) -> 파일명, 옵션, 권한
	// file, err := os.OpenFile("/Users/jejeongmin/documents/go/src/NaverWebtoon_crawler/Scrape/"+fn+".txt", os.O_CREATE|os.O_RDWR, os.FileMode(0777))

	//에러체크
	// errCheck(err)

	//메소드 종료 시 파일 닫아야 함
	// defer file.Close()

	//쓰기 버퍼 선언
	// w := bufio.NewWriter(file)

	recentList := scrape.FindAll(root, parseSubNodes)

	//parseSubNodes 함수를 사용해서 원하는 노드 순회(평점 긁어오기)하면서 출력
	//<td class="title"> 최근 10화 웹툰 사이트 태그
	for _, link := range recentList {
		//동기화
		wg.Add(1)
		//해당 데이터 출력
		go scrapRatingType(scrape.Attr(link, "href"), fn)
		//파싱 데이터 -> 버퍼에 기록
		// w.WriteString(scrape.Text(g) + "\r\n")
	}
	// w.Flush()
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

	//웹툰들 URL 추출
	fmt.Println("urlRoot : ", urlRoot)
	fmt.Printf("response.Body : %#v", response.Body)
	fmt.Println("root : ", root)
	urlList := scrape.FindAll(root, parseMainNodes)

	//월요일 웹툰들 전부 긁어와서 for range로 순회
	for _, link := range urlList {
		// //대상 URL 출력
		// fmt.Println("Mon Link : ", link, idx)
		// //웹툰 제목
		// fmt.Println(scrape.Attr(link, "title"))

		fileName := scrape.Attr(link, "title")
		fmt.Println("filename is : ", fileName)

		//작업 대기열에 고루틴이 다 끝날때까지 기다릴 수 있게 추가
		wg.Add(1) //Done 개수와 일치

		//고루틴 시작 -> 작업 대기열 개수와 같아야 함.
		go scrapContents(scrape.Attr(link, "href"), fileName)
	}
	wg.Wait()
	fmt.Println("크롤링 종료")
}
