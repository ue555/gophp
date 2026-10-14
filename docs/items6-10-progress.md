# 項目6〜10の進捗（2026-10-14）

## 現在の判定：622件が公式PHPの3版すべてで一致

ユーザーがMacで実行した[修正後の全件比較ログ](validation/named-exception-622-compare-all-mac.log)を確認し、元ログとバイト一致で保存した。PHP 8.3.35・8.4.26・8.5.11の各版で**622 matched・0 different・0 errors**、`make compare-all`終了0。fixtureは正常318・実行時230・構文74。605件から追加した17件の公式比較待ちは解消し、名前付き引数、文字列キーの引数展開、例外基盤、レビュー修正3件、束縛失敗時のトレース修正は収録fixtureの範囲で比較済みである。今回は保存済みログを確認・保存したもので、比較を再実行していない。

項目9・10と項目6〜10全体は部分実装のまま。Traversable展開、参照捕捉、残るThrowable API・継承規則、set_error_handler／restore_error_handler／trigger_error、高度な型、Composer autoload、Go製PHPでのLaravel起動は残る。fixture外の境界条件やPHP全体の互換性は未検証である。

## 622件修正前・接続制限時の記録（履歴）

ユーザーのMacでの[622件全件比較ログ](validation/named-exception-622-user-before-fix.log)を元ログからバイト単位で保存した。PHP 8.3.35・8.4.26・8.5.11は各**620 matched・2 different・0 errors**、`make compare-all`終了2。不一致は`runtime/named_duplicate.php`と`runtime/named_unknown.php`だけで、名前付き引数の束縛に失敗した時点では公式PHPのトレースが`#0 {main}`のみだった。Goは呼び出し先フレームを早く作り、余計な`target(...)`行を出していた。

`internal/php/class.go`の`invokeWithCapture`で、束縛をスコープ・フレーム生成前へ移した。両fixtureの期待stderrは、提供ログにある公式PHPの実出力を根拠に訂正し、メッセージ、stdout、終了コードは維持した。関数内で起きるエラーの既存フレームは維持する。新しいGo単体テストは、トップレベルの束縛エラーに呼び出し先フレームがなく、入れ子の呼び出しでは呼び出し元フレームが残ることを確認する。[保存済み公式出力との再照合](validation/named-exception-binding-saved-output-recheck.log)は2件×3設定のstdout・stderr・終了コードが一致（終了0）。これは修正後の公式PHP再実行ではない。

[ローカルテスト](validation/named-exception-binding-local.log)、[vet](validation/named-exception-binding-vet.log)、[ビルド](validation/named-exception-binding-build.log)、[差分検査](validation/named-exception-binding-diff-check.log)は各終了0。[修正後の全件比較試行](validation/named-exception-binding-compare-all-attempt.log)はPodmanソケット接続拒否で`make compare-all`終了2、各版の版確認が終了125。この試行は各版**0一致・0不一致・1環境エラー・622件未実施**であり、fixture差異とは扱わない。Macでの再比較はこの文書の下記手順に従い、`make compare-all`の終了コードをパイプ直後に記録する。

## 622件比較前の記録（履歴）

レビューで確認された3件を修正した。`parser.arguments`は名前付き引数の後にある`...`を未選択分岐でもコンパイル時に拒否する。`try ... finally {}`は空ブロックのnilとfinally節の不在を別に判定する。`catch`変数は`variableSlot`へ書き込み、先に実行された`global`／`static`の束縛を尊重する。根拠は[PHP 8.4.26のコンパイラ](https://github.com/php/php-src/blob/php-8.4.26/Zend/zend_compile.c)、[PHPのtry/catch/finally](https://www.php.net/manual/en/language.exceptions.php)および変数スコープ規則。`parser/named_unpack_after.php`、`parser/named_unpack_skipped.php`、`success/exception_empty_finally.php`、`success/exception_catch_bindings.php`の4件を追加し、現在は**622件**（成功318・実行時230・構文74）。既存618件のfixtureと期待値は維持した。

[レビュー修正後のローカルテスト](validation/named-exception-review-local.log)、[vet](validation/named-exception-review-vet.log)、[ビルド](validation/named-exception-review-build.log)、[差分検査](validation/named-exception-review-diff-check.log)は各終了0。Pythonは12件。[対象比較試行](validation/named-exception-review-targeted-attempt.log)と[全件比較試行](validation/named-exception-review-compare-all-attempt.log)はPodmanソケット拒否で各版の版確認が終了125、比較器・makeとも終了2。今回の比較は各版**0一致・0不一致・1環境エラー・622件未実施**。新規4件の期待値は暫定。Macでは停止中の既存コンテナのみ起動し、次を実行する。終了コードは各パイプの直後に保存する。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/named-exception-review-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/named-exception-review-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
for file in tests/fixtures/parser/named_unpack_after.php tests/fixtures/parser/named_unpack_skipped.php tests/fixtures/success/exception_empty_finally.php tests/fixtures/success/exception_catch_bindings.php; do
  ./containers/php-version/compare-all.sh "$file" ./bin/gophp 2>&1 | tee -a docs/validation/named-exception-review-targeted-mac.log
  target_exit=$?
  echo "$file exit=$target_exit" | tee -a docs/validation/named-exception-review-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/named-exception-review-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/named-exception-review-compare-all-mac.log
```

## 618件時点の記録（履歴）

ユーザーがMacで実行した[605件の成功ログ](validation/diamond-types-fix-compare-all-mac.log)は、PHP 8.3.35・8.4.26・8.5.11で各**605 matched、0 different、0 errors**、`make compare-all exit=0`を記録する。以下の旧605件時点の「3件再比較待ち」はこの成功によって解消した。過去の失敗・Podman接続拒否ログは当時の履歴として維持する。

今回、名前付き引数と文字列キー配列の展開、基本的なthrow/try/catch/finallyを加え、成功9・実行時3・構文1の**13件**を追加した。現在の合計は**618件**（成功316・実行時230・構文72）。既存605件のソース・期待値は変更していない。新規13件の期待値はGoローカル実行とPHP公式資料に基づく**暫定値**で、公式PHP3版でのバイト単位比較は未実施。比較不能をfixture一致とは数えない。

| 項目 | 現在の範囲と残件 |
| --- | --- |
| 6 複数ファイル | ローカル読込と読込先からの基本例外伝播をGoテスト・fixtureで確認。include_path、wrapper、解析失敗・once履歴・ファイル間診断の境界が残る |
| 7 クラス | 従来の部分対応を維持。型付き契約・複雑なtrait競合等が残る |
| 8 型 | 単純な引数・戻り値型とstrict_typesの部分対応。プロパティ／定数型、高度な型、readonly、enum等が残る |
| 9 呼び出し | 名前付き引数を宣言順と独立に束縛し、展開配列の整数キーを位置引数、文字列キーを名前付き引数として扱う。評価はソース順。関数・メソッド・staticメソッド・コンストラクタ・クロージャ・一部callable・対応済み組み込み関数で共通処理を利用。参照捕捉、Traversable展開、エラー直前の副作用順序を含む未収録境界が残る |
| 10 例外 | **部分実装**。throw文／式、try/catch/finally、複数型catch、再throw、Exception/Error/TypeError/ArgumentCountError/ValueError/DivisionByZeroError/UnhandledMatchError、getMessage/getCode/getFile/getLine/getTraceAsStringの基礎を追加。set_error_handler等、完全なThrowable規則、宣言エラーとの境界、未捕捉診断の全形式は残る |

今回の[ローカル検証ログ](validation/named-exception-local.log)はGoとPython 12件を含む`make test`終了0。[vet](validation/named-exception-vet.log)、[ビルド](validation/named-exception-build.log)、[差分検査](validation/named-exception-diff-check.log)も終了0。[対象比較試行](validation/named-exception-targeted-attempt.log)と[全件比較試行](validation/named-exception-compare-all-attempt.log)はPodmanソケット接続拒否により比較器／`make compare-all`終了2。各版の版確認`podman exec`は終了125であり、**各版0一致・0不一致・1環境エラー、618件未実施**。これはPHPプログラムの終了125やfixture差異ではない。Macでは既存コンテナを削除せず、停止中なら起動し、`set -o pipefail`のもとで次を実行する。追加分を対象比較した後、全件を比較する。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/named-exception-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/named-exception-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
for file in tests/fixtures/success/named_*.php tests/fixtures/success/exception_*.php tests/fixtures/runtime/named_*.php tests/fixtures/runtime/exception_uncaught.php tests/fixtures/parser/named_*.php; do
  ./containers/php-version/compare-all.sh "$file" ./bin/gophp 2>&1 | tee -a docs/validation/named-exception-targeted-mac.log
  target_exit=$?
  echo "$file exit=$target_exit" | tee -a docs/validation/named-exception-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/named-exception-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/named-exception-compare-all-mac.log
```

## 605件時点の記録（履歴）

### 当時の判定：605件のGoローカル検証済み、修正した3件は公式再比較待ち

今回の[追加fixture一覧](validation/items6-10-types-new-fixtures.txt)は正常11・実行時4・構文2の計17件。現時点の総数は**605件**（正常307・実行時227・構文71）。比較済みの旧588件のPHPソースと期待値は変更していない。旧588件が公式PHP 8.3.35・8.4.26・8.5.11で各588一致した記録は[直前の成功ログ](validation/item7-588-user-comparison-after-fix.log)に残る。ユーザー提供の[605件比較ログ](validation/argument-unpack-compare-all-mac.log)は各版602一致・2不一致・1検証エラー、`make compare-all`終了2だった。追加17件中、14件はこの実比較で一致した。修正した3件は公式再比較待ちであり、605件全体の比較成功とは扱わない。

失敗した3件は次のとおり。`interface_diamond_arity.php`では親interface間の非互換を、クラスの実装より先に、最初の継承元のメソッド宣言行で検出するよう変更した。`type_strict_arg.php`では提供ログの公式PHPの診断にある`and defined`をユーザー定義関数の引数TypeErrorへ追加した。前二者の修正後Go出力は[保存済み公式出力との照合](validation/diamond-types-fix-saved-output-recheck-final.log)で3版×2件のstdout・stderr・終了コードが一致した。これは公式PHPの再実行ではない。`success/interface_diamond_contracts.php`は第二引数を必須とする旧fixture自体が公式PHPで宣言エラーだったため、第二引数を任意にして両interfaceの契約が両立する成功例へ修正した。旧期待値の訂正は[修正前の公式出力](validation/argument-unpack-compare-all-mac.log)を根拠とし、成功例の新しい宣言内容は公式PHPで再確認が必要である。

今回の作業を依存順に分け、型なしのinterface契約の保持、単純型とファイル単位のstrict_types、動的なクラス名・配列callable、複数ファイルの型付き呼び出し、値捕捉クロージャとarrow関数、可変長パラメータと数値キー配列の引数展開の一部まで実装した。続くレビューで、展開後の通常位置引数を受理していた不具合を修正した。`f(...[1], 2)`と未選択分岐の同形は解析時に拒否し、複数の展開と展開前の位置引数は維持する。[PHP公式RFC](https://wiki.php.net/rfc/argument_unpacking)は展開後の通常引数を禁止している。例外・高度な型・名前付き引数と参照捕捉にはまだ着手できていない。個々のPHP診断と型変換の境界は公式PHP実比較で確定する必要がある。

| 項目 | 今回を含む対応状態 | 残件・次の実装箇所 |
| --- | --- | --- |
| 6 複数ファイル | **部分実装**。既存のローカル読込に加え、読込元と宣言元の`strict_types`を別々に保持し、別ファイル関数の型付き呼出しをGo単体テストで検証 | `include.go`の失敗診断、例外を伴う読込、ファイル間トレース、include_path・stream wrapperの必要範囲 |
| 7 クラス | **部分実装**。複数interfaceから継承した同名メソッドの各契約をクラス実装時に検査。文字列変数による`new`と`instanceof`、オブジェクト／静的メソッドの配列callableを追加 | `class.go`のtrait競合・宣言時点・動的staticメンバー、複合型署名と型付きプロパティ／定数の互換性 |
| 8 型 | **部分実装・公式比較待ち**。`int`/`float`/`string`/`bool`/`array`/`object`/`mixed`/クラスの単純な引数・戻り値、`?T`、`declare(strict_types=0/1)`、基本的な型変換とTypeError、単純なメソッド型互換性 | `types.go`と`class.go`でunion/intersection/DNF、void/never、完全な変換・警告、型付きプロパティ／定数、readonly、昇格、enum、完全な分散検査。未対応構文は構文解析で拒否 |
| 9 呼び出し | **部分実装・公式比較待ち**。文字列変数関数と二要素配列callable、値捕捉クロージャ、arrow関数、staticクロージャ、可変長パラメータ、数値キー配列の引数展開を追加 | `function.go`・`parser.go`・`closure.go`で参照捕捉、名前付き引数と文字列キーの展開、Traversableの展開、一般のcallable、呼出し境界 |
| 10 例外 | **未実装**。TypeErrorは既存の致命エラー診断経路へ接続したが、PHPの`catch`で捕捉できない | `eval.go`の制御フローへthrow/try/catch/finallyを加え、Throwable階層、ハンドラ、再throwとスタックトレースを統合 |

[修正後のGo・Pythonテスト](validation/diamond-types-fix-local-test.log)、[vet](validation/diamond-types-fix-vet.log)、[ビルド](validation/diamond-types-fix-build.log)、[差分検査](validation/diamond-types-fix-diff-check.log)は各終了0。Pythonは12件。[修正後の605件全件比較試行](validation/diamond-types-fix-compare-all-attempt.log)はPodmanのPHPバージョン確認で停止し、`make compare-all`終了2。**この再比較試行は各版0一致・0不一致・1環境エラー・605件未実施**。PHPプログラムの終了125やfixture差分ではない。先行する[新規2件の対象比較試行](validation/argument-unpack-targeted-attempt.log)と[605件の全件比較試行](validation/argument-unpack-compare-all-attempt.log)、[可変長引数の対象比較試行](validation/items6-10-variadic-targeted-attempt.log)も履歴として残す。

固定vendorでは`vendor/composer/ClassLoader.php:80,551`のinclude_path利用フラグは既定falseで、PSR-0のフォールバック時だけ`stream_resolve_include_path`へ進む。Laravelの`config/logging.php:102`には`php://stderr`があるが、初期autoloadでそのログ設定まで到達するかは未確認。したがってinclude_pathとstream wrapperは現段階では未対応とし、起動トレースで必要性を決める。PHP 8.4の非対称可視性・プロパティフックは固定vendorの`.php`を文字列検索した範囲で候補なしだが、構文解析による網羅検査は未実施。8.3設定でこれらを受理する実装はない。

Macでは既存コンテナを削除せず、状態確認後に停止中のものだけ起動する。`tee`直後の終了コードを保存する。追加分を先に比較し、その後に全件比較する。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/diamond-types-fix-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/diamond-types-fix-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
: > docs/validation/diamond-types-fix-targeted-mac.log
for fixture in tests/fixtures/runtime/interface_diamond_arity.php tests/fixtures/runtime/type_strict_arg.php tests/fixtures/success/interface_diamond_contracts.php; do
  ./containers/php-version/compare-all.sh "$fixture" ./bin/gophp 2>&1 | tee -a docs/validation/diamond-types-fix-targeted-mac.log
  target_exit=$?
  echo "$fixture exit=$target_exit" | tee -a docs/validation/diamond-types-fix-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/diamond-types-fix-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/diamond-types-fix-compare-all-mac.log
```

## 直前の基準：収録588件は3版とも一致（履歴）

ユーザーがMacで実行した[修正後の公式PHP全件ログ](validation/item7-588-user-comparison-after-fix.log)を確認・保存した。PHP 8.3.35・8.4.26・8.5.11の各版で**588 matched、0 different、0 errors**、`make compare-all`終了コード0。fixtureの実数は正常296・実行時223・構文69で計588件。新規27件を含む収録fixtureの「公式PHP比較待ち」は解消した。これは今回新たに比較を実行した結果ではなく、ユーザーのMacでの実行記録である。

経緯は、[修正前の不一致](validation/item7-588-user-comparison-before-fix.log)（8.3/8.4で各8件、8.5で7件）→下記の8件の修正→今回の全件一致、の順である。途中の[Podman接続拒否](validation/item7-588-fix-compare-all-attempt.log)と保存済み出力への再照合も、その時点の記録として残す。fixture外の境界条件や項目7全体・項目6〜10全体を完了扱いにしない。[依存監査](dependency-audit.md)にも現在の機能別判定を記録した。

| 項目 | 現在の対応範囲 | 主な残件 |
| --- | --- | --- |
| 6 複数ファイル | ローカル`include` / `require` / `*_once`、`__FILE__` / `__DIR__`、読み込みスコープ・履歴の部分集合 | include_path・stream wrapper、読み込み先例外とファイル間トレースの統合 |
| 7 クラス | 継承・interface・trait・staticメンバー・型なしクラス定数・`self` / `parent` / `static`の部分集合 | 型付き署名の互換性、動的名、未収録の宣言競合・診断 |
| 8 型 | 未実装 | 引数・戻り値・プロパティの型、strict_types、readonly、enum等 |
| 9 呼び出し | `global`・`static`ローカル、文字列の変数関数の部分集合 | クロージャ・参照捕捉、可変長・名前付き引数、引数展開等 |
| 10 例外 | 未実装 | throw、try/catch/finally、Throwable、エラーハンドラ等 |

Composer生成autoload、Reflection・Attribute、HTTP・PDO、Go製PHPでのLaravel/API Platform起動も後続の残件であり、公式PHPでの調査アプリ成功とは分ける。

文書検索で見つかった古い記述は、README・互換性文書・依存監査表・比較コンテナREADMEの旧件数や未対応断定を**要更新**として改めた。下の各時点の接続拒否、不一致、暫定期待値、および第1版・項目3〜5時点の件数は**履歴として維持**する。調査アプリのキャッシュ診断、Go製コンテナ自体の未実施検証、Composer autoload以降の未実装は**別件の未実施・残件**として扱う。

## 過去の修正・試行履歴

以下の節の「現在」「今回」「未実施」「暫定」は各節の記録時点を指す。上の588件全一致という現在の比較結果とは区別する。

## 588件の公式比較で判明した差異の修正（2026-10-14）

ユーザー提供の[修正前の全件ログ](validation/item7-588-user-comparison-before-fix.log)は `make compare-all` 終了2、PHP 8.3.35・8.4.26が各 **580一致・8不一致・0環境エラー**、8.5.11が **581一致・7不一致・0環境エラー**を示す。旧561件の全一致は維持され、新規27件の「公式比較待ち」はこの実行によって解消した。ただし588件全一致には至っていない。以下の「未実施・暫定」は以前の試行履歴として残す。

差異8件は、final定数の上書き、interface実装の引数名付き診断、traitメソッド競合の相手名、PHP 8.3/8.4におけるtraitと親プロパティの非互換診断、宣言の発生時点、未実行分岐内の相対クラス名とstaticプロパティunsetの早過ぎる拒否だった。`class_final_constant.php`、`class_relative_outside.php`、`class_static_unset.php`、`interface_arity.php`、`interface_missing_implementation.php`、`trait_method_collision.php`、`trait_separate_use_collision.php`、`trait_parent_property_version.php` の期待値のみ、提供ログの公式PHP出力に合わせて訂正した。以前の561件のソース・期待値は変更していない。fixture追加は **0件**、総数は **588件**。

クラス本体の宣言検査は既存の前方参照を維持し、実行文を経たinterface実装・trait合成は宣言到達時に処理する。存在しない依存を事前に検出する既存の検査も維持する。`self::class` とstaticプロパティの `unset` は未選択分岐では実行しない。実行されたstaticプロパティの `unset` は引き続き拒否する。[ローカル検証](validation/item7-588-fix-local.log)では `make test`、`make vet`、`make build`、`git diff --check` がすべて終了0（Python 12件）。[保存済み公式出力への再照合](validation/item7-588-fix-saved-output-recheck.log)は8件×3設定で一致し、元の不一致23組のPHP出力をバイト単位で確認した。これは**修正後の公式PHP再実行ではない**。

この環境の[全件比較試行](validation/item7-588-fix-compare-all-attempt.log)はPodman接続拒否により、各版ともバージョン確認の終了125で停止した。`make compare-all` 終了2、**今回の再実行は各版0一致・0不一致・1環境エラー・588件未実施**。PHPプログラム自身の終了125やfixture不一致とは数えない。Macでの再比較は以下を実行し、パイプ直後の終了コードを保存する。既存コンテナは削除しない。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make compare-all 2>&1 | tee docs/validation/item7-588-fix-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/item7-588-fix-compare-all-mac.log
```

## 分割したtrait useの競合修正（2026-10-14）

前回の判定は`decl.methods`の存在だけを見ていたため、先行する`use`で取り込んだメソッドもクラス本体の定義と誤認した。`methodDecl.traitSource`で取り込み元を保持し、後続の`use`と同名なら競合として拒否する。クラス本体へ直接書いたメソッドは引き続きtraitより優先し、`insteadof`による同一use内の解決も維持する。[PHP公式trait規則](https://www.php.net/manual/en/language.oop5.traits.php)では未解決の同名traitメソッドは致命エラーとなる。

`runtime/trait_separate_use_collision.php`と対比用の`success/trait_separate_use_class_method.php`を追加した。Goの8.3・8.4・8.5設定では前者がstdout空・終了255、後者が`3\n`・終了0。前者の診断文・行番号・PHP 8.5のトレースは公式PHP実比較前の**暫定値**。旧561件のソースと期待値は変更していない。現在は**588 fixture**（旧561件＋追加27件）。

[ローカルログ](validation/item7-separate-use-local.log)では`make test`、`make vet`、`make build`、`git diff --check`が各終了0（Python 12件）。[Podman状態](validation/item7-separate-use-podman.log)はソケット接続拒否で終了125。新規失敗系・正常系の比較器はそれぞれ終了2で公式PHPの版確認前に停止（[失敗系](validation/item7-separate-use-targeted.log)、[正常系](validation/item7-separate-use-success-targeted.log)）。[全件試行](validation/item7-separate-use-compare-all.log)も環境エラーで`make compare-all`終了2。**今回の実比較は各版0一致・0不一致・1環境エラー・588件未実施**。保存済みの旧561件一致と混同しない。

Macでは既存コンテナの状態を確認し、停止中のものだけ起動してから次を実行する。終了コードは各コマンド直後に保存する。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/item7-separate-use-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/item7-separate-use-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
: > docs/validation/item7-separate-use-targeted-mac.log
for file in tests/fixtures/runtime/trait_separate_use_collision.php tests/fixtures/success/trait_separate_use_class_method.php; do
  ./containers/php-version/compare-all.sh "$file" ./bin/gophp 2>&1 | tee -a docs/validation/item7-separate-use-targeted-mac.log
  target_exit=$?
  echo "$file exit=$target_exit" | tee -a docs/validation/item7-separate-use-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/item7-separate-use-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/item7-separate-use-compare-all-mac.log
```

## レビュー指摘2件の修正（前段階・586件時点）

`self::`・`parent::`・`static::`による静的メソッド呼び出しは、メソッドを探すクラスと遅延静的束縛の呼出し元クラスを分けた。明示的な`A::`は新しい呼出し元を設定する。複数traitの同名メソッドについては、利用クラス自身に同名のメソッドがあればそれを優先し、trait候補同士の衝突を拒否しない。クラスに実装がない場合の競合検査と`insteadof`・`as`の処理は維持した。根拠は[PHPの遅延静的束縛](https://www.php.net/manual/en/language.oop5.late-static-bindings.php)と[traitの優先順位](https://www.php.net/manual/en/language.oop5.traits.php)。

再現2件を`success/class_forwarding_static_calls.php`、`success/trait_conflict_class_method.php`として追加し、3世代の`self::`・`parent::`・`static::`と明示クラス名呼出しの差をGo単体テストでも確認した。旧561件のソースと期待値は変更していない。現在は**586件**（前段階の追加23件と今回の2件を含む）で、追加25件の期待値は公式PHP実比較前の暫定値。[追加一覧](validation/item7-fixtures.txt)を参照。

[ローカルログ](validation/item7-review-local.log)では`make test`、`make vet`、`make build`、`git diff --check`が各終了0。Podman状態確認はソケット接続拒否で終了125（[ログ](validation/item7-review-podman.log)）。新規2件の対象比較は各版のPHP確認前に停止して比較器終了2（[遅延束縛](validation/item7-review-targeted.log)、[trait](validation/item7-review-trait-targeted.log)）。`make compare-all`も同じ環境エラーで終了2（[全件ログ](validation/item7-review-compare-all.log)）。今回の公式PHP比較は各版**0一致・0不一致・1環境エラー・586件未実施**。既存561件の保存済み比較結果を今回の新コードの比較結果として扱わない。

Macでは既存コンテナの状態を確認し、停止中のものだけ起動してから、下記を実行する。各終了コードはパイプの直後に保存する。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/item7-review-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/item7-review-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
: > docs/validation/item7-review-targeted-mac.log
for file in tests/fixtures/success/class_forwarding_static_calls.php tests/fixtures/success/trait_conflict_class_method.php; do
  ./containers/php-version/compare-all.sh "$file" ./bin/gophp 2>&1 | tee -a docs/validation/item7-review-targeted-mac.log
  target_exit=$?
  echo "$file exit=$target_exit" | tee -a docs/validation/item7-review-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/item7-review-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/item7-review-compare-all-mac.log
```

## 項目7の継続実装（前段階・584件時点）

開始基準の[561件比較ログ](validation/class-private-diagnostic-compare-all-mac.log)には、PHP 8.3.35・8.4.26・8.5.11の各版で561一致・0不一致・0環境エラー、`make compare-all`終了0がある。過去の「14件比較待ち」は解消した。今回、旧561件を変更せず、[追加一覧](validation/item7-fixtures.txt)の正常系14件・失敗系9件を加えて**584件**にした。新規23件の期待値は公式PHP実比較前の**暫定値**で、Go出力を含むローカル動作とPHP公式マニュアルの規則で検討したもの。Goの期待値一致を公式比較済みとは扱わない。

| 領域 | 今回動くようにした範囲 | 残る範囲・確認事項 |
| --- | --- | --- |
| interface | 宣言、複数extends、複数implements、親から継承する契約、公開メソッドの実装義務、static一致、型なし引数数の基本検査、`instanceof` | 型付き署名、複数interface間の競合・重複・循環の完全検査、診断文言と宣言時点の3版照合 |
| trait | 宣言、trait内use、複数取込、メソッド・プロパティ・定数のクラスごとの複製、`insteadof`・`as`の基本、抽象メソッド、`as final`、staticメンバーの取込 | adapt指定の全競合、プロパティ・定数の厳密な互換性、trait直参照の非推奨警告、循環と条件付き読込の診断。[PHP公式マニュアルの8.5差](https://www.php.net/manual/en/language.oop5.traits.php)に沿って親とtraitの同名プロパティの優先順位を切り替えたが、実比較待ち |
| staticメンバー | 宣言・初期化・読取・代入、複合代入・増減・配列要素更新、`isset`・`empty`・`??`・`??=`、継承共有と再宣言分離、static呼出し、`$this`なしのフレーム、直接`unset`の拒否 | 可変クラス名・動的staticメンバー、継承時の全アクセス・診断、trait内staticの版境界 |
| クラス定数 | class/interface/traitの型なし定数、可視性、継承、final、`self::`等からの参照と既定値の一部 | 型付き定数（項目8）、すべての定数式・衝突とバージョン差、正確な宣言診断 |
| self/parent/static | `::class`、クラス定数・staticプロパティ・メソッド、新規生成、継承メソッド内の遅延静的束縛、クラス外の字句検査 | 動的クラス名と特殊な呼出し文脈、親なし・クラス外の正確な診断 |
| 宣言検査 | final/abstract/可視性に加えて、静的性の不一致と型なし引数の最小互換性を検査 | 型宣言と戻り値、interface/traitの全競合・循環、正確なバイト列・発生時点 |

実装は`classDecl`に種類とinterface/trait依存を持たせ、前方登録と実行時登録の双方で依存を解決する。traitメンバーは利用クラスごとに複製して宣言スコープ・privateキー・staticセルを分ける。staticプロパティはクラスごとのセルを持ち、継承時は同一セル、再宣言時は新セルを使う。`self`は宣言元、`static`は呼び出されたクラスを参照する。既存の配列値コピーとオブジェクト共有には変更を加えていない。

[ローカルログ](validation/item7-local.log)で`make test`・`make vet`・`make build`・`git diff --check`は各終了0（Go全パッケージ、Python 12件、584 fixtureの版別期待値）。[Podman状態ログ](validation/item7-podman.log)はソケット接続拒否で終了125。対象比較は比較器終了2（[対象ログ](validation/item7-targeted-compare-attempt.log)）、`make compare-all`も3版の公式PHP確認で停止して終了2（[全件ログ](validation/item7-compare-all-attempt.log)）。**今回追加分の公式比較は各版0一致・0不一致・1環境エラー・23件未実施**。旧561件の保存済み一致は有効だが、新コードでの再比較は各版584件未実施。Podmanの125をPHPプログラムの終了コードに数えない。

Macでは既存コンテナの稼働状態を確認し、停止中のものだけ起動してから次を実行する。各パイプの直後に終了コードを保存する。

```sh
cd /Users/kouji/dev/php
podman ps -a --format '{{.Names}} {{.Status}}'
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/item7-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/item7-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
: > docs/validation/item7-targeted-mac.log
while IFS= read -r fixture; do
  ./containers/php-version/compare-all.sh "$fixture" ./bin/gophp 2>&1 | tee -a docs/validation/item7-targeted-mac.log
  fixture_exit=$?
  echo "$fixture exit=$fixture_exit" | tee -a docs/validation/item7-targeted-mac.log
done < docs/validation/item7-fixtures.txt
make compare-all 2>&1 | tee docs/validation/item7-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/item7-compare-all-mac.log
```

次は対象比較の不一致から宣言診断と版差を修正し、interfaceの複数継承競合・型なし署名、trait adaptationの衝突、staticの可視性と継承を追加fixtureで確定する。型付き署名は項目8に残す。項目7も項目6〜10全体も完了扱いにしない。

## 561件の公式比較確定

Macで再実行した[比較ログ](validation/class-private-diagnostic-compare-all-mac.log)を確認した。PHP 8.3.35・8.4.26・8.5.11の各版で **561 matched、0 different、0 errors**、`make compare-all`終了0。直前まで「公式PHP比較待ち」だった14件も含む。以下の比較待ち・失敗・接続拒否は当時の履歴として残す。今回追加するfixtureにはこの561件の確認結果を流用しない。

## 561件の公式比較と宣言診断の修正

ユーザーがMacで実行した `make compare-all` の[全ログ](validation/class-private-561-user-comparison-before-fix.log)を確認した。修正前はPHP 8.3.35・8.4.26が各 **560一致・1不一致・0環境エラー**、8.5.11が **556一致・5不一致・0環境エラー**（各561件）、`make compare-all`終了2。未比較だった14件は実行され、同名privateプロパティを含むほかの新規ケースも一致した。以前の「新規14件は未比較」は、以下に過去の試行履歴として残す。

共通の不一致は `class_final_parent.php` の古い文言で、公式PHPは `Class FinalDescendant cannot extend final class FinalParent` と出す。PHP 8.5のみ、宣言段階のFatal error 5件に `Stack trace:\n#0 {main}\n` が付き、抽象メソッド1件の文言が `remaining method`（単数）になる。公式ログのバイト列を根拠に実装と版別期待値を修正した。クラス宣言の検査エラーを `fatalDiagnostic` として扱い、既存のPHP 8.5コンパイル時診断経路でトレースを付ける。その他の実行時エラーには一律にトレースを追加しない。修正後のGoローカル検証結果と公式PHP再比較は別記録とし、後者が未実施なら全一致とはしない。

修正後、`make test`、`make vet`、`make build`、`git diff --check` は各終了0。[ローカル検証ログ](validation/class-private-diagnostic-local.log)にコマンドと終了コードを保存した。修正前の公式比較ログにある不一致7組（8.3が1件、8.4が1件、8.5が5件）のPHP側stdout・stderr・終了コードを抽出し、修正後の`./bin/gophp --php-version=...`と照合して**7組すべてバイト一致**した。[保存済み公式出力との再照合ログ](validation/class-private-diagnostic-official-log-recheck.log)。これは保存済み出力への照合であり、新たな公式PHP実行ではない。

`podman ps -a` はソケット接続拒否で終了125（[環境ログ](validation/class-private-diagnostic-podman.log)）。修正後の`make compare-all`はGoローカルテスト・ビルド後、3版とも公式PHPのバージョン確認段階で停止し、make終了2（[全件比較試行ログ](validation/class-private-diagnostic-compare-all-attempt.log)）。**修正後の公式PHP実比較は各版0一致・0不一致・1環境エラー・561件未実施**。修正前の560/560/556一致と混同しない。Macでは既存コンテナの稼働状態を確認して必要なものだけ起動し、`set -o pipefail; make compare-all 2>&1 | tee docs/validation/class-private-diagnostic-compare-all-mac.log; compare_exit=$?; echo "make compare-all exit=$compare_exit" | tee -a docs/validation/class-private-diagnostic-compare-all-mac.log` を実行する。

## 宣言失敗系と同名privateプロパティ（修正前の記録）

前回の550件を変更せず、宣言失敗6件、クラス外からのprotectedアクセス失敗1件、同名privateプロパティの正常系4件を追加した。現在**561件**。宣言失敗は `runtime/class_final_parent.php`、`class_final_override.php`、`class_narrow_method.php`、`class_narrow_property.php`、`class_abstract_unimplemented.php`、`class_abstract_instantiate.php`。前5件は宣言前の`echo`が出ないこと、最後は生成前の`echo`だけが出ることをGoで確認した。正常系は `success/class_private_slots.php`、`class_private_public.php`、`class_private_three_levels.php`、`class_private_quiet.php`。クラス外のprotectedアクセスは`runtime/class_private_protected_outside.php`で確認する。期待値と診断は**公式PHP比較前の暫定値**である。

従来のオブジェクトプロパティは名前だけをキーにしたため、親子で同名のprivate宣言を保持できなかった。宣言元クラスと名前でprivate用の内部キーを作り、public/protectedは従来の名前キーを共有する。メソッド内では宣言元スコープのprivateを優先し、それ以外は子側から可視の宣言を探す。初期化、通常読取・代入、複合代入・増減、配列内更新、`isset`・`empty`・`unset`・`??`・`??=`、nullsafeの読取は共通のキー解決を通る。オブジェクト比較は内部の値表を比較するため、異なるprivateスロットも含む。オブジェクトの`foreach`は現実装で未対応のまま。内部キーは警告文やメソッド出力に使わない。参照・Reflection・マジックメソッドは対象外。

公式PHPの[可視性規則](https://www.php.net/manual/en/language.oop5.visibility.php)は、親のprivateプロパティを子の同名宣言が置き換えず、別の内部名を持つと明記している。[final規則](https://www.php.net/manual/en/language.oop5.final.php)にも基本診断例がある。これらは実装方針の根拠だが、**今回のfixtureのバイト単位比較結果ではない**。Podman状態確認は接続拒否で終了125。宣言1件とprivate正常系1件の対象比較はどちらも各版のバージョン確認で止まり、比較器終了2。全fixtureの再比較も未実施。既存547件の3版全一致は下記の保存済みログで確認できるが、新規14件を含む561件の一致とは扱わない。

Macでの対象比較と全件比較は、停止中の既存コンテナだけを起動して次を実行する。各終了コードはパイプの直後に保存する。コンテナの削除・初期化はしない。

```sh
cd /Users/kouji/dev/php
podman ps -a --format '{{.Names}} {{.Status}}'
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
./containers/php-version/check.sh 2>&1 | tee docs/validation/class-private-check-mac.log
check_exit=$?
echo "check.sh exit=$check_exit" | tee -a docs/validation/class-private-check-mac.log
GOCACHE=/tmp/gophp-go-cache make build 2>&1 | tee docs/validation/class-private-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/class-private-build-mac.log
: > docs/validation/class-private-targeted-mac.log
for file in \
  tests/fixtures/success/class_override_visibility.php \
  tests/fixtures/success/class_final_members.php \
  tests/fixtures/success/class_abstract_basic.php \
  tests/fixtures/runtime/class_final_parent.php \
  tests/fixtures/runtime/class_final_override.php \
  tests/fixtures/runtime/class_narrow_method.php \
  tests/fixtures/runtime/class_narrow_property.php \
  tests/fixtures/runtime/class_abstract_unimplemented.php \
  tests/fixtures/runtime/class_abstract_instantiate.php \
  tests/fixtures/runtime/class_private_protected_outside.php \
  tests/fixtures/success/class_private_slots.php \
  tests/fixtures/success/class_private_public.php \
  tests/fixtures/success/class_private_three_levels.php \
  tests/fixtures/success/class_private_quiet.php; do
  ./containers/php-version/compare-all.sh "$file" ./bin/gophp 2>&1 | tee -a docs/validation/class-private-targeted-mac.log
  compare_exit=$?
  echo "$file compare exit=$compare_exit" | tee -a docs/validation/class-private-targeted-mac.log
done
GOCACHE=/tmp/gophp-go-cache make compare-all 2>&1 | tee docs/validation/class-private-all-mac.log
all_exit=$?
echo "make compare-all exit=$all_exit" | tee -a docs/validation/class-private-all-mac.log
```

Goの期待値テストは561件を全設定で通過した。公式PHP比較ができたら、まず新規14件の実診断・発生時点と前回3件の期待値を照合し、版差を修正する。その後、全件を再比較する。今回のクラス部分集合と項目7全体は区別し、項目6〜10全体も未完了のまま。

| 今回の検証 | 終了コードと結果 | ログ |
| --- | --- | --- |
| `make test` / `make vet` / `make build` / `git diff --check` | 全て0。Go全パッケージ・Python 12件、561 fixtureのローカル期待値が成功 | [ローカル](validation/class-private-local.log) |
| `podman ps -a` | 125。ソケット接続拒否でコンテナ状態は未確認 | [Podman](validation/class-private-podman.log) |
| 宣言失敗1件の3版対象比較 | 各版のバージョン確認で環境エラー、比較器終了2、fixture実行0 | [宣言](validation/class-private-declaration-compare-attempt.log) |
| private正常系1件の3版対象比較 | 同じ環境エラー、比較器終了2、fixture実行0 | [private](validation/class-private-targeted-compare-attempt.log) |
| `make compare-all` | ローカルテスト・ビルド後に各版のバージョン確認で環境エラー、make終了2 | [全件](validation/class-private-all-compare-attempt.log) |

今回の全件比較試行では、PHP 8.3.35・8.4.26・8.5.11の各版が**一致0・不一致0・環境エラー1・未実施561件**。Podmanの終了125をPHPプログラム自身の終了コードとして数えない。Goの失敗系診断は公式PHPの正確なバイト列・行番号・版差を未確認のため暫定である。

## 547件比較後の継続実装

検証済み547件のソースと期待値は変更せず、`success/class_override_visibility.php`、`success/class_final_members.php`、`success/class_abstract_basic.php` を追加した。現在は**550 fixture**。順にpublic/protectedプロパティの合法な再宣言とprivateメソッドの呼び分け、final宣言、抽象メソッドを子が実装して親の具体メソッドから呼ぶ動作をGo側で確認する。新規3件の期待値は公式PHP比較前の**暫定値**である。

クラス登録時に親のpublic/protectedメンバーと子の可視性を比較し、狭める再宣言を拒否する。finalクラスの継承とfinalメソッドの上書きも拒否する。親のprivateメソッド内からの同名呼出しは宣言元クラスのprivateメソッドへ解決する。`abstract class`と抽象メソッドの最小構文、具象子クラスの実装義務、抽象クラス直接生成の拒否を追加した。親のprivate**プロパティ**と子の同名プロパティは、現在の値表が単一名で格納しているため、誤った結果で実行せず明示的に拒否する。この衝突を正しく扱うには宣言元を含むプロパティ格納キーが必要。型付き宣言、interface、trait、抽象宣言の全制約、staticメンバー、クラス定数、self/parent/static、型互換性は残る。項目7は部分実装、項目8・10は未実装、項目6・9も部分実装である。

| 検証 | 結果 | 記録 |
| --- | --- | --- |
| `make test`、`make vet`、`make build` | 各終了0。Go全パッケージとPython 12件成功。buildのstat cache書込警告は終了0 | [ローカル](validation/items6-10-continuation-local.log) |
| `podman ps -a` | ソケット接続拒否、終了125。コンテナの稼働状態は未確認 | [環境](validation/items6-10-continuation-podman.log) |
| `class_override_visibility.php`の対象比較試行 | 各版のバージョン確認で環境エラー。比較器終了2、fixture実行0 | [対象比較](validation/items6-10-continuation-targeted.log) |
| `class_abstract_basic.php`の対象比較試行 | 各版のバージョン確認で環境エラー。比較器終了2、fixture実行0 | [抽象クラス比較](validation/items6-10-continuation-abstract-targeted.log) |
| `make compare-all` | ローカルテストとビルド後、各版のバージョン確認で環境エラー。make終了2、fixture実行0 | [全件比較](validation/items6-10-continuation-compare-all.log) |

今回新たに実行できた公式PHP比較は、全件比較の各版で**一致0・不一致0・環境エラー1・未実施550件**。終了125はPodmanの接続エラーで、PHPプログラムの終了コードではない。先にユーザーから提供された547件の全一致・make終了0は、下記の保存ログで別に確認済み。今回の新規3件の公式PHP比較と全550件の再比較は未実施である。

MacでPodmanの状態を確認し、停止している比較コンテナだけ起動してから次を実行する。VM起動済みなら`podman machine start`は行わない。各終了コードは`tee`直後に保存する。

```sh
cd /Users/kouji/dev/php
podman ps -a --format '{{.Names}} {{.Status}}'
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
./containers/php-version/check.sh 2>&1 | tee docs/validation/items6-10-continuation-check-mac.log
check_exit=$?
echo "check.sh exit=$check_exit" | tee -a docs/validation/items6-10-continuation-check-mac.log
GOCACHE=/tmp/gophp-go-cache make build 2>&1 | tee docs/validation/items6-10-continuation-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/items6-10-continuation-build-mac.log
./containers/php-version/compare-all.sh tests/fixtures/success/class_override_visibility.php ./bin/gophp 2>&1 | tee docs/validation/items6-10-continuation-override-mac.log
override_exit=$?
echo "override compare exit=$override_exit" | tee -a docs/validation/items6-10-continuation-override-mac.log
./containers/php-version/compare-all.sh tests/fixtures/success/class_final_members.php ./bin/gophp 2>&1 | tee docs/validation/items6-10-continuation-final-mac.log
final_exit=$?
echo "final compare exit=$final_exit" | tee -a docs/validation/items6-10-continuation-final-mac.log
./containers/php-version/compare-all.sh tests/fixtures/success/class_abstract_basic.php ./bin/gophp 2>&1 | tee docs/validation/items6-10-continuation-abstract-mac.log
abstract_exit=$?
echo "abstract compare exit=$abstract_exit" | tee -a docs/validation/items6-10-continuation-abstract-mac.log
GOCACHE=/tmp/gophp-go-cache make compare-all 2>&1 | tee docs/validation/items6-10-continuation-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/items6-10-continuation-compare-all-mac.log
```

次は親privateプロパティの格納単位を宣言元別に直し、interface/abstractと型情報を同時に設計する。さらにstrict_typesのファイル文脈、クロージャ捕捉・引数束縛、throw/try/finallyとPHPエラーオブジェクトの順に最小fixtureを比較しながら進める。項目6の読込失敗・例外・スタックトレースは例外基盤と統合して検証する。Laravelの起動は引き続き未対応であり、固定lock・vendorは変更していない。

## 547件の公式PHP比較確定

ユーザー提供の `/tmp/php-compare-547.log` を[保存](validation/items6-10-547-user-comparison.log)した。ログは `make compare-all` の終了コード0、PHP 8.3.35・8.4.26・8.5.11のそれぞれ **547 matched、0 different、0 errors** を示す。新規20件も含め、現在の547件について「公式PHP比較未実施・期待値暫定」という記述は解消した。以下の未実施・環境エラー記録は、その時点の試行履歴として残す。これらの一致は項目6〜10の未実装機能まで検証した意味ではない。

## 親クラス読込の再修正（2026-10-13）

前回の事前登録は`require`式を走査していたため、関数呼び出し経由と`switch`内の読込を見落とした。現在は**未知の親を持つクラスの結合を保留**し、前の実行文が関数呼び出し・分岐・switch・include等を通じてクラスを定義し得る場合、クラス宣言に到達した時点の登録表で親を解決する。同じファイルの無条件な親子は依存順に前方登録する。純粋な先行`echo`しかない場合の既存の早期エラー検査は維持する。この保留は任意の呼び出し先を静的に追う方式ではなく、実行文を保守的に境界として扱う方式である。

`success/include_parent_via_function.php` と `success/include_parent_via_switch.php` を追加し、合計は**547件**。補助ファイルは既存の `support/review_base.inc` を再利用した。Go単体テストでは両ケースが`7`・終了0で、読込が実行されない場合は親なしで終了255、同一ファイルの前方参照も維持した。今回の2件を含む新規20件は公式PHP3版との実比較が**未実施**であり、期待値は暫定である。項目6〜10は引き続き未完了。

| 今回の検証 | 結果・終了コード | ログ |
| --- | --- | --- |
| `make test` / `make vet` / `make build` / `git diff --check` | 全て0。Goテスト全パッケージ、Python 12件 | [ローカル](validation/items6-10-parent-load-local.log) |
| `podman ps -a` | 接続拒否、125 | [環境](validation/items6-10-parent-load-container.log) |
| 追加2件の3版比較試行 | 各版ともバージョン確認で環境エラー。各比較器終了2、fixture実行0 | [対象比較](validation/items6-10-parent-load-targeted.log) |
| `make compare-all` | ローカルテスト・ビルド後、各版のバージョン確認で環境エラー。make終了2、fixture実行0 | [全件比較](validation/items6-10-parent-load-compare-all.log) |

今回の各版の実比較結果は**一致0・不一致0・環境エラー1・未実施547件**。終了125はPodmanの接続エラーで、PHPプログラムの終了コードではない。Mac側では停止中の比較コンテナだけを起動し、以下を実行する。VM起動済みなら再起動しない。

```sh
cd /Users/kouji/dev/php
podman ps -a --format '{{.Names}} {{.Status}}'
# 停止を確認した場合のみ: podman start php-version-83 php-version-84 php-version-85
set -o pipefail
./containers/php-version/check.sh 2>&1 | tee docs/validation/items6-10-parent-load-check-mac.log
parent_check_exit=$?
echo "check.sh exit=$parent_check_exit" | tee -a docs/validation/items6-10-parent-load-check-mac.log
GOCACHE=/tmp/gophp-go-cache make compare-all 2>&1 | tee docs/validation/items6-10-parent-load-compare-all-mac.log
parent_compare_exit=$?
echo "make compare-all exit=$parent_compare_exit" | tee -a docs/validation/items6-10-parent-load-compare-all-mac.log
```

## レビュー指摘への追補（2026-10-13）

以下の「538件」はレビュー修正前の履歴である。この段階では7件を追加し、当時は**545件**（従来の527件と前回11件を維持）となった。追加したのは `include_parent_class.php`、`include_parent_timing.php`、`class_forward_parent.php`、`conditional_parent_class.php`、`static_global_order.php`、`global_static_order.php`、`include_global_after_static.php`。補助ファイル `review_base.inc`、`review_bind.inc` は独立fixtureに数えない。追加7件の期待値は公式PHPとの実比較前の**暫定値**である。

レビュー指摘1の原因は、`executeProgram` が主ファイルの子クラスを実行前に登録し、その時点では `require` が提供する親が存在しなかったこと。同一ファイルの無条件宣言は親から再帰的に前方登録し、読込や条件付き宣言が親を提供し得る場合だけ子の宣言位置まで登録を保留する。前方参照、読込後の継承、条件付き親、存在しない親、循環、読込前の出力をGo単体テストで確認した。循環・不正な継承の**公式PHP互換診断と発生時点は未確認**である。

レビュー指摘2の原因は、同じ名前が `staticNames` と `globalNames` の両方に残り、読書きが常にstaticを優先したこと。フレームごとの単一 `bindings` 表に置き換え、最後に実行した宣言で格納先を選ぶ。`unset` は現在の束縛だけを外し、ルート値・関数別static値を保持する。複数回呼出し、独立実行、逆順、unset後の再宣言、include先のglobalをGo単体テストで確認した。[PHPの変数スコープ](https://www.php.net/manual/en/language.variables.scope.php)は両宣言を参照束縛として説明するが、**同名での宣言順と正確な出力は公式PHP3版で未確認**であり、Goの `11|11` は暫定である。

| 項目 | 現在の実装済み範囲 | 主な残件 |
| --- | --- | --- |
| 6 複数ファイル | 既存のローカルinclude群に加え、読込先を親にした継承と呼出元の束縛変更をGoで検証 | include_path、ラッパー、診断の厳密一致、例外との統合、複数ファイルの完全なスタックトレース |
| 7 クラス | 既存の単純な継承に、同一ファイル前方参照と読込後の親登録を追加 | interface、trait、abstract/final、staticメンバー、クラス定数、self/parent/static、private衝突、宣言互換性 |
| 8 型 | 未実装 | 引数・戻り値・プロパティ型、strict_types、union/intersection/DNF、readonly、昇格、enum、版差 |
| 9 呼出し | global/staticの宣言順による束縛切替を追加。文字列変数関数は前回から維持 | 参照捕捉、クロージャ、arrow、callable、可変長・名前付き・展開引数 |
| 10 例外 | 未実装 | throw、try/catch/finally、Throwable階層、ハンドラ、警告との統合 |

次は公式PHP3版で今回7件と前回11件の実比較を取り、診断・宣言時点を確定する。実装は型なしのクラス宣言検査、型メタデータとstrict_types、参照捕捉を含む呼出し、例外制御の順で小さなfixtureごとに進める。既存527件の過去の一致を新規18件の実比較済みという意味に拡張しない。

| 今回の検証 | 結果・終了コード | ログ |
| --- | --- | --- |
| `make test` / `make vet` / `make build` / `git diff --check` | いずれも0。Goテスト全パッケージ、Python 12件。ビルドのstat cache権限警告は非致命 | [ローカル](validation/items6-10-review-local.log) |
| `podman ps -a` | ソケット接続拒否、125 | [コンテナ](validation/items6-10-review-container.log) |
| 新規1件の3版比較試行 | 各版のPHPバージョン確認で環境エラー。比較器終了2、fixture実行0 | [対象比較](validation/items6-10-review-targeted.log) |
| `make compare-all` | ローカルテスト・ビルド後、各版のPHPバージョン確認で環境エラー。make終了2、fixture実行0 | [全件比較](validation/items6-10-review-compare-all.log) |

今回の公式PHP実比較は各版 **一致0・不一致0・環境エラー1・未実施545件**。比較器がバージョン確認で停止したため、終了125をPHPプログラムの終了コードやfixtureのDIFFとして扱わない。比較コンテナの稼働状態そのものも確認できなかった。Macでは下の既存コマンドのログ名を `items6-10-review-compare-all-mac.log` に変えて再実行できる。停止しているコンテナだけ起動し、既存コンテナは削除しない。

今回の状態を確認するためのMac側コマンドは次のとおり。`podman ps -a` で停止を確認したコンテナだけ `podman start` し、起動済みVMへ `podman machine start` は行わない。`make compare-all` は545件すべてを3版で比較し、正常系も終了コード0を要求する。各コマンドの終了コードは `tee` の直後に保存する。

```sh
cd /Users/kouji/dev/php
podman ps -a --format '{{.Names}} {{.Status}}'
# 必要な場合のみ: podman start php-version-83 php-version-84 php-version-85
set -o pipefail
./containers/php-version/check.sh 2>&1 | tee docs/validation/items6-10-review-check-mac.log
review_check_exit=$?
echo "check.sh exit=$review_check_exit" | tee -a docs/validation/items6-10-review-check-mac.log
GOCACHE=/tmp/gophp-go-cache make compare-all 2>&1 | tee docs/validation/items6-10-review-compare-all-mac.log
review_compare_exit=$?
echo "make compare-all exit=$review_compare_exit" | tee -a docs/validation/items6-10-review-compare-all-mac.log
```

主基準はPHP 8.4.26、追加対象は8.3.35・8.5.11。固定済みのLaravel 12.69.3 / API Platform Laravel 5.0.1 / 94パッケージのlock・vendorは変更していない。着手前の527 fixtureと期待値は維持した。今回の実装はGoが公式PHPを呼び出さない。**項目6〜10全体は未完了**であり、GoでのLaravel起動にも至っていない。

## 依存順と今回の到達点

| 順 | 単位 | 到達点 | 残件・次の最小fixture |
| ---: | --- | --- | --- |
| 1 | ファイル文脈と読込 | `include` / `require` / `_once`、`__FILE__` / `__DIR__`、ローカル実パス、読込元変数スコープ、読込先return、入れ子・循環中のonce、ファイル別の名前空間と事前検査を実装。`include_*.php` 3件とGo単体テスト | include_path・ストリームラッパー、パス探索の全条件、失敗時の警告とスタック、読込先での例外、読み込み済み一覧の公開。`include_missing.php` と `require_parse_error.php` を公式PHP3版で採取してから実装を詰める |
| 2 | 参照を要する呼び出し基盤 | `global $name` のルート束縛、関数呼び出しをまたぐ更新、`unset` によるローカル束縛解除を実装。`static $name` の呼び出し間保持・再帰共有・`unset`後のローカル再束縛も実装。文字列による変数関数呼び出しも追加。`global_*.php`・`static_local_*.php`・`function_variable_call.php` 計5件 | 参照セル、参照捕捉、static初期化の宣言境界・全型、callable形式・評価順。`global_missing.php` と名前空間での変数関数を比較する |
| 3 | クラス階層と宣言検査 | 単純な `extends`、継承したpublic/protectedプロパティ・メソッド・コンストラクタ、親クラスへの`instanceof`を実装。`class_extends_*.php`・`class_instanceof_parent.php` 計3件 | interface、trait、abstract/final、staticメンバー、クラス定数、self/parent/static、private名の衝突、完全なアクセス・宣言競合。型なしオーバーライドと不正継承を次に比較する |
| 4 | 型情報 | 未着手 | 引数・戻り値・プロパティ、nullable/union/intersection/DNF、strict_types、readonly、昇格、enum、版差。先に親子メソッドの型なし階層、次にスカラー引数とファイル別strict_typesを比較する |
| 5 | 呼び出し表現と例外制御 | 未着手 | クロージャ、arrow、変数関数、callable、可変長・名前付き・展開引数、throw/try/catch/finally、Throwable階層、エラーハンドラ。捕捉値と再throw/finallyを各1件から比較する |

上表の順3〜5は項目7〜10の相互依存に合わせて分割した。項目6・7・9は**部分実装**、項目8・10は**未実装**である。`include`や`extends`をパーサで受理しただけで完了とはせず、公式PHPの厳密比較が終わるまで新規ケースの期待値は暫定とする。

`evaluator.loadedFiles` は独立実行ごとに生成し、主ファイルと成功した読込先を実パスで記録する。`*_once` は現在実行中のファイルも含めて再読込を防ぐ。入れ子の読込では事前登録状態を保存・復元し、関数・クラス・定数の実行単位の登録表は共有する。`__FILE__` / `__DIR__` は各ファイルの解析時に解決する。失敗位置は `sourceError` で内側のファイル名を保持する。引数・添字・配列の既存の評価規則を変更せず、読込元の変数写像をそのまま使う。深さ256超はGo独自エラーでありPHP互換ではない。PHPの[include](https://www.php.net/manual/en/function.include.php)、[include_once](https://www.php.net/manual/en/function.include-once.php)、[マジック定数](https://www.php.net/manual/en/language.constants.magic.php)を設計の参照にした。これは資料に基づく実装判断であり、新規fixtureの実比較結果ではない。

`global` はルート写像を束縛し、`static` は関数宣言ごとの永続写像を束縛する実装で、[PHPの変数スコープ](https://www.php.net/manual/en/language.variables.scope.php)と[unset](https://www.php.net/manual/en/function.unset.php)の記述を参照した。参照一般、`$GLOBALS`、参照代入、参照foreachは未対応。今回、後続項目11の参照基盤全体を先行実装したわけではない。

固定vendorでは `research/laravel-api/vendor/autoload.php:20-22` の `require_once` を処理する入口を実装したが、その直後の静的メソッド呼び出しは未対応で止まる。`vendor/composer/autoload_real.php:5-48` のstaticメンバー、autoload登録、`Closure::bind`、`$GLOBALS` も残る。Composer生成autoloadの実行（項目16）は未完了である。固定vendorを文字列検索した限り、PHP 8.4の非対称可視性とプロパティフック形式の該当は見つからなかったが、構文解析による網羅検査ではない。8.3モードで8.4固有構文を受理する対応は行っていない。

## 検証と再実行

今回の追加fixtureは `success/include_once_scope.php`、`success/include_namespace.php`、`success/include_nested_once.php`、`success/global_binding.php`、`success/global_unset_binding.php`、`success/static_local_recursive.php`、`success/static_local_unset.php`、`success/class_extends_basic.php`、`success/class_extends_constructor.php`、`success/class_instanceof_parent.php`、`success/function_variable_call.php` の11件。補助ファイルは `tests/fixtures/support/*.inc` とし、`.php` のみを走査する比較器・Goテストで独立実行されない。合計538件。既存527件のソースと期待値は変更していない。

| 検証 | 結果・終了コード | 記録 |
| --- | --- | --- |
| `GOCACHE=/tmp/gophp-go-cache make test` | 初回は主ファイル診断パスの回帰で終了2。修正後はGo全パッケージとPython 12件成功、0 | [ローカルログ](validation/items6-10-local.log) |
| `GOCACHE=/tmp/gophp-go-cache make vet` | 成功、0 | [ローカルログ](validation/items6-10-local.log) |
| `GOCACHE=/tmp/gophp-go-cache go build -o bin/gophp ./cmd/gophp` | 成功、0。モジュールstat cacheへの権限警告あり | [ローカルログ](validation/items6-10-local.log) |
| `git diff --check` | 成功、0 | [ローカルログ](validation/items6-10-local.log) |
| `./containers/php-version/check.sh` | Podman接続拒否、125 | [ログ](validation/items6-10-container-check.log) |
| 新規1件の3版比較 | バージョン確認で各版停止、総合2 | [ログ](validation/items6-10-targeted.log) |
| `GOCACHE=/tmp/gophp-go-cache make compare-all` | ローカルテスト・ビルド後、各版のバージョン確認で停止。make終了2 | [ログ](validation/items6-10-compare-all.log) |

各版とも今回の**実比較は0一致・0不一致・バージョン確認の環境エラー1・fixture未実施538**。既存527件の過去の3版一致は[提供済みログ](validation/phase2-527-comparison-provided.log)の結果であり、今回再実行した結果ではない。新規11件の期待値はGo側テストに用いる暫定値である。Podmanの接続拒否をfixtureのDIFFやPHP自身の終了コード125には数えない。

Macで比較コンテナの状態を確認し、停止中のコンテナだけ起動してから次を実行する。すでに起動したVMに `podman machine start` は不要。終了コードはパイプの直後に保存する。

```sh
cd /Users/kouji/dev/php
podman ps -a --format '{{.Names}} {{.Status}}'
# 停止している名前があれば、その名前だけ podman start する
# 例: podman start php-version-83 php-version-84 php-version-85
set -o pipefail
./containers/php-version/check.sh 2>&1 | tee docs/validation/items6-10-container-check-mac.log
check_exit=$?
echo "check.sh exit=$check_exit" | tee -a docs/validation/items6-10-container-check-mac.log
GOCACHE=/tmp/gophp-go-cache make compare-all 2>&1 | tee docs/validation/items6-10-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/items6-10-compare-all-mac.log
```

比較が失敗したら新規fixtureのstdout・stderr・終了コードを個別に比較し、公式PHPの実結果を根拠に原因を修正する。保存済み期待値を書き換えて差異を隠さない。調査アプリのキャッシュ書込エラーは別件として継続追跡する。
