This is a library to interact with Mikrotik managed switches running SwOS

There is no official documentation for this API so it may not work for other versions or hardware.

It was tested with 	CSS106-1G-4P-1S running SwOs 2.18

Current status of implementation:

| Page in UI  | Page in API          | Status      |
|-------------|----------------------|-------------|
| Link        | LinkPage (link.b)    | Implemented |
| System      | SysPage (sys.b)      | Implemented |
| SFP         | SfpPage (sfp.b)      | Implemented |
| RSTP        | SysPage (sfp.b)      | Implemented |
| RSTP        | RstpPage (rstp.b)    | Implemented |
| Forwarding  | FwdPage (fwd.b)      | Implemented |
| VLAN        | FwdPage (fwd.b)      | Implemented |
| VLANs       | VlanPage (vlan.b)    | Implemented |
| Hosts       | HostPage (host.b)    | Missing     |
| Hosts       | HostPage (!dhost.b)  | Missing     |
| IGMP Groups | IgmpPage (!igmp.b)   | Missing     |
| SNMP        | SnpPage (!snmp.b)    | Missing     |
| ACL         | AclPage (!acl.b)     | Missing     |
| Statistics  | StatsPage (!stats.b) | Missing     |
| Errors      | StatsPage (!stats.b) | Missing     |

Status description:
- *Missing* - No implementation
- *Implemented* - Implemented but not tested
- *Complete* - Tested and documented
