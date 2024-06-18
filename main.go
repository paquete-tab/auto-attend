package main

import (
  "fmt"
  "time"
  "github.com/go-rod/rod"
  "github.com/go-rod/rod/lib/input"
)

type Course struct{
  id string
  weekday time.Weekday
  period int
}

func main() {
  loginURL := "https://moodle.s.kyushu-u.ac.jp/login/index.php"
  username := "USERNAME"
  password := "PASSWORD"
  courses := []*Course{
    {"12345", time.Monday, 1}, //月曜1限

    {"56789", time.Wednesday, 3}, //水曜3,4限
    {"56789", time.Wednesday, 4},
  }

  periodtime := [][][]int{
    {{8,40},{10,10}}, //1限
    {{10,30},{12,0}}, //2限
    {{13,0},{14,30}}, //3限
    {{14,50},{16,20}}, //4限
    {{16,40},{18,10}}, //5限
  }
  now := time.Now()
  current := -1 //授業時間外の場合

  //現在の時限
  for id, t := range periodtime{
    start := time.Date(now.Year(), now.Month(), now.Day(), t[0][0], t[0][1], now.Second(), now.Nanosecond(), now.Location())
    end := time.Date(now.Year(), now.Month(), now.Day(), t[1][0], t[1][1], now.Second(), now.Nanosecond(), now.Location())

    if start.Unix() <= now.Unix() && end.Unix() >= now.Unix(){
      current = id + 1
    }
  }

  //fmt.Println(current)
  
  //ログイン処理
  browser := rod.New().MustConnect()
  defer browser.MustClose()
  page := browser.MustPage(loginURL).MustWaitStable()
  page.MustElement("#username").MustInput(username)
  page.MustElement("#password").MustInput(password).MustType(input.Enter)
  page.MustWaitStable()

  //出席
  for _, c := range courses{
    if c.weekday == now.Weekday() && c.period == current{
      courseURL := "https://moodle.s.kyushu-u.ac.jp/course/view.php?id=" + c.id
      page := browser.MustPage(courseURL)
      page.MustWaitStable()
      page.MustScreenshot("moodle.png")
      fmt.Println(courseURL)
    }
  }
}