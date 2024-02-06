## Table of Contents

- [Ntp.proto](#Ntp.proto)
    - [Ntp](#siemens.iedge.dmapi.ntp.v1.Ntp)
    - [PeerDetails](#siemens.iedge.dmapi.ntp.v1.PeerDetails)
    - [Status](#siemens.iedge.dmapi.ntp.v1.Status)
  
    - [NtpService](#siemens.iedge.dmapi.ntp.v1.NtpService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="Ntp.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## Ntp.proto



<a name="siemens.iedge.dmapi.ntp.v1.Ntp"></a>

### Ntp
Type contains an array of ntp server addresses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ntpServer | [string](#string) | repeated | array of multiple ntp server address. |






<a name="siemens.iedge.dmapi.ntp.v1.PeerDetails"></a>

### PeerDetails
Peer Details from ntpq -p output


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| remoteServer | [string](#string) |  | NTP server address |
| referenceID | [string](#string) |  | Reference id for the NTP server |
| stratum | [string](#string) |  | Stratum for the NTP Server |
| type | [string](#string) |  | Type of server (local, unicast, multicast, or broadcast) |
| poll | [int32](#int32) |  | How frequently to query server (in seconds) |
| when | [int32](#int32) |  | How many seconds passed after the last poll. |
| reach | [string](#string) |  | octal bitmask of success or failure of last 8 queries (left-shifted). eg:375 |
| delay | [float](#float) |  | network round trip time (in milliseconds) |
| offset | [float](#float) |  | difference between local clock and remote clock (in milliseconds) |
| jitter | [float](#float) |  | Difference of successive time values from server (in milliseconds) |






<a name="siemens.iedge.dmapi.ntp.v1.Status"></a>

### Status
Type for ntp current sync status


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| isNtpServiceRunning | [bool](#bool) |  | indicates that ntp service is running or not |
| isSynced | [bool](#bool) |  | indicates NTP server synced or not |
| lastConfigurationTime | [string](#string) |  | time of the last performed iedk ntp configuration. |
| lastSyncTime | [string](#string) |  | time of the last ntp sync operation. |
| peerDetails | [PeerDetails](#siemens.iedge.dmapi.ntp.v1.PeerDetails) | repeated | NTPQ peer information array. Only exist after ntp configuration done. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="siemens.iedge.dmapi.ntp.v1.NtpService"></a>

### NtpService
Ntp service ,uses a UNIX Domain Socket "/var/run/devicemodel/ntp.sock" for GRPC communication.
protoc  generates both client and server instance for this Service.
GRPC Status codes : https://developers.google.com/maps-booking/reference/grpc-api/status_codes .

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SetNtpServer | [Ntp](#siemens.iedge.dmapi.ntp.v1.Ntp) | [.google.protobuf.Empty](#google.protobuf.Empty) | Set ntp server |
| GetNtpServer | [.google.protobuf.Empty](#google.protobuf.Empty) | [Ntp](#siemens.iedge.dmapi.ntp.v1.Ntp) | Returns ntp servers |
| GetStatus | [.google.protobuf.Empty](#google.protobuf.Empty) | [Status](#siemens.iedge.dmapi.ntp.v1.Status) | Returns NTP Status message. |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |
