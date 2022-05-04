package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

/*
POST /create_event
POST /update_event
POST /delete_event
GET /events_for_day
GET /events_for_week
GET /events_for_month
*/

const (
	ud    = "user_id"
	d     = "date"
	df    = "2006-01-02"
	crEv  = "/create_event"
	upEv  = "/update_event"
	delEv = "/delete_event"
	ed    = "/events_for_day"
	ew    = "/events_for_week"
	em    = "/events_for_month"
)

func dayEq(t1, t2 time.Time) bool {
	return t1.Sub(t2) >= 0 && t1.Sub(t2) <= 24*time.Hour || t1.Sub(t2) <= 0 && t1.Sub(t2) <= -24*time.Hour
}

type event struct {
	UserId int       `json:"user_id"`
	Date   time.Time `json:"date"`
	Uid    int       `json:"uid"`
}

type result struct {
	Events []*event `json:"result"`
}

type errStr struct {
	Err string `json:"error"`
}

type rangeFunc func(*event) bool

func newEvent(userId int, date time.Time, uid int) *event {
	return &event{UserId: userId, Date: date, Uid: uid}
}

type storage map[int]*event

type base struct {
	storage
	number int
	*sync.RWMutex
}

func newBase(number int) *base {
	storage := make(map[int]*event)
	rwmu := new(sync.RWMutex)
	return &base{storage, number, rwmu}
}

func (b *base) createLine(userId int, date time.Time) int {
	b.Lock()
	b.number++
	b.storage[b.number] = newEvent(userId, date, b.number)
	b.Unlock()
	b.RLock()
	x := b.number
	b.RUnlock()
	return x
}

func (b *base) readLine(uid int) (ev *event, ok bool) {
	b.RLock()
	ev, ok = b.storage[uid]
	b.RUnlock()
	return
}

func (b *base) forRange(rf rangeFunc) []*event {
	res := make([]*event, 0)
	b.RLock()
	for _, v := range b.storage {
		if rf(v) {
			res = append(res, v)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Uid < res[j].Uid
	})
	b.RUnlock()
	return res
}

func (b *base) updateLine(uid int, newUserId int, newDate ...time.Time) bool {
	i := 0
	b.RLock()
	_, ok := b.storage[uid]
	b.RUnlock()
	if ok {
		b.Lock()
		if len(newDate) == 1 {
			b.storage[uid].Date = newDate[0]
			i++
		}
		if newUserId > 0 {
			b.storage[uid].UserId = newUserId
			i++
		}
		b.Unlock()
	}

	return i > 0
}

func (b *base) deleteLine(uid int) bool {
	//todo попробовать удалять без uid по комбинации user + date
	b.RLock()
	if _, ok := b.storage[uid]; ok {
		b.RUnlock()
		b.Lock()
		delete(b.storage, uid)
		b.Unlock()
		return true
	}
	b.RUnlock()
	return false
}

func (b *base) create(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	if !postAndFormCheck(r, w) {
		return
	}
	if _, ok := r.Form[ud]; !ok {
		checkError(errors.New("no field "+ud), 400, w)
		return
	}
	if _, ok := r.Form[d]; !ok {
		checkError(errors.New("no field "+d), 400, w)
		return
	}
	userId, err := strconv.Atoi(r.Form.Get(ud))
	if !checkError(err, 400, w) {
		return
	}
	date, err := time.Parse(df, r.Form.Get(d))
	if !checkError(err, 400, w) {
		return
	}
	ev := 0
	if userId >= 0 {
		ev = b.createLine(userId, date)
	} else {
		checkError(errors.New("wrong userId"), 400, w)
		return
	}
	res, err := json.Marshal(result{[]*event{b.storage[ev]}})
	if !checkError(err, 500, w) {
		return
	}

	w.WriteHeader(200)
	_, _ = fmt.Fprintln(w, string(res))
}

func (b *base) update(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	if !postAndFormCheck(r, w) {
		return
	}
	_, ok := r.Form["uid"]
	if !ok {
		checkError(errors.New("wrong format data for update"), 400, w)
		return
	}
	uid, err := strconv.Atoi(r.Form.Get("uid"))
	if !checkError(err, 400, w) {
		return
	}
	_, ok1 := r.Form[ud]
	_, ok2 := r.Form[d]
	if ok1 || ok2 {
		if ok1 && ok2 {
			NewUserId, err1 := strconv.Atoi(r.Form.Get(ud))
			NewDate, err2 := time.Parse(df, r.Form.Get(d))
			if err1 != nil && err2 != nil {
				checkError(errors.New("wrong data for update"), 400, w)
				return
			}
			if err1 != nil {
				if !b.updateLine(uid, -1, NewDate) {
					checkError(errors.New("cant find event by uid"), 400, w)
					return
				}
				re := new(result)
				re.Events = append(re.Events, b.storage[uid])
				res, err := json.Marshal(re)
				if !checkError(err, 500, w) {
					return
				}
				w.WriteHeader(200)
				_, _ = fmt.Fprintln(w, string(res))
				return
			}
			if err2 != nil {
				if !b.updateLine(uid, NewUserId) {
					checkError(errors.New("cant find event by uid"), 400, w)
					return
				}
				re := new(result)
				re.Events = append(re.Events, b.storage[uid])
				res, err := json.Marshal(re)
				if !checkError(err, 500, w) {
					return
				}
				w.WriteHeader(200)
				_, _ = fmt.Fprintln(w, string(res))
				return
			}
			if !b.updateLine(uid, NewUserId, NewDate) {
				checkError(errors.New("cant find event by uid"), 400, w)
				return
			}
			re := new(result)
			re.Events = append(re.Events, b.storage[uid])
			res, err := json.Marshal(re)
			if !checkError(err, 500, w) {
				return
			}
			w.WriteHeader(200)
			_, _ = fmt.Fprintln(w, string(res))
			return
		} else if !ok1 {
			NewDate, err := time.Parse(df, r.Form.Get(d))
			if !checkError(err, 400, w) {
				return
			}
			if !b.updateLine(uid, -1, NewDate) {
				checkError(errors.New("cant find event by uid"), 400, w)
				return
			}
			re := new(result)
			re.Events = append(re.Events, b.storage[uid])
			res, err := json.Marshal(re)
			if !checkError(err, 500, w) {
				return
			}
			w.WriteHeader(200)
			_, _ = fmt.Fprintln(w, string(res))
		} else if !ok2 {
			NewUserId, err := strconv.Atoi(r.Form.Get(ud))
			if !checkError(err, 400, w) {
				return
			}
			if NewUserId <= 0 {
				checkError(errors.New("wrong userId"), 400, w)
				return
			}
			if !b.updateLine(uid, NewUserId) {
				checkError(errors.New("cant find event by uid"), 400, w)
				return
			}
			re := new(result)
			re.Events = append(re.Events, b.storage[uid])
			res, err := json.Marshal(re)
			if !checkError(err, 500, w) {
				return
			}
			w.WriteHeader(200)
			_, _ = fmt.Fprintln(w, string(res))
		}
	} else {
		checkError(errors.New("no data for update"), 400, w)
	}
	/*
		_, ok = r.Form[d]
		if !ok {
			var userId int
			user := r.Form.Get(ud)
			if user == "" {
				checkError(errors.New("userId is wrong, cant update"), 400, w)
				return
			} else {
				userId, err = strconv.Atoi(user)
				if !checkError(err, 400, w) {
					return
				}
			}
			if ch := b.updateLine(uid, userId); ch {
				res, err := json.Marshal(result{[]*event{b.storage[uid]}})
				if !checkError(err, 500, w) {
					return
				}

				w.WriteHeader(200)
				_, _ = fmt.Fprintln(w, string(res))
				return
			} else {
				checkError(errors.New("this uid not found or userId is wrong, cant update"), 400, w)
				return
			}
		} else {
			date, err := time.Parse(df, r.Form.Get(d))
			if !checkError(err, 400, w) {
				return
			}
			var userId int
			user := r.Form.Get(ud)
			if user == "" {
				userId = -1
			} else {
				userId, err = strconv.Atoi(user)
				if !checkError(err, 400, w) {
					return
				}
			}
			if ch := b.updateLine(uid, userId, date); ch {
				res, err := json.Marshal(result{[]*event{b.storage[uid]}})
				if !checkError(err, 500, w) {
					return
				}

				w.WriteHeader(200)
				_, _ = fmt.Fprintln(w, string(res))
				return
			} else {
				checkError(errors.New("this uid not found or userId is wrong, cant update"), 400, w)
				return
			}
		}
	*/
}

func (b *base) readDay(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	if r.Method != "GET" {
		checkError(errors.New("wrong method"), 400, w)
		return
	}
	_, ok := r.URL.Query()[d]
	if !ok {
		checkError(errors.New("cant found events with no date"), 400, w)
		return
	}
	dates := r.URL.Query().Get(d)
	date, err := time.Parse(df, dates)
	if !checkError(err, 400, w) {
		return
	}
	_, ok = r.URL.Query()[ud]
	if !ok {
		ires := b.forRange(func(ev *event) bool {
			return dayEq(date, ev.Date)
		})
		if len(ires) == 0 {
			checkError(errors.New("no event for this date"), 400, w)
			return
		}
		re := new(result)
		re.Events = ires
		res, err := json.Marshal(re)
		if !checkError(err, 500, w) {
			return
		}
		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, string(res))
	} else {
		userId, err := strconv.Atoi(r.URL.Query().Get(ud))
		if !checkError(err, 400, w) {
			return
		}
		ires := b.forRange(func(ev *event) bool {
			return dayEq(date, ev.Date) && ev.UserId == userId
		})
		if len(ires) == 0 {
			checkError(errors.New("no event for this date and userId"), 400, w)
			return
		}
		re := new(result)
		re.Events = ires
		res, err := json.Marshal(re)
		if !checkError(err, 500, w) {
			return
		}
		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, string(res))
	}
}

func (b *base) readWeek(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	if r.Method != "GET" {
		checkError(errors.New("wrong method"), 400, w)
		return
	}
	_, ok := r.URL.Query()[d]
	if !ok {
		checkError(errors.New("cant found events with no date"), 400, w)
		return
	}
	dates := r.URL.Query().Get(d)
	date, err := time.Parse(df, dates)
	if !checkError(err, 400, w) {
		return
	}
	_, ok = r.URL.Query()[ud]
	if !ok {
		ires := b.forRange(func(ev *event) bool {
			for i := 0; i < 7; i++ {
				if dayEq(date.AddDate(0, 0, i), ev.Date) {
					return true
				}
			}
			return false
		})
		if len(ires) == 0 {
			checkError(errors.New("no event for this date"), 400, w)
			return
		}
		re := new(result)
		re.Events = ires
		res, err := json.Marshal(re)
		if !checkError(err, 500, w) {
			return
		}
		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, string(res))
	} else {
		userId, err := strconv.Atoi(r.URL.Query().Get(ud))
		if !checkError(err, 400, w) {
			return
		}
		ires := b.forRange(func(ev *event) bool {
			for i := 0; i < 7; i++ {
				if dayEq(date.AddDate(0, 0, i), ev.Date) && ev.UserId == userId {
					return true
				}
			}
			return false
		})
		if len(ires) == 0 {
			checkError(errors.New("no event for this date and userId"), 400, w)
			return
		}
		re := new(result)
		re.Events = ires
		w.WriteHeader(200)
		res, err := json.Marshal(re)
		if !checkError(err, 500, w) {
			return
		}
		_, _ = fmt.Fprintln(w, string(res))
	}

}

func (b *base) readMonth(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	if r.Method != "GET" {
		checkError(errors.New("wrong method"), 400, w)
		return
	}
	_, ok := r.URL.Query()[d]
	if !ok {
		checkError(errors.New("cant found events with no date"), 400, w)
		return
	}
	dates := r.URL.Query().Get(d)
	date, err := time.Parse(df, dates)
	if !checkError(err, 400, w) {
		return
	}
	_, ok = r.URL.Query()[ud]
	if !ok {
		ires := b.forRange(func(ev *event) bool {
			return date.Month() == ev.Date.Month()
		})
		if len(ires) == 0 {
			checkError(errors.New("no event for this date"), 400, w)
			return
		}
		re := new(result)
		re.Events = ires
		res, err := json.Marshal(re)
		if !checkError(err, 500, w) {
			return
		}
		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, string(res))
	} else {
		userId, err := strconv.Atoi(r.URL.Query().Get(ud))
		if !checkError(err, 400, w) {
			return
		}
		ires := b.forRange(func(ev *event) bool {
			return date.Month() == ev.Date.Month() && ev.UserId == userId
		})
		if len(ires) == 0 {
			checkError(errors.New("no event for this date and userId"), 400, w)
			return
		}
		re := new(result)
		re.Events = ires
		w.WriteHeader(200)
		res, err := json.Marshal(re)
		if !checkError(err, 500, w) {
			return
		}
		_, _ = fmt.Fprintln(w, string(res))
	}
}

func (b *base) delete(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	if !postAndFormCheck(r, w) {
		return
	}
	_, ok := r.Form["uid"]
	if !ok {
		checkError(errors.New("wrong format data for delete"), 400, w)
		return
	}
	uid, err := strconv.Atoi(r.Form.Get("uid"))
	if !checkError(err, 400, w) {
		return
	}
	if b.deleteLine(uid) {
		res, err := json.Marshal(newEvent(0, time.Unix(0, 0), uid))
		if !checkError(err, 500, w) {
			return
		}
		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, string(res))
		return
	} else {
		checkError(errors.New("cant find this uid"), 400, w)
		return
	}
}

func checkError(err error, code int, w http.ResponseWriter) bool {
	if err == nil {
		return true
	}
	w.WriteHeader(code)
	errB, jsonErr := json.Marshal(errStr{err.Error()})
	if jsonErr != nil {
		_, _ = fmt.Fprintln(w, "something wrong with json\n", jsonErr)
	}
	_, _ = fmt.Fprintln(w, string(errB))
	return false
}

func postAndFormCheck(r *http.Request, w http.ResponseWriter) bool {
	if r.Method != "POST" {
		checkError(errors.New("wrong method"), 400, w)
		return false
	}
	err := r.ParseForm()
	if err != nil {
		checkError(err, 503, w)
		return false
	}
	return true
}

func main() {
	P := flag.String("p", "80", "port")
	flag.Parse()
	b := newBase(0)
	http.HandleFunc(crEv, b.create)
	http.HandleFunc(upEv, b.update)
	http.HandleFunc(delEv, b.delete)
	http.HandleFunc(ed, b.readDay)
	http.HandleFunc(ew, b.readWeek)
	http.HandleFunc(em, b.readMonth)
	log.Println("server:", http.ListenAndServe("localhost:"+*P, nil))

}
