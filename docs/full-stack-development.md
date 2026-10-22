# 固定検証アプリへ向けた6段階の連続実装

2026-10-22追記：既存の固定アプリ完成記録を維持し、別対象のLaravel 13.38.0（開発依存込み110パッケージ）をGo8.5で実起動した。autoload、CLI、標準migration、実HTTP、Blade冷/温キャッシュ、静的アセット、非公開経路、20連続/4並列要求を確認。現行1201 fixtureは公式PHP3版と全件一致。今回の共通実装修正後に、固定Laravel/API PlatformはGo3設定で96連続/8並列、タスクAPIは公式/Go3組で各191要求、Laravel12は8.5の公式/Go各111要求を再検証した。[変更範囲・実行証拠・部分対応の制限](validation/laravel13-7oa6dbgc/REPORT.md)。以下の2000要求や旧件数は当時の履歴であり、今回の実測数とは区別する。

2026-10-21現在。主基準はPHP 8.4.26、比較対象は8.3.35・8.5.11。完成判定は[固定アプリのHTTP・DB契約](laravel-api-target.md)と、Go製PHPによる`vendor/autoload.php`、Laravel CLI・HTTP、API Platformの実行結果に従う。固定した`research/laravel-api/composer.lock`と`vendor/`は変更しない。公式PHPでのアプリ成功はGoでの実行成功ではない。

追加の実用アプリ検証として[research/task-api](../research/task-api/README.md)を作成した。Laravel／Eloquent／SQLiteによるTask CRUDの実HTTP比較は公式PHP・Goの3組、各191要求で成功。新規4 fixtureを含む現行**1173件（正常745・実行時294・構文134）**も公式3版すべて一致し、既存BookアプリのGo3設定実TCP回帰も成功した。[タスクAPIの実装範囲・検証結果・制限・ログ](validation/task-api-results.md)。以下の6段階の完成記録と既存成果は保持する。


タスクAPI追加前のfixture・期待値は**1169件**（正常741・実行時294・構文134）。公式PHP [8.3.35](validation/full-stack-1169-oracle-nulls-83-mac.log)・[8.4.26](validation/full-stack-1169-oracle-nulls-84-mac.log)・[8.5.11](validation/full-stack-1169-oracle-nulls-85-mac.log)との版別`make compare`は、各1169件の個別MATCH、差異0・エラー0、終了0を記録した。`check.sh`の3版PASSはユーザー報告であり、`make compare-all`の成功記録ではない。固定アプリはGo製PHPの実TCPサーバーで[8.3](validation/full-stack-real-http-83-mac.log)・[8.4](validation/full-stack-real-http-84-mac.log)・[8.5](validation/full-stack-real-http-85-mac.log)の各設定を検証し、CLI migration、HTTP・SQLite CRUD、入力検証、API Platform IRI／OpenAPI、2000回連続・8並列要求が各版終了0だった。**固定アプリについて定義した全6段階の完成条件は達成**した。PHP全機能の対応完了を示すものではない。

| 段階 | 依存順の実装単位 | 現在の状態 | 統合確認 |
| --- | --- | --- | --- |
| 1：参照・反復 | 参照引数・戻り値・捕捉 → 参照foreach → Iterator／IteratorAggregate／Traversable → Countable／ArrayAccess／Generator | 基本の`yield`・`yield from`と遅延実行、ArrayAccess、FilterIteratorまで追加。参照引数の各種格納先、参照戻り値、Generatorの全制御、SPLの再帰ディレクトリ反復は基本形を実装したが、全制御は未完了 | GeneratorとFilterIteratorの小fixtureはGoで実行済み。収録fixtureは公式3版比較済み |
| 2：オブジェクト操作・DI | 動的呼び出し・マジックメソッド・clone → Reflection・属性 | clone、Reflectionのクラス・メソッド・引数の一部、属性構文の保持を追加。ReflectionParameterの属性列挙と生成を接続。完全なDIは未対応 | Container経路でReflectionParameterを使用。自作属性からサービスを生成する小fixtureはGoで成功。固定vendorのContainer属性DIとAPI Platform属性リソース発見をGoの3設定で実行済み。Reflection一般は未完了 |
| 3：ライブラリ読み込み | 固定vendorの経路に必要な標準関数・ファイル操作 → autoload登録 → PSR-4／classmap／files | autoloadと標準関数を実アプリの経路に沿って追加。PCRE・iconv・日時APIは一部のみ対応 | 実`vendor/autoload.php`と`artisan list`はGoで終了0。収録fixtureは公式3版比較済み、実vendorの全経路は比較対象外 |
| 4：Laravel CLI | 設定・コンテナ・Console起動時の不足を補完 | Container、Symfony Console、設定読込の経路を順に実行して不足を修正中 | `artisan list`・`about`・`route:list`が実際の一覧や状態を出力し終了0。別プロセスのGo製PHPで3設定それぞれのmigration・`gophp:probe`・`api:openapi:export`・DBテーブル作成を確認 |
| 5：HTTP・DB | HTTP状態分離 → PDO／SQLite → バリデーション・Eloquent | 固定アプリの経路を実装。PDO/SQLite一般の全操作は未対応 | 一時DBでJSON応答、作成・取得・無効入力422、状態分離をプロセス内と実TCPサーバーの両方で確認。各版の実TCPで2000回連続・8並列要求が成功。収録fixtureも公式3版比較済み |
| 6：API Platform | リソース・属性・シリアライズ → CRUD・API仕様 → 継続実行 | 固定BookリソースのPOST/GET/PATCH/DELETE、IRI、OpenAPIを実装。API Platform一般の全機能は対象外 | GoのHTTP統合テストと8回の追加CRUD、JSON-LDページ分割を確認。各版の実TCPでも固定リソースCRUD・OpenAPI・DB状態が成功。収録fixtureは版別公式比較済み |

## 固定アプリの完成判定（2026-10-21）

この判定は固定したアプリのCLI・HTTP・DB・API契約を対象とし、PHP全機能の実装完了を意味しない。現行1169件への変更後に、次を確認した。

| 要件 | 現在の証拠 | 判定 |
| --- | --- | --- |
| 版別CLI・DB・OpenAPI | [別プロセスのGo3設定検証](validation/full-stack-current-cli-matrix-local.log)。各版のmigration、`gophp:probe`、`api:openapi:export`、SQLiteテーブル確認がPASS、スクリプト終了0 | ローカル確認済み |
| HTTP・DB・API | [固定アプリの対象統合テスト](validation/full-stack-current-app-gate-local.log)が終了0。API Platform CRUDは8.3・8.5、手動routeの422・OpenAPI・状態分離は8.4、JSON-LDページ分割は3設定で確認 | プロセス内HTTPハンドラで確認済み |
| 実サーバー経由のHTTP契約 | [8.3](validation/full-stack-real-http-83-mac.log)・[8.4](validation/full-stack-real-http-84-mac.log)・[8.5](validation/full-stack-real-http-85-mac.log)の版別生ログで、JSON、クエリ、422、SQLite CRUD、IRI、OpenAPI、2000回連続・8並列要求が各版PASS・終了0 | 固定アプリの契約を実TCPで確認済み |
| 現行fixtureの公式PHP比較 | [8.3](validation/full-stack-1169-oracle-nulls-83-mac.log)・[8.4](validation/full-stack-1169-oracle-nulls-84-mac.log)・[8.5](validation/full-stack-1169-oracle-nulls-85-mac.log)は各1169件の個別MATCH、0差異・0エラー、終了0 | 収録fixtureの範囲で達成 |

実TCPの版別生ログには、Go製PHPの`--serve`起動、`/api/probe`とクエリ付きGET、手動routeの作成201・取得200・入力エラー422、API Platformの作成201・返却IRIの取得200・更新200・削除204・削除後404、OpenAPIの200、最終SQLite行の検査後のPASSがある。各版で2000回の連続GET（probe 1750回・保存済みBook取得250回）と8並列GETを確認し、`exit_code: 0`で終わった。[ログの照合記録](validation/full-stack-completion-evidence-review.log)に個別結果と件数の検査を残す。`server_exit_code: -15`は検証後の`terminate()`による停止で、検証コマンドの失敗ではない。サーバー出力は各版ログの`server_log_begin`／`server_log_end`に収録され、今回は空だった。

検証スクリプトはmigration成功後、サーバー起動前に**検証専用の一時storage**内の`framework/cache/data`だけを削除・再作成する。migration前にAPI Platformが空のBookメタデータをキャッシュするための準備で、公式PHP用`verify.py`と同じ順序である。初回Mac試行は失敗したが、そのログは再実行で上書きされた。初回の症状とこのキャッシュ削除による解消はユーザー報告として記録し、初回の出力・終了コードを再構成しない。以前のCodex制限環境でのTCP bind拒否は別の[保存済み試行](validation/full-stack-1169-real-http-84-local.log)に残す。

以上により、**固定したLaravel・API Platform検証アプリの完成条件は達成**。対象外のPHP機能（Reflection・Generator・SPL・PDOなどの全境界）、実運用条件での長時間負荷や複数DBは完了と扱わない。

初回比較より前の接続試行（履歴）：[接続ログ](validation/full-stack-gate-1169-connectivity-2026-10-21.log)では`podman ps`とPHP 8.3ラッパーの`-v`がともに接続拒否で終了125となり、その時点では版別`make compare`を開始できなかった。当時のfixture・期待値は各1169件（正常745・実行時290・構文134）だった。

### 1169件の初回公式比較と修正

ユーザーがMacで実行した版別[8.3生ログ](validation/full-stack-1169-83-mac.log)・[8.4生ログ](validation/full-stack-1169-84-mac.log)・[8.5生ログ](validation/full-stack-1169-85-mac.log)は、順に**1159一致・7差異・3エラー**、**1160一致・7差異・2エラー**、**1160一致・7差異・2エラー**で、いずれも`exit_code: 2`。`check.sh`の3版PASSはユーザー報告。エラーは環境障害ではなく、`success`に置かれた公式PHPの致命終了ケースである。

Go側では、`yield from`が委譲先の数値キーを外側の次の自動キーへ持ち込まないようにし、`IteratorAggregate`経由で取得したGeneratorへの`send`と戻り値は直接委譲と区別した。再includeの無名クラスには別のクラス同一性を与えた。戻り値型のない`IteratorAggregate::getIterator()`の非推奨診断、PDOの未実装属性の`errorInfo`で欠ける第3要素、`call_user_func`系のstrict型、`LimitIterator`のlimit 0でのrewind例外とフレームを修正した。PDOの`ATTR_ORACLE_NULLS`は[PHP 8.4.26のPDO実装](https://raw.githubusercontent.com/php/php-src/php-8.4.26/ext/pdo/pdo_dbh.c)が値を検査せず、[2ビットの格納先](https://raw.githubusercontent.com/php/php-src/php-8.4.26/ext/pdo/php_pdo_driver.h)へ代入するため、9の読み戻しは1となる。`array_any()`は[PHP 8.4以降](https://www.php.net/manual/en/function.array-any.php)の関数なので、共通fixtureは8.3でも有効な分岐を使う。

公式結果に従い、非推奨診断を出す`generator_reference_aggregate`と`pdo_fetch_attributes`、致命終了する`internal_callback_weak_strict`と`limit_iterator_basic`を`runtime/`へ移した。4件は名前と検証内容を維持し、期待値のstdout・stderr・終了コードを初回生ログから訂正した。`http_core_primitives`の8.3用分岐と合わせて、既存の暫定期待値だけを正解として扱わない。fixture総数は1169件のまま、現在の内訳は正常741・実行時294・構文134。[make test](validation/full-stack-1169-fix-test-local.log)、[vet](validation/full-stack-1169-fix-vet-local.log)、[build](validation/full-stack-1169-fix-build-local.log)、差分検査はローカルで終了0。[Podman再確認](validation/full-stack-1169-fix-podman-local.log)は接続拒否で終了125のため、修正後の公式3版比較は各1169件とも未実施で、期待値は再比較まで暫定のまま。

その後の版別再比較（2026-10-21）：[8.3生ログ](validation/full-stack-1169-after-fix-83-mac.log)と[8.4生ログ](validation/full-stack-1169-after-fix-84-mac.log)は各**1169一致・0差異・0エラー、終了0**。[8.5生ログ](validation/full-stack-1169-after-fix-85-mac.log)は**1168一致・1差異・0エラー、終了2**。`success/pdo_oracle_nulls_version.php`のみ、公式8.5の`accepted|9`に対してGoは`accepted|1`だった。8.3・8.4の公式結果は`accepted|1`。`check.sh`の3版PASSはユーザー報告。GoのPDO属性格納を8.5では設定整数の保持、8.3・8.4では2ビット格納に分け、同じfixtureの期待値も版別にした。[make test](validation/full-stack-1169-pdo85-test-local.log)・[vet](validation/full-stack-1169-pdo85-vet-local.log)・[build](validation/full-stack-1169-pdo85-build-local.log)・[差分検査](validation/full-stack-1169-pdo85-diff-local.log)は終了0。[Go3設定の対象ログ](validation/full-stack-1169-pdo85-target-local.log)は順に`accepted|1`・`accepted|1`・`accepted|9`。この環境の[Podman接続確認](validation/full-stack-1169-pdo85-podman-local.log)は終了125のため、8.5の当該1件と変更後の全件一致は公式再比較待ち。初回の差異と接続拒否は上記の履歴として維持する。

初回差異後に用意した版別再比較手順（履歴）。公式比較は後述の1169件一致で達成済み。実HTTP検証は版ごとに独立した一時SQLiteを使う。

当時のPDO版差だけを8.5で再確認する対象コマンドは、`make build`後の`python3 skills/php-go-interpreter/scripts/compare_php.py --php ./containers/php-version/bin/php85 --php-version 8.5.11 --fixtures tests/fixtures/success/pdo_oracle_nulls_version.php --interpreter ./bin/gophp --php-version=8.5 '{file}'`だった。全1169件は下記ループで版別に比較し、後述の成功ログを保存した。

```zsh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  read -r tag version target <<< "$spec"
  log="docs/validation/full-stack-gate-1169-compare-${tag}-mac.log"
  printf 'command: make compare PHP=./containers/php-version/bin/php%s PHP_VERSION=%s\n' "$tag" "$version" > "$log"
  make compare "PHP=./containers/php-version/bin/php${tag}" "PHP_VERSION=${version}" 2>&1 | tee -a "$log"
  result=$?; printf 'exit_code: %s\n' "$result" >> "$log"
  app_log="docs/validation/full-stack-gate-app-${tag}-mac.log"
  printf 'command: python3 research/laravel-api/verify-gophp.py --version %s\n' "$target" > "$app_log"
  python3 research/laravel-api/verify-gophp.py --version "$target" 2>&1 | tee -a "$app_log"
  result=$?; printf 'exit_code: %s\n' "$result" >> "$app_log"
done
```

その後の公式比較結果（2026-10-21）：[8.3生ログ](validation/full-stack-1169-oracle-nulls-83-mac.log)・[8.4生ログ](validation/full-stack-1169-oracle-nulls-84-mac.log)・[8.5生ログ](validation/full-stack-1169-oracle-nulls-85-mac.log)は、各1169件の個別MATCH・1169一致・0差異・0エラー・`exit_code: 0`を記録する。版ごとの`make compare`であり、`make compare-all`成功の証拠ではない。`check.sh`の3版PASSはユーザー報告。前段落までの不一致・接続拒否・期待値暫定は当時の履歴で、現行1169件の公式比較待ちは解消した。

制限環境での以前の試行（履歴）：`python3 research/laravel-api/verify-gophp.py --version 8.4`は[実行ログ](validation/full-stack-1169-real-http-84-local.log)のとおり、一時SQLiteへのmigrationを終了0で実行後、空きTCPポートのbindで`Operation not permitted`となり、スクリプト終了2。当時はGoサーバーが起動せず、HTTP要求・応答とサーバーログは未取得だった。[検証スクリプトのヘルプ](validation/full-stack-1169-http-verifier-check-local.log)、[CLI単独の8.4再確認](validation/full-stack-1169-http-cli-check-local.log)、[Python12件](validation/full-stack-1169-http-python-local.log)、[差分検査](validation/full-stack-1169-real-http-diff-local.log)は終了0。Go実装は変更していない。[既存の2000回・8並列のin-process成功ログ](validation/full-stack-http-2000-stress-local.log)は終了0だが、重複実行の[中断ログ](validation/full-stack-1169-http-stress-2000-local.log)は終了130として残す。これらの当時の記録は、後続のMac実TCP成功ログに置き換えない。

当時は公式比較だけが達成済みで、実TCPの完成判定が残っていた。後続の版別Macログで実TCP条件も満たした。対象外のPHP全機能実装を、この固定アプリ判定へ追加しない。Podmanの自動再起動・VM再作成は行わない。

実TCP検証の再現手順は次のとおり。各版で一時SQLite・storageを作成し、Go製PHPでmigration、サーバー起動、待機、HTTP契約、2000回のGETと8並列GET、DB確認を行う。スクリプトは自分が起動したサーバーを停止し、一時ディレクトリを削除する。版別ログにはHTTP結果、サーバー出力、サーバー終了コード、検証コマンドの終了コードを残す。

```zsh
cd /Users/kouji/dev/php
set -o pipefail
make build 2>&1 | tee docs/validation/full-stack-real-http-build-mac.log
build_rc=$?
printf 'exit_code: %s\n' "$build_rc" | tee -a docs/validation/full-stack-real-http-build-mac.log
if (( build_rc != 0 )); then exit "$build_rc"; fi
for spec in '83 8.3' '84 8.4' '85 8.5'; do
  read -r tag version <<< "$spec"
  log="docs/validation/full-stack-real-http-${tag}-mac.log"
  printf 'command: python3 research/laravel-api/verify-gophp.py --version %s --stress-requests 2000\n' "$version" > "$log"
  python3 research/laravel-api/verify-gophp.py --version "$version" --stress-requests 2000 2>&1 | tee -a "$log"
  verify_rc=$?
  printf 'exit_code: %s\n' "$verify_rc" | tee -a "$log"
done
```

以下の機能別記録は作業当時の履歴である。各節に残る「現行件数」「公式比較待ち」はその時点を指し、冒頭の1169件・3版一致の現状を上書きしない。

## 2026-10-21：IteratorAggregateの有限な多段委譲

`yield from`の委譲先が`IteratorAggregate::getIterator()`を連続して返すとき、Go側の深さ16制限が有限の24段連鎖も拒否していた。[再現fixture](../tests/fixtures/success/generator_delegate_aggregate_depth.php)を加え、安全上限を1024段へ広げた。Goの3設定では葉の`ArrayIterator`まで到達して`9`、終了0を[対象ログ](validation/full-stack-generator-aggregate-depth-target-local.log)に記録した。[IteratorAggregateの公式契約](https://www.php.net/manual/en/class.iteratoraggregate.php)は`getIterator(): Traversable`で、16段の制限を規定しない。ただし公式PHP3版の実比較は未実施で、この期待値は暫定。自己を返す循環経路はGo側で上限により終了255となり、公式PHPの診断との一致は未確認である。

現行fixture・期待値は**1169件**（正常745・実行時290・構文134）。直近の公式比較済み997件から追加172件と変更後全件は未比較。[Podmanの現時点の接続確認](validation/full-stack-current-podman-check.log)は終了125で`operation not permitted`であり、fixtureの一致・不一致には数えない。今回の各版は一致0・不一致0・全1169件未実施で、版共通のPodman接続確認1回のみ環境エラー（終了125）だった。[全ローカルtest](validation/full-stack-generator-aggregate-depth-test.log)（Python12件）・[vet](validation/full-stack-generator-aggregate-depth-vet.log)・[build](validation/full-stack-generator-aggregate-depth-build.log)・[差分検査](validation/full-stack-generator-aggregate-depth-diff.log)とGo対象ログは終了0。Macでは次節の版別対象比較コマンドのfixtureを`generator_delegate_aggregate_depth.php`に替え、続いて版別の`make compare`で全件比較する。

## 2026-10-21：固定Laravel LazyCollectionのGenerator統合

固定vendorの`Illuminate\Support\LazyCollection`を実際に読み込み、`flatten()`・`flatMap()`・`chunk()`、一部だけ消費した後の`remember()`再走査を[統合テスト](../cmd/gophp/main_test.go)へ追加した。配列とGeneratorへの`yield from`、クロージャの参照捕捉、複数回の反復を組み合わせ、Goの8.3・8.4・8.5設定で[対象テスト](validation/full-stack-lazy-collection-test-local.log)が終了0になった。同じ3設定で[`artisan list`の版別実行](validation/full-stack-lazy-artisan-list-83-local.log)（[8.4](validation/full-stack-lazy-artisan-list-84-local.log)、[8.5](validation/full-stack-lazy-artisan-list-85-local.log)）もコマンド一覧を出力して終了0となった。[全ローカルtest](validation/full-stack-lazy-collection-all-test-local.log)（Python12件）・[vet](validation/full-stack-lazy-vet-local.log)・[build](validation/full-stack-lazy-build-local.log)・[差分検査](validation/full-stack-lazy-diff-local.log)は終了0。これはGoのローカル統合証拠であり、公式PHP3版での一致やGenerator全制御の実装完了を示さない。現行1168件の公式比較と固定アプリの実HTTPサーバー検証は引き続き未実施である。

## 2026-10-21：深い有限配列のarray_walk_recursive

固定vendorのLaravel TranslatorとSymfony Translationは`array_walk_recursive()`の参照コールバックを使う。Go実装の再帰処理は深さ64を超える有限配列も循環として拒否していた。再帰中の配列の同一性を追跡し、固定の早い拒否を除いた。[70段の回帰fixture](../tests/fixtures/success/array_walk_recursive_depth.php)は葉の値を1回だけ参照更新し、Goの8.3・8.4・8.5設定で`8|1`と終了0を記録した。[PHP 8.3.35](https://raw.githubusercontent.com/php/php-src/php-8.3.35/ext/standard/array.c)・[8.4.26](https://raw.githubusercontent.com/php/php-src/php-8.4.26/ext/standard/array.c)・[8.5.11](https://raw.githubusercontent.com/php/php-src/php-8.5.11/ext/standard/array.c)の公式実装は再帰中の配列を検出するが、深さ64の固定上限は設けない。この結果は公式PHPでの実比較ではなく、期待値は暫定である。循環配列は既存の参照基盤がまだ拒否し、1024段超はGoの安全上限で拒否する。固定vendorのLaravel Translatorのネストした翻訳文も[再現可能なGo統合テスト](../cmd/gophp/main_test.go)と[対象ログ](validation/full-stack-laravel-translation-test.log)で3設定とも`Hi Ada`と`Bye Ada`への置換を確認したが、公式PHPとの一致とは扱わない。

現行**1168件**（正常744・実行時290・構文134）から追加171件と変更後全件の公式PHP3版比較は未実施。既知の[Podman接続拒否](validation/full-stack-continuation-podman.log)後に環境変化がないため再試行していない。今回の各版の一致0・不一致0・環境エラー0・未実施1168件。[Go対象ログ](validation/full-stack-array-walk-depth-target-local.log)、[全ローカルtest](validation/full-stack-array-walk-depth-test.log)（Python12件）・[vet](validation/full-stack-array-walk-depth-vet.log)・[build](validation/full-stack-array-walk-depth-build.log)・[差分検査](validation/full-stack-array-walk-depth-diff.log)はいずれも終了0。[固定Laravel Translatorの統合テスト](validation/full-stack-laravel-translation-test.log)と[`artisan list`再実行](validation/full-stack-array-walk-depth-artisan-list-local.log)も成功した。Macでは下節の版別対象比較で`array_walk_recursive_depth.php`を指定し、その後に版別の`make compare`で全件比較する。独立して進める次の作業は、固定vendorが使う`array_walk()`の通常配列・オブジェクト境界とGeneratorの未確認制御を調べ、固定アプリ経路で必要な差異を修正することである。

## 2026-10-21：extractの入力変数と同名キー

`extract($data, EXTR_REFS)`で入力配列に`data`キーがあると、最初の要素の参照取り込みで`$data`の束縛が付け替わる。従来のGo実装は次の要素で付け替え後の整数を配列として使い、Fatalになっていた。元の配列を保持する参照セルを走査中に使うよう修正した。[回帰fixture](../tests/fixtures/success/extract_refs_shadow.php)は同名キーと後続キー、元配列に残る別名、取り込み後の双方向更新を確認する。[Go3設定の対象ログ](validation/full-stack-extract-shadow-target-local.log)は出力`2|10|20|10|11|21`と終了0を記録する。期待値は[PHP 8.3.35](https://raw.githubusercontent.com/php/php-src/php-8.3.35/ext/standard/array.c)・[8.4.26](https://raw.githubusercontent.com/php/php-src/php-8.4.26/ext/standard/array.c)・[8.5.11](https://raw.githubusercontent.com/php/php-src/php-8.5.11/ext/standard/array.c)の`php_extract_ref_overwrite()`が元の配列を走査し、シンボル表の変数束縛を要素参照へ付け替える実装から判断した暫定値で、公式PHPの実比較結果ではない。

現行fixture・期待値は**1167件**（正常743・実行時290・構文134）。追加170件と変更後全件の公式PHP3版比較は未実施。各版の今回の一致0・不一致0・環境エラー0・未実施1167件。最新の[全ローカルtest](validation/full-stack-extract-shadow-test.log)・[vet](validation/full-stack-extract-shadow-vet.log)・[build](validation/full-stack-extract-shadow-build.log)・[差分検査](validation/full-stack-extract-shadow-diff.log)は別ログに記録する。Macでは下節の対象比較コマンドのfixtureを`tests/fixtures/success/extract_refs_shadow.php`へ替え、続けて版別`make compare`で全件比較する。動的変数名による同様の束縛変更や、引数式の途中で元配列が変更される場合は未検証。

## 2026-10-21：array_walkの参照格納先と引数評価

`array_walk()`／`array_walk_recursive()`でも、配列引数の受け手・添字を呼び出し時に1回評価し、各要素の参照コールバックへ同じ格納先から作った子要素の参照を渡す。`extract()`で導入した評価済み引数の格納先を、参照引数の束縛にも接続した。[回帰fixture](../tests/fixtures/success/array_walk_source_once.php)は、受け手と添字を返す関数の呼び出し回数、元配列への書き戻し、入れ子配列の再帰版を確認する。[Go3設定の対象ログ](validation/full-stack-array-walk-place-target-local.log)は出力と終了0を記録した。最新の[全ローカルtest](validation/full-stack-array-walk-place-test.log)、[vet](validation/full-stack-array-walk-place-vet.log)、[build](validation/full-stack-array-walk-place-build.log)、[差分検査](validation/full-stack-array-walk-place-diff.log)を検証記録とする。

当時のfixture・期待値は**1166件**（正常742・実行時290・構文134）。追加169件と変更後全件の公式PHP3版比較は未実施で、期待値は暫定。各版の今回の一致0・不一致0・環境エラー0・未実施1166件。コールバック中に配列構造を変更する場合や通常オブジェクトを`array_walk()`へ渡す場合は未検証・未対応。Macでは下節の版別比較コマンドの対象を`tests/fixtures/success/array_walk_source_once.php`へ替えて対象比較し、その後に版別の`make compare`で全件を比較する。

## 2026-10-21：extractの参照取り込み

固定vendorの`laravel/serializable-closure/src/Serializers/Native.php`は`extract($this->code['use'], EXTR_OVERWRITE | EXTR_REFS)`を使用する。`extract()`の`EXTR_REFS`と`EXTR_SKIP | EXTR_REFS`を、参照可能な配列引数について既存の参照セルへ接続した。名前付き引数は仮引数名で元の式を選び、逆順指定にも対応する。[基本fixture](../tests/fixtures/success/extract_refs_basic.php)では変数間の双方向更新、コピーとの分離、ローカルスコープ、オブジェクト内の配列、`EXTR_SKIP`を確認した。[参照を返す関数のfixture](../tests/fixtures/success/extract_refs_source_once.php)は評価済み参照セルを使い、[受け手と添字に副作用があるfixture](../tests/fixtures/success/extract_refs_lvalue_once.php)は引数評価時に確定した格納先を使う。後者では受け手と添字を返す関数を各1回だけ呼び、参照から元の配列を更新する。Go3設定の[基本ログ](validation/full-stack-extract-refs-target-local.log)・[参照元関数ログ](validation/full-stack-extract-source-once-local.log)・[格納先評価ログ](validation/full-stack-extract-lvalue-once-local.log)で対象出力と終了0を確認した。最新の[全ローカルtest](validation/full-stack-extract-lvalue-once-test.log)、[vet](validation/full-stack-extract-lvalue-once-vet.log)、[build](validation/full-stack-extract-lvalue-once-build.log)、[差分検査](validation/full-stack-extract-lvalue-once-diff.log)を検証記録とする。

当時は**1165件**（正常741・実行時290・構文134）。[公式PHPの`extract()`仕様](https://www.php.net/function.extract.php)を参照したが、今回の新規fixtureと変更後全件の公式PHP3版比較は未実施で、期待値は暫定。各版の今回の一致0・不一致0・環境エラー0・未実施1165件。受け手・キー式を1回評価する経路はローカル確認済みだが、抽出中に元の配列変数と同名のキーでその束縛が変わる場合や、参照元の引数展開は未検証。`EXTR_PREFIX_*`等の他フラグも未対応。固定vendorでの使用箇所は確認したが、そのシリアライザの実行経路は未検証。

Macでは次節の版別比較コマンドを3回使い、対象fixtureを`extract_refs_basic.php`・`extract_refs_source_once.php`・`extract_refs_lvalue_once.php`へそれぞれ替えて対象比較する。続いて`make compare`で全件比較し、各コマンド直後の終了コードをログへ保存する。

## 2026-10-21：array_walkの名前付き引数と参照コールバック

`array_walk()`／`array_walk_recursive()`は参照コールバックを既存の配列要素の参照格納先へ渡していたが、`callback:`を先に指定すると最初の呼び出し式を配列の参照元と誤認していた。評価済み引数の`array`仮引数に対応する元の式を選ぶよう修正した。[回帰fixture](../tests/fixtures/success/array_walk_named_reference.php)は、逆順の名前付き引数、参照での要素更新、再帰版を確認する。[Go3設定の対象ログ](validation/full-stack-array-walk-named-target-local.log)と[全ローカルtest](validation/full-stack-array-walk-named-test.log)（Python12件）・[vet](validation/full-stack-array-walk-named-vet.log)・[build](validation/full-stack-array-walk-named-build.log)・[差分検査](validation/full-stack-array-walk-named-diff.log)は終了0。現行**1162件**（正常738・実行時290・構文134）で、この期待値は[PHPのarray_walkの引数・参照コールバックの説明](https://www.php.net/manual/en/function.array-walk.php)とGoローカル結果に基づく暫定値。公式PHP3版の実比較は環境変化のないPodman接続拒否のため未試行で、各版の今回の一致0・不一致0・環境エラー0・未実施1162件。引数の式に副作用がある場合の参照元再評価、配列構造をコールバック中に変更する場合のPHP未定義動作は、このfixtureで検証していない。Macでは下節の版別比較コマンドでこのfixtureを対象比較し、続けて全件比較する。

Macでの対象・全件比較はリポジトリのルートで次を実行する。各コマンドの直後に終了コードを保存する。

```zsh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  read -r tag version target <<< "$spec"
target_log="docs/validation/full-stack-array-walk-named-target-${tag}-mac.log"
printf 'command: compare_php.py php%s %s array_walk_named_reference.php --php-version=%s\n' "$tag" "$version" "$target" > "$target_log"
python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php${tag}" --php-version "$version" --fixtures tests/fixtures/success/array_walk_named_reference.php --interpreter ./bin/gophp "--php-version=${target}" '{file}' 2>&1 | tee -a "$target_log"
compare_rc=$?; printf 'exit_code: %s\n' "$compare_rc" >> "$target_log"
all_log="docs/validation/full-stack-array-walk-named-all-${tag}-mac.log"
printf 'command: make compare PHP=./containers/php-version/bin/php%s PHP_VERSION=%s\n' "$tag" "$version" > "$all_log"
make compare "PHP=./containers/php-version/bin/php${tag}" "PHP_VERSION=${version}" 2>&1 | tee -a "$all_log"
compare_rc=$?; printf 'exit_code: %s\n' "$compare_rc" >> "$all_log"
done
```

## 2026-10-21：固定アプリの応答契約を追加検査

[HTTP統合テスト](../cmd/gophp/main_test.go)で、API Platformの作成・取得応答の`application/ld+json`と、OpenAPI JSON内の`/api/books`のGET・POST定義を確認するようにした。実サーバー用の[版別検証スクリプト](../research/laravel-api/verify-gophp.py)にも同じ確認を加えた。対象Goテストは単独で通過し、[全ローカルtest](validation/full-stack-http-contract-strengthened-test.log)（Python12件）・[vet](validation/full-stack-http-contract-strengthened-vet.log)・[build](validation/full-stack-http-contract-strengthened-build.log)・[Python構文検査](validation/full-stack-http-contract-strengthened-python.log)・[差分検査](validation/full-stack-http-contract-strengthened-diff.log)は終了0。fixture数は**1161件**のまま。この追加検査は既存のGo出力を確認したもので、公式PHPのOpenAPI出力との同一性や実HTTPサーバーの成功を証明しない。公式PHP3版の現行1161件比較は既知のPodman接続拒否後の環境変化がなく未試行で、各版の今回の一致0・不一致0・環境エラー0・未実施1161件。次はMacで実HTTP検証と版別の対象・全件比較を行い、API仕様の経路と診断を照合する。

## 2026-10-21：実サーバーでの固定アプリ検証準備

[Go実サーバー検証スクリプト](../research/laravel-api/verify-gophp.py)を加えた。版別に一時SQLiteとstorageを作り、Go製PHPで`artisan migrate --force`を実行してから`--serve`で固定アプリを起動する。続いて公式PHP側の[既存検証スクリプト](../research/laravel-api/verify.py)と同じ初期HTTP契約を観測し、API PlatformのPATCH・DELETE・削除後404、OpenAPI JSON、最終DB行まで追加確認する。両スクリプトで使うHTTP要求関数は、PATCH用にContent-Typeを明示できるよう拡張した。従来の公式PHP検証呼び出しの引数と既定動作は変えていない。

Pythonの[構文・ヘルプ検査](validation/full-stack-go-server-verifier-local.log)は終了0。[8.4設定の実行試行](validation/full-stack-go-server-84-sandbox-attempt.log)では、Goの別プロセスで`artisan migrate --force`が一時SQLiteにmigrationを作成して終了0だった。その直後、127.0.0.1への[待ち受けが権限エラー](validation/full-stack-real-server-sandbox-check.log)となり、スクリプトは環境エラーとして終了2を記録した。実HTTPサーバーのリクエストは**未実施**で、in-processのGo HTTP統合テスト成功を読み替えない。Podmanの接続拒否も変わらず、現行1161件の公式PHP3版比較は未実施。今回の各版の公式fixture一致0・不一致0・環境エラー0・未実施1161件。Macではビルド後に次を実行し、実際の終了コードを保存する。

待ち受けを使わない`--cli-only`でも同じ版別一時環境を検査できるようにした。[別プロセスの3設定生ログ](validation/full-stack-go-cli-subprocess-3versions-local.log)は、各版で`artisan migrate --force`と`artisan gophp:probe`が終了0、`probe-ok`を出力し、SQLiteに`books`テーブルができたことを記録する。スクリプト全体も終了0。これはGo製PHPのCLI統合結果であり、公式PHPとのfixture厳密比較や実HTTPの成功ではない。

CLI単独モードへ`api:openapi:export`のJSON解析も加えた。[再実行した3設定生ログ](validation/full-stack-go-cli-openapi-subprocess-3versions-local.log)には各コマンドの終了0と生成JSON、`openapi`版・`/api/books`のGET／POST定義の検査結果、スクリプトの終了0を保存した。以前のCLI単独ログは履歴として維持する。これも実HTTPサーバーや公式PHPの出力一致を示すものではない。

CLIのみの再実行は`python3 research/laravel-api/verify-gophp.py --cli-only`で行う。`--version 8.4`のように指定するとその設定だけを実行する。

```zsh
cd /Users/kouji/dev/php
make build
set -o pipefail
python3 research/laravel-api/verify-gophp.py 2>&1 | tee docs/validation/full-stack-go-real-server-mac.log
verify_rc=$?
printf 'exit_code: %s\n' "$verify_rc" | tee -a docs/validation/full-stack-go-real-server-mac.log
```

このログはGo製PHPの実アプリ実行を確認する。公式PHP3版とのfixture厳密比較は、環境回復後に下節の版別`make compare`を別に実行する。

## 2026-10-21：PDOStatementの未定義パラメータとエラー状態

`PDOStatement::execute()`の開始時にStatementの古いSQLSTATEをクリアする。[PHP 8.4のPDO実装](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo/pdo_stmt.c)にある`PDO_STMT_CLEAR_ERR()`を参照した。SQLiteの未定義名前パラメータは[実行前の束縛フック](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo_sqlite/sqlite_statement.c)で位置0へ変わり、SQLiteの`SQLITE_RANGE`（25）は[PHPのSQLiteエラー表](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo_sqlite/sqlite_driver.c)で`HY000`になる。ローカルSQLiteでも位置0への束縛が25を返すことを確認した。Go側は位置0を手前で汎用PDOExceptionとして返しており、Statementの`errorCode()`を更新できなかったため、既存のSQLite束縛・エラー変換経路へ通した。[回帰fixture](../tests/fixtures/success/pdo_parameter_error_state.php)は例外とStatementの`HY000`、`errorInfo`のSQLite番号25と詳細、成功再実行後の`00000`を観測する。今回の期待値は公式ソースとGo実行に基づく**暫定値**であり、公式PHP3版との実比較は未実施。PDOの全エラーモードでの警告数、`errorInfo`のほかの境界、ほかのドライバは未検証として残す。

今回1件を加え、現行**1161件**（正常737・実行時290・構文134）。[Go3設定の対象ログ](validation/full-stack-pdo-param-error-target-local.log)と[全ローカルtest](validation/full-stack-pdo-param-error-test.log)（Python12件を含む）・[vet](validation/full-stack-pdo-param-error-vet.log)・[build](validation/full-stack-pdo-param-error-build.log)・[差分検査](validation/full-stack-pdo-param-error-diff.log)は終了0。最終比較済み997件から追加164件と変更後全件は公式比較待ち。Podmanの同じ接続拒否から環境変化がなく、この単位では試行していない。各版の今回の公式一致0・不一致0・環境エラー0・未実施1161件。Macでは下記の版別対象比較に`pdo_parameter_error_state.php`を追加し、その後に全件比較を行う。

## 2026-10-21：固定アプリの2000リクエスト継続試験

固定アプリの`/api/probe`と`/api/books-manual/{id}`を交互に要求する既存のHTTP・SQLite統合テストを`GOPHP_HTTP_STRESS_REQUESTS=2000`で実行した。[生ログ](validation/full-stack-http-2000-stress-local.log)には2000件、18分45秒、GC後の保持ヒープ865448→744720バイト、test終了0を記録した。これはテスト用の同一Goプロセス・一時SQLiteでのローカル結果で、運用環境での負荷試験や公式PHP3版との厳密比較ではない。試験は下記のHTTPクエリ解析変更を加える前に開始しており、変更後の通常の全ローカルテストは別ログに記録した。

## 2026-10-21：HTTPの配列形式クエリとフォーム入力

HTTP入口の`$_GET`・`$_POST`作成だけがGoの`url.ParseQuery`を使っており、`filter[name]`や`items[]`を平坦なキーとして格納していた。既存の`parse_str()`用の順序付き配列解析を共通化し、GETクエリと`application/x-www-form-urlencoded`のPOST本文へ接続した。`$_POST`の生成はPOSTメソッドとフォームの正確なメディア型に限り、PUT本文と接尾辞の異なるContent-Typeを誤認しないようにした。[外部変数のPHP公式説明](https://www.php.net/manual/en/language.variables.external.php)はフォームの配列形式とキー名中のドット・空白の変換を、[`parse_str()`の説明](https://www.php.net/manual/en/function.parse-str.php)はURLデコードと`[]`による追加を示す。新しい[HTTP統合テスト](../cmd/gophp/main_test.go)は、入れ子・追記・空白のデコード・ドットの変換と、POST以外の本文をリクエストを通して検査し、[対象ログ](validation/full-stack-http-query-array-target-local.log)でGoの実行結果と終了0を確認した。

この時点のfixtureは**1160件**だった。[全ローカルtest](validation/full-stack-http-query-array-test.log)（Python12件を含む）・[vet](validation/full-stack-http-query-array-vet.log)・[build](validation/full-stack-http-query-array-build.log)・[差分検査](validation/full-stack-http-query-array-diff.log)は終了0。この時点のHTTP入力経路と追加163件は公式PHP3版との実比較を行っていない。クエリ解析の`max_input_vars`、multipart/form-data、設定可能な引数区切り、複雑な不正キーの全境界は未対応・未検証として残す。Podmanの既知の接続拒否から環境変化を確認できず、同じ接続試行は繰り返していない。当時の各版の公式一致0・不一致0・環境エラー0・未実施1160件。

## 2026-10-21：LOBを束縛したStatementの再実行

[PHP 8.4のPDO SQLiteパラメータ処理](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo_sqlite/sqlite_statement.c)は、`PDO::PARAM_LOB`へstream resourceを渡すと初回実行時に残りの内容を読み、登録パラメータをその文字列へ置き換える。Go側は登録値をresourceのまま保持していたため、次の`execute()`で終端から読み直して空BLOBを束縛していた。束縛後の変換値をStatementへ保存し、同じStatementを2回実行する[回帰fixture](../tests/fixtures/success/pdo_bind_stream_repeat.php)を追加した。Goの3設定では、両回ともNULを含む同じBLOBを取得し終了0である。公式PHP3版の実比較は未実施で、期待値は公式ソースとGo実行に基づく暫定値である。

この時点で1件を加え、当時**1160件**（正常736・実行時290・構文134）。比較済み997件から追加163件と変更後全件は公式比較待ちだった。直近の[Podman接続拒否](validation/full-stack-podman-2026-10-21-recheck.log)から環境変化を確認できず再試行していない。当時の各版の公式一致0・不一致0・環境エラー0・未実施1160件。[Go対象ログ](validation/full-stack-pdo-lob-repeat-target-local.log)と[test](validation/full-stack-pdo-lob-repeat-test.log)・[vet](validation/full-stack-pdo-lob-repeat-vet.log)・[build](validation/full-stack-pdo-lob-repeat-build.log)・[差分](validation/full-stack-pdo-lob-repeat-diff.log)、[固定アプリCLI・API統合テスト](validation/full-stack-pdo-lob-repeat-app-local.log)を保存する。Macでは下節の対象比較へ`pdo_bind_stream_repeat.php`を追加し、その後に版別全件比較を実行する。全6段階の完成条件は未達である。

## 2026-10-21：PDO名前付き引数の解決時点

[PHP 8.4のPDO SQLiteパラメータ処理](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo_sqlite/sqlite_statement.c)は、名前付きパラメータのSQLite内の位置を実行準備時に求める。Go側が`bindValue(':missing', ...)`で即座に`false`を返していた経路を改め、登録済みの名前・数値位置と順序をStatement内で保持し、`execute()`時に位置を解決して束縛する。存在しない名前による実行失敗を捕捉した後、`execute(['known' => 'ready'])`で以前の束縛を置き換える[回帰fixture](../tests/fixtures/success/pdo_bind_unknown_name.php)を追加した。Goの3設定では出力・終了0を確認したが、公式PHP3版との実比較は未実施で、この期待値は暫定である。

当時1件を加え、**1159件**（正常735・実行時290・構文134）。公式比較済みの997件から追加162件と変更後全件は未比較。既知の[Podman接続拒否](validation/full-stack-podman-2026-10-21-recheck.log)後の環境変化を確認できず、今回は公式比較を再試行していない。各版の今回の一致0・不一致0・環境エラー0・未実施1159件。[Go対象ログ](validation/full-stack-pdo-deferred-bind-target-local.log)と[test](validation/full-stack-pdo-deferred-bind-test.log)・[vet](validation/full-stack-pdo-deferred-bind-vet.log)・[build](validation/full-stack-pdo-deferred-bind-build.log)・[差分](validation/full-stack-pdo-deferred-bind-diff.log)、[`artisan list`の3設定再実行](validation/full-stack-pdo-deferred-bind-artisan-list.log)を保存する。Macでは下節の版別対象比較へ`pdo_bind_unknown_name.php`を追加し、版別全件比較を実施する。引き続き公式の実比較と固定アプリの長時間継続検証が必要で、全6段階は未完了。

Macでの対象比較は、接続回復後に次を実行する。全件比較は同じ版の対象比較後、下節の`make compare`コマンドを使う。`pipefail`を設定し、比較コマンド自身の終了コードを直後に記録する。

```zsh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  read -r tag version target <<< "$spec"
  for name in pdo_statement_argument_validation pdo_bind_stream_order pdo_bind_unknown_name pdo_bind_stream_repeat pdo_parameter_error_state; do
    log="docs/validation/full-stack-pdo-${name}-${tag}-mac.log"
    printf 'command: compare_php.py --php ./containers/php-version/bin/php%s --php-version %s --fixtures tests/fixtures/success/%s.php --interpreter ./bin/gophp --php-version=%s {file}\n' "$tag" "$version" "$name" "$target" > "$log"
    python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php${tag}" --php-version "$version" --fixtures "tests/fixtures/success/${name}.php" --interpreter ./bin/gophp "--php-version=${target}" '{file}' 2>&1 | tee -a "$log"
    compare_rc=$?; printf 'exit_code: %s\n' "$compare_rc" >> "$log"
  done
done
```

## 2026-10-21：PDOパラメータの検査と登録順

`PDOStatement::bindValue()`の`$param`・`$type`、`execute()`の`$params`を共通引数束縛へ接続した。0以下の位置と空の名前はValueError、不正な値型はTypeErrorとして拒否し、`execute(params: [...])`も扱う。[PHP 8.4のPDO実装](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo/pdo_stmt.c)では登録パラメータを配列の挿入順に処理する。Go側のmap走査順による揺れを除くため、Statementに登録順を保持した。引数の拒否後に従前の束縛を使えることを[引数fixture](../tests/fixtures/success/pdo_statement_argument_validation.php)で、同一streamを逆順に束縛した場合を[順序fixture](../tests/fixtures/success/pdo_bind_stream_order.php)で観測する。Goの8.3・8.4・8.5設定では対象3件の出力・終了0を確認した。公式PHPの実比較はまだ行っておらず、この2件の期待値も暫定である。

当時2件を加え、**1158件**（正常734・実行時290・構文134）。最後の公式比較済み997件から追加161件と変更後全件は公式PHP3版で未比較。前回の[Podman接続拒否](validation/full-stack-podman-2026-10-21-recheck.log)後、環境変化を確認できず今回は再試行していないため、各版の今回の公式一致0・不一致0・環境エラー0・未実施1158件。[Go対象ログ](validation/full-stack-pdo-bind-order-target-local.log)と[test](validation/full-stack-pdo-bind-order-test.log)・[vet](validation/full-stack-pdo-bind-order-vet.log)・[build](validation/full-stack-pdo-bind-order-build.log)・[差分検査](validation/full-stack-pdo-bind-order-diff.log)、[`artisan list`の3設定再実行](validation/full-stack-pdo-bind-order-artisan-list.log)を保存した。Macで接続が回復したら、下節の版別比較コマンドで`pdo_statement_argument_validation.php`と`pdo_bind_stream_order.php`も対象比較し、版別`make compare`で全件を比較する。次はPDOパラメータのさらに細かな変換・失敗後状態と固定アプリの長時間継続条件を確認する。全6段階は未完了である。

## 2026-10-21：PDO `bindValue()`の型指定とBLOB書き込み

固定Laravelの[`Connection::bindValues()`](../research/laravel-api/vendor/laravel/framework/src/Illuminate/Database/Connection.php)は`PDO::PARAM_INT`・`PARAM_LOB`・`PARAM_STR`を指定する。[PHP 8.3](https://github.com/php/php-src/blob/PHP-8.3/ext/pdo_sqlite/sqlite_statement.c)・[8.4](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo_sqlite/sqlite_statement.c)・[8.5](https://github.com/php/php-src/blob/PHP-8.5/ext/pdo_sqlite/sqlite_statement.c)のPDO SQLiteパラメータ処理に合わせ、`bindValue()`の`PARAM_NULL`・`PARAM_INT`・`PARAM_BOOL`・`PARAM_STR`・`PARAM_LOB`をSQLiteのnull・整数・text・blob格納へ接続した。LOBには文字列と読み取り可能なstream resourceの基本形を含む。`execute(array)`は[PDO本体](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo/pdo_stmt.c)に合わせ、以前の束縛を置き換えて各引数を`PARAM_STR`として登録し、次の引数なし`execute()`でもその値を再利用する。文字列型への変換は`bindValue()`の登録時に行い、再実行時に`__toString()`を余分に呼ばない。NULを含むBLOB、空BLOB、文字列／整数のSQLite格納型を[対象fixture](../tests/fixtures/success/pdo_bind_value_types.php)の`typeof()`・`hex()`で観測した。Goの3設定で同じ出力と終了0をローカル確認した。名前付き引数の順序と登録時の`__toString()`呼出回数も同じfixtureで観測した。その他のPDOパラメータ型、全resource種別、変換・診断の境界は未実装または未検証である。

当時1件を加え、**1156件**（正常732・実行時290・構文134）。対象fixtureと変更後全件の公式PHP3版比較はPodmanの既知の接続拒否後、環境変化が確認できないため今回未試行。各版の今回の公式一致0・不一致0・環境エラー0・未実施1156件で、期待値は公式ソースとGoローカル結果に基づく暫定値。[対象ローカルログ](validation/full-stack-pdo-bind-types-target-local.log)と[test](validation/full-stack-pdo-bind-types-test.log)・[vet](validation/full-stack-pdo-bind-types-vet.log)・[build](validation/full-stack-pdo-bind-types-build.log)・[差分検査](validation/full-stack-pdo-bind-types-diff.log)、[固定アプリのCLI・API統合テスト](validation/full-stack-pdo-bind-types-app-local.log)を保存した。Macで接続回復後は、下の版別対象比較のfixtureを`tests/fixtures/success/pdo_bind_value_types.php`に指定し、その後に版別`make compare`を実行する。固定アプリのDB経路とCLI・HTTP・API契約はGoローカルで再実行したが、公式PHP3版の変更後比較は残る。全6段階は未完了。Macでの版別対象・全件比較コマンドは次のとおり。実行する場合は各コマンドの終了コードを直後に保存し、Podman終了125をfixture差異へ数えない。

```zsh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  read -r tag version target <<< "$spec"
  target_log="docs/validation/full-stack-pdo-bind-types-target-${tag}-mac.log"
  printf 'command: compare_php.py --php ./containers/php-version/bin/php%s --php-version %s --fixtures tests/fixtures/success/pdo_bind_value_types.php --interpreter ./bin/gophp --php-version=%s {file}\n' "$tag" "$version" "$target" > "$target_log"
  python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php${tag}" --php-version "$version" --fixtures tests/fixtures/success/pdo_bind_value_types.php --interpreter ./bin/gophp "--php-version=${target}" '{file}' 2>&1 | tee -a "$target_log"
  compare_rc=$?; printf 'exit_code: %s\n' "$compare_rc" >> "$target_log"
  all_log="docs/validation/full-stack-pdo-bind-types-all-${tag}-mac.log"
  printf 'command: make compare PHP=./containers/php-version/bin/php%s PHP_VERSION=%s\n' "$tag" "$version" > "$all_log"
  make compare "PHP=./containers/php-version/bin/php${tag}" "PHP_VERSION=${version}" 2>&1 | tee -a "$all_log"
  compare_rc=$?; printf 'exit_code: %s\n' "$compare_rc" >> "$all_log"
done
```

## 2026-10-21：SQLite BLOBの取得値

`sqliteStatement.columnValue()`が`SQLITE_BLOB`をnullに落としていたため、[PDO SQLite 8.4の列取得実装](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo_sqlite/sqlite_statement.c)と同様にバイト長を保ったPHP文字列へ変換した。NULを含むBLOBと空BLOBを[対象fixture](../tests/fixtures/success/pdo_sqlite_blob_fetch.php)で別々に観測し、SQLのNULLとも区別する。Goの3設定で出力・終了0を確認した。LOBパラメータの型指定・resourceを使う書き込みはこの変更で実装していない。

当時1件を加え、**1155件**（正常731・実行時290・構文134）。最後の公式比較済み997件から追加158件と変更後全件は公式PHP3版で未比較、期待値は公式ソースとGoローカル実行に基づく暫定値である。[対象ログ](validation/full-stack-pdo-blob-target-local.log)と[test](validation/full-stack-pdo-blob-test.log)・[vet](validation/full-stack-pdo-blob-vet.log)・[build](validation/full-stack-pdo-blob-build.log)・[差分](validation/full-stack-pdo-blob-diff.log)を保存する。Podmanは前記の接続拒否・終了125から状態変化を確認できず再試行していない。今回の公式比較は各版0件試行、確認済み差異0、1155件未実施。Macでは下の版別対象比較で`tests/fixtures/success/pdo_sqlite_blob_fetch.php`を指定し、版別の`make compare`で全件を比較する。次は固定Laravelが使う`bindValue()`の型指定とBLOB書き込みを公式PDO SQLiteの経路に照らして検証し、アプリのCLI・HTTP・DB・API契約を新しい変更後のGo実装で再確認する。

## 2026-10-21：PDO Statementの列取得と取得モードの版差

`PDOStatement::fetchColumn($column)`が列番号を無視していたため、指定列を返し、負数・列範囲外は行を進めた後でValueErrorを出す経路を追加した。`PDO::FETCH_DEFAULT`と`setFetchMode()`も接続した。[PHP 8.3のStatement実装](https://github.com/php/php-src/blob/PHP-8.3/ext/pdo/pdo_stmt.c)では`setFetchMode(PDO::FETCH_DEFAULT)`後の`fetch()`が`FETCH_BOTH`になるが、[8.4](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo/pdo_stmt.c)・[8.5](https://github.com/php/php-src/blob/PHP-8.5/ext/pdo/pdo_stmt.c)は接続の現在の既定モードを採用するため、Goの対象版設定で分けた。`fetchColumn()`の[実装](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo/pdo_stmt.c)は次の行を取得してから列を読む。

[回帰fixture](../tests/fixtures/success/pdo_statement_fetch_column.php)は2列目・1列目・終了、無効な列、`FETCH_DEFAULT`の版差を観測する。Goの3設定では期待出力・終了0を確認した。今回1件を加え、当時の**1154件**（正常730・実行時290・構文134）。最後の公式比較済み997件から追加157件と変更後全件は公式PHP3版で未比較、期待値は公式ソースとGoローカル実行に基づく暫定値である。[対象ログ](validation/full-stack-pdo-fetch-column-target-local.log)と[test](validation/full-stack-pdo-fetch-column-test.log)・[vet](validation/full-stack-pdo-fetch-column-vet.log)・[build](validation/full-stack-pdo-fetch-column-build.log)・[差分](validation/full-stack-pdo-fetch-column-diff.log)を保存する。[Podman接続拒否](validation/full-stack-podman-2026-10-21-recheck.log)を今回の1回の状態確認で再確認した（終了125）。これはfixture不一致ではない。今回の公式比較は各版0件試行、確認済み差異0、当時の1154件未実施。Macでは下の版別対象比較で`tests/fixtures/success/pdo_statement_fetch_column.php`を指定し、続けて版別の`make compare`で全件を比較する。

## 2026-10-21：PDO SQLiteのタイムアウト属性

既定取得モードの先行実装は、Goが未実装のモード値を`setAttribute()`時に`false`へしていた。これは[PHP 8.4のPDO属性処理](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo/pdo_dbh.c)と異なるため、0以外の整数を設定・取得で保持し、未実装の取得形は実際の`fetch()`で明示的に拒否するよう訂正した。[既定取得モードfixture](../tests/fixtures/success/pdo_default_fetch_mode.php)へ99の保持を追加した。公式PHP3版の実比較は未実施で、この期待値は暫定。

`PDO::ATTR_TIMEOUT`（値2）の直接設定とconstructor optionsを`sqlite3_busy_timeout()`へ接続した。接続時の初期値は60秒。[PHP 8.4のPDO SQLiteドライバ](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo_sqlite/sqlite_driver.c)では設定のみを受け付け、取得は`IM001`となるため、その区別を維持した。[回帰fixture](../tests/fixtures/success/pdo_sqlite_timeout.php)は`PRAGMA busy_timeout`で初期60000ミリ秒、設定後2000ミリ秒、constructor optionsの1000ミリ秒を観測し、取得をPDOExceptionとして捕捉する。実際の競合ロック時の待機・タイムアウト診断は未検証。

Goの8.3・8.4・8.5設定で対象出力・終了0を確認した。今回1件を加え、当時の**1153件**（正常729・実行時290・構文134）。最後の公式比較済み997件から追加156件と変更後全件は公式PHP3版で未比較、期待値は公式ソースとGoローカル出力に基づく暫定値である。[対象ログ](validation/full-stack-pdo-timeout-target-local.log)と[test](validation/full-stack-pdo-timeout-test.log)・[vet](validation/full-stack-pdo-timeout-vet.log)・[build](validation/full-stack-pdo-timeout-build.log)・[差分](validation/full-stack-pdo-timeout-diff.log)、[`artisan list`](validation/full-stack-pdo-timeout-artisan-list.log)を保存する。今回の公式比較は各版0件試行、環境エラー0、当時の1153件未実施。Macでは下の版別対象比較で`tests/fixtures/success/pdo_sqlite_timeout.php`を指定し、版別の`make compare`で全件を比較する。Podmanの接続状態に変化を確認できていないため、既知の接続拒否を繰り返し試していない。

## 2026-10-21：PDOの既定取得モードをStatementへ固定

`PDO::ATTR_DEFAULT_FETCH_MODE`の既定値を`FETCH_BOTH`とし、`setAttribute()`・`getAttribute()`とconstructor optionsへ接続した。対応する取得形は`FETCH_ASSOC`・`FETCH_NUM`・`FETCH_BOTH`・`FETCH_OBJ`・`FETCH_COLUMN`。未対応の取得形は設定成功を偽装せず、実装境界として残す。`FETCH_DEFAULT`（値0）はValueErrorを返す。[PHP 8.4 PDOの属性処理](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo/pdo_dbh.c)ではStatement準備時に接続の既定値をコピーするため、Go側も作成済みStatementに後からの接続属性変更を反映しない。

[回帰fixture](../tests/fixtures/success/pdo_default_fetch_mode.php)は初期値、設定前後のStatement、無効な0を確認する。Goの3設定で出力・終了0を確認した。今回1件を加え、当時の**1152件**（正常728・実行時290・構文134）。最後の公式比較済み997件から追加155件と変更後全件は公式PHP3版で未比較、期待値は公式ソースとGoローカル出力に基づく暫定値である。[対象ログ](validation/full-stack-pdo-default-fetch-target-local.log)と[test](validation/full-stack-pdo-default-fetch-test.log)・[vet](validation/full-stack-pdo-default-fetch-vet.log)・[build](validation/full-stack-pdo-default-fetch-build.log)・[差分](validation/full-stack-pdo-default-fetch-diff.log)、[`artisan list`](validation/full-stack-pdo-default-fetch-artisan-list.log)を保存する。今回の公式比較は各版0件試行、環境エラー0、当時の1152件未実施。Macでは下の版別対象比較で`tests/fixtures/success/pdo_default_fetch_mode.php`を指定し、版別の`make compare`で全件を比較する。Podmanの接続状態に変化が確認できていないため、既知の接続拒否を繰り返し試していない。

## 2026-10-21：固定vendorが設定するPDO取得属性

Laravel Connectorの既定オプションにある`ATTR_CASE`・`ATTR_ORACLE_NULLS`・`ATTR_STRINGIFY_FETCHES`をSQLite接続へ保持し、行の列名と値の取得へ反映した。列名の大小文字変換はStatementの初回実行時に確定し、その後の接続属性変更では既存Statementの列名を変えない。[PHP 8.4の汎用属性処理](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo/pdo_dbh.c)、[取得値処理](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo/pdo_stmt.c)と[SQLiteドライバの属性表](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo_sqlite/sqlite_driver.c)を参照した。SQLiteが受け付けない`ATTR_EMULATE_PREPARES`の直接設定は`false`とし、未対応属性の取得は`IM001`のPDOExceptionへ運ぶ。PHP 8.3の`ATTR_ORACLE_NULLS`不正値検査は8.4以降と異なるため版を分けた。

[PDO属性fixture](../tests/fixtures/runtime/pdo_fetch_attributes.php)は列名、数値の文字列化、空文字列とnull、constructor options、未対応属性、Statementの列名確定時点を観測する。[版別の不正値fixture](../tests/fixtures/success/pdo_oracle_nulls_version.php)は8.3と8.4以降の`ATTR_ORACLE_NULLS`宣言検査差を確認する。Goの3設定で対象出力・終了0を確認した。追加2件、当時の**1151件**（正常727・実行時290・構文134）。最後の公式比較済み997件から追加154件と変更後全件の公式PHP3版比較は未実施で、期待値は公式ソースとGoローカル実行に基づく暫定値である。この時点では`ATTR_DEFAULT_FETCH_MODE`、SQLiteのタイムアウト属性、属性の全診断境界、Blobや複数ドライバの値変換が未対応・未検証だった。

[Go対象ログ](validation/full-stack-pdo-fetch-attributes-target-local.log)、[test](validation/full-stack-pdo-fetch-attributes-test.log)・[vet](validation/full-stack-pdo-fetch-attributes-vet.log)・[build](validation/full-stack-pdo-fetch-attributes-build.log)・[差分](validation/full-stack-pdo-fetch-attributes-diff.log)、[`artisan list`](validation/full-stack-pdo-fetch-attributes-artisan-list.log)を保存する。Podman接続拒否から状態変化を確認できず再試行していない。各版の今回の公式比較は0件実施、確認済み差異0、今回の環境エラー0、変更後1151件未実施。Macでは下の版別コマンドの対象一覧へ`pdo_fetch_attributes`・`pdo_oracle_nulls_version`を加える。

## 2026-10-21：PDOのエラーモード

固定vendorのSymfony `PdoSessionHandler`は`PDO::getAttribute(PDO::ATTR_ERRMODE)`を読み、LaravelのConnectorは`PDO::ATTR_ERRMODE`を設定する。[PHP 8.4のPDO実装](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo/pdo_dbh.c)では既定が`ERRMODE_EXCEPTION`で、`ERRMODE_SILENT`・`ERRMODE_WARNING`・`ERRMODE_EXCEPTION`を受け付ける。Go側のSQLite接続にエラーモードを保持し、constructor optionsと`setAttribute()`・`getAttribute()`の基本経路へ接続した。SQL実行失敗時に例外、警告、`false`を分け、警告ハンドラの例外で後続処理を中断する。

[モード切替fixture](../tests/fixtures/success/pdo_errmode_modes.php)は既定値、silent／warning、SQLSTATE保持、不正値、constructor optionsを、[警告例外fixture](../tests/fixtures/success/pdo_errmode_handler_throw.php)はハンドラが投げた後の中断を、[Statementとqueryのfixture](../tests/fixtures/success/pdo_errmode_statement_query.php)はsilent時の戻り値と接続・StatementのSQLSTATE分離を確認する。Goの8.3・8.4・8.5設定で3件の出力・終了0を確認した。追加3件、現行**1149件**（正常725・実行時290・構文134）。最後の公式比較済み997件から追加152件と変更後全件は未比較で、今回の期待値は公式ソースとGoローカル実行に基づく暫定値である。この時点では`PDO::ATTR_ERRMODE`以外の属性一般、warningの全診断文・位置、全PDO操作とドライバでのモード動作は未対応または未検証。既存の他属性の簡略処理をPHP互換と扱わない。

[Go対象ログ](validation/full-stack-pdo-errmode-target-local.log)、[test](validation/full-stack-pdo-errmode-test.log)・[vet](validation/full-stack-pdo-errmode-vet.log)・[build](validation/full-stack-pdo-errmode-build.log)・[差分](validation/full-stack-pdo-errmode-diff.log)、[`artisan list`](validation/full-stack-pdo-errmode-artisan-list.log)を保存する。Podmanの接続状態に変化を確認できず、今回の公式比較は各版0件実施、確認済み差異0、今回の環境エラー0、変更後1149件未実施。Macでは下の版別比較コマンドの対象一覧へ`pdo_errmode_modes`・`pdo_errmode_handler_throw`・`pdo_errmode_statement_query`を加え、版ごとの全件比較を行う。

## 2026-10-21：`PDO::query()`失敗時の接続エラー照会

[PHP 8.4のPDO::query()実装](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo/pdo_dbh.c)は実行に失敗したStatementを接続ハンドルの`query_stmt`へ保持し、`PDO::errorCode()`／`errorInfo()`はそのStatementの状態を参照する。次の接続操作では[PDOの状態クリア処理](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo/php_pdo_error.h)が保持先を外す。Go側も`query()`実行失敗後の接続照会をこの経路へ接続した。準備済みStatementの`execute()`失敗が接続側へ漏れない従来の区別は維持する。

[query失敗fixture](../tests/fixtures/success/pdo_query_error_state.php)は一意制約違反後の接続SQLSTATE・SQLite番号と、次の`exec()`成功による状態クリアを観測する。Goの8.3・8.4・8.5設定で出力・終了0を確認した。追加1件、現行**1146件**（正常722・実行時290・構文134）。公式PHP3版の対象・全件比較は未実施で、期待値は公式ソースとGoローカル実行に基づく暫定値である。[対象ローカルログ](validation/full-stack-pdo-query-state-target-local.log)と[test](validation/full-stack-pdo-query-state-test.log)・[vet](validation/full-stack-pdo-query-state-vet.log)・[build](validation/full-stack-pdo-query-state-build.log)・[差分](validation/full-stack-pdo-query-state-diff.log)を保存する。Podmanは直近の接続拒否から状態変化が確認できず再試行していない。今回の各版比較は0件実施、確認済み不一致0、今回の環境エラー0、変更後1146件未実施。`PDO::ATTR_ERRMODE`の全モード、`query()`の他の失敗位置や他ドライバは未検証。

接続が回復したMacでは、PDO状態と属性の現行11件を対象比較し、同じ版の全件比較を順番に行う。各比較の実際の終了コードを直後に記録する。

```zsh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  read -r tag version target <<< "$spec"
  for fixture in tests/fixtures/success/{pdo_error_state_handles,pdo_query_error_state,pdo_errmode_modes,pdo_errmode_handler_throw,pdo_errmode_statement_query,pdo_oracle_nulls_version,pdo_default_fetch_mode,pdo_sqlite_timeout,pdo_statement_fetch_column,pdo_sqlite_blob_fetch}.php tests/fixtures/runtime/pdo_fetch_attributes.php; do
    log="docs/validation/full-stack-pdo-target-$(basename "$fixture" .php)-${tag}-mac.log"
    python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php${tag}" --php-version "$version" --fixtures "$fixture" --interpreter ./bin/gophp "--php-version=${target}" '{file}' 2>&1 | tee "$log"
    check_rc=$?; printf 'exit_code: %s\n' "$check_rc" >> "$log"
  done
  log="docs/validation/full-stack-pdo-all-${tag}-mac.log"
  make compare "PHP=./containers/php-version/bin/php${tag}" "PHP_VERSION=${version}" 2>&1 | tee "$log"
  check_rc=$?; printf 'exit_code: %s\n' "$check_rc" >> "$log"
done
```

## 2026-10-21：PDO接続とStatementのエラー状態

[PDO::errorCode()](https://www.php.net/manual/en/pdo.errorcode.php)と[PDO::errorInfo()](https://www.php.net/manual/en/pdo.errorinfo.php)は接続ハンドルへ直接行った操作の状態を返し、準備済みStatementの操作は別の[PDOStatement::errorInfo()](https://www.php.net/manual/en/pdostatement.errorinfo.php)で読む。Go側のSQLite接続とStatementへ独立した状態を追加し、`errorCode()`と`errorInfo()`を組み込みメソッド経路へ接続した。接続直後・未実行Statementの`errorCode()`はnull、`errorInfo()[0]`は空文字列とする。直接SQL・準備済みStatementの成功後は`00000`、SQLite失敗後はそれぞれのSQLSTATEと番号・詳細を返す。成功した後は、そのハンドルの古い失敗情報を消す。準備済みStatementの`execute()`失敗を接続側へ写さない。[PHP 8.4のPDO実装](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo/pdo_dbh.c)にある`PDO::query()`の接続状態特例はこの時点では未実装で、上の後続記録で扱う。

[ハンドル分離fixture](../tests/fixtures/success/pdo_error_state_handles.php)は接続直後、直接SQL失敗、成功後のクリア、準備済みStatementの一意制約失敗、接続側状態の維持、Statement再実行後のクリアを順に確認する。Goの3設定で対象出力・終了0を確認した。追加1件、現行**1145件**（正常721・実行時290・構文134）。最後の公式PHP3版比較済み997件からの追加148件と変更後全件は未比較で、新規期待値は公式マニュアルとGoローカル実行に基づく暫定値である。`PDO::ATTR_ERRMODE`の全モード、すべてのSQLite操作後の状態遷移、その他ドライバは未対応・未検証。

[対象ローカルログ](validation/full-stack-pdo-handle-state-target-local.log)、[test](validation/full-stack-pdo-handle-state-test.log)・[vet](validation/full-stack-pdo-handle-state-vet.log)・[build](validation/full-stack-pdo-handle-state-build.log)・[差分](validation/full-stack-pdo-handle-state-diff.log)、[`artisan list`](validation/full-stack-pdo-handle-state-artisan-list.log)を記録する。Podmanは直近の[接続拒否・終了125](validation/full-stack-podman-recheck.log)から状態変化を確認できず、今回は再試行していない。公式PHP各版の今回の実比較0件、確認済み不一致0、今回の環境エラー0、変更後1145件未実施。Macでは次節の版別`compare_php.py`コマンドへ`pdo_error_state_handles`を追加して対象比較し、続けて版別`make compare`で全件を比較する。

## 2026-10-21：SQLiteのSQLSTATEとLaravel一意制約例外

PHP公式の[8.3](https://github.com/php/php-src/blob/PHP-8.3/ext/pdo_sqlite/sqlite_driver.c)・[8.4](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo_sqlite/sqlite_driver.c)・[8.5](https://github.com/php/php-src/blob/PHP-8.5/ext/pdo_sqlite/sqlite_driver.c)ではSQLiteの`SQLITE_ERROR`を`HY000`、`SQLITE_CONSTRAINT`を`23000`へ写す。各版の[PDO 8.3](https://github.com/php/php-src/blob/PHP-8.3/ext/pdo/pdo_dbh.c)・[8.4](https://github.com/php/php-src/blob/PHP-8.4/ext/pdo/pdo_dbh.c)・[8.5](https://github.com/php/php-src/blob/PHP-8.5/ext/pdo/pdo_dbh.c)は、例外の`code`にSQLSTATE文字列を、`errorInfo`にSQLSTATE・SQLiteエラー番号・詳細文を格納する。Go側も実際のSQLite失敗を共通Throwable経路へ運ぶ際に、この情報を保持するよう変更した。手動で作る例外やGo固有の未対応PDO操作へSQLSTATEを捏造しない。

[PDOの失敗fixture](../tests/fixtures/success/pdo_exception_sqlstate.php)はテーブル未発見と一意制約違反を独立して捕捉し、`getCode()`と`errorInfo`の3要素を観測する。[固定Laravelのfixture](../tests/fixtures/success/eloquent_unique_constraint_vendor.php)はSQLiteの一意制約違反が`UniqueConstraintViolationException`となり、`PDOException`を`previous`に保持する経路を確認する。Goの8.3・8.4・8.5設定で対象2件の出力・終了0を確認した。追加2件、当時**1144件**（正常720・実行時290・構文134）、最後の公式比較済み997件からの追加147件と変更後全件の公式PHP3版比較は未実施で期待値は暫定。SQLiteの全エラー番号、PDOの全エラーモード・SQLSTATE表・接続失敗診断は未対応範囲として残る。

[対象ローカルログ](validation/full-stack-pdo-state-target-local.log)、[test](validation/full-stack-pdo-state-test.log)・[vet](validation/full-stack-pdo-state-vet.log)・[build](validation/full-stack-pdo-state-build.log)・[差分](validation/full-stack-pdo-state-diff.log)、[`artisan list`](validation/full-stack-pdo-state-artisan-list.log)を記録する。Podmanの直近の[接続拒否・終了125](validation/full-stack-podman-recheck.log)から環境変化が確認できず、同じ失敗を再試行していない。今回の実比較は各版とも0件、確認済み不一致0、今回の環境エラー0、変更後1144件未実施。Macでは次節の版別比較コマンドの`name`一覧へ`pdo_exception_sqlstate`と`eloquent_unique_constraint_vendor`を追加し、対象比較後に版別全件比較を実行する。

## 2026-10-21：Laravel/EloquentカーソルとPDO例外

固定vendorの`Illuminate\\Database\\Connection::cursor()`はGeneratorを返し、最初の反復までSQL実行を遅らせる。Query Builderの`cursor()`はこのGeneratorを`LazyCollection`から委譲する。[遅延実行fixture](../tests/fixtures/success/eloquent_cursor_vendor.php)では、2つのカーソル生成後にSQLite表と行を作り、直接カーソルの束縛引数とQuery Builderの全行走査を確認した。Goの8.3・8.4・8.5設定では期待出力と一致した。

SQL失敗時の`PDOException`がThrowable階層に未登録だったため、Laravelの`catch (Exception)`を通らず未捕捉になっていた。`PDOException extends RuntimeException`を組み込みクラスへ登録し、`Exception`の保護された`message`・`code`と`PDOException`の公開`errorInfo`をプロパティ経路に加えた。継承先が`message`・`code`を書いた場合も取得メソッドの内容へ反映する。[直接PDO例外](../tests/fixtures/success/pdo_exception_catch.php)と[LaravelのQueryException連鎖](../tests/fixtures/success/eloquent_cursor_failure_vendor.php)を追加し、Goの3設定で余分な診断なしに捕捉できることを確認した。当時はPDOのSQLSTATE、`errorInfo`の非null内容が残件だった。上節で代表的なSQLite失敗の基本形を補完したが、DB例外診断の全面互換は残る。

追加3件、現行**1142件**（正常718・実行時290・構文134）。[対象ローカルログ](validation/full-stack-eloquent-cursor-target-local.log)、[test](validation/full-stack-eloquent-cursor-test.log)・[vet](validation/full-stack-eloquent-cursor-vet.log)・[build](validation/full-stack-eloquent-cursor-build.log)・[差分検査](validation/full-stack-eloquent-cursor-diff.log)および[`artisan list`](validation/full-stack-eloquent-cursor-artisan-list.log)を記録した。公式PHP3版の今回の対象比較・変更後全件比較はいずれも未実施。最後の[Podman接続拒否・終了125](validation/full-stack-podman-recheck.log)以降の環境変化を確認できないため、同じ接続失敗を繰り返し試していない。各版とも今回の確認済み一致0・不一致0・環境エラー0・未実施1142件で、新規期待値は暫定。Macで接続回復後は、下記の既存版別`compare_php.py`へ`eloquent_cursor_vendor.php`、`eloquent_cursor_failure_vendor.php`、`pdo_exception_catch.php`をそれぞれ`--fixtures tests/fixtures/success/...php`として指定し、続けて版別`make compare`で全件を比較する。節目では固定アプリのCLI・HTTP・DB・API契約を同じ入力で公式PHPと照合する。

```zsh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  read -r tag version target <<< "$spec"
  for name in eloquent_cursor_vendor eloquent_cursor_failure_vendor pdo_exception_catch; do
    log="docs/validation/full-stack-eloquent-cursor-${name}-${tag}-mac.log"
    printf 'command: compare_php.py --php ./containers/php-version/bin/php%s --php-version %s --fixtures tests/fixtures/success/%s.php --interpreter ./bin/gophp --php-version=%s {file}\n' "$tag" "$version" "$name" "$target" > "$log"
    python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php${tag}" --php-version "$version" --fixtures "tests/fixtures/success/${name}.php" --interpreter ./bin/gophp "--php-version=${target}" '{file}' 2>&1 | tee -a "$log"
    check_rc=$?; printf 'exit_code: %s\n' "$check_rc" >> "$log"
  done
  log="docs/validation/full-stack-eloquent-cursor-all-${tag}-mac.log"
  printf 'command: make compare PHP=./containers/php-version/bin/php%s PHP_VERSION=%s\n' "$tag" "$version" > "$log"
  make compare "PHP=./containers/php-version/bin/php${tag}" "PHP_VERSION=${version}" 2>&1 | tee -a "$log"
  check_rc=$?; printf 'exit_code: %s\n' "$check_rc" >> "$log"
done
```

## 2026-10-21：継承メソッドのトレースとGeneratorの開始順

PHPの[8.3](https://raw.githubusercontent.com/php/php-src/PHP-8.3/Zend/zend_builtin_functions.c)・[8.4](https://raw.githubusercontent.com/php/php-src/PHP-8.4/Zend/zend_builtin_functions.c)・[8.5](https://raw.githubusercontent.com/php/php-src/PHP-8.5/Zend/zend_builtin_functions.c)の`zend_fetch_debug_backtrace()`は、メソッドフレームの`class`を関数の宣言スコープから取得し、`object`キーを`type`・`args`より先に格納する。Go側は継承先の実オブジェクト名を表示し、キー順も異なっていた。呼び出しフレームのクラス名を宣言元に修正し、配列生成順を揃えた。[継承メソッドのfixture](../tests/fixtures/success/debug_backtrace_declaring_class.php)は宣言元名、受け手の同一性、キー順を観測する。固定vendorの`once()`とEloquent再帰防止traitを含む既存対象も再実行した。

Generatorの[8.3](https://raw.githubusercontent.com/php/php-src/PHP-8.3/Zend/zend_generators.c)・[8.4](https://raw.githubusercontent.com/php/php-src/PHP-8.4/Zend/zend_generators.c)・[8.5](https://raw.githubusercontent.com/php/php-src/PHP-8.5/Zend/zend_generators.c)の`getReturn()`は未開始時に最初の`yield`まで実行してから戻り値を判定する。マニュアルの説明だけから「実行せずに例外を返す」と誤判断したが、公式ソースで確認して元のGo動作を維持した。[副作用を観測するfixture](../tests/fixtures/success/generator_get_return_lazy.php)では最初の`yield`で停止する場合と、実行時に`yield`へ到達せず直ちにreturnする場合を分けた。既存の[`generator_initial_methods.php`](../tests/fixtures/success/generator_initial_methods.php)の期待値は変更していない。この記録は公式ソースからの判断で、公式PHPの実出力による確認ではない。

当時1139件（正常715・実行時290・構文134）、最後の公式比較済み基準997件からの追加142件と変更後全件は3版とも実比較待ち。[Goの3設定で対象6件](validation/full-stack-trace-scope-target-local.log)、[make test](validation/full-stack-trace-scope-test.log)・[vet](validation/full-stack-trace-scope-vet.log)・[build](validation/full-stack-trace-scope-build.log)・[差分検査](validation/full-stack-trace-scope-diff.log)、変更後の[`artisan list`](validation/full-stack-trace-scope-artisan-list.log)は終了0。Podmanは直近の[接続拒否・終了125](validation/full-stack-podman-recheck.log)から状態変化が確認できず再試行していない。今回は公式PHP3版で各0件実施・0件の確認済み不一致・1139件未実施。Macでは次節の版別コマンドの対象一覧に`debug_backtrace_declaring_class.php`・`generator_get_return_lazy.php`・`generator_initial_methods.php`を追加し、版ごとの全件比較を実行する。

## 2026-10-21：呼び出し元オブジェクトのバックトレース

固定vendorの`Illuminate\Support\Onceable::tryFromTrace()`は`debug_backtrace(DEBUG_BACKTRACE_PROVIDE_OBJECT, 2)`の呼び出し元フレームに含まれる`object`で`once()`キャッシュをオブジェクト別に分ける。従来のGo実装は最上位フレームだけで現在の`$this`を探索しており、メソッドから別関数を呼んだ場合に呼び出し元の`object`が欠けた。このため異なる2オブジェクトの`once()`が同じ値を共有した。メソッド呼び出しフレームへ受け手オブジェクトを保持し、`debug_backtrace()`の各フレームへ`DEBUG_BACKTRACE_PROVIDE_OBJECT`指定時だけ渡すよう修正した。[PHP公式マニュアル](https://www.php.net/manual/en/function.debug-backtrace.php)の`object`キーの説明と固定vendorの実コードに基づく。`object`を表示用文字列から復元せず、元のオブジェクトの同一性を維持する。

[呼び出し元フレームのfixture](../tests/fixtures/success/debug_backtrace_caller_object.php)はインスタンスメソッドの`object`同一性、静的メソッドでの省略、オプション0での省略を検査する。[固定vendorの`once()` fixture](../tests/fixtures/success/laravel_once_object_scope.php)は同じオブジェクトへの2回の呼び出しと別オブジェクトへの呼び出しを分ける。さらに[固定vendorのEloquent再帰防止trait](../tests/fixtures/success/eloquent_recursion_cache_vendor.php)を直接実行し、再帰時の既定値、オブジェクト間の分離、例外後のキャッシュ解除を確認した。現行1137件（正常713・実行時290・構文134）のうち最後の公式比較済み997件からの追加140件は、3版とも未比較で期待値は暫定。最初の2件の[Go対象6実行](validation/full-stack-backtrace-object-target-local.log)、追加のEloquent fixtureの[Go対象3実行](validation/full-stack-eloquent-recursion-target-local.log)、追加後の[make test](validation/full-stack-eloquent-recursion-test.log)・[vet](validation/full-stack-eloquent-recursion-vet.log)・[build](validation/full-stack-eloquent-recursion-build.log)・[差分検査](validation/full-stack-eloquent-recursion-diff.log)は終了0。[変更後の`artisan list`](validation/full-stack-backtrace-object-artisan-list.log)も終了0。新しい[Podman試行](validation/full-stack-podman-recheck.log)は接続拒否で終了125。これはfixture不一致に含めない。今回の公式PHP3版は各0件比較、0件の確認済み差異、変更後全1137件が未実施。

Macで接続回復後、まず現行の対象6件を比較し、その後版ごとの全件比較を行う。コマンドの終了コードをパイプライン直後に保存する。

```zsh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  read -r tag version target <<< "$spec"
  for fixture in tests/fixtures/success/{debug_backtrace_caller_object,laravel_once_object_scope,eloquent_recursion_cache_vendor,debug_backtrace_declaring_class,generator_get_return_lazy,generator_initial_methods}.php; do
    log="docs/validation/full-stack-backtrace-target-$(basename "$fixture" .php)-${tag}-mac.log"
    python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php${tag}" --php-version "$version" --fixtures "$fixture" --interpreter ./bin/gophp "--php-version=${target}" '{file}' 2>&1 | tee "$log"
    check_rc=$?; printf 'exit_code: %s\n' "$check_rc" >> "$log"
  done
  log="docs/validation/full-stack-backtrace-all-${tag}-mac.log"
  make compare "PHP=./containers/php-version/bin/php${tag}" "PHP_VERSION=${version}" 2>&1 | tee "$log"
  check_rc=$?; printf 'exit_code: %s\n' "$check_rc" >> "$log"
done
```

現在の固定アプリでの`once()`オブジェクト別キャッシュは小fixtureでGoローカル確認した。固定アプリの全経路でのキャッシュ寿命、Closure呼び出しフレームの`object`境界、PHP公式3版との厳密比較は残る。次の作業はこれらの境界と同一入力によるCLI・HTTP・DB・API出力比較を進める。段階1〜6は未完了。

## 2026-10-21：API Platformの版別CRUDとキャッシュの値コピー

従来、固定アプリのAPI Platform CRUDを一連で確認していたGo統合テストは8.4設定だった。[新しい版別統合テスト](../cmd/gophp/main_test.go)は8.3・8.5設定でも、独立した一時SQLite・storageを用いてPOST→GET→PATCH→DELETE→削除後404とDB行数0を検査し、[実行ログ](validation/full-stack-api-version-matrix-local.log)では両設定とも終了0だった。これにより固定アプリの基本CRUDはGoの3設定でローカル確認済みだが、公式PHPとの現行実装の比較は未実施である。

CLIの`artisan route:list`もGoの[8.3](validation/full-stack-route-list-8.3-local.log)・[8.4](validation/full-stack-route-list-8.4-local.log)・[8.5](validation/full-stack-route-list-8.5-local.log)設定で個別に終了0を確認した。これらは版別のGo実行ログであり、公式PHP CLIとの出力比較ではない。

`api:openapi:export`も新しい一時storageでGoの[8.3](validation/full-stack-openapi-8.3-local.log)・[8.4](validation/full-stack-openapi-8.4-local.log)・[8.5](validation/full-stack-openapi-8.5-local.log)設定を個別に実行し、各終了0、JSONの`openapi=3.2.0`と`/api/books`のパスを確認した。Goの3設定で生成したJSON本文は各16029バイトでSHA-256も一致した。公式PHPの出力との厳密比較はまだ行っていない。

`CachingIterator::getCache()`の結果は配列の値コピーとして扱い、返却後の配列変更と本体への`[]`書き込み・削除が相互に漏れないことを[新規fixture](../tests/fixtures/success/caching_iterator_array_copy.php)で確認した。[Goの3設定の実行ログ](validation/full-stack-caching-array-copy-local.log)はそれぞれ同じ暫定期待値・終了0を記録する。[PHP 8.4のSPL実装](https://github.com/php/php-src/blob/PHP-8.4/ext/spl/spl_iterators.c)は`getCache()`で配列をコピーして返す。現行1134件（正常710・実行時290・構文134）のうち、最後の公式比較済み基準997件からの追加137件と変更後全件は未比較。次の版別対象比較ではこのfixtureも下記の4件へ加え、全件は既存の`make compare`を使う。

追加後の[make test](validation/full-stack-api-version-matrix-test.log)（Python 12件を含む）・[vet](validation/full-stack-api-version-matrix-vet.log)・[build](validation/full-stack-api-version-matrix-build.log)・[差分検査](validation/full-stack-api-version-matrix-diff.log)はいずれも終了0。CLI・APIのGoローカル統合は進んだが、公式PHP3版の新規対象・全件比較と固定アプリの同一入力による版別比較は残る。

## 2026-10-21：CachingIteratorのArrayAccess境界

[キャッシュのArrayAccess fixture](../tests/fixtures/success/caching_iterator_array_access.php)で、`offsetExists`・`offsetGet`・`offsetSet`・`offsetUnset`、件数更新と未知キーの`E_WARNING`を確認した。警告文とレベルの期待値は[PHP 8.4のSPL実装](https://github.com/php/php-src/blob/PHP-8.4/ext/spl/spl_iterators.c)から設定し、Goの対象テストは通過した。当時1133件（正常709・実行時290・構文134）。新規1件を含む136件の追加分と変更後全件は公式PHP3版未比較で、期待値は暫定。Macでの対象比較では下記コマンドのfixture列へ`tests/fixtures/success/caching_iterator_array_access.php`も加える。通常の配列構文でのArrayAccess、引数型、反復中のキャッシュ変更は別の境界として残る。

## 2026-10-21：CachingIteratorの先読み順と内部フラグ

[PHP 8.4のSPL実装](https://github.com/php/php-src/blob/PHP-8.4/ext/spl/spl_iterators.c)の`spl_dual_it_fetch()`は、先読み時に内側の`current()`を`key()`より先に呼ぶ。Go側は逆だったため、[呼び出し順fixture](../tests/fixtures/success/caching_iterator_call_order.php)で副作用を観測して修正した。さらにSPLの`CIT_PUBLIC=0xffff`相当のマスクを適用し、下位16ビットの公開値は保持しつつ、利用者が`CIT_VALID`内部ビット（65536）をコンストラクタや`setFlags()`へ渡して有効状態を偽装できないようにした。[フラグ変更](../tests/fixtures/success/caching_iterator_flags_transition.php)と[内部ビット](../tests/fixtures/success/caching_iterator_private_flags.php)の2件も追加した。現在1132件（正常708・実行時290・構文134）。3件の期待値は公式ソースからの推論とGoローカル実行に基づく暫定値で、公式PHPの実比較結果ではない。

[対象のGo3設定](validation/full-stack-caching-boundary-local.log)、[make test](validation/full-stack-caching-boundary-test.log)・[vet](validation/full-stack-caching-boundary-vet.log)・[build](validation/full-stack-caching-boundary-build.log)、変更後の[`artisan list`](validation/full-stack-caching-boundary-artisan-list.log)は終了0。当時の新規3件と変更後1132件の公式PHP3版比較は未実施で、接続環境は前回の[Podman終了125](validation/full-stack-caching-podman.log)から変化がないため再試行していない。各版の今回の実比較0件、確認済み差異0件、変更後全件1132件未実施。Macでは次の版別コマンドで5件の`tests/fixtures/success/caching_iterator_{call_order,flags_transition,private_flags,array_access,array_copy}.php`を個別に`--fixtures`へ渡し、その後`make compare`を3版で実行する。次はArrayAccess診断と反復中変更の境界を調べ、固定アプリの実経路に必要な不足を実装する。

下位16ビットの保持を加えた後も[make test](validation/full-stack-caching-public-bits-test.log)・[vet](validation/full-stack-caching-public-bits-vet.log)・[build](validation/full-stack-caching-public-bits-build.log)は終了0。これはGoローカル検証であり、公式PHPの実比較は未実施。

```zsh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  read -r tag version target <<< "$spec"
  for fixture in tests/fixtures/success/caching_iterator_{call_order,flags_transition,private_flags,array_access,array_copy}.php; do
    log="docs/validation/full-stack-caching-boundary-$(basename "$fixture" .php)-${tag}-mac.log"
    python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php${tag}" --php-version "$version" --fixtures "$fixture" --interpreter ./bin/gophp "--php-version=${target}" '{file}' 2>&1 | tee "$log"
    check_rc=$?; printf 'exit_code: %s\n' "$check_rc" >> "$log"
  done
  log="docs/validation/full-stack-caching-boundary-all-${tag}-mac.log"
  make compare "PHP=./containers/php-version/bin/php${tag}" "PHP_VERSION=${version}" 2>&1 | tee "$log"
  check_rc=$?; printf 'exit_code: %s\n' "$check_rc" >> "$log"
done
```

## 2026-10-21：固定アプリの1024リクエスト継続検証

一時SQLite・一時storageを使う[`TestFixedLaravelProbeThroughHTTPBridge`](../cmd/gophp/main_test.go)を`GOPHP_HTTP_STRESS_REQUESTS=1024`で再実行した。[生ログ](validation/full-stack-http-1024-caching.log)は追加の1024 GET、所要7分59秒、GC後保持ヒープ908936→705504バイト、テストとコマンドの終了0を記録する。同テストでは先行する手動API・API PlatformのCRUD、422、OpenAPI JSON、並列8 GETも検査する。これで以前の512 GETより長いGoローカルの状態分離確認が増えた。公式PHP3版の現行fixture比較や、時間・負荷条件を変えた運用上の長時間試験を代替しない。

## 2026-10-21：CachingIteratorとLaravel Collection

固定vendorの`Illuminate\Collections\Traits\EnumeratesValues::getCachingIterator()`は`CachingIterator`を生成する。Go側へ基本の`CachingIterator`を登録し、Iteratorを1件先読みしつつ現在値・キーを保持する。`hasNext()`は先読み済みの内部Iteratorの位置を読む。`FULL_CACHE`の要素保存、`getCache()`・`count()`、ArrayAccess、flagsの基本形、文字列化の基本形を追加した。先読み順とキャッシュ挿入位置は[PHP 8.4のSPL実装](https://github.com/php/php-src/blob/PHP-8.4/ext/spl/spl_iterators.c)を参照した。Laravel Collectionから`getCachingIterator()`を呼ぶ[固定vendor fixture](../tests/fixtures/success/caching_collection_vendor.php)と、通常走査・フルキャッシュ・先読み位置・失敗系の5件を追加した。当時1129件（正常705・実行時290・構文134）。[対象fixtureのGo4設定](validation/full-stack-caching-fixtures-local.log)、[test](validation/full-stack-caching-test.log)・[vet](validation/full-stack-caching-vet.log)・[build](validation/full-stack-caching-build.log)・[差分検査](validation/full-stack-caching-diff.log)、変更後の[`artisan list`](validation/full-stack-caching-artisan-list.log)・[`about`](validation/full-stack-caching-artisan-about.log)・[`route:list`](validation/full-stack-caching-artisan-route-list.log)・[OpenAPI export](validation/full-stack-caching-openapi-export.log)はいずれも終了0を記録する。

当時、公式PHP 8.3.35・8.4.26・8.5.11との新規5件および変更後1129件の実比較は[Podman終了125の接続拒否](validation/full-stack-caching-podman.log)が継続し未実施。各版0件の新規実比較、0件の確認済み差異、1129件の変更後全件比較は未実施で、期待値は暫定である。`CachingIterator`の`setFlags()`境界、全ArrayAccess診断、反復中の変更・参照、`RecursiveCachingIterator`は未検証／未実装。Macで接続回復後は下記の対象比較コマンドの`--fixtures`へ`tests/fixtures/success/caching_{iterator_basic,iterator_cache,iterator_reject,iterator_lookahead,collection_vendor}.php`をそれぞれ個別指定し、同じ版別`make compare`で全件を比較する。再開時は`CachingIterator`の`setFlags()`・ArrayAccessの診断と反復中変更を小さなケースで検査し、固定アプリで必要な差異を補う。Podman接続回復後は新規5件と全1129件を版別比較し、固定アプリのCLI・HTTP・DB・API契約の公式PHPとの差も確認する。

## 2026-10-21：AppendIteratorとFinderの複数探索先

`AppendIterator`の基本操作（`append()`、`rewind()`、`valid()`、`key()`、`current()`、`next()`、`getInnerIterator()`、`getIteratorIndex()`、`getArrayIterator()`）をIteratorの共通呼び出し経路へ接続した。複数の内部イテレータを宣言順に走査し、各キーを保持する。Symfony Finderの複数`in()`先を実行して、`RecursiveDirectoryIterator::getChildren()`がサブクラスの型付きプロパティ初期値を失う不具合を発見した。子を具体クラスの通常コンストラクタ経路で生成するよう修正し、[Finderの複数探索先](../tests/fixtures/success/finder_multi_vendor.php)と[子サブクラスの初期値](../tests/fixtures/success/spl_child_subclass_properties.php)を回帰fixtureにした。[Finderの初回停止ログ](validation/full-stack-finder-multi-local.log)は修正前の履歴である。

新規5件（AppendIterator 3件と上記2件）で当時1124件（正常700・実行時290・構文134）。[Goの3設定での新規対象15実行](validation/full-stack-append-fixtures-local.log)とローカル[test](validation/full-stack-append-test.log)・[vet](validation/full-stack-append-vet.log)・[build](validation/full-stack-append-build.log)・[差分](validation/full-stack-append-diff.log)、変更後の[`artisan list`](validation/full-stack-append-artisan-list.log)を記録した。固定vendorのFinder複数探索先はGoで3ファイルを列挙した。PHP公式3版との比較はPodman接続拒否が続いているため今回未試行で、追加127件の期待値は暫定である。各版とも今回0件比較実施、0件の確認済み差異、1124件の変更後全件比較は未実施。最後の公式比較済み基準は997件のまま。`AppendIterator`の反復中`append()`や全診断境界、SPL子クラスのコンストラクタ例外は未検証である。

Macで接続回復後は、リポジトリルートから対象5件へ`python3 skills/php-go-interpreter/scripts/compare_php.py --php ./containers/php-version/bin/php84 --php-version 8.4.26 --fixtures tests/fixtures/success/append_iterator_basic.php --interpreter ./bin/gophp --php-version=8.4 '{file}'`を実行する。PHP 8.3.35・8.5.11では`php83`／`8.3.35`／`8.3`、`php85`／`8.5.11`／`8.5`へ対応する3箇所を替え、残り4件のファイル名も個別に指定する。全件は版ごとに`make compare PHP=./containers/php-version/bin/php84 PHP_VERSION=8.4.26`（8.3／8.5も対応する版名へ変更）を実行し、個別結果・集計・実際の終了コードを版別生ログに保存する。この時点で未実装だった`CachingIterator`は後の節で基本形を追加し、固定vendor fixtureでもGoローカル実行済みである。次はSPL境界と固定アプリのCLI・HTTP・DB・API経路の不足を依存順に補う。

## 2026-10-21：Iteratorラッパーと固定vendorの配列ページング

固定vendorの`ApiPlatform\\State\\Pagination\\ArrayPaginator`は`LimitIterator`と`iterator_count()`を使うことを実コードで確認した。Go側に`IteratorIterator`と`LimitIterator`の基本ラッパーを追加し、既存の`Iterator`／`IteratorAggregate`呼び出し経路から内側のイテレータへ委譲する。`LimitIterator`はoffset・limitを適用し、キーを再採番しない。`iterator_count()`は配列とTraversableを受け、実際に反復して件数を数える。公式PHPの[LimitIterator](https://www.php.net/manual/en/class.limititerator.php)・[IteratorIterator](https://www.php.net/manual/en/class.iteratoriterator.php)・[iterator_count()](https://www.php.net/manual/en/function.iterator-count.php)を参照したが、実比較は未実施である。`IteratorIterator`の第2引数での独自クラス指定、ラッパーの反復中変更・全`seek()`境界・キャッシュ挙動は未対応／未検証で、独自クラス指定はGo固有のエラーとして拒否する。

新規6 fixture（IteratorIterator基本・IteratorAggregate経由、LimitIterator基本・不正値、iterator_count、固定vendorのArrayPaginator）を追加した。当時1119件（正常695・実行時290・構文134）のローカル結果は[新規対象](validation/full-stack-iterator-wrappers-fixtures-local.log)、[test](validation/full-stack-iterator-wrappers-test.log)・[vet](validation/full-stack-iterator-wrappers-vet.log)・[build](validation/full-stack-iterator-wrappers-build.log)・[差分](validation/full-stack-iterator-wrappers-diff.log)へ記録する。変更後の[`artisan list`](validation/full-stack-iterator-wrappers-artisan-list.log)も終了0。固定vendorのArrayPaginatorを一時スクリプトでも直接実行し、`1:20;|1|3`・終了0を[生ログ](validation/full-stack-array-paginator-local-after.log)で確認した。これはGoローカル証拠であり、PHP3版の互換証拠ではない。最後の公式比較済みは各997件のまま。Podmanは同一の[接続失敗](validation/full-stack-object-foreach-podman.log)から状態変化がないため再試行せず、各版0件比較実施・0件不一致・1119件未実施。Macで接続回復後は次の対象6件も前節の版別対象比較ループへ加え、版別`make compare`で全件を確認する。 当時の次の実装単位はSymfony Finderの複数探索先で使う`AppendIterator`と、今回のIteratorラッパーの反復中変更・診断境界の検証である。これと並行して固定アプリのHTTP/API契約を広げ、Podman接続回復後は1119件の版別公式比較を行う。

```zsh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  read -r tag version target <<< "$spec"
  for fixture in tests/fixtures/success/{iterator_wrapper_basic,iterator_wrapper_aggregate,limit_iterator_invalid,iterator_count_basic,array_paginator_vendor}.php tests/fixtures/runtime/limit_iterator_basic.php; do
    log="docs/validation/full-stack-wrapper-target-$(basename "$fixture" .php)-${tag}-mac.log"
    python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php${tag}" --php-version "$version" --fixtures "$fixture" --interpreter ./bin/gophp "--php-version=${target}" '{file}' 2>&1 | tee "$log"
    check_rc=$?; printf 'exit_code: %s\n' "$check_rc" >> "$log"
  done
  log="docs/validation/full-stack-wrapper-all-${tag}-mac.log"
  make compare "PHP=./containers/php-version/bin/php${tag}" "PHP_VERSION=${version}" 2>&1 | tee "$log"
  check_rc=$?; printf 'exit_code: %s\n' "$check_rc" >> "$log"
done
```

## 2026-10-21：通常オブジェクトのforeach

[PHPのオブジェクト反復](https://www.php.net/manual/en/language.oop5.iterations.php)で規定される可視プロパティの基本形を実装した。`Iterator`／`IteratorAggregate`／Generatorは従来の反復経路を優先する。通常オブジェクトでは格納済みのプロパティを宣言順に調べ、実行中のクラススコープで見えないprivate／protectedスロットと未初期化の型付きプロパティを除く。privateスロットの内部キーは反復キーへ漏らさない。値渡しの反復は既存の配列反復処理へ接続し、参照反復は既存の参照格納先・型制約を通す。`stdClass`の動的プロパティも対象。反復中の追加・削除・値変更、参照Iterator一般、マジックプロパティの挙動は未検証で、対応完了とは扱わない。

新規4 fixture（`foreach_object_visibility.php`、`foreach_object_private_slots.php`、`foreach_object_dynamic.php`、`foreach_object_reference.php`）を追加し、当時1113件（正常689・実行時290・構文134）となった。[Goの4設定での対象結果](validation/full-stack-object-foreach-fixtures-local.log)、[test](validation/full-stack-object-foreach-test.log)・[vet](validation/full-stack-object-foreach-vet.log)・[build](validation/full-stack-object-foreach-build.log)・[差分検査](validation/full-stack-object-foreach-diff.log)、変更後の[`artisan list`](validation/full-stack-object-foreach-artisan-list.log)はいずれも終了0。公式PHPの実比較は[Podman接続失敗の記録](validation/full-stack-object-foreach-podman.log)後には実施していない。8.3.35・8.4.26・8.5.11は各0件比較実施、0件不一致、1113件未実施で、Podman終了125はfixtureエラーには含めない。新規期待値は公式マニュアルに基づく暫定値。Macで接続が回復したら、まず4件を対象比較し、その後、既存の版別`make compare`で全件を比較する。版・Go設定をそろえるコマンドは下記のとおり。

```zsh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  read -r tag version target <<< "$spec"
  for fixture in tests/fixtures/success/foreach_object_{visibility,private_slots,dynamic,reference}.php; do
    log="docs/validation/full-stack-object-target-$(basename "$fixture" .php)-${tag}-mac.log"
    python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php${tag}" --php-version "$version" --fixtures "$fixture" --interpreter ./bin/gophp "--php-version=${target}" '{file}' 2>&1 | tee "$log"
    check_rc=$?; printf 'exit_code: %s\n' "$check_rc" >> "$log"
  done
  log="docs/validation/full-stack-object-all-${tag}-mac.log"
  make compare "PHP=./containers/php-version/bin/php${tag}" "PHP_VERSION=${version}" 2>&1 | tee "$log"
  check_rc=$?; printf 'exit_code: %s\n' "$check_rc" >> "$log"
done
```

## 2026-10-21：JSON-LDコレクションと長文xxh3

ページングの件数だけを検査していた段階では、`application/ld+json`の一覧がJSON配列になり、`@type=Collection`・`totalItems`・`member`が欠けていた。固定vendorのSymfony Serializerは`is_subclass_of($type, '\\Traversable')`で正規化器を選ぶが、Goの`is_subclass_of()`はクラス名文字列の先頭`\\`を除いていなかった。これを修正するとJSON-LD正規化へ進み、続いてAPI Platformのメタデータキャッシュキーに使う240バイト超の`hash('xxh3')`で停止した。xxHash公式実装の長文経路では64ビット出力と128ビット出力の下位64ビットが同じ累積値から生成されるため、既存の長文`xxh128`計算を64ビット側へ接続した。これは[xxHash公式ソース](https://github.com/Cyan4973/xxHash/blob/dev/xxhash.h)に基づく実装で、PHP3版の実比較は未実施である。

[修正後のGo統合ログ](validation/full-stack-pagination-jsonld.log)はGoの8.3・8.4・8.5設定それぞれで独立した31行の一時SQLiteを使い、3ページすべてでHTTP 200・診断なし、JSON-LDの`Collection`・`totalItems=31`・`member`が30／1／0件になることを確認する。先頭`\\`を持つクラス・インターフェース名を[`is_a`／`is_subclass_of`のfixture](../tests/fixtures/success/is_subclass_of_fqn.php)で、241・256・1024・1025バイトの長文ハッシュを[別のfixture](../tests/fixtures/success/hash_xxh3_long.php)で確認する。加えて[xxHash公式テストベクトル](https://github.com/Cyan4973/xxHash/blob/dev/cli/xsum_sanity_check.c)の長さ403・512・2048・2099・2240・2367バイトを[Go単体テスト](../internal/php/hash_xxh3_test.go)で照合した。[新規2件のGo3設定ローカル結果](validation/full-stack-jsonld-fixtures-local.log)は各設定2件とも暫定期待値に一致した。fixture期待値は公式PHPの実出力による確定前の暫定値である。[変更後のtest](validation/full-stack-jsonld-test.log)・[vet](validation/full-stack-jsonld-vet.log)・[build](validation/full-stack-jsonld-build.log)・[差分](validation/full-stack-jsonld-diff.log)を記録する。当時1109件の全件・新規2件は3版とも未比較、環境エラーは前回のPodman接続拒否として別記録する。Macでは下記の版別コマンドを用い、対象比較の`--fixtures`を`tests/fixtures/success/is_subclass_of_fqn.php`と`tests/fixtures/success/hash_xxh3_long.php`の各ファイルへ置き換えて実行する。全件は同じ版指定の`make compare`を使い、各コマンドの終了コードを直後に保存する。

## 2026-10-21：内部関数コールバックの型モードとAPIページング

31行を持つ一時SQLiteで`GET /api/books?page=1`を実行すると、以前のGo製PHPは500となった。[失敗ログ](validation/full-stack-pagination-probe.log)の最初のTypeErrorは固定vendorの`ApiPlatform\State\Util\RequestParser::parseRequestParams()`にある`array_map(static fn (string $key) => ..., array_keys($params))`で発生した。`page`の16進表現が数値キーになるため、コールバックにはintが渡る。Goは内部関数`array_map()`からの呼び出しにも呼び出し元の`strict_types=1`を適用していた。[PHPの型宣言マニュアル](https://www.php.net/manual/en/language.types.declarations.php)では内部関数からの呼び出しはその宣言の影響を受けない。内部関数が行うコールバック呼び出しへ共通の弱い型モードを追加し、array系・`preg_replace_callback()`・`call_user_func()`系の既存経路を接続した。戻り値の型検査はコールバックの宣言元で行う既存処理を維持する。

[修正後の固定アプリ対象ログ](validation/full-stack-pagination-after-callback.log)は1ページ目30件、2ページ目1件、3ページ目0件の200応答と、31件すべての重複しないID・空でないtitleを記録する。この時点ではJSON-LD Collection構造を検査していなかった。[小さなstrict型の回帰fixture](../tests/fixtures/runtime/internal_callback_weak_strict.php)はGoの3設定で[ローカル期待値と一致](validation/full-stack-pagination-callback-local.log)した。これは公式PHPの実比較前の暫定期待値である。変更後の[test](validation/full-stack-pagination-callback-test.log)・[vet](validation/full-stack-pagination-callback-vet.log)・[build](validation/full-stack-pagination-callback-build.log)・[差分検査](validation/full-stack-pagination-callback-diff.log)、[`artisan list`の再実行](validation/full-stack-pagination-callback-artisan-list.log)を別に記録する。当時1107件（正常683・実行時290・構文134）のうち、最後の公式比較済み基準は997件。追加110件と変更後全件の3版比較は未実施であり、Goのローカルページング成功を公式互換の証拠として扱わない。

## 2026-10-21：固定アプリの512リクエスト継続検証

固定アプリのHTTP統合テストに任意の`GOPHP_HTTP_STRESS_REQUESTS`を追加した。既定の96回は維持し、指定時のみ最大2000回まで同じhandlerと一時SQLiteへ繰り返しGETする。各応答の本文・ステータス・診断を検査し、直前のCRUDで作った行も定期的に読み取る。GC後の保持ヒープが開始前から64 MiBを超えて増えた場合は失敗させる。これは緩い漏れ検出であり、全メモリ寿命や本番負荷の保証ではない。

[内部コールバック修正後の512回対象実行ログ](validation/full-stack-http-512-after-callback.log)はコマンドと終了0、追加GETの所要時間3分20秒、保持ヒープ934768→720384バイトを記録する。[修正前の同テスト](validation/full-stack-http-512-stress.log)も終了0だが、現在の実装の証拠には修正後のログを使う。テスト全体は既存の作成・取得・無効入力・API Platform CRUD・仕様JSON・並列8 GETも通過した。固定アプリの経路をGoでローカル実行した証拠で、公式PHP3版の当時の1106 fixture比較や長時間運用検証を置き換えない。`make test`・`make vet`・`make build`・`git diff --check`の変更後結果は[検証ログ](validation/full-stack-http-512-test.log)等に保存する。

## 2026-10-21：参照Generator内の委譲宣言検査

[PHP公式バグ記録 #71252](https://bugs.php.net/bug.php?id=71252)は、`yield from`を参照Generator内で使うことがPHP 7.0.4以降禁止されたと記す。Go側は従来、式の実行時に独自エラーで拒否していた。呼び出されない関数内でも宣言の事前検査で`Fatal error: Cannot use "yield from" inside a by-reference generator`とするよう変更した。入れ子の通常クロージャの`yield from`は外側の参照Generatorに含めず受理する。追加2 fixtureのGo3設定での出力は[対象ログ](validation/full-stack-reference-yield-from-local.log)、共有基盤の検証は[test](validation/full-stack-reference-yield-from-test.log)・[vet](validation/full-stack-reference-yield-from-vet.log)・[build](validation/full-stack-reference-yield-from-build.log)・[差分検査](validation/full-stack-reference-yield-from-diff.log)、固定アプリCLIの再実行は[`artisan list`ログ](validation/full-stack-reference-yield-from-artisan-list.log)に記録する。当時1106件（正常682・実行時290・構文134）の公式PHP3版比較は未実施で、各版0一致・0不一致・0 fixture環境エラー・1106件未実施。Podmanの接続失敗は別の環境エラーとして[前回の試行ログ](validation/full-stack-generator-reference-podman.log)に記録済み。版別診断・期待値は暫定。残る参照Generatorの境界は、変数以外を`yield`したときの診断と反復中変更である。

## 2026-10-21：通常関数・メソッド・Closureの参照戻り値（1104件時点）

通常関数、インスタンス／staticメソッド、Closure、アロー関数の`&`戻り値を既存の参照セルへ接続した。`$alias =& source()`はセルへ束縛し、通常の`$copy = source()`は値として読む。参照を返す別関数への委譲、配列要素、プロパティ、staticローカル変数を対象にした。動的Closure呼び出しは通常評価と参照取得で共通の呼び出し処理を使い、通常評価時に内部セルを露出させない。[PHPの参照戻り値](https://www.php.net/manual/en/language.references.return.php)と[アロー関数](https://www.php.net/manual/en/functions.arrow.php)の資料を確認したが、今回の期待値は公式PHP実比較前の暫定値である。

[7 fixture・Goの3設定での実行](validation/full-stack-reference-return-local.log)、[test](validation/full-stack-reference-return-test.log)・[vet](validation/full-stack-reference-return-vet.log)・[build](validation/full-stack-reference-return-build.log)、[`artisan list`の再実行](validation/full-stack-reference-return-artisan-list.log)を記録した。当時1104件（正常681・実行時290・構文133）のうち、最後に公式比較したのは997件で、追加107件と変更後全件の3版比較は未実施。Podmanの最後の接続失敗は[終了125のログ](validation/full-stack-generator-reference-podman.log)に記録済みで、状態変化の情報がないため今回は再試行していない。参照を返す呼び出しを参照引数へ渡す基本形は一度だけ評価し、別名として渡す。非変数式と括弧の参照戻り値は未検証・未対応の境界として残す。

Macでの対象比較は、リポジトリルートから各版について例えば `python3 skills/php-go-interpreter/scripts/compare_php.py --php ./containers/php-version/bin/php84 --php-version 8.4.26 --fixtures tests/fixtures/success/reference_return_closure.php --interpreter ./bin/gophp --php-version=8.4 '{file}'` と実行する。全件は `make compare PHP=./containers/php-version/bin/php84 PHP_VERSION=8.4.26` とする。8.3と8.5ではそれぞれ `php83`／`8.3.35`、`php85`／`8.5.11` を指定し、Go設定も `--php-version=8.3`／`8.5` にする。版別ログにはコマンド、個別結果、集計、コマンド直後の実際の終了コードを保存する。

## 2026-10-21：Generatorの参照yield（参照戻り値追加前の履歴）

[PHPのGenerator構文](https://www.php.net/manual/en/language.generators.syntax.php)にある`function &name()`の基本形を既存の参照セルへ接続した。変数・配列要素・通常プロパティを`yield`し、`foreach (... as &$value)`から共有値を変更できる。`IteratorAggregate::getIterator()`がこのGeneratorを返す場合も通る。値渡しの`foreach`と`Generator::current()`・`send()`・`throw()`の結果は参照セルをPHPの値として露出させず、値を読み取る。参照yieldを宣言していないGeneratorへの参照foreachは開始前に拒否する。参照戻り値を持つ通常関数は引き続き未対応であり、参照Generatorでの`yield from`と非変数式の参照yieldもGo固有の未対応エラーで拒否する。これらをPHP本来の診断と混同しない。

[追加5件のGo3設定での実行ログ](validation/full-stack-generator-reference-local.log)は各出力・終了0・ローカル期待値との一致を記録する。[test](validation/full-stack-generator-reference-test.log)・[vet](validation/full-stack-generator-reference-vet.log)・[build](validation/full-stack-generator-reference-build.log)・[差分検査](validation/full-stack-generator-reference-diff.log)・[fixture race検査](validation/full-stack-generator-reference-race.log)も終了0。変更後の[`artisan list`](validation/full-stack-generator-reference-artisan-list.log)はコマンド一覧と終了0を記録した。[Podman接続の再確認](validation/full-stack-generator-reference-podman.log)は終了125の環境エラーであり、公式PHP3版の対象・全件比較は未実施。今回の5件を含む追加100件の期待値は暫定で、収録1097件を公式一致済みとは扱わない。次は参照Generatorの委譲と参照戻り値、固定アプリの長時間継続・公式比較を進める。

## 2026-10-20 の継続実装と次の検証点

Generator停止点を越えた後、Finder、`xxh128`、Reflection、ArrayObject、PCRE2への限定的フォールバック、日時API、SQLiteなど実vendor経路で不足した機能を追加した。[`artisan list`の生ログ](validation/full-stack-http-api-artisan-list.log)はコマンド一覧と`exit_code: 0`を記録する。[OpenAPI export](validation/full-stack-http-api-openapi-export.log)は新しい一時storageでJSONを生成し終了0。[固定アプリのHTTP統合テスト](../cmd/gophp/main_test.go)は一時SQLiteを用い、probe、手動の作成・取得・422、API Platformの作成・取得・一覧・PATCH・DELETE・削除後404・OpenAPI JSON取得を一連の別リクエストで確認した。ここはGo製PHPのローカル実行証拠で、公式PHP3版との比較結果ではない。初期契約の基本経路はローカルで通ったが、全6段階の完成条件は未達。

DELETE経路で、前のHTTPリクエストから残るSQLite SELECTカーソルが`database is locked`を起こした。各実行が所有するPDOステートメントをfinalizeしてから接続を閉じるよう修正し、連続リクエストの統合テストで再発しないことを確認した。Symfony Responseの空本文処理で必要な`ini_set('default_mimetype', '')`、リクエスト形式の判定に必要な`strrchr()`を追加した。`ini_set()`の対応項目は`default_mimetype`と`include_path`に限定し、他は成功を偽らずfalseを返す。PCRE2とSQLiteはcgo・実行環境に依存する。`xxh128`の長文入力、Reflection捕捉状態、request_parse_bodyの全Content-Type、Generatorの`send()`・`throw()`基本形と委譲中の転送を追加した。Generatorの全制御、継続実行の状態漏れと負荷は未検証・未実装の境界が残る。

[ローカルtest](validation/full-stack-enum-unserialize-test.log)・[vet](validation/full-stack-enum-unserialize-vet.log)・[build](validation/full-stack-enum-unserialize-build.log)・[差分検査](validation/full-stack-enum-unserialize-diff.log)は実コマンドと終了0を記録する。新規100件と現行1097件の公式PHP比較は3版とも未実施（各版0一致・0不一致・0 fixture環境エラー・1097件未実施）。Podmanの接続拒否は別の環境エラーとして記録し、環境が変わらないため繰り返し試していない。Macでの版別対象・全件比較は下記の既存`compare_php.py`と`make compare`の手順を使い、対象には`enum_serialization.php`・`enum_unserialize_failure.php`・`get_debug_type_basic.php`、`http_*.php`、`generator_*.php`を含める。

Generatorの`valid()`・`current()`・`key()`で未開始の本体を最初の`yield`まで進め、`next()`はそこから次の地点へ進める。最初の`yield`にいる間の`rewind()`は許可し、それを越えた後は`Exception`を投げる。`getReturn()`も未開始の本体を開始し、戻り値がまだ確定していなければ`Exception`を投げる。これは[PHP 8.4のGenerator実装](https://github.com/php/php-src/blob/PHP-8.4/Zend/zend_generators.c)の制御分岐と[公式rewind説明](https://www.php.net/manual/en/generator.rewind.php)を照合した。新規[制御fixture](../tests/fixtures/success/generator_control_methods.php)・[初回呼出fixture](../tests/fixtures/success/generator_initial_methods.php)は[Goの3設定でローカル期待値に一致](validation/full-stack-generator-control-local.log)。[test](validation/full-stack-generator-control-test.log)・[vet](validation/full-stack-generator-control-vet.log)・[build](validation/full-stack-generator-control-build.log)・[差分検査](validation/full-stack-generator-control-diff.log)は各コマンド終了0。変更後の[`artisan list`再実行](validation/full-stack-generator-control-artisan-list.log)もコマンド一覧と終了0を記録した。公式PHPの実行比較は未実施で、期待値は暫定。参照yieldや全Generator制御を完了扱いにしない。

Generatorの`yield from`について、内側のGeneratorへ外側の`send()`値・`throw()`例外を転送し、既存の`Iterator`／`IteratorAggregate`を使う委譲も基本形を追加した。`IteratorAggregate::getIterator()`がGeneratorを返す場合も接続した。[4件のGoローカル結果](validation/full-stack-generator-delegation-local.log)は3設定で期待値と一致。[test](validation/full-stack-generator-delegation-test.log)・[vet](validation/full-stack-generator-delegation-vet.log)・[build](validation/full-stack-generator-delegation-build.log)・[差分検査](validation/full-stack-generator-delegation-diff.log)は終了0。新規4件の期待値は[PHPのGenerator委譲の説明](https://www.php.net/manual/en/language.generators.syntax.php#control-structures.yield.from)に基づく暫定値で、公式PHP実比較待ち。参照yieldとGeneratorの全制御は引き続き未対応・未検証である。全件テスト中に`stream_select()`のシステムコールが一度中断され、既存fixtureに警告が出たため、待機時間を保って`EINTR`時に再試行するよう修正した。対象fixtureの再実行と全件テストは成功したが、この経路も公式PHP比較待ちである。

固定vendorの属性経路を直接実行したところ、PHP 8.5設定だけ`vendor/symfony/polyfill-uuid/bootstrap85.php`の属性付きグローバル定数で停止した。PHP 8.5の[定数属性](https://www.php.net/deprecated)を構文解析・定数メタデータへ追加し、`#[Deprecated]`定数の読み取りを警告ハンドラへ接続した。続くLaravel Container内の`array_last()`未定義も、PHP 8.5限定の[array_first()](https://www.php.net/manual/en/function.array-first.php)・[array_last()](https://www.php.net/manual/en/function.array-last.php)を追加して解消した。[固定vendorの属性DI・リソース発見テスト](validation/full-stack-fixed-vendor-attribute-di.log)はGoの3設定で終了0。[8.5設定のartisan list](validation/full-stack-php85-artisan-list.log)も実際のコマンド一覧を表示して終了0。新規2件を含むローカル[test](validation/full-stack-php85-vendor-test.log)・[vet](validation/full-stack-php85-vendor-vet.log)・[build](validation/full-stack-php85-vendor-build.log)・[差分検査](validation/full-stack-php85-vendor-diff.log)は終了0。定数属性の完全なターゲット検査・ReflectionConstant・`Deprecated`メッセージの細部は未対応で、新規期待値と変更後全件の公式比較は未実施。

同じHTTPハンドラと一時SQLiteで追加8回のPOST/PATCH/DELETEを実行すると、当初2回目のPOSTが「stringを要求しstringを受けた」という型判定失敗になった。原因はGoの`serialize()`がenum caseを通常の`O`オブジェクトとして記録し、vendorキャッシュからの`unserialize()`後にcaseの単一性を失ったこと。公式PHPの[enumシリアライズ仕様](https://www.php.net/manual/en/language.enumerations.serialization.php)に沿って`E`形式で記録し、既存caseへの復元を実装した。[繰り返しCRUDを含むローカルtest](validation/full-stack-serialization-test.log)は終了0、[新規fixtureのGo3設定](validation/full-stack-serialization-local-versions.log)も終了0。`get_debug_type()`も[公式の型名](https://www.php.net/manual/en/function.get-debug-type.php)に沿って基本形を追加した。新規2件の公式PHP3版比較は未実施で、期待値は暫定。続いて未知のenum/caseや非enumを復元する際の警告・false返却を追加した。警告ハンドラが例外を投げたら追加警告や後続処理へ進まない。[Goの3設定での対象実行](validation/full-stack-enum-unserialize-local.log)と[test](validation/full-stack-enum-unserialize-test.log)は終了0。警告数と期待値は[PHP 8.3.35の公式ソース](https://github.com/php/php-src/blob/php-8.3.35/ext/standard/var_unserializer.re)および[enumシリアライズのマニュアル](https://www.php.net/manual/en/language.enumerations.serialization.php)に基づく暫定値で、8.4・8.5を含む実比較は未実施。一般の参照タグ、カスタムシリアライズ、破損データの全診断は未対応の境界として残る。

固定アプリの`routes/console.php`をLaravelの`withRouting(commands: ...)`で登録し、独自`gophp:probe`コマンドを追加した。[CLI生ログ](validation/full-stack-custom-artisan-command.log)は`probe-ok`・終了0を記録し、Go統合テストでは8.3・8.4・8.5設定の各々で同じ出力を確認した。固定lock・vendorと既存HTTP契約は維持した。公式PHP3版でのこの独自コマンドの実比較は未実施。

[継続リクエストのGoローカル検証](validation/full-stack-sustained-http-test.log)では、同じHTTPハンドラ・一時SQLiteで8回の追加CRUD後に96回のGET（probeと保存済みBook取得）と8並列GETを実行し、すべて200・期待本文・診断なしを確認した。通常実行での96回の測定値は31.91秒、GC後の確保ヒープは開始766,576バイト、終了745,992バイト。[race detector付きの同じテスト](validation/full-stack-sustained-race.log)も終了0（96回の測定値は3分34.91秒）。これらの短い観測は長時間負荷やメモリ上限の保証ではない。[Go3設定の固定アプリHTTP起動](validation/full-stack-http-versions-local.log)と、前リクエストのステータス・ヘッダーが残らない単体テストも通過した。[全ローカルtest](validation/full-stack-sustained-test.log)・[vet](validation/full-stack-sustained-vet.log)・[build](validation/full-stack-sustained-build.log)・[差分検査](validation/full-stack-sustained-diff.log)は終了0。固定アプリの版別公式比較は依然未実施である。

次の実装・検証単位は、固定アプリのHTTP契約を公式PHPの同じ入力で比較し、より長い継続実行の状態分離・PDOカーソル寿命・未対応Content-Typeを確認すること。公式比較環境が戻るまでは、Generator/Reflection/標準関数の既知の限定対応を拡張し、独立して進められるローカル統合検証を続ける。以下の古い停止点は当時の履歴である。

`artisan about`の起動経路で`phpversion()`、実子プロセスとpipeを使う`proc_open()`・`proc_get_status()`・`proc_close()`、`php://temp`、`stream_select()`、出力バッファ、PHPとHTMLの混在テンプレートを追加した。さらにTermwindが使う`DOMDocument::loadHTML()`・`getElementsByTagName()`・`saveXML()`、ノードの属性・兄弟・子の基本走査、`strip_tags()`・`html_entity_decode()`・`mb_strimwidth()`、PCRE2を使う`preg_split()`の先読みを接続した。[`artisan about`](validation/full-stack-http-api-artisan-about.log)は環境・ドライバ情報を出力して終了0。続く[`artisan route:list`](validation/full-stack-http-api-artisan-route-list.log)では`hexdec()`を追加し、8経路の一覧を出力して終了0。いずれもログにLaravelの整形エラーがないことを確認した。DOMはTermwindが通るHTML断片を処理する限定実装で、完全なHTML5解析やDOM APIではない。一般のprocess descriptor、signal、ストリームオプション、HTMLテンプレート境界、マルチバイト文字幅の全Unicode境界は未検証である。

一時SQLiteを指定した[`artisan migrate --force`](validation/full-stack-http-api-artisan-migrate.log)で、移行表と固定アプリの`books`表を作成できた。初回試行で`number_format()`不足を補い、移行ファイルを再読込した匿名クラスが新しいASTの未登録クラス情報を使って親メソッドを失う不具合を修正した。同じソース位置の登録済み匿名クラスを再利用し、再includeの回帰fixtureを追加した。CLIの`migrate:status`はDB設定なしではTypeErrorになり、`:memory:`では移行表なしと表示するため、移行コマンドの成功証拠とは扱わない。

固定アプリのHTTPブリッジへクエリ付きGETを12回続ける回帰検証を追加した。最初はSymfony HttpFoundationの`urldecode()`・`parse_str()`・`hex2bin()`・`http_build_query()`が不足していたため失敗し、これらの基本形とURLエンコードを追加した後、既存の手動/API Platform CRUDを含む統合テスト全体が成功した。値の保持・ヘッダー・出力の状態分離を複数リクエストで確認した証拠だが、長時間負荷やメモリ上限の達成証拠ではない。`parse_str()`の参照出力は位置引数の基本形を実装し、名前付き引数や入力上限等は未検証である。

段階2の独立した検証では、`ReflectionParameter::getAttributes()`が常に空配列を返していたことを確認した。宣言済みパラメータ属性を既存のReflectionAttribute処理へ接続し、自作属性の`newInstance()`からコンストラクタ引数を作る小fixtureがGoで成功した。固定vendorの完全なDI実行経路や属性の全フィルタ規則は未検証である。

## 以前のSPL起動停止点（履歴）

当時、FinderのSPLディレクトリ反復、設定配列の操作、参照foreachの式対象、プロパティの参照引数、HTTP Request入口の関数、ReflectionFunction／ReflectionMethodの一部、属性の名前付き引数構文を追加した。[SPLと配列の回帰fixture](../tests/fixtures/success/generator_spl_directory.php)など3件を収録し、当時の1055件のGoローカル[test](validation/full-stack-spl-test.log)・[vet](validation/full-stack-spl-vet.log)・[build](validation/full-stack-spl-build.log)・[差分検査](validation/full-stack-spl-diff.log)は終了0。`artisan list`はFinderの旧停止点を越えてAPI Platformのリソース探索まで進み、当時は`hash('xxh3')`を要求する。追加モジュール取得はDNS制限で失敗した。異なるハッシュで代用せず、xxh3の正しい実装とその後の起動経路を当時の再開地点とした。RecursiveIteratorIteratorは現在の配列スナップショット型の基本走査のみで、反復中の変更・`getSubIterator()`の位置等は未検証。PCREの先読みはFinderの先頭ドット判定に限った実装であり、一般PCRE互換ではない。`array_walk()`の参照コールバック、`extract()`のEXTR_REFSも未対応として明示的に拒否する。

## Generator後の継続実装（前回の停止点の履歴）

`yield`・`yield from`を遅延実行するGenerator、配列とTraversableの委譲、`getReturn()`、終了時の`finally`をGo側へ追加し、`artisan list`を繰り返し実行した。起動経路で見つかった不足に対し、動的メソッド名、属性構文、`eval`、一部のReflection、DateTime、WeakMap、PCRE系関数、ファイル／ストリーム関数、Unicodeエスケープ、SPL例外、FilterIteratorを段階的に追加した。現在の[生の試行ログ](validation/full-stack-generator-artisan-latest.log)は、Symfony Finderの`RecursiveDirectoryIterator`が未実装であることを示す。Laravelが例外を整形表示するためプロセス終了コードは0だが、コマンド一覧は出ていない。実行経路で`ReflectionParameter::$name`と継承元privateプロパティへの書き込みの誤警告も発見し、修正した。

新しい[正規表現](../tests/fixtures/success/generator_console_regex.php)、[Unicode](../tests/fixtures/success/generator_unicode_escape.php)、[FilterIterator](../tests/fixtures/success/generator_filter_iterator.php)、[継承元private書込](../tests/fixtures/success/generator_inherited_private_write.php)等のfixture期待値は暫定。PCREはGoのregexpで表せる一部の構文だけを扱い、一般のpossessive量指定子の意味や後読みを完全には再現しない。iconvはUTF-8・ASCII・Latin-1の基本形のみ、DateTimeは限定したコンストラクタ／メソッドのみ、WeakMapはGC寿命を再現しない。`Reflection::getAttributes()`はまだ属性情報を返さない。これらをLaravel全体の対応済み機能とは扱わない。定数初期値が未登録クラスを参照する場合の遅延評価は、Symfony polyfillの循環したクラス定数参照を越えるため追加したが、PHPの全遅延評価規則は未検証。

ローカルの[test](validation/full-stack-generator-test.log)・[vet](validation/full-stack-generator-vet.log)・[build](validation/full-stack-generator-build.log)・[差分検査](validation/full-stack-generator-diff.log)を実行した。公式PHPの直近比較済み証拠は997件のまま。[Podman接続拒否の試行](validation/full-stack-continuation-podman.log)以降、接続状態が変わらないため版別比較を繰り返していない。1052件中、997件の過去の一致は維持するが、55件の追加期待値と変更後の全件比較は未確認である。Macで接続が回復したら、リポジトリルートで次を版ごとに実行し、各コマンドの終了コードを直後に保存する。対象比較では`generator_*.php`を含め、全件比較は`make compare`を用いる。

```zsh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  read -r tag version target <<< "$spec"
  for fixture in tests/fixtures/success/generator_*.php tests/fixtures/success/http_*.php; do
    log="docs/validation/full-stack-target-$(basename "$fixture" .php)-${tag}-mac.log"
    printf 'command: compare_php.py --php ./containers/php-version/bin/php%s --php-version %s --fixtures %s --interpreter ./bin/gophp --php-version=%s {file}\n' "$tag" "$version" "$fixture" "$target" > "$log"
    python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php${tag}" --php-version "$version" --fixtures "$fixture" --interpreter ./bin/gophp "--php-version=${target}" '{file}' 2>&1 | tee -a "$log"
    rc=$?; printf 'exit_code: %s\n' "$rc" >> "$log"
  done
  log="docs/validation/full-stack-all-${tag}-mac.log"
  printf 'command: make compare PHP=./containers/php-version/bin/php%s PHP_VERSION=%s\n' "$tag" "$version" > "$log"
  make compare "PHP=./containers/php-version/bin/php${tag}" "PHP_VERSION=${version}" 2>&1 | tee -a "$log"
  rc=$?; printf 'exit_code: %s\n' "$rc" >> "$log"
done
```

固定vendorの優先経路として、`laravel/framework/src/Illuminate/Container/Container.php`の参照引数、`Illuminate/Container/RewindableGenerator.php`の反復関連、`symfony/console/Application.php`の参照利用をソース検索で確認した。これらの実行経路にGoが到達した証拠ではない。[依存監査](dependency-audit.md)と[アプリ完成条件](laravel-api-target.md)に照らし、次は参照引数を配列要素・プロパティへ広げ、参照戻り値と参照foreachへ進む。

以下は`yield`停止時の履歴。実vendorの入口はGoで繰り返し実行し、[以前の匿名クラス停止点](validation/full-stack-vendor-autoload-probe-after-substr.log)を越え、[現在の`vendor/autoload.php`実行](validation/full-stack-continuation-vendor-autoload.log)は終了0となった。これはComposerのfiles読み込みとautoload登録が通った証拠であり、autoload対象クラスの網羅やLaravel起動の証拠ではない。[`artisan list`実行](validation/full-stack-continuation-artisan.log)は`Illuminate/Container/Container.php:678`の`yield`で終了255。クラス継承・interface・trait探索でComposer登録済みautoloadを呼ぶようにし、PSR-4経由のtraitを単独で読み込めることも確認した。次はGeneratorと反復契約を実装し、この停止点を越える。`$GLOBALS`配列全体はPHP 8.1以降読み取り用コピーであり、現行実装の直下キー読み書き・ネスト操作以外は未検証。[PHP公式の説明](https://www.php.net/manual/en/reserved.variables.globals.php)に従い、完全なスーパーグローバル対応とは扱わない。`$_SERVER`にはCLIの`argv`・`argc`だけを設定し、`$_ENV`は空配列から開始する。HTTP環境変数は未設定。`extension_loaded()`はGo側に拡張を実装していないためfalseを返し、Core/standardを含む公式PHPの拡張一覧を再現しない。追加した標準関数の引数型・境界診断も公式比較待ちである。

今回の実装では匿名クラスの宣言をASTに保持し、式へ到達した時点で既存クラス登録・trait合成・生成経路へ渡す。同一宣言を繰り返し評価した場合は同じ内部クラスを使う。配列分解代入は右辺を一度評価し、左辺を順に既存格納先へ書く。`goto`はラベルを持つ外側の実行ブロックへ制御を戻す基本形で、ループ・switch境界と未定義ラベルの完全なコンパイル時検査は未対応。`__call`・`__callStatic`・`__toString`は基本の呼び出しを実装したが、マジックメソッド署名・可視性や`Stringable`契約の全検査は残る。式結果への添字代入はPHP構文として受理するが、格納先が得られない実行経路はGo側エラーとなり、`ArrayAccess`等の次の基盤へ引き継ぐ。以上の境界を公式PHPと一致済みとは扱わない。[匿名クラス](https://www.php.net/manual/en/language.oop5.anonymous.php)、[配列分解代入](https://www.php.net/manual/en/language.types.array.php)、[マジックメソッド](https://www.php.net/manual/en/language.oop5.magic.php)、[goto](https://www.php.net/manual/en/control-structures.goto.php)の公式説明を参照したが、新規fixtureの版別実比較は未実施である。

現在のGo実装ではクロージャ`use (&$x)`が外側の変数と参照セルを共有する。関数・メソッド・クロージャの`&$parameter`は、既存の変数を直接渡す位置引数・名前付き引数に限り共通の束縛処理へ接続した。`foreach ($array as &$value)`は配列変数に限り、追記された要素も走査する。Iteratorの`rewind`・`valid`・`key`・`current`・`next`、IteratorAggregateの`getIterator`、Countableの`count`を既存のメソッド呼び出しへ接続した。`iterable`型と`is_iterable()`は配列またはTraversable実装を受理する。配列・オブジェクト以外への`foreach`と`strlen(null)`で警告ハンドラが投げた例外を伝播し、後続処理を中断する。版別に固定された`PHP_VERSION_ID`も追加した。配列要素・プロパティ・未定義変数・展開・可変長の参照引数、配列式の参照foreach、一般オブジェクトのpublicプロパティ反復、ArrayAccess、Generatorは未対応であり、対応済みとして扱わない。新規期待値は[PHPの参照引数](https://www.php.net/manual/en/language.references.pass.php)、[無名関数](https://www.php.net/manual/en/functions.anonymous.php)、[foreach](https://www.php.net/manual/en/control-structures.foreach.php)、[Countable](https://www.php.net/manual/en/class.countable.php)、[Iterator](https://www.php.net/manual/en/class.iterator.php)を参照した暫定値で、公式実比較を待つ。

比較環境は今回も`podman ps`が終了125（[ソケット接続拒否](validation/full-stack-continuation-podman.log)）。同じ状態で比較を繰り返さない。[Goの新規7件・3設定のログ](validation/full-stack-continuation-local-versions.log)は公式比較ではない。共有基盤の[test](validation/full-stack-continuation-test.log)・[vet](validation/full-stack-continuation-vet.log)・[build](validation/full-stack-continuation-build.log)・[差分検査](validation/full-stack-continuation-diff.log)は終了0。今回の追加7件を含む新規38件と全1035件について、PHP 8.3.35・8.4.26・8.5.11の公式比較はいずれも未実施である。Macで接続が回復したら、リポジトリルートから次の版別コマンドを順に実行し、stdout・stderr・終了コードを既存比較器で厳密に比較する。`make compare`はtestとbuildも実行する。各コマンドの終了コードは直後に記録する。今回の7件は既存の対象リストに追加する。

```zsh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  read -r tag version target <<< "$spec"
  for fixture in tests/fixtures/success/reference_closure_capture*.php tests/fixtures/success/reference_parameter_*.php tests/fixtures/success/reference_foreach_*.php tests/fixtures/success/countable_basic.php tests/fixtures/success/iterator_*.php tests/fixtures/success/reference_iteration_miniapp.php tests/fixtures/success/iterable_type.php tests/fixtures/success/foreach_handler_throw.php tests/fixtures/success/strlen_handler_throw.php tests/fixtures/success/php_version_id.php tests/fixtures/success/enum_serialization.php tests/fixtures/success/enum_unserialize_failure.php tests/fixtures/success/get_debug_type_basic.php tests/fixtures/success/globals_scope.php tests/fixtures/success/autoload_basic.php tests/fixtures/success/array_cast_basic.php tests/fixtures/success/dirname_basic.php tests/fixtures/success/closure_bind_basic.php tests/fixtures/success/error_suppression.php tests/fixtures/success/function_exists_basic.php tests/fixtures/success/extension_loaded_unknown.php tests/fixtures/success/hex_literal.php tests/fixtures/success/exit_finally.php tests/fixtures/runtime/exit_status.php tests/fixtures/success/server_superglobal_basic.php tests/fixtures/success/class_exists_basic.php tests/fixtures/success/strtr_basic.php tests/fixtures/success/strrpos_basic.php tests/fixtures/success/substr_basic.php; do
    log="docs/validation/full-stack-target-$(basename "$fixture" .php)-${tag}-mac.log"
    python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php${tag}" --php-version "$version" --fixtures "$fixture" --interpreter ./bin/gophp "--php-version=${target}" '{file}' 2>&1 | tee "$log"
    rc=$?; printf 'exit_code: %s\n' "$rc" >> "$log"
  done
  for fixture in tests/fixtures/success/anonymous_class_constructor_trait.php tests/fixtures/success/array_destructuring_basic.php tests/fixtures/success/object_cast_env.php tests/fixtures/success/goto_retry_basic.php tests/fixtures/success/magic_call_basic.php tests/fixtures/success/shebang_file_exists.php tests/fixtures/success/cli_arguments_basic.php; do
    log="docs/validation/full-stack-target-$(basename "$fixture" .php)-${tag}-mac.log"
    python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php${tag}" --php-version "$version" --fixtures "$fixture" --interpreter ./bin/gophp "--php-version=${target}" '{file}' 2>&1 | tee "$log"
    rc=$?; printf 'exit_code: %s\n' "$rc" >> "$log"
  done
  log="docs/validation/full-stack-all-${tag}-mac.log"
  make compare "PHP=./containers/php-version/bin/php${tag}" "PHP_VERSION=${version}" 2>&1 | tee "$log"
  rc=$?; printf 'exit_code: %s\n' "$rc" >> "$log"
done
```

当時の再開地点：Symfony Finderが必要とする`RecursiveDirectoryIterator`・`SplFileInfo`・`RecursiveIteratorIterator`の実装を進め、`artisan list`を再実行する。公式比較待ちは独立した実装の停止条件ではない。未比較・未実行の段階を完了と記録しない。以前の再開地点として`yield`を先に実装し、ContainerのGenerator使用箇所は越えた。同時に参照引数の格納先（配列要素・プロパティ等）、参照戻り値、`ArrayAccess`を進める。Generatorを構文だけ受理して実行成功に見せることはしない。その後`artisan`を再試行し、Reflection・属性・その他のマジックメソッドへ進む。
