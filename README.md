# Naver Webtoon All rating crawler

⭐요일별 모든 웹툰들의 전체 화의 평점을 각 파일별로 나누어 크롤링해오는 프로그램⭐

```go
//스크래핑 대상 URL

const (
	urlRoot    = "https://comic.naver.com/webtoon/weekday.nhn"
	urlSubRoot = "https://comic.naver.com"
)
```


## Stack
|           |     Crawler      |
|:---------:|:---------:|
| Developer | 제정민 | 
| Develop Language | GO |  
| Develop Tool     | Visual Studio Code|

### 프로젝트 진행 단계   
> 1일차 - 🌟네이버 웹툰 메인페이지서 모든 웹툰을 고루틴(쓰레드)를 통해 이동 후 페이지 평점 정보를 가지고 옴.<br/>
> 2일차 - 🤹생각보다 너무 쉬워서 모든 웹툰의 최근 10화를 고루틴(쓰레드)를 통해 2차 이동 후 참여자 수도 가지고 옴.<br/>
> 3일차 - 🐛자꾸 에러가 떠서 로그를 찍어보다가 19금 웹툰에서 막히는 것을 발견하고 고침<br/>
> 4일차 - 👨‍💻마지막으로 파일 입출력 처리에 관한 문제를 발견하고 오류 수정<br/>


### 🧵 모든 웹툰을 고루틴을 통해 이동하는 코드

```go
//class가 title인(모든) 웹툰들 대상으로 원하는 URL파싱 후 반환하는 함수
func parseMainNodes(n *html.Node) bool {
	if n.DataAtom == atom.A && n.Parent != nil {
		return scrape.Attr(n, "class") == "title"
	}
	return false
}
```

```go
urlList := scrape.FindAll(root, parseMainNodes)

	//웹툰들 전부 긁어와서 for range로 순회
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
```
<img width="1440" alt="Screen Shot 2019-12-19 at 20 26 54" src="https://user-images.githubusercontent.com/52072077/71170170-f131bd00-229d-11ea-935c-0daf82572cc9.png">

### ✅ 고루틴이 있다면 동기화가 필수
```go
//동기화를 위한 작업 그룹 선언
var wg sync.WaitGroup
var mutex = &sync.Mutex{}
```

### 👉 GO언어를 사용한다면 에러체크도 필수 
```go
//에러체크할 함수 선언하여 계속 사용해줌.
func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
```

### 🧵 모든 웹툰의 최근 10화 웹툰을 고루틴을 통해 들어가는 코드
```go
//<td class="title"> 최근 10화 웹툰 사이트들 태그 긁어오는 함수
func parseSubNodes(n *html.Node) bool {
	// return n.Parent.DataAtom == atom.Td && scrape.Attr(n, "class") == "title"
	if n.DataAtom == atom.A && n.Parent != nil {
		return scrape.Attr(n.Parent, "class") == "title"
	}
	return false
}
```

```go
recentList := scrape.FindAll(root, parseSubNodes)

	//parseSubNodes 함수를 사용해서 원하는 노드 순회(평점 긁어오기)하면서 출력
	//<td class="title"> 최근 10화 웹툰 사이트 태그
	for _, link := range recentList {
		//동기화
		wg.Add(1)
		mutex.Lock()
		//해당 데이터 출력
		go scrapRatingType(scrape.Attr(link, "href"), fn)
		//파싱 데이터 -> 버퍼에 기록
		// w.WriteString(scrape.Text(g) + "\r\n")
	}
	// w.Flush()
```

<img width="1440" alt="Screen Shot 2019-12-19 at 20 26 19" src="https://user-images.githubusercontent.com/52072077/71170174-f2fb8080-229d-11ea-9d89-ad3484194c3d.png">

### 👨‍💻 웹툰들의 평점과 참여자 수 정보를 가지고오는 코드

```go
//<div class="..." id="topTotalStarPoint"> 가 있는 코드들 긁어오는 함수
func parseStarNodes(n *html.Node) bool {
	return n.DataAtom == atom.Div && scrape.Attr(n, "id") == "topTotalStarPoint"
}
```

<img width="1440" alt="Screen Shot 2019-12-19 at 20 26 41" src="https://user-images.githubusercontent.com/52072077/71170177-f4c54400-229d-11ea-8e2f-49b2800e23df.png">

### 🗒️ 웹툰이름.txt 파일에 평점과 참여자 수를 자동으로 작성해주는 코드
```go
//파일 스크림 생성(열기) -> 파일명, 옵션, 권한
	file, err := os.OpenFile("/Users/jejeongmin/documents/go/src/NaverWebtoon_crawler/Scrape/"+fn+".txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.FileMode(0777))
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
```


<img width="1440" alt="Screen Shot 2019-12-19 at 19 59 29" src="https://user-images.githubusercontent.com/52072077/71168499-318f3c00-229a-11ea-9004-84e12504b3dc.png">



