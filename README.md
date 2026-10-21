# gophp — 機能を限定したPHP CLI処理系

固定Laravel・API Platformアプリの[全6段階の完成判定](docs/full-stack-development.md)は、定義された対象範囲で達成した。公式PHP 8.3.35・8.4.26・8.5.11との版別比較は各1173件一致・差異0・エラー0・終了0。Go製PHPの実TCPサーバーでも3設定すべてでJSON・入力検証・SQLite CRUD・API Platform・OpenAPI・2000回連続／8並列要求が成功し、検証コマンドは各版終了0。PHP全機能の対応完了を意味しない。

追加した[Laravelタスク管理API](research/task-api/README.md)は、SQLite・Eloquent・標準入力検証を使う認証なしのCRUD。公式PHP／Goの3組で各191要求の実TCP比較を行う。[実装時の証拠](docs/validation/task-api-results.md)と[リリース再構築・CIの記録](docs/task-api-release.md)を区別して保存する。新規環境では両アプリの固定lockからvendorを導入し、`scripts/verify-release.py`で全fixture比較・タスクHTTP比較・既存アプリの2000連続／8並列GETまで実行できる。

以下の機能別記録に残る「比較待ち」と旧件数は当時の履歴であり、現行1173件の公式比較は完了している。

`yield from`で`IteratorAggregate`を24段委譲する有限の連鎖を、Go固有の深さ16制限で拒否しないよう修正した。[回帰fixture・公式比較待ち](docs/full-stack-development.md)。

`array_walk_recursive()`の有限配列を深さ64で誤って循環と判定していた制限を修正した。70段の回帰ケースと固定Laravel Translatorのネストした置換をGoの3設定で確認した。[詳細・公式比較待ち](docs/full-stack-development.md)。

`extract(..., EXTR_REFS)`で入力配列と同名のキーを取り込む場合も、元配列の参照セルを保持して後続要素を処理するよう修正した。[回帰fixture・検証待ち](docs/full-stack-development.md)。

`array_walk()`／再帰版の参照コールバックで、受け手・添字を1回評価した配列格納先を使うよう修正した。[回帰fixture・検証待ち](docs/full-stack-development.md)。

`extract(..., EXTR_REFS)`の参照可能な配列引数を既存の参照セルに接続し、Goの3設定で基本動作を確認した。[回帰fixture・制限・公式比較待ち](docs/full-stack-development.md)。

PDO SQLiteの接続・Statement別`errorCode()`／`errorInfo()`を追加し、Goの3設定で基本の状態遷移を確認した。[対応範囲と未比較範囲](docs/full-stack-development.md)。

`PDOStatement::bindValue()`の`PARAM_INT`・`PARAM_BOOL`・`PARAM_STR`・`PARAM_LOB`・`PARAM_NULL`とSQLite BLOB書き込みをGoで追加し、登録時の文字列変換と`execute(array)`の再束縛をローカル検証した。[記録](docs/full-stack-development.md)。

最新のDB経路では、PDO SQLiteのSQLSTATE／`errorInfo`とLaravelの一意制約例外をGoの3設定でローカル確認した。[対応範囲と公式比較待ち](docs/full-stack-development.md)。

1167件時点の履歴（2026-10-21）：fixture・期待値は**1167件**（正常743・実行時290・構文134）。公式PHP3版と各997件一致した版別生ログが最後の比較済み基準で、追加170件と変更後全件はPodman接続拒否で公式比較未実施、期待値は暫定。Go製PHPでは通常オブジェクトの`foreach`、`IteratorIterator`・`LimitIterator`・`AppendIterator`・`CachingIterator`・`iterator_count()`、Generatorの基本形、通常関数・メソッド・Closureの参照戻り値、委譲中の`send()`・`throw()`、Iterator委譲から固定vendorの起動経路を進め、`artisan list`・`about`・`route:list`・独自`gophp:probe`・一時SQLiteでの`migrate --force`とOpenAPI exportの終了0、固定アプリのHTTP/SQLite/API Platform基本CRUD、および同一プロセスでの8回の追加CRUD・[512回](docs/validation/full-stack-http-512-after-callback.log)・[1024回の追加GET](docs/validation/full-stack-http-1024-caching.log)・JSON-LDコレクション形式の[ページ分割](docs/validation/full-stack-pagination-jsonld.log)・8並列GETを一時DBのローカル統合テストで確認した。[現在の範囲・ログ・制限](docs/full-stack-development.md)。項目11と全6段階の完成条件は未完了。

HTTP入口ではGETの配列形式クエリとPOSTのフォーム本文を`parse_str()`と共通の解析経路へ接続した。POST本文の解析はメソッドとContent-Typeを検査する。[対象テストと未検証の境界](docs/full-stack-development.md)を参照。公式PHPとのHTTP経路の実比較は未実施。

固定アプリのHTTP・SQLite経路は、同一Goプロセスの[2000リクエスト継続試験](docs/validation/full-stack-http-2000-stress-local.log)でもローカル終了0だった。運用環境での負荷とPHP3版比較は引き続き未検証。

PDOStatementの未定義パラメータ実行ではSQLiteの`HY000`をStatementのエラー状態へ記録し、成功再実行でクリアする[回帰fixture](tests/fixtures/success/pdo_parameter_error_state.php)を加えた。Go3設定の対象実行と全ローカルテストは通過したが、公式PHP3版の実比較は待機中。[実装・ログ・制限](docs/full-stack-development.md)を参照。

固定アプリを実際の`gophp --serve`プロセスで検証する[版別スクリプト](research/laravel-api/verify-gophp.py)を追加した。3設定の別プロセスによるCLI migration・独自コマンド・OpenAPI出力・DBテーブル作成は成功した。この環境では127.0.0.1の待ち受けが拒否されるため、実サーバーのHTTP・SQLite・API Platform応答は未収録。[Macでの実行手順](docs/full-stack-development.md)を参照。

以前の状態（履歴）：[公式PHP3版と各997件一致した生ログ](docs/item11-progress.md)を基準に、参照・反復の入口、匿名クラス、配列分解代入、基本マジックメソッド、autoloadに必要なファイル確認、CLI引数などを追加した。fixture・期待値は**1035件**（正常613・実行時290・構文132）。追加38件と全件の公式比較はPodman接続拒否で未実施、期待値は暫定。実`vendor/autoload.php`はGoで終了0、`artisan list`はGenerator構文で停止している。[全6段階の進捗と残件](docs/full-stack-development.md)。項目11とアプリ完成条件は未完了。

今回の追加前の比較済み基準：[PHP 8.3](docs/validation/reference-property-974-83-mac.log)・[8.4](docs/validation/reference-property-974-84-mac.log)・[8.5](docs/validation/reference-property-974-85-mac.log)の版別生ログは各974件の個別MATCH、**974一致・0不一致・0エラー、終了0**を記録する。959件から追加した型なしプロパティ参照15件は収録fixtureの範囲で公式比較済み。`check.sh`のPASSはユーザー報告。版ごとの`make compare`であり、`make compare-all`成功記録ではない。以下の974件「未比較」は比較前の履歴。

974件比較前の基準（履歴、2026-10-18）：[版別生ログ](docs/item11-progress.md)で公式PHP3版が各959一致・0不一致・0エラー、終了0。944件以降の追加15件と訂正した期待値は収録範囲で比較済み。以下の959件「未比較」は比較前の履歴。

959件比較前（履歴）：呼び出し元のローカル変数に残る配列要素参照を共有判定へ加え、fixture・期待値は**959件**（正常539・実行時288・構文132）。Goローカル検証済み。公式PHP3版の対象・全件比較はPodman接続拒否で未実施、新規期待値は暫定。[修正・試行ログ](docs/item11-progress.md)。項目11全体は部分実装。

今回の修正前（履歴、2026-10-18）：配列コピーの参照分離を修正し、fixture・期待値は**958件**（正常538・実行時288・構文132）。Goローカル検証は完了したが、公式PHP3版はPodman接続拒否で未比較。新規3件と訂正した既存1件の期待値は暫定。[修正・公式ソース根拠・試行ログ](docs/item11-progress.md)。項目11全体は部分実装。

958件への修正前（履歴、2026-10-18）：配列要素への参照を追加し、fixture・期待値は**955件**（正常535・実行時288・構文132）。Goローカル検証は成功したが、新規11件と全955件の公式PHP3版比較はPodman接続拒否で未実施、期待値は暫定。[項目11の設計・ログ・再実行手順](docs/item11-progress.md)。項目11全体と項目6〜10の全仕様対応は未完了。

955件追加前の基準（履歴）：fixture・期待値は**944件**（正常524・実行時288・構文132）。公式PHP 8.3.35・8.4.26・8.5.11との[版別生ログ](docs/item11-progress.md)で各版**944一致・0不一致・0エラー、終了0**を確認した。変数間の参照代入を含む追加7件の比較待ちは解消した。今回は保存済みログの確認と文書更新で、比較を再実行していない。対応は変数間の`=&`に限り、項目11全体と項目6〜10の全仕様対応は未完了。[対応範囲と残件](docs/item11-progress.md)。

944件比較前の状況（履歴、2026-10-18）：公式PHP3版の[版別生ログ](docs/items6-10-integration-gate.md)で既存**937件が各版937一致・0不一致・0エラー、終了0**と確認でき、項目11へ移行可能となった。項目11では変数間の参照代入を追加し、fixture・期待値は**944件**（正常524・実行時288・構文132）。Goローカル検証は成功したが、新規7件と修正後全件の公式PHP比較はPodman接続拒否で未実施、期待値は暫定。[参照の設計・制限・検証ログ](docs/item11-progress.md)。項目6〜10および項目11全体は部分実装。

937件比較前の履歴（2026-10-18）：`define()`第3引数の`null`変換中に非推奨診断のハンドラが投げた場合、定数登録前に中断するよう修正した。`true`指定時の警告は別経路として維持する。fixture・期待値は**937件**（正常517・実行時288・構文132）。Goの[test](docs/validation/warning-audit-third-null-test.log)・[vet](docs/validation/warning-audit-third-null-vet.log)・[build](docs/validation/warning-audit-third-null-build.log)は終了0。公式PHP3版の対象・全件比較はPodman接続拒否で未実施、追加2件の期待値は暫定。[監査記録・比較ログ・項目11への移行条件](docs/items6-10-integration-gate.md)。項目6〜10全体は部分実装。

936件時点の履歴（2026-10-18）：警告中断経路を追加監査し、`define(null, ...)`の診断ハンドラが例外を投げた後に定数を登録するGo側の経路を修正した。fixture・期待値は**936件**（正常516・実行時288・構文132）。[935件の版別生ログ](docs/items6-10-integration-gate.md)はこの修正前の基準。新規1件と修正後全936件の公式PHP3版比較はPodman接続拒否で未実施、期待値は暫定。Goの[test](docs/validation/warning-audit-test.log)・[vet](docs/validation/warning-audit-vet.log)・[build](docs/validation/warning-audit-build.log)は終了0。

936件追加前の履歴（2026-10-18）：fixture・期待値は**935件**（正常515・実行時288・構文132）。PHP 8.3.35・8.4.26・8.5.11の[版別生ログ](docs/items6-10-integration-gate.md)は各版935 MATCH・0不一致・0エラーを記録し、追加した`.=`対比も一致した。版別コマンドの終了0と`check.sh`のPASSはユーザー報告で、生ログには記載がない。型付きプロパティの複合代入エラー後の値は演算と左辺型によって異なり、`.=`も一律に「値を保つ」とは扱わない。公式比較の移行条件は満たしたが、未収録の警告中断経路の監査が項目11移行前の別条件として残る。[判断と履歴](docs/items6-10-integration-gate.md)。項目6〜10全体は部分実装。

934件時点の履歴（2026-10-18）：型付きプロパティの`/= 0`・`%= 0`・`+= []`・警告中断の`+= '2x'`に対するGo修正後、Podman接続拒否で公式PHPとの再比較ができなかった。[修正前の934件集計](docs/validation/items6-10-compound-934-before-fix-summary-mac.log)は各版933一致・1不一致・0エラー。[933件一致の集計](docs/validation/items6-10-compound-933-after-fix-summary-mac.log)も保存した。両集計に個別stdoutと終了コードはない。

今回の追加前（履歴、2026-10-17）：fixtureは**933件**（正常513・実行時288・構文132）。`%=`の整数変換で警告ハンドラが投げた例外を伝播し、変数・配列要素は元の値を保持し、型付きプロパティは結果0を格納する動作だった。Goの3版設定でローカル検証済み。公式PHP3版の対象・全件比較はPodman接続拒否でfixture実行前に停止し、各版0一致・0不一致・1環境エラー・933件未実施。追加2件の期待値は暫定。[項目11への移行条件](docs/items6-10-integration-gate.md)。

統合検証の前段階（履歴、2026-10-17）：fixtureは**931件**（正常511・実行時288・構文132）。項目6〜10の既存機能を組み合わせた7件を追加し、警告ハンドラの例外後に動的プロパティなどの格納が進む不具合を修正した。Goローカル検証は通過したが、今回の公式PHP3版比較はPodman接続拒否で未実施。新規期待値は暫定で、直前の[924件版別一致集計](docs/validation/callable-924-after-fix-summary-mac.log)は変更前の履歴。[統合検証一覧・残件・項目11への移行判定](docs/items6-10-integration-gate.md)を参照。項目6〜10全体は部分実装のまま。

統合検証前の検証履歴（2026-10-17）：fixtureは**924件**（正常504・実行時288・構文132）。ユーザーが版ごとに`make compare`を実行した[6行の集計ログ](docs/validation/callable-924-after-fix-summary-mac.log)で、PHP 8.3.35・8.4.26・8.5.11はそれぞれ**924一致・0不一致・0エラー**。908件から追加したcallable／Closure関連16件の公式比較待ちは、収録fixtureの範囲で解消した。`check.sh`のPASSはユーザー報告であり、ログには個別fixture行と各コマンドの終了コードがない。今回はログを保存・照合したもので、比較を再実行していない。[対応範囲と残件](docs/items6-10-progress.md)を参照。項目9と項目6〜10全体は未完了。

924件追加前の検証履歴（2026-10-16）：fixtureは**908件**（正常490・実行時287・構文131）。ユーザーがPHP 8.3.35・8.4.26・8.5.11で版別に`make compare`を実行した[修正後の集計ログ](docs/validation/special-types-908-after-fix-summary-mac.log)は、各版**908一致・0不一致・0エラー**を記録する。888件から追加した20件の公式比較待ちは、収録fixtureの範囲で解消した。ログは6行の集計で、個別fixture行・各コマンドの終了コードは含まれない。修正6件の解消と`check.sh`の各版PASSはユーザー報告であり、今回は保存済みログの確認・保存のみ行った。[修正前の結果と残件](docs/items6-10-progress.md)を参照。項目8と項目6〜10全体は未完了。

修正前の記録（履歴）：[908件の比較](docs/validation/special-types-908-before-fix-mac.log)は8.3.35・8.4.26が各902一致／4不一致／2テストエラー、8.5.11が903一致／3不一致／2テストエラーだった。前方型参照の2件を失敗系へ移し、宣言時の型解決と診断をGo側で修正した。修正直後はPodman接続拒否で公式PHP比較を実施できず、期待診断を暫定とした。後続の上記集計で収録fixtureの比較待ちは解消した。

修正前の追加時点（履歴）：`self`／`parent`／`static`の宣言スコープに基づく型解決と、autoloadを使わない前方型参照の互換性検査をGo側に追加した。当時は正常492・実行時285・構文131の908件だった。既存888件は[版別集計](docs/validation/enum-888-after-fix-summary-mac.log)で公式PHP3版の各888件一致を確認済み。新規20件は当時[8.4の版別比較試行](docs/validation/special-types-compare-84-attempt.log)などがPodman接続拒否で終了2となり、公式比較待ちだった。その後、上記の修正前比較が得られた。[進捗と再実行手順](docs/items6-10-progress.md)を参照。

enum比較時点の履歴（2026-10-16）：fixtureは**888件**（正常477・実行時281・構文130）。ユーザーがMacで版ごとに`make compare`を実行した後の[修正後の集計ログ](docs/validation/enum-888-after-fix-summary-mac.log)には、PHP 8.3.35・8.4.26・8.5.11の各版で**888一致・0不一致・0エラー**と記録されている。878件以降の追加11件の公式比較待ちは収録fixtureの範囲で解消した。ログは6行の集計で、個別fixture行と各コマンドの終了コードは含まれない。修正4件の解消と`check.sh`の各版PASSはユーザー報告であり、今回は比較を再実行していない。[修正前の4件不一致](docs/validation/enum-888-before-fix-mac-summary.log)と[進捗・残件](docs/items6-10-progress.md)を参照。項目8と項目6〜10全体は未完了。

前回の884件時点（2026-10-16）：enumの通常／backed宣言、caseの同一性、`name`／`value`、`cases()`、`from()`／`tryFrom()`、基本的な型・interface・trait連携を部分実装した。fixtureは**884件**（正常475・実行時279・構文130）。既存878件は[版別集計](docs/validation/readonly-878-after-fix-summary-mac.log)で公式PHP 8.3.35・8.4.26・8.5.11とそれぞれ一致。追加6件は[Podman接続拒否](docs/validation/enum-compare-all-attempt.log)により公式比較未実施で、期待値は暫定である。項目8と項目6〜10全体は未完了。[進捗と制限](docs/items6-10-progress.md)を参照。

enum追加前の878件時点（2026-10-16）：コンストラクタのプロパティ昇格とreadonlyの基本範囲をGo側へ追加した。fixtureは**878件**（正常471・実行時278・構文129）。ユーザーのMacで版ごとに`make compare`を実行した後の[集計ログ](docs/validation/readonly-878-after-fix-summary-mac.log)では、PHP 8.3.35・8.4.26・8.5.11の各版で**878一致・0不一致・0エラー**。851件以降の追加27件も収録範囲で比較済みとなった。このログには個別fixture行とコマンド終了コードがなく、指摘2件の解消と`check.sh`のPASSはユーザー報告である。今回はログを保存しただけで比較を再実行していない。[修正前の876一致・1不一致・1テストエラー](docs/validation/readonly-878-before-fix-mac.log)は履歴として維持する。項目8と項目6〜10全体は未完了。[実装範囲・検証履歴](docs/items6-10-progress.md)を参照。

851件時点の履歴：動的プロパティ代入の受け手再束縛を修正した。[修正前ログ](docs/validation/dynamic-property-851-before-fix-mac.log)は3版共通で1件不一致。修正後の3版各851件一致は上記の別実行による確認で、こちらで比較を再実行した結果ではない。

Goで字句解析・構文解析・評価を行うPHPインタプリタです。基準は **PHP 8.4.26 / 64ビット環境**。
同じGo実装とPHP fixtureを、8.3.35・8.4.26・8.5.11と比較します。実行時に公式PHPを呼び出しません。

2026-10-15時点の確定履歴：**817件**（正常433・実行時272・構文112）。ユーザーがMacで実行した[修正後の全件比較](docs/validation/void-never-817-after-fix-mac.log)は、PHP 8.3.35・8.4.26・8.5.11の各版で**817一致・0不一致・0エラー**、終了コード0でした。787件以降に追加した30件の公式比較待ちは、収録fixtureの範囲で解消しました。今回は保存済みログの確認・保存であり、比較を再実行していません。[修正前の3件の不一致と残件](docs/items6-10-progress.md)を参照してください。

787件の追加直後（履歴）：Go側のローカル検証は成功しましたが、[当時の全件比較試行](docs/validation/type-variance-float-compare-all-attempt.log)はPodman接続拒否で終了2でした。この時点での「新規18件は未実施・期待値暫定」という判定は、上記の後続成功ログで解消しています。

782件時点の開発記録（履歴）：継承メソッドの引数反変・戻り値共変をGo側へ追加し、新規13件はGoローカル検証済みでした。[当時の比較試行](docs/validation/type-variance-compare-all-attempt.log)はPodman接続拒否で終了2。宣言時の`int`と`float`の関係には、後続レビューで上記の修正が必要になりました。

765件時点の開発記録（履歴）：追加26件はGoの保存済み期待値テストを通過しましたが、[当時の比較試行](docs/validation/dnf-compare-all-attempt.log)はPodman接続拒否で終了2でした。この時点ではintersection同士の包含関係を宣言時に検査していませんでした。

750件時点の確定履歴：ユーザーがMacで実行した[修正後の全件比較](docs/validation/union-diagnostic-750-after-fix-mac.log)で、PHP 8.3.35・8.4.26・8.5.11の各版が**750一致・0不一致・0エラー**、終了コード0でした。[対象比較](docs/validation/union-nullable-syntax-targeted-mac.log)でも`union_nullable_syntax.php`が各版1件一致しました。727件から追加した23件の公式比較待ちは、収録fixtureの範囲で解消しています。この記録はMacの保存済みログを確認・保存したもので、比較を再実行していません。[修正前の不一致と残件](docs/items6-10-progress.md)を参照してください。項目8全体、項目6〜10全体、Go製PHPでのLaravel起動は未完了です。

727件時点の確定履歴：ユーザーがMacで実行した[公式PHP3版の全件比較](docs/validation/static-throw-727-compare-all-mac.log)は、8.3.35・8.4.26・8.5.11で各727一致・0不一致・0エラー、終了コード0でした。直接の`Class::{throw ...}()`を事前に拒否する修正と追加9件は、収録fixtureの範囲で公式比較済みです。この記録は保存済みログの確認・保存であり、今回の比較結果ではありません。

718件比較時点の履歴：**修正前に比較済みのfixtureは718件**（正常393・実行時250・構文75）。ユーザーがMacで実行した[修正後の全件比較](docs/validation/item7-dynamic-static-718-after-fix-mac.log)は、PHP 8.3.35・8.4.26・8.5.11で各718一致・0不一致・0エラー、終了0でした。追加11件の公式比較待ちは解消しました。今回は保存済みログを確認・保存したもので、比較を再実行していません。`Class::{throw ...}()`をGoが受理する差は、この718件の比較時点では収録fixture外の既知不具合でした。[修正前の比較](docs/validation/item7-dynamic-static-718-before-fix-mac.log)と[進捗・残件](docs/items6-10-progress.md)を参照してください。項目7・項目6〜10全体とGo製PHPでのLaravel起動は未完了です。Laravel/API Platform調査アプリは公式PHPの3版でHTTP・SQLite・API作成/取得が成功しています。[証拠と残件](docs/validation/laravel-api-evidence-2026-10-13.md)を参照してください。

727件追加直後はPodman接続制限で公式比較できず、期待値を暫定としていました。この[当時の試行と後続の成功](docs/items6-10-progress.md)は履歴として区別します。
[今回の対応範囲](docs/phase2-compatibility.md)、[依存調査](docs/dependency-audit.md)、[Laravel/API Platformの初期目標](docs/laravel-api-target.md)、[実行時設計](docs/phase2-design.md) を参照してください。

## 実行

Go 1.22以上、Python 3、CコンパイラとSQLite・PCRE2開発ライブラリを使用します。開発時の確認環境はGo 1.27.0・Python 3.9.6・darwin/arm64です。固定アプリの依存導入とCIの手順は[v0.3.0の再構築・CI・公開記録](docs/task-api-release.md)を参照してください。

```sh
cd /Users/kouji/dev/php
make build
./bin/gophp examples/report.php
./bin/gophp --php-version=8.5 examples/report.php
./bin/gophp --help
```

[集計スクリプト](examples/report.php)は、配列を関数に渡してforeachで選別し、結果の配列を返します。
isset・empty・unset・キャスト・標準関数を使い、元配列の保持とStatsオブジェクトの共有も確認します。
出力は以下です。

```text
a:tea:12:3
c:水:3:3
2:2:1:0:integer
```

CLIは `gophp [--php-version=8.3|8.4|8.5] <file.php> [arguments...]`。既定は8.4です。
標準出力はプログラムの出力、標準エラーは診断とヘルプです。
終了コードは成功0、ファイル読み込み等の失敗1、引数の誤り2、PHP構文・致命エラー255です。
ファイル名の前に `--` を指定できます。PHPの `-r`・`-l`、stdin実行は未対応です。

## Go製PHPをコンテナで実行

```sh
make container-build
make container-check
make container-run
make container-run GOPHP_TARGET=8.5 GOPHP_FILE=tests/fixtures/success/hello.php
```

PodmanでGo製インタプリタ自体をビルドし、PHPファイルをマウントして実行します。
podman-composeでの実行と全fixtureのチェックは [gophpコンテナの案内](containers/gophp/README.md) にあります。
このCodexセッションではPodman接続が拒否されたため、コンテナの実ビルド・実行は未確認です。

## 対応範囲

| 項目 | 内容 |
| --- | --- |
| 値・演算 | int・float・string・bool・null・配列・オブジェクト、四則演算・剰余・連結・比較・論理短絡、三項・null合体・ビット演算 |
| 代入 | 式内・連鎖・複合代入、`??=`、ビット複合代入、前置／後置の増減、配列要素・プロパティの更新 |
| 制御 | if / elseif / else if / else、switch、match、while / for / do while / foreach、レベル付きbreak / continue、return |
| 関数・クラス | 通常関数、位置引数・既定値・再帰・ローカルスコープ、`global`・`static`ローカル、文字列と二要素配列の一部の変数callable。継承・interface・trait・staticメンバー・型なしクラス定数・`self`/`parent`/`static`、動的クラス名の部分集合 |
| 型宣言 | 単純・複合の引数／戻り値型、型付きプロパティ／クラス定数の基本形、ファイル単位の`declare(strict_types=0/1)`、型変換とTypeErrorの部分集合。動的インスタンスプロパティ名は収録fixtureの範囲で公式比較済み |
| 引数・例外 | 名前付き引数、整数／文字列キー配列の展開、可変長引数、throw・try/catch/finally、基本的なThrowable階層の部分集合。収録fixtureの範囲で公式PHP3版と比較済み |
| 配列 | 作成・キー変換・挿入順・添字・追加・ネスト・値コピー・削除・値渡しforeach |
| 言語構造 | isset（複数対象・短絡）、empty（通常式も可）、unset |
| キャスト | `(bool)`・`(int)`・`(float)`・`(string)` |
| 標準関数 | strlen、count、gettype、is_int、is_float、is_string、is_bool、is_null、is_array、is_object、is_numeric |
| 文字列 | 二重引用符内の限定した変数展開、heredoc、nowdoc |
| 名前・定数 | 名前空間と各種use、`__NAMESPACE__`、`Foo::class`、`const`、`define` / `defined` / `constant`、COUNT_NORMAL / COUNT_RECURSIVE |
| 複数ファイル | ローカルファイルの`include` / `require` / `*_once`、`__FILE__` / `__DIR__`、`include_path`探索と`get_include_path` / `set_include_path` / `ini_restore('include_path')`の部分集合 |

foreachは配列の値渡し形式が対象です。参照形式・オブジェクト反復・分解代入は明示的に拒否します。
参照全般・型付きプロパティ／定数の未収録形式・上表以外の標準関数と定数は対象外です。クロージャは値捕捉・arrow・static形式の部分集合です。参照捕捉、Traversableの引数展開、一般のcallable、例外APIとエラーハンドラの残る形式は未対応です。動的インスタンスプロパティ名の式は部分対応で、マジックメソッドや未収録の名前型変換・評価順序は残ります。クラス機能と複数ファイル実行にも未対応の境界条件があります。
キャストの別名やその他のキャスト、文字列添字の更新にも対応していません。
通常関数とメソッドを合わせた呼び出し深度は256です。
詳細な仕様・診断・版差・制約は [互換性文書](docs/compatibility.md) に記載しています。

依存調査用の[専用アプリ](research/laravel-api/README.md)には、固定したComposer依存と最小HTTP検証コードを置いています。公式PHPの3版でアプリ検証に成功していますが、既存fixtureの一致もこの公式PHPでの成功も、Go製PHPのフレームワーク起動を示すものではありません。

## 検証

```sh
make test vet build
make compare       # 主基準 PHP 8.4.26、全fixture
make compare-all   # 8.3.35 / 8.4.26 / 8.5.11、全fixture
make compare-parser
```

比較環境の起動方法は [コンテナの案内](containers/php-version/README.md) を参照してください。
[php-go-interpreterスキル](skills/php-go-interpreter/SKILL.md) の比較スクリプトを使用し、
stdout・stderrをバイト単位、終了コードを厳密に比較します。trim・警告の削除・fixtureの除外は行いません。
正常系がどちらか一方でも非ゼロ終了なら検証エラーとします。

以下は追加前の第1版に関する検証記録です：

- Goテスト、Pythonテスト11件、go vet、ビルドは成功。Pythonにはコンテナ停止の検出と誤検出防止の回帰テスト2件を追加しています。
- 既存337件のPHPソースと保存済み期待値を維持し、129件追加。全466件（正常208・実行時198・構文／コンパイル60）。
- ビルド済みGo CLIと保存済み期待値の照合は3設定とも466一致・差分0・検証エラー0。正常系208件は全設定で終了コード0。
- ユーザー環境での公式PHP実比較は8.3.35・8.4.26・8.5.11の各466一致・差分0・検証エラー0。正常系208件は各版とも両処理系で終了コード0。
- fixtureの補助関数名の修正後、停止していた3コンテナを起動し、`check.sh` と `make compare-all` が成功しました。

[修正後のローカル検証ログ](docs/validation/v1-followup-local-validation.txt)、
[ユーザー提供の全件実比較と完成確認](docs/validation/v1-user-comparison-complete.md) を区別して記録しています。
新規129件の期待値は公式資料・PHPソースを根拠に作成し、全件を3版の公式PHP実比較で確認しました。
比較スクリプトはPodman接続障害・コンテナ停止を検証エラーとして扱い、その版の残りを未実行と明示します。
[全件成功前の修正記録](docs/validation/v1-user-comparison-followup.md) と
[エージェント環境で接続を拒否された試行ログ](docs/validation/v1-followup-php-comparison.txt) は履歴として残しています。

追加前の337件は、ユーザー提供の公式比較ログで各版337一致・差分0・検証エラー0を確認済みです。
これは今回の変更前の履歴です。[以前の仕様・検証記録](docs/validation/pre-v1-compatibility.md) を参照してください。
PHP全機能との互換性を示すものではありません。

リリース前の再構築手順・CI・公開状況は[v0.3.0の再構築・CI・公開記録](docs/task-api-release.md)を参照してください。
