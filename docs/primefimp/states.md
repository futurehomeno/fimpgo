
`pt:j1/mt:cmd/rt:app/rn:vinculum/ad:1`

```json
{
  "corid": "",
  "ctime": "2020-02-18T19:28:13.691351",
  "props": null,
  "serv": "vinculum",
  "tags": null,
  "type": "cmd.pd7.request",
  "uid": "6c6fccb0-527c-11ea-b474-a7d82dbdcf44",
  "val_t": "object",
  "ver": "1",
  "val": {
    "cmd": "get",
    "component": null,
    "id": null,
    "client": null,
    "param": {
      "components": [
        "state"
      ]
    },
    "requestId": "158205049369172"
  },
  "resp_to": "pt:j1/mt:rsp/rt:cloud/rn:remote-client/ad:smarthome-app",
  "src": "app"
}
```


Response : 

`pt:j1/mt:rsp/rt:cloud/rn:remote-client/ad:smarthome-app`

```json
{
  "corid": "6c6fccb0-527c-11ea-b474-a7d82dbdcf44",
  "ctime": "2020-02-18T20:03:07+0100",
  "props": {},
  "serv": "vinculum",
  "tags": [],
  "type": "evt.pd7.response",
  "uid": "eddf3175-435c-484c-96b1-35d5df4941f4",
  "val": {
    "errors": null,
    "param": {
      "state": {
        "devices": [
          {
            "id": 10,
            "services": []
          },
          {
            "id": 11,
            "services": [
              {
                "addr": "/rt:dev/rn:virtual/ad:1/sv:sensor_temp/ad:yr_report",
                "attributes": [
                  {
                    "name": "sensor",
                    "values": [
                      {
                        "props": {
                          "unit": "C"
                        },
                        "ts": "2020-02-18 19:57:32 +0100",
                        "val": 5.9,
                        "val_t": "float"
                      }
                    ]
                  }
                ],
                "name": "sensor_temp"
              }
            ]
          },
          {
            "id": 12,
            "services": []
          },
          {
            "id": 13,
            "services": []
          },
          {
            "id": 15,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:16_0",
                "attributes": [
                  {
                    "name": "state",
                    "values": [
                      {
                        "ts": "2020-01-29 12:18:08 +0100",
                        "val": "UP",
                        "val_t": "string"
                      }
                    ]
                  },
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NO_ACK",
                          "src": "nodeId=16_0;service=battery"
                        },
                        "ts": "2020-02-11 15:35:28 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 23,
            "services": []
          },
          {
            "id": 24,
            "services": []
          },
          {
            "id": 35,
            "services": []
          },
          {
            "id": 36,
            "services": []
          },
          {
            "id": 37,
            "services": []
          },
          {
            "id": 38,
            "services": []
          },
          {
            "id": 39,
            "services": []
          },
          {
            "id": 45,
            "services": [
              {
                "addr": "/rt:dev/rn:flow/ad:1/sv:sensor_contact/ad:3YBhqMcMkf47G0F_0",
                "attributes": [
                  {
                    "name": "open",
                    "values": [
                      {
                        "ts": "2020-02-12 13:17:00 +0100",
                        "val": true,
                        "val_t": "bool"
                      }
                    ]
                  }
                ],
                "name": "sensor_contact"
              },
              {
                "addr": "/rt:dev/rn:flow/ad:1/sv:out_bin_switch/ad:3YBhqMcMkf47G0F_0",
                "attributes": [
                  {
                    "name": "binary",
                    "values": [
                      {
                        "ts": "2020-02-12 13:17:00 +0100",
                        "val": true,
                        "val_t": "bool"
                      }
                    ]
                  }
                ],
                "name": "out_bin_switch"
              }
            ]
          },
          {
            "id": 49,
            "services": [
              {
                "addr": "/rt:dev/rn:flow/ad:1/sv:alarm_fire/ad:1A9tYzUMoI_0",
                "attributes": [
                  {
                    "name": "alarm",
                    "values": [
                      {
                        "ts": "2020-02-10 12:12:50 +0100",
                        "val": {
                          "event": "smoke",
                          "status": "deactiv"
                        },
                        "val_t": "str_map"
                      }
                    ]
                  }
                ],
                "name": "alarm_fire"
              }
            ]
          },
          {
            "id": 56,
            "services": []
          },
          {
            "id": 63,
            "services": []
          },
          {
            "id": 70,
            "services": [
              {
                "addr": "/rt:dev/rn:flow/ad:1/sv:sensor_temp/ad:uFTrELNx4JPpuSt_0",
                "attributes": [
                  {
                    "name": "sensor",
                    "values": [
                      {
                        "props": {
                          "unit": "C"
                        },
                        "ts": "2020-02-18 19:57:33 +0100",
                        "val": 41.698,
                        "val_t": "float"
                      }
                    ]
                  }
                ],
                "name": "sensor_temp"
              }
            ]
          },
          {
            "id": 73,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:86_0",
                "attributes": [
                  {
                    "name": "state",
                    "values": [
                      {
                        "ts": "2020-01-30 07:23:39 +0100",
                        "val": "DOWN",
                        "val_t": "string"
                      }
                    ]
                  },
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NOROUTE",
                          "src": "nodeId=86_0;service=scene_ctrl"
                        },
                        "ts": "2020-01-18 17:32:17 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 74,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:door_lock/ad:92_0",
                "attributes": [
                  {
                    "name": "lock",
                    "values": [
                      {
                        "props": {
                          "timeout_s": "254",
                          "unsecured_desc": ""
                        },
                        "ts": "2020-02-11 15:38:12 +0100",
                        "val": {
                          "bolt_is_locked": false,
                          "door_is_closed": false,
                          "is_secured": false,
                          "latch_is_closed": false
                        },
                        "val_t": "bool_map"
                      }
                    ]
                  }
                ],
                "name": "door_lock"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:92_0",
                "attributes": [
                  {
                    "name": "state",
                    "values": [
                      {
                        "ts": "2020-01-29 17:35:54 +0100",
                        "val": "DOWN",
                        "val_t": "string"
                      }
                    ]
                  },
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NO_ACK",
                          "src": "nodeId=92_0;service=door_lock"
                        },
                        "ts": "2020-01-23 14:58:07 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 76,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:thermostat/ad:98_0",
                "attributes": [
                  {
                    "name": "mode",
                    "values": [
                      {
                        "ts": "2020-02-11 15:36:24 +0100",
                        "val": "heat",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "thermostat"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:98_0",
                "attributes": [
                  {
                    "name": "state",
                    "values": [
                      {
                        "ts": "2020-01-29 17:36:34 +0100",
                        "val": "DOWN",
                        "val_t": "string"
                      }
                    ]
                  },
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NOROUTE",
                          "src": "nodeId=98_0;service=thermostat"
                        },
                        "ts": "2020-02-11 15:36:54 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:meter_elec/ad:98_0",
                "attributes": [
                  {
                    "name": "meter",
                    "values": [
                      {
                        "props": {
                          "unit": "kWh"
                        },
                        "ts": "2020-02-11 15:36:27 +0100",
                        "val": 51.9000015258789,
                        "val_t": "float"
                      }
                    ]
                  }
                ],
                "name": "meter_elec"
              }
            ]
          },
          {
            "id": 77,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:thermostat/ad:98_1",
                "attributes": [
                  {
                    "name": "mode",
                    "values": [
                      {
                        "ts": "2020-02-11 15:36:35 +0100",
                        "val": "heat",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "thermostat"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:98_1",
                "attributes": [
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NOROUTE",
                          "src": "nodeId=98_1;service=thermostat"
                        },
                        "ts": "2020-02-11 15:36:39 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 78,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:98_2",
                "attributes": [
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NOROUTE",
                          "src": "nodeId=98_2;service=sensor_temp"
                        },
                        "ts": "2020-02-11 15:36:44 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 79,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:98_3",
                "attributes": [
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NOROUTE",
                          "src": "nodeId=98_3;service=sensor_temp"
                        },
                        "ts": "2020-02-11 15:36:50 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 80,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:meter_elec/ad:98_4",
                "attributes": [
                  {
                    "name": "meter",
                    "values": [
                      {
                        "props": {
                          "unit": "kWh"
                        },
                        "ts": "2020-02-11 15:36:58 +0100",
                        "val": 51.9000015258789,
                        "val_t": "float"
                      }
                    ]
                  }
                ],
                "name": "meter_elec"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:98_4",
                "attributes": [
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NOROUTE",
                          "src": "nodeId=98_4;service=meter_elec"
                        },
                        "ts": "2020-02-11 15:37:03 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 81,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:meter_elec/ad:100_0",
                "attributes": [
                  {
                    "name": "meter",
                    "values": [
                      {
                        "props": {
                          "unit": "kWh"
                        },
                        "ts": "2019-09-27 12:36:18 +0200",
                        "val": 0,
                        "val_t": "float"
                      }
                    ]
                  }
                ],
                "name": "meter_elec"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:100_0",
                "attributes": [
                  {
                    "name": "state",
                    "values": [
                      {
                        "ts": "2019-09-19 21:18:58 +0200",
                        "val": "DOWN",
                        "val_t": "string"
                      }
                    ]
                  },
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NOROUTE",
                          "src": "nodeId=100_0;service=out_bin_switch"
                        },
                        "ts": "2019-09-27 12:36:18 +0200",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 82,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:meter_elec/ad:100_1",
                "attributes": [
                  {
                    "name": "meter",
                    "values": [
                      {
                        "props": {
                          "unit": "kWh"
                        },
                        "ts": "2019-09-27 12:36:06 +0200",
                        "val": 0,
                        "val_t": "float"
                      }
                    ]
                  }
                ],
                "name": "meter_elec"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:100_1",
                "attributes": [
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NOROUTE",
                          "src": "nodeId=100_1;service=out_bin_switch"
                        },
                        "ts": "2019-09-27 12:36:18 +0200",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 83,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:100_2",
                "attributes": [
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NOROUTE",
                          "src": "nodeId=100_2;service=sensor_temp"
                        },
                        "ts": "2019-09-27 12:36:18 +0200",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 84,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:100_3",
                "attributes": [
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NOROUTE",
                          "src": "nodeId=100_3;service=sensor_temp"
                        },
                        "ts": "2019-09-27 12:36:17 +0200",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 85,
            "services": []
          },
          {
            "id": 89,
            "services": []
          },
          {
            "id": 90,
            "services": []
          },
          {
            "id": 91,
            "services": []
          },
          {
            "id": 93,
            "services": []
          },
          {
            "id": 94,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:103_0",
                "attributes": [
                  {
                    "name": "state",
                    "values": [
                      {
                        "ts": "2020-01-29 18:07:30 +0100",
                        "val": "DOWN",
                        "val_t": "string"
                      }
                    ]
                  },
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NOROUTE",
                          "src": "nodeId=103_0;service=battery"
                        },
                        "ts": "2020-02-11 15:37:36 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 95,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:104_0",
                "attributes": [
                  {
                    "name": "state",
                    "values": [
                      {
                        "ts": "2020-02-07 23:53:39 +0100",
                        "val": "UP",
                        "val_t": "string"
                      }
                    ]
                  },
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NOROUTE",
                          "src": "nodeId=104_0;service=battery"
                        },
                        "ts": "2020-02-11 15:37:17 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 96,
            "services": []
          },
          {
            "id": 97,
            "services": []
          },
          {
            "id": 99,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:111_0",
                "attributes": [
                  {
                    "name": "state",
                    "values": [
                      {
                        "ts": "2020-01-29 17:37:54 +0100",
                        "val": "DOWN",
                        "val_t": "string"
                      }
                    ]
                  },
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NO_ACK",
                          "src": "nodeId=111_0;service=sensor_humid"
                        },
                        "ts": "2020-02-11 15:40:27 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 103,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:118_0",
                "attributes": [
                  {
                    "name": "state",
                    "values": [
                      {
                        "ts": "2020-01-29 18:07:31 +0100",
                        "val": "DOWN",
                        "val_t": "string"
                      }
                    ]
                  },
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "ZW_ERR_OP_FAILED",
                          "src": "nodeId=118_0;service=battery"
                        },
                        "ts": "2020-02-11 15:37:39 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 105,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:meter_elec/ad:120_0",
                "attributes": [
                  {
                    "name": "meter",
                    "values": [
                      {
                        "props": {
                          "unit": "kWh"
                        },
                        "ts": "2020-02-11 15:36:07 +0100",
                        "val": 0,
                        "val_t": "float"
                      }
                    ]
                  }
                ],
                "name": "meter_elec"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:120_0",
                "attributes": [
                  {
                    "name": "state",
                    "values": [
                      {
                        "ts": "2020-01-29 17:38:34 +0100",
                        "val": "DOWN",
                        "val_t": "string"
                      }
                    ]
                  },
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NO_ACK",
                          "src": "nodeId=120_0;service=meter_elec"
                        },
                        "ts": "2020-02-11 15:36:12 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 106,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:meter_elec/ad:120_1",
                "attributes": [
                  {
                    "name": "meter",
                    "values": [
                      {
                        "props": {
                          "unit": "kWh"
                        },
                        "ts": "2020-02-11 15:36:18 +0100",
                        "val": 0,
                        "val_t": "float"
                      }
                    ]
                  }
                ],
                "name": "meter_elec"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:120_1",
                "attributes": [
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NO_ACK",
                          "src": "nodeId=120_1;service=meter_elec"
                        },
                        "ts": "2020-02-11 15:36:24 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 107,
            "services": []
          },
          {
            "id": 108,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:121_0",
                "attributes": [
                  {
                    "name": "state",
                    "values": [
                      {
                        "ts": "2020-01-29 17:39:14 +0100",
                        "val": "DOWN",
                        "val_t": "string"
                      }
                    ]
                  },
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_ROUTING_NOT_IDLE",
                          "src": "nodeId=121_0;service=out_bin_switch"
                        },
                        "ts": "2020-02-10 12:12:48 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 109,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:122_0",
                "attributes": [
                  {
                    "name": "state",
                    "values": [
                      {
                        "ts": "2020-02-08 00:30:41 +0100",
                        "val": "UP",
                        "val_t": "string"
                      }
                    ]
                  },
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NO_ACK",
                          "src": "nodeId=122_0;service=battery"
                        },
                        "ts": "2020-02-11 15:37:34 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:sensor_temp/ad:122_0",
                "attributes": [
                  {
                    "name": "sensor",
                    "values": [
                      {
                        "props": {
                          "unit": "C"
                        },
                        "ts": "2020-02-12 09:23:58 +0100",
                        "val": 16.4400005340576,
                        "val_t": "float"
                      }
                    ]
                  }
                ],
                "name": "sensor_temp"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:alarm_burglar/ad:122_0",
                "attributes": [
                  {
                    "name": "alarm",
                    "values": [
                      {
                        "ts": "2020-02-10 09:43:26 +0100",
                        "val": {
                          "event": "tamper_removed_cover",
                          "status": "activ"
                        },
                        "val_t": "str_map"
                      }
                    ]
                  }
                ],
                "name": "alarm_burglar"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:battery/ad:122_0",
                "attributes": [
                  {
                    "name": "lvl",
                    "values": [
                      {
                        "ts": "2020-01-13 03:45:31 +0100",
                        "val": 53,
                        "val_t": "int"
                      }
                    ]
                  }
                ],
                "name": "battery"
              }
            ]
          },
          {
            "id": 111,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:124_0",
                "attributes": [
                  {
                    "name": "state",
                    "values": [
                      {
                        "ts": "2020-01-30 09:23:49 +0100",
                        "val": "DOWN",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 116,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:127_0",
                "attributes": [
                  {
                    "name": "state",
                    "values": [
                      {
                        "ts": "2020-01-29 17:39:54 +0100",
                        "val": "DOWN",
                        "val_t": "string"
                      }
                    ]
                  },
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NOROUTE",
                          "src": "nodeId=127_0;service=battery"
                        },
                        "ts": "2020-02-11 15:37:59 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 119,
            "services": [
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:out_lvl_switch/ad:130_0",
                "attributes": [
                  {
                    "name": "lvl",
                    "values": [
                      {
                        "ts": "2020-02-10 11:17:29 +0100",
                        "val": 99,
                        "val_t": "int"
                      }
                    ]
                  }
                ],
                "name": "out_lvl_switch"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:out_bin_switch/ad:130_0",
                "attributes": [
                  {
                    "name": "binary",
                    "values": [
                      {
                        "ts": "2020-02-11 15:35:47 +0100",
                        "val": true,
                        "val_t": "bool"
                      }
                    ]
                  }
                ],
                "name": "out_bin_switch"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:meter_elec/ad:130_0",
                "attributes": [
                  {
                    "name": "meter",
                    "values": [
                      {
                        "props": {
                          "unit": "kWh"
                        },
                        "ts": "2020-02-11 15:35:39 +0100",
                        "val": 1.30400002002716,
                        "val_t": "float"
                      }
                    ]
                  }
                ],
                "name": "meter_elec"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:color_ctrl/ad:130_0",
                "attributes": [
                  {
                    "name": "color",
                    "values": [
                      {
                        "ts": "2020-02-10 12:12:48 +0100",
                        "val": {
                          "blue": 255,
                          "green": 255,
                          "red": 255
                        },
                        "val_t": "int_map"
                      }
                    ]
                  }
                ],
                "name": "color_ctrl"
              },
              {
                "addr": "/rt:dev/rn:zw/ad:1/sv:dev_sys/ad:130_0",
                "attributes": [
                  {
                    "name": "state",
                    "values": [
                      {
                        "ts": "2020-02-12 08:04:44 +0100",
                        "val": "DOWN",
                        "val_t": "string"
                      }
                    ]
                  },
                  {
                    "name": "config",
                    "values": [
                      {
                        "ts": "2020-02-10 13:23:06 +0100",
                        "val": {
                          "85": "5;1"
                        },
                        "val_t": "str_map"
                      }
                    ]
                  },
                  {
                    "name": "error",
                    "values": [
                      {
                        "props": {
                          "msg": "TRANSMIT_COMPLETE_NO_ACK",
                          "src": "nodeId=130_0;service=out_bin_switch"
                        },
                        "ts": "2020-02-11 15:35:39 +0100",
                        "val": "TX_ERROR",
                        "val_t": "string"
                      }
                    ]
                  }
                ],
                "name": "dev_sys"
              }
            ]
          },
          {
            "id": 120,
            "services": [
              {
                "addr": "/rt:dev/rn:flow/ad:1/sv:sensor_presence/ad:ZpIa3heg12uaBN1_0",
                "attributes": [
                  {
                    "name": "presence",
                    "values": [
                      {
                        "ts": "2020-02-12 17:46:44 +0100",
                        "val": false,
                        "val_t": "bool"
                      }
                    ]
                  }
                ],
                "name": "sensor_presence"
              },
              {
                "addr": "/rt:dev/rn:flow/ad:1/sv:out_bin_switch/ad:ZpIa3heg12uaBN1_0",
                "attributes": [
                  {
                    "name": "binary",
                    "values": [
                      {
                        "ts": "2020-02-12 17:46:44 +0100",
                        "val": false,
                        "val_t": "bool"
                      }
                    ]
                  }
                ],
                "name": "out_bin_switch"
              }
            ]
          },
          {
            "id": 121,
            "services": []
          }
        ]
      }
    },
    "requestId": "158205049369172",
    "success": true
  },
  "val_t": "object",
  "ver": "1"
}

```