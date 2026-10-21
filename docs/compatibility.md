# PHP互換性と検証状況

現行1169件について、公式PHP 8.3.35・8.4.26・8.5.11との版別`make compare`は各1169一致・0差異・0エラー・終了0。Go製PHPの固定Laravel・API Platformアプリは、3設定の[実TCP検証](full-stack-development.md)でもHTTP・SQLite CRUD・OpenAPI・2000回連続／8並列要求を通過し、各版終了0だった。修正前の不一致と制限環境の接続拒否は履歴として残す。これらの一致をPHP全機能の互換と扱わない。

以下の機能別記録にある古い件数や「比較待ち」は作業当時の履歴であり、現行1169件の比較状況は上記を基準とする。

`IteratorAggregate::getIterator()`を経由した`yield from`の24段有限連鎖をGoの3設定で確認した。[対応範囲・循環時の制限](full-stack-development.md)。初回1169件比較では別のGenerator委譲境界に差異が出て修正し、現行fixtureの範囲では後続の公式3版比較で一致した。循環時の診断など未収録の境界は未確認。

`array_walk_recursive()`は再帰中の配列を追跡し、70段の有限配列をGoの3設定で処理した。[対象ケースと制限](full-stack-development.md)。現行1168件の公式PHP3版比較は未実施で、期待値は暫定。

`extract(..., EXTR_REFS)`で入力配列の変数名と同名キーを取り込む場合の参照セルを修正した。[Goの対象出力とPHP公式ソースからの暫定判断](full-stack-development.md)。当時1167件の公式PHP3版実比較は未実施。

`array_walk()`／`array_walk_recursive()`の配列引数の格納先を呼び出し時に保持し、参照コールバックへの要素渡しに再利用する。[評価回数のfixtureと公式比較待ち](full-stack-development.md)。Go3設定の対象実行は成功したが、当時1167件の公式PHP3版比較は未実施。

`extract(..., EXTR_REFS)`の基本形をGoの3設定でローカル確認した。[fixture・未検証の境界](full-stack-development.md)。当時1167件の公式PHP3版比較は未実施で、期待値は暫定。

`array_walk()`／`array_walk_recursive()`は、参照コールバックへ渡す配列を逆順の名前付き引数`callback:`・`array:`から解決するよう修正した。[回帰fixtureと公式比較待ち](full-stack-development.md)。Goの3設定で対象出力・終了0を確認したが、当時1167件の公式PHP3版比較は未実施。

PDO束縛の追加範囲（2026-10-21）：`bindValue()`の主要`PDO::PARAM_*`型、NULを含むLOB、stream resourceの基本形、名前付き引数と`execute(array)`の文字列束縛をGoの3設定でローカル確認した。[対象fixtureと制限](full-stack-development.md)。公式PHP3版との実比較は未実施。

PDO状態の追加範囲（2026-10-21）：SQLite接続と準備済みStatementへ独立した`errorCode()`／`errorInfo()`、接続の`ATTR_ERRMODE`、列名・取得値に作用する`ATTR_CASE`・`ATTR_ORACLE_NULLS`・`ATTR_STRINGIFY_FETCHES`の基本経路を追加した。Goの3設定で対象fixtureを確認した。[対象と制限](full-stack-development.md)。当時1167件の公式PHP3版比較は未実施で、期待値は暫定。

SQLite失敗の追加範囲（2026-10-21）：`PDOException::getCode()`へSQLSTATE文字列を、`errorInfo`へSQLSTATE・番号・詳細を保持し、固定Laravelの一意制約例外へ接続した。[2 fixtureと公式ソースに基づく暫定期待値](full-stack-development.md)。Goの3設定では確認済みだが、公式PHP3版との当時1167件の実比較は未実施。

最新の追加範囲（2026-10-21）：固定vendorのLaravel/Eloquentカーソルの遅延実行・束縛引数・LazyCollection委譲と、`PDOException`の捕捉・`QueryException`への連鎖をGoの3設定で確認した。[3 fixtureと検証の限界](full-stack-development.md)を参照。当時1144件の公式PHP3版比較は未実施で、期待値は暫定。項目11と全6段階は部分実装のままである。

1167件時点の履歴（2026-10-21）：fixture・期待値は**1167件**（正常743・実行時290・構文134）。公式PHP 8.3.35・8.4.26・8.5.11との最後の比較済み基準は各997件一致・終了0の[版別生ログ](item11-progress.md)。追加170件と変更後の全件は公式比較未実施で、期待値は暫定。Goローカルでは通常オブジェクトの`foreach`とIteratorラッパー、固定vendorの`ArrayPaginator`・Finderの複数探索先・Laravel Collectionの`CachingIterator`、Generator、`artisan list`・`about`・`route:list`・独自`gophp:probe`、OpenAPI export、固定アプリのHTTP/SQLite/API Platform基本CRUDと追加8回の連続CRUD、[512回](validation/full-stack-http-512-after-callback.log)・[1024回の追加GET](validation/full-stack-http-1024-caching.log)、JSON-LD Collection形式の[ページ分割](validation/full-stack-pagination-jsonld.log)・8並列GETを実行した。enum caseの`E`形式シリアライズと`get_debug_type()`も追加した。[統合テストと制限](full-stack-development.md)。実行成功はPHP3版との互換性証明ではなく、項目11と全6段階は未完了。

HTTP入力の追加範囲：GETの`filter[name]`・`items[]`などとPOSTフォーム本文を既存の`parse_str()`解析へ接続した。PUT本文と別メディア型を`$_POST`へ取り込まない。[GoのHTTP対象テスト](validation/full-stack-http-query-array-target-local.log)は終了0。HTTP経路の公式PHP3版実比較とmultipart/form-dataなどの全境界は未実施・未対応であり、[詳細](full-stack-development.md)に分けて記録する。

固定アプリの同一Goプロセスでの[2000リクエスト試験](validation/full-stack-http-2000-stress-local.log)は終了0で、保持ヒープはGC後865448→744720バイトだった。これはローカル統合試験であり、PHP版別の厳密比較と運用条件の負荷検証とは区別する。

PDOStatementの未定義パラメータ実行でSQLiteが返す`HY000`を保持し、次の成功した`execute()`でエラー状態をクリアする。Go3設定の[対象fixtureログ](validation/full-stack-pdo-param-error-target-local.log)は終了0。診断とエラーモードの版別厳密比較は未実施で、[詳細](full-stack-development.md)に暫定範囲を記録する。

実HTTPサーバーの検証範囲：固定アプリのCLI migrationからAPI Platform CRUDまでを`gophp --serve`経由で確認する[版別スクリプト](../research/laravel-api/verify-gophp.py)を準備した。構文検査と[3設定の別プロセスでのmigration・独自コマンド・OpenAPI・テーブル作成](validation/full-stack-go-cli-openapi-subprocess-3versions-local.log)は終了0。現在のsandboxは[ローカル待ち受けを拒否](validation/full-stack-real-server-sandbox-check.log)するため、実HTTP応答は未検証。in-processのHTTPテスト成功と区別する。

以前の状態（履歴）：収録997件は[版別生ログ](item11-progress.md)で公式PHP3版それぞれと一致、終了0。現行fixture・期待値は**1035件**（正常613・実行時290・構文132）。今回の匿名クラス、配列分解代入、`goto`、`__call`／`__callStatic`／`__toString`、CLI shebangと引数、`file_exists()`、`(object)`、`$_ENV`などを含む追加38件と全件の公式比較はPodman接続拒否で未実施、期待値は暫定。[全6段階の実装状態](full-stack-development.md)。項目11とアプリ完成条件は未完了。

今回の追加前の比較済み基準：[PHP 8.3](validation/reference-property-974-83-mac.log)・[8.4](validation/reference-property-974-84-mac.log)・[8.5](validation/reference-property-974-85-mac.log)の版別生ログは各974件の個別MATCH、**974一致・0不一致・0エラー、終了0**を記録する。追加15件は収録fixtureの範囲で比較済み。版ごとの`make compare`で、`check.sh`のPASSはユーザー報告。以下の974件「未比較」は比較前の履歴。

974件比較前の基準（履歴）：[版別生ログ](item11-progress.md)で公式PHP3版が各959一致・0不一致・0エラー、終了0。追加15件と訂正した期待値は収録範囲で比較済み。以下の959件時点の「未比較」は当時の履歴。

959件比較前（履歴）：呼び出し元スコープの要素参照を共有判定へ加え、fixture・期待値は**959件**（正常539・実行時288・構文132）。Goローカル検証済み、公式PHP3版はPodman接続拒否で未比較。新規期待値は暫定。[項目11の詳細とログ](item11-progress.md)。項目11全体は部分実装。

今回の修正前（履歴、2026-10-18）：配列コピー時の参照分離を修正し、fixture・期待値は**958件**（正常538・実行時288・構文132）。Goローカル検証済み、公式PHP3版はPodman接続拒否で未比較。既存1件の訂正値と新規3件の期待値は暫定。[項目11の詳細とログ](item11-progress.md)。項目11全体は部分実装。

958件への修正前（履歴、2026-10-18）：配列要素への参照を加えたfixture・期待値は**955件**（正常535・実行時288・構文132）。Goのtest・vet・build・差分検査は終了0。公式PHP3版の対象・全件比較はPodman接続拒否で未実施、新規期待値は暫定。[項目11の対応範囲・試行ログ](item11-progress.md)。項目11全体は部分実装。

955件追加前の基準（履歴）：fixture・期待値は**944件**（正常524・実行時288・構文132）。[版別生ログ](item11-progress.md)でPHP 8.3.35・8.4.26・8.5.11の各版が**944一致・0不一致・0エラー、終了0**。変数間の参照代入を含む追加7件は収録fixtureの範囲で公式比較済み。今回は保存済みログの確認で、比較は再実行していない。項目11全体・項目6〜10の全仕様対応は未完了。[項目11の設計・制限](item11-progress.md)。以下の「比較待ち」は各時点の履歴。

944件比較前の状況（履歴、2026-10-18）：既存937件は[版別生ログ](items6-10-integration-gate.md)で公式PHP 8.3.35・8.4.26・8.5.11の各版と全件一致し、比較終了コードも0。変数間の参照代入を追加した現在のfixture・期待値は**944件**（正常524・実行時288・構文132）。追加7件と修正後全件の公式比較はPodman接続拒否で未実施、期待値は暫定。[項目11の設計と試行ログ](item11-progress.md)。項目6〜10の全仕様対応および項目11全体は未完了。以下は937件以前の履歴。

## `define()`第3引数の診断中断（2026-10-18、937件比較前の履歴）

`null`の`bool`変換で非推奨診断ハンドラが投げた場合は定数登録前に中断する。`true`の警告処理は別経路とする。fixture・期待値は**937件**（正常517・実行時288・構文132）。Goローカル検証は成功し、公式PHP3版の対象・全件比較はPodman接続拒否で未実施。追加2件の期待値は暫定。[監査・比較試行ログ](items6-10-integration-gate.md)を参照。以下は936件以前の履歴。

## 警告中断経路の追加監査（2026-10-18、936件時点の履歴）

`define(null, ...)`の非推奨診断でハンドラが投げても定数を登録していたGo側の経路を修正し、診断後の定数状態と後続副作用を観測するfixtureを1件追加した。総数は**936件**（正常516・実行時288・構文132）。既存935件の[版別生ログ](items6-10-integration-gate.md)は変更前の証拠。新規・全件の公式PHP3版再比較はPodman接続拒否で未実施、期待値は暫定。[監査範囲・ログ・移行条件](items6-10-integration-gate.md)。以下は935件時点以前の履歴。

## 複合代入エラー時の格納（2026-10-18、936件追加前の履歴）

fixture・期待値は**935件**（正常515・実行時288・構文132）。公式PHP 8.3.35・8.4.26・8.5.11の版別生ログで各版935一致・0不一致・0エラー。型付きプロパティの`/= 0`・`%= 0`・`+= []`・警告ハンドラ例外を伴う`+= '2x'`の失敗後の値0、および`.=`の左辺型別の違いを収録fixtureの範囲で確認した。`.=`では文字列型プロパティが値を保つ例と、数値型・nullable・union・mixedのプロパティが0になる例があり、一律の値保持規則とはしない。通常変数・配列要素の警告中断は型付きプロパティと分けて扱う。ビット複合代入への一般化は未検証。[版別ログと移行条件](items6-10-integration-gate.md)。比較コマンドの終了0と`check.sh`のPASSはユーザー報告で、ログには終了コードがない。以下の934件・933件の記述は当時の履歴。

934件時点の履歴：修正前の[集計](validation/items6-10-compound-934-before-fix-summary-mac.log)は3版とも933一致・1不一致・0エラー。Go側修正後はPodman接続拒否で比較できず、期待値を暫定としていた。[933件一致集計](validation/items6-10-compound-933-after-fix-summary-mac.log)は元ログとバイト一致で保存済み。

## 剰余複合代入の警告中断（2026-10-17、934件追加前の履歴）

現行fixtureは**933件**（正常513・実行時288・構文132）。`%=`の整数変換でハンドラが例外を投げた場合、格納前に中断する。Goの3版設定で検証済み。公式PHPの[8.3](validation/items6-10-modulo-compare-83.log)・[8.4](validation/items6-10-modulo-compare-84.log)・[8.5](validation/items6-10-modulo-compare-85.log)の比較試行はPodman接続拒否で各終了2。各版0一致・0不一致・1環境エラー・933件未実施。追加期待値は暫定。[移行判定](items6-10-integration-gate.md)。次節の931件は修正前の履歴。

## 項目6〜10の統合検証（2026-10-17、933件への追加前の履歴）

fixtureは**931件**（正常511・実行時288・構文132）。前回の[924件版別一致集計](validation/callable-924-after-fix-summary-mac.log)は今回の変更前の記録。新規7件はGoの3版設定と期待値テストで検証し、警告ハンドラ例外後の動的プロパティ・配列キー・増減の格納中断を修正した。公式PHPの[8.3](validation/items6-10-integration-compare-83.log)・[8.4](validation/items6-10-integration-compare-84.log)・[8.5](validation/items6-10-integration-compare-85.log)の版別比較試行はPodman接続拒否で各コマンド終了2。各版0一致・0不一致・1環境エラー・931件未実施で、期待値は暫定。[統合検証と移行判定](items6-10-integration-gate.md)。下の924件全一致は今回の変更前の履歴。

## callable／Closureの924件（2026-10-17、今回の変更前の比較済み履歴）

fixtureは**924件**（正常504・実行時288・構文132）。ユーザーが版ごとに`make compare`を実行した[修正後の6行集計](validation/callable-924-after-fix-summary-mac.log)ではPHP 8.3.35・8.4.26・8.5.11の各版が**924一致・0不一致・0エラー**。`call_user_func_array`と`call_user_func`自身の名前付き引数の修正を含む、新規16件の比較待ちは収録fixtureの範囲で解消した。`check.sh`のPASSはユーザー報告であり、集計ログに個別fixture行と各コマンドの終了コードはない。今回はログの保存・照合のみで比較は再実行していない。[以前の対象比較試行](validation/callable-named-targeted-84-attempt.log)と版別全件比較試行のPodman接続拒否は当時の履歴として残す。項目9と項目6〜10全体は部分実装。[実装・未対応の境界](items6-10-progress.md)を参照。

## 特殊型・未登録型の908件比較（2026-10-16、今回の変更前の履歴）

fixtureは**908件**（正常490・実行時287・構文131）。ユーザーが版別に`make compare`を実行した後の[修正後6行集計ログ](validation/special-types-908-after-fix-summary-mac.log)で、PHP 8.3.35・8.4.26・8.5.11がそれぞれ**908一致・0不一致・0エラー**と確認できた。888件からの追加20件は収録fixtureの範囲で公式比較済み。ログに個別fixture行と終了コードは含まれず、修正6件の解消と`check.sh`のPASSはユーザー報告である。今回は既存ログの保存・照合であり、比較は再実行していない。以下の不一致と接続拒否は当時の履歴として維持する。[未対応範囲](items6-10-progress.md)。

## 特殊型のレビュー差異を修正（2026-10-16、修正直後の履歴）

[ユーザー提供の修正前集計ログ](validation/special-types-908-before-fix-mac.log)は8.3.35・8.4.26で各902一致／4不一致／2テストエラー、8.5.11で903一致／3不一致／2テストエラー。`self`を宣言元クラス名として署名診断に表示し、親なし`parent`の引用符と未登録型名を修正した。継承比較に必要な後方型宣言を先行登録せず、2件を`success`から`runtime`へ移した。総数は**908件**（正常490・実行時287・構文131）。修正直後はPodman接続拒否で公式PHP3版比較を実施できず、期待診断を暫定としていた。後続の冒頭908件集計で収録範囲の比較待ちは解消した。[進捗](items6-10-progress.md)を参照。

## 特殊型と未登録型（2026-10-16、修正前の履歴）

この追加時点のfixtureは**908件**（正常492・実行時285・構文131）。既存888件は[版別の成功集計](validation/enum-888-after-fix-summary-mac.log)でPHP 8.3.35・8.4.26・8.5.11の各版と一致済み。新規20件はGo側で実行・期待値テスト済みだったが、公式PHPは[8.3](validation/special-types-compare-83-attempt.log)・[8.4](validation/special-types-compare-84-attempt.log)・[8.5](validation/special-types-compare-85-attempt.log)ともPodman接続拒否でバージョン確認に失敗し、各コマンド終了2だった。後続の修正前比較と今回の修正は上記を参照。`self`／`parent`は宣言元、`static`戻り値は呼び出し先のクラスで照合する。[対象範囲と制限](items6-10-progress.md)。

## enum宣言診断の修正（2026-10-16、収録fixtureの比較済み）

[修正前のMac比較集計](validation/enum-888-before-fix-mac-summary.log)では、各版884一致・4不一致・0エラーだった。Go側でenumの予約メソッドの宣言時点、プロパティ拒否文言、trait由来`from()`／`tryFrom()`の扱いを修正した。現在のfixtureは**888件**（正常477・実行時281・構文130）。ユーザーが版別に`make compare`を実行した後の[修正後集計ログ](validation/enum-888-after-fix-summary-mac.log)は、PHP 8.3.35・8.4.26・8.5.11で各**888一致・0不一致・0エラー**を記録する。追加11件の比較待ちは収録fixtureの範囲で解消した。この6行ログに個別fixture行と各コマンドの終了コードはなく、修正4件の解消と`check.sh`のPASSはユーザー報告である。今回は保存済みログの確認・保存のみ行い、比較を再実行していない。[Podman接続拒否による修正後試行](validation/enum-888-diagnostic-compare-all-attempt.log)は履歴として残す。[進捗](items6-10-progress.md)を参照。

## 通常enumのfrom／tryFrom宣言を許す修正（2026-10-16、修正前の履歴）

通常enumの`from()`／`tryFrom()`はユーザー定義メソッドとして扱う。backed enumでは既存の組み込みAPIとの競合として拒否し、`cases()`は両方で拒否する。直接宣言とtrait取り込み後の検査を修正した。戻り値の`self`／`?self`は今回の再現例に必要なメソッド範囲で追加した。fixtureは**888件**（正常476・実行時280・構文132）で、既存878件のみ[保存済み集計](validation/readonly-878-after-fix-summary-mac.log)で公式PHP3版と一致済み。新規11件の[対象比較](validation/enum-method-targeted-attempt.log)と[全件比較](validation/enum-method-compare-all-attempt.log)は環境エラーにより未実施で、期待値は暫定である。[進捗](items6-10-progress.md)を参照。

## enum基本機能の884件時点の検証状況（2026-10-16）

Go側へ通常／int・string backed enum、case同一性、`name`／`value`、`cases()`、`from()`／`tryFrom()`、型検査、interface・プロパティを持たないtraitの基本経路を追加した。fixtureは**884件**（正常475・実行時279・構文130）。既存878件は[版別集計](validation/readonly-878-after-fix-summary-mac.log)で3版とも一致済み。新規6件は[対象比較](validation/enum-targeted-attempt.log)と[全件比較試行](validation/enum-compare-all-attempt.log)がPodman接続拒否で終了2となり、各版0一致・0不一致・1環境エラー、6件未実施。期待値は暫定である。宣言制約、定数式、組み込みAPIの引数・診断、比較・変換の境界は残る。[進捗](items6-10-progress.md)を参照。

## readonlyとプロパティ昇格の878件時点の検証状況（2026-10-16）

fixtureは**878件**（正常471・実行時278・構文129）。コンストラクタ引数からのプロパティ宣言と本体実行前の初期代入、readonlyの一度だけの初期化・再代入／unset／配列書き換え拒否、readonly classの型付きプロパティと動的プロパティ作成禁止を部分実装した。readonly親クラスを通常の子クラスが継承する誤受理も修正した。PHP 8.3ではreadonlyの初期化を宣言元クラスに限り、8.4・8.5では公開／保護プロパティの子クラス初期化を許す。ユーザーが版別に`make compare`を実行した後の[6行の集計ログ](validation/readonly-878-after-fix-summary-mac.log)で各版**878一致・0不一致・0エラー**を確認し、追加27件の比較待ちは収録範囲で解消した。集計ログに個別fixture行・終了コードはなく、指摘2件の解消と`check.sh`のPASSはユーザー報告である。比較を今回再実行していない。[修正前の876一致・1不一致・1テストエラー](validation/readonly-878-before-fix-mac.log)と[接続拒否による修正後の試行](validation/readonly-878-followup-compare-all-attempt.log)は履歴として維持する。clone、参照、フック、非対称可視性、未収録の宣言診断は残る。[進捗と制限](items6-10-progress.md)を参照。

## 851件の公式比較済み基準（別実行による確認）

ユーザーのMacでの[全件比較](validation/readonly-baseline-851-compare-all-mac.log)は、8.3.35・8.4.26が各851一致／0不一致／0エラー、8.5.11が830一致／0不一致／21環境エラーで`make compare-all`終了2。VM再起動後の[8.5.11単独比較](validation/readonly-baseline-851-85only-mac.log)は851一致／0不一致／0エラー、終了0だった。3版それぞれで851件一致を確認したが、1回の`make compare-all`が成功したわけではない。Podman資源枯渇は未確認の仮説である。

## 動的プロパティ追加直後の記録（履歴）

動的インスタンスプロパティ名`$object->$name`・`$object->{$name}`・`$object->{式}`とnullsafeの読み取り、通常／複合代入、増減、`isset`・`empty`・`??`・`??=`・`unset`、配列要素操作を既存の宣言スロットへ接続した。fixtureは**851件**（正常457・実行時276・構文118）。ユーザーのMacでの[修正前の全件比較](validation/dynamic-property-851-before-fix-mac.log)は8.3.35・8.4.26が各850一致／1不一致／0エラー、8.5.11が829一致／1不一致／21環境エラー、終了2。不一致の`dynamic_property_receiver_change.php`は、Goが代入先の受け手を名前式・右辺より早く固定していたため。Go側を修正し、公式PHPのstdout`1|9|1|10`に期待値を訂正した。[修正後のローカルテスト](validation/dynamic-property-receiver-fix-test.log)は成功したが、[対象比較](validation/dynamic-property-receiver-fix-targeted-attempt.log)・[全件比較](validation/dynamic-property-receiver-fix-compare-all-attempt.log)はPodman接続拒否で終了2。**修正後の公式比較は未実施**である。旧[比較試行](validation/dynamic-property-compare-all-attempt.log)も履歴として維持する。名前式の非文字列変換、評価順序の未収録ケース、マジックメソッド等は残る。項目7・8と項目6〜10全体は未完了。[詳細](items6-10-progress.md)を参照。

### 833件の比較済み基準（履歴）

fixtureは**833件**（正常443・実行時273・構文117）。型付きインスタンス／staticプロパティ、型付きクラス／interface／trait定数の基本形、および複合代入式の格納値返却を追加した。ユーザーのMacでの[修正後の全件比較](validation/typed-members-833-retry-mac.log)は、公式PHP 8.3.35・8.4.26・8.5.11の各版で**833一致・0不一致・0エラー**、ログ記載の終了コード0。元ログとのバイト一致を確認して保存した。817件から追加した16件の公式比較待ちは収録fixtureの範囲で解消した。途中の[診断2件の不一致](validation/typed-members-833-before-fix-mac.log)と[8.5の一時的なタイムアウト](validation/typed-members-833-after-mac.log)は履歴として残す。今回はログの確認・保存で、比較を再実行していない。動的インスタンスプロパティアクセスなどは未対応で、項目8と項目6〜10全体は部分実装のままである。[進捗・制限](items6-10-progress.md)。

## 817件比較時点（履歴）

2026-10-15時点は**817件**（正常433・実行時272・構文112）。ユーザーの[修正後の公式PHP3版全件比較](validation/void-never-817-after-fix-mac.log)は各版**817一致・0不一致・0エラー**、終了コード0。787件からの追加30件は収録fixtureの範囲で公式比較済み。今回は保存済みログを確認・保存したもので、比較は再実行していない。[修正前の不一致](validation/void-never-817-before-fix-mac.log)は履歴として維持する。型付きプロパティ・定数、未登録型の解決などは残る。[詳細](items6-10-progress.md)。

### 787件追加直後の比較試行（履歴）

当時は新規18件がGoローカルテストを通過したが、[全件比較試行](validation/type-variance-float-compare-all-attempt.log)はPodman接続拒否で終了2。各版0一致・0不一致・1環境エラーであり、この時点の新規期待値・診断は暫定だった。後続の上記787件成功ログで収録範囲の比較待ちは解消した。

### 782件追加直後の記録（履歴）

当時は13件を追加し、Goローカルテストが通過した。[当時の全件比較試行](validation/type-variance-compare-all-attempt.log)はPodman接続拒否で終了2。`int`を`float`の宣言時下位型とした誤りは、後続レビューを受けて上記のとおり修正した。

### 769件の確定記録（履歴）

ユーザーのMacで実行した`/tmp/php-compare-769-after.log`を[保存済みログ](validation/dnf-769-after-fix-mac.log)へバイト単位でコピーした。各版769一致・0不一致・0エラー、ログ末尾の終了コード0。以下の765件・769件追加直後の接続拒否と暫定期待値は当時の履歴であり、収録769件の比較待ちはこの成功で解消した。今回はログの確認・保存で、比較を再実行していない。

### 765件時点の追加直後（履歴）

当時は追加15件がGoローカル照合を通過し、[比較試行](validation/dnf-compare-all-attempt.log)はPodman接続拒否で終了2、公式PHP3版との比較は未実施だった。この時点ではintersection同士の包含関係を宣言時に検査していなかった。

### 750件時点の確定履歴

当時の収録fixtureは**750件**（正常405・実行時257・構文88）。主基準PHP 8.4.26、追加対象8.3.35・8.5.11。ユーザーのMacでの[修正後全件比較](validation/union-diagnostic-750-after-fix-mac.log)は各版**750一致・0不一致・0エラー**、ログ記載の終了コード0。[`union_nullable_syntax.php`の対象比較](validation/union-nullable-syntax-targeted-mac.log)も各版1一致・0不一致・0エラーで、比較コマンド終了0はユーザー報告による確認である。追加23件の比較待ちは収録fixtureの範囲で解消した。これは保存済みログの確認・保存であり、今回の比較結果ではない。項目8と項目6〜10全体は部分実装で、Go製PHPでのLaravel起動は未対応。[詳細](items6-10-progress.md)を参照。

### 750件の修正前不一致（履歴）

[修正前の全件比較](validation/union-diagnostic-750-before-fix-mac.log)は各版749一致・1不一致・0エラー、終了2だった。唯一の差は`parser/union_nullable_syntax.php`の診断文で、Go側の解析エラーと期待値を公式PHPの実出力に合わせて修正した。この環境からの[対象比較試行](validation/union-nullable-syntax-targeted-attempt.log)と[全件比較試行](validation/union-nullable-syntax-compare-all-attempt.log)はPodman接続拒否で終了2だったが、後続の上記Mac成功ログで収録fixtureの差異解消を確認した。

### 727件の全一致（直前の確定履歴）

ユーザーのMacでの[修正後全件比較](validation/static-throw-727-compare-all-mac.log)は各版727一致・0不一致・0エラー、終了0。新規9件を含む全件がMATCHし、直接`Class::{throw ...}()`の事前拒否は収録fixtureの範囲で公式比較済みである。この記録は保存済みログの確認・保存であり、今回の比較結果ではない。

### 718件比較時点（履歴）

修正前に比較済みの収録fixtureは**718件**（正常393・実行時250・構文75）。主基準PHP 8.4.26、追加対象8.3.35・8.5.11。ユーザーのMacでの[修正後全件比較](validation/item7-dynamic-static-718-after-fix-mac.log)は各版718一致・0不一致・0エラー、終了0。`success/class_dynamic_static_order.php`も3版すべてで一致し、追加11件の比較待ちは解消した。今回は保存済みログの確認・保存であり、比較は再実行していない。[修正前の失敗](validation/item7-dynamic-static-718-before-fix-mac.log)は履歴として残す。直接`Class::{throw ...}()`をGoが受理する差は、この718件の比較時点では収録fixture外の既知不具合だった。[進捗・対象範囲](items6-10-progress.md)を参照。項目7と項目6〜10全体は部分実装で、Go製PHPでのLaravel起動は未対応。

### 727件追加直後の比較待ち（履歴）

2026-10-15追補：当時は**727件**（正常394・実行時251・構文82）。直接の動的staticメソッド名`{throw ...}`とその括弧付き形を事前検査で拒否する修正を加え、未選択分岐、未呼び出し関数・メソッド、動的クラス名、`parent::`、読み込み先、合法な三項式と関数呼び出しの対比を9件に収録した。Go側の期待値テストは成功。この時点では新規9件の公式PHP3版比較がPodman接続拒否で**未実施**で、期待診断は暫定だった。後続の冒頭727件成功ログで比較待ちは解消した。[対象比較試行](validation/static-throw-targeted-attempt.log)と[全件比較試行](validation/static-throw-compare-all-attempt.log)は一致・不一致の実測値ではない。上記718件の全一致は変更前の確定履歴である。[詳細と再実行手順](items6-10-progress.md)を参照。項目7と項目6〜10全体は引き続き部分実装。

### 641件の修正前比較（履歴）

ユーザーの[修正前全件比較](validation/include-integration-641-user-before-fix.log)は、8.3.35・8.4.26が各640一致・1不一致、8.5.11が638一致・3不一致、環境エラー各0、終了2だった。弱い`int`引数への小数変換の非推奨警告と8.5の読み込み先宣言エラートレースを修正し、後続の上記641件成功ログで解消を確認した。修正途中の[保存済み公式stderrとの照合](validation/include-integration-641-saved-output-recheck.log)と[比較試行](validation/include-integration-641-review-compare-all-attempt.log)は、当時の記録として維持する。

## 622件の修正前比較（履歴）

[修正前全件ログ](validation/named-exception-622-user-before-fix.log)は、各版620一致・2不一致・0環境エラー、終了2だった。名前付き引数の重複・未知名の2件を公式PHPのスタックトレースに合わせて修正し、[保存済み公式出力との照合](validation/named-exception-binding-saved-output-recheck.log)は2件×3設定で一致した。当時は修正後の全件比較をPodman接続制限で実施できなかった。上記の後続成功ログで比較待ちは解消した。

## 605件以前の比較・実装記録（履歴）

当時の605件の[修正前比較ログ](validation/argument-unpack-compare-all-mac.log)では各版602一致・2不一致・1検証エラー、終了2だった。その3件を修正した後、この環境からの[比較試行](validation/diamond-types-fix-compare-all-attempt.log)はPodman接続拒否で終了2だった。これらは後続の605件全一致ログで解消した過去の状況であり、当時の結果を書き換えない。

## 過去の比較・実装記録

以下の「現在」「今回」「未実施」は各段階の記録時点を指す。588件全一致は605件の比較より前に確定した結果である。

以下の561件の記述は追加前の確定結果である。

ユーザーの[再比較ログ](validation/class-private-diagnostic-compare-all-mac.log)で、PHP 8.3.35・8.4.26・8.5.11の各561 fixtureが全一致（差異・環境エラー0、`make compare-all`終了0）した。以前の14件の「公式比較待ち」は解消した。今回以後に追加するfixtureの比較結果は別に記録する。

以下は再比較前の履歴である。

当時の比較対象は**561 fixture**。ユーザー提供の[修正前の全件比較ログ](validation/class-private-561-user-comparison-before-fix.log)では、PHP 8.3.35・8.4.26が各560一致・1不一致・0環境エラー、8.5.11が556一致・5不一致・0環境エラーで、`make compare-all`は終了2だった。全版共通のfinalクラス継承エラー文言と、PHP 8.5の宣言時Fatal errorのスタックトレース・単数形を公式出力に合わせて修正した。修正後のGoローカル検証と修正前の公式PHP出力7組への再照合は成功した。修正後の公式PHP全件再比較はPodman接続拒否のため未実施である。詳細は[項目6〜10の進捗](items6-10-progress.md)を参照する。

以下は修正前の履歴である。

同名privateプロパティの独立格納と宣言失敗系を追加し、現在は**561 fixture**。今回追加の11件と直前の3件はGoローカル期待値テストのみ成功、公式PHP3版の比較はPodman接続拒否で未実施。既存547件の一致記録とは分けて扱う。詳細とMac再比較コマンドは[項目6〜10の進捗](items6-10-progress.md)にある。

547件比較後に型なしクラス継承のメンバー検査とfinal・abstractの部分集合を実装し、現在550 fixture。追加3件はGo側の期待値テストのみ成功し、公式PHP3版との比較はPodmanソケット接続拒否により未実施。[継続実装と検証ログ](items6-10-progress.md)を参照。以下の547件全一致は追加前の結果である。

ユーザー提供の[547件比較ログ](validation/items6-10-547-user-comparison.log)で、PHP 8.3.35・8.4.26・8.5.11の各547件が一致（差異0・環境エラー0、`make compare-all`終了0）した。以下の「新規20件は未実施」はこの比較以前の履歴であり、現在の判定ではない。

2026-10-13の2回のレビュー修正で計9件を追加し、現在は547 fixtureである。追加9件と前回追加の11件はGo側の期待値テストのみ通過し、公式PHP3版との実比較はPodman接続拒否のため未実施。既存527件の過去の一致記録と区別する。詳細は[項目6〜10の進捗](items6-10-progress.md)。以下の「538件」はレビュー修正前の記録として読む。

基本構文・名前空間追加範囲と当時の各507件一致は[現段階の互換性](phase2-compatibility.md)に記録した。その後のnullsafe修正で527件となり、ユーザー提供ログでは公式PHP 3版で全件一致している。Composerの共通lockとvendorは生成・監査済みで、公式PHPの最小アプリ検証も成功した。[初期目標](laravel-api-target.md)と[依存監査](dependency-audit.md)を参照。さらに今回11件を追加し、Go側の全538件の期待値テストが通過したが、新規分の公式PHP実比較は未実施である。変数・シンボル・エラーの構造は[実行時設計](phase2-design.md)を参照。以下は **第1版466件の完了時点の記録** である。

## 第1版の状態（2026-10-12）

**限定した対応範囲を持つCLI処理系の第1版が完成。** 項目9〜12を追加し、既存の項目1〜8を維持しています。
ユーザー環境での公式PHP実比較ログにより、8.3.35・8.4.26・8.5.11の各466件、
合計1,398件の一致（差分0・検証エラー0）を確認しました。
第1版の完成条件である「追加fixtureを含む全件が3バージョンの公式PHP実比較で成功」は達成済みです。
正常系208件は各版とも両処理系で終了コード0です。
[完成確認の記録](validation/v1-user-comparison-complete.md) に、ユーザー提供の実比較と
エージェント環境でのローカル検証・接続を拒否された試行を区別して記載しています。

| 項目 | 状況 |
| --- | --- |
| 主基準 | PHP 8.4.26、64ビット整数のCLI環境 |
| 追加比較対象 | PHP 8.3.35 / 8.5.11 |
| Go対象設定 | `--php-version=8.3\|8.4\|8.5`、既定8.4 |
| 実装・fixture | 全3設定で共通。版差は評価・診断の対象設定で切り替える |
| 開発環境 | Go 1.27.0、Python 3.9.6、darwin/arm64。go.modはGo 1.22以上 |
| 公式PHPの役割 | 比較基準。Goインタプリタの実行依存ではない |

旧文書の過去の検証結果・当時の仕様説明は [履歴](validation/pre-v1-compatibility.md) に移しました。
以下は第1版時点の対応範囲です。今回追加した構文は上記の現段階の互換性文書を参照してください。

## 構文と型

| 分類 | 対応する部分集合 |
| --- | --- |
| ファイル | 先頭の `<?php` と空白、末尾の省略可能な `?>`、LF / CRLF / CR |
| コメント | `//`、`#`、`/* ... */` |
| 値 | int64、float64、バイト列string、bool、null、順序付き配列、オブジェクト |
| リテラル | 10進整数、10進実数・指数表記、単一／二重引用符、true / false / null |
| 代入先 | 変数、配列要素、インスタンスプロパティ |
| 代入 | `=`、`+= -= *= /= %= .=`、前置／後置 `++ --`、式内・右結合の連鎖代入 |
| 二項演算 | `+ - * / % . < <= > >= == != <> <=> === !== && \|\|` |
| 単項演算 | `+ - !`、上記増減、`(bool) (int) (float) (string)` |
| 文 | echo（複数式可）、式文、空文、ブロック、if / elseif / else if / else、return |
| ループ | while、for、do / while、配列の値渡しforeach、レベル付きbreak / continue |
| 配列 | `[]` / `array()`、キー付き要素、末尾カンマ、添字読み書き、末尾追加、ネスト、結合・比較 |
| 関数 | 名前による通常関数の定義・呼び出し、位置引数・既定値・値渡し・再帰 |
| クラス | 継承なしのclass、stdClass、プロパティ・メソッド・__construct・$this・アクセス制御 |
| 言語構造 | isset、empty、unset |
| 定数 | COUNT_NORMAL = 0、COUNT_RECURSIVE = 1（名前は大小文字を区別） |

ファイル全体を解析・検査してから実行します。未選択分岐、短絡される右辺、呼ばれない関数、
終了済みループの後半にある構文エラーや無効なbreakレベルも拒否します。
選択されなかった式・本体の副作用、警告、呼び出し、ゼロ除算は実行しません。
末尾の終了タグの直後は改行1つまでを許し、インラインHTMLや複数PHPブロックは拒否します。

文字列はUTF-8を含むバイト列です。単一引用符の `\\`・`\'`、二重引用符の
`\n \r \t \v \e \f \\ \" \$`、1〜2桁の16進エスケープ、255以下の8進エスケープに対応します。
変数展開・heredoc / nowdoc・Unicodeコードポイントエスケープは対象外です。
文字列添字の読み取りとisset / emptyに対応し、書き込みは拒否します。文字列添字のunsetはPHP形式のErrorです。

### 演算・比較・代入

高い順に添字・プロパティ・呼び出し、増減・単項符号・キャスト、`!`、`* / %`、`+ -`、`.`、
`< <= > >=`、`== != <> === !== <=>`、`&&`、`||`、代入です。
算術・連結・論理演算は左結合、代入は右結合です。大小比較群と等価比較群はそれぞれ非結合です。
`false || $a = 7` は `false || ($a = 7)` として扱います。
代入式は代入後の値、前置増減は更新後、後置増減は更新前の型と値を返します。
代入先の添字・受け手の関数／メソッド呼び出しを重複評価しません。

整数の加減乗算は範囲を超えるとfloatになります。整数同士の割り切れる除算はint、それ以外はfloatです。
最小整数 / -1もfloatです。剰余は整数変換してintを返し、符号は被除数に従います。
最小整数 % -1は0です。変換後のゼロで割る場合もDivisionByZeroErrorになります。
算術で数値文字列を認識し、先頭だけ数値の場合は警告、非数値文字列はTypeErrorです。
bool・nullの算術変換、floatの精度喪失診断、整数境界、負のゼロ、INF・NaNを扱います。
配列同士の `+` は左側のキーを優先する結合です。オブジェクトの `%` / `%=` はTypeErrorです。

緩い比較では数値文字列、数値と非数値文字列、nullと文字列、bool変換、配列・オブジェクトを区別します。
厳密比較は型を変換しません。NaNは自身とも等しくありません。
配列の厳密比較はキー・値・型・順序、緩い比較は順序を無視したキーと値を比較します。
配列・オブジェクトの大小比較には全順序にならない組み合わせがあります。
オブジェクトの厳密比較は同一インスタンス、緩い比較はクラスとプロパティによります。
文字列化はPHPのprecision=14を目標とし、出力・連結でのtrueは `1`、false / nullは空文字列です。
非数値文字列の増減やboolの増減など、既存の警告・非推奨診断も維持しています。

### 関数・メソッドとスコープ

関数・メソッド・クラス名は大小文字を区別せず、変数・プロパティ名は区別します。
位置引数、配列を含む対応定数式の既定値、値渡し、余分な位置引数の評価、再帰を扱います。
引数不足はArgumentCountError、未定義関数はError、重複定義はFatal errorです。
関数のトップレベル直下の宣言は記述位置より前から使用できます。
分岐・ループ・関数内の宣言はそこへ到達した時点で登録され、再実行すると重複定義になります。
各呼び出しの変数スコープは独立し、通常関数は呼び出し元メソッドの$thisやprivateアクセス権を引き継ぎません。
returnはネストしたループも終了させ、値省略・末尾到達の戻り値はnullです。

クラスはpublic / protected / privateまたはvarによる型なしプロパティを持ちます。
同じクラスのメソッドから非publicメンバーを使用できます。継承はありません。
__constructを呼び出し、その戻り値はnewでは無視します。明示的なメソッド呼び出しでは戻り値を受け取れます。
独自クラスの動的プロパティ作成はDeprecated、stdClassでは診断なしです。
未定義プロパティ、可視性違反、メソッド不足、引数不足、呼び出し履歴の既存診断を維持しています。

配列は代入・引数・戻り値で値コピーになります。実装では公開した配列を変更せず、更新経路をコピーします。
配列中のオブジェクトを含め、オブジェクトは同じインスタンスを共有します。
キーの型変換、挿入順序、同じキーの更新位置、負のキーからの自動採番、整数最大キーも扱います。

## foreach（項目9）

`foreach ($expr as $value)` と `foreach ($expr as $key => $value)` を実装しました。
式を開始時に1回だけ評価し、その配列の挿入順で反復します。値はコピー、キーはint / stringのままです。
反復中の元配列の更新・追加・unset・再代入は反復対象を変更しません。
ループ変数内の配列更新も元配列に影響せず、オブジェクトのプロパティ変更は共有されます。
値の代入がキーの代入より先です。添字・プロパティも反復時の代入先にできます。

空配列ではループ変数を代入しません。終了時には最後に代入・変更された変数の値が残ります。
本体でunsetした変数はそのまま未定義です。関数内ではその関数のローカル変数になります。
while / for / do / foreachを混在してネストでき、break / continueの省略値は1です。
正の整数リテラルだけをレベルに指定できます。ループ外・不正値・ネストを超える指定は実行前に拒否します。
continue先がforなら更新式、do / whileなら条件式を評価します。break後の不要な式は評価しません。

非配列のスカラー・nullでは `foreach() argument must be of type array|object, ... given` を警告し、
本体を実行せず後続文へ進みます。未定義変数を渡す場合は通常の未定義変数警告も発生します。
参照形式と分解代入は構文解析時、オブジェクト反復は実行時の明確なgophp診断で拒否します。
これら対象外形式はPHP自身では有効なため、Goの拒否テストで検査し、互換性一致の主張に含めません。

## isset・empty・unset・キャスト（項目10）

isset・emptyは専用の式、unsetは専用の文です。通常関数への変換では実装していません。

| 構文 | 挙動 |
| --- | --- |
| `isset($a, $b[0], $o->p)` | 存在し、かつnullでない対象だけtrue。左から評価し、最初のfalseで短絡 |
| `empty($x)` | 未定義またはPHPの真偽値でfalseとなる値ならtrue |
| `empty(通常の式)` | 式を1回評価して真偽値を反転。式中の通常の警告は保持 |
| `unset($a, $b[0], $o->p)` | 左から対象を削除。対象式の呼び出し・添字を重複評価しない |

isset / emptyは未定義変数、存在しないキー・プロパティ、アクセスできないprivate / protectedプロパティを
通常の読み取り警告なしで検査します。nullはissetでfalse、emptyでtrueです。
添字の式自体は通常評価されるため、未定義の添字変数、関数呼び出しの失敗、無効なキー型などの診断は残ります。
スカラーへの存在確認や文字列の範囲外添字も扱います。issetは通常式そのものを対象にできません。
$this自体のisset / emptyと、オブジェクトコンテキスト外での$thisのプロパティアクセスを区別します。

unsetは未定義の変数そのものや欠落キー・プロパティを削除しても作成しません。
一方 `unset($missing[0])` は未定義の配列変数の警告を出します。
nullの添字削除は何もせず、falseは非推奨診断、その他の非配列や文字列の添字削除はErrorです。
private / protectedの削除はアクセス違反になります。宣言済みプロパティを削除して再代入しても動的プロパティにはしません。
配列の残りの順序を維持し、削除したキーの再挿入は末尾になります。末尾追加の採番は削除で巻き戻りません。
最大整数キーを削除した場合は、その最大キーへの追加を再び許可します。
値コピー済み配列の削除は他の変数へ影響せず、オブジェクトのプロパティ削除は共有されます。

キャストは `(bool)`・`(int)`・`(float)`・`(string)` の4つに限定しています。
`(boolean)`・`(integer)`・`(double)`・`(real)`・`(binary)` の別名と、
`(array)`・`(object)`・`(unset)` は診断で拒否します。大文字小文字や括弧内の空白は許可します。

明示的な数値キャストでは、先頭の数値部分を使い、非数値文字列を0へ変換します。
暗黙の算術変換に伴う「非数値文字」警告や小数部切り捨ての非推奨を、明示的キャストへ流用しません。
有限の範囲外floatからintは64ビットで折り返し、数値文字列からintは有限範囲外なら端値へ飽和します。
INF / NaNからintは0です。空配列の数値変換は0、非空は1です。
オブジェクトのint / float変換は警告と1、string変換はErrorになります。
配列のstring変換は警告と `Array`。bool変換では配列の空／非空、文字列 `''` / `'0'` などPHPの規則を使います。
PHP 8.5の範囲外floatやNaNの警告は明示的キャストでも対象設定に従って出します。

## 限定した標準関数（項目11）

| 関数 | 引数と戻り値・診断 |
| --- | --- |
| `strlen($string)` | 引数1つ。バイト数をintで返す。int / float / boolは文字列へ変換、nullはDeprecatedと0。array / objectはTypeError |
| `count($value, $mode = COUNT_NORMAL)` | 引数1〜2。配列を数える。通常は直下の要素数、COUNT_RECURSIVEでは配列の子要素数も加える。オブジェクトの中身は辿らない |
| `gettype($value)` | 引数1つ。`NULL / boolean / integer / double / string / array / object` の文字列 |
| `is_int / is_float / is_string / is_bool / is_null / is_array / is_object` | 各引数1つ。型を変換せずboolを返す |
| `is_numeric($value)` | 引数1つ。int / floatと、全体が数値形式の文字列でtrue。前後のASCII空白を許可し、末尾の非数値文字は不許可 |

is_numericはfloatのINF / NaNでもtrue、文字列 `'INF'` / `'NAN'`、空文字列、空白のみ、
bool / null / array / objectはfalseです。strlenのUTF-8文字列 `日本語` は9バイトです。
標準関数名は大小文字を区別しません。上表の組み込み関数と同名のユーザー関数は重複定義の診断です。
メソッド名には同じ名前を使えます。位置引数は左から評価し、余分な引数にも副作用があります。
指定数を外れる場合はArgumentCountErrorです。未定義関数の既存動作は維持しています。

countはこの部分集合では配列だけを受け付けます。PHPのCountableオブジェクトはinterface機能とともに対象外です。
非配列は `Countable|array` を要求するTypeErrorになります。
モードはint、bool、数値文字列、範囲内floatを受け付け、小数部の切り捨てではDeprecatedを出します。
nullはDeprecatedと0、非数値文字列・配列・オブジェクト・INF・NaNはTypeErrorです。
変換後の0 / 1以外はValueErrorです。モードの変換と値チェックは、第1引数の型チェックより先です。
参照形式がないため循環配列は作成できず、PHPの循環配列に対するcountの再帰警告は対象外です。

診断の履歴はPHPの専用命令に合わせ、1引数のstrlen / countの型エラーでは組み込み呼び出しフレームを追加しません。
引数個数違反と2引数countのエラーではフレームを追加します。呼び出し元の通常関数・メソッドの履歴を維持します。

## CLI仕様と実行例（項目12）

```sh
make build
./bin/gophp examples/report.php
./bin/gophp --php-version=8.3 examples/report.php
./bin/gophp --php-version 8.4 -- examples/report.php
./bin/gophp --help
```

[report.php](../examples/report.php) は [共通fixture](../tests/fixtures/success/v1_report.php) と同じ内容です。
配列引数をforeachで処理して配列を返し、isset / empty / unset、キャスト、標準関数、
ローカルスコープ、配列コピー、オブジェクト共有、短絡評価を組み合わせています。
[再帰の例](../tests/fixtures/success/v1_recursive_array_processing.php) はネスト配列を再帰的に集計します。

CLIの書式は `gophp [--php-version=8.3|8.4|8.5] <file.php>` です。
対象を省略すると8.4。オプションはファイル名より前に置き、`--` でオプション解釈を終了できます。
相対パスは絶対パスへ変換し、診断にはそのファイルパスと行番号を出します。
空白・日本語を含むパスでも処理します。ファイルは1つだけ受け付けます。

| 終了コード | 意味 |
| --- | --- |
| 0 | 正常実行、非致命警告を伴う実行、トップレベルreturn、ヘルプ表示 |
| 1 | ファイル読み込み・出力等のI/O失敗 |
| 2 | 引数不足・過多、未知のオプション、無効な対象バージョン |
| 255 | Parse error、コンパイル時Fatal error、実行時の致命エラー、未対応構文・操作の拒否 |

stdoutはプログラム出力、stderrは診断・ヘルプです。ファイル末尾の改行やechoの空白を補正しません。
PHPの `-r`・`-l`・`--version`、stdinによるプログラム実行、複数ファイル、スクリプト引数には対応していません。
`$argc`・`$argv`・スーパーグローバルも拒否します。php.iniや拡張は読み込みません。

Go製インタプリタ自体のコンテナ構成は [containers/gophp](../containers/gophp/README.md) にあります。
同じCLI仕様・対象設定で実行し、`run.sh` はホストと同じ絶対パスを読み取り専用でマウントします。
コンテナ版のチェックは保存済み期待値との照合です。2026-10-12のエージェント環境では
Podman接続が拒否され、コンテナの実ビルド・実行・公式PHPとのコンテナ版実比較は未実施です。

## バージョン差と対応状況

| 差異 | Go側の扱いと根拠 |
| --- | --- |
| 8.4からの括弧付きnew直後のアクセス | 8.3は構文エラー。8.4 / 8.5は許可。既存fixtureで過去に実比較済み |
| 関数重複定義 | 8.3は `Cannot redeclare`、8.4 / 8.5は `Cannot redeclare function`。既存3fixtureは過去に実比較済み。新しい組み込み名衝突もユーザー提供ログで3版とも一致 |
| コンパイルエラーと循環比較 | 8.5のスタックトレース・エラー形式を対象設定で変更。既存fixtureで過去に実比較済み |
| optional引数をrequiredの前に置く宣言 | 8.4 / 8.5の診断に関数名を追加。過去に実比較済み |
| NaNのbool / string変換 | 8.5でWarning。既存の条件・論理・連結は過去に実比較済み。追加したキャスト・empty・strlenはユーザー提供ログで3版とも一致 |
| floatからintへの範囲外変換 | 8.5でWarning。NaNの暗黙変換では精度喪失Deprecatedも残す。明示的キャスト・isset / unsetの添字はユーザー提供ログで3版とも一致 |
| nullの配列キー | 8.5の読取り・書込み・isset / emptyはDeprecated。UNSET_DIMでは警告しない。追加した存在確認・削除はユーザー提供ログで3版とも一致 |
| 非数値文字列の増分 | 8.5のDeprecated文言を切り替える。既存fixtureで過去に実比較済み |

共通のPHPソースをすべての版に渡します。期待値の版差は `stderr_by_version`・`stdout_by_version`・
`exit_by_version` に明示し、比較時に診断を捨てたり出力をtrimしたりしません。

## 検証方法と今回の結果

```sh
make test vet build
make compare-all
make compare-parser
```

[php-go-interpreterスキル](../skills/php-go-interpreter/SKILL.md) とその比較スクリプトを使用します。
公式PHPはiniを読み込まず、診断をstderrへ送り、ログを無効化し、タイムゾーンUTCで実行します。
各fixtureを別プロセスで動かし、stdinは空、実比較のタイムアウトは15秒です。
stdout・stderrのバイト列と終了コードの3つを比較し、どれかが違えばDIFFになります。
正常系fixtureが両処理系で同じエラーになってもMATCHにはせず、非ゼロ終了を検証エラーにします。
公式PHP側の終了125に加えて、Podman固有の接続エラーまたはstderr先頭の
`Error: can only create exec sessions on running containers` を確認した場合はPHPを実行できなかった検証エラーとし、
stdout・stderr・終了コードを表示してその版の比較を停止します。終了125だけなら通常どおり厳密比較します。
この環境エラー判定は公式PHP側に限り、Go側の同じ終了コード・診断は通常の比較結果として扱います。
PHP側のタイムアウト・実行障害では版確認を再度試し、成功すれば残りを続行、失敗すれば停止します。
停止時は残りの件数を `not run` と明示し、一致数へ含めず終了コード2にします。全件成功には再実行が必要です。

既存337件はPHPファイルのSHA-256と保存済み期待値を変更していません。
新規129件（正常48・実行時72・構文／コンパイル9）を加え、計466件です。
現在の分類は正常208・実行時198・構文／コンパイル60。fixtureの削除・除外はありません。
参照foreach・オブジェクト反復・キャスト別名などのGo独自拒否はGo単体テストに置いています。
既存の「foreachを未対応として拒否」するテストは「参照foreachを拒否」するテストへ更新しました。

Goテスト、Pythonテスト11件、go vet、ビルドは成功しました。
Goテストでは既定・8.3・8.4・8.5の設定で全fixtureを照合します。
ビルド済みCLIでも5秒のタイムアウトを付けて全件照合し、空白・日本語のパス、ヘルプ、引数、終了コードを確認しました。
Pythonテストは比較スクリプトの厳密比較・正常終了チェック・版指定・エラー集約と、接続障害／タイムアウト後の未実行の扱いを検査します。PHP実比較の代用ではありません。
コンテナ停止の検出と誤検出防止の回帰テスト2件を追加し、既存9件を維持しました。
停止の模擬では、最初のfixtureが一致した後にPHP側が終了125と停止診断を返し、
そのfixtureをERRORとし、残りを実行せず `not run` に計上することを確認します。
PHP自身の終了125、終了コードが異なる同じ文言、stdout上の文言、Go側の診断では通常の比較を続けます。
この比較器の修正は以下の全466件×3版のユーザー提供実比較より後に行い、修正後の検証には模擬プロセス結果を使用しています。

**保存済み期待値とのCLI照合（公式PHPとの実比較ではない）：**

| Go対象 | 件数 | 一致 | 不一致 | 検証エラー | 正常系の終了0 |
| --- | ---: | ---: | ---: | ---: | ---: |
| 8.3 | 466 | 466 | 0 | 0 | 208 |
| 8.4 | 466 | 466 | 0 | 0 | 208 |
| 8.5 | 466 | 466 | 0 | 0 | 208 |

**ユーザー提供の公式PHP比較（補助関数改名・コンテナ起動後の全件実行）：**

| PHP | 対象件数 | 一致 | 不一致 | 検証エラー | 正常系の両処理系終了0 |
| --- | ---: | ---: | ---: | ---: | ---: |
| 8.3.35 | 466 | 466 | 0 | 0 | 208 |
| 8.4.26 | 466 | 466 | 0 | 0 | 208 |
| 8.5.11 | 466 | 466 | 0 | 0 | 208 |

ユーザーは停止していた `php-version-83`・`php-version-84`・`php-version-85` を起動し、
`check.sh` の3版PASSに続き `make compare-all` を実行しました。
Goテスト・Pythonテスト9件・ビルドもその実行で成功しています。go vetは先のローカル検証で成功済みです。
修正した `success/isset_short_circuit_order.php` を含め全件一致し、未実行のfixtureはありません。
この実比較はユーザー環境のログを確認した結果です。エージェント環境からの実行結果とは扱いません。
[提供ログの抜粋と完成条件](validation/v1-user-comparison-complete.md) を参照してください。

全件成功前のユーザー実比較（8.3・8.4が各465一致、8.5が179一致後に接続障害）と
補助関数の改名・比較器の修正は [追補修正の履歴](validation/v1-user-comparison-followup.md) に残しています。
エージェント環境での公式比較試行は3版とも版確認時に接続を拒否され、fixtureの実比較は0件でした。
過去の失敗・未実行を一致へ読み替えず、今回の全件再実行で完成を確認しています。

[修正後のローカル検証ログ](validation/v1-followup-local-validation.txt) と
[修正後の実比較試行ログ](validation/v1-followup-php-comparison.txt) を保存しています。
初回の [ローカル検証](validation/v1-local-validation.txt)・[実比較試行](validation/v1-php-comparison.txt) は履歴として残しています。
新規129件の期待値は公式資料と対象版PHPソースから手動で作成しました。Goの出力や公式PHP実行から採取したものではありません。
今回のユーザー提供実比較で、新規129件と既存337件の全件を3版とも確認しました。
第1版の完成条件に関する未実施項目はありません。確認した互換性の範囲は収録したfixtureと以下の対応範囲・制約に限ります。

## 第1版時点の未対応機能と既知の制約

この一覧は第1版時点の履歴である。その後に実装したswitch/match、名前空間、ビット演算、文字列展開等の対応範囲は[第2段階の互換性文書](phase2-compatibility.md)を参照。当時の527 fixtureは3版とも一致し、その後の追加分も現在の588件の比較に含まれる。Laravel/API PlatformのGo実行は未対応である。

- 参照（代入・引数・戻り値・foreach）、resource、Countable / Iterator / ArrayAccess、オブジェクト反復。
- 配列展開・分解代入、文字列添字の書き込み、動的な変数／プロパティ／メソッド名。
- 可変長／名前付き引数、型宣言、既定値のnew、クロージャ・変数関数、global / static変数。
- 継承・interface・trait、staticメンバー・クラス定数・属性・匿名クラス・clone・名前空間・autoload、__construct以外のマジックメソッド。
- switch / match、例外処理、include / require、eval、print、and / or / xor、ビット演算・累乗、代替コロン構文。
- 上表以外の標準関数・定数と定数定義、キャスト別名・配列／オブジェクトへのキャスト。
- 16 / 8 / 2進数、数値の桁区切り、文字列補間・heredoc / nowdoc、インラインHTML・複数PHPブロック。
- Web実行、C拡張、JIT、php.ini、PHP全体の標準実行環境。

通常関数・メソッドの合計呼び出し深度は256で、超過するとgophp独自のErrorです。
CLIにループ回数制限はありません。fuzzテストは内部の実行予算、比較スクリプトはプロセスのタイムアウトで停止します。
配列更新時のコピーは性能優先のZend実装とは異なり、大きな配列や深いネストの性能互換を目指しません。
PHPのメモリ制限・スタック制限・ini依存の書式や環境設定を再現していません。
収録していない入力、複雑な式で同じ値を書き換える評価順序、浮動小数点の全境界、構文エラー全種類の文言には未検証領域があります。
未対応構文はstderr診断と非ゼロ終了で拒否し、文言にはgophp独自のものが残ります。
組み込み名の重複定義チェックは対応する11関数に限定しています。PHPの対象外標準関数名（`key` など）を
Goでユーザー関数名に使用した場合の互換性も対象外です。fixtureの補助関数には標準関数と衝突しない名前を使います。
収録した466件では3版ともstdout・stderr・終了コードが一致しました。PHP全機能や未収録の入力の互換性を保証するものではありません。

## 今回参照した公式資料

- [foreach](https://www.php.net/manual/en/control-structures.foreach.php)、[配列](https://www.php.net/manual/en/language.types.array.php)
- [isset](https://www.php.net/manual/en/function.isset.php)、[empty](https://www.php.net/manual/en/function.empty.php)、[unset](https://www.php.net/manual/en/function.unset.php)
- [PHP 8.4.26のコンパイル](https://github.com/php/php-src/blob/php-8.4.26/Zend/zend_compile.c)（foreachの代入順序・専用言語構造・標準関数の専用命令）
- [PHP 8.4.26の引数処理](https://github.com/php/php-src/blob/php-8.4.26/Zend/zend_API.c)（個数・型の診断、弱い数値／文字列変換）
- [PHP 8.4.26のcount](https://github.com/php/php-src/blob/php-8.4.26/ext/standard/array.c)
- [PHP 8.5.11の実行命令](https://github.com/php/php-src/blob/php-8.5.11/Zend/zend_vm_def.h)（isset / empty / unset、strlen / count、foreach）
- [PHP 8.5.11の添字処理](https://github.com/php/php-src/blob/php-8.5.11/Zend/zend_execute.c)、[診断](https://github.com/php/php-src/blob/php-8.5.11/Zend/zend.c)
- [PHP 8.5.11の型変換](https://github.com/php/php-src/blob/php-8.5.11/Zend/zend_operators.c)、[float整数変換](https://github.com/php/php-src/blob/php-8.5.11/Zend/zend_operators.h)

資料からの実装判断と保存済み期待値の一致は、公式PHPでプログラムを実行して得た結果と区別しています。

## 項目6〜10の追加範囲（662件時点の履歴）

現在のGo実装は、ローカルファイル読み込み、継承・interface・trait・staticメンバー・型なしクラス定数、単純な引数・戻り値型、ファイル単位strict_types、値捕捉クロージャ・arrow関数、可変長・名前付き引数、配列の整数／文字列キー展開、基本的な例外制御を**部分対応**とする。読み込み先のParseError捕捉、読み込みを含むトレース、失敗別の`_once`履歴を追加した。通常の`include`／`require`は解析成功後、`*_once`は解析前に読込履歴へ登録する。例外連鎖関連の新規21件を含む662件は[修正後の3版全一致ログ](validation/throwable-chain-662-compare-all-mac.log)で検証済み。高度な型、参照捕捉、Traversable展開、残るThrowable APIと継承規則、エラーハンドラ、include_path全般・ストリームラッパー、Composer生成autoload、Reflection・Attribute、HTTP・DB統合はこの記録時点では残件だった。後続のエラーハンドラは684件で比較済み、include_pathのローカル探索は700件時点では公式比較待ちだった。その後の702件全一致は冒頭の成功ログで確認済み。項目6・10および項目6〜10全体は部分実装。[進捗と制限](items6-10-progress.md)を参照。上の第1版時点の未対応一覧と、以前の[538件時点の試行ログ](validation/items6-10-compare-all.log)は履歴である。
