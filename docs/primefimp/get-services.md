

topic: pt:j1/mt:cmd/rt:app/rn:vinculum/ad:1
 
```json
{
  "ctime": "2019-05-31 17:36:31 +0200",
  "props": {},
  "resp_to": "pt:j1/mt:rsp/rt:app/rn:pf-ex-client/ad:r8117wyq",
  "serv": "vinculum",
  "src": "prime-fimp-example-client",
  "tags": [],
  "type": "cmd.pd7.request",
  "uid": "09936904-8ec8-4280-9285-945d77bd0a74",
  "val": {
    "cmd": "get",
    "component": null,
    "id": null,
    "param": {
      "components": [
        "service"
      ]
    },
    "requestId": 7294000000007
  },
  "val_t": "object",
  "ver": "1"
} 
```

Reset smarthub : 

```json
{
  "ctime": "2019-05-31 17:36:31 +0200",
  "props": {},
  "resp_to": "pt:j1/mt:rsp/rt:cloud/rn:backend-service/ad:tf-factory-hub-reset",
  "serv": "vinculum",
  "src": "prime-fimp-example-client",
  "tags": [],
  "type": "cmd.pd7.request",
  "uid": "09936904-8ec8-4280-9285-945d77bd0a74",
  "val": {
    "cmd": "delete", 
    "component": "config", 
    "id": null, 
    "param": {}, 
    "requestId": 1 
  },
  "val_t": "object",
  "ver": "1"
} 


```

response 

topic: pt:j1/mt:rsp/rt:app/rn:pf-ex-client/ad:r8117wyq

```json
{
  "corid": "09936904-8ec8-4280-9285-945d77bd0a74",
  "ctime": "2019-07-03T13:31:22+0200",
  "props": {},
  "serv": "vinculum",
  "tags": [],
  "type": "evt.pd7.response",
  "uid": "44ee61c9-4b67-4c10-afd3-279d3c6edc91",
  "val": {
    "errors": null,
    "param": {
      "service": {
        "fireAlarm": {
          "alarmState": "ready",
          "appliances": {
            "24": {
              "power": "off"
            },
            "45": {
              "power": "off"
            },
            "49": {
              "power": "off"
            },
            "89": {
              "power": "off"
            },
            "90": {
              "power": "off"
            },
            "91": {
              "power": "on"
            },
            "92": {
              "power": "off"
            }
          },
          "areas": {},
          "confirmed": false,
          "delay": 120,
          "enabled": true,
          "lastTested": "2018-06-27T09:27:14Z",
          "lastTriggered": "2019-06-18T07:11:12Z",
          "lights": {
            "64": {
              "power": "off"
            },
            "113": {
              "power": "on"
            },
            "114": {
              "power": "on"
            }
          },
          "locks": {
            "40": {
              "lockState": null
            },
            "74": {
              "lockState": null
            }
          },
          "sirens": {
            "73": {
              "siren": null
            },
            "99": {
              "siren": null
            },
            "111": {
              "siren": null
            }
          },
          "supported": true,
          "triggeredDevices": [],
          "triggeredRooms": []
        }
      }
    },
    "requestId": 7294000000007,
    "success": true
  },
  "val_t": "object",
  "ver": "1"
}

```

