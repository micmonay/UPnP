# UPnP in go

Easy for dialog with a UPnP device.

> (Sorry if my english is bad in the project or the examples, I try)

I built a tool at    [https://github.com/micmonay/UpnpTools](https://github.com/micmonay/UpnpTools)

This exemple is for get the IPV4 in your router. For exemple your router has duty a UPnP active

```go
up := upnp.NewUPNP(upnp.SERVICE_GATEWAY_IPV4_V2) // or upnp.SERVICE_GATEWAY_IPV4_V1
Interface, err := upnp.GetInterfaceByName("eth0")
if err != nil {
	log.Println(err)
	return
}
// Get all devices compatible for the service name (timeout 1 second)
devices := up.GetAllCompatibleDevice(Interface, 1)
if len(devices) == 0 {
	return
}
// Get services
services := devices[0].GetServicesByType(upnp.SERVICE_GATEWAY_IPV4_V2) // or upnp.SERVICE_GATEWAY_IPV4_V1
if len(services) == 0 {
	return
}
// if you have a one routeur it's ok, other ...
service := services[0]
if service == nil {
	log.Println("not found service")
	return
}
// send request
response, err := service.GetAction("GetExternalIPAddress").Send()
if err != nil {
	log.Println(err)
	return
}
// get response argument
ip, err := response.GetValueArgument("NewExternalIPAddress")
if err != nil {
	log.Println(err)
	return
}
	fmt.Println("Your WAN ip address is " + ip)
```

If your action containe arguments at sending.

```go
Action := Services[0].GetAction("ActionName")
Action.AddVariable("nameOfArgument","value")
Action.Send()
```
