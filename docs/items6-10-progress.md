# 項目6〜10の進捗（2026-10-15）

現在の検証状況（2026-10-18）：項目11の型なしプロパティ参照を追加し、fixture・期待値は**974件**。Goローカル検証済み、公式PHP3版はPodman接続拒否で未比較。項目6〜10からの残件を維持する。[項目11の現状](item11-progress.md)。

追加前の比較済み基準：959件は[版別生ログ](item11-progress.md)で公式PHP3版が各959一致・終了0。以下の959件「未比較」は比較前の履歴。

959件比較前（履歴）：項目11の呼び出し元スコープの参照共有を修正し、fixture・期待値は**959件**。Goローカル検証済み、公式PHP3版はPodman接続拒否で未比較。項目6〜10からの残件を維持する。[項目11の現状](item11-progress.md)。

今回の修正前（履歴）：項目11の配列コピー参照境界を修正し、fixture・期待値は**958件**。Goローカル検証済み、公式PHP3版はPodman接続拒否で未比較。項目6〜10からの残件を維持する。[項目11の現状](item11-progress.md)。

958件への修正前（履歴、2026-10-18）：項目11の配列要素参照を加えた**955件**はGoローカル検証済み、公式PHP3版はPodman接続拒否で未比較。追加11件の期待値は暫定。[項目11の進捗・ログ](item11-progress.md)。項目6〜10の残件と項目11全体の部分実装は維持する。

955件追加前の基準（履歴）：項目11の変数間参照を加えた**944件**は[版別生ログ](item11-progress.md)で公式PHP3版とも944一致・0不一致・0エラー、終了0。追加7件の比較待ちは解消した。項目6〜10からの残件は後続へ引き継ぐ。今回は保存済みログの確認で比較を再実行していない。[項目11の対応範囲](item11-progress.md)。以下は各時点の履歴。

944件比較前の移行状況（履歴、2026-10-18）：[版別生ログ](items6-10-integration-gate.md)で937件の公式PHP3版比較は各版全件一致・終了0と確認でき、項目11への移行条件を満たした。項目6〜10の未対応機能は後続へ引き継ぐ。変数間参照の追加7件を含む944件の公式比較は環境エラーで未実施。[項目11の進捗](item11-progress.md)。以下の「比較待ち」は各時点の履歴。

## `define()`第3引数の警告中断（2026-10-18、937件比較前の履歴）

第3引数`null`の非推奨診断ハンドラが例外を投げた際、定数を登録しないよう修正。`true`指定時の警告は別経路として対比するfixtureを追加した。総数は**937件**（正常517・実行時288・構文132）。[test](validation/warning-audit-third-null-test.log)・[vet](validation/warning-audit-third-null-vet.log)・[build](validation/warning-audit-third-null-build.log)は終了0。公式PHP3版の[対象・全件比較](items6-10-integration-gate.md)は環境エラーで未実施、追加2件の期待値は暫定。項目11への移行は公式比較を条件とする。以下は936件以前の履歴。

## 警告中断経路の追加監査（2026-10-18、936件時点の履歴）

`constant_builtin.go`の`define(null, ...)`で診断のerrorを捨てて定数登録に進む経路を修正した。新規`integration_warning_define_state.php`は定数未登録、後続の正常な登録、`strlen`・`count`・cast・`foreach`でハンドラが投げた後の副作用回数を観測する。fixture・期待値は**936件**（正常516・実行時288・構文132）。[test](validation/warning-audit-test.log)・[vet](validation/warning-audit-vet.log)・[build](validation/warning-audit-build.log)は終了0。公式PHP3版の[対象・全件比較](items6-10-integration-gate.md)はPodman接続拒否で各終了2、修正後936件は未比較で期待値は暫定。項目11移行はこの比較が条件。次節は935件時点の履歴。

## 複合代入の失敗結果を区別（2026-10-18、936件追加前の履歴）

追加2件後の**935件**（正常515・実行時288・構文132）は[8.3](validation/items6-10-compound-final-compare-83-mac.log)・[8.4](validation/items6-10-compound-final-compare-84-mac.log)・[8.5](validation/items6-10-compound-final-compare-85-mac.log)の生ログで各版935一致・0不一致・0エラー。`integration_compound_failure_state.php`と`integration_concat_handler_state.php`も3版でMATCH。型付きプロパティの`/= 0`・`%= 0`・`+= []`・警告中断の`+= '2x'`で値0を確認し、`.=`は左辺型によって値保持と0格納が分かれる。通常変数・配列要素の警告中断は別に確認した。比較コマンド終了0と`check.sh`のPASSはユーザー報告で、ログには含まれない。ローカルのtest・vet・build・差分検査もユーザー報告で各終了0。ビット演算への一般化は未検証。[移行条件・過去の接続拒否](items6-10-integration-gate.md)を参照。次段落は934件時点の履歴。

934件時点の履歴：前回の「4ケースでは格納しない」という判断を撤回してGo側を修正したが、当時の[対象比較](validation/items6-10-compound-recheck-targeted-84.log)・[全件比較](validation/items6-10-compound-recheck-all-84.log)はPodman接続拒否で終了2、修正後の公式一致は未確認だった。[修正前集計](validation/items6-10-compound-934-before-fix-summary-mac.log)と[933件一致集計](validation/items6-10-compound-933-after-fix-summary-mac.log)を履歴として保存している。

## 剰余複合代入の例外伝播（2026-10-17、934件追加前の履歴）

追加2件で**933件**（正常513・実行時288・構文132）。`integerOperand`の警告・非推奨診断からのエラーを`modulo`と`evalCompound`へ返し、格納前に中断する。共通変換を使うbitwiseも同じ伝播に接続した。変数・配列要素の値保持、型付きプロパティへの結果0の格納、ハンドラが例外を投げない通常`%=`をGoの3設定で確認。[test](validation/items6-10-modulo-test.log)・[vet](validation/items6-10-modulo-vet.log)・[build](validation/items6-10-modulo-build.log)は終了0。公式PHPの[対象比較](validation/items6-10-modulo-targeted-84.log)と[全件比較](validation/items6-10-modulo-compare-84.log)はPodman接続拒否で各終了2。各版0一致・0不一致・1環境エラー・933件未実施。追加期待値は暫定。[移行条件](items6-10-integration-gate.md)。次節の931件は修正前の履歴。

## 統合検証と項目11への移行判定（2026-10-17、933件への追加前の履歴）

追加7件で現在**931件**（正常511・実行時288・構文132）。読み込み×型・callable、readonly／昇格×型エラー、警告ハンドラ×格納、finally×ClosureをGoで統合検証した。警告ハンドラの例外後も動的プロパティ作成などが進む不具合を修正した。[検証一覧、既存fixtureとの対応、A／B／Cの残件分類と移行条件](items6-10-integration-gate.md)を参照。[ローカルテスト](validation/items6-10-integration-test.log)は通過。公式PHP3版との[対象比較](validation/items6-10-integration-targeted-84.log)と[版別全件比較](validation/items6-10-integration-compare-84.log)はPodman接続拒否で停止し、各版0一致・0不一致・1環境エラー・931件未実施。新規期待値は暫定。下の924件全一致は今回の修正前の履歴で、項目6〜10の全仕様対応は未完了。

## callable／Closureの基本経路（2026-10-17、今回の変更前に収録924件を公式比較済み）

ユーザーが版ごとに`make compare`を実行した[6行集計ログ](validation/callable-924-after-fix-summary-mac.log)はPHP 8.3.35・8.4.26・8.5.11の各版で**924一致・0不一致・0エラー**を記録する。fixture実数と期待値も924件（正常504・実行時288・構文132）。908件以降の追加16件について、収録範囲の公式比較待ちは解消した。`check.sh`の各版PASSはユーザー報告であり、集計ログに個別fixture行・各コマンドの終了コードは含まれない。今回は元ログをそのまま保存・照合したもので、比較を再実行していない。下記のPodman接続拒否と暫定期待値は修正直後の履歴である。

### Go実装直後の検証記録（履歴）

当時の回帰修正：`call_user_func_array(args: [2, 3], callback: 'add')`がコールバックを配列と誤認していた。組み込み関数自身の`callback`・`args`を共通の`bindCallArguments`で束縛してから呼び出す。`call_user_func`も同じ処理を使い、可変長引数の名前を転送する。指定順を逆にした正常例と、不明名・二重指定の捕捉例を2 fixtureとGo単体テストに追加した。[ローカルテスト](validation/callable-named-test.log)・[vet](validation/callable-named-vet.log)・[build](validation/callable-named-build.log)・[差分検査](validation/callable-named-diff-check.log)は終了0。公式PHP3版の[8.3](validation/callable-named-targeted-83-attempt.log)・[8.4](validation/callable-named-targeted-84-attempt.log)・[8.5](validation/callable-named-targeted-85-attempt.log)の対象比較と[全件比較](validation/callable-named-compare-all-attempt.log)はPodman終了125で止まり、比較コマンド終了2。追加2件の期待値は暫定である。

固定vendorではSymfony Consoleの`Attribute/Option.php:54`に`is_callable($suggestedValues)`と`$suggestedValues(...)`、Laravel `Session/Middleware/StartSession.php:307`に`call_user_func`、Composer `autoload_real.php:32`に`call_user_func`と同`:37`に`Closure::bind`、API Platform Laravel `routes/api.php:58`に`__invoke`を対象とする`is_callable`がある。これは**ソース上の使用確認**であり、Goでvendor経路を実行した証拠ではない。

今回補完したのは関数名・`Class::method`文字列、二要素配列、Closureの共通callable解決、引数／戻り値の`callable`型、`Closure`型、参照出力引数を使わない`is_callable`、`Closure::fromCallable`、関数・オブジェクトメソッド・staticメソッドの第一級callable、基本的な`call_user_func`／`call_user_func_array`である。解決時に可視性とstatic性を確認し、第一級callableは対象メソッド・受け手・called scopeをClosure値に保持する。呼び出し時は既存の名前付き引数・可変長・型検査を通す。値捕捉・staticクロージャとエラーハンドラ接続も維持する。`callable`型のプロパティは宣言時に拒否し、昇格プロパティでも同じ制限を使う。

基本経路の追加時点では**14件**で合計**922件**（正常502・実行時288・構文132）。[当時のGo・Pythonテスト](validation/callable-local-test.log)、[vet](validation/callable-local-vet.log)、[build](validation/callable-local-build.log)、[差分検査](validation/callable-diff-check.log)を実施した。[当時の対象比較](validation/callable-targeted-84-attempt.log)と[版別全件比較](validation/callable-compare-84-attempt.log)はPodmanソケット接続拒否で公式PHPの版確認に失敗し、比較コマンド終了2。当時の回帰2件を加えた実数は**924件**（正常504・実行時288・構文132）。8.3・8.4・8.5の各版で**0一致・0不一致・1環境エラー・924件未実施**（うち新規16件）。この修正直後には新規期待値・診断を暫定とし、[908件一致](validation/special-types-908-after-fix-summary-mac.log)を新しい924件へ拡張しなかった。後続の冒頭集計で比較待ちは解消した。

未対応・未収録：`__invoke`一般、`Closure::bind`／`bindTo`／`call`、参照を要する`is_callable`第3引数、参照捕捉・参照引数・参照戻り値、Traversable展開、マジック`__call`／`__callStatic`、組み込み関数全域のcallback対応、callableエラー診断とトレースの境界、namespaceと継承callableの追加境界。`call_user_func`系は対応済み組み込み関数への部分到達に限り、vendorの`array_merge`等は未対応。項目11へ参照基盤を引き継ぎ、項目6〜10の統合検証では複数ファイル・strict_types・エラーハンドラ・継承・例外を通した実行と公式3版比較を残す。項目9・項目6〜10全体は部分実装のまま。

Macでの対象・版別全件比較（終了コードはコマンド直後に記録）：

```sh
cd /Users/kouji/dev/php
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  set -- $=spec
  python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php$1" --php-version "$2" --fixtures tests/fixtures/success/callable_builtin_named_order.php --interpreter ./bin/gophp "--php-version=$3" '{file}' 2>&1 | tee "docs/validation/callable-named-targeted-$1-mac.log"
  result=$?; printf 'exit_code=%s\n' "$result" >> "docs/validation/callable-named-targeted-$1-mac.log"
  make compare "PHP=./containers/php-version/bin/php$1" "PHP_VERSION=$2" 2>&1 | tee "docs/validation/callable-compare-$1-mac.log"
  result=$?; printf 'exit_code=%s\n' "$result" >> "docs/validation/callable-compare-$1-mac.log"
done
```

## 特殊型の908件公式比較済み範囲（2026-10-16、今回の変更前の履歴）

fixture実数と期待値は**908件**（正常490・実行時287・構文131）。ユーザーが版別に`make compare`を実行した[修正後6行集計ログ](validation/special-types-908-after-fix-summary-mac.log)には、PHP 8.3.35・8.4.26・8.5.11でそれぞれ**908一致・0不一致・0エラー**と記録されている。888件から追加した20件の公式比較待ちは、収録fixtureの範囲で解消した。元ログと保存先はバイト一致。今回はこちらで比較を再実行していない。ログに個別fixture行・コマンド終了コードはなく、修正6件の解消と`check.sh`の3版PASSはユーザー報告として区別する。過去の不一致とPodman接続拒否は下の履歴として維持する。

今回確定したのは収録908件の版別集計上の一致であり、項目8・項目6〜10全体の完成やfixture外の境界条件の検証完了ではない。autoload、未知階層の再検査、`callable`・`iterable`型の実行時意味、参照・Traversable、固定vendor経路のGo実行は引き続き残る。

## 特殊型の修正前公式比較と修正（2026-10-16、当時の履歴）

[ユーザーのMacでの修正前集計ログ](validation/special-types-908-before-fix-mac.log)は、8.3.35・8.4.26で各**902一致・4不一致・2テストエラー**、8.5.11で**903一致・3不一致・2テストエラー**。`check.sh`の各版PASSはユーザー報告であり、この集計ログは個別のstderr全文やコマンド終了コードを含まない。`special_final_static_85`だけ8.5.11で一致していた。今回、親なし`parent`診断の引用符、署名診断の`self`を宣言元クラス名にする処理、階層未確認時の子側の型名を修正した。継承比較のためだけに同一ファイル後方の型クラスを先行登録する処理を除き、`special_forward_variance`と`special_namespace_forward`を`success`から`runtime`へ移した。総数は**908件**（正常490・実行時287・構文131）のまま。既存888件のソース・期待値は変えていない。[make test](validation/special-types-review-test.log)、[vet](validation/special-types-review-vet.log)、[build](validation/special-types-review-build.log)は終了0。修正後の[8.3](validation/special-types-review-compare-83-attempt.log)・[8.4](validation/special-types-review-compare-84-attempt.log)・[8.5](validation/special-types-review-compare-85-attempt.log)の版別比較と[対象比較](validation/special-types-review-targeted-84-attempt.log)はPodman接続拒否により各終了2で、**各版0件一致・0件不一致・1件の環境エラー、908件未実施**。Goローカル検証の成功を公式比較済みとは扱わず、修正対象の期待診断は実比較まで暫定である。前節の「正常492・実行時285」は修正前の分類を指す。

修正直後は下記のMac再比較手順による対象6件と版別全件908件の比較が次の作業だった。後続の冒頭908件集計で収録fixtureの一致を確認したが、集計ログには各fixtureのstdout・stderr・終了コードは残っていない。項目8・項目6〜10全体は部分実装のまま。

## 特殊型・未登録型・宣言検査（2026-10-16、修正前の追加時点の履歴）

fixtureは**908件**（正常492・実行時285・構文131）。従来の888件は[版別6行集計](validation/enum-888-after-fix-summary-mac.log)で3版それぞれ888一致・差異0・エラー0。今回は20件を追加し、Goローカルでは3設定の期待値テストが通った。公式PHPとの[8.3](validation/special-types-compare-83-attempt.log)・[8.4](validation/special-types-compare-84-attempt.log)・[8.5](validation/special-types-compare-85-attempt.log)の版別全件比較と[対象比較](validation/special-types-targeted-84-attempt.log)は、すべてバージョン確認時のPodman終了125で比較コマンド終了2となった。**各版0件一致・0件不一致・1件の環境エラー、908件未実施**（うち新規20件）。この試行を公式PHPのfixture不一致に数えない。今回の期待値・診断は暫定であり、前の888件の比較済み範囲も今回の修正後に公式PHPで再検証されたわけではない。

`function.go`／`types.go`で複合型内を含む`self`・`parent`と戻り値専用`static`を宣言元クラスに結び、実行時の`static`はcalled scopeで検査する。`class.go`の引数反変・戻り値共変、型付きプロパティ・定数の比較にも宣言元を渡す。同一ファイルで後方にある無条件の型宣言は既存の前方登録へ接続し、実行文が先行していて未登録型を読み込む可能性がある場合は宣言登録を到達時まで遅らせる。未選択分岐やincludeを型検査のために実行しない。同じ未登録型名の契約は比較可能だが、未知の階層関係を要する不一致は`Could not check compatibility`診断にする。PHP 8.5のfinal子クラスによる`static`戻り値の例外も局所的に扱う。実装・期待診断の正確さは公式実比較待ちである。

既知の制限：autoload／spl_autoload_register、未知階層の再検査全般、循環・条件付き宣言の未収録境界、`callable`・`iterable`型の実行時意味、参照・Traversableは未対応。循環する親宣言の基本失敗例は収録したが、公式比較前の診断は暫定である。固定vendorの`static`／`self`型はソースで確認したが、その経路のGo実行は未確認。項目8および項目6〜10全体は部分実装のまま。次は公式PHP3版の対象・版別全件比較で宣言時点と診断を確定し、差異を修正してから未知階層の保留／再検査と`callable`・`iterable`の最小型判定を検討する。

Macでの再実行例（各版の生ログには個別fixture行と終了コードを保存する）：

```sh
cd /Users/kouji/dev/php
set -o pipefail
for spec in '83 8.3.35' '84 8.4.26' '85 8.5.11'; do
  set -- $=spec
  make compare "PHP=./containers/php-version/bin/php$1" "PHP_VERSION=$2" 2>&1 | tee "docs/validation/special-types-compare-$1-mac.log"
  result=$?
  printf 'command exit: %s\n' "$result" | tee -a "docs/validation/special-types-compare-$1-mac.log"
done
```

対象だけの確認には各版で`python3 skills/php-go-interpreter/scripts/compare_php.py --php ./containers/php-version/bin/php84 --php-version 8.4.26 --fixtures tests/fixtures/success/special_inherited_scope.php --interpreter ./bin/gophp --php-version=8.4 '{file}'`を使い、`php83`／8.3.35／8.3、`php85`／8.5.11／8.5へそれぞれ置き換える。終了コードは各コマンド直後に保存する。[型宣言の公式資料](https://www.php.net/manual/en/language.types.relative-class-types.php)、[採択済みstatic戻り値RFC](https://wiki.php.net/rfc/static_return_type)、[継承規則](https://www.php.net/manual/en/language.oop5.variance.php)を参照したが、資料確認は公式PHP実比較と区別する。

## enum宣言診断の4件を修正（2026-10-16、収録fixtureの比較済み）

ユーザーがMacで版別に`make compare`を実行した後の[修正後6行集計ログ](validation/enum-888-after-fix-summary-mac.log)には、PHP 8.3.35・8.4.26・8.5.11で各**888一致・0不一致・0エラー／888件**と記録されている。878件以降の追加11件の公式比較待ちは収録fixtureの範囲で解消した。ログには個別fixture行と各`make compare`の終了コードがない。修正4件の解消と`check.sh`の各版PASSはユーザー報告であり、今回はログの保存・照合のみで比較を再実行していない。項目8と項目6〜10全体は未完了。

修正前の[Mac集計ログ](validation/enum-888-before-fix-mac-summary.log)は、PHP 8.3.35・8.4.26・8.5.11の各版で**884一致・4不一致・0エラー／888件**。以下の修正直後はPodman接続拒否で再比較できなかったが、その後に上記の集計が得られた。4件の具体的な差異はユーザー報告を根拠とする。

Go側では、enumに直接宣言された予約メソッドの検査を解析時から宣言登録時へ移し、`before`の後にenum宣言行でFatalを出すようにした。プロパティ拒否文言を`Enum <name> cannot include properties`へ直し、trait由来の`from()`／`tryFrom()`はbacked enumでも許容した。2件を`parser/`から`runtime/`へ、traitの正常例を`runtime/`から`success/`へ移した。総数は**888件**（正常477・実行時281・構文130）のまま。[Go・Pythonテスト](validation/enum-888-diagnostic-test.log)、[vet](validation/enum-888-diagnostic-vet.log)、[build](validation/enum-888-diagnostic-build.log)、[差分検査](validation/enum-888-diagnostic-diff-check.log)は終了0。[修正後のGo個別実行ログ](validation/enum-888-diagnostic-go-individual.log)にはenum 10件×3設定のstdout・stderr・終了コードを保存した。修正直後の[対象比較試行](validation/enum-888-diagnostic-targeted-attempt.log)と[全件比較試行](validation/enum-888-diagnostic-compare-all-attempt.log)は、Podmanのバージョン確認が終了125となり、各コマンド終了2だった。この失敗は当時の履歴であり、後続の修正後集計とは区別する。

## 通常enumのfrom／tryFrom宣言修正（2026-10-16、修正前の履歴）

解析時とtrait合成後の再宣言検査で、`from()`／`tryFrom()`の禁止をbacked enumへ限定した。通常enumはユーザー定義のstaticメソッドを既存の呼び出し経路で実行する。`cases()`は通常／backed双方で禁止する。再現例の`: self`とtrait内の`?self`に必要なメソッド戻り値の最小対応も追加した。[公式のenumメソッド説明](https://www.php.net/manual/en/language.enumerations.methods.php)と[BackedEnum説明](https://www.php.net/manual/en/class.backedenum.php)を仕様根拠とし、公式PHP実比較とは区別する。

追加**4件**で総数**888件**（正常476・実行時280・構文132）。既存878件は過去の版別集計で3版一致。前回追加6件と今回追加4件、計10件は公式比較待ち。[Go・Pythonテスト](validation/enum-method-test.log)、[vet](validation/enum-method-vet.log)、[build](validation/enum-method-build.log)、[差分検査](validation/enum-method-diff-check.log)を参照。enum 10件をGoの3設定で個別実行した[stdout・stderr・終了コードの生ログ](validation/enum-888-go-individual.log)も保存した。これは公式PHP比較ではない。[対象比較](validation/enum-method-targeted-attempt.log)と[全件比較](validation/enum-method-compare-all-attempt.log)はPodman接続拒否のため終了2で、各版0一致・0不一致・1環境エラー、10件未実施。期待値・診断は暫定。項目8および項目6〜10全体は未完了。Macでの版別コマンドは下のenum基本経路節に記したものを使用する。

## enum基本経路の追加（2026-10-16、884件時点の履歴）

宣言情報にbacking型とcase宣言順を持ち、実行単位のクラス定数表へcaseごとの単一`phpObject`を登録した。配列格納・引数渡しでも同一性を維持する。caseを参照する基本的なクラス定数別名も処理する。`UnitEnum`／`BackedEnum`を組み込みinterfaceとして型判定へ接続し、`cases()`、backed enumの`from()`／`tryFrom()`、`name`／`value`、基本的なinterface・trait・名前空間・include、`new`とプロパティ宣言の拒否を追加した。

追加**6件**で総数**884件**（正常475・実行時279・構文130）。既存878件は[従来の版別集計](validation/readonly-878-after-fix-summary-mac.log)で3版とも一致済み。新規6件はGo期待値のみで暫定。[対象比較](validation/enum-targeted-attempt.log)と[全件比較試行](validation/enum-compare-all-attempt.log)はPodman接続拒否で各終了2。各版0一致・0不一致・1環境エラー、6件未実施。[Podman確認](validation/enum-podman.log)は終了125。ローカルの[test](validation/enum-test.log)、[vet](validation/enum-vet.log)、[build](validation/enum-build.log)、[差分検査](validation/enum-diff-check.log)を参照。

残件：backing値がクラス定数を参照する場合を含む定数式の相互依存・遅延解決、backed値重複の正確な拒否時点、組み込みAPIの共通引数束縛・トレースと診断、enum固有の比較／変換境界、trait経由の禁止メンバー、禁止マジックメソッドの網羅、未選択分岐とinclude先の宣言エラーの公式比較。項目8と項目6〜10全体を完了扱いにしない。固定vendorでの使用確認は[依存監査](dependency-audit.md)に記録し、Go実行到達は未確認。

Macでの再比較（終了コードはコマンド直後に保存）：

```sh
cd /Users/kouji/dev/php
set -o pipefail
make build
make compare PHP=./containers/php-version/bin/php83 PHP_VERSION=8.3.35 PHP_TARGET=8.3 2>&1 | tee docs/validation/enum-compare-83-mac.log
result=$?; printf 'exit_code=%s\n' "$result" >> docs/validation/enum-compare-83-mac.log
make compare PHP=./containers/php-version/bin/php84 PHP_VERSION=8.4.26 PHP_TARGET=8.4 2>&1 | tee docs/validation/enum-compare-84-mac.log
result=$?; printf 'exit_code=%s\n' "$result" >> docs/validation/enum-compare-84-mac.log
make compare PHP=./containers/php-version/bin/php85 PHP_VERSION=8.5.11 PHP_TARGET=8.5 2>&1 | tee docs/validation/enum-compare-85-mac.log
result=$?; printf 'exit_code=%s\n' "$result" >> docs/validation/enum-compare-85-mac.log
```

## readonly追補の878件が版ごとの公式PHP比較で一致（2026-10-16）

ユーザーのMacで版ごとに`make compare`を実行した後の[6行の集計ログ](validation/readonly-878-after-fix-summary-mac.log)を`/tmp/php-compare-878-after-readonly.log`から内容を変えず保存し、`cmp`でバイト一致を確認した。PHP 8.3.35・8.4.26・8.5.11でそれぞれ**878 matched・0 different・0 errors**。851件から追加した27件の公式比較待ちは、**収録fixtureの範囲で解消**した。これは保存済み集計の確認であり、今回比較を再実行した結果ではない。`make compare-all`は使われていない。集計ログには**個別fixture行と各`make compare`の終了コードがない**ため、指摘2件の解消と`check.sh`の各版PASSはユーザー報告として区別する。修正後の版別生ログは、この環境では見つからなかった。修正前の876一致・1不一致・1テストエラーは下節の履歴として維持する。項目8と項目6〜10全体は部分実装で、clone・参照・フック等は残る。

## readonly traitプロパティ・親private名の追補（2026-10-16、修正直後の履歴）

ユーザーがVM再起動後に各版を個別比較した[修正前の878件ログ](validation/readonly-878-before-fix-mac.log)は、8.3.35・8.4.26・8.5.11の各版で**876 matched・1 different・1 error**。`check.sh`の3版PASSはユーザー報告であり、この集計ログ自体にコマンド終了コードはない。`ERROR`は環境障害ではなく、成功fixture`readonly_class_trait_property.php`が公式PHPで終了255になったための分類である。`DIFF`は`readonly_private_child_scope.php`。この時点の比較は全件一致ではない。

Go側ではtraitから取り込む非readonlyプロパティをreadonly classの宣言時に拒否し、前者を`runtime/`へ移した。親のprivateプロパティは子メソッドから見えないため、同名への書き込みは別の動的プロパティとして扱い、作成時に非推奨診断を出す。後者もstderrを伴うので`runtime/`へ移した。移動以外のfixture総数は**878件**（正常471・実行時278・構文129）。[make test](validation/readonly-878-followup-test.log)、[vet](validation/readonly-878-followup-vet.log)、[build](validation/readonly-878-followup-build.log)は終了0。修正後の[trait対象比較](validation/readonly-878-followup-trait-targeted-attempt.log)・[private対象比較](validation/readonly-878-followup-private-targeted-attempt.log)・[全件比較](validation/readonly-878-followup-compare-all-attempt.log)は、こちらのPodman接続拒否により各終了2。**修正後の公式PHP3版比較は未実施**で、今回の診断期待値は暫定である。

Macでの再比較（各コマンドの終了コードを直後に保存）：

```sh
cd /Users/kouji/dev/php
set -o pipefail
make build
./containers/php-version/compare-all.sh tests/fixtures/runtime/readonly_class_trait_property.php ./bin/gophp 2>&1 | tee docs/validation/readonly-878-followup-trait-targeted-mac.log
trait_exit=$?
echo "trait targeted exit=$trait_exit" | tee -a docs/validation/readonly-878-followup-trait-targeted-mac.log
./containers/php-version/compare-all.sh tests/fixtures/runtime/readonly_private_child_scope.php ./bin/gophp 2>&1 | tee docs/validation/readonly-878-followup-private-targeted-mac.log
private_exit=$?
echo "private targeted exit=$private_exit" | tee -a docs/validation/readonly-878-followup-private-targeted-mac.log
make compare-all 2>&1 | tee docs/validation/readonly-878-followup-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/readonly-878-followup-compare-all-mac.log
```

## readonly class継承検査の修正（2026-10-16、追加直後の履歴）

通常クラスがreadonlyな親クラスを継承する場合も、宣言時に`Non-readonly class ... cannot extend readonly class ...`で拒否する。以前の877件時点ではこの方向を許していた。PHPの[クラス解説](https://www.php.net/manual/en/language.oop5.basic.php)に加え、[readonly amendments RFC](https://wiki.php.net/rfc/readonly_amendments)の提案1は**反対12・賛成7で否決**され、[実装案のPR](https://github.com/php/php-src/pull/9827)も不採用だったことを確認した。RFC全体の「Implemented」は採択された別の提案を含む状態であり、提案1の採択を意味しない。

旧`success/readonly_parent_mutable_child.php`を`parser/readonly_parent_mutable_child.php`へ移し、readonly親子の合法な`success/readonly_parent_readonly_child.php`を追加した。fixtureは**878件**（正常473・実行時276・構文129）、851件の比較済み基準から27件追加。Goの3設定の期待値検査は通過したが、[対象比較試行](validation/readonly-inheritance-targeted-attempt.log)と[全件比較試行](validation/readonly-inheritance-compare-all-attempt.log)はPodman接続拒否により各終了2。各版で**今回0一致・0不一致・1環境エラー、27件未実施**。宣言診断の期待値は公式PHP実比較まで暫定である。[make test](validation/readonly-inheritance-test.log)、[vet](validation/readonly-inheritance-vet.log)、[build](validation/readonly-inheritance-build.log)は終了0。項目8・項目6〜10全体は未完了。

Macでの対象比較と全件比較：

```sh
cd /Users/kouji/dev/php
set -o pipefail
make build
./containers/php-version/compare-all.sh tests/fixtures/parser/readonly_parent_mutable_child.php ./bin/gophp 2>&1 | tee docs/validation/readonly-inheritance-targeted-mac.log
target_exit=$?
echo "targeted comparison exit=$target_exit" | tee -a docs/validation/readonly-inheritance-targeted-mac.log
make compare-all 2>&1 | tee docs/validation/readonly-inheritance-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/readonly-inheritance-compare-all-mac.log
```

## readonly・コンストラクタプロパティ昇格（2026-10-16、877件時点の履歴）

昇格引数を既存の仮引数と宣言済みプロパティの両方として保持し、引数束縛・型検査後、コンストラクタ本体より先に一度だけ格納する。既定値は仮引数側だけに持つ。readonlyは各プロパティの宣言情報とreadonly classのフラグに保持し、共通の`storeObjectProperty`で未初期化からの初回代入、2回目の書き込み、型不一致後の状態を扱う。代入箇所のstrict_types、privateスロット分離、動的名アクセス、配列の値コピーは既存経路を使う。8.3の初期化は宣言元スコープ、8.4・8.5では公開／保護プロパティを子クラスからも初期化可能とする。当時はRFC提案本文を採択済みと誤読し、通常クラスがreadonly親クラスを継承できると記録した。この記述は誤りで、上節の修正を現在の状態とする。[PHPプロパティの公式説明](https://www.php.net/manual/en/language.oop5.properties.php)、[readonly amendments RFC](https://wiki.php.net/rfc/readonly_amendments)、[昇格のRFC](https://wiki.php.net/rfc/constructor_promotion)を仕様根拠とした。

新規**26件**（正常16、構文失敗10）を加え、fixtureは**877件**（正常473・実行時276・構文128）。既存851件のPHPソース・期待値は維持した。[Go・Pythonテスト](validation/readonly-promotion-877-test.log)、[vet](validation/readonly-promotion-877-vet.log)、[ビルド](validation/readonly-promotion-877-build.log)、[差分検査](validation/readonly-promotion-877-diff-check.log)は各終了0。[Podman確認](validation/readonly-promotion-podman.log)は終了125、[対象比較](validation/readonly-promotion-877-targeted-attempt.log)と[全件比較](validation/readonly-promotion-877-compare-all-attempt.log)は各終了2。各版**今回0一致・0不一致・1環境エラー、新規26件未実施**で、診断文と期待値は暫定。項目8・項目6〜10全体は未完了。

追加fixture一覧：`success/promotion_untyped.php`, `success/readonly_child_scope.php`, `success/readonly_class_basic.php`, `success/readonly_class_trait_property.php`, `success/readonly_include_finally.php`, `success/readonly_parent_mutable_child.php`, `success/readonly_private_child_scope.php`, `success/readonly_promotion_basic.php`, `success/readonly_promotion_composite.php`, `success/readonly_promotion_inheritance_trait.php`, `success/readonly_promotion_parent_call.php`, `success/readonly_property_operations.php`, `success/readonly_type_failure_state.php`, `success/readonly_uninitialized_coalesce.php`, `success/readonly_uninitialized_unset.php`, `success/readonly_write_forms.php`, `parser/promotion_duplicate.php`, `parser/promotion_implicit_nullable.php`, `parser/promotion_outside_constructor.php`, `parser/promotion_variadic.php`, `parser/readonly_class_inheritance.php`, `parser/readonly_class_unselected_invalid.php`, `parser/readonly_class_untyped.php`, `parser/readonly_static.php`, `parser/readonly_unselected_invalid.php`, `parser/readonly_untyped.php`。

残件：参照付き昇格、clone内のreadonly再初期化、未初期化readonlyのunsetの未収録境界、PHP 8.4以降の非対称可視性による明示的なset権限、プロパティフック、enum、Reflection、魔法メソッド。型付きクラスのtrait・継承に関する宣言診断と未選択分岐の時点も、公式PHP比較で確認する。これらを対応済みとしない。

Macでの再比較（稼働コンテナは再起動しない）：

```sh
cd /Users/kouji/dev/php
podman ps -a --format '{{.Names}} {{.Status}}'
for name in php-version-83 php-version-84 php-version-85; do
  if [ "$(podman inspect -f '{{.State.Running}}' "$name")" = false ]; then podman start "$name"; fi
done
set -o pipefail
make build
./containers/php-version/compare-all.sh tests/fixtures/success/readonly_promotion_basic.php ./bin/gophp 2>&1 | tee docs/validation/readonly-promotion-targeted-mac.log
target_exit=$?
echo "target exit=$target_exit" | tee -a docs/validation/readonly-promotion-targeted-mac.log
make compare-all 2>&1 | tee docs/validation/readonly-promotion-compare-all-mac.log
comparison_exit=$?
echo "compare-all exit=$comparison_exit" | tee -a docs/validation/readonly-promotion-compare-all-mac.log
# compare-allの途中で8.5だけPodman終了125になった場合は、原因を確認した後に別実行として記録する。
PHP=./containers/php-version/bin/php85 PHP_VERSION=8.5.11 make compare 2>&1 | tee docs/validation/readonly-promotion-85only-mac.log
php85_exit=$?
echo "8.5-only exit=$php85_exit" | tee -a docs/validation/readonly-promotion-85only-mac.log
```

直前の851件は、ユーザーのMacでの[全件比較ログ](validation/readonly-baseline-851-compare-all-mac.log)で8.3.35・8.4.26が各851一致、8.5.11が830一致／21環境エラー、コマンド終了2だった。VM再起動後の[8.5単独ログ](validation/readonly-baseline-851-85only-mac.log)は851一致、終了0。元ログと保存先の`cmp`によるバイト一致を確認した。8.5のPodman資源枯渇は未確認の推測であり、恒久対策を今回の言語実装へ混ぜない。以下の動的プロパティ比較待ち・失敗記録は当時の履歴である。

## 動的プロパティの受け手再束縛を修正（851件時点の履歴）

ユーザーがMacで実行した[修正前の851件全件比較](validation/dynamic-property-851-before-fix-mac.log)を元ログ`/tmp/php-compare-latest7.log`から内容を変えず保存し、`cmp`でバイト一致を確認した。`make compare-all`終了2。PHP 8.3.35・8.4.26は各**850一致・1不一致・0エラー**、8.5.11は**829一致・1不一致・21環境エラー**。8.5の21件はPodmanの終了125であり、fixture不一致ではない。不一致は3版共通の`success/dynamic_property_receiver_change.php`で、公式stdoutは`1|9|1|10`、Goは修正前に`9|2|10|2`を出した。

原因は`prepareTarget`が動的名のプロパティ代入で受け手オブジェクトを名前式・右辺より先に固定したこと。変数の受け手は実際の書き込み時に解決するよう修正し、同fixtureの期待stdoutを公式PHPの実結果へ訂正した。PHPソースと他の期待値は維持した。Goの3設定では対象が`1|9|1|10`、[Go・Pythonテスト](validation/dynamic-property-receiver-fix-test.log)・[vet](validation/dynamic-property-receiver-fix-vet.log)・[ビルド](validation/dynamic-property-receiver-fix-build.log)・[差分検査](validation/dynamic-property-receiver-fix-diff-check.log)は各終了0。こちらの[Podman確認](validation/dynamic-property-receiver-fix-podman.log)は終了125、修正後の[対象比較](validation/dynamic-property-receiver-fix-targeted-attempt.log)と[全件比較](validation/dynamic-property-receiver-fix-compare-all-attempt.log)は各終了2で、各版とも**修正後0一致・0不一致・1環境エラー、全851件未実施**。この修正を公式比較済みとは扱わない。

Macでの再比較（既存コンテナの稼働を確認し、停止中のものだけ起動する）：

```sh
cd /Users/kouji/dev/php
podman ps -a --format '{{.Names}} {{.Status}}'
for name in php-version-83 php-version-84 php-version-85; do
  if [ "$(podman inspect -f '{{.State.Running}}' "$name")" = false ]; then podman start "$name"; fi
done
set -o pipefail
make build
./containers/php-version/compare-all.sh tests/fixtures/success/dynamic_property_receiver_change.php ./bin/gophp 2>&1 | tee docs/validation/dynamic-property-receiver-fix-targeted-mac.log
target_exit=$?
echo "target exit=$target_exit" | tee -a docs/validation/dynamic-property-receiver-fix-targeted-mac.log
make compare-all 2>&1 | tee docs/validation/dynamic-property-receiver-fix-compare-all-mac.log
comparison_exit=$?
echo "compare-all exit=$comparison_exit" | tee -a docs/validation/dynamic-property-receiver-fix-compare-all-mac.log
```

## 動的インスタンスプロパティ名の部分実装（追加直後の記録・履歴）

`$object->$name`・`$object->{$name}`・`$object->{式}`を構文木の名前式として保持し、`?->`にも接続した。読み取り・代入・複合代入・増減・`isset`・`empty`・`??`・`??=`・`unset`・配列要素操作が、固定名と同じ`propertyKey`、可視性・privateスロット・型付き格納経路を通る。当初は代入先の受け手を名前式と右辺より先に一度保持したが、上記のMac実比較で不一致が判明し、後続の修正で撤回した。nullsafe短絡時は名前式と後続添字を実行しない。[PHPの変数プロパティ名の説明](https://www.php.net/manual/en/language.variables.variable.php)を構文の参照資料とした。この節の比較待ちは追加直後の状態を指す。

新規fixtureは`success/dynamic_property_*.php`14件、実行時失敗の`runtime/dynamic_property_*.php`2件と警告を出す`runtime/dynamic_property_creation.php`1件、`parser/dynamic_property_nullsafe_write.php`1件の**18件**。総数は**851件**（正常457・実行時276・構文118）。既存833件のソースと期待値は変更していない。[Go・Pythonテスト](validation/dynamic-property-test.log)・[vet](validation/dynamic-property-vet.log)・[ビルド](validation/dynamic-property-build.log)・[差分検査](validation/dynamic-property-diff-check.log)は終了0。[Podman状態確認](validation/dynamic-property-podman.log)はソケット接続拒否で終了125、[対象比較](validation/dynamic-property-targeted-attempt.log)と[全件比較](validation/dynamic-property-compare-all-attempt.log)は終了2。各版とも**今回0一致・0不一致・1環境エラー、851件未実施**である。新規18件の期待値は暫定で、Goの期待値テストを公式PHPとの一致と扱わない。直前に全件一致した833件の記録は次節の[成功ログ](validation/typed-members-833-retry-mac.log)にある。

未対応・未確認：動的メソッド呼び出し、`__get`等のマジックメソッド、参照経由、readonly、プロパティ昇格、プロパティフック、非対称可視性、および名前式に配列・オブジェクト等を使う場合の診断や未収録の副作用順序。項目7・8および項目6〜10全体は部分実装のまま。次はMacで対象と全件を比較し、差があれば最小再現例で修正する。既存コンテナが停止している場合だけ起動し、比較コマンド自身の終了コードをパイプ直後に保存する。

```sh
cd /Users/kouji/dev/php
podman ps -a --format '{{.Names}} {{.Status}}'
for name in php-version-83 php-version-84 php-version-85; do
  if [ "$(podman inspect -f '{{.State.Running}}' "$name")" = false ]; then podman start "$name"; fi
done
set -o pipefail
make build
./containers/php-version/compare-all.sh tests/fixtures/success/dynamic_property_basic.php ./bin/gophp 2>&1 | tee docs/validation/dynamic-property-targeted-mac.log
comparison_exit=$?; echo "target exit=$comparison_exit" | tee -a docs/validation/dynamic-property-targeted-mac.log
make compare-all 2>&1 | tee docs/validation/dynamic-property-compare-all-mac.log
comparison_exit=$?; echo "compare-all exit=$comparison_exit" | tee -a docs/validation/dynamic-property-compare-all-mac.log
```

## 型付きメンバーの833件が公式PHP3版と一致（2026-10-16）

ユーザーがMacで実行した[修正後の全件比較](validation/typed-members-833-retry-mac.log)を元の`/private/tmp/php-compare-833-retry.log`とバイト単位で照合して保存した。PHP 8.3.35・8.4.26・8.5.11の各版で**833 matched・0 different・0 errors**、ログ末尾の終了コード0。fixture内訳は正常443・実行時273・構文117。817件以降の新規16件と型付きメンバーの診断修正2件は、収録fixtureの範囲で公式比較済みとなった。これは保存済みログの確認で、今回比較を再実行した結果ではない。動的インスタンスプロパティアクセス`$object->{$name}`、readonly、プロパティ昇格、プロパティフック、非対称可視性、enum、参照経由の型制約、未登録型の解決などは残る。項目8全体と項目6〜10全体は未完了。

### 修正前と途中停止の記録（履歴）

修正前の[Mac全件比較ログ](validation/typed-members-833-before-fix-mac.log)は、PHP 8.3.35・8.4.26・8.5.11の各版で**831 matched・2 different・0 errors**、終了コード2だった。相違は`parser/typed_constant_covariance.php`と`parser/typed_property_trait_conflict.php`のstderrだけで、stdout・行番号・終了コードは一致した。診断修正後の[最初の実行](validation/typed-members-833-after-mac.log)は8.3・8.4が各833一致、8.5が832一致・0不一致・1環境エラーで終了2。8.5の公式PHP実行と版確認が15秒でタイムアウトしたもので、fixture差異ではない。上記の再試行で833件全一致した。

定数型の継承診断を`Type of Child::VALUE must be compatible with Parent::VALUE of type int`形式へ変更し、traitとクラスの型付き同名プロパティの競合診断に`However, the definition differs and is considered incompatible. Class was composed`を追加した。該当2件の暫定期待値は報告された公式PHPの出力に合わせて訂正した。[Go・Pythonテスト](validation/typed-member-diagnostic-test.log)、[vet](validation/typed-member-diagnostic-vet.log)、[ビルド](validation/typed-member-diagnostic-build.log)は終了0。当時のこちらの環境ではPodmanが[接続拒否](validation/typed-member-diagnostic-podman.log)・終了125で、[対象比較試行](validation/typed-member-diagnostic-targeted-attempt.log)と[全件比較試行](validation/typed-member-diagnostic-compare-all-attempt.log)は終了2だった。後続のMacでの再試行により収録範囲の比較待ちは解消した。以下は当時用意した再実行コマンドである。

```sh
cd /Users/kouji/dev/php
set -o pipefail
make build
./containers/php-version/compare-all.sh tests/fixtures/parser/typed_constant_covariance.php ./bin/gophp 2>&1 | tee docs/validation/typed-constant-covariance-mac.log
comparison_exit=$?; echo "target exit=$comparison_exit" | tee -a docs/validation/typed-constant-covariance-mac.log
./containers/php-version/compare-all.sh tests/fixtures/parser/typed_property_trait_conflict.php ./bin/gophp 2>&1 | tee docs/validation/typed-property-trait-conflict-mac.log
comparison_exit=$?; echo "target exit=$comparison_exit" | tee -a docs/validation/typed-property-trait-conflict-mac.log
make compare-all 2>&1 | tee docs/validation/typed-member-diagnostic-compare-all-mac.log
comparison_exit=$?; echo "compare-all exit=$comparison_exit" | tee -a docs/validation/typed-member-diagnostic-compare-all-mac.log
```

## 型付きプロパティの複合代入戻り値を追加した時点（履歴・2026-10-16）

`internal/php/arithmetic.go`の`evalCompound`が、保存先の型検査・変換後にも保存前の演算値を返していた。成功した書き込みの後に、評価済みの格納先から最終値を取得するよう変更した。[PHPの型付きプロパティRFC](https://wiki.php.net/rfc/typed_properties_v2)は、代入式が型変換後の格納値を返す規則を示す。`success/typed_property_compound_result.php`と`success/typed_static_property_compound_result.php`を追加し、`int`プロパティへの`.=`について、式の結果と格納値がともに`integer|12`となることをGoの3設定で確認した。

fixtureは**833件**（正常443・実行時273・構文117）。817件以降の新規16件は[Go・Pythonテスト](validation/typed-compound-test.log)で成功し、[vet](validation/typed-compound-vet.log)・[ビルド](validation/typed-compound-build.log)・[差分検査](validation/typed-compound-diff-check.log)も終了0。当時の公式PHP3版の[新規2件の対象比較試行](validation/typed-compound-targeted-attempt.log)と[全件比較試行](validation/typed-compound-compare-all-attempt.log)はPodman接続拒否で各終了2。[Podman状態確認](validation/typed-compound-podman.log)は終了125。この時点では各版0一致・0不一致・1環境エラー、833件未実施で期待値は暫定だった。後続の冒頭の成功ログで収録範囲の比較待ちは解消した。

Macでの対象比較は、既存コンテナが停止している場合だけ起動した後に次を実行する。`tee`の直後に元コマンドの終了コードを保存する。

```sh
cd /Users/kouji/dev/php
podman ps -a --format '{{.Names}} {{.Status}}'
for name in php-version-83 php-version-84 php-version-85; do
  if [ "$(podman inspect -f '{{.State.Running}}' "$name")" = false ]; then podman start "$name"; fi
done
set -o pipefail
make build
./containers/php-version/compare-all.sh tests/fixtures/success/typed_property_compound_result.php ./bin/gophp 2>&1 | tee docs/validation/typed-compound-property-mac.log
comparison_exit=$?; echo "property exit=$comparison_exit" | tee -a docs/validation/typed-compound-property-mac.log
./containers/php-version/compare-all.sh tests/fixtures/success/typed_static_property_compound_result.php ./bin/gophp 2>&1 | tee docs/validation/typed-compound-static-mac.log
comparison_exit=$?; echo "static exit=$comparison_exit" | tee -a docs/validation/typed-compound-static-mac.log
make compare-all 2>&1 | tee docs/validation/typed-compound-compare-all-mac.log
comparison_exit=$?; echo "compare-all exit=$comparison_exit" | tee -a docs/validation/typed-compound-compare-all-mac.log
```

## 型付きプロパティ・型付きクラス定数を追加した時点（履歴・831件）

817件の成功記録を維持したまま14件を追加し、当時**831件**（正常441・実行時273・構文117）。[Go・Pythonテスト](validation/typed-members-test.log)、[vet](validation/typed-members-vet.log)、[ビルド](validation/typed-members-build.log)は各終了0。公式PHP3版の[対象比較](validation/typed-members-targeted-attempt.log)と[全件比較](validation/typed-members-compare-all-attempt.log)はそれぞれ終了2。各版の版確認で`podman exec`が接続拒否・終了125となったため、各版は**今回0一致・0不一致・1環境エラー・831件未実施**であり、追加fixtureの期待値と診断は暫定である。[Podman状態確認](validation/typed-members-podman-status.log)も終了125。以下の817件比較結果は追加前の履歴として維持する。

[PHPの型宣言](https://www.php.net/manual/en/language.types.declarations.php)、[プロパティ](https://www.php.net/manual/en/language.oop5.properties.php)、[クラス定数](https://www.php.net/manual/en/language.oop5.constants.php)を仕様参照とした。ただし今回の新規fixtureは実比較前である。`propertyDecl`／`classConstantDecl`に解決済みの型名とnullable情報を保持する。型付きインスタンスプロパティは初期値がない場合に配列スロットを作らず、`unset`後も同じ未初期化状態へ戻る。型付きstaticプロパティは継承時に共有するスロットに初期化フラグを持つ。`isset`・`empty`・`??`は未初期化を静かに扱い、通常読み取りと複合更新は`Error`とする。書き込みは格納前に代入元ファイルの`strict_types`を使い、弱い型変換・精度喪失の非推奨診断後に反映する。クラス定数は宣言値を型検査し、intからfloatへの特例だけ変換する。親子の通常プロパティ型は双方向包含で不変、定数型は共変として比較する。traitの同名メンバーの型も比較する。

対応は既存の型パーサーが受け付ける型に限る。動的staticプロパティ名の既存形式`Class::${$name}`は共通の型検査へつながるが、インスタンスの`$object->{$name}`は構文解析が未対応で明示的に拒否する。readonly、プロパティ昇格、プロパティフック、非対称可視性、enum、参照経由の型制約、マジックメソッド、Reflection、autoloadを含む未登録型の解決は未対応。固定vendorの[API Platform Collection](../research/laravel-api/vendor/api-platform/hydra/Collection.php)には初期値なしの`public int $totalItems`、[Symfony ServiceLocatorTrait](../research/laravel-api/vendor/symfony/service-contracts/ServiceLocatorTrait.php)には`private array $loading = []`がある。固定vendor内の型付きクラス定数は今回の対象検索では確認できなかった。vendorの実行経路到達は未確認。次の実装単位は公式比較で判明する差異の修正、その後にreadonly・プロパティ昇格である。項目8全体、項目6〜10全体は未完了。

Macでの比較手順（停止した既存コンテナだけ起動し、VMが起動済みなら再起動しない）：

```sh
cd /Users/kouji/dev/php
podman ps -a --format '{{.Names}} {{.Status}}'
for name in php-version-83 php-version-84 php-version-85; do
  if [ "$(podman inspect -f '{{.State.Running}}' "$name")" = false ]; then podman start "$name"; fi
done
set -o pipefail
make build
./containers/php-version/compare-all.sh tests/fixtures/success/typed_property_basics.php ./bin/gophp 2>&1 | tee docs/validation/typed-members-targeted-mac.log
comparison_exit=$?; echo "targeted exit=$comparison_exit" | tee -a docs/validation/typed-members-targeted-mac.log
make compare-all 2>&1 | tee docs/validation/typed-members-compare-all-mac.log
comparison_exit=$?; echo "compare-all exit=$comparison_exit" | tee -a docs/validation/typed-members-compare-all-mac.log
```

## 817件比較時点（履歴）：公式PHP3版で一致

ユーザーがMacで実行した`/tmp/php-compare-817-after.log`を[修正後の全件比較ログ](validation/void-never-817-after-fix-mac.log)に内容を変えず保存し、`cmp`でバイト一致を確認した。PHP 8.3.35・8.4.26・8.5.11は各**817 matched・0 different・0 errors**、`make compare-all`終了0。`never_nullable.php`、`void_nullable.php`、`void_constructor.php`、`namespace_new_fqcn.php`、`void_never_namespace_case.php`も3版でMATCHした。787件から追加した30件の公式比較待ちは収録fixtureの範囲で解消した。今回は保存済みログの確認・保存であり、こちらで比較を再実行していない。項目8全体・項目6〜10全体、未登録型の後続解決、型付きプロパティ・定数、Go製PHPでのLaravel起動などは未完了である。

## 診断3件の修正直後（履歴・当時は全件再比較待ち）

ユーザーがMacで実行した`/tmp/php-compare-latest5.log`を[修正前の公式比較ログ](validation/void-never-817-before-fix-mac.log)へバイト一致で保存した。PHP 8.3.35・8.4.26・8.5.11はいずれも**814 matched・3 different・0 errors**、`make compare-all`終了2。差異は`parser/never_nullable.php`、`parser/void_nullable.php`、`parser/void_constructor.php`のstderrのみだった。公式PHPは`?never`を`Fatal error: never can only be used as a standalone type`、`?void`を`Fatal error: Void can only be used as a standalone type`、コンストラクタの戻り値型を`Fatal error: Method Invalid::__construct() cannot declare a return type`と報告する。8.5では各診断に`Stack trace:\n#0 {main}\n`が付く。

`function.go`でnullableな`void`／`never`を宣言時の致命エラーとし、`class.go`のコンストラクタ診断を`Method`へ変更した。既存の3件の期待値だけを[公式PHPログ](validation/void-never-817-before-fix-mac.log)に基づき訂正した。[保存済み公式出力と修正後Go出力の照合](validation/void-never-diagnostic-targeted-local.log)は3件×3版のstdout・stderr・終了コードで一致し、終了0。[`make test`](validation/void-never-diagnostic-test.log)、[`make vet`](validation/void-never-diagnostic-vet.log)、[`make build`](validation/void-never-diagnostic-build.log)も終了0。修正後の[全件比較試行](validation/void-never-diagnostic-compare-all-attempt.log)はPodman接続拒否で終了2となり、各版の版確認は終了125。**修正後817件の全件一致は未確認**であり、Macでの再比較には次節の`make compare-all`手順を使う。過去の787件成功記録は別の履歴である。

## 名前空間内の完全修飾クラス生成の修正時点（履歴）

式の先頭の`new`と後続の`\Exception`を、空白・改行・コメントで区切っても単一の名前として連結していた。`lexer.go`にトークン開始位置を保持し、`namespace.go`では離れた`new`と`\`を結合せず、`new`キーワードと完全修飾クラス名として解析するよう修正した。PHPの[名前空間の公式説明](https://www.php.net/manual/en/language.namespaces.basics.php)に`new \Exception`の使用例がある。先行する`success/void_never_namespace_case.php`を直接の`new \Exception`へ戻し、新たに`success/namespace_new_fqcn.php`で空白・改行・コメントを挟む形式を検証する。既存の816件の期待値は変更していない。

fixtureは**817件**（正常433・実行時272・構文112）。[Go・Pythonテスト](validation/namespace-new-test.log)は終了0。新規1件の期待値と、先行29件の公式PHP比較は暫定である。[対象比較試行](validation/namespace-new-targeted-attempt.log)・[全件比較試行](validation/namespace-new-compare-all-attempt.log)は各終了2で、PHP 8.3.35・8.4.26・8.5.11の版確認がPodmanソケット接続拒否（`podman exec`終了125）となった。各版**今回0一致・0不一致・1環境エラー・817件未実施**。以前の787件成功記録は有効な履歴として維持する。Macでは以下を実行し、既存コンテナが停止している場合だけ起動する。

```sh
cd /Users/kouji/dev/php
podman ps -a --format '{{.Names}} {{.Status}}'
podman start php-version-83 php-version-84 php-version-85 # 停止中のものだけ指定
set -o pipefail
make build
./containers/php-version/compare-all.sh tests/fixtures/success/namespace_new_fqcn.php ./bin/gophp 2>&1 | tee docs/validation/namespace-new-targeted-mac.log
comparison_exit=$?; echo "targeted comparison exit=$comparison_exit" | tee -a docs/validation/namespace-new-targeted-mac.log
make compare-all 2>&1 | tee docs/validation/namespace-new-compare-all-mac.log
comparison_exit=$?; echo "make compare-all exit=$comparison_exit" | tee -a docs/validation/namespace-new-compare-all-mac.log
```

## void／never戻り値型を追加した時点の記録（履歴・公式PHP比較待ち）

既存の[787件成功ログ](validation/type-variance-787-compare-all-mac.log)はPHP 8.3.35・8.4.26・8.5.11で各787一致・0不一致・0エラー、終了0。今回29件を加え、**816件**（正常432・実行時272・構文112）。新規29件の期待値・診断はGoローカル実行と[PHPのvoid説明](https://www.php.net/manual/en/language.types.void.php)、[never説明](https://www.php.net/manual/en/language.types.never.php)、[PHP 8.3.35のコンパイラ実装](https://fossies.org/windows/www/php-8.3.35-src.zip/Zend/zend_compile.c)、[PHP 8.4のメソッド診断変更](https://github.com/php/php-src/commit/d84ed03a68b1e12d6e269125e29b6a4a26bd0204)から置いた**暫定値**。公式PHP3版の実比較は未実施であり、816件一致とは記載しない。

`void`は関数・メソッド・staticメソッド・クロージャで`return;`または末尾到達を許し、呼び出し結果を`null`とする。値付き`return`（`return null`を含む）は未選択分岐でも実行前に拒否する。`never`は例外を投げる経路を許し、明示的な`return`を実行前に拒否し、末尾到達は`TypeError`とする。アロー関数の`void`は値を返すため拒否し、`never`のthrow式を許す。入れ子の関数・クロージャのreturnは外側の型検査から分離した。継承時は`never`を戻り値の最小型、`void`を値型や`null`と区別し、型なし親の子に`void`を許す。メソッドの`void`／`never`診断は8.3と8.4以降で分ける。`__construct`の戻り値型は拒否する。`exit`／`die`、型付きプロパティ・定数、未登録型の後続解決等は残る。項目8と項目6〜10全体は未完了。

追加fixture：正常系`void_never_basic`、`void_never_nested`、`void_never_finally`、`void_never_variance`、`void_never_include`、`void_never_methods`、`void_never_variance_edges`、`void_never_namespace_case`、`never_implicit_catch`。実行時`never_implicit`、`never_method_implicit`、`void_variance_mixed`、`never_variance_widen`。構文`void_return_value`、`void_return_null`、`never_return`、`never_return_unselected`、`void_parameter`、`never_parameter`、`void_union`、`never_union`、`void_nullable`、`never_nullable`、`void_arrow`、`void_constructor`、`void_method_return`、`never_method_return`、`void_intersection`、`never_intersection`。補助ファイル`support/void_never_defs.inc`はfixture件数に含めない。元の787件は変更していない。

[`make test`](validation/void-never-continuation-test.log)、[`make vet`](validation/void-never-continuation-vet.log)、[`make build`](validation/void-never-continuation-build.log)は816件時点で終了0。追加前の[対象比較試行](validation/void-never-targeted-attempt.log)と[全件比較試行](validation/void-never-compare-all-attempt.log)は各終了2。いずれも3版の版確認でPodman接続拒否・`podman exec`終了125となった。816件時点の再試行は下記に記録する。これまでの787件の成功履歴と区別する。

816件時点の[対象比較試行](validation/void-never-continuation-targeted-attempt.log)と[全件比較試行](validation/void-never-continuation-compare-all-attempt.log)も、各版の版確認がPodmanソケット接続拒否（`podman exec`終了125）となり、コマンドは各終了2だった。各版**今回0一致・0不一致・1環境エラー・816件未実施**。対象比較の`void_never_namespace_case.php`も未実施だった。当時は`new \Exception`を現Goパーサーが`ReturnTypes\new\Exception()`と解釈する未修正問題があり、fixtureでは`use Exception`を使っていた。上記の修正後は直接の完全修飾名へ戻した。

Macでの再比較は次のとおり。VMが起動済みなら再度起動せず、停止した既存コンテナだけ起動する。`tee`直後に元コマンドの終了コードを保存する。

```sh
cd /Users/kouji/dev/php
podman ps -a --format '{{.Names}} {{.Status}}'
podman start php-version-83 php-version-84 php-version-85 # 停止している場合のみ
set -o pipefail
make build
./containers/php-version/compare-all.sh tests/fixtures/success/void_never_namespace_case.php ./bin/gophp 2>&1 | tee docs/validation/void-never-continuation-targeted-mac.log
comparison_exit=$?; echo "targeted comparison exit=$comparison_exit" | tee -a docs/validation/void-never-continuation-targeted-mac.log
make compare-all 2>&1 | tee docs/validation/void-never-continuation-compare-all-mac.log
comparison_exit=$?; echo "make compare-all exit=$comparison_exit" | tee -a docs/validation/void-never-continuation-compare-all-mac.log
```

以下の「現在」と書かれた節は各節の記録時点の状態を指す履歴である。

## 787件時点：複合型の継承署名、公式PHP3版で一致（履歴）

ユーザーがMacで実行した`/tmp/php-compare-latest4.log`を[成功ログ](validation/type-variance-787-compare-all-mac.log)へ内容を変えずに保存し、`cmp`とSHA-256でバイト一致を確認した。PHP 8.3.35・8.4.26・8.5.11は各**787 matched・0 different・0 errors**、ログ末尾の終了コードは0。`variance_float_return.php`、`variance_float_param.php`、2件のunion不互換、実行時変換の正常系も各版でMATCHした。769件から追加した18件の公式比較待ちは、収録fixtureの範囲で解消した。今回は保存済みの比較ログを確認したもので、比較を再実行していない。

引数反変・戻り値共変の対応はこの収録範囲で一区切りとする。未登録型を後から解決する場合の宣言時点・診断、autoloadが関与する型階層、未収録の複合型境界は残る。項目8全体と項目6〜10全体は未完了。次の実装単位は`void`・`never`を含む戻り値型の補完とし、型付きプロパティ・定数、readonly、enum等とは区別する。

## 787件追加直後：宣言時のint／float型互換性修正（履歴）

`typeTermAccepts()`から`float ← int`の特例を除去した。PHP 8.4の[`zend_perform_covariant_type_check()`](https://raw.githubusercontent.com/php/php-src/PHP-8.4/Zend/zend_inheritance.c)は宣言の組み込み型マスクに追加された型を不互換とし、引数の反変は引数の順序を入れ替えて同じ検査を使う。実行時の`coerceType()`にある整数からfloatへの変換は変更していない。回帰fixtureは`runtime/variance_float_return.php`、`variance_float_param.php`、`variance_float_union_return.php`、`variance_float_union_param.php`と、`success/variance_float_runtime_coercion.php`の5件。正常1・実行時4を加えて**787件**（正常423・実行時268・構文96）。既存782件のPHPソース・期待値は変更していない。

[`make test`](validation/type-variance-float-test.log)、[`make vet`](validation/type-variance-float-vet.log)、[`make build`](validation/type-variance-float-build.log)、[`git diff --check`](validation/type-variance-float-diff-check.log)は各終了0。[Podman確認](validation/type-variance-float-podman.log)は終了125。[対象比較試行](validation/type-variance-float-targeted-attempt.log)と[全件比較試行](validation/type-variance-float-compare-all-attempt.log)は各終了2で、PHP3版の版確認が`podman exec`終了125（ソケット接続拒否）。各版とも**今回0一致・0不一致・1環境エラー、全787件未実施**。新規18件はまだ公式PHPと実比較できず、期待値・診断・拒否時点は暫定。直前に公式比較済みなのは下記769件である。

Macでは`cd /Users/kouji/dev/php`、`podman ps`を確認し、停止している既存コンテナだけ`podman start php-version-83 php-version-84 php-version-85`で起動する。`set -o pipefail`、`make build`の後、`./containers/php-version/compare-all.sh tests/fixtures/runtime/variance_float_return.php ./bin/gophp 2>&1 | tee docs/validation/type-variance-float-targeted-mac.log`を実行し、直後に`comparison_exit=$?; echo "targeted comparison exit=$comparison_exit" | tee -a docs/validation/type-variance-float-targeted-mac.log`を記録する。全件は`make compare-all 2>&1 | tee docs/validation/type-variance-float-compare-all-mac.log`、直後に`comparison_exit=$?; echo "make compare-all exit=$comparison_exit" | tee -a docs/validation/type-variance-float-compare-all-mac.log`を記録する。残りの新規fixtureも第1引数を各ファイルへ替えて対象比較する。

## 782件追加直後の状態（履歴）：複合型の継承署名検査

PHP 8.4.26を主基準とし、8.3.35・8.5.11も対象に、`types.go`の宣言型を「union選択肢の中のintersection構成型」として包含比較する。引数は子の許容範囲が親を含むか、戻り値は親の許容範囲が子を含むかを調べる。クラス継承・interface実装／継承、`object`、`mixed`、`bool`と`true`／`false`、`null`、型なし宣言を扱う。弱い実行時スカラー変換はこの宣言時判定に使用しない。DNF宣言の冗長性検査は別処理のままで、autoloadを起こさない。

新規13件は正常10件・実行時失敗3件。総数は**782件**（正常422・実行時264・構文96）。元の769件のソース・期待値は変更していない。[`make test`](validation/type-variance-test.log)、[`make vet`](validation/type-variance-vet.log)、[`make build`](validation/type-variance-build.log)、[`git diff --check`](validation/type-variance-diff-check.log)は各終了0。[対象](validation/type-variance-targeted-attempt.log)と[全件](validation/type-variance-compare-all-attempt.log)の公式比較試行は各終了2で、版確認の`podman exec`が各版125（ソケット接続拒否）。各版**今回0一致・0不一致・1環境エラー、今回の全782件未実施**。このうち769件には直前の成功履歴があり、**新規13件は一度も公式比較されていない**。成功系のGoローカル実行は終了0だが、公式PHPとの一致は未確認。新規fixtureの期待値、失敗時の診断と発生時点は暫定である。

正常系：`variance_class`、`variance_union`、`variance_intersection`、`variance_dnf`、`variance_mixed_object`、`variance_abstract_trait`、`variance_diamond`、`variance_variadic`、`variance_cross_file`、`variance_namespace`。失敗系：`variance_param_narrow`、`variance_return_wide`、`variance_dnf_narrow`。いずれも`tests/fixtures/`の該当分類にある。`support/variance_contract.inc`は独立fixtureに数えない。

固定vendorの`api-platform/state` v5.0.3の`ObjectMapper/ObjectMapper.php:21-42`ではinterface実装で`object|string|null`を引数に、`object`を戻り値に宣言する。`Provider/ParameterProvider.php:161`のintersection引数と`symfony/uid` v7.4.20の`Factory/TimeBasedUuidFactory.php:28`のintersection戻り値も確認した。これらのコードのGo実行到達は未確認であり、今回のfixtureは宣言包含判定の小例を対象にする。未解決の型名は同名なら同一型と扱い、階層が未登録なら包含を推測せず拒否する。公式PHPが後から型を読み込む境界と、その警告・診断時点は未検証。`self`／`parent`／`static`型、`void`／`never`、`iterable`／`callable`、参照、型付きプロパティ・定数は次の単位に残る。項目8全体と項目6〜10全体は未完了。

Macで比較する場合：`cd /Users/kouji/dev/php`で`podman ps`を確認し、既存比較コンテナが停止していれば`podman start php-version-83 php-version-84 php-version-85`で起動する。`set -o pipefail`後、`make build`、`./containers/php-version/compare-all.sh tests/fixtures/success/variance_class.php ./bin/gophp 2>&1 | tee docs/validation/type-variance-targeted-mac.log`、直後に`comparison_exit=$?; echo "targeted comparison exit=$comparison_exit" | tee -a docs/validation/type-variance-targeted-mac.log`を実行する。全件は`make compare-all 2>&1 | tee docs/validation/type-variance-compare-all-mac.log`、直後に`comparison_exit=$?; echo "make compare-all exit=$comparison_exit" | tee -a docs/validation/type-variance-compare-all-mac.log`を実行する。対象を追加する際は`tests/fixtures/{success,runtime}/variance_*.php`を一つずつ比較スクリプトの第1引数へ渡す。

## 769件の修正後成功（履歴）

ユーザーがMacで実行した`/tmp/php-compare-769-after.log`はPHP 8.3.35・8.4.26・8.5.11が各**769 matched・0 different・0 errors**、末尾の終了コード0。[保存ログ](validation/dnf-769-after-fix-mac.log)と元ログのバイト一致を確認した。追加19件の公式比較待ちは、収録769件の範囲で解消した。今回は保存済みログを確認したもので、比較を再実行していない。次節の不一致・接続拒否は過去の履歴である。

## DNF診断2件の修正（当時の公式再比較待ち・履歴）

[修正前のMac比較](validation/dnf-769-before-fix-mac.log)は3版とも767一致・2不一致・0エラー、全件コマンド終了2。重複intersectionの診断を `Type B&A is redundant with type A&B` に、引数型の非括弧intersection後の `|` 診断を `expecting variable` 付きに修正した。該当する既存2件の期待値だけを公式ログに基づき訂正し、fixture総数は769件のまま。戻り値型の診断経路は分けて維持する。修正後の公式PHP再比較は未実施。


## 769件追加直後の判定：intersection／DNFの追加19件は当時公式比較待ち（履歴）

後続レビューで、`(A&B)|(A&B&C)`を実行前に拒否しない差を確認した。`typeAlternative.moreRestrictiveThan()`は名前解決済みの構成型だけで厳密な包含関係を比較し、選択肢の前後を問わず、より制約の強いintersectionを冗長として拒否する。クラス読み込みや継承関係の推測は加えない。[公式DNF RFC](https://wiki.php.net/rfc/dnf_types)と[PHP 8.4の宣言検査](https://github.com/php/php-src/blob/PHP-8.4/Zend/zend_compile.c)を根拠とする。追加4件は`parser/dnf_intersection_subset.php`、`parser/dnf_intersection_subset_reverse.php`、`parser/dnf_intersection_subset_reordered.php`、`success/dnf_intersection_non_subset.php`。前3件は宣言前の`echo`も実行されないこと、最後は包含関係のない2選択肢を受理することを確認する。

当時のfixtureは**769件**（正常412・実行時261・構文96）。既存765件と期待値は変更していない。[`make test`](validation/dnf-subset-test.log)、[`make vet`](validation/dnf-subset-vet.log)、[`make build`](validation/dnf-subset-build.log)は終了0。[対象比較試行](validation/dnf-subset-targeted-attempt.log)と[全件比較試行](validation/dnf-subset-compare-all-attempt.log)はPodmanソケット接続拒否で終了2。3版とも版確認の終了125で停止し、各版**0一致・0不一致・1環境エラー・769件未実施**。追加4件の期待診断は公式ソースに基づく暫定値であり、公式PHPによる実比較済みとは記載しない。Macでは`set -o pipefail`の後、`./containers/php-version/compare-all.sh tests/fixtures/parser/dnf_intersection_subset.php ./bin/gophp 2>&1 | tee docs/validation/dnf-subset-targeted-mac.log`を実行し、直後に`result=$?; echo "targeted comparison exit=$result" | tee -a docs/validation/dnf-subset-targeted-mac.log`を記録する。全件は`make compare-all 2>&1 | tee docs/validation/dnf-subset-compare-all-mac.log`、直後に`result=$?; echo "make compare-all exit=$result" | tee -a docs/validation/dnf-subset-compare-all-mac.log`を記録する。

## 765件時点のintersection／DNF追加（履歴・公式比較待ち）

PHP 8.4.26を主基準に、8.3.35・8.5.11も対象とする。`typeDeclaration`で純粋な`A&B`と括弧付きintersectionを含む`(A&B)|C`を読み、`types.go`の各union選択肢をintersection構成型の列として扱う。実行時はintersectionの全構成型、DNFのいずれかの選択肢への一致を要求し、既存unionの弱いスカラー変換とファイル単位のstrict_typesを維持する。引数・戻り値・可変長、クロージャ・arrow・staticメソッド、名前付き／展開済み引数、同等な型集合の継承署名を対象とした。完全な共変・反変と型付きプロパティ・定数は残件。

fixtureは**765件**（正常411・実行時261・構文93）。既存750件とその期待値は変更せず、新規15件の期待値は**Goローカル照合済み・公式PHP比較待ちの暫定値**である。[`make test`](validation/dnf-local-test.log)、[`make vet`](validation/dnf-local-vet.log)、[`make build`](validation/dnf-local-build.log)は終了0。[対象比較試行](validation/dnf-targeted-attempt.log)と[全件比較試行](validation/dnf-compare-all-attempt.log)はPodmanソケット接続拒否で終了2。3版ともPHP版確認が終了125、各版**0一致・0不一致・1環境エラー・765件未実施**。fixtureの不一致やPHPプログラムの終了125ではない。下の750件の成功は変更前の確定履歴であり、765件の公式比較結果ではない。

追加した正常系は`intersection_basic`、`dnf_branches`、`dnf_call_forms`、`dnf_namespace_alias`、`dnf_cross_file_strict`、`dnf_inheritance`。実行時失敗は`intersection_partial_type`、`dnf_return_type`、`dnf_partial_type`、`intersection_implicit_nullable`。構文・宣言失敗は`intersection_builtin`、`dnf_ungrouped`、`dnf_duplicate`、`dnf_default_invalid`、`intersection_nullable_shorthand`で、いずれも`tests/fixtures/`配下の`.php`である。補助ファイル`support/dnf_strict.inc`は独立fixtureに数えない。型の既定値・冗長型の宣言診断とPHP 8.4以降の暗黙nullable非推奨は、公式PHP実比較まで暫定とする。

Macでの再比較：`cd /Users/kouji/dev/php`、`podman info`で接続を確認し、停止している既存コンテナだけ`podman start php-version-83 php-version-84 php-version-85`で起動する。`set -o pipefail`の後、`./containers/php-version/compare-all.sh tests/fixtures/success/intersection_basic.php ./bin/gophp 2>&1 | tee docs/validation/dnf-targeted-mac.log`を実行し、**直後に**`result=$?; echo "targeted comparison exit=$result" | tee -a docs/validation/dnf-targeted-mac.log`を記録する。全件は`make compare-all 2>&1 | tee docs/validation/dnf-compare-all-mac.log`、直後に`result=$?; echo "make compare-all exit=$result" | tee -a docs/validation/dnf-compare-all-mac.log`を記録する。新規全件を対象にする場合は`tests/fixtures/{success,runtime,parser}/{intersection_*,dnf_*}.php`を個別に指定して同様に実行する。差異があれば公式PHPの実出力を先に確認する。

## 750件時点の判定：収録fixtureは公式PHP3版で全一致（履歴）

ユーザーがMacで実行した[修正後の全件比較ログ](validation/union-diagnostic-750-after-fix-mac.log)を`/tmp/php-compare-750-after.log`から内容を変えずに保存し、`cmp`でバイト一致を確認した。PHP 8.3.35・8.4.26・8.5.11は各**750 matched・0 different・0 errors**、ログ末尾の終了コードは0。[`parser/union_nullable_syntax.php`の対象比較](validation/union-nullable-syntax-targeted-mac.log)も`/tmp/union-targeted.log`とバイト一致し、各版**1 matched・0 different・0 errors**。対象比較コマンドの終了0はユーザー報告による確認であり、対象ログ自体には終了コードの行がない。727件から追加した23件の公式比較待ちは、収録fixtureの範囲で解消した。今回は保存済みログの確認・保存で、比較を再実行していない。union型の完全な共変・反変、intersection、DNF、型付きプロパティ・定数などは引き続き項目8の残件である。項目6〜10全体とGo製PHPによるLaravel起動も未完了。

## 750件の修正前不一致と比較試行（履歴）

ユーザーのMacで実行した[`make compare-all`修正前ログ](validation/union-diagnostic-750-before-fix-mac.log)を内容を変えずに保存し、元の`/tmp/php-compare-750.log`と`cmp`でバイト一致を確認した。PHP 8.3.35・8.4.26・8.5.11は各**749 matched・1 different・0 errors**、ログ記載の終了コードは2。唯一の差は`parser/union_nullable_syntax.php`の`?int|string`直後の`|`で、公式PHPは`Parse error: syntax error, unexpected token "|", expecting variable`、Go側は独自文言だった。`typeDeclaration()`でnullable短縮形の直後の`|`を検出して公式文言を返し、fixture期待値を公式PHPの実出力に基づき訂正した。既存の未選択分岐内のfixtureを使い、実行前の拒否とstdoutが空になることを単体テストでも確認した。

修正後は[`make test`](validation/union-nullable-syntax-test.log)、[`make vet`](validation/union-nullable-syntax-vet.log)、[`make build`](validation/union-nullable-syntax-build.log)、[`git diff --check`](validation/union-nullable-syntax-diff-check.log)が終了0。[対象比較試行](validation/union-nullable-syntax-targeted-attempt.log)と[全件比較試行](validation/union-nullable-syntax-compare-all-attempt.log)はいずれもPodman接続拒否で終了2、3版ともPHP版確認の終了125以降は実行されていない。修正後の公式PHP一致・不一致は未測定であり、750件全件一致も未確認。Macでは既存コンテナを確認・起動し、下記の手順で対象と全件を再比較する。

### union型診断表示の追補（修正前の実比較に含まれる）

union型の診断表示に、PHPの`zend_type_to_string_resolved()`と同じ共通順序を適用した。クラス名は宣言順、組み込み型は固定順とし、`int|string`を`string|int`、`int|string|null`を`string|int|null`と表示する。引数・戻り値・継承メソッド署名・既定値検査は同じ`declaredTypeText()`を使う。[PHP 8.4の実装](https://sources.debian.org/src/php8.4/8.4.26-1/Zend/zend_compile.c/)を根拠とする。この追補を記した時点では新規3件と既存のunion診断fixtureの**公式PHP実比較は未実施**だったが、後続の冒頭750件成功ログで収録範囲の一致を確認した。既存727件の比較済みfixtureは変更していない。既存のunion診断3件（引数・戻り値・不正な既定値）の暫定期待値は、公式ソースの表示順に基づいて訂正した。今回のfixture追加後は**750件**（正常405・実行時257・構文88）である。[`make test`](validation/union-diagnostic-test.log)、[`make vet`](validation/union-diagnostic-vet.log)、[`make build`](validation/union-diagnostic-build.log)、[`git diff --check`](validation/union-diagnostic-diff-check.log)は終了0。[`make compare-all`試行](validation/union-diagnostic-compare-all-attempt.log)はPodmanソケット接続拒否で終了2、各版の版確認が終了125となり、各版0一致・0不一致・1環境エラー・750件未実施。公式PHPの終了125やfixtureの不一致とは区別する。

Macでは下記の既存コンテナ確認の後、`./containers/php-version/compare-all.sh tests/fixtures/runtime/union_argument_type.php ./bin/gophp`などの対象比較と`make compare-all`を実行する。対象ログを残す場合は`set -o pipefail`の下で`... 2>&1 | tee docs/validation/union-diagnostic-targeted-mac.log`とし、直後の`comparison_exit=$?`を保存する。全件ログも同様に`make compare-all 2>&1 | tee docs/validation/union-diagnostic-compare-all-mac.log`の直後に終了コードを保存する。

### 747件追加直後（履歴）

項目8の関数・メソッド・staticメソッド・クロージャ・arrow関数の引数／戻り値に、`int|string`などのunion型を追加した。`function.go`で構成型を名前解決して重複・冗長な宣言を検査し、`types.go`では順序付きの構成型として扱う。実行時は既存型に一致する値を保持し、弱いスカラー変換だけPHPの優先順位で選ぶ。引数は呼び出し元、戻り値は宣言元の`strict_types`を維持する。継承・interfaceの同じ構成型集合は記述順が異なっても互換と判定する。完全なunionの共変・反変、intersection、DNF、型付きプロパティ／定数は未対応である。固定vendorの`api-platform/laravel`の`Routing/IriConverter.php:100`に`object|string`、Laravelの`Filesystem.php:585`に`array|string|int`を目視確認した。これらのvendor実行経路がGoで通ることは未確認である。

新規**20件**（正常11・実行時3・構文6）で総数**747件**（正常405・実行時254・構文88）。既存727件のPHPソースと期待値は維持した。新規期待値はGoローカル実行とPHP公式の[型宣言](https://www.php.net/manual/en/language.types.declarations.php)・[弱いunion型変換](https://wiki.php.net/rfc/union_types_v2)・[暗黙nullableの非推奨](https://wiki.php.net/rfc/deprecate-implicitly-nullable-types)に基づく**暫定値**であり、公式PHPとのバイト単位一致は未確認。[対象比較試行](validation/union-targeted-attempt.log)は比較器終了2、PHP3版の版確認は各`podman exec`終了125（ソケット接続拒否）だった。各版0一致・0不一致・1環境エラー・新規20件未実施。全747件も[全件比較試行](validation/union-compare-all-attempt.log)で終了2、比較自体は未実施である。PHPプログラムの終了125やfixture差異とは区別する。

Macで既存コンテナの状態を確認し、停止していれば起動した上で次を実行する。終了コードは各コマンド直後に保存する。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
for fixture in tests/fixtures/{success,runtime,parser}/union_*.php; do
  ./containers/php-version/compare-all.sh "$fixture" ./bin/gophp 2>&1 | tee -a docs/validation/union-targeted-mac.log
  comparison_exit=$?
  echo "$fixture comparison exit=$comparison_exit" | tee -a docs/validation/union-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/union-compare-all-mac.log
comparison_exit=$?
echo "make compare-all exit=$comparison_exit" | tee -a docs/validation/union-compare-all-mac.log
```

後続のローカル追補：`?T`と`T|null`を同じ構成型集合として扱う継承署名テストを追加した。PHP 8.4の[`zend_compile_typename_ex`](https://github.com/php/php-src/blob/PHP-8.4/Zend/zend_compile.c)にある`true|false`と`object|クラス型`の宣言診断へ文言を合わせた。[`make test`](validation/union-followup-test.log)、[`make vet`](validation/union-followup-vet.log)、[`make build`](validation/union-followup-build.log)、[`git diff --check`](validation/union-followup-diff-check.log)は各終了0。[`make compare-all`再試行](validation/union-followup-compare-all-attempt.log)はPodmanソケット接続拒否で終了2、3版とも版確認の終了125以降は実行されず、0一致・0不一致・1環境エラー、747件未実施である。**公式PHPでの実比較は未実施**。残る境界はunionを含む完全な分散、定数式による既定値、すべての弱い型変換・非推奨診断、版別の宣言診断。これらを公式PHP比較で確認して次の型実装単位へ進む。項目8全体と項目6〜10全体は未完了である。

## 727件の全一致（直前の確定履歴）

ユーザーがMacで実行した`/tmp/php-compare-latest2.log`を[727件の成功ログ](validation/static-throw-727-compare-all-mac.log)へ内容を変えずに保存し、`cmp`でバイト一致を確認した。PHP 8.3.35・8.4.26・8.5.11は各**727 matched・0 different・0 errors**、ログ末尾の終了コードは0。直接の`Class::{throw ...}()`、括弧付き、未選択分岐、未呼び出し関数・メソッド、動的クラス名、`parent::`、読み込み先、合法な対比を含む追加9件も3版すべてでMATCHした。今回新たに比較を実行した結果ではなく、保存済みのユーザー実行ログの確認である。追加9件の公式比較待ちは解消したが、未収録の式形式、項目7全体、項目6〜10全体、およびGo製PHPでのLaravel起動は未完了。

## 727件追加直後の判定（履歴）：Go側修正・公式比較待ち

2026-10-15に`internal/php/validate.go`で`staticCallExpr`のメソッド名式が直接`throwExpr`である場合、実行前の検査で`Method name must be a string`を返すようにした。括弧は同じASTに畳まれる。PHP 8.3・8.4・8.5の[`zend_compile_static_call`](https://github.com/php/php-src/blob/PHP-8.4/Zend/zend_compile.c)はメソッド名式をコンパイルし、その結果が定数かつ非文字列なら`E_COMPILE_ERROR`にする。`zend_compile_throw`は式としてコンパイルされると定数の真を返す。したがって直接形はコンパイル時の失敗となる。三項式の内部にある`throw`やメソッド名を返す関数呼び出しまで一律に拒否しない。主ファイルの未選択分岐・未呼び出し関数／メソッドでは実行前、読み込み先では読み込み開始後のコンパイル時に失敗する。後者は`catch`や`finally`で打ち消せない既存の宣言エラー経路を使う。

新規9件（構文7・実行時1・正常1）を加え、総数は**727件**（正常394・実行時251・構文82）。既存718件のソースと期待値は維持した。[`make test`](validation/static-throw-make-test.log)、[`make vet`](validation/static-throw-make-vet.log)、[`make build`](validation/static-throw-make-build.log)は各終了0。新規期待値は保存済みの旧公式比較で直接形の文言と8.5のトレースを確認し、その他の境界は公式コンパイラソースとGoローカル結果から設定した**暫定値**である。新規9件の公式実比較は未実施。既存718件の3版全一致は後続節の過去の確定結果であり、727件の全一致を意味しない。

[対象比較試行](validation/static-throw-targeted-attempt.log)と[全件比較試行](validation/static-throw-compare-all-attempt.log)は、各版の確認でPodman接続が終了125となり、両比較コマンドは終了2だった。各版の公式比較は0件実施、727件未実施である。環境エラーをfixtureの一致・不一致に数えない。Macでは既存コンテナを確認・起動し、次のコマンドで対象9件と全件を比較する。各終了コードはコマンド直後に保存する。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
make build || exit $?
set -o pipefail
for fixture in tests/fixtures/parser/class_static_throw_*.php tests/fixtures/runtime/class_static_throw_include.php tests/fixtures/success/class_static_throw_nested.php; do
  ./containers/php-version/compare-all.sh "$fixture" ./bin/gophp 2>&1 | tee -a docs/validation/static-throw-targeted-mac.log
  result=$?
  echo "$fixture comparison exit=$result" | tee -a docs/validation/static-throw-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/static-throw-compare-all-mac.log
result=$?
echo "make compare-all exit=$result" | tee -a docs/validation/static-throw-compare-all-mac.log
```

この時点では全件一致が未確認で、「Go側実装・ローカル検証済み、公式PHP比較待ち」とした。後続の727件全一致は冒頭の成功ログで確認した。式形式全般の定数畳み込みや未収録の境界、項目7全体と項目6〜10全体は残件である。

## 718件時点の判定（履歴）：公式PHP3版で全一致

ユーザーがMacで実行した`/tmp/php-dynamic-static-718-after-fix.log`を[修正後の全件比較ログ](validation/item7-dynamic-static-718-after-fix-mac.log)へ内容を変えずに保存し、`cmp`でバイト一致を確認した。PHP 8.3.35・8.4.26・8.5.11は各**718 matched・0 different・0 errors**、`make compare-all`終了0。`success/class_dynamic_static_order.php`も3版すべてでMATCHし、新規11件の公式比較待ちは解消した。今回は保存済みログの確認・保存であり、こちらで比較を再実行していない。直接`Class::{throw ...}()`をGoが受理する差は、収録fixtureの全一致と別に残る既知の互換性不具合である。項目7と項目6〜10全体は未完了。

### 修正前の718件比較（履歴）：717一致・1検証エラー

ユーザーがMacで実行した`/tmp/php-dynamic-static-718.log`を[修正前の全件比較ログ](validation/item7-dynamic-static-718-before-fix-mac.log)へそのまま保存し、`cmp`でバイト一致を確認した。PHP 8.3.35・8.4.26・8.5.11は各**717 matched・0 different・1 verification error**、比較終了2だった。`success/class_dynamic_static_order.php`の11行目で、メソッド名式へ直接`throw new Exception('stop')`を書くと、公式PHPは実行前に`Fatal error: Method name must be a string`として終了255にする。一方Goはcatchして正常終了したため、成功fixtureの条件を満たさなかった。他の717件は3版ともMATCHした。今回の公式比較はユーザーのMacでの実行結果であり、こちらで再実行できたものではない。

評価順序を確認する成功fixtureとして、11行目のメソッド名式を`orderThrow()`へ変更し、その関数内で例外を投げる形にした。Goの3設定で`class|member|arg|7|class|property|rhs|3|class|caught\n`、stderr空、終了0を確認した。保存済み期待値の変更は不要だった。修正後の[`make test`](validation/item7-dynamic-static-order-fix-test.log)、[`make vet`](validation/item7-dynamic-static-order-fix-vet.log)、[`make build`](validation/item7-dynamic-static-order-fix-build.log)、[`git diff --check`](validation/item7-dynamic-static-order-fix-diff-check.log)は各終了0。[対象比較試行](validation/item7-dynamic-static-order-fix-targeted-attempt.log)と[全件比較試行](validation/item7-dynamic-static-order-fix-compare-all-attempt.log)は当時Podman接続拒否で終了2となり、この時点では修正後の公式PHP再比較が未実施だった。後続の718件全一致は冒頭の成功ログで確認した。`Class::{throw ...}()`をGoが受理する差は未修正の既知不具合であり、別の失敗系fixtureと宣言段階の検査が必要。

Macではコンテナ稼働を確認してから、`set -o pipefail`のもとで次を実行し、比較コマンド直後に終了コードを保存する。既存ログは上書きしない。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
./containers/php-version/compare-all.sh tests/fixtures/success/class_dynamic_static_order.php ./bin/gophp 2>&1 | tee docs/validation/item7-dynamic-static-order-fix-targeted-mac.log
target_exit=$?
echo "targeted exit=$target_exit" | tee -a docs/validation/item7-dynamic-static-order-fix-targeted-mac.log
make compare-all 2>&1 | tee docs/validation/item7-dynamic-static-order-fix-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/item7-dynamic-static-order-fix-compare-all-mac.log
```

### 718件追加時点（履歴）：当時は新規11件を公式比較待ちとした

ユーザーがMacで実行した`/tmp/php-compare-latest.log`を[707件の成功ログ](validation/item7-declarations-707-compare-all-mac.log)へ内容を変えずに保存し、`cmp`でバイト一致を確認した。PHP 8.3.35・8.4.26・8.5.11で各**707 matched・0 different・0 errors**、`make compare-all`終了0。追加5件の比較待ちは解消した。これは保存済みログの確認であり、今回新たに公式PHPを実行した結果ではない。

今回、`$class::method()`、`ClassName::$method()`、`$class::$method()`、`ClassName::{expr}()`、`ClassName::${expr}`の動的指定を、既存のstaticメンバー探索・可視性・呼び出し・引数束縛へ接続した。裸の`ClassName::$name`は固定名のstaticプロパティ、呼び出し付きの`ClassName::$name()`は変数`$name`の値をメソッド名とする。代入先のクラス式・名前式は右辺より前に一度評価して保持する。既存の`self`・`parent`・`static`の転送と明示的なクラス指定の遅延静的束縛を維持した。[PHPのstaticメンバー](https://www.php.net/manual/en/language.oop5.static.php)と[変数関数の規則](https://www.php.net/manual/en/functions.variable-functions.php)を構文の根拠とする。

新規fixtureは`success/class_dynamic_static_{binding,call,errors,keywords,order,property,trait,unselected,visibility}.php`と`runtime/class_dynamic_static_{missing_method,private}.php`の**11件**。現在718件（正常393・実行時250・構文75）。既存707件のPHPソース・期待値は変更していない。新規期待値はGoローカル出力を保存した**暫定値**で、公式PHP3版との実比較前である。追加11件時点の[`make test`](validation/item7-dynamic-static-test.log)、[`make vet`](validation/item7-dynamic-static-vet.log)、[`make build`](validation/item7-dynamic-static-build.log)は終了0だった。今回、文字列変数内の`self`・`parent`・`static`を特殊なスコープ名として解決しないよう修正し、追加1件とGo単体テストを加えた。[最終`make test`](validation/item7-dynamic-static-followup-final-test.log)（GoとPython12件）、[`make vet`](validation/item7-dynamic-static-followup-final-vet.log)、[`make build`](validation/item7-dynamic-static-followup-final-build.log)は終了0。[対象比較試行](validation/item7-dynamic-static-followup-targeted-attempt.log)と[全件比較試行](validation/item7-dynamic-static-followup-compare-all-attempt.log)はPodmanソケット接続拒否により各終了2。PHP版確認の`podman exec`は各終了125だった。各版0一致・0不一致・1環境エラー、718件未実施であり、fixtureの不一致には数えない。

残件は動的クラス定数、`__callStatic`、未収録の式形式・評価順序・診断、型付きstaticプロパティ、traitとprivate staticの複雑な継承境界である。項目7全体、項目6〜10全体、Go製PHPでのLaravel起動は未完了。Macでは次を実行する（既存コンテナは停止時だけ起動し、ログは新しい名前で保存する）。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/item7-dynamic-static-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/item7-dynamic-static-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
for file in tests/fixtures/success/class_dynamic_static_*.php tests/fixtures/runtime/class_dynamic_static_*.php; do
  ./containers/php-version/compare-all.sh "$file" ./bin/gophp 2>&1 | tee -a docs/validation/item7-dynamic-static-targeted-mac.log
  target_exit=$?
  echo "$file exit=$target_exit" | tee -a docs/validation/item7-dynamic-static-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/item7-dynamic-static-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/item7-dynamic-static-compare-all-mac.log
```

### 707件の追加実装時点（履歴）：当時は追加5件を公式PHP比較待ちとした

項目7の今回の追加では、抽象親クラスが保留したinterface契約を具象子クラスでも検査するようにした。interface本体のメソッド検査はソース順に行う。traitの`as`で異なるメソッドに同じ別名を与えた場合は、後からの別名を黙って捨てず競合として拒否する。クラス本体の同名メソッドが別名より優先する動作は維持した。既存の`insteadof`・選ばれなかったメソッドの別名・分割use・入れ子traitを再利用し、[PHPのinterface規則](https://www.php.net/manual/en/language.oop5.interfaces.php)と[traitの競合規則](https://www.php.net/manual/en/language.oop5.traits.php)を根拠にした。

新規fixtureは`success/interface_abstract_inherited_contract.php`、`runtime/interface_abstract_child_missing.php`、`success/trait_nested_resolution.php`、`runtime/trait_alias_collision.php`、`success/trait_alias_class_priority.php`の5件。現在707件（正常384・実行時248・構文75）で、既存702件のPHPソース・期待値は維持した。新規5件の診断を含む期待値は公式PHP3版との実比較前の**暫定値**。[`make test`](validation/item7-declarations-test.log)（GoとPython12件）、[`make vet`](validation/item7-declarations-vet.log)、[`make build`](validation/item7-declarations-build.log)、[`git diff --check`](validation/item7-declarations-diff-check.log)は各終了0。Podmanの[状態確認](validation/item7-declarations-podman.log)はソケット接続拒否で終了125。[対象比較試行](validation/item7-declarations-targeted-attempt.log)と[全件比較試行](validation/item7-declarations-compare-all-attempt.log)は各版のPHP確認前に環境エラーとなり、全件比較は終了2。今回の公式PHP実比較は各版0一致・0不一致・1環境エラー・707件未実施である。比較前の[702件全一致ログ](validation/include-path-702-after-fix-mac.log)は新コードの結果として扱わない。

固定vendorの`Illuminate/Database/Query/Builder.php:39-41`は複数traitと`__call as macroCall`を使い、`Illuminate/Database/DatabaseManager.php:23-25`も別名指定を使う。実コードを確認したが、Goでそれらのファイルを起動・実行できるとは未確認。複合型署名、循環・重複した依存の全診断、trait適応句の未収録境界、動的staticメンバーは残件であり、項目7全体と項目6〜10全体は未完了。

Macでの対象・全件比較は次のとおり。既存コンテナが停止している場合だけ起動し、各コマンド直後に終了コードを保存する。先行ログは上書きしない。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/item7-declarations-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/item7-declarations-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
for file in tests/fixtures/success/interface_abstract_inherited_contract.php tests/fixtures/runtime/interface_abstract_child_missing.php tests/fixtures/success/trait_nested_resolution.php tests/fixtures/runtime/trait_alias_collision.php tests/fixtures/success/trait_alias_class_priority.php; do
  ./containers/php-version/compare-all.sh "$file" ./bin/gophp 2>&1 | tee -a docs/validation/item7-declarations-targeted-mac.log
  target_exit=$?
  echo "$file exit=$target_exit" | tee -a docs/validation/item7-declarations-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/item7-declarations-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/item7-declarations-compare-all-mac.log
```

### 702件成功時点の判定（履歴）

ユーザーがMacで実行した[修正後の全件比較ログ](validation/include-path-702-after-fix-mac.log)を保存し、元の`/tmp/php-compare-702-after.log`とのバイト一致を確認した。PHP 8.3.35・8.4.26・8.5.11は各**702 matched・0 different・0 errors**、`make compare-all`終了0。`include_path_empty.php`も3版でMATCHした。追加18件の公式比較待ちは解消した。今回は保存済みログの確認・保存であり、比較を再実行していない。収録外の境界条件、項目6全体と項目6〜10全体、Go製PHPでのLaravel起動は未完了。

### 702件の修正前比較と接続拒否（履歴）

ユーザーがMacで実行した[修正前の702件全件比較ログ](validation/include-path-702-before-fix-mac.log)では、PHP 8.3.35・8.4.26・8.5.11の各版で701一致・1不一致・0エラー、`make compare-all`終了2だった。不一致は`success/include_path_empty.php`のstdoutだけで、公式PHPは`bad|missing|first|restored\n`、Goは`empty|missing|first|restored\n`だった。`set_include_path('')`が失敗して`false`を返し、既存のパスを維持するよう修正した。fixtureのソースは維持し、公式PHPの実出力を根拠に保存済み期待値だけを訂正した。戻り値が`false`で設定値が変わらないことを3設定で確認するGo単体テストも追加した。

修正後の[`make test`](validation/include-path-empty-fix-test.log)、[`make vet`](validation/include-path-empty-fix-vet.log)、[`make build`](validation/include-path-empty-fix-build.log)、[`git diff --check`](validation/include-path-empty-fix-diff-check.log)は各終了0。Goの3設定で当該fixtureのstdoutは公式ログと同じになった。[修正後の対象比較試行](validation/include-path-empty-fix-targeted-attempt.log)と[全件比較試行](validation/include-path-empty-fix-compare-all-attempt.log)はPodman接続拒否により各版の版確認で終了125、比較コマンドは各終了2。各版0一致・0不一致・1環境エラーで、修正後の公式比較は未実施である。以下のMac手順で対象比較と全件比較を再実行する。過去の701一致・1不一致は修正前の実測値として維持する。

### 厳密モードのnull引数修正と700件時点の記録（履歴）

2026-10-14の追補として、`set_include_path(null)`と`ini_restore(null)`を`declare(strict_types=1)`で呼ぶとき、非推奨診断と空文字列への変換に進まずTypeErrorを返すよう修正した。PHPの[組み込み関数の型規則](https://www.php.net/manual/en/functions.internal.php)では、nullの暗黙変換は弱い型モードだけに適用される。設定値を変えないことを確認する2件のfixtureとGo単体テストを追加した。**当時702件（正常381・実行時246・構文75）**で、684件以降の追加18件は公式PHP3版で未比較、期待値は暫定だった。今回の[`make test`](validation/include-path-strict-null-test.log)、[`make vet`](validation/include-path-strict-null-vet.log)、[`make build`](validation/include-path-strict-null-build.log)、[`git diff --check`](validation/include-path-strict-null-diff-check.log)は各終了0。[対象比較試行](validation/include-path-strict-null-targeted-attempt.log)と[全件比較試行](validation/include-path-strict-null-compare-all-attempt.log)は各終了2で、3版とも版確認時のPodman接続失敗（`podman exec`終了125）が1環境エラー、0一致・0不一致、702件未実施である。PHPプログラムの終了コード125と混同しない。

Macでは下の手順で既存コンテナの稼働を確認し、`tests/fixtures/success/include_path_strict_null_{set,restore}.php`を対象比較した後、`make compare-all`を実行する。既存ログを上書きせず、`set -o pipefail`のもとで各比較コマンド直後の終了コードを保存する。

#### 700件時点の記録（履歴）

2026-10-14に項目6のローカルパス探索と`include_path`を追加した。現在700件（正常379・実行時246・構文75）で、[`make test`](validation/include-path-make-test.log)（GoとPython12件）、[`make vet`](validation/include-path-make-vet.log)、[`make build`](validation/include-path-make-build.log)、[`git diff --check`](validation/include-path-diff-check.log)は各終了0。公式PHPとの直近の確定結果は追加前の684件で、下記の[Mac成功ログ](validation/error-handler-684-fix-compare-all-mac.log)による各版684一致・0不一致・0エラー、終了0である。新規16件の期待値は暫定であり、公式PHP3版との実比較は未実施。今回の[対象比較試行](validation/include-path-targeted-attempt.log)と[全件比較試行](validation/include-path-compare-all-attempt.log)は、各版の版確認でPodman接続拒否（終了125）となり、比較コマンドは両方とも終了2。各版0一致・0不一致・1環境エラー、新規16件を含む700件が未実施である。PHP fixtureの差異ではない。

`include`、`require`、両`_once`を共通の探索処理へ接続した。絶対・`./`・`../`を現在の作業ディレクトリから直接解決し、通常の相対名は`include_path`の各要素を順に探索してから実行中ファイルのディレクトリを試す。相対`include_path`要素は作業ディレクトリ基準。実パスへ正規化して既存の`_once`履歴を照合し、解析失敗時の登録タイミング、読み込み元のスコープとファイル文脈、例外伝播を維持した。`get_include_path()`、`set_include_path()`、`ini_restore('include_path')`を実行単位の状態として追加した。比較コンテナの`php -n`に合わせた初期値は`.:/usr/local/lib/php`。`restore_include_path()`はPHP 8.0で削除済みなので、8.3〜8.5では未定義関数として拒否する。[設計と公式ソース](phase2-design.md)を参照。

新規fixtureは`success/include_path_{empty,explicit,handler_throw,named,nested,once,order,relative,require_catch,restore,retry,unselected}.php`、`runtime/include_path_{get_extra,invalid_type,restore_removed,warning_false}.php`。補助の`.inc`と相対シンボリックリンクは`tests/fixtures/support/include_path/`に置き、単独fixtureに含めない。ローカルテストも比較器も作業ディレクトリをリポジトリのルートへ揃えた。未収録の読み取り権限・ストリームラッパー・他のini設定・全OSパス規則は残件で、項目6全体と項目6〜10全体を完了扱いにしない。

Macでの対象比較と全件比較は、既存コンテナの稼働を確認したうえで次を実行する。各コマンド自身の終了コードを`tee`直後に保存する。先行ログは上書きしない。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/include-path-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/include-path-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
for file in tests/fixtures/success/include_path_*.php tests/fixtures/runtime/include_path_*.php; do
  ./containers/php-version/compare-all.sh "$file" ./bin/gophp 2>&1 | tee -a docs/validation/include-path-targeted-mac.log
  target_exit=$?
  echo "$file exit=$target_exit" | tee -a docs/validation/include-path-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/include-path-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/include-path-compare-all-mac.log
```

## 684件成功時点の判定（履歴）

ユーザーがMacで実行した[修正後の全件比較ログ](validation/error-handler-684-fix-compare-all-mac.log)を確認した。PHP 8.3.35・8.4.26・8.5.11の各版で**684 matched・0 different・0 errors**、`make compare-all`終了0。fixture実数と期待値登録数はともに684件（正常367・実行時242・構文75）。追加したエラーハンドラ関連22件の公式比較待ちは解消した。`error_handler_invalid_callable.php`、`error_handler_invalid_level.php`、`error_handler_mask.php`、`error_handler_user_fatal.php`も3版でMATCHした。今回は保存済みのユーザー実行ログを確認したもので、比較を再実行していない。

この一致は収録fixtureの範囲である。対応callableの未収録形式、ハンドラ中の登録・復元境界、すべての警告発生箇所からの例外中断、`error_reporting()`・`@`・`set_exception_handler()`などは残る。項目10全体と項目6〜10全体、Go製PHPでのLaravel起動は未完了である。

### 2026-10-14：684件Mac比較で見つかった4件（修正前の履歴）

ユーザーがMacで実行した[対象比較](validation/error-handler-null-restore-targeted-mac.log)は新規のnull登録fixtureが3版とも1一致・0不一致・0エラー、終了0。[修正前の全件比較](validation/error-handler-null-restore-compare-all-mac.log)はPHP 8.3.35・8.4.26が各681一致・3不一致・0エラー、8.5.11が680一致・4不一致・0エラー、`make compare-all`終了2だった。保存ログでは`set_error_handler`の不正callableと`trigger_error`の不正レベルの呼び出しフレームは3版とも存在し、`E_USER_ERROR`のFatal errorに呼び出しフレームが付くのは8.5だけだった。

公式ログを根拠に、`set_error_handler`のTypeError文言へ`or null`を追加し、組み込み関数自身が発生させた引数エラーのフレームと8.5の`E_USER_ERROR`のフレームを記録する。単独の`$missing;`は値が使われず、公式PHPでは未定義変数の警告が出ないため、その式の評価を省いた。4件の既存期待値はログのPHP出力に合わせて訂正した。この節の記録時点では修正後の公式PHP実比較は未実施だった。

修正後の[`make test`](validation/error-handler-684-fix-test.log)（GoとPython12件）、[`make vet`](validation/error-handler-684-fix-vet.log)、[`make build`](validation/error-handler-684-fix-build.log)、[`git diff --check`](validation/error-handler-684-fix-diff-check.log)は終了0。[保存済み公式出力との再照合](validation/error-handler-684-saved-output-recheck.log)では、差があった10件のstdout／stderrが修正後のGo実行とバイト一致し、各fixtureの終了コードも一致した。これは保存ログとの照合であり、公式PHPの再実行ではない。[修正後の全件比較試行](validation/error-handler-684-fix-compare-all-attempt.log)はPodmanソケット接続拒否で終了2。3版とも版確認の`podman exec`が終了125となったため、各版0一致・0不一致・1環境エラー、684件未実施である。

Macでの修正後の全件再比較は、既存コンテナが停止している場合だけ起動し、次を実行する。終了コードは`tee`の直後に取得する。修正前のMacログは上書きしない。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make compare-all 2>&1 | tee docs/validation/error-handler-684-fix-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/error-handler-684-fix-compare-all-mac.log
```

### 2026-10-14：ハンドラ内のnull登録後の復元（比較前の履歴）

実行中のハンドラ内で`set_error_handler(null)`を呼んだ場合、登録操作があっても終了時に登録が空なら元のハンドラを復元するよう修正した。[回帰fixture](../tests/fixtures/success/error_handler_register_null_inside.php)とGo単体テストを追加した。期待値は公式PHPで未確認の暫定値である。前回の再登録・長文メッセージのfixtureと合わせて、今回の公式比較対象は3件となる。

修正後の[`make test`](validation/error-handler-null-restore-test.log)（Go・Python12件）、[`make vet`](validation/error-handler-null-restore-vet.log)、[`make build`](validation/error-handler-null-restore-build.log)、[`git diff --check`](validation/error-handler-null-restore-diff-check.log)は終了0。fixtureと期待値は684件（正常367・実行時242・構文75）。[対象比較](validation/error-handler-null-restore-targeted-attempt.log)と[全件比較](validation/error-handler-null-restore-compare-all-attempt.log)は各終了2で、3版ともPodman接続拒否により版確認前に停止した。各版0一致・0不一致・1環境エラー、684件未実施。

### 2026-10-14：ハンドラ再登録と長いメッセージの修正（前回の履歴）

レビューで見つかった2件を修正した。実行中のハンドラを一時的に外すため、その中の`set_error_handler()`は`null`を返し、新しく登録されたハンドラが入れ子の警告を扱う。警告メッセージの1024バイト切り詰めを除去し、1100バイトをそのまま渡す。[再登録](../tests/fixtures/success/error_handler_register_inside.php)と[長いメッセージ](../tests/fixtures/success/error_handler_long_message.php)を追加した。期待値は公式PHPで未確認の暫定値である。

修正後の[`make test`](validation/error-handler-review-test.log)（Go・Python12件）、[`make vet`](validation/error-handler-review-vet.log)、[`make build`](validation/error-handler-review-build.log)、[`git diff --check`](validation/error-handler-review-diff-check.log)はいずれも終了0。fixtureと期待値は683件（正常366・実行時242・構文75）。[対象比較](validation/error-handler-review-targeted-attempt.log)は2ファイルとも終了2、[全件比較](validation/error-handler-review-compare-all-attempt.log)は終了2。各版の`podman exec`は接続拒否で終了125となり、公式PHPの版確認前に停止した。各版0一致・0不一致・1環境エラー、683件未実施。PHPプログラムの終了125やfixture差異ではない。

Macでは既存コンテナの稼働を確認してから、`make build`後に`./containers/php-version/compare-all.sh tests/fixtures/success/error_handler_register_inside.php ./bin/gophp`と同じコマンドの`error_handler_long_message.php`版を実行する。続けて`make compare-all`を実行し、いずれも`set -o pipefail`、`tee docs/validation/...`、直後の`status=$?`でログと終了コードを保存する。下記の全エラーハンドラfixture用手順も利用できる。

`set_error_handler()`・`restore_error_handler()`・`trigger_error()`の部分集合、エラーレベル定数、警告／非推奨診断とハンドラの接続を追加した。関数名・クロージャ・二要素のオブジェクト／静的メソッド配列callableを受け、レベルマスク、登録履歴、`null`登録、`false`による通常診断へのフォールバックを扱う。ハンドラ例外を式・読み込み元まで伝える。実行中のハンドラを一時的に登録から外すため、新しい登録がなければ入れ子の警告は通常診断へ出る。実行中に登録した別ハンドラは入れ子の警告を処理し、`set_error_handler()`はこの場面で`null`を返す。1100バイトの警告メッセージは切り詰めず渡す。未定義変数・配列キー・プロパティ、読み込み失敗、既存の非推奨診断を接続した。PHP 8.4以降の`E_USER_ERROR`非推奨と`E_ALL`値差を対象版で切り替える。構文・宣言エラー、TypeError等のThrowableは警告ハンドラへ流さない。

修正前のローカル検証時点では**684件**（正常367・実行時242・構文75）で、追加22件とGo単体テストを検証した。当時の公式PHP3版との追加分・全件実比較はPodmanソケット接続拒否で未実施であり、期待値は暫定だった。追加前の662件は下記の保存済み成功ログで各版すべて一致した。今回の機能の厳密な完了判定は公式比較後に行う。対応callableの可視性・厳密な型診断、ハンドラ中の未収録の登録・復元境界、すべての既存警告発生箇所における例外中断は未確認・残件。`error_reporting()`、`@`、`set_exception_handler()`は対象外。項目10全体と項目6〜10全体は引き続き部分実装である。

初回のエラーハンドラ追加時点（681件）の記録：[`make test`](validation/error-handler-make-test.log)（Python12件）、[`make vet`](validation/error-handler-make-vet.log)、[`make build`](validation/error-handler-make-build.log)、[`git diff --check`](validation/error-handler-diff-check.log)は終了0。当時のfixture実数と期待値登録数はともに681。公式PHPの[対象比較試行](validation/error-handler-targeted-attempt.log)と[全件比較試行](validation/error-handler-compare-all-attempt.log)は各終了2。3版とも版確認の`podman exec`は終了125で、各版0一致・0不一致・1環境エラー、当時の681件が未実施。これはPHPコードの終了125ではない。追加fixtureのPHP公式実出力は未取得である。

Macでの再比較は、既存コンテナが停止している場合だけ起動してから次を実行する。`tee`の直後に終了コードを保存する。対象比較は今回のファイルだけを列挙し、全件比較で既存分の回帰を確認する。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/error-handler-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/error-handler-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
for file in tests/fixtures/success/error_handler_*.php tests/fixtures/runtime/error_handler_*.php; do
  ./containers/php-version/compare-all.sh "$file" ./bin/gophp 2>&1 | tee -a docs/validation/error-handler-targeted-mac.log
  target_exit=$?
  echo "$file exit=$target_exit" | tee -a docs/validation/error-handler-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/error-handler-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/error-handler-compare-all-mac.log
```

## 662件成功時点の判定（履歴）

ユーザーがMacで実行した`/tmp/php-compare-662-after.log`を[662件の修正後成功ログ](validation/throwable-chain-662-compare-all-mac.log)へ内容を変えずに保存し、`cmp`でバイト一致を確認した。PHP 8.3.35・8.4.26・8.5.11の各版で**662一致・0不一致・0エラー**、`make compare-all`終了0。`exception_trace_call_kinds.php`も各版でMATCHした。今回は保存済みログの確認・保存であり、比較は再実行していない。先行する[641件の成功ログ](validation/include-exception-641-compare-all-mac.log)以降に追加した21件の公式比較待ちは解消した。

現在**662件**（正常350・実行時237・構文75）。例外連鎖、`getPrevious()`・`getTrace()`等のThrowable取得API、名前付き・可変長引数のトレース、派生例外からの`parent::getMessage()`等は、収録fixtureの範囲で3版の公式PHPと一致した。項目10全体と項目6〜10全体は部分実装であり、未収録の境界条件やLaravel起動を保証しない。先行するPodman接続拒否は当時の[対象比較試行](validation/throwable-chain-targeted-attempt.log)と[全件比較試行](validation/throwable-chain-compare-all-attempt.log)、クロージャ名の修正前不一致は[修正前ログ](validation/variadic-trace-662-before-fix-mac.log)に履歴として残す。

対応した組み込みクラスは`Exception`と`Error`、および既存の`TypeError`、`ArgumentCountError`、`ValueError`、`DivisionByZeroError`、`UnhandledMatchError`、`CompileError`、`ParseError`。取得メソッドは`getMessage`、`getCode`、`getFile`、`getLine`、`getPrevious`、`getTrace`、`getTraceAsString`の部分集合。コンストラクタの`message`、`code`、`previous`を位置・名前付き引数で受け、対応済みstrict_typesで検査する。`Throwable`の直接実装を拒否し、ユーザーの`Exception`／`Error`派生クラスでは取得メソッドのfinal上書きを拒否する。`set_error_handler`など、未対応のThrowable API、参照、Reflection、Composer autoload、Laravel起動は残る。

以下は662件成功前に用意した当時の再比較手順である。各終了コードは`tee`の直後に取得する。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/throwable-chain-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/throwable-chain-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
for file in tests/fixtures/success/exception_previous_chain.php tests/fixtures/success/exception_trace_creation.php tests/fixtures/success/exception_finally_chain.php tests/fixtures/success/exception_user_parent.php tests/fixtures/success/exception_previous_type.php tests/fixtures/success/exception_state_separation.php tests/fixtures/success/exception_trace_include.php tests/fixtures/success/exception_trace_values.php tests/fixtures/success/exception_trace_call_kinds.php tests/fixtures/success/exception_finally_existing.php tests/fixtures/success/exception_constructor_strict.php tests/fixtures/success/exception_finally_nested_chain.php tests/fixtures/success/exception_finally_type_error.php tests/fixtures/runtime/exception_uncaught_chain.php tests/fixtures/parser/exception_final_override.php; do
  ./containers/php-version/compare-all.sh "$file" ./bin/gophp 2>&1 | tee -a docs/validation/throwable-chain-targeted-mac.log
  target_exit=$?
  echo "$file exit=$target_exit" | tee -a docs/validation/throwable-chain-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/throwable-chain-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/throwable-chain-compare-all-mac.log
```

## PHP 8.4以降のクロージャのトレース名（修正前の履歴）

ユーザーがMacで実行した`/tmp/php-compare-662.log`を[修正前の全件ログ](validation/variadic-trace-662-before-fix-mac.log)へ内容を変えずに保存し、`cmp`でバイト一致を確認した。PHP 8.3.35は**662一致・0不一致・0エラー**、8.4.26・8.5.11は各**661一致・1不一致・0エラー**、`make compare-all`終了2。唯一の不一致は`success/exception_trace_call_kinds.php`のクロージャの`getTrace()[0]['function']`で、8.4以降の公式PHPは`{closure:ファイルパス:宣言行}`、Goは`{closure}`を返していた。8.3では従来の`{closure}`が一致した。

呼び出しフレームのメソッド表示名をクロージャの宣言元ソース情報と対象版から作り、`getTrace()`と`getTraceAsString()`で共有する。既存fixtureのソースは維持し、公式PHPログを根拠に8.4を基準とするバージョン別stdout期待値を追加した。Goのfixtureテストではパスの`{file}`を実パスへ展開する。比較器のstdoutは加工しない。[保存済み公式stdoutとの対象再照合](validation/closure-trace-saved-output-recheck.log)は8.4・8.5の2件でバイト一致、Goのstderr空・終了0も確認した。これは公式PHPの再実行ではない。修正後の3版全件比較は未実施で、現時点の一致・不一致件数を修正前の結果から更新しない。

[`make test`](validation/closure-trace-test.log)（Python12件を含む）、[`make vet`](validation/closure-trace-vet.log)、[`make build`](validation/closure-trace-build.log)、[`git diff --check`](validation/closure-trace-diff-check.log)は各終了0。修正後の[対象比較試行](validation/closure-trace-targeted-attempt.log)と[全件比較試行](validation/closure-trace-compare-all-attempt.log)はPodmanソケット接続拒否で各終了2、3版とも版確認の`podman exec`は終了125である。各版0一致・0不一致・1環境エラー、662件未実施。環境の終了125をPHPプログラムの終了コードとして扱わない。

Macでの修正後比較は既存コンテナが停止している場合だけ起動し、次を実行する。終了コードは`tee`直後に取得する。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/closure-trace-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/closure-trace-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
./containers/php-version/compare-all.sh tests/fixtures/success/exception_trace_call_kinds.php ./bin/gophp 2>&1 | tee docs/validation/closure-trace-targeted-mac.log
target_exit=$?
echo "targeted exit=$target_exit" | tee -a docs/validation/closure-trace-targeted-mac.log
make compare-all 2>&1 | tee docs/validation/closure-trace-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/closure-trace-compare-all-mac.log
```

## 可変長引数のトレース値の修正（2026-10-14）

`snapshotCallFrame`が例外生成時の可変長ローカル配列を読み直していたため、関数内の要素変更・追加・`unset`が`getTrace()[0]['args']`に混入した。仮引数の束縛時に型変換後の可変長引数配列をフレームへ保持し、トレース生成ではその配列を展開する。通常の仮引数の既定値処理は維持した。新規fixtureは`exception_trace_variadic_changed.php`、`exception_trace_variadic_appended.php`、`exception_trace_variadic_unset.php`の3件。Go単体テストでは3設定すべてで渡した`10, 20`を保持することを確認した。

この段階では**662件**（正常350・実行時237・構文75）で、641件は[既存の公式PHP3版成功ログ](validation/include-exception-641-compare-all-mac.log)で全一致、後続21件の期待値は暫定だった。[`make test`](validation/variadic-trace-test.log)（Python12件を含む）、[`make vet`](validation/variadic-trace-vet.log)、[`make build`](validation/variadic-trace-build.log)、[`git diff --check`](validation/variadic-trace-diff-check.log)は各終了0。[対象比較試行](validation/variadic-trace-targeted-attempt.log)と[全件比較試行](validation/variadic-trace-compare-all-attempt.log)は各終了2。Podmanソケット接続拒否により、各版の版確認は終了125、各版0一致・0不一致・1環境エラー、662件未実施だった。PHPコードの終了125ではない。後続の662件実比較結果は上の節を参照。

Macでの新規3件と全件の再比較コマンドは次の通り。既存コンテナが停止している場合だけ起動する。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/variadic-trace-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/variadic-trace-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
for file in tests/fixtures/success/exception_trace_variadic_changed.php tests/fixtures/success/exception_trace_variadic_appended.php tests/fixtures/success/exception_trace_variadic_unset.php; do
  ./containers/php-version/compare-all.sh "$file" ./bin/gophp 2>&1 | tee -a docs/validation/variadic-trace-targeted-mac.log
  target_exit=$?
  echo "$file exit=$target_exit" | tee -a docs/validation/variadic-trace-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/variadic-trace-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/variadic-trace-compare-all-mac.log
```

## 例外トレースの引数と親アクセサのレビュー修正（2026-10-14）

`callFrame`に呼び出し時の引数と仮引数への束縛結果を保持し、例外生成時に仮引数順のトレースを作る。名前付き指定で飛ばされた既定値を含め、可変長引数はローカル配列を一つの引数として記録せず各要素へ展開する。文字列名で受けた可変長引数は`getTrace()`の`args`でもキーを保持する。`parent::getMessage()`など、派生例外からの親アクセサ呼び出しを既存の組み込み取得処理へ接続した。新規fixtureは`exception_trace_named_default.php`、`exception_trace_variadic.php`、`exception_parent_accessors.php`の3件。当時の期待値は暫定だったが、後続の662件比較ではこの3件が3版とも一致した。

[`make test`](validation/throwable-trace-review-make-test.log)、[`make vet`](validation/throwable-trace-review-make-vet.log)、[`make build`](validation/throwable-trace-review-make-build.log)、[`git diff --check`](validation/throwable-trace-review-diff-check.log)は各終了0。Goの3設定では新規3件が正常終了した。[対象比較試行](validation/throwable-trace-review-targeted-attempt.log)と[全件比較試行](validation/throwable-trace-review-compare-all-attempt.log)は各終了2で、3版とも`podman exec`終了125の接続拒否。各版0一致・0不一致・1環境エラー、659件未実施である。これはPHPプログラムの終了コード125ではない。

Macでの今回の対象比較と全件比較は次を実行する。既存コンテナが停止している場合だけ起動し、`tee`直後に終了コードを保存する。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/throwable-trace-review-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/throwable-trace-review-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
for file in tests/fixtures/success/exception_trace_named_default.php tests/fixtures/success/exception_trace_variadic.php tests/fixtures/success/exception_parent_accessors.php; do
  ./containers/php-version/compare-all.sh "$file" ./bin/gophp 2>&1 | tee -a docs/validation/throwable-trace-review-targeted-mac.log
  target_exit=$?
  echo "$file exit=$target_exit" | tee -a docs/validation/throwable-trace-review-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/throwable-trace-review-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/throwable-trace-review-compare-all-mac.log
```

## 641件成功前の比較・修正（履歴）

ユーザーがMacで実行した[修正前の全件ログ](validation/include-integration-641-user-before-fix.log)を確認し、元ログとバイト一致で保存した。PHP 8.3.35・8.4.26は各**640一致・1不一致・0環境エラー**、8.5.11は**638一致・3不一致・0環境エラー**、`make compare-all`終了2。不一致は`include_integration_strict_restore.php`の弱い`int`引数変換で非推奨警告が欠けたこと、8.5の`include_integration_declaration.php`と`include_integration_fatal_finally.php`でコンパイル時宣言エラートレースに`require()`が余計に入ったことによる。

引数変換成功時の小数精度喪失警告を呼び出し先の宣言行に出し、8.5の読み込み先宣言エラーでは`require()`フレームだけを除外した。実行時エラーや明示的に投げた例外のトレースは従来どおりにした。警告付きの正常終了fixtureは、ソースを変えず`success`から`runtime`へ移し、公式PHPログを根拠に期待stderrを訂正した。現在は**641件**（正常331・実行時236・構文74）。[保存済み公式stderrとの再照合](validation/include-integration-641-saved-output-recheck.log)は対象5比較で一致（終了0）したが、公式PHPの再実行ではない。[`make test`](validation/include-integration-641-review-final-make-test.log)、[`make vet`](validation/include-integration-641-review-make-vet.log)、[`make build`](validation/include-integration-641-review-make-build.log)は各終了0。

[対象比較試行](validation/include-integration-641-review-targeted-attempt.log)と[全件比較試行](validation/include-integration-641-review-compare-all-attempt.log)は、各版の版確認がPodmanソケット接続拒否（`podman exec`終了125）で止まり、各比較コマンドは終了2。修正後は**各版0一致・0不一致・1環境エラー、641件未実施**であり、修正前の実比較件数を修正後の成功件数として流用しない。Macでの対象比較と全件比較は次を実行する。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/include-integration-641-review-make-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/include-integration-641-review-make-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
for file in tests/fixtures/runtime/include_integration_strict_restore.php tests/fixtures/runtime/include_integration_declaration.php tests/fixtures/runtime/include_integration_fatal_finally.php; do
  ./containers/php-version/compare-all.sh "$file" ./bin/gophp 2>&1 | tee -a docs/validation/include-integration-641-review-targeted-mac.log
  target_exit=$?
  echo "$file exit=$target_exit" | tee -a docs/validation/include-integration-641-review-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/include-integration-641-review-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/include-integration-641-review-compare-all-mac.log
```

## 641件の初回比較前（履歴）

前回の`_once`履歴修正に続き、通常の`include`／`require`では解析・宣言検査が成功した後で履歴に登録するよう変更した。`include_once`／`require_once`は従来どおり、ファイルを開いた後、解析前に登録する。通常形式の解析失敗後に`*_once`を使うと再解析して`ParseError`となり、`*_once`の解析失敗後は同じ実パスへの`*_once`が`1`を返す。今回追加した3 fixtureとGo単体テストで両経路を確認した。

fixtureは**641件**（正常332・実行時235・構文74）。追加前の622件だけが[ユーザーの公式PHP3版成功ログ](validation/named-exception-622-compare-all-mac.log)で全一致している。新規19件の期待値はGoローカル検証による暫定値。[`make test`](validation/include-history-form-final-make-test.log)、[`make vet`](validation/include-history-form-final-make-vet.log)、[`make build`](validation/include-history-form-final-make-build.log)は各終了0。今回の[Podman確認](validation/include-history-form-podman-check.log)は終了125。[対象比較](validation/include-history-form-targeted-attempt.log)と[全件比較](validation/include-history-form-compare-all-attempt.log)の試行は、各版の版確認がPodman接続拒否（`podman exec`終了125）で止まり、比較コマンドは各終了2。**各版0一致・0不一致・1環境エラー、641件未実施**である。以下の638件時点の説明とログは履歴として読む。

Macで追加3件と全件を実比較する手順は次の通り。既存コンテナは停止中の場合だけ起動し、`tee`直後に終了コードを取得する。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/include-history-form-make-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/include-history-form-make-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
for file in tests/fixtures/success/include_integration_parse_normal_include.php tests/fixtures/success/include_integration_parse_normal_require.php tests/fixtures/success/include_integration_parse_once_history.php; do
  ./containers/php-version/compare-all.sh "$file" ./bin/gophp 2>&1 | tee -a docs/validation/include-history-form-targeted-mac.log
  target_exit=$?
  echo "$file exit=$target_exit" | tee -a docs/validation/include-history-form-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/include-history-form-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/include-history-form-compare-all-mac.log
```

## 638件時点のレビュー修正と比較試行（履歴）

項目6・10の読み込み・例外統合について、既存のローカルinclude/requireと例外伝播に、読み込み先の捕捉可能な`ParseError`、ファイル間トレース、失敗時の`_once`履歴を追加した。レビュー後、ファイルを開いた直後、コンパイル前に実パスを履歴へ登録するよう修正した。構文エラーを捕捉しても、2度目の同一ファイルへの`include_once`は`1`を返す。ファイル未発見は履歴に入れない。捕捉できない`fatalDiagnostic`は`finally`の`return`や`throw`で打ち消さず、診断と終了255を維持する。読み込み先での通常のthrow、`return`、関数スコープ・global/static束縛、ファイル文脈、strict_typesの既存動作は維持する。これらは[追加16件の一覧](validation/include-exception-fixtures.txt)とGo単体テストでローカル確認した。`include_path`全般、ストリームラッパーは未対応。

fixtureは**638件**（正常329・実行時235・構文74）。追加前の622件はユーザーがMacで実行した[公式PHP3版の成功ログ](validation/named-exception-622-compare-all-mac.log)で全一致。新規16件の期待値はGo側で確認した**暫定値**である。レビュー修正後の[`make test`](validation/include-exception-review-make-test.log)、[`make vet`](validation/include-exception-review-make-vet.log)、[`make build`](validation/include-exception-review-make-build.log)は終了0。今回の[対象比較試行](validation/include-exception-review-targeted-attempt.log)と[全件比較試行](validation/include-exception-review-compare-all-attempt.log)は、PHP 8.3.35・8.4.26・8.5.11の版確認時にPodmanソケット接続拒否（`podman exec`終了125）で停止し、各比較コマンドは終了2。**各版0一致・0不一致・1環境エラー、638件未実施**であり、PHPプログラムの終了125やfixtureの差異ではない。修正前の[対象](validation/include-exception-targeted-attempt.log)・[全件](validation/include-exception-compare-all-attempt.log)試行とその当時の暫定期待値は履歴として維持する。

Macでは既存コンテナが停止しているときだけ起動し、次を実行する。対象比較は補助ファイルを単独fixtureに含めない。各パイプの直後に終了コードを保存する。

```sh
cd /Users/kouji/dev/php
for name in php-version-83 php-version-84 php-version-85; do
  running=$(podman inspect --format '{{.State.Running}}' "$name") || exit 1
  if [ "$running" != true ]; then podman start "$name" || exit 1; fi
done
set -o pipefail
make build 2>&1 | tee docs/validation/include-exception-review-make-build-mac.log
build_exit=$?
echo "make build exit=$build_exit" | tee -a docs/validation/include-exception-review-make-build-mac.log
if [ "$build_exit" -ne 0 ]; then exit "$build_exit"; fi
for file in tests/fixtures/success/include_integration_*.php tests/fixtures/runtime/include_integration_*.php; do
  ./containers/php-version/compare-all.sh "$file" ./bin/gophp 2>&1 | tee -a docs/validation/include-exception-review-targeted-mac.log
  target_exit=$?
  echo "$file exit=$target_exit" | tee -a docs/validation/include-exception-review-targeted-mac.log
done
make compare-all 2>&1 | tee docs/validation/include-exception-review-compare-all-mac.log
compare_exit=$?
echo "make compare-all exit=$compare_exit" | tee -a docs/validation/include-exception-review-compare-all-mac.log
```

項目6・10および項目6〜10全体は部分実装のまま。公式比較後も未収録の診断・制御境界、残るThrowable API・エラーハンドラ、include_pathとラッパーが残る。以下は追加前に確定した622件の記録である。

## 622件が公式PHPの3版すべてで一致した記録（履歴）

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
