package networkOperators

//func RunAmericas() {
//	getHtml("https://en.wikipedia.org/wiki/List_of_mobile_network_operators_of_the_Americas")
//}
//
//func RunAsiaPacificRegion() {
//	getHtml("https://en.wikipedia.org/wiki/List_of_mobile_network_operators_of_the_Asia_Pacific_region")
//}
//
//func RunEurope() {
//	getHtml("https://en.wikipedia.org/wiki/List_of_mobile_network_operators_of_Europe")
//}
//
//func RunMiddleEastAndAfrica() {
//	getHtml("https://en.wikipedia.org/wiki/List_of_mobile_network_operators_of_the_Middle_East_and_Africa")
//}
//
//func getHtml(url string) {
//	client := netpoll.HTTP(&tls.Config{InsecureSkipVerify: true})
//	do, _ := httpApi.Do(
//		client, httpApi.Url(url),
//		httpApi.Method(hertzConst.MethodGet),
//		httpApi.Header("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.67 Safari/537.36"),
//	)
//
//	htmlContent := string(do)
//	root, _ := htmlquery.Parse(strings.NewReader(htmlContent))
//	title := htmlquery.Find(root, "//*[@id=\"mw-content-text\"]/div[1]/h2")
//	table := htmlquery.Find(root, "//*[@id=\"mw-content-text\"]/div[1]/table")
//
//	filterTable := make([]*html.Node, 0)
//	for _, node := range table {
//		operatorNodeList := htmlquery.Find(node, "/tbody/tr")
//		if len(operatorNodeList) > 1 {
//			filterTable = append(filterTable, node)
//		}
//	}
//
//	var stringBuilder strings.Builder
//
//	for i, node := range filterTable {
//		operatorNodeList := htmlquery.Find(node, "/tbody/tr")
//		countryNameList := htmlquery.Find(title[i], "/span[1]")
//
//		countryName := htmlquery.InnerText(countryNameList[0])
//		if len(countryName) == 0 {
//			countryNameList = htmlquery.Find(title[i], "/span[2]")
//			countryName = htmlquery.InnerText(countryNameList[0])
//		}
//
//		stringBuilder.SerializeString("\n")
//		stringBuilder.SerializeString(countryName)
//		stringBuilder.SerializeString("||")
//
//		for j, tr := range operatorNodeList {
//			operatorNode := htmlquery.Find(tr, "/td[2]/a")
//			if len(operatorNode) == 0 {
//				operatorNode = htmlquery.Find(tr, "/td[2]")
//			}
//
//			if len(operatorNode) == 0 {
//				continue
//			}
//
//			operatorText := htmlquery.InnerText(operatorNode[0])
//			if len(operatorText) == 0 {
//				continue
//			}
//
//			operatorText = strings.ReplaceAll(operatorText, "\r", "")
//			operatorText = strings.ReplaceAll(operatorText, "\n", "")
//			operatorText = strings.ReplaceAll(operatorText, "\t", "")
//
//			stringBuilder.SerializeString(operatorText)
//			if j < len(operatorNodeList)-1 {
//				stringBuilder.SerializeString("||")
//			}
//		}
//	}
//
//	fmt.Println(stringBuilder.String())
//
//	time.Sleep(time.Second)
//}
//
//func genData() {
//	type networkOperation struct {
//		CountryName     string `json:"countryName"`
//		NetworkOperator string `json:"operator"`
//	}
//
//	f, _ := os.Open("./network_operators.txt")
//	sc := bufio.NewScanner(f)
//
//	totalFormatList := make([][]string, 0)
//
//	for sc.Scan() {
//		text := sc.Text()
//
//		splitText := strings.Split(text, "||")
//		formatList := make([]string, 0)
//
//		for _, v := range splitText {
//			formatList = append(formatList, v)
//			if len(formatList) == 2 {
//				break
//			}
//		}
//
//		if len(formatList) != 2 {
//			fmt.Println(formatList)
//		} else {
//			totalFormatList = append(totalFormatList, formatList)
//		}
//	}
//
//	isoMap := make(map[string]bool)
//	imsiList := make([]networkOperation, 0)
//
//	for _, v := range totalFormatList {
//		iso := v[0]
//
//		if _, ok := isoMap[iso]; !ok {
//			imsiList = append(imsiList, networkOperation{
//				CountryName:     iso,
//				NetworkOperator: v[1],
//			})
//
//			isoMap[iso] = true
//		}
//	}
//
//	fmt.Println(len(imsiList))
//	b, _ := json.Marshal(imsiList)
//	buf := bytes.NewBuffer(make([]byte, 0))
//	w := brotli.NewWriterLevel(buf, brotli.BestCompression)
//	w.Write(b)
//	w.Flush()
//
//	fmt.Println(buf.Bytes())
//}
