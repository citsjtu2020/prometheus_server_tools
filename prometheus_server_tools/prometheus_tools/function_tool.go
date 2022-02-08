package prometheus_tools

//func Check(err error){
//	if err != nil{
//		panic(err.Error())
//	}
//}
//
//func Do(url string) (string,error){
//	defer func() {
//		if err := recover(); err != nil{
//			fmt.Println(err)
//		}
//	}()
//	body, err := Get(url)
//	//Check(err)
//	if err != nil{
//		return "", err
//	}
//	plaintext := string(body)
//	return plaintext,nil
//}
//
//func Get(url string) ([]byte,error){
//	defer func() {
//		if err := recover(); err != nil{
//			fmt.Println(err)
//		}
//	}()
//	resp,err := http.Get(url)
//	Check(err)
//	if resp.StatusCode != 200{
//		return []byte(""),fmt.Errorf("status wrong: %v",resp.StatusCode)
//	}
//	return ioutil.ReadAll(resp.Body)
//}

//func main() {
//
//}