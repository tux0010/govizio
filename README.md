# govizio
Go Library to control Vizio Smartcast TVs.

Example:

```
// Authentication
sc := govizio.NewSmartCast("192.168.1.4", "test123", "govizio example client")
spResp, err := sc.StartPairing()
if err != nil {
	log.Fatal(err)
}

reader := bufio.NewReader(os.Stdin)
fmt.Print("Enter response code: ")
code, _ := reader.ReadString('\n')

r, err := sc.SubmitChallenge(spResp, strings.TrimSpace(code))
if err != nil {
	log.Fatal(err)
}

sc.SetAuthToken(r.AuthToken)


// Toggle Power
err = sc.KeyCommand(govizio.PowerToggle)
if err != nil {
	log.Fatal(err)
}
```

### References
This library is built using the definitions found in https://github.com/exiva/Vizio_SmartCast_API
