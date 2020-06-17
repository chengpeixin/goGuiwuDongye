package main

import (
	"fmt"
	"sync"
	"github.com/PuerkitoBio/goquery"
	"time"

	//"golang.org/x/net/html"
	"net/http"
	//"net/url"
)
//初始url
const startUrl string = "http://dongyeguiwu.zuopinj.com/index.html"
var wg sync.WaitGroup
// 共多少本书，共多少章
var bookNum int = 0
var bookChapterNum int =0
//协程的数量
var syncNum int =0

//首先获取所有的书的地址
var bookAllLink = []string{}

//获取所有的书籍
func getBookUriList (uri string,linkList []string) []string {
	doc := getDocument(uri)
	doc.Find(".tab-detail").First().Find(".zp-book-item").Each(func (i int,s *goquery.Selection){
		href,_ := s.Find("a").Attr("href")
		linkList = append(linkList, href)
	})
	nextPageHref,exists := doc.Find("#lg_nextpage").Attr("href")
	if exists==false {
		return linkList
	}else {
		return getBookUriList(nextPageHref,linkList)
	}
}

// 获取章节
func getBookContent (uri string){
	doc := getDocument(uri)
	doc.Find(".book_list li").Each(func(i int, selection *goquery.Selection) {
		bookChapterNum++
		//chapterTitle := selection.Find("a").Text()
		chapterHref,_ := selection.Find("a").Attr("href")
		wg.Add(1)
		syncNum++
		go getBookChapterContent(chapterHref)
	})
	wg.Done()
}

// 获取所有书
func getAllBook(uriList []string){
	// 开启协程去获取书的内容
	for i:=0;i<len(uriList);i++{
		bookHref := uriList[i]
		wg.Add(1)
		syncNum++
		go getBookContent(bookHref)
	}
}


// 获取章节的内容
func getBookChapterContent(uri string){
	doc := getDocument(uri)
	srcbox := doc.Find("title").Text()
	fmt.Println(srcbox)
	wg.Done()
}

// 获取一个经过goquery process的文档
func getDocument(uri string) goquery.Document {
	client := &http.Client{}
	req,err := http.NewRequest("GET",uri,nil)
	if err != nil{
		fmt.Println(err)
	}
	req.Header.Add("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
	resp,err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	doc,err := goquery.NewDocumentFromReader(resp.Body)
	if err !=nil {
		fmt.Println(err)
	}
	return *doc
}
// 统计并打印耗时和数据
func statisticsData(startTime time.Time){
	fmt.Printf("共爬取%d本书,共%d章，共开启%d个协程,耗时%s\n",bookNum,bookChapterNum,syncNum,time.Since(startTime))
}
func main(){
	startTime := time.Now()
	bookListHref := getBookUriList(startUrl,bookAllLink)
	bookNum = len(bookListHref)
	// 获取书的章节和内容
	getAllBook(bookListHref)
	wg.Wait()
	statisticsData(startTime)
}