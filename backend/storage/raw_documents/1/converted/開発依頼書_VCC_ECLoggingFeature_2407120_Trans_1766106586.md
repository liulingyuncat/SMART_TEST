# 開発依頼書_VCC_ECLoggingFeature_2407120.pptx

## 文档信息
- **文件名**: 開発依頼書_VCC_ECLoggingFeature_2407120.pptx
- **文件大小**: 529317 字节
- **MIME类型**: application/vnd.openxmlformats-officedocument.presentationml.presentation
- **文档类型**: PowerPoint文档(PPTX)
- **转换时间**: 2025-12-19 09:09:46

## 文档内容

### 第 1 页

VAIO ユーティリティ/アプリケーション開発依頼書

1

依頼

確認

承認

PC-Soft

古谷

2024/3/22

PC-Soft

小林

2024/3/25

PC-Soft

中山

2024/3/25

開発名

EC Logging Feature

対象機種

Spark

開発概要：

ECが保持している情報（バッテリーの充放電、LID Open/Close等）をテキストファイルに保存する機能。

スケジュール：

   Beta    Spark Pilot4 Build (2024/06/06)

RC   Spark RC1 Build (2024/07/25)

開発総工数：

開発工数見込み　　1人月

評価工数見込み　   1人月

他機種展開予定

Spark以降のSNCX EC Logging Feature (21Ah)

をサポートする機種全て

補足：

### 第 2 页

EC Logging Feature

SNCX EC Logging Feature (21Ah)

BIOS要求仕様書

\\vsv0054\fsroot\folders\SNTF0218\ALL\WG\BIOS\Requirements\StdBIOSSpec\Feature\SNCX EC Logging Feature

Standard-BIOS-SNCXECLoggingFeature-01.00.00.pdf

SNCX Command等に変更はありません。

どのBIOSから対応されるかは別途豊田さんと調整をお願いします。

SNCX System Charging Control Feature 3 (215h)

BIOS要求仕様書

\\vsv0054\fsroot\folders\SNTF0218\ALL\WG\BIOS\Requirements\StdBIOSSpec\Feature\SNCX System Charging Control Feature

Standard-BIOS-SNCXSystemChargingControlFeature-03.00.00.pdf

218hは廃止して215hの仕様に戻してください。

### 第 3 页

EC Logging Feature

SNCX EC Logging Feature (21Ah)

EC Log構成

ヘッダー4Byte + 6Byte単位のECイベントログで構成されており、ECが確保したRing Buffer Size分（Sparkは252Byte）のログが取得可能です。

ECイベントログはEvent Code(1Byte)、Parameter1(1Byte)、Parameter2(4Byte)の計6Byteで構成されています。

ECイベントログはRing Buffer方式で記録され、1番古いログは最新のログで上書きされます。

3

Current Write Index

(Next Log Index) is 5

LogDataのイメージ図

初期値は0 Fill

EC Event Log Format

### 第 4 页

EC Logging Feature

SNCX EC Logging Feature (21Ah)

Command 10h Get Current Time

各ECイベントログの発生日時を計算するための [Current Time Information] を取得します。

このTime InformationはBatteryが稼働している間、秒単位で増加する時間情報（Battery FW Runtime）です。

各ECイベントログに記録された [Log Time Information] と [Current Time Information] の差分を取ることで、そのイベントログが現在の時刻から何秒前に発生したログなのかを計算することができます。

[イベント発生日時] ＝ [現在のWindows日時] ー ([Current Time Information] ー [Log Time Information] )

                                             ※ 現在のWindowsの日時が狂っている状態でも、そのまま現在の時刻を使用して計算してください。

4

### 第 5 页

EC Logging Feature

SNCX EC Logging Feature (21Ah)

Command 11h Get Log Data

SNCX Event

S0中はLogが更新されるたびにSNCX Eventが発生します。

SNCX Eventが発生しない状況（パワマネ後、Utility起動時等）ではUtility側で自発的にCommand 11hを発行してください。

5

### 第 6 页

EC Logging Feature

Utilityへの実装要求

ダブルチェック処理

ECがLog書き込み中ではないことを判定するために100msec間隔で2回Command 11hを実行して同じ結果が得られることを確認してください。

リトライ処理

ダブルチェック処理がNGの場合、3秒待ってから再度ダブルチェックを最大10回まで繰り返してください。

リトライ処理中に割り込み(SNCX Event等)が発生した場合は、リトライ処理を中断して次の処理に移ってください。

Log 記録順計算処理

Ring Bufferにある各ECイベントログの記録順をCurrent Write Indexの値から逆算してください。

重複ログ除外処理

ダブルチェック処理がOKの場合、テキストファイルに保存済みのログを除外してください。

前回のCommand 11hの結果をメモリ配列上に保存しておいて比較するだけでは不十分なので注意してください。

シャットダウン状態でもECイベントログが生成されるため、Windows（Utility）起動直後の1回目の処理はテキストファイルに保存済みの結果と比較する必要があります。

6

Latest -2 Log

(Log Index is 41)

Latest -1 Log

(Log Index is 0)

Latest Log

(Log Index is 1)

例：Current Write Indexが2

RingBufferSizeが252の場合

### 第 7 页

EC Logging Feature

Utilityへの実装要求

[イベント発生日時]計算処理

Command 10hで[Current Time Information]を取得し、各ログの[イベント発生日時]を計算してください。

[イベント発生日時] ＝ [現在のWindows日時] ー ( [Current Time Information] ー [Log Time Information] )

下記のように正しく計算できない場合はテキストファイルに出力する[イベント発生日時]の値は無しにしてください

[Current Time Information]あるいは[LogTime Information]が0xFFFFFFFE(Unknown)の場合

バッテリーが取り外されている状態ではTime Informationが0xFFFFFFFEになります。

[Current Time Information] < [LogTime Information]の場合

基本的には遭遇しないはずですが、 ECの電源が入ったまま新しいバッテリーに交換されると遭遇します。

ログのParameter2がTime Informationとして定義されていない場合

現時点ではEvent Code 0xD8(Type-C Device identifier)のLogのみが該当します。

7

### 第 8 页

EC Logging Feature

Utilityへの実装要求

テキストファイル出力処理

保存場所

VCCの設定変更ログ機能でもLogフォルダを作成するようなので、そこと同じフォルダにします。

テキストファイルはBIOSバージョン毎に月単位で分けて保存してください。

%ProgramData%\VAIO Corporation\VAIO Control Center\Log\BIOSLog_[BIOSVersion]_[YYYYMM].txt

出力フォーマット

重複していない新規ログは下記のフォーマットでテキストファイルに追記出力してください。

下記の各列はカンマで区切って出力してください。

例

2024/03/13 08:19:08, 2024/03/13 08:19:10, DD 01 70 28 DB 00

, 2024/03/14 12:11:16, DD 00 FE EE EEEE, (Time InformationがUnknownの場合は[イベント発生日時]は無し)

8

[イベント発生日時] 

テキストファイル保存日時

EC Log 

Raw Data

Format

YYYY/MM/DD hh:mm:ss

YYYY/MM/DD hh:mm:ss

Binary

(hex)

(例: BIOSLog_R0181ST_202403.txt)

### 第 9 页

9

Flow Chart

Start VES

(Boot Windows) 

Resume from ModS/Hib

SNCX Event

Calculate Log Order

Wait 3 Sec

Double Check Process

(Call command 11h twice with 100msec interval 

and check if results are same)

Check retry count

NG

OK

> 10 times

< =10 times

Check Next Event

(SNCX event or 

resume from ModS/Hib)

Stop retry process 

and wait next event

No Event

Event occurred 

Exclude Duplicate Log

(Check if log has already been outputted text file)

Calculate Log Datetime 

Output Log to text file

Transition to next event

---
*本文档由自动转换工具生成*
