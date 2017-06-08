package main

import (
    "fmt"
    "github.com/vsco/tail"
    "regexp"
    "errors"
    "encoding/json"
    "os"
    "sync"
    "./index"
    "flag"
    "time"
    "strconv"
)

const (
  PARSE_REGEX_NGINX = `(?P<ipaddress>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) - - \[(?P<dateandtime>\d{2}\/[A-Za-z]{3}\/\d{4}:\d{2}:\d{2}:\d{2} (\+|\-)\d{4})\] ((\"(?P<method>[A-Z]+) )(?P<url>.+) (HTTP\/...")|\"-\") (?P<statuscode>\d+) (?P<bytessent>\d+) (["](?P<referer>(\-)|(.+))["]|\"-\") (["](?P<useragent>.+)["]|\"-\")`
)

func ParseLine(line string) (map[string]interface{}, error) {
    // regex throws away first part of log until the request and ignores everything past status code
    // save request route and status code
    fmt.Println("Parsing: " + line)
    regex := regexp.MustCompile(PARSE_REGEX_NGINX)
    groupNames := regex.SubexpNames()
	matches := regex.FindStringSubmatch(line)
	if matches == nil {
		return nil, errors.New("regex parse failed, skipping line:" + line)
	}

    mapped_matches := make(map[string]interface{})
    for i, match := range matches {
        groupName := groupNames[i]
        switch groupName {
        case "":
            continue
        case "bytessent", "statuscode":
            mapped_matches[groupName], _ = strconv.Atoi(match)
        default:
            mapped_matches[groupName] = match
        }
    }

	return mapped_matches, nil
}

var wg sync.WaitGroup

func main(){
  index.ElasticHost = flag.String("elastichost", "localhost", "The elasticsearch host name.")
  index.ElasticIndex = flag.String("index", "log-shipper", "The elasticsearch index name.")
  flag.Parse()

  files := flag.Args()
  for _, file := range files {
    wg.Add(1)
    go follow(file)
  }
  wg.Wait()
}

const ISO_FORMAT string = "2006-01-02T15:04:05.000-0700"

func follow(file string) {
  defer wg.Done()
  t, _ := tail.TailFile(file, tail.Config{
      Location: &tail.SeekInfo{
        Offset: 0,
        Whence: os.SEEK_END,
      },
      Follow: true,
      ReOpen: true})
  for line := range t.Lines {
      //fmt.Println(line.Text)
      parsed, err := ParseLine(line.Text)
      if err != nil {
        fmt.Println(err.Error())
        continue
      }

      timestamp, err := time.Parse("02/Jan/2006:15:04:05 -0700", parsed["dateandtime"].(string))
      if err == nil {
        timestampString := timestamp.Format(ISO_FORMAT)
        parsed["@timestamp"] = timestampString
      } else {
        fmt.Println("Error parsing date! " + err.Error())
      }
      //parsed["@timestamp"] = time.Now().Format(ISO_FORMAT)

      fmt.Println(parsed["statuscode"])
      jsonString, _ := json.Marshal(parsed)
      fmt.Println(string(jsonString))
      index.Index(parsed)
  }
}
