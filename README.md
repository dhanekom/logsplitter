# logsplitter
A small library that splits log files into its separate parts.

# Installation
```
$ go get https://github.com/dhanekom/logsplitter
```

# Usage
Given a log file with the following content:
```
2021/08/28 19:41:15.740|ERROR|3|2222|An error message
2021/08/30 19:41:15.740|INFO|2|1553|A log message
```
the following program will produce the following result
```
2021/08/28 19:41:15.740
ERROR
3
2222
An error message
-----
2021/08/30 19:41:15.740
INFO
2
1553
Starting SKU Refresh
```


Example code:
```
// Open a file
file, err := os.Open("test.log")
if err != nil {
  log.Fatal(err)
}
defer file.Close()

regexPattern := "^(.+)\|(.+)\|(.+)\|(.+)\|(.+)$"

// Create a string splitter
stringSplitter, err := logsplitter.NewRegexStringSplitter(regexPattern)
if err != nil {
  log.Fatal(err)
}

parser := logsplitter.NewParser(stringSplitter)
reader := logsplitter.NewParseReader(file, parser)

for {
  // Read log a log file line by line and return a logsplitter.Fields
  // fields is a slice of logsplitter.Field
  fields, err := reader.Read()
  if err != nil {
    if err == io.EOF {
      break
    }

    log.Fatal(err)
  }
	
  // Use the fields slice
  for _, field := range fields {
    fmt.Println(field.Value)
  }
}
```